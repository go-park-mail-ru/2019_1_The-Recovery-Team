package models

//easyjson:json
type Profile struct {
	ID uint64 `json:"id,omitempty"`
	ProfileInfo
	Score
}

//easyjson:json
type Score struct {
	Record int `json:"record,omitempty"`
	Win    int `json:"win,omitempty"`
	Loss   int `json:"loss,omitempty"`
}

//easyjson:json
type ProfileInfo struct {
	Email    string `json:"email,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Password string `json:"password,omitempty"`
}

//easyjson:json
type ProfileRegistration struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

//easyjson:json
type Profiles []Profile
