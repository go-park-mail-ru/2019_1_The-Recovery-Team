package models

type Profile struct {
	ID       uint   `json:"id"`
	Email    string `json:"email,omitempty"`
	Nickname string `json:"nickname"`
	Password string `json:"password,omitempty"`
	Score
}

type Score struct {
	Record int `json:"record"`
	Win    int `json:"win"`
	Loss   int `json:"loss"`
}
