// This file is safe to edit. Once it exists it will not be overwritten
package config

import (
	"events/internal/apiui"
	"events/internal/pubsub"
	"events/internal/slogger"
)

// Config is the configuration for the application.
type Config struct {
	// web server configuration
	HTTPServer `yaml:"http_server"`

	// pubsub transport configuration
	PubSub pubsub.PubSubConfig `yaml:"pubsub"`

	// Swagger - конфигурация UI swagger-endpointa
	Swagger apiui.SwaggerConfig `yaml:"swagger"`

	// Log configuration
	Logger slogger.LogConfig `yaml:"logger"`
}

// HTTPServer - web server configuration
type HTTPServer struct {
	// Scheme - схема сервера: http, https or unix
	Scheme string `yaml:"scheme"`

	// Host - хост сервера: localhost
	Host string `yaml:"host"`

	// Port - порт сервера: 8080
	Port int `yaml:"port"`
}
