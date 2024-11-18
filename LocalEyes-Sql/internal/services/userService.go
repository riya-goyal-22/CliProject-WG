package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
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

func (s *UserService) Signup(username, password string, dwellingAge float64, answer string) error {
	hashedPassword := HashPassword(password)
	tag := utils.SetTag(dwellingAge)
	hashedAnswer := HashPassword(answer)
	fmt.Println("answer is :", answer, " hashed answer is :", hashedAnswer)
	user := &models.User{
		//UId:          uuid.New().String(),
		Username:         username,
		Password:         hashedPassword,
		City:             "delhi",
		Notification:     []string{},
		IsActive:         true,
		DwellingAge:      dwellingAge,
		Tag:              tag,
		SecurityQuestion: "Your nickname + your favourite food",
		SecurityAnswer:   hashedAnswer,
	}
	err := s.Repo.Create(user)
	return err
}

func (s *UserService) Login(username, password string) (*models.User, error) {
	hashedPassword := HashPassword(password)
	user, err := s.Repo.FindByUsernamePassword(username, hashedPassword)
	if err != nil {
		return nil, utils.InvalidAccountCredentials
	} else if user == nil {
		return nil, utils.InvalidAccountCredentials
	} else if user.IsActive == false {
		return nil, utils.InactiveUser
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

func (s *UserService) GetNotifications(uid string) ([]string, error) {
	user, err := s.GetUserById(uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NoUser
		}
		return nil, err
	}
	//err = s.Repo.ClearNotification(uid)
	//if err != nil {
	//	return nil, err
	//}
	return user.Notification, nil
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

func (s *UserService) PasswordReset(resetUser models.ResetPasswordUser) error {
	user, err := s.Repo.FindByUsername(resetUser.Username)
	if err != nil {
		return utils.NoUser
	}
	fmt.Println("user : ", user)
	if HashPassword(resetUser.SecurityAnswer) != user.SecurityAnswer {
		fmt.Println(".....", HashPassword(resetUser.SecurityAnswer), ".......", user.SecurityAnswer, "......")
		return utils.InvalidAnswer
	}
	hashedPassword := HashPassword(resetUser.NewPassword)
	var userUpdated = models.User{
		Username:         user.Username,
		Password:         hashedPassword,
		IsActive:         user.IsActive,
		UId:              user.UId,
		City:             user.City,
		DwellingAge:      user.DwellingAge,
		Tag:              user.Tag,
		Notification:     user.Notification,
		SecurityQuestion: user.SecurityQuestion,
		SecurityAnswer:   user.SecurityAnswer,
	}

	err = s.Repo.UpdateUser(user.UId, &userUpdated)
	if err != nil {
		fmt.Println("new user", userUpdated, "error : ", err.Error())
		return err
	}
	return nil
}

func (s *UserService) UpdateUser(uId string, requestUser *models.Client) error {
	var dwellingAge = (requestUser.LivingSince.Days / 365.0) + (requestUser.LivingSince.Years) + (requestUser.LivingSince.Months / 12.0)
	var hashedPassword = HashPassword(requestUser.Password)
	var hashedAnswer = HashPassword(requestUser.Answer)
	var tag = utils.SetTag(dwellingAge)
	user := models.User{
		Username:       requestUser.Username,
		Password:       hashedPassword,
		City:           requestUser.City,
		DwellingAge:    dwellingAge,
		SecurityAnswer: hashedAnswer,
		Tag:            tag,
	}
	err := s.Repo.UpdateUser(uId, &user)
	if err != nil {
		return err
	}
	return nil
}
