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
	"tender/internal/models"
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

type Token struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

func NewFCMHandler(client *messaging.Client, db *sql.DB) *FCMHandler {
	return &FCMHandler{Client: client, DB: db}
}

func (h *FCMHandler) SendMessage(ctx context.Context, token string, UserId int, title, body string) error {
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
				"apns-priority": "10",
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
	} else {
		err = h.CreateNotify(UserId, title, body)
		if err != nil {
			log.Printf("Ошибка при отправке уведомления: %v", err)
			return err
		}

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
		err = h.SendMessage(ctx, token, req.UserId, req.Title, req.Body)
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
	var newToken Token

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
        ( user_id, token) 
        VALUES ( ?, ?);`

	_, err := h.DB.Exec(stmt1, clientID, token)
	if err != nil {
		fmt.Println("safjdajs")
		return err
	}
	return nil
}

func (h *FCMHandler) DeleteToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get(":id")

	if token == "" {
		http.Error(w, "Failed to fetch tokens", http.StatusInternalServerError)
		return
	}

	err := h.DeleteTokenRep(token)
	if err != nil {
		http.Error(w, "Failed to delete tokens", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FCMHandler) DeleteTokenRep(token string) error {
	stmt := `DELETE FROM notify_tokens WHERE token = ?`
	_, err := h.DB.Exec(stmt, token)
	if err != nil {
		return err
	}

	return nil
}

func (h *FCMHandler) CreateNotify(clientID int, title, body string) error {

	stmt1 := `
        INSERT INTO notify_history 
        ( user_id, title, body) 
        VALUES ( ?, ?, ?);`

	_, err := h.DB.Exec(stmt1, clientID, title, body)
	if err != nil {
		return err
	}
	return nil
}

func (h *FCMHandler) ShowNotifyHistory(w http.ResponseWriter, r *http.Request) {
	var newNotify models.Notify

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&newNotify)
	if err != nil {
		http.Error(w, "Failed to fetch tokens", http.StatusBadRequest)
		return
	}

	notify, err := h.GetNotifyHistory(newNotify.UserID)
	if err != nil {
		log.Printf("Error fetching notifications: %v", err)
		http.Error(w, "Failed to fetch notification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notify)
}

func (h *FCMHandler) GetNotifyHistory(userId int) ([]models.Notify, error) {
	if h.DB == nil {
		log.Print("h.DB is nil")
		return []models.Notify{}, fmt.Errorf("database connection is not initialized")
	}

	var notifications []models.Notify
	query := "SELECT * FROM notify_history WHERE user_id = ?"
	rows, err := h.DB.Query(query, userId)
	if err != nil {
		return []models.Notify{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var notification models.Notify
		if err := rows.Scan(&notification); err != nil {
			return []models.Notify{}, err
		}
		notifications = append(notifications, notification)
	}

	if err = rows.Err(); err != nil {
		return []models.Notify{}, err
	}

	return notifications, nil
}

func (h *FCMHandler) DeleteNotifyHistory(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")

	if id == "" {
		http.Error(w, "Failed to fetch tokens", http.StatusInternalServerError)
		return
	}

	err := h.DeleteNotifyRep(id)
	if err != nil {
		http.Error(w, "Failed to delete tokens", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FCMHandler) DeleteNotifyRep(id string) error {
	stmt := `DELETE FROM notify_history WHERE id = ?`
	_, err := h.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
