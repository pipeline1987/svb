package database

import (
	"context"
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
	"github.com/pipeline1987/SVB/models"
)

type PsqlRepository struct {
	db *sql.DB
}

func NewPsqlRepository(url string) (*PsqlRepository, error) {
	db, instanceError := sql.Open("postgres", url)

	if instanceError != nil {
		return nil, instanceError
	}

	return &PsqlRepository{db}, nil
}

func (repo *PsqlRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	existingUser, existingError := repo.db.QueryContext(
		ctx,
		"SELECT id FROM users WHERE email = $1",
		user.Email,
	)

	var preSavedUser = models.User{}

	for existingUser.Next() {
		if existingError = existingUser.Scan(&preSavedUser.Id); existingError == nil {
			return nil, errors.New("there are a user with this email")
		}
	}

	_, insertError := repo.db.ExecContext(
		ctx,
		"INSERT INTO users (id, email, full_name, password) VALUES ($1, $2, $3, $4)",
		user.Id, user.Email, user.FullName, user.Password,
	)

	if insertError != nil {
		return nil, insertError
	}

	result, getError := repo.db.QueryContext(
		ctx,
		"SELECT id FROM users WHERE id = $1",
		user.Id,
	)

	var savedUser = models.User{}

	for result.Next() {
		if getError = result.Scan(&savedUser.Id); getError == nil {
			return &savedUser, nil
		}
	}

	if getError = result.Err(); getError != nil {
		return nil, getError
	}

	defer func() {
		getError = result.Close()

		if getError != nil {
			log.Fatal(getError)
		}
	}()

	return &savedUser, nil
}

func (repo PsqlRepository) ReadUser(ctx context.Context, id string) (*models.User, error) {
	result, getError := repo.db.QueryContext(
		ctx,
		"SELECT id, email, full_name FROM users WHERE id = $1",
		id,
	)

	var user = models.User{}

	for result.Next() {
		if getError = result.Scan(&user.Id, &user.Email, &user.FullName); getError == nil {
			return &user, nil
		}
	}

	if getError = result.Err(); getError != nil {
		return nil, getError
	}

	defer func() {
		getError = result.Close()

		if getError != nil {
			log.Fatal(getError)
		}
	}()

	return &user, nil
}

func (repo PsqlRepository) ReadUserByEmail(ctx context.Context, email string) (*models.User, error) {
	result, getError := repo.db.QueryContext(
		ctx,
		"SELECT id, email, password FROM users WHERE email = $1",
		email,
	)

	var user = models.User{}

	for result.Next() {
		if getError = result.Scan(&user.Id, &user.Email, &user.Password); getError == nil {
			return &user, nil
		}
	}

	if getError = result.Err(); getError != nil {
		return nil, getError
	}

	defer func() {
		getError = result.Close()

		if getError != nil {
			log.Fatal(getError)
		}
	}()

	return &user, nil
}

func (repo *PsqlRepository) CreateBankAccount(
	ctx context.Context,
	bankAccount *models.BankAccount,
) (*models.BankAccount, error) {
	existingBankAccount, existingError := repo.db.QueryContext(
		ctx,
		"SELECT id FROM bank_accounts WHERE user_id = $1 AND name = $2",
		bankAccount.UserId,
		bankAccount.Name,
	)

	var preSavedBankAccount = models.BankAccount{}

	for existingBankAccount.Next() {
		if existingError = existingBankAccount.Scan(&preSavedBankAccount.Id); existingError == nil {
			return nil, errors.New("there are a bank account with this user_id and name")
		}
	}

	_, insertError := repo.db.ExecContext(
		ctx,
		"INSERT INTO bank_accounts (id, user_id, name, balance, state) VALUES ($1, $2, $3, $4, $5)",
		bankAccount.Id, bankAccount.UserId, bankAccount.Name, 0, "active",
	)

	if insertError != nil {
		return nil, insertError
	}

	result, getError := repo.db.QueryContext(
		ctx,
		"SELECT id FROM bank_accounts WHERE id = $1",
		bankAccount.Id,
	)

	var savedBankAccount = models.BankAccount{}

	for result.Next() {
		if getError = result.Scan(&savedBankAccount.Id); getError == nil {
			return &savedBankAccount, nil
		}
	}

	if getError = result.Err(); getError != nil {
		return nil, getError
	}

	defer func() {
		getError = result.Close()

		if getError != nil {
			log.Fatal(getError)
		}
	}()

	return &savedBankAccount, nil
}

