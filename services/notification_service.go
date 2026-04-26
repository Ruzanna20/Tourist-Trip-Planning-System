package services

import (
	"fmt"
	"travel-planning/models"
	"travel-planning/repository"
)

type NotificationService struct {
	Repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		Repo: repo,
	}
}

func (s *NotificationService) GetNotifications(userID int) ([]models.Notification, error) {
	return s.Repo.GetUserNotifications(userID)
}

func (s *NotificationService) CreateNotification(userID int, tripID int, message string, msgType string) (int, error) {
	if message == "" {
		return 0, fmt.Errorf("message cannot be empty")
	}
	return s.Repo.SaveNotification(userID, tripID, message, msgType)
}

func (s *NotificationService) MarkAsRead(notificationID int, userID int) error {
	return s.Repo.MarkAsRead(notificationID, userID)
}

func (s *NotificationService) GetUnreadCount(userID int) (int, error) {
	return s.Repo.GetUnreadCount(userID)
}

func (s *NotificationService) Delete(notificationID int, userID int) error {
	return s.Repo.DeleteNotification(notificationID, userID)
}
