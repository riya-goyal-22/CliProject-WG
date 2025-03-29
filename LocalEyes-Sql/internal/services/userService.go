package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"localEyes/internal/interfaces"
	"localEyes/internal/models"
	"localEyes/utils"
	"os"
	"strconv"
)

type UserService struct {
	Repo    interfaces.UserRepository
	OtpRepo interfaces.OTPRepoInterface
}

func NewUserService(repo interfaces.UserRepository, OtpRepo interfaces.OTPRepoInterface) *UserService {
	return &UserService{
		Repo:    repo,
		OtpRepo: OtpRepo,
	}
}

func (s *UserService) Signup(username, password, email string, dwellingAge float64) error {
	hashedPassword := HashPassword(password)
	tag := utils.SetTag(dwellingAge)
	user := &models.User{
		Username:     username,
		Password:     hashedPassword,
		City:         "delhi",
		Notification: []string{},
		IsActive:     true,
		DwellingAge:  dwellingAge,
		Tag:          tag,
		Email:        email,
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
func (s *UserService) ValidateEmail(email string) bool {
	if email == "localeyes22@gmail.com" {
		return false
	}
	user, err := s.Repo.FindByUserMail(email)
	if user == nil || err != nil {
		return true
	}
	return false
}

func (s *UserService) SendOtp(email string) error {
	otp, err := s.OtpRepo.GenerateOTP()
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", os.Getenv("SMTPSenderEmail"))
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Go SMTP Test")
	message.SetBody("text/plain", "Hello,\r\nThis is your otp to reset password: "+otp)

	port, _ := strconv.Atoi(os.Getenv("SMTPPort"))
	// Create a dialer with SMTP server information
	dialer := gomail.NewDialer(
		os.Getenv("SMTPServer"),
		port,
		os.Getenv("SMTPSenderEmail"),
		os.Getenv("SMTPSenderPassword"),
	)
	if os.Getenv("SMTPServer") == "" || os.Getenv("SMTPPort") == "" || os.Getenv("SMTPSenderEmail") == "" || os.Getenv("SMTPSenderPassword") == "" {
		return fmt.Errorf("missing required environment variables for SMTP configuration")

	}

	// Send the email via the dialer
	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	s.OtpRepo.SaveOTP(email, otp)
	return nil
}

func (s *UserService) PasswordReset(resetUser models.ResetPasswordUser) error {
	if s.OtpRepo.ValidateOTP(resetUser.Email, resetUser.OTP) {
		user, err := s.Repo.FindByUserMail(resetUser.Email)
		if err != nil {
			return utils.WrongOTP
		}
		hashedPassword := HashPassword(resetUser.NewPassword)
		var userUpdated = models.User{
			Username:     user.Username,
			Password:     hashedPassword,
			IsActive:     user.IsActive,
			UId:          user.UId,
			City:         user.City,
			DwellingAge:  user.DwellingAge,
			Tag:          user.Tag,
			Notification: user.Notification,
			Email:        user.Email,
		}
		err = s.Repo.UpdateUser(&userUpdated)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *UserService) UpdateUser(uId string, requestUser *models.Client) error {
	var dwellingAge = (requestUser.LivingSince.Days / 365.0) + (requestUser.LivingSince.Years) + (requestUser.LivingSince.Months / 12.0)
	var hashedPassword = HashPassword(requestUser.Password)
	var tag = utils.SetTag(dwellingAge)
	user := models.User{
		UId:         uId,
		Username:    requestUser.Username,
		Password:    hashedPassword,
		City:        requestUser.City,
		DwellingAge: dwellingAge,
		Tag:         tag,
	}
	err := s.Repo.UpdateUser(&user)
	if err != nil {
		return err
	}
	return nil
}
