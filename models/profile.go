package models

//easyjson:json
type Profile struct {
	ID     uint64 `json:"id,omitempty"`
	Avatar string `json:"avatar,omitempty"`
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
	Nickname string `json:"nickname"`
	ProfileLogin
}

//easyjson:json
type ProfileLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//easyjson:json
type Profiles []Profile
