package app

import (
	"events/internal/generated/models"
	apiEvents "events/internal/generated/restapi/operations/events"
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func (srv *Service) CreatePaymentEventHandler(params apiEvents.CreatePaymentEventParams) middleware.Responder {
	event := &PaymentEvent{
		PaymentID:  *params.Body.PaymentID,
		UserID:     *params.Body.UserID,
		Amount:     *params.Body.Amount,
		Status:     *params.Body.Status,
		Timestamp:  time.Time(*params.Body.Timestamp),
		MethodType: params.Body.MethodType,
	}

	err := srv.Routers.EventBus.Publish(srv.ctx, event)
	if err != nil {
		responder := apiEvents.NewCreatePaymentEventInternalServerError()
		return responder
	}

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	timestamp := strfmt.DateTime(time.Now())
	eventType := "payment-event"
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

	responder := apiEvents.NewCreatePaymentEventCreated()
	responder.SetPayload(payload)

	return responder
}
