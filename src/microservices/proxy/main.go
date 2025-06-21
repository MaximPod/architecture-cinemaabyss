package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MonolithURL            string `envconfig:"MONOLITH_URL" default:"http://monolith:8080"`
	MoviesServiceURL       string `envconfig:"MOVIES_SERVICE_URL" default:"http://movies-service:8081"`
	MoviesMigrationPercent int    `envconfig:"MOVIES_MIGRATION_PERCENT" default:"0"`
	GradualMigration       bool   `envconfig:"GRADUAL_MIGRATION" default:"false"`
}

func main() {
	// Кофигурация сервиса
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации из переменных окружения: %s", err)
	}

	// Создаем новый генератор случайных чисел
	rng := rand.New(rand.NewSource(rand.Int63()))

	// Адреса целевых сервисов
	monolithURL, err := url.Parse(cfg.MonolithURL)
	if err != nil {
		log.Fatal(err)
	}

	moviesServiceURL, err := url.Parse(cfg.MoviesServiceURL)
	if err != nil {
		log.Fatal(err)
	}

	//  Если фичефлаг не поднят, все идет на монолит
	if cfg.GradualMigration {
		log.Println("Gradual migration is enabled")
		log.Printf("MOVIES_MIGRATION_PERCENT is %d\n", cfg.MoviesMigrationPercent)
	} else {
		log.Println("Gradual migration is disabled")
		cfg.MoviesMigrationPercent = 0
	}

	// Создаем reverse proxy для каждого сервиса
	monolith := httputil.NewSingleHostReverseProxy(monolithURL)
	moviesService := httputil.NewSingleHostReverseProxy(moviesServiceURL)

	// Обработчик входящих запросов
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		switch r.URL.Path {
		// В случае запроса к /api/movies отправляем его на сервисы фильмов/монолит в зависимости от фичефлага и процента
		case "/api/movies":
			// Генерируем случайное число от 0 до 99
			randomNum := rng.Intn(99)
			
			// Распределяем запросы по случайному числу
			if randomNum < cfg.MoviesMigrationPercent {
				log.Printf("Routing to moviesService (%d < %d)", randomNum, cfg.MoviesMigrationPercent)
				moviesService.ServeHTTP(w, r)
			} else {
				log.Printf("Routing to monolith (%d >= %d)", randomNum, cfg.MoviesMigrationPercent)
				monolith.ServeHTTP(w, r)
			}
			return
		// В случае запроса к любому другому адресу отправляем его на монолит
		default:
			log.Println("Routing to monolith")
			monolith.ServeHTTP(w, r)
			return
		}

	})

	// Запускаем сервер
	log.Println("Proxy server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
