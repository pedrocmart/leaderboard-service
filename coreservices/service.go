package coreservices

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/pedrocmart/leaderboard-service/models"
)

//NewCoreService - will return an implementation of the http interface which is a collection of http functions.
//It will also add it to the core
func NewCoreService(core *models.Core) models.Service {
	basicService := BasicService{
		Core: core,
	}
	core.Service = &basicService
	return &basicService
}

type BasicService struct {
	Core *models.Core
}

func (bhs *BasicService) HandleSubmitScore(ctx context.Context, request *models.SubmitScoreRequest, userId string) (*models.SubmitScoreResponse, error) {
	var err error
	var score, currentScore int
	var isAbsolute bool

	if request.Score != "" && request.Total != nil {
		err := fmt.Errorf("You can only submit the absolute score or the relative score.")
		return nil, err
	}

	request.UserID, err = strconv.Atoi(userId)
	if err != nil {
		err := fmt.Errorf("User_Id must be an integer.")
		return nil, err
	}

	exists, err := bhs.Core.StoreService.DoesUserExist(ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	if request.Total != nil {
		score = *request.Total
		isAbsolute = true
	} else {
		//makes sure the user only inputs + or - in the beginning, followed by a number
		//eg: -100 or +100
		re := regexp.MustCompile(`^[+|-](\d+)$`)
		subMatchAll := re.FindAllString(request.Score, -1)
		if len(subMatchAll) > 0 {
			score, err = strconv.Atoi(subMatchAll[0])
			if err != nil {
				return nil, err
			}
		} else {
			err = errors.New("Wrong format for the relative score. It must start with a [+] or [-] symbol.")
			return nil, err
		}
	}

	if !exists {
		err := bhs.Core.StoreService.CreateUser(ctx, request.UserID, score)
		if err != nil {
			return nil, err
		}
		currentScore = score
	} else {
		if isAbsolute {
			err := bhs.Core.StoreService.UpdateAbsoluteUserScore(ctx, request.UserID, score)
			if err != nil {
				return nil, err
			}
			currentScore = score
		} else {
			err = bhs.Core.StoreService.UpdateRelativeUserScore(ctx, request.UserID, score)
			if err != nil {
				return nil, err
			}

			user, err := bhs.Core.StoreService.GetUserById(ctx, request.UserID)
			if err != nil {
				return nil, err
			}
			currentScore = user.Score
		}
	}
	response := new(models.SubmitScoreResponse)
	response.UserID = request.UserID
	response.Score = currentScore

	return response, nil
}

func (bhs *BasicService) HandleGetRanking(ctx context.Context, rankingType string) (*models.GetRankingResponse, error) {
	var ranking []models.Ranking

	//regex to make sure the user inputs "top" and a number after it
	//eg: top100
	isTopType, err := regexp.MatchString(`(?i)^top(\d+)$`, rankingType)
	if err != nil {
		return nil, err
	}

	if isTopType {
		//gets the number from rankingType
		re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		subMatchAll := re.FindAllString(rankingType, -1)
		topPositions := 0
		if len(subMatchAll) > 0 {
			topPositions, err = strconv.Atoi(subMatchAll[0])
			if err != nil {
				return nil, err
			}
		}
		if topPositions <= 0 {
			return nil, errors.New("The position must be greater than 0.")
		}

		ranking, err = bhs.Core.StoreService.GetUsers(ctx, topPositions)
		if err != nil {
			return nil, err
		}

	} else {
		//regex to make sure the user inputs "at", a number after, followed by slash and another number
		//eg: at100/3
		isAtType, err := regexp.MatchString(`(?i)^at(\d+)/(\d+)$`, rankingType) //At100/3
		if err != nil {
			return nil, err
		}

		if !isAtType {
			return nil, errors.New("The only formats accepted for type are: Top100 and At100/3.")
		}

		//gets the two numbers from rankingType
		re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		subMatchAll := re.FindAllString(rankingType, -1)
		topPositions, around := 0, 0

		if len(subMatchAll) > 0 {
			topPositions, err = strconv.Atoi(subMatchAll[0])
			if err != nil {
				return nil, err
			}

			around, err = strconv.Atoi(subMatchAll[1])
			if err != nil {
				return nil, err
			}
		}

		if topPositions <= 0 || around <= 0 {
			return nil, errors.New("The positions must be greater than 0.")
		}

		//get lower and upper values
		ranking, err = bhs.Core.StoreService.GetUsersBetween(ctx, topPositions, around)
		if err != nil {
			return nil, err
		}
	}

	response := new(models.GetRankingResponse)
	response.Ranking = ranking

	return response, nil
}
