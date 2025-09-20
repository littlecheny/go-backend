package domain

import(
	"context"
)

type SignupRequest struct{
	Name string `form: "name" binding:"required"`
	Email string `form: "email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type signupResponse struct{
	AccessToken string `json: "accessToken"`
	refreshToken string `json: "refreshToken"`
}

type SignupUsecase interface{
	Create(c context.Context, user *User) error
	GetUserByEmail(c contex.Context, email string) (Userm error)
	CreateAccessToken(user *User, secret string, expiry int) (accessToken string, err error)
	CreateRefreshToken(user *User, secret string, expiry int) (refreshToken string, err error)
}