package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"travel-planning/services"

	"github.com/gorilla/mux"
)

type NotificationHandlers struct {
	NotService *services.NotificationService
}

func NewNotificationHandlers(notService *services.NotificationService) *NotificationHandlers {
	return &NotificationHandlers{
		NotService: notService,
	}
}

// GetMyNotifications godoc
// @Summary Get user notifications
// @Tags Notifications
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Notification
// @Router /api/notifications [get]
func (h *NotificationHandlers) GetMyNotifications(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		http.Error(w, "Unauthorized: User ID missing", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusUnauthorized)
		return
	}

	notifications, err := h.NotService.GetNotifications(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// MarkAsRead godoc
// @Summary Mark a notification as read
// @Tags Notifications
// @Security BearerAuth
// @Param id path int true "Notification ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/notifications/{id}/read [post]
func (h *NotificationHandlers) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, _ := strconv.Atoi(userIDStr)

	err = h.NotService.MarkAsRead(id, userID)
	if err != nil {
		http.Error(w, "Could not update notification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// GetUnreadCount godoc
// @Summary Get count of unread notifications
// @Tags Notifications
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]int
// @Router /api/notifications/unread-count [get]
func (h *NotificationHandlers) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, _ := strconv.Atoi(userIDStr)

	count, err := h.NotService.GetUnreadCount(userID)
	if err != nil {
		http.Error(w, "Failed to get count", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"unread_count": count})
}

// DeleteNotification godoc
// @Summary Delete a notification
// @Tags Notifications
// @Security BearerAuth
// @Param id path int true "Notification ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/notifications/{id} [delete]
func (h *NotificationHandlers) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, _ := strconv.Atoi(userIDStr)

	err = h.NotService.Delete(id, userID)
	if err != nil {
		http.Error(w, "Could not delete notification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
