package app

import (
	apiHealth "events/internal/generated/restapi/operations/health"

	"github.com/go-openapi/runtime/middleware"
)

func (srv *Service) GetEventsServiceHealthHandler(params apiHealth.GetEventsServiceHealthParams) middleware.Responder {

	payload := &apiHealth.GetEventsServiceHealthOKBody{
		Status: true,
	}

	responder := apiHealth.NewGetEventsServiceHealthOK()
	responder.SetPayload(payload)

	return responder
}
