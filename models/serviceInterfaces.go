package models

import (
	"context"
	"net/http"
)

type HandlersService interface {
}

//go:generate moq -out ../mocks/service.go -pkg mocks  . Service
type Service interface {
	HandleSubmitScore(context.Context, *SubmitScoreRequest, string) (*SubmitScoreResponse, error)
	HandleGetRanking(context.Context, string) (*GetRankingResponse, error)
}

//go:generate moq -out ../mocks/storeService.go -pkg mocks  . StoreService
type StoreService interface {
	CreateUser(ctx context.Context, id int, total int) error
	UpdateRelativeUserScore(ctx context.Context, id int, score int) error
	UpdateAbsoluteUserScore(ctx context.Context, id int, score int) error
	GetUsers(ctx context.Context, top int) ([]Ranking, error)
	GetUserById(ctx context.Context, id int) (*User, error)
	GetUsersBetween(ctx context.Context, lower, upper int) ([]Ranking, error)
	DoesUserExist(ctx context.Context, id int) (bool, error)
}

//go:generate moq -out ../mocks/requestResponse.go -pkg mocks  . RequestResponse
type RequestResponse interface {
	HandleError(err error, w http.ResponseWriter, r *http.Request, status int)
	HandleResponse(body interface{}, w http.ResponseWriter, r *http.Request, status int)
	ReadBodyAsJSON(req *http.Request, dest interface{}) (err error)
}
