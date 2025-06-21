package app

import (
	"events/internal/generated/models"
	apiEvents "events/internal/generated/restapi/operations/events"
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func (srv *Service) CreateMovieEventHandler(params apiEvents.CreateMovieEventParams) middleware.Responder {
	event := &MovieEvent{
		MovieID:     *params.Body.MovieID,
		Title:       *params.Body.Title,
		Action:      *params.Body.Action,
		UserID:      params.Body.UserID,
		Rating:      params.Body.Rating,
		Genres:      params.Body.Genres,
		Description: params.Body.Description,
	}

	err := srv.Routers.EventBus.Publish(srv.ctx, event)
	if err != nil {
		responder := apiEvents.NewCreateMovieEventInternalServerError()
		return responder
	}

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	timestamp := strfmt.DateTime(time.Now())
	eventType := "movie-event"
	status := "success"

	payload := &models.EventResponse{
		Event: &models.Event{
			ID:        &id,
			Payload:   nil,
			Timestamp: &timestamp,
			Type:      &eventType,
		},
		Status: &status,
	}
	responder := apiEvents.NewCreateMovieEventCreated()
	responder.SetPayload(payload)

	return responder
}
