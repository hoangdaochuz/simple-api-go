package services

import "example.com/go-api/internal/repository"


func GetAllEvents() ([]repository.Event,error){
	return repository.GetAllEvents()
}

func GetEventById(id uint) (repository.Event,error){
	return repository.GetEventById(id)
}

func CreateEvent(event *repository.Event) error {
	return repository.CreateEvent(event)
}

func UpdateEvent(event *repository.Event) error{
	return repository.UpdateEvent(event)
}

func DeleteEvent(id uint) error{
	return repository.DeleteEvent(id)
}


func SearchEvents(name string, limit,offset int) ([]repository.Event,int64,bool,error) {
	return repository.SearchEvents(name,limit,offset)
}
