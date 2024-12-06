package service

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/config"

	"github.com/G9QBootcamp/qoli-survey/internal/auth/models"
	"github.com/G9QBootcamp/qoli-survey/internal/auth/repository"
	userRepository "github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"golang.org/x/net/context"
)

type IAuthService interface {
	SaveOTP(context.Context, uint) (string, error)
	SendOTPEmail(context.Context, uint, string) error
	VerifyOTP(context.Context, uint, string) (bool, error)
}

type AuthService struct {
	conf     *config.Config
	repo     repository.IAuthRepository
	userRepo userRepository.IUserRepository
	logger   logging.Logger
}

func New(conf *config.Config, repo repository.IAuthRepository, userRepo userRepository.IUserRepository, logger logging.Logger) *AuthService {
	return &AuthService{conf: conf, repo: repo, userRepo: userRepo, logger: logger}
}

func (s *AuthService) SaveOTP(c context.Context, userID uint) (string, error) {
	otp := fmt.Sprintf("%06d", randInt(100000, 999999))

	expiresAt := time.Now().Add(10 * time.Minute)

	otpRecord := models.OTP{
		UserID:    userID,
		Code:      otp,
		ExpiresAt: expiresAt,
		IsValid:   true,
	}

	err := s.repo.CreateOTP(c, otpRecord)
	if err != nil {
		s.logger.Error(logging.Internal, "FailedToSaveOTP", "Error saving OTP", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return "", err
	}

	return otp, nil
}

func (s *AuthService) VerifyOTP(c context.Context, userID uint, otp string) (bool, error) {
	otpRecord, err := s.repo.GetOTPByUserID(c, userID)
	fmt.Println(err != nil, otpRecord == nil, !otpRecord.IsValid, otpRecord.Code != otp, otpRecord.ExpiresAt.Before(time.Now()))
	if err != nil || otpRecord == nil || !otpRecord.IsValid || otpRecord.Code != otp || otpRecord.ExpiresAt.Before(time.Now()) {
		return false, errors.New("invalid or expired OTP")
	}

	user, err := s.userRepo.GetUserByID(c, userID)
	if err != nil || user == nil {
		return false, errors.New("user not found")
	}

	user.EmailVerified = true
	s.userRepo.UpdateUser(c, user)

	otpRecord.IsValid = false
	err = s.repo.UpdateOTP(c, otpRecord)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *AuthService) SendOTPEmail(c context.Context, userID uint, otp string) error {
	user, err := s.userRepo.GetUserByID(c, userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	subject := "Your OTP Code"
	body := fmt.Sprintf("Hello %s %s,\n\nYour OTP code is: %s\n\nIt is valid for 10 minutes.", user.FirstName, user.LastName, otp)
	msg := []byte("To: " + user.Email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", s.conf.Email.SMTPUser, s.conf.Email.SMTPPass, s.conf.Email.SMTPServer)

	err = smtp.SendMail(s.conf.Email.SMTPServer+":"+s.conf.Email.SMTPPort, auth, s.conf.Email.FromEmail, []string{user.Email}, msg)
	if err != nil {
		log.Printf("Error sending OTP email to user %s: %v", user.Email, err)
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	log.Printf("OTP email sent to user %s", user.Email)
	return nil
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}
