package model

type Vote struct {
	UserName      string `json:"user_name"`
	CandidateName string `json:"candidate_name"`
}
