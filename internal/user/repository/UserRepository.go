package repository

import (
	"context"
	"errors"

	"github.com/G9QBootcamp/qoli-survey/internal/db"
	surveyModels "github.com/G9QBootcamp/qoli-survey/internal/survey/models"
	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUsers(ctx context.Context, filters dto.UserFilters) ([]models.User, error)
	GetUserByID(ctx context.Context, userID uint) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id uint) error
	IsEmailOrNationalIDTaken(ctx context.Context, email, nationalID string) bool
	UpdateMaxSurveys(ctx context.Context, userID string, maxSurveys int) error
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserCount(ctx context.Context) (int64, error)
	GetRoleByName(ctx context.Context, roleName string) (*models.Role, error)
	CreateTransaction(ctx context.Context, transaction *models.Transaction) error
	GetBalance(ctx context.Context, userID uint) (float64, error)
	Deposit(ctx context.Context, userID uint, amount float64) error
	Withdraw(ctx context.Context, userID uint, amount float64) error
	Transfer(ctx context.Context, senderID, receiverID uint, amount float64) error
	GetVoterID(ctx context.Context, voteID uint) (uint, error)
	UpdateVoteVoter(ctx context.Context, voterID uint, voteID uint) error
}

type UserRepository struct {
	db     db.DbService
	logger logging.Logger
}

func NewUserRepository(db db.DbService, logger logging.Logger) *UserRepository {
	return &UserRepository{db: db, logger: logger}
}

func (r *UserRepository) GetUsers(ctx context.Context, filters dto.UserFilters) ([]models.User, error) {
	var users []models.User

	query := r.db.GetDb().WithContext(ctx).Model(&models.User{})

	if filters.Name != "" {
		query = query.Where("CONCAT(first_name, ' ', last_name) ILIKE ?", "%"+filters.Name+"%")
	}
	if filters.Email != "" {
		query = query.Where("email ILIKE ?", "%"+filters.Email+"%")
	}

	if filters.NationalID != "" {
		query = query.Where("national_id ILIKE ?", "%"+filters.NationalID+"%")
	}

	if filters.City != "" {
		query = query.Where("city ILIKE ?", "%"+filters.City+"%")
	}

	if filters.YearOfBirth > 0 {
		query = query.Where("EXTRACT(YEAR FROM date_of_birth) = ?", filters.YearOfBirth)
	}
	query = query.Where("deleted_at is null")

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}

	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}
	if err := query.Find(&users).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return []models.User{}, nil
		}
		r.logger.Error(logging.Database, logging.Select, "get users error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})

		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User

	err := r.db.GetDb().WithContext(ctx).First(&user, userID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.db.GetDb().WithContext(ctx).Preload("GlobalRole").Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	user.WalletBalance = 100
	err := r.db.GetDb().WithContext(ctx).Create(&user).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Insert, "create user error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}

	return user, err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uint) error {
	return r.db.GetDb().WithContext(ctx).Where("ID = ?", id).Delete(&models.User{}).Error
}

func (r *UserRepository) IsEmailOrNationalIDTaken(ctx context.Context, email, nationalID string) bool {
	var user models.User
	err := r.db.GetDb().WithContext(ctx).Where("email = ? OR national_id = ?", email, nationalID).First(&user).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (r *UserRepository) UpdateMaxSurveys(ctx context.Context, userID string, maxSurveys int) error {
	err := r.db.GetDb().WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("max_surveys", maxSurveys).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Update, "failed to update max surveys", map[logging.ExtraKey]interface{}{
			logging.ErrorMessage: err.Error(),
		})
	}
	return err
}
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.GetDb().WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.GetDb().WithContext(ctx).Model(&models.User{}).Count(&count).Error
	if err != nil {
		if err != nil {
			r.logger.Error(logging.Database, logging.Select, "get user count error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		}
		return 0, err
	}
	return count, nil
}

func (r *UserRepository) GetRoleByName(ctx context.Context, roleName string) (*models.Role, error) {
	var role models.Role
	err := r.db.GetDb().WithContext(ctx).Where("name = ?", roleName).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// new transaction
func (r *UserRepository) CreateTransaction(ctx context.Context, transaction *models.Transaction) error {
	return r.db.GetDb().WithContext(ctx).Create(transaction).Error
}

// finding balance using transaction table
func (r *UserRepository) GetBalance(ctx context.Context, userID uint) (float64, error) {
	var result struct {
		Total float64
	}
	err := r.db.GetDb().WithContext(ctx).
		Model(&models.Transaction{}).
		Select("SUM(amount) as total").
		Where("user_id = ?", userID).
		Scan(&result).Error
	if err != nil {
		return 0, err
	}
	return result.Total, nil
}

// Deposit - add money to user's wallet
func (r *UserRepository) Deposit(ctx context.Context, userID uint, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		r.logger.Error(logging.Database, logging.Update, "deposit failed to find user", map[logging.ExtraKey]interface{}{
			logging.ErrorMessage: err.Error(),
		})
		return err
	}

	user.WalletBalance += amount

	if err := r.db.GetDb().WithContext(ctx).Save(user).Error; err != nil {
		r.logger.Error(logging.Database, logging.Update, "deposit failed to update wallet balance", map[logging.ExtraKey]interface{}{
			logging.ErrorMessage: err.Error(),
		})
		return err
	}

	transaction := &models.Transaction{
		UserID: userID,
		Amount: amount,
		Reason: "deposit",
	}
	r.CreateTransaction(ctx, transaction)

	return nil
}

// Withdraw - subtract money from user's wallet
func (r *UserRepository) Withdraw(ctx context.Context, userID uint, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		r.logger.Error(logging.Database, logging.Update, "withdraw failed to find user", map[logging.ExtraKey]interface{}{
			logging.ErrorMessage: err.Error(),
		})
		return err
	}

	if user.WalletBalance < amount {
		return errors.New("insufficient balance")
	}

	user.WalletBalance -= amount

	if err := r.db.GetDb().WithContext(ctx).Save(user).Error; err != nil {
		r.logger.Error(logging.Database, logging.Update, "withdraw failed to update wallet balance", map[logging.ExtraKey]interface{}{
			logging.ErrorMessage: err.Error(),
		})
		return err
	}

	transaction := &models.Transaction{
		UserID: userID,
		Amount: -amount,
		Reason: "withdraw",
	}
	r.CreateTransaction(ctx, transaction)

	return nil
}

