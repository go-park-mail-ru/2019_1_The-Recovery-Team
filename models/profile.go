package models

//easyjson:json
type Profile struct {
	ID     uint64 `json:"id,omitempty" example:"1"`
	Avatar string `json:"avatar,omitempty example:"upload/img/1.png"`
	ProfileInfo
	Score
}

//easyjson:json
type Score struct {
	Record int `json:"record" example:"1500"`
	Win    int `json:"win" example:"100"`
	Loss   int `json:"loss" example:"50"`
}

//easyjson:json
type ProfileInfo struct {
	Email    string `json:"email,omitempty" example:"test@mail.ru"`
	Nickname string `json:"nickname,omitempty" example:"test"`
	Password string `json:"password,omitempty" example:"password"`
}

//easyjson:json
type ProfileRegistration struct {
	Nickname string `json:"nickname" example:"test"`
	ProfileLogin
}

//easyjson:json
type ProfileLogin struct {
	Email    string `json:"email" example:"test@mail.ru`
	Password string `json:"password" example:"password"`
}

//easyjson:json
type Profiles []Profile
