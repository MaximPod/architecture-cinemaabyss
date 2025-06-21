// This file is safe to edit. Once it exists it will not be overwritten
package apiui

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-openapi/spec"
)

// SwaggerConfig - конфигурация swagger
type SwaggerConfig struct {
	// Enable - флаг доступности Swagger UI
	Enable bool `yaml:"enable"`
	// Host - хост swagger-сервера: localhost
	Host string `yaml:"host"`
	// Port - порт swagger-сервера: 8080
	Port int `yaml:"port"`
	// Endpoint - эндпойнт Swagger UI ("/docs", "/swagger/index.html")
	Endpoint string `yaml:"endpoint"`
}

// Swagger - сущность сервиса swagger
type Swagger struct {
	Enable   bool   // флаг доступности Swagger UI
	Host     string // хост swagger-сервера: localhost
	Port     int    // порт swagger-сервера: 8080
	Endpoint string // эндпойнт Swagger UI ("/docs", "/swagger/index.html")
}

func NewSwagger(cfg SwaggerConfig) *Swagger {
	return &Swagger{
		Enable:   cfg.Enable,
		Host:     cfg.Host,
		Port:     cfg.Port,
		Endpoint: cfg.Endpoint,
	}
}

// UpdateSwaggerHostFromConfig обновляет swagger host из конфига.
// В swagger-спецификации прописывется dev-host для безопасной работы с эндпойнтами
// из любой среды. В конфиге указывается необходимый для задач: localhost:8080
func (s *Swagger) UpdateSwaggerHostFromConfig(currentSpec json.RawMessage) json.RawMessage {
	if (s.Host == "") || (s.Port == 0) {
		return currentSpec
	}

	nextSpec := new(spec.Swagger)
	err := json.Unmarshal(currentSpec, nextSpec)
	if err != nil {
		slog.Error("Apiui: json.Unmarshal swagger spec", "error", err)
		return currentSpec
	}

	nextSpec.Host = fmt.Sprintf("%s:%d", s.Host, s.Port)
	res, err := nextSpec.MarshalJSON()
	if err != nil {
		slog.Error("Apiui: json.Marshal next spec", "error", err)
		return currentSpec

	}

	return res
}

// SecurityControl - проверяет доступность к ендпойнту Swagger UI и /swagger.json по флагу Enable
func (s *Swagger) SecurityControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case s.Endpoint, "/swagger.json":
			// проверка флага swagger disabled
			if !s.Enable {
				w.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Fprintf(w, "405: Swagger UI disabled")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
