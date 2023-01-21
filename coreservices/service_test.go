package coreservices

import (
	"context"
	"fmt"
	"testing"

	"github.com/pedrocmart/leaderboard-service/mocks"
	"github.com/pedrocmart/leaderboard-service/models"
	"github.com/stretchr/testify/assert"
)

func TestNewCoreService(t *testing.T) {
	cases := []struct {
		description  string
		panics       bool
		version      string
		expectedCore *models.Core
		core         *models.Core
	}{
		{
			description:  "should validate the version environment variable. Should not panic since it is not empty",
			core:         &models.Core{},
			expectedCore: &models.Core{},
		},
	}
	for _, tc := range cases {
		NewCoreService(tc.core)
		tc.expectedCore.Service = &BasicService{Core: tc.expectedCore}
		assert.Equal(t, tc.expectedCore, tc.core, tc.description)
	}
}

func TestBasicService_HandleSubmitScore(t *testing.T) {
	cases := []struct {
		description                  string
		basicAPIService              BasicService
		ctx                          context.Context
		request                      *models.SubmitScoreRequest
		userIdRequest                string
		expectedResponse             *models.SubmitScoreResponse
		expectedError                error
		createUserFuncError          error
		doesUserExistFuncError       error
		doesUserExistFunc            bool
		updateAbsoluteUserScoreError error
		updateRelativeUserScoreError error
		getUserByIdError             error
		getUserById                  *models.User
	}{
		{
			description: "should insert and return user with absolute score",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:           context.Background(),
			userIdRequest: "1",
			request: &models.SubmitScoreRequest{
				Total: &[]int{320}[0],
			},
			doesUserExistFunc:   false,
			createUserFuncError: nil,
			expectedResponse: &models.SubmitScoreResponse{
				UserID: 1,
				Score:  320,
			},
		},
		{
			description: "should return error when converting user_id to int",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:           context.Background(),
			userIdRequest: "abc",
			request: &models.SubmitScoreRequest{
				Total: &[]int{320}[0],
			},
			doesUserExistFunc:   false,
			createUserFuncError: nil,
			expectedResponse:    nil,
			expectedError:       fmt.Errorf("User_Id must be an integer."),
		},
		{
			description: "should return error when sending score and total at the same time",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:           context.Background(),
			userIdRequest: "1",
			request: &models.SubmitScoreRequest{
				Total: &[]int{320}[0],
				Score: "-100",
			},
			doesUserExistFunc:   false,
			createUserFuncError: nil,
			expectedResponse:    nil,
			expectedError:       fmt.Errorf("You can only submit the absolute score or the relative score."),
		},
		{
			description: "should return error when checks DoesUserExist ",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:           context.Background(),
			userIdRequest: "1",
			request: &models.SubmitScoreRequest{
				Score: "-100",
			},
			doesUserExistFunc:      false,
			createUserFuncError:    nil,
			expectedResponse:       nil,
			expectedError:          fmt.Errorf("mock-error"),
			doesUserExistFuncError: fmt.Errorf("mock-error"),
		},
		{
			description: "should return error when CreateUser",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:           context.Background(),
			userIdRequest: "1",
			request: &models.SubmitScoreRequest{
				Score: "-100",
			},
			doesUserExistFunc:   false,
			expectedResponse:    nil,
			expectedError:       fmt.Errorf("mock-error"),
			createUserFuncError: fmt.Errorf("mock-error"),
		},
		{
			description: "should return error when UpdateAbsoluteUserScore",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:           context.Background(),
			userIdRequest: "1",
			request: &models.SubmitScoreRequest{
				Total: &[]int{320}[0],
			},
			doesUserExistFunc:            true,
			expectedResponse:             nil,
			expectedError:                fmt.Errorf("mock-error"),
			updateAbsoluteUserScoreError: fmt.Errorf("mock-error"),
		},
		{
			description: "should return error when UpdateRelativeUserScore",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:           context.Background(),
			userIdRequest: "1",
			request: &models.SubmitScoreRequest{
				Score: "-100",
			},
			doesUserExistFunc:            true,
			expectedResponse:             nil,
			expectedError:                fmt.Errorf("mock-error"),
			updateRelativeUserScoreError: fmt.Errorf("mock-error"),
		},
		{
			description: "should return error when GetUserById",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:           context.Background(),
			userIdRequest: "1",
			request: &models.SubmitScoreRequest{
				Score: "-100",
			},
			doesUserExistFunc: true,
			expectedResponse:  nil,
			expectedError:     fmt.Errorf("mock-error"),
			getUserByIdError:  fmt.Errorf("mock-error"),
		},
		{
			description: "should update and return user with absolute score",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:           context.Background(),
			userIdRequest: "1",
			request: &models.SubmitScoreRequest{
				Total: &[]int{320}[0],
			},
			doesUserExistFunc: true,
			expectedResponse: &models.SubmitScoreResponse{
				UserID: 1,
				Score:  320,
			},
		},
		{
			description: "should update and return user with relative score",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:           context.Background(),
			userIdRequest: "1",
			request: &models.SubmitScoreRequest{
				Score: "-100",
			},
			doesUserExistFunc: true,
			getUserById: &models.User{
				UserID: 1,
				Score:  -100,
			},
			expectedResponse: &models.SubmitScoreResponse{
				UserID: 1,
				Score:  -100,
			},
		},
		{
			description: "should return error when converting score regex",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:           context.Background(),
			userIdRequest: "1",
			request: &models.SubmitScoreRequest{
				Score: "-100a",
			},
			doesUserExistFunc: true,
			expectedResponse:  nil,
			expectedError:     fmt.Errorf("Wrong format for the relative score. It must start with a [+] or [-] symbol."),
			getUserByIdError:  fmt.Errorf("Wrong format for the relative score. It must start with a [+] or [-] symbol."),
		},
	}
	for _, tc := range cases {
		mockedStoreService := mocks.StoreServiceMock{
			CreateUserFunc: func(ctx context.Context, id int, total int) error {
				return tc.createUserFuncError
			},
			DoesUserExistFunc: func(ctx context.Context, id int) (bool, error) {
				return tc.doesUserExistFunc, tc.doesUserExistFuncError
			},
			GetUserByIdFunc: func(ctx context.Context, id int) (*models.User, error) {
				return tc.getUserById, tc.getUserByIdError
			},
			UpdateAbsoluteUserScoreFunc: func(ctx context.Context, id int, score int) error {
				return tc.updateAbsoluteUserScoreError
			},
			UpdateRelativeUserScoreFunc: func(ctx context.Context, id int, score int) error {
				return tc.updateRelativeUserScoreError
			},
		}

		tc.basicAPIService.Core.StoreService = &mockedStoreService

		res, err := tc.basicAPIService.HandleSubmitScore(tc.ctx, tc.request, tc.userIdRequest)
		assert.Equal(t, tc.expectedResponse, res, tc.description)
		assert.Equal(t, tc.expectedError, err, tc.description)
	}
}

func TestBasicService_HandleGetRanking(t *testing.T) {
	cases := []struct {
		description          string
		basicAPIService      BasicService
		ctx                  context.Context
		request              string
		expectedResponse     *models.GetRankingResponse
		expectedError        error
		getUsersBetweenError error
		getUsersBetween      []models.Ranking
		getUsersError        error
		getUsers             []models.Ranking
	}{
		{
			description: "should return ranking using type Top",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:     context.Background(),
			request: "top100",
			getUsers: []models.Ranking{
				{
					Position: 1,
					UserID:   1,
					Score:    100,
				},
				{
					Position: 2,
					UserID:   200,
					Score:    50,
				},
			},
			expectedResponse: &models.GetRankingResponse{
				Ranking: []models.Ranking{
					{
						Position: 1,
						UserID:   1,
						Score:    100,
					},
					{
						Position: 2,
						UserID:   200,
						Score:    50,
					},
				},
			},
		},
		{
			description: "should return ranking using type At",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:     context.Background(),
			request: "At100/1",
			getUsersBetween: []models.Ranking{
				{
					Position: 1,
					UserID:   1,
					Score:    100,
				},
				{
					Position: 2,
					UserID:   200,
					Score:    50,
				},
			},
			expectedResponse: &models.GetRankingResponse{
				Ranking: []models.Ranking{
					{
						Position: 1,
						UserID:   1,
						Score:    100,
					},
					{
						Position: 2,
						UserID:   200,
						Score:    50,
					},
				},
			},
		},
		{
			description: "should return error using type top position 0",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:              context.Background(),
			request:          "top0",
			expectedResponse: nil,
			expectedError:    fmt.Errorf("The position must be greater than 0."),
		},
		{
			description: "should return error using invalid type",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:              context.Background(),
			request:          "invalid0",
			expectedResponse: nil,
			expectedError:    fmt.Errorf("The only formats accepted for type are: Top100 and At100/3."),
		},
		{
			description: "should return error when GetUsers",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:              context.Background(),
			request:          "top10",
			expectedResponse: nil,
			getUsersError:    fmt.Errorf("mock-error"),
			expectedError:    fmt.Errorf("mock-error"),
		},
		{
			description: "should return error when GetUsersBetween",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:                  context.Background(),
			request:              "At100/1",
			expectedResponse:     nil,
			getUsersBetweenError: fmt.Errorf("mock-error"),
			expectedError:        fmt.Errorf("mock-error"),
		},
		{
			description: "should return error using position 0",
			basicAPIService: BasicService{
				Core: &models.Core{},
			},
			ctx:                  context.Background(),
			request:              "At100/0",
			expectedResponse:     nil,
			getUsersBetweenError: fmt.Errorf("The positions must be greater than 0."),
			expectedError:        fmt.Errorf("The positions must be greater than 0."),
		},
		// {
		// 	description: "should return error when converting user_id to int",
		// 	basicAPIService: BasicService{
		// 		Core: &models.Core{},
		// 	},
		// 	ctx:           context.Background(),
		// 	userIdRequest: "abc",
		// 	request: &models.SubmitScoreRequest{
		// 		Total: &[]int{320}[0],
		// 	},
		// 	doesUserExistFunc:   false,
		// 	createUserFuncError: nil,
		// 	expectedResponse:    nil,
		// 	expectedError:       fmt.Errorf("User_Id must be an integer."),
		// },
		// {
		// 	description: "should return error when sending score and total at the same time",
		// 	basicAPIService: BasicService{
		// 		Core: &models.Core{},
		// 	},
		// 	ctx:           context.Background(),
		// 	userIdRequest: "1",
		// 	request: &models.SubmitScoreRequest{
		// 		Total: &[]int{320}[0],
		// 		Score: "-100",
		// 	},
		// 	doesUserExistFunc:   false,
		// 	createUserFuncError: nil,
		// 	expectedResponse:    nil,
		// 	expectedError:       fmt.Errorf("You can only submit the absolute score or the relative score."),
		// },
		// {
		// 	description: "should return error when checks DoesUserExist ",
		// 	basicAPIService: BasicService{
		// 		Core: &models.Core{},
		// 	},
		// 	ctx:           context.Background(),
		// 	userIdRequest: "1",
		// 	request: &models.SubmitScoreRequest{
		// 		Score: "-100",
		// 	},
		// 	doesUserExistFunc:      false,
		// 	createUserFuncError:    nil,
		// 	expectedResponse:       nil,
		// 	expectedError:          fmt.Errorf("mock-error"),
		// 	doesUserExistFuncError: fmt.Errorf("mock-error"),
		// },
		// {
		// 	description: "should return error when CreateUser",
		// 	basicAPIService: BasicService{
		// 		Core: &models.Core{},
		// 	},
		// 	ctx:           context.Background(),
		// 	userIdRequest: "1",
		// 	request: &models.SubmitScoreRequest{
		// 		Score: "-100",
		// 	},
		// 	doesUserExistFunc:   false,
		// 	expectedResponse:    nil,
		// 	expectedError:       fmt.Errorf("mock-error"),
		// 	createUserFuncError: fmt.Errorf("mock-error"),
		// },
		// {
		// 	description: "should return error when UpdateAbsoluteUserScore",
		// 	basicAPIService: BasicService{
		// 		Core: &models.Core{},
		// 	},
		// 	ctx:           context.Background(),
		// 	userIdRequest: "1",
		// 	request: &models.SubmitScoreRequest{
		// 		Total: &[]int{320}[0],
		// 	},
		// 	doesUserExistFunc:            true,
		// 	expectedResponse:             nil,
		// 	expectedError:                fmt.Errorf("mock-error"),
		// 	updateAbsoluteUserScoreError: fmt.Errorf("mock-error"),
		// },
		// {
		// 	description: "should return error when UpdateRelativeUserScore",
		// 	basicAPIService: BasicService{
		// 		Core: &models.Core{},
		// 	},
		// 	ctx:           context.Background(),
		// 	userIdRequest: "1",
		// 	request: &models.SubmitScoreRequest{
		// 		Score: "-100",
		// 	},
		// 	doesUserExistFunc:            true,
		// 	expectedResponse:             nil,
		// 	expectedError:                fmt.Errorf("mock-error"),
		// 	updateRelativeUserScoreError: fmt.Errorf("mock-error"),
		// },
		// {
		// 	description: "should return error when GetUserById",
		// 	basicAPIService: BasicService{
		// 		Core: &models.Core{},
		// 	},
		// 	ctx:           context.Background(),
		// 	userIdRequest: "1",
		// 	request: &models.SubmitScoreRequest{
		// 		Score: "-100",
		// 	},
		// 	doesUserExistFunc: true,
		// 	expectedResponse:  nil,
		// 	expectedError:     fmt.Errorf("mock-error"),
		// 	getUserByIdError:  fmt.Errorf("mock-error"),
		// },
		// {
		// 	description: "should update and return user with absolute score",
		// 	basicAPIService: BasicService{
		// 		Core: &models.Core{},
		// 	},
		// 	ctx:           context.Background(),
		// 	userIdRequest: "1",
		// 	request: &models.SubmitScoreRequest{
		// 		Total: &[]int{320}[0],
		// 	},
		// 	doesUserExistFunc: true,
		// 	expectedResponse: &models.SubmitScoreResponse{
		// 		UserID: 1,
		// 		Score:  320,
		// 	},
		// },
		// {
		// 	description: "should update and return user with relative score",
		// 	basicAPIService: BasicService{
		// 		Core: &models.Core{},
		// 	},
		// 	ctx:           context.Background(),
		// 	userIdRequest: "1",
		// 	request: &models.SubmitScoreRequest{
		// 		Score: "-100",
		// 	},
		// 	doesUserExistFunc: true,
		// 	getUserById: &models.User{
		// 		UserID: 1,
		// 		Score:  -100,
		// 	},
		// 	expectedResponse: &models.SubmitScoreResponse{
		// 		UserID: 1,
		// 		Score:  -100,
		// 	},
		// },
		// {
		// 	description: "should return error when converting score regex",
		// 	basicAPIService: BasicService{
		// 		Core: &models.Core{},
		// 	},
		// 	ctx:           context.Background(),
		// 	userIdRequest: "1",
		// 	request: &models.SubmitScoreRequest{
		// 		Score: "-100a",
		// 	},
		// 	doesUserExistFunc: true,
		// 	expectedResponse:  nil,
		// 	expectedError:     fmt.Errorf("Wrong format for the relative score. It must start with a [+] or [-] symbol."),
		// 	getUserByIdError:  fmt.Errorf("Wrong format for the relative score. It must start with a [+] or [-] symbol."),
		// },
	}
	for _, tc := range cases {
		mockedStoreService := mocks.StoreServiceMock{
			GetUsersFunc: func(ctx context.Context, top int) ([]models.Ranking, error) {
				return tc.getUsers, tc.getUsersError
			},
			GetUsersBetweenFunc: func(ctx context.Context, lower, upper int) ([]models.Ranking, error) {
				return tc.getUsersBetween, tc.getUsersBetweenError
			},
		}

		tc.basicAPIService.Core.StoreService = &mockedStoreService

		res, err := tc.basicAPIService.HandleGetRanking(tc.ctx, tc.request)
		assert.Equal(t, tc.expectedResponse, res, tc.description)
		assert.Equal(t, tc.expectedError, err, tc.description)
	}
}
