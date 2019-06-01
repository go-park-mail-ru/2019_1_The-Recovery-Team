package profile

//easyjson:json
type Profile struct {
	Info
	Email    string `json:"email,omitempty" example:"test@mail.ru" valid:"required~Email is required,email~Incorrect email"`
	Password string `json:"-,omitempty" example:"password" valid:"required~Password is required,stringlength(4|32)~Incorrect password length(4-32)"`
}

//easyjson:json
type Info struct {
	ID       uint64 `json:"id,omitempty" example:"1"`
	Nickname string `json:"nickname,omitempty" example:"test" valid:"required~Nickname is required,stringlength(4|20)~Incorrect nickname length(4-20)"`
	Avatar   string `json:"avatar,omitempty" example:"upload/img/1.png"`
	Oauth    string `json:"oauth,omitempty"`
	OauthId  string `json:"oauthId,omitempty"`
	Score
}

//easyjson:json
type Score struct {
	Position uint64 `json:"position" example:"1"`
	Record   int64  `json:"record" example:"1500"`
	Win      int64  `json:"win" example:"100"`
	Loss     int64  `json:"loss" example:"50"`
}

//easyjson:json
type Profiles struct {
	List  []Info
	Total int64 `json:"total" example:"50"`
}

//easyjson:json
type ID struct {
	Id uint64 `json:"id,omitempty" example:"1"`
}

//easyjson:json
type Avatar struct {
	Path string `json:"avatar,omitempty" example:"upload/img/1.png"`
}

//easyjson:json
type Registration struct {
	Nickname string `json:"nickname" example:"test" valid:"required~Nickname is required,stringlength(4|20)~Incorrect nickname length(4-20)"`
	Login
}

//easyjson:json
type Login struct {
	Email    string `json:"email" example:"test@mail.ru" valid:"required~Email is required,email~Incorrect email"`
	Password string `json:"password" example:"password" valid:"required~Password is required,stringlength(4|32)~Incorrect password length(4-32)"`
}

//easyjson:json
type Create struct {
	Email    string `json:"email,omitempty" example:"test@mail.ru" valid:"required~Email is required,email~Incorrect email"`
	Nickname string `json:"nickname,omitempty" example:"test" valid:"required~Nickname is required,stringlength(4|20)~Incorrect nickname length(4-20)"`
	Password string `json:"password,omitempty" example:"password" valid:"required~Password is required,stringlength(4|32)~Incorrect password length(4-32)"`
}

type CreateOauth struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
	Avatar Avatar `json:"avatar"`
	Oauth  string `json:"oauth"`
}

//easyjson:json
type Created struct {
	ID       uint64 `json:"id,omitempty" example:"1"`
	Email    string `json:"email,omitempty" example:"test@mail.ru" valid:"required~Email is required,email~Incorrect email"`
	Nickname string `json:"nickname,omitempty" example:"test" valid:"required~Nickname is required,stringlength(4|20)~Incorrect nickname length(4-20)"`
	Avatar   string `json:"avatar,omitempty" example:"upload/img/1.png"`
	Score
}

//easyjson:json
type UpdateInfo struct {
	Email    string `json:"email,omitempty" example:"test@mail.ru" valid:"required~Email is required,email~Incorrect email"`
	Nickname string `json:"nickname,omitempty" example:"test" valid:"required~Nickname is required,stringlength(4|20)~Incorrect nickname length(4-20)"`
}

//easyjson:json
type UpdatePassword struct {
	Password    string `json:"password,omitempty" example:"password" valid:"required~Password is required,stringlength(4|32)~Incorrect password length(4-32)"`
	PasswordOld string `json:"password_old" example:"password_old" valid:"required~Old password is required,stringlength(4|32)~Incorrect password length(4-32)"`
}

//easyjson:json
type OauthAccessTokenRaw struct {
	Token  string `json:"access_token"`
	UserId uint64 `json:"user_id"`
}

//easyjson:json
type OauthAccessToken struct {
	Token  string `json:"access_token"`
	UserId string `json:"user_id"`
}

//easyjson:json
type OauthResponse struct {
	Response []OauthAvatar `json:"response"`
}

//easyjson:json
type OauthAvatar struct {
	Avatar string `json:"photo_100"`
}
