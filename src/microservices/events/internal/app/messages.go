package app

import "time"

// Коллекция сообщений сервиса.

// event messages
// события состоят из имени сущности и (опционально) глагола в прошедшем времени

// MovieEvent - Регистрирует новое событие, связанное с фильмом
type MovieEvent struct {
	MovieID     int64    `json:"movie_id"`
	Title       string   `json:"title"`
	Action      string   `json:"action"`
	UserID      int64    `json:"user_id"`
	Rating      float32  `json:"rating"`
	Genres      []string `json:"genres"`
	Description string   `json:"description"`
}

// UserEvent - Регистрирует новое событие, связанное с пользователем
type UserEvent struct {
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}

// PaymentEvent - Регистрирует новое событие, связанное с платежом
type PaymentEvent struct {
	PaymentID  int64     `json:"payment_id"`
	UserID     int64     `json:"user_id"`
	Amount     float32   `json:"amount"`
	Status     string    `json:"status"`
	Timestamp  time.Time `json:"timestamp"`
	MethodType string    `json:"method_type"`
}
