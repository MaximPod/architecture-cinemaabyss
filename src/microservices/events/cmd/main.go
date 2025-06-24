// This file is safe to edit. Once it exists it will not be overwritten
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-openapi/loads"

	"events/internal/apiui"
	"events/internal/app"
	"events/internal/config"
	"events/internal/generated/restapi"
	"events/internal/generated/restapi/operations"
	"events/internal/slogger"

	apiEvents "events/internal/generated/restapi/operations/events"

	apiHealth "events/internal/generated/restapi/operations/health"
)

const ApplicationName = "events"

func main() {
	// context
	ctx := context.Background()

	// configuration
	config, err := config.Load(ApplicationName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// logger
	slogger.InitLogger(config.Logger)

	// service-application
	service, err := app.New(ctx, config)
	if err != nil {
		slog.Error("Main: failed to create service application", "error", err)
		os.Exit(1)
	}

	// swagger
	swagger := apiui.NewSwagger(config.Swagger)
	restapi.SwaggerJSON = swagger.UpdateSwaggerHostFromConfig(restapi.SwaggerJSON)

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		slog.Error("Main: failed to load swagger specification", "error", err)
		os.Exit(1)
	}

	// restapi
	api := operations.NewEventsAPI(swaggerSpec)

	api.EventsCreateMovieEventHandler = apiEvents.CreateMovieEventHandlerFunc(service.CreateMovieEventHandler)
	api.EventsCreatePaymentEventHandler = apiEvents.CreatePaymentEventHandlerFunc(service.CreatePaymentEventHandler)
	api.EventsCreateUserEventHandler = apiEvents.CreateUserEventHandlerFunc(service.CreateUserEventHandler)
	api.HealthGetEventsServiceHealthHandler = apiHealth.GetEventsServiceHealthHandlerFunc(service.GetEventsServiceHealthHandler)

	// web server
	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.Port = config.HTTPServer.Port

	// add swagger security control feature
	server.SetHandler(swagger.SecurityControl(server.GetHandler()))

	// configure server api
	server.ConfigureAPI()

	// callback from server: shutdown event send to service
	api.PreServerShutdown = service.PreServerShutdown
	api.ServerShutdown = service.OnShutdown

	// PubSub service handlers
	err = service.PubSubInitHandlers()
	if err != nil {
		slog.Error("Main: failed to add pubsub command handlers", "error", err)
		os.Exit(1)
	}

	//SignalNotifyRun отключен т.к. за сигналами оси в этой имплементации следит server
	//service.AddRunners(service.SignalNotifyRun)

	service.AddRunners(service.HandleShutdown)
	service.AddRunners(server.RunServerFunc)
	service.AddRunners(service.Routers.MessageRouter.Run)

	err = service.StartRunners()

	service.Exit(err)
}
