package repository

import (
	"time"

	"example.com/go-api/internal/database"
)

type Event struct{
	ID uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
	Name string `json:"name"`
	Description string `json:"description"`
	StartAt *time.Time `json:"startAt"`
}


func Migrate(){
	database.DB.AutoMigrate(&Event{})
}

func CreateEvent(event *Event) error{
	return database.DB.Create(event).Error
}

func GetAllEvents() ([]Event, error)  {
	var events []Event
	err :=  database.DB.Find(&events).Error
	return events, err
}



func GetEventById(id uint) (Event,error) {
	var event Event
	err :=database.DB.First(&event,id).Error
	return event, err
}

func UpdateEvent(event *Event) error {
	return database.DB.Save(event).Error
}

func DeleteEvent (id uint) error{
	return database.DB.Delete(&Event{},id).Error
}

func SearchEvents(name string, limit,offset int) ([]Event,int64, bool, error) {
	var events []Event
	var total int64
 	err := database.DB.Where("name LIKE ?","%"+name+"%").Limit(limit).Offset(offset).Find(&events).Count(&total).Error
	isLast := offset + len(events) >= int(total)
	return events,total, isLast, err
}

