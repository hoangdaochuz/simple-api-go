package services

import "example.com/go-api/internal/repository"

func GetUserByEmailService(email string) (repository.User,error){
	return repository.GetUserByEmail(email)
}

func CreateUser(user *repository.User) error{
	return repository.CreateUser(user)
}