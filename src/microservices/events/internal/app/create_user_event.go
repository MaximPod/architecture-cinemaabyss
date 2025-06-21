package app

import (
	"events/internal/generated/models"
	apiEvents "events/internal/generated/restapi/operations/events"
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func (srv *Service) CreateUserEventHandler(params apiEvents.CreateUserEventParams) middleware.Responder {
	event := &UserEvent{
		UserID:    *params.Body.UserID,
		Username:  params.Body.Username,
		Email:     params.Body.Email,
		Action:    *params.Body.Action,
		Timestamp: time.Time(*params.Body.Timestamp),
	}

	err := srv.Routers.EventBus.Publish(srv.ctx, event)
	if err != nil {
		responder := apiEvents.NewCreateUserEventInternalServerError()
		return responder
	}

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	timestamp := strfmt.DateTime(time.Now())
	eventType := "user-event"
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

	responder := apiEvents.NewCreateUserEventCreated()
	responder.SetPayload(payload)

	return responder
}
