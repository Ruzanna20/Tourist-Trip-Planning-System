package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"travel-planning/repository"

	"github.com/gorilla/mux"
)

type NotificationHandlers struct {
	NotRepo *repository.NotificationRepository
}

func NewNotificationHandlers(notRepo *repository.NotificationRepository) *NotificationHandlers {
	return &NotificationHandlers{
		NotRepo: notRepo,
	}
}

func (h *NotificationHandlers) GetMyNotifications(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("userID")
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDVal.(int)

	notifications, err := h.NotRepo.GetUserNotifications(userID)
	if err != nil {
		http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func (h *NotificationHandlers) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notifID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	userIDVal := r.Context().Value("userID")
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDVal.(int)

	err = h.NotRepo.MarkAsRead(notifID, userID)
	if err != nil {
		http.Error(w, "Failed to update notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
}

func (h *NotificationHandlers) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("userID")
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDVal.(int)

	count, err := h.NotRepo.GetUnreadCount(userID)
	if err != nil {
		http.Error(w, "Failed to get count", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"unread_count": count})
}
