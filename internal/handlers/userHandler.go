package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"tender/internal/models"
	"tender/internal/services"
)

type UserHandler struct {
	Service *services.UserService
}

type BalanceUpdateRequest struct {
	Amount float64 `json:"balance"`
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdUser, err := h.Service.SignUp(r.Context(), user)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) || errors.Is(err, models.ErrDuplicatePhone) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (h *UserHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userInfo, err := h.Service.LogIn(r.Context(), user)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) || errors.Is(err, models.ErrInvalidPassword) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(userInfo)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.Service.GetUserByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var balanceUpdate BalanceUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&balanceUpdate)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.Service.UpdateBalance(r.Context(), id, balanceUpdate.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	balance, err := h.Service.GetBalance(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
	w.WriteHeader(http.StatusOK)

}

func (h *UserHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	balance, err := h.Service.GetBalance(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}

func (h *UserHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteUserByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrPermissionNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user.ID = id

	updatedUser, err := h.Service.UpdateUser(r.Context(), user)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.Service.ChangePassword(r.Context(), id, req.OldPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else if errors.Is(err, models.ErrInvalidPassword) {
			http.Error(w, "Old password is incorrect", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password updated successfully"))
}

func (h *UserHandler) SendRecoveryHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	{
		var passwordRecoveryRequest struct {
			Email string `json:"email"`
		}

		body, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Failed to fetch tokens", http.StatusBadRequest)
			return
		}

		user.ID, err = h.Service.FindUserByEmail(r.Context(), passwordRecoveryRequest.Email)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "no client found"):
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	SendRecovery(user.Email, user.ID)
	h.sendResponseWithMessage(w, fmt.Sprintf("ссылка для восстанавления успешно отправлено в почту пользователя %s", user.Email), http.StatusOK)
}

func SendRecovery(to string, userId int) {
	from := "tendercommunitykgz@gmail.com"
	password := "tendercommunity12"
	subject := "Восстановление пароля в системе TENDER-COMMUNITY"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	link := fmt.Sprintf("http://176.126.166.163:4000/password/recovery/mail?user_id=%d", userId)

	body := fmt.Sprintf(`
		<html>
		<body>
			<p>Здравствуйте,</p>
			<p>Для восстановления пароля, пожалуйста, перейдите по следующей ссылке:</p>
			<p><a href="%s">Восстановить пароль</a></p>
			<p>Если вы не запрашивали восстановление пароля, просто игнорируйте это письмо.</p>
		</body>
		</html>`, link)

	message := fmt.Sprintf("From: %s\r\n", from) +
		fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		fmt.Sprintf("MIME-Version: 1.0\r\n") +
		fmt.Sprintf("Content-Type: text/html; charset=\"UTF-8\"\r\n") +
		"\r\n" + body

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))

	if err != nil {
		fmt.Println("Ошибка отправки письма:", err)
		return
	}

	fmt.Println("Письмо успешно отправлено!")
}

func (h *UserHandler) sendResponseWithMessage(w http.ResponseWriter, message string, status int) {
	messageResponse := models.MessageResponse{Message: message}
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(messageResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) PasswordRecoveryHandler(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Query().Get("user_id")
	if user_id == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "tendercommunity://reset_password?hash="+user_id, http.StatusSeeOther)
}
