package service

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/util"

	"github.com/G9QBootcamp/qoli-survey/internal/auth/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/auth/models"
	"github.com/G9QBootcamp/qoli-survey/internal/auth/repository"
	userModels "github.com/G9QBootcamp/qoli-survey/internal/user/models"
	userRepository "github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"golang.org/x/net/context"
)

type IAuthService interface {
	SaveOTP(context.Context, uint) (string, error)
	SendOTPEmail(context.Context, uint, string) error
	VerifyOTP(context.Context, dto.VerifyOTPRequest) (bool, error)
	Signup(c context.Context, req dto.SignupRequest) (*userModels.User, error)
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

func (s *AuthService) Signup(c context.Context, req dto.SignupRequest) (*userModels.User, error) {
	user := userModels.User{
		NationalID:   req.NationalID,
		Email:        req.Email,
		PasswordHash: req.Password,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		City:         req.City,
	}

	if req.DateOfBirth != "" {
		dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			s.logger.Error(logging.Internal, logging.FailedToParseDate, "failed to parse date", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})

			return nil, errors.New("invalid date format")
		}
		user.DateOfBirth = dateOfBirth
	}

	if s.userRepo.IsEmailOrNationalIDTaken(c, user.Email, user.NationalID) {
		return nil, errors.New("email or national ID already in use")
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		s.logger.Error(logging.Internal, logging.HashPassword, "failed to hash password", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})

		return nil, err
	}
	user.PasswordHash = hashedPassword

	userCount, err := s.userRepo.GetUserCount(c)
	if err != nil {
		s.logger.Error(logging.Internal, logging.FailedToGetUserCount, "error in get user count", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})

		return nil, err
	}

	if userCount == 0 {
		superAdminRole, err := s.userRepo.GetRoleByName(c, "Super Admin")
		if err != nil {
			s.logger.Error(logging.Database, logging.FailedToGetRole, "failed to fetch super admin role", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})
		}
		user.GlobalRole = *superAdminRole
	} else {
		userRole, err := s.userRepo.GetRoleByName(c, "User")
		if err != nil {
			s.logger.Error(logging.Database, logging.FailedToGetRole, "failed to fetch user role", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})
			return nil, err
		}
		user.GlobalRole = *userRole
	}

	_, err = s.userRepo.CreateUser(c, &user)
	if err != nil {
		s.logger.Error(logging.Internal, logging.FailedToCreateUser, "error in create user", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})

		return nil, err
	}

	return &user, nil
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

func (s *AuthService) VerifyOTP(c context.Context, req dto.VerifyOTPRequest) (bool, error) {
	user, err := s.userRepo.GetUserByEmail(c, req.Email)
	if err != nil || user == nil {
		return false, errors.New("user not found")
	}

	otpRecord, err := s.repo.GetOTPByUserID(c, user.ID)
	fmt.Println(err != nil, otpRecord == nil, !otpRecord.IsValid, otpRecord.Code != req.OTP, otpRecord.ExpiresAt.Before(time.Now()))
	if err != nil || otpRecord == nil || !otpRecord.IsValid || otpRecord.Code != req.OTP || otpRecord.ExpiresAt.Before(time.Now()) {
		return false, errors.New("invalid or expired OTP")
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