func (repo PsqlRepository) GetBankAccountById(ctx context.Context, id string, userId string) (*models.BankAccount, error) {
	result, getError := repo.db.QueryContext(
		ctx,
		"SELECT id, name, balance, state FROM bank_accounts WHERE id = $1 AND user_id = $2",
		id,
		userId,
	)

	var bankAccount = models.BankAccount{}

	for result.Next() {
		if getError = result.Scan(
			&bankAccount.Id,
			&bankAccount.Name,
			&bankAccount.Balance,
			&bankAccount.State,
		); getError == nil {
			return &bankAccount, nil
		}
	}

	if getError = result.Err(); getError != nil {
		return nil, getError
	}

	defer func() {
		getError = result.Close()

		if getError != nil {
			log.Fatal(getError)
		}
	}()

	return &bankAccount, nil
}

func (repo PsqlRepository) UpdateBankAccountById(
	ctx context.Context,
	id string,
	userId string,
	bankAccount *models.BankAccount,
) (*models.BankAccount, error) {
	execResult, execErr := repo.db.ExecContext(ctx,
		"UPDATE bank_accounts SET name = $1, state = $2 WHERE id = $3 AND user_id = $4",
		bankAccount.Name,
		bankAccount.State,
		id,
		userId)

	if execErr != nil {
		return nil, execErr
	}

	n, err := execResult.RowsAffected()

	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, sql.ErrNoRows
	}

	result, getError := repo.db.QueryContext(
		ctx,
		"SELECT id, name, balance, state FROM bank_accounts WHERE id = $1",
		id,
	)

	var updatedBankAccount = models.BankAccount{}

	for result.Next() {
		if getError = result.Scan(
			&updatedBankAccount.Id,
			&updatedBankAccount.Name,
			&updatedBankAccount.Balance,
			&updatedBankAccount.State,
		); getError == nil {
			return &updatedBankAccount, nil
		}
	}

	defer func() {
		getError = result.Close()

		if getError != nil {
			log.Fatal(getError)
		}
	}()

	return &updatedBankAccount, nil
}

func (repo PsqlRepository) DeleteBankAccountById(ctx context.Context, id string, userId string) error {
	execResult, execErr := repo.db.ExecContext(
		ctx,
		"DELETE FROM bank_accounts WHERE id = $1 AND user_id = $2",
		id,
		userId,
	)

	if execErr != nil {
		return execErr
	}

	n, err := execResult.RowsAffected()

	if err != nil {
		return err
	}

	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (repo PsqlRepository) GetAllBankAccountsByUserId(ctx context.Context, userId string) ([]*models.BankAccount, error) {
	result, getError := repo.db.QueryContext(
		ctx,
		"SELECT id, name, balance, state FROM bank_accounts WHERE user_id = $1",
		userId,
	)

	defer func() {
		getError = result.Close()

		if getError != nil {
			log.Fatal(getError)
		}
	}()

	var bankAccounts []*models.BankAccount

	for result.Next() {
		var bankAccount = models.BankAccount{}

		if getError = result.Scan(
			&bankAccount.Id,
			&bankAccount.Name,
			&bankAccount.Balance,
			&bankAccount.State,
		); getError == nil {
			bankAccounts = append(bankAccounts, &bankAccount)
		}
	}

	if getError = result.Err(); getError != nil {
		return nil, getError
	}

	return bankAccounts, nil
}

func (repo *PsqlRepository) Close() error {
	return repo.db.Close()
}
