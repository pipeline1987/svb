package repositories

import (
	"context"

	"github.com/pipeline1987/SVB/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	ReadUser(ctx context.Context, id string) (*models.User, error)
	ReadUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateBankAccount(ctx context.Context, bankAccount *models.BankAccount) (*models.BankAccount, error)
	GetBankAccountById(ctx context.Context, id string, userId string) (*models.BankAccount, error)
	UpdateBankAccountById(
		ctx context.Context,
		id string,
		userId string,
		bankAccount *models.BankAccount,
	) (*models.BankAccount, error)
	DeleteBankAccountById(ctx context.Context, id string, userId string) error
	GetAllBankAccountsByUserId(ctx context.Context, userId string) ([]*models.BankAccount, error)
	Close() error
}

var implementation Repository

func SetRepository(repository Repository) {
	implementation = repository
}

func CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return implementation.CreateUser(ctx, user)
}

func ReadUser(ctx context.Context, id string) (*models.User, error) {
	return implementation.ReadUser(ctx, id)
}

func ReadUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementation.ReadUserByEmail(ctx, email)
}

func CreateBankAccount(ctx context.Context, bankAccount *models.BankAccount) (*models.BankAccount, error) {
	return implementation.CreateBankAccount(ctx, bankAccount)
}

func GetBankAccountById(ctx context.Context, id string, userId string) (*models.BankAccount, error) {
	return implementation.GetBankAccountById(ctx, id, userId)
}

func UpdateBankAccountById(ctx context.Context, id string, userId string, bankAccount *models.BankAccount) (*models.BankAccount, error) {
	return implementation.UpdateBankAccountById(ctx, id, userId, bankAccount)
}

func DeleteBankAccountById(ctx context.Context, id string, userId string) error {
	return implementation.DeleteBankAccountById(ctx, id, userId)
}

func GetAllBankAccountsByUserId(ctx context.Context, userId string) ([]*models.BankAccount, error) {
	return implementation.GetAllBankAccountsByUserId(ctx, userId)
}

func Close() error {
	return implementation.Close()
}