// Transfer - transfer money from sender to receiver
func (r *UserRepository) Transfer(ctx context.Context, senderID, receiverID uint, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	if senderID == receiverID {
		return errors.New("cannot transfer to the same user")
	}
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	sender, err := r.GetUserByID(ctx, senderID)
	if err != nil {
		return err
	}

	receiver, err := r.GetUserByID(ctx, receiverID)
	if err != nil {
		return err
	}

	if sender.WalletBalance < amount {
		return errors.New("insufficient balance")
	}

	sender.WalletBalance -= amount
	receiver.WalletBalance += amount

	if err := r.db.GetDb().WithContext(ctx).Save(sender).Error; err != nil {
		return err
	}
	if err := r.db.GetDb().WithContext(ctx).Save(receiver).Error; err != nil {
		return err
	}

	transaction := &models.Transaction{
		UserID: senderID,
		Amount: -amount,
		Reason: "Transfer",
	}
	r.CreateTransaction(ctx, transaction)

	transaction = &models.Transaction{
		UserID: receiverID,
		Amount: amount,
		Reason: "Transfer",
	}
	r.CreateTransaction(ctx, transaction)

	return nil
}

func (r *UserRepository) GetVoterID(ctx context.Context, voteID uint) (uint, error) {
	var vote surveyModels.Vote
	err := r.db.GetDb().WithContext(ctx).First(&vote, voteID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	return vote.VoterID, err
}

func (r *UserRepository) UpdateVoteVoter(ctx context.Context, voterID uint, voteID uint) error {
	err := r.db.GetDb().WithContext(ctx).Model(&surveyModels.Vote{}).Where("id = ?", voteID).Update("voter_id", voterID).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Update, "failed to update max surveys", map[logging.ExtraKey]interface{}{
			logging.ErrorMessage: err.Error(),
		})
	}
	return err
}
