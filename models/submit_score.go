package models

type SubmitScoreRequest struct {
	UserID int    `json:"user,omitempty"`
	Total  *int   `json:"total,omitempty"`
	Score  string `json:"score,omitempty"`
}

type SubmitScoreResponse struct {
	UserID int `json:"user_id,omitempty"`
	Score  int `json:"score,omitempty"`
}
