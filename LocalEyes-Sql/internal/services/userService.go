package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"localEyes/constants"
	"localEyes/internal/interfaces"
	"localEyes/internal/models"
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
		//UId:          primitive.NewObjectID(),
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

func (s *UserService) Login(Username, password string) (*models.User, error) {
	hashedPassword := HashPassword(password)
	user, err := s.Repo.FindByUsernamePassword(Username, hashedPassword)
	//user.NotifyChannel = make(chan string, 5)
	if err != nil {
		return nil, errors.New(constants.Red + "Invalid Account credentials" + constants.Reset)
	} else if user == nil {
		return nil, errors.New(constants.Red + "Invalid Account credentials" + constants.Reset)
	} else if user.IsActive == false {
		return nil, errors.New(constants.Red + "InActive Account" + constants.Reset)
	}
	return user, nil
}

func (s *UserService) DeActivate(UId int) error {
	err := s.Repo.UpdateActiveStatus(UId, false)
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

func (s *UserService) NotifyUsers(UId int, title string) error {
	return s.Repo.PushNotification(UId, title)
}

func (s *UserService) UnNotifyUsers(UId int) error {
	err := s.Repo.ClearNotification(UId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}
