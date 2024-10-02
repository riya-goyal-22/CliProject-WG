package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"localEyes/internal/interfaces"
	"localEyes/internal/models"
	"localEyes/utils"
)

type UserService struct {
	Repo interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) Signup(username, password string, dwellingAge int, tag string) error {
	hashedPassword := HashPassword(password)

	user := &models.User{
		//UId:          uuid.New().String(),
		Username:     username,
		Password:     hashedPassword,
		City:         "delhi",
		Notification: []string{},
		IsActive:     true,
		DwellingAge:  dwellingAge,
		Tag:          tag,
	}
	err := s.Repo.Create(user)
	return err
}

func (s *UserService) Login(username, password string) (*models.User, error) {
	hashedPassword := HashPassword(password)
	user, err := s.Repo.FindByUsernamePassword(username, hashedPassword)
	if err != nil {
		return nil, errors.New("invalid Account credentials")
	} else if user == nil {
		return nil, errors.New("invalid Account credentials")
	} else if user.IsActive == false {
		return nil, errors.New("inActive Account")
	}
	return user, nil
}

func (s *UserService) DeActivate(uid string) error {
	err := s.Repo.UpdateActiveStatus(uid, false)
	if err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func (s *UserService) NotifyUsers(uid, title string) error {
	return s.Repo.PushNotification(uid, title)
}

//func (s *UserService) UnNotifyUsers(uid string) error {
//	err := s.Repo.ClearNotification(uid)
//	if errors.Is(err, sql.ErrNoRows) {
//		return nil
//	}
//	return err
//}

func (s *UserService) GetNotifications(uid string) (*[]string, error) {
	user, err := s.GetUserById(uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NoUser
		}
		return nil, err
	}
	err = s.Repo.ClearNotification(uid)
	if err != nil {
		return nil, err
	}
	return &user.Notification, nil
}

func (s *UserService) GetUserById(uid string) (*models.User, error) {
	user, err := s.Repo.FindByUId(uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) ValidateUsername(username string) bool {
	if username == "admin" || username == "Admin" {
		return false
	}
	user, err := s.Repo.FindByUsername(username)
	if user == nil || err != nil {
		return true
	}
	return false
}
