package app

type tokenDTO struct {
	Token string `json:"token"`
	Exp int64 `json:"expire"`
	Login string `json:"login"`
}
