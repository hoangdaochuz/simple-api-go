package transport

import (
	"net/http"

	"example.com/go-api/internal/dtos"
	"example.com/go-api/internal/repository"
	"github.com/gin-gonic/gin"
)

// GetUserMyProfileHandler handles the HTTP GET request to retrieve the profile of the authenticated user.
// @Summary Get user profile
// @Description Retrieve the profile information of the authenticated user.
// @Tags User
// @Produce json
// @Success 200 {object} dtos.UserResponse "User profile retrieved successfully"
// @Failure 500 {object} map[string]string "User not found"
// @Router /users/my-profile [get]
func GetUserMyProfileHandler(ctx *gin.Context){
	user,isExist := ctx.Get("user")
	if !isExist {
		ctx.JSON(500,gin.H{"error":"User not found"})
		return
	}

	ctx.JSON(http.StatusOK, dtos.UserResponse{
		ID: user.(repository.User).ID,
		UserName: user.(repository.User).UserName,
		Email: user.(repository.User).Email,
	})	
}