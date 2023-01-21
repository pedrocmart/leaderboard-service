package coreservices

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pedrocmart/leaderboard-service/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestBasicStoreService_CreateUser(t *testing.T) {
	cases := []struct {
		description   string
		core          *models.Core
		context       context.Context
		userId        int
		score         int
		expectedError error
		err           error
		query         string
	}{
		{
			description: "Should create a row in db",
			core:        &models.Core{},
			context:     context.Background(),
			userId:      1,
			score:       100,
			query:       `INSERT INTO users (id, score) VALUES ($1, $2)`,
		},
		{
			description:   "Should return an error",
			core:          &models.Core{},
			context:       context.Background(),
			userId:        1,
			score:         100,
			query:         `INSERT INTO users (id, score) VALUES ($1, $2)`,
			err:           fmt.Errorf("mock-error"),
			expectedError: errors.Wrapf(fmt.Errorf("mock-error"), "create user"),
		},
	}
	for _, tc := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		tc.core.DB = db
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(tc.query)).
			WithArgs(
				tc.userId,
				tc.score,
			).
			WillReturnResult(sqlmock.NewResult(1, 1)).
			WillReturnError(tc.err)

		basicStore := NewStoreService(tc.core, tc.core.DB)
		// now we execute our method
		err = basicStore.CreateUser(tc.context, tc.userId, tc.score)
		if tc.err != nil {
			assert.Error(t, err)
			return
		}

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestBasicStoreService_UpdateRelativeUserScore(t *testing.T) {
	cases := []struct {
		description   string
		core          *models.Core
		context       context.Context
		userId        int
		score         int
		expectedError error
		err           error
		query         string
	}{
		{
			description: "Should update user with a relative score",
			core:        &models.Core{},
			context:     context.Background(),
			userId:      1,
			score:       100,
			query: `UPDATE users 
			SET score = score + $1
			WHERE id = $2`,
		},
		{
			description: "Should return an error",
			core:        &models.Core{},
			context:     context.Background(),
			userId:      100,
			score:       1,
			query: `UPDATE users 
			SET score = score + $1
			WHERE id = $2`,
			err:           fmt.Errorf("mock-error"),
			expectedError: errors.Wrapf(fmt.Errorf("mock-error"), "create user"),
		},
	}
	for _, tc := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		tc.core.DB = db
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(tc.query)).
			WithArgs(
				tc.score,
				tc.userId,
			).
			WillReturnResult(sqlmock.NewResult(1, 1)).
			WillReturnError(tc.err)

		basicStore := NewStoreService(tc.core, tc.core.DB)
		// now we execute our method
		err = basicStore.UpdateRelativeUserScore(tc.context, tc.userId, tc.score)
		if tc.err != nil {
			assert.Error(t, err)
			return
		}

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestBasicStoreService_UpdateAbsoluteUserScore(t *testing.T) {
	cases := []struct {
		description   string
		core          *models.Core
		context       context.Context
		userId        int
		score         int
		expectedError error
		err           error
		query         string
	}{
		{
			description: "Should update user with an absolute score",
			core:        &models.Core{},
			context:     context.Background(),
			userId:      1,
			score:       100,
			query: `UPDATE users 
			SET score = $1
			WHERE id = $2`,
		},
		{
			description: "Should return an error",
			core:        &models.Core{},
			context:     context.Background(),
			userId:      100,
			score:       1,
			query: `UPDATE users 
			SET score = $1
			WHERE id = $2`,
			err:           fmt.Errorf("mock-error"),
			expectedError: errors.Wrapf(fmt.Errorf("mock-error"), "create user"),
		},
	}
	for _, tc := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		tc.core.DB = db
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(tc.query)).
			WithArgs(
				tc.score,
				tc.userId,
			).
			WillReturnResult(sqlmock.NewResult(1, 1)).
			WillReturnError(tc.err)

		basicStore := NewStoreService(tc.core, tc.core.DB)
		// now we execute our method
		err = basicStore.UpdateAbsoluteUserScore(tc.context, tc.userId, tc.score)
		if tc.err != nil {
			assert.Error(t, err)
			return
		}

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestBasicStoreService_DoesUserExist(t *testing.T) {
	cases := []struct {
		description   string
		core          *models.Core
		context       context.Context
		userId        int
		expectedError error
		err           error
		query         string
		rows          *sqlmock.Rows
	}{
		{
			description: "Should check if user exists",
			core:        &models.Core{},
			context:     context.Background(),
			userId:      1,
			query:       "SELECT id FROM users WHERE id = $1",
			rows: sqlmock.NewRows(([]string{
				"id",
			})).AddRow(1),
		},
		{
			description: "Should return an error",
			core:        &models.Core{},
			context:     context.Background(),
			userId:      1,
			rows: sqlmock.NewRows(([]string{
				"score",
			})).AddRow("a"),
			query:         "SELECT id FROM users WHERE id = $1",
			err:           fmt.Errorf("mock-error"),
			expectedError: errors.Wrapf(fmt.Errorf("mock-error"), "create user"),
		},
	}
	for _, tc := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		tc.core.DB = db

		mock.ExpectQuery(regexp.QuoteMeta(tc.query)).
			WithArgs(
				tc.userId,
			).
			WillReturnRows(tc.rows)

		basicStore := NewStoreService(tc.core, tc.core.DB)
		// now we execute our method
		_, err = basicStore.DoesUserExist(tc.context, tc.userId)
		if tc.err != nil {
			assert.Error(t, err)
			return
		}

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestBasicStoreService_GetUsers(t *testing.T) {
	cases := []struct {
		description    string
		core           *models.Core
		context        context.Context
		userId         int
		expectedError  error
		err            error
		query          string
		rows           *sqlmock.Rows
		expectedResult []models.Ranking
	}{
		{
			description: "Should get users",
			core:        &models.Core{},
			context:     context.Background(),
			userId:      1,
			query:       "SELECT id, score FROM users ORDER BY score DESC LIMIT $1",
			rows: sqlmock.NewRows(([]string{
				"id",
				"score",
			})).AddRow(1, 100),
			expectedResult: []models.Ranking{
				{
					Position: 1,
					UserID:   1,
					Score:    100,
				},
			},
		},
		{
			description: "Should return an error",
			core:        &models.Core{},
			context:     context.Background(),
			userId:      1,
			rows: sqlmock.NewRows(([]string{
				"score",
			})).AddRow("a"),
			query:         "SELECT id, score FROM users ORDER BY score DESC LIMIT $1",
			err:           fmt.Errorf("mock-error"),
			expectedError: errors.Wrapf(fmt.Errorf("mock-error"), "create user"),
		},
	}
	for _, tc := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		tc.core.DB = db

		mock.ExpectQuery(regexp.QuoteMeta(tc.query)).
			WithArgs(
				tc.userId,
			).
			WillReturnRows(tc.rows)

		basicStore := NewStoreService(tc.core, tc.core.DB)
		// now we execute our method
		result, err := basicStore.GetUsers(tc.context, tc.userId)
		if tc.err != nil {
			assert.Error(t, err)
			return
		}

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		assert.Equal(t, tc.expectedResult, result, tc.description)
		assert.Equal(t, tc.err, err, tc.description)
	}
}

func TestBasicStoreService_GetUsersBetween(t *testing.T) {
	cases := []struct {
		description    string
		core           *models.Core
		context        context.Context
		pos            int
		around         int
		expectedError  error
		err            error
		query          string
		rows           *sqlmock.Rows
		expectedResult []models.Ranking
	}{
		{
			description: "Should get users with offset and limit",
			core:        &models.Core{},
			context:     context.Background(),
			pos:         7,
			around:      3,
			query:       "SELECT id, score FROM users ORDER BY score DESC LIMIT $1 OFFSET $2",
			rows: sqlmock.NewRows(([]string{
				"id",
				"score",
			})).AddRow(1, 100),
			expectedResult: []models.Ranking{
				{
					Position: 4,
					UserID:   1,
					Score:    100,
				},
			},
		},
	}
	for _, tc := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		tc.core.DB = db

		mock.ExpectQuery(regexp.QuoteMeta(tc.query)).
			WithArgs(
				tc.pos,
				tc.around,
			).
			WillReturnRows(tc.rows)
		basicStore := NewStoreService(tc.core, tc.core.DB)
		// now we execute our method
		result, err := basicStore.GetUsersBetween(tc.context, tc.pos, tc.around)
		if tc.err != nil {
			assert.Error(t, err)
			return
		}

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		assert.Equal(t, tc.expectedResult, result, tc.description)
		assert.Equal(t, tc.err, err, tc.description)
	}
}

func TestBasicStoreService_GetUserById(t *testing.T) {
	cases := []struct {
		description    string
		core           *models.Core
		context        context.Context
		userId         int
		score          int
		expectedError  error
		err            error
		query          string
		rows           *sqlmock.Rows
		expectedResult *models.User
	}{
		{
			description: "Should get an user by id",
			core:        &models.Core{},
			context:     context.Background(),
			userId:      1,
			score:       100,
			query:       "SELECT id, score FROM users WHERE id = $1",
			rows: sqlmock.NewRows(([]string{
				"id",
				"score",
			})).AddRow(1, 100),
			expectedResult: &models.User{
				UserID: 1,
				Score:  100,
			},
		},
		{
			description: "Should return an error",
			core:        &models.Core{},
			context:     context.Background(),
			userId:      1,
			rows: sqlmock.NewRows(([]string{
				"score",
			})).AddRow("a"),
			query:         "SELECT id, score FROM users WHERE id = $1",
			err:           fmt.Errorf("mock-error"),
			expectedError: errors.Wrapf(fmt.Errorf("mock-error"), "create user"),
		},
	}
	for _, tc := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		tc.core.DB = db

		mock.ExpectQuery(regexp.QuoteMeta(tc.query)).
			WithArgs(
				tc.userId,
			).
			WillReturnRows(tc.rows)

		basicStore := NewStoreService(tc.core, tc.core.DB)
		// now we execute our method
		result, err := basicStore.GetUserById(tc.context, tc.userId)
		if tc.err != nil {
			assert.Error(t, err)
			return
		}

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		assert.Equal(t, tc.expectedResult, result, tc.description)
		assert.Equal(t, tc.err, err, tc.description)
	}
}
