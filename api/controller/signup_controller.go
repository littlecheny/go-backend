package controller

import (
	"net/http"

	"github.com/littlecheny/go-backend/bootstrap"
	"github.com/littlecheny/go-backend/domain"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type SignupController struct{
	SignupUsercase domain.SignupUsecase
	Env *bootstrap.Env
}

func (sc *SignupController) Signup(c *gin.Context) {
	var request domain.SignupRequest

	err := c.ShouldBlind(&request)
	if err != nil{
		c.JSON(http.StattusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	_, err = sc.SignupUsecase.GetUserByEmail(c, request.Email)
	if err == nil {
		c.JSON(http.StatusConflict, domain.ErrorResponse{Message:"User already exists with given email"})
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCose,
	)
	if err != nil{
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	request.Password = string(encryptedPassword)

	user := domain.User{
		ID: primitive.NewObjectionID(),
		Name: request.Name,
		Email: request.Email
		Password: request.Password
	}

	err = sc.SignupUsecase.Create(c, &user)

	if err != nil{
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message:err.Error()})
		return
	}

	accessToken, err := sc.SignupUsecase.CreateAccessToken(&user, sc.Env.AccessTokenSecret, sc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	refreshToken, err := sc.SignupUsecase.CreateRefreshToken(&user,sc.Env.RefreshTokenSecret, sc.Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	signupResponse := domain.signupResponse{
		AccessToken: accessToken,
		refreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, signupResponse)
}