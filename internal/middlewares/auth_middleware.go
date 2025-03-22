package middlewares

import (
	"net/http"
	"os"
	"strings"
	"time"

	"example.com/go-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)
func CheckAuth(ctx *gin.Context){
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	authToken := strings.Split(authHeader," ")
	if len(authToken) != 2  || authToken[0] != "Bearer" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	tokenString := authToken[1]
	accToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_,ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("SECRET_KEY")),nil
	})
	if err != nil || !accToken.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims,ok := accToken.Claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64){
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userModel,err:= services.GetUserByEmailService(claims["email"].(string))
	if userModel.ID == 0 || err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.Set("user", userModel)
	ctx.Next()
}	