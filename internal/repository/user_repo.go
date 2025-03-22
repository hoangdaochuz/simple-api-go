package repository

import (
	"time"

	"example.com/go-api/internal/database"
)
type User struct {
	ID uint `json:"id" gorm:"primaryKey"`
	UserName string `json:"userName"`
	Email string `json:"email"`
	Password string `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}


// func Migrate(){
// 	database.DB.AutoMigrate(&User{});
// }

func GetUserByEmail(email string) (User,error){
	var user User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return user, err
}

func CreateUser(user *User) (error){
	return database.DB.Create(user).Error
}