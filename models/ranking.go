package models

type GetRankingResponse struct {
	Ranking []Ranking `json:"ranking"`
}

type Ranking struct {
	Position int `json:"position"`
	UserID   int `json:"user_id"`
	Score    int `json:"score"`
}
