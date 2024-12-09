package service

import (
	"errors"
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"golang.org/x/net/context"

	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/util"
	"github.com/G9QBootcamp/qoli-survey/pkg/jwtutils"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
)

type IUserService interface {
	GetUsers(context.Context, dto.UserGetRequest) ([]*dto.UserResponse, error)
	SetMaxSurveys(ctx context.Context, userID string, maxSurveys int) error
	Login(c context.Context, req dto.LoginRequest) (string, time.Time, error)
	UpdateUserProfile(c context.Context, userID uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	GetUser(c context.Context, id uint) (*dto.UserResponse, error)
	GetBalance(ctx context.Context, userID uint) (float64, error)
	ProcessTransaction(ctx context.Context, buyerID, sellerID uint, amount float64, voteCount int) error
	Deposit(ctx context.Context, userID uint, amount float64) error
	Transfer(ctx context.Context, senderID, receiverID uint, amount float64) error
	Withdraw(ctx context.Context, userID uint, amount float64) error
}
type UserService struct {
	conf   *config.Config
	repo   repository.IUserRepository
	logger logging.Logger
}

func New(conf *config.Config, repo repository.IUserRepository, logger logging.Logger) *UserService {
	return &UserService{conf: conf, repo: repo, logger: logger}
}

func (s *UserService) GetUsers(c context.Context, r dto.UserGetRequest) ([]*dto.UserResponse, error) {
	userFilters := dto.UserFilters{Name: r.Name}
	users, err := s.repo.GetUsers(c, userFilters)
	if err != nil {
		s.logger.Error(logging.Internal, logging.FailedToGetUsers, "error in get users", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})
		return nil, err
	}
	usersResponse := []*dto.UserResponse{}

	for _, user := range users {
		usersResponse = append(usersResponse, ToUserResponse(&user))
	}
	return usersResponse, nil
}

func (s *UserService) UpdateUserProfile(c context.Context, userID uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetUserByID(c, userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.DateOfBirth != "" {
		dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth) // Assuming format "YYYY-MM-DD"
		if err != nil {
			return nil, errors.New("invalid date format")
		}
		if time.Since(user.CreatedAt) > 24*time.Hour {
			return nil, errors.New("date of birth cannot be updated after 24 hours of registration")
		}
		user.DateOfBirth = dateOfBirth
	}
	if req.City != "" {
		user.City = req.City
	}

	updatedUser, err := s.repo.UpdateUser(c, user)
	if err != nil {
		s.logger.Error(logging.Internal, logging.FailedToUpdateUser, "error in updating user", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})

		return nil, err
	}

	return ToUserResponse(updatedUser), nil
}

func ToUserResponse(user *models.User) *dto.UserResponse {
	var dateOfBirth string

	if !user.DateOfBirth.IsZero() {
		formattedDate := user.DateOfBirth.Format("2006-01-02")
		dateOfBirth = formattedDate
	}

	return &dto.UserResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DateOfBirth: dateOfBirth,
		City:        user.City,
	}
}

func (s *UserService) Login(c context.Context, req dto.LoginRequest) (string, time.Time, error) {
	user, err := s.repo.GetUserByEmail(c, req.Email)
	if err != nil {
		s.logger.Error(logging.Internal, logging.UserNotAuthorized, "invalid email", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})
		return "", time.Time{}, errors.New("invalid email")
	}

	if !user.EmailVerified {
		s.logger.Info(logging.Internal, logging.UserNotAuthorized, "email is not verified", map[logging.ExtraKey]interface{}{logging.Service: "UserService"})
		return "", time.Time{}, errors.New("email is not verified")
	}

	if err := util.CheckPassword(req.Password, user.PasswordHash); err != nil {
		s.logger.Error(logging.Internal, logging.UserNotAuthorized, "invalid email", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})
		return "", time.Time{}, errors.New("invalid password")
	}

	expiresAt := time.Now().Add(time.Duration(s.conf.JWT.ExpireMinutes) * time.Minute)
	token, err := jwtutils.GenerateToken(user.ID, user.GlobalRole.Name, s.conf.JWT.SecretKey, s.conf.JWT.ExpireMinutes)
	if err != nil {
		s.logger.Error(logging.Internal, logging.FailedToGenerateToken, "generate token error", map[logging.ExtraKey]interface{}{logging.Service: "UserService", logging.ErrorMessage: err.Error()})
		return "", time.Time{}, errors.New("failed to generate token")
	}

	return token, expiresAt, nil
}

func (s *UserService) GetUser(c context.Context, id uint) (*dto.UserResponse, error) {
	survey, err := s.repo.GetUserByID(c, id)

	if err != nil {
		return nil, err
	}

	sResponse := dto.UserResponse{}

	err = util.ConvertTypes(s.logger, survey, &sResponse)

	if err != nil {
		return nil, err
	}

	return &sResponse, nil
}

func (s *UserService) SetMaxSurveys(ctx context.Context, userID string, maxSurveys int) error {
	if err := s.repo.UpdateMaxSurveys(ctx, userID, maxSurveys); err != nil {
		s.logger.Error(logging.Internal, logging.Update, "failed to update max surveys", map[logging.ExtraKey]interface{}{
			logging.Service:      "UserService",
			logging.ErrorMessage: err.Error(),
		})
		return err
	}
	return nil
}

// Deposit money to user's wallet
func (s *UserService) Deposit(ctx context.Context, userID uint, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.WalletBalance += amount

	_, err = s.repo.UpdateUser(ctx, user)
	return err
}

// Withdraw money from user's wallet
func (s *UserService) Withdraw(ctx context.Context, userID uint, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.WalletBalance < amount {
		return errors.New("insufficient balance")
	}

	user.WalletBalance -= amount

	_, err = s.repo.UpdateUser(ctx, user)
	return err
}

func (s *UserService) Transfer(ctx context.Context, senderID, receiverID uint, amount float64) error {
	if senderID == receiverID {
		return errors.New("cannot transfer to the same user")
	}
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	sender, err := s.repo.GetUserByID(ctx, senderID)
	if err != nil {
		return err
	}

	if sender.WalletBalance < amount {
		return errors.New("insufficient balance")
	}

	receiver, err := s.repo.GetUserByID(ctx, receiverID)
	if err != nil {
		return err
	}

	sender.WalletBalance -= amount
	receiver.WalletBalance += amount

	// Update both users
	_, err = s.repo.UpdateUser(ctx, sender)
	if err != nil {
		return err
	}

	_, err = s.repo.UpdateUser(ctx, receiver)
	if err != nil {
		return err
	}

	transaction := &models.Transaction{
		BuyerID:  receiverID,
		SellerID: senderID,
		Amount:   amount,
	}
	return s.repo.CreateTransaction(ctx, transaction)
}
