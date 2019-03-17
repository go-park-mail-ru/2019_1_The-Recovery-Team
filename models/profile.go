package models

//easyjson:json
type Profile struct {
	ProfileInfo
	Email    string `json:"email,omitempty" example:"test@mail.ru"`
	Password string `json:"-,omitempty" example:"password"`
}

//easyjson:json
type ProfileInfo struct {
	ProfileID
	Nickname string `json:"nickname,omitempty" example:"test"`
	ProfileAvatar
	Score
}

//easyjson:json
type ProfileID struct {
	ID uint64 `json:"id,omitempty" example:"1"`
}

//easyjson:json
type ProfileAvatar struct {
	Avatar string `json:"avatar,omitempty" example:"upload/img/1.png"`
}

//easyjson:json
type Score struct {
	Record int `json:"record" example:"1500"`
	Win    int `json:"win" example:"100"`
	Loss   int `json:"loss" example:"50"`
}

//easyjson:json
type Profiles struct {
	List  []ProfileInfo
	Total int64 `json:"total" example:"50"`
}

//easyjson:json
type ProfileRegistration struct {
	Nickname string `json:"nickname" example:"test"`
	ProfileLogin
}

//easyjson:json
type ProfileLogin struct {
	Email    string `json:"email" example:"test@mail.ru"`
	Password string `json:"password" example:"password"`
}

//easyjson:json
type ProfileCreate struct {
	Email    string `json:"email,omitempty" example:"test@mail.ru"`
	Nickname string `json:"nickname,omitempty" example:"test"`
	Password string `json:"password,omitempty" example:"password"`
}

//easyjson:json
type ProfileCreated struct {
	ProfileID
	Email    string `json:"email,omitempty" example:"test@mail.ru"`
	Nickname string `json:"nickname,omitempty" example:"test"`
	Avatar   string `json:"avatar,omitempty" example:"upload/img/1.png"`
}

//easyjson:json
type ProfileUpdate struct {
	Email    string `json:"email,omitempty" example:"test@mail.ru"`
	Nickname string `json:"nickname,omitempty" example:"test"`
}

//easyjson:json
type ProfileUpdatePassword struct {
	Password    string `json:"password,omitempty" example:"password"`
	PasswordOld string `json:"password_old" example:"password_old"`
}
