package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"tender/internal/models"
)

type UserRepository struct {
	Db *sql.DB
}

var (
	ErrDuplicateEmail = errors.New("Пользователь с таким адресом электронной почты уже существует")
	ErrDuplicatePhone = errors.New("Пользователь с таким номером телефона уже существует")
	ErrNotFound       = func(errorMessage string) error {
		return errors.New(fmt.Sprintf("no client found with the given %s", errorMessage))
	}
)

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {

	rows, err := r.Db.QueryContext(ctx, "SELECT id, name, last_name, email, phone, inn, balance, password, status FROM users WHERE id != 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.LastName, &user.Email, &user.Phone, &user.INN, &user.Balance, &user.Password, &user.Status); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) SignUp(ctx context.Context, user models.User) (models.User, error) {
	var exists int
	emailCheckQuery := "SELECT COUNT(*) FROM users WHERE email= ?"
	phoneCheckQuery := "SELECT COUNT(*) FROM users WHERE phone IS NOT NULL AND phone = ? "

	err := r.Db.QueryRow(emailCheckQuery, user.Email).Scan(&exists)
	if err != nil {
		return models.User{}, err
	}
	if exists > 0 && user.Email != "" {
		return models.User{}, ErrDuplicateEmail
	}

	err = r.Db.QueryRow(phoneCheckQuery, user.Phone).Scan(&exists)
	if err != nil {
		return models.User{}, err
	}
	if exists > 0 && user.Phone != "" {
		return models.User{}, ErrDuplicatePhone
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return models.User{}, err
	}

	result, err := r.Db.ExecContext(ctx, "INSERT INTO users(name, last_name, email, phone, inn, password, balance) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.Name, user.LastName, user.Email, user.Phone, user.INN, hashedPassword, 0)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return models.User{}, models.ErrDuplicateEmail
		}
		return models.User{}, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return models.User{}, err
	}

	user.ID = int(userID)
	user.Balance = 0

	return user, nil
}

func (r *UserRepository) LogIn(ctx context.Context, user models.User) (models.User, error) {
	var storedUser models.User

	query := "SELECT id, name, last_name, email, phone, inn, password, balance,status FROM users WHERE email = ? OR phone = ?"
	err := r.Db.QueryRowContext(ctx, query, user.Email, user.Phone).Scan(
		&storedUser.ID,
		&storedUser.Name,
		&storedUser.LastName,
		&storedUser.Email,
		&storedUser.Phone,
		&storedUser.INN,
		&storedUser.Password,
		&storedUser.Balance,
		&storedUser.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return models.User{}, models.ErrInvalidPassword
		}
		return models.User{}, err
	}

	if storedUser.Status == 2 {
		return models.User{}, errors.New("Пользователь заблокирован")
	}

	return storedUser, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (models.User, error) {
	var user models.User

	query := "SELECT id, name, last_name, email, phone, inn, balance, password, status FROM users WHERE id = ?"
	err := r.Db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.INN,
		&user.Balance,
		&user.Password,
		&user.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepository) UpdateBalance(ctx context.Context, id int, amount float64) error {
	if id != 1 {
		_, err := r.Db.ExecContext(ctx, "UPDATE users SET balance = balance + ? WHERE id = ?", amount, id)
		_, err = r.Db.ExecContext(ctx, "UPDATE users SET balance = balance - ? WHERE id = ?", amount, 1)
		return err
	} else {
		_, err := r.Db.ExecContext(ctx, "UPDATE users SET balance = balance + ? WHERE id = ?", amount, 1)
		return err
	}

}

func (r *UserRepository) GetBalance(ctx context.Context, id int) (float64, error) {
	var balance float64

	query := "SELECT balance FROM users WHERE id = ?"
	err := r.Db.QueryRowContext(ctx, query, id).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrUserNotFound
		}
		return 0, err
	}

	return balance, nil
}

// DeleteUserByID removes a user from the database by ID.
func (r *UserRepository) DeleteUserByID(ctx context.Context, id int) error {

	_, err := r.Db.ExecContext(ctx, "UPDATE transactions SET user_id = NULL WHERE user_id = ?", id)
	if err != nil {
		return err
	}

	result, err := r.Db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	query := "UPDATE users SET"
	params := []interface{}{}

	if user.Name != "" {
		query += " name = ?,"
		params = append(params, user.Name)
	}
	if user.LastName != "" {
		query += " last_name = ?,"
		params = append(params, user.LastName)
	}
	if user.Email != "" {
		query += " email = ?,"
		params = append(params, user.Email)
	}
	if user.Phone != "" {
		query += " phone = ?,"
		params = append(params, user.Phone)
	}
	if user.INN != "" {
		query += " inn = ?,"
		params = append(params, user.INN)
	}
	if user.Balance != 0 {
		query += " balance = ?,"
		params = append(params, user.Balance)
	}
	if user.Password != "" {
		query += " password = ?,"
		params = append(params, user.Password)
	}

	if user.Status != 0 {
		query += " status = ?,"
		params = append(params, user.Status)

	}

	// Trim the last comma from the query
	query = query[:len(query)-1]
	query += " WHERE id = ?"
	params = append(params, user.ID)

	_, err := r.Db.ExecContext(ctx, query, params...)
	if err != nil {
		return models.User{}, err
	}

	// Retrieve the updated user
	row := r.Db.QueryRowContext(ctx, "SELECT id, name, last_name, email, phone, inn, balance, password, status FROM users WHERE id = ?", user.ID)
	var updatedUser models.User
	err = row.Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.LastName, &updatedUser.Email, &updatedUser.Phone, &updatedUser.INN, &updatedUser.Balance, &updatedUser.Password, &updatedUser.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, err
	}

	return updatedUser, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID int, hashedPassword string) error {
	_, err := r.Db.ExecContext(ctx, "UPDATE users SET password = ? WHERE id = ?", hashedPassword, userID)
	return err
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (int, error) {
	stmt := `SELECT id FROM users WHERE email = ?`
	var userId int
	err := r.Db.QueryRowContext(ctx, stmt, email).Scan(&userId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, ErrNotFound(email)
		default:
			return 0, err
		}
	}
	return userId, nil
}

func (r *UserRepository) GetUserTransactionDifferences(ctx context.Context) ([]models.UserTransactionDifference, error) {
	query := `
        SELECT
            u.id,
            u.name,
            u.last_name,
            COALESCE(SUM(CASE WHEN et.status = 2 THEN et.total ELSE 0 END), 0) AS difference
        FROM
            users u
                JOIN
            extra_transactions et ON u.id = et.user_id
        GROUP BY
            u.id, u.name, u.last_name
    `

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.UserTransactionDifference
	for rows.Next() {
		var result models.UserTransactionDifference
		if err := rows.Scan(&result.ID, &result.Name, &result.LastName, &result.Difference); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
