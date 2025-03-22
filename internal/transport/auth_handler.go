package transport

import (
	"net/http"
	"os"
	"time"

	"example.com/go-api/internal/dtos"
	"example.com/go-api/internal/repository"
	"example.com/go-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserRegisterHandler handles the HTTP POST request for user registration.
// @Summary Register a new user
// @Description Register a new user with a unique email and hashed password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dtos.AuthRegisterInput true "User registration details"
// @Success 201 {object} repository.User "User registered successfully"
// @Failure 400 {object} map[string]string "Invalid input or email already exists"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /auth/register [post]
func UserRegisterHandler(ctx *gin.Context){
	var userRegisterInput dtos.AuthRegisterInput
	err := ctx.ShouldBindBodyWithJSON(&userRegisterInput);
	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":err})
		return
	}

	email := userRegisterInput.Email
	user,_ := services.GetUserByEmailService(email)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError,gin.H{"error":err})
	// 	return
	// }
	if user.ID != 0{
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"Email already exists"})
		return
	}

	passwordHashed,err := bcrypt.GenerateFromPassword([]byte(userRegisterInput.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err})
		return
	}

	var userModel = &repository.User{
		UserName: userRegisterInput.UserName,
		Email: userRegisterInput.Email,
		Password: string(passwordHashed),
	}
	services.CreateUser(userModel);
	ctx.JSON(http.StatusCreated,&dtos.UserResponse{
		ID: userModel.ID,
		UserName: userModel.UserName,
		Email: userModel.Email,
	})
}

// Start of Selection

// UserLoginHandler handles the HTTP POST request for user login.
// @Summary Login a user
// @Description Authenticate a user with email and password, returning access and refresh tokens.
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dtos.AuthLoginInput true "User login details"
// @Success 200 {object} map[string]string "Tokens generated successfully"
// @Failure 400 {object} map[string]string "Invalid email or password"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /auth/login [post]
func UserLoginHandler(ctx *gin.Context){

	var userLoginInput dtos.AuthLoginInput
	err := ctx.ShouldBindBodyWithJSON(&userLoginInput);
	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":err})
		return
	}

	userModel,err := services.GetUserByEmailService(userLoginInput.Email)

	if userModel.ID == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"Invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userModel.Password),[]byte(userLoginInput.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"Invalid email or password"})
		return
	}

	generateAccToken:= jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"email": userModel.Email,
		"userName": userModel.UserName,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})

	generateRefToken:= jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"email": userModel.Email,
		"userName": userModel.UserName,
		"exp": time.Now().Add(time.Hour * 24*7).Unix(),
	})

	accessToken,err := generateAccToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err})
		return
	}
	refressToken,err := generateRefToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err})
		return
	}
	ctx.SetCookie("refreshToken",refressToken,60*60*24*7,"/","localhost",false,true)
	ctx.JSON(http.StatusOK,gin.H{"accessToken":accessToken})
}

// UserRefreshTokenHandler handles the HTTP POST request to refresh the user's access token.
// @Summary Refresh user access token
// @Description Refresh the access token using a valid refresh token stored in cookies.
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string "New access token generated successfully"
// @Failure 400 {object} map[string]string "Invalid or expired refresh token"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /auth/refresh [get]
func UserRefreshTokenHandler(ctx *gin.Context){
	refreshTokenString,err := ctx.Cookie("refreshToken")
	if refreshTokenString == "" || err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"Invalid refresh token"})
		return
	}
	// verify the refresh token
	refreshToken, err := jwt.Parse(refreshTokenString, func(t *jwt.Token) (interface{}, error) {
		_,ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("SECRET_KEY")),nil
	})
	if err != nil || !refreshToken.Valid{
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"Invalid refresh token"})
		return
	}
	claims,ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"Invalid refresh token"})
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64){
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"Refresh token expired"})
		return
	}	

	userModel,err := services.GetUserByEmailService(claims["email"].(string))
	if userModel.ID == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"User not found"})
		return
	}
	// generate new access token and refresh token. Then set the new refresh token to the cookie and return the new access token
	generateAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"email": userModel.Email,
		"userName": userModel.UserName,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})
	generateRefToken:= jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"email": userModel.Email,
		"userName": userModel.UserName,
		"exp": time.Now().Add(time.Hour * 24*7).Unix(),
	})
	accessToken,err := generateAccessToken.SignedString([]byte(os.Getenv("SECRET_KEY")));
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err})
		return
	}
	refToken,err := generateRefToken.SignedString([]byte(os.Getenv("SECRET_KEY")));
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err})
		return
	}

	ctx.SetCookie("refreshToken",refToken,60*60*24*7,"/","localhost",false,true)
	ctx.JSON(http.StatusOK, gin.H{"accessToken":accessToken})
}

// UserLogoutHandler handles the HTTP POST request to log out the user.
// @Summary Log out user
// @Description Log out the authenticated user by clearing the refresh token cookie.
// @Tags Auth
// @Success 200 {object} map[string]string "Logout successfully"
// @Router /auth/logout [get]
func UserLogoutHandler(ctx *gin.Context){
	ctx.SetCookie("refreshToken","",-1,"/","localhost",false,true)
	ctx.JSON(http.StatusOK,gin.H{"message":"Logout successfully"})
}