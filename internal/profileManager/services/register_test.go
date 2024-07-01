package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ravilock/goduit/internal/profileManager/models"
	"github.com/ravilock/goduit/mocks"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRegister(t *testing.T) {
	userRegister := mocks.NewUserRegister(t)
	hasher := mocks.NewHasher(t)
	service := &registerProfileService{userRegister, hasher}

	t.Run("Should register user successfully", func(t *testing.T) {
		// Arrange
		password := "testing-password"
		hashedPassword := "hashed-password"
		user := generateValidUser()
		user.PasswordHash = &hashedPassword
		*user.ID = primitive.NewObjectID()
		ctx := context.Background()
		hasher.EXPECT().Hash(password).Return(hashedPassword, nil).Once()
		userRegister.EXPECT().RegisterUser(ctx, user).Return(user, nil).Once()

		// Act
		token, err := service.Register(ctx, user, password)

		// Assert
		require.NotZero(t, token)
		require.NoError(t, err)
	})

	t.Run("Should return error if hasher returns error", func(t *testing.T) {
		// Arrange
		password := "testing-password"
		expectedError := errors.New("Password could not be hashed")
		ctx := context.Background()
		hasher.EXPECT().Hash(password).Return("", expectedError).Once()

		// Act
		token, err := service.Register(ctx, nil, password)

		// Assert
		require.Zero(t, token)
		require.ErrorIs(t, err, ErrFailedToGeneratePasswordHash)
		require.ErrorIs(t, err, expectedError)
	})

	t.Run("Should return error if repository fails to persist user", func(t *testing.T) {
		// Arrange
		password := "testing-password"
		hashedPassword := "hashed-password"
		user := generateValidUser()
		user.PasswordHash = &hashedPassword
		expectedError := errors.New("User could not be registered")
		ctx := context.Background()
		hasher.EXPECT().Hash(password).Return(hashedPassword, nil).Once()
		userRegister.EXPECT().RegisterUser(ctx, user).Return(nil, expectedError).Once()

		// Act
		token, err := service.Register(ctx, user, password)

		// Assert
		require.Zero(t, token)
		require.ErrorIs(t, err, ErrFailedToRegisterUser)
		require.ErrorIs(t, err, expectedError)
	})
}

func generateValidUser() *models.User {
	username := "testing-username"
	email := "testing_email@email.com"
	createdAt := time.Now()
	return &models.User{
		ID:           nil,
		Username:     &username,
		Email:        &email,
		PasswordHash: nil,
		Bio:          new(string),
		Image:        new(string),
		CreatedAt:    &createdAt,
		UpdatedAt:    nil,
		LastSession:  &createdAt,
	}
}
