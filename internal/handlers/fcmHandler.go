package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"firebase.google.com/go/messaging"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Add the appropriate database driver
	"io"
	"log"
	"net/http"
)

type FCMHandler struct {
	Client *messaging.Client
	DB     *sql.DB
}

type NotificationRequest struct {
	Id     int    `json:"id"`
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func NewFCMHandler(client *messaging.Client, db *sql.DB) *FCMHandler {
	return &FCMHandler{Client: client, DB: db}
}

func (h *FCMHandler) SendMessage(ctx context.Context, token, title, body string) error {
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				ChannelID: "high_priority_channel",
			},
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10", // Immediate delivery
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: title,
						Body:  body,
					},
					Sound: "default",
				},
			},
		},
	}

	response, err := h.Client.Send(ctx, message)
	if err != nil {
		log.Printf("Ошибка при отправке уведомления: %v", err)
		return err
	}

	log.Printf("Отправка уведомления выполнена успешно: %s\n", response)
	return nil
}

func (h *FCMHandler) NotifyChange(w http.ResponseWriter, r *http.Request) {
	var req NotificationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received notification request: %+v", req)

	ctx := r.Context()
	tokens, err := h.GetTokensByClientID(req.UserId)
	if err != nil {
		log.Printf("Error fetching tokens: %v", err)
		http.Error(w, "Failed to fetch tokens", http.StatusInternalServerError)
		return
	}

	// Send notifications to each token
	for _, token := range tokens {
		err = h.SendMessage(ctx, token, req.Title, req.Body)
		if err != nil {
			log.Printf("Error sending notification to token %s: %v", token, err)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification sent successfully"))
}

func (h *FCMHandler) GetTokensByClientID(clientID int) ([]string, error) {
	if h.DB == nil {
		log.Print("h.DB is nil")
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var tokens []string
	query := "SELECT token FROM notify_tokens WHERE user_id = ?"
	rows, err := h.DB.Query(query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	fmt.Println("good")
	return tokens, nil
}

func (h *FCMHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
	var newToken NotificationRequest

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&newToken)
	if err != nil {
		http.Error(w, "Failed to fetch tokens", http.StatusBadRequest)
		return
	}

	err = h.InsertToken(newToken.UserId, newToken.Token)
	if err != nil {
		http.Error(w, "Failed to insert tokens", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *FCMHandler) InsertToken(clientID int, token string) error {

	stmt1 := `
        INSERT INTO notify_tokens 
        (id, user_id, token) 
        VALUES (1, ?, ?);`

	_, err := h.DB.Exec(stmt1, clientID, token)
	if err != nil {
		return err
	}
	return nil
}
