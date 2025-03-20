package transport

import (
	"net/http"
	"strconv"

	"example.com/go-api/internal/repository"
	"example.com/go-api/internal/services"
	"github.com/gin-gonic/gin"
)

// GetAllEventsHandler godoc
// @Summary Get all events
// @Description Retrieve a list of all events
// @Produce json
// @Tags Events
// @Success 200 {array} repository.Event
// @Failure 500 {object} map[string]string
// @Router /events [get]
func GetAllEventsHandler(ctx *gin.Context){
	events,err:= services.GetAllEvents()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK,events)
}

// GetEventById handles the HTTP GET request to retrieve an event by its ID.
// @Summary Get an event by ID
// @Description Retrieve an event's details using its unique identifier.
// @Accept json
// @Produce json
// @Tags Events
// @Param id path int true "Event ID"
// @Success 200 {object} repository.Event "Event details"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /events/{id} [get]
func GetEventById(ctx *gin.Context){
	idParam := ctx.Params.ByName("id")
	id,err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	event,err := services.GetEventById(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error":"Event not found"})
		return
	}
	ctx.JSON(http.StatusOK,event)
}


// SearchEventsHandler godoc
// @Summary Search events
// @Description Search for events by name with pagination support
// @Accept json
// @Produce json
// @Tags Events
// @Param name query string false "Event name"
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Success 200 {object} map[string]interface{} "Search results with pagination meta"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /events/search [get]
func SearchEventsHandler(ctx *gin.Context){
	name := ctx.Query("name")

	var err error
	var limit int
	var offset int
	var events []repository.Event
	limit,err = strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}
	offset, err = strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}
	var isLast bool
	var total int64
	events,total,isLast,err = services.SearchEvents(name,limit,offset)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err})
		return
	}
	ctx.JSON(http.StatusOK,gin.H{
		"data": events,
		"meta": gin.H{
			"isLast": isLast,
			"total": total,
		},
	})
}



// CreateEventHandler handles the HTTP POST request to create a new event.
// @Summary Create a new event
// @Description Create a new event with the provided details.
// @Accept json
// @Produce json
// @Tags Events
// @Param event body repository.Event true "Event details"
// @Success 201 {object} repository.Event "Created event"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events [post]
func CreateEventHandler(ctx *gin.Context){
	var event repository.Event
	err := ctx.ShouldBindBodyWithJSON(&event)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":err})
		return
	}
	err = services.CreateEvent(&event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err})
		return
	}
	ctx.JSON(http.StatusCreated,event)
}

// UpdateEventHandler handles the HTTP PUT request to update an existing event.
// @Summary Update an existing event
// @Description Update an existing event with the provided details.
// @Accept json
// @Produce json
// @Tags Events
// @Param id path int true "Event ID"
// @Param event body repository.Event true "Event details"
// @Success 200 {object} repository.Event "Updated event"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events/{id} [put]
func UpdateEventHandler(ctx *gin.Context){
		idParams := ctx.Params.ByName("id")
		id,err := strconv.Atoi(idParams)
		if err != nil {
			ctx.JSON((http.StatusBadRequest), gin.H{"error":"Invalid ID"})
			return
		}
		var event repository.Event
		err = ctx.ShouldBindBodyWithJSON(&event)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error":err})
			return
		}
		event.ID = uint(id)

		err = services.UpdateEvent(&event)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":err})
			return
		}
		ctx.JSON(http.StatusOK,event)
}

// DeleteEventHandler handles the HTTP DELETE request to delete an existing event.
// @Summary Delete an existing event
// @Description Delete an existing event by its ID.
// @Tags Events
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]string "Event deleted"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /events/{id} [delete]
func DeleteEventHandler(ctx *gin.Context){
	id,err := strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"Invalid ID"})
		return
	}
	idEvent := uint(id)
	err = services.DeleteEvent(idEvent)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message":"Event deleted"})
}