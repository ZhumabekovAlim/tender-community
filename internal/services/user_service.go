package services

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"tender/internal/models"
	"tender/internal/repositories"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.Repo.GetAllUsers(ctx)
}

func (s *UserService) SignUp(ctx context.Context, user models.User) (models.User, error) {
	return s.Repo.SignUp(ctx, user)
}

func (s *UserService) LogIn(ctx context.Context, user models.User) (models.User, error) {
	return s.Repo.LogIn(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (models.User, error) {
	return s.Repo.GetUserByID(ctx, id)
}

func (s *UserService) UpdateBalance(ctx context.Context, id int, amount float64) error {
	return s.Repo.UpdateBalance(ctx, id, amount)
}

func (s *UserService) GetBalance(ctx context.Context, id int) (float64, error) {
	return s.Repo.GetBalance(ctx, id)
}

func (s *UserService) DeleteUserByID(ctx context.Context, id int) error {
	return s.Repo.DeleteUserByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	if user.Password != "" {
		newPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			return models.User{}, err
		}
		user.Password = string(newPass)
	}
	return s.Repo.UpdateUser(ctx, user)
}

func (s *UserService) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	user, err := s.Repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return models.ErrInvalidPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	return s.Repo.UpdatePassword(ctx, userID, string(hashedPassword))
}

func (s *UserService) FindUserByEmail(ctx context.Context, email string) (int, error) {
	return s.Repo.FindUserByEmail(ctx, email)
}
