package usecase_test

import (
	"context"
	"errors"
	"testing"
	"wealth-vault/user-service/internal/domain"
	usecase "wealth-vault/user-service/internal/usecase/user"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	mock_repo "wealth-vault/user-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserUsecase_CreateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - create user", func(t *testing.T) {
		mockRepo := new(mock_repo.MockUserRepository)

		uc := usecase.NewUserUsecase(mockRepo, nil, nil, nil, nil)

		req := &pb.CreateUserRequest{
			Email:    "test@example.com",
			Username: "gemini",
		}

		mockRepo.
			On("CreateUser", ctx, mock.AnythingOfType("*domain.User")).
			Return(nil).
			Run(func(args mock.Arguments) {
				user := args.Get(1).(*domain.User)
				user.ID = uuid.New() // simulate DB assigning ID
			}).
			Once()

		res, err := uc.CreateUser(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.True(t, res.Success)
		assert.Equal(t, "gemini", res.User.Username)
		assert.Equal(t, "test@example.com", res.User.Email)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository returns error", func(t *testing.T) {
		mockRepo := new(mock_repo.MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil, nil, nil, nil)

		req := &pb.CreateUserRequest{
			Email:    "fail@example.com",
			Username: "fail",
		}

		mockRepo.
			On("CreateUser", ctx, mock.Anything).
			Return(errors.New("db error")).
			Once()

		res, err := uc.CreateUser(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)

		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_GetUser(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - user found", func(t *testing.T) {
		mockRepo := new(mock_repo.MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil, nil, nil, nil)

		id := uuid.New()

		expectedUser := &domain.User{
			ID:       id,
			Username: "gemini",
			Email:    "gemini@example.com",
		}

		mockRepo.
			On("GetUser", ctx, id).
			Return(expectedUser, nil).
			Once()

		res, err := uc.GetUser(ctx, &pb.GetUserByIDRequest{
			Id: id.String(),
		})

		assert.NoError(t, err)
		assert.True(t, res.Success)
		assert.Equal(t, "gemini", res.User.Username)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		mockRepo := new(mock_repo.MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil, nil, nil, nil)

		res, err := uc.GetUser(ctx, &pb.GetUserByIDRequest{
			Id: "invalid-uuid",
		})

		assert.Error(t, err)
		assert.Nil(t, res)

		mockRepo.AssertNotCalled(t, "GetUser", mock.Anything, mock.Anything)
	})

	t.Run("Repository returns error", func(t *testing.T) {
		mockRepo := new(mock_repo.MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil, nil, nil, nil)

		id := uuid.New()

		mockRepo.
			On("GetUser", ctx, id).
			Return(nil, errors.New("not found")).
			Once()

		res, err := uc.GetUser(ctx, &pb.GetUserByIDRequest{
			Id: id.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, res)

		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_GetUsersByEmail(t *testing.T) {
	ctx := context.Background()
	email := "gemini@example.com"

	t.Run("Success - users found", func(t *testing.T) {
		mockRepo := new(mock_repo.MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil, nil, nil, nil)

		expectedUsers := []*domain.User{
			{
				ID:       uuid.New(),
				Username: "gemini_1",
				Email:    email,
			},
			{
				ID:       uuid.New(),
				Username: "gemini_2",
				Email:    email,
			},
		}

		mockRepo.
			On("GetUsersByEmail", ctx, email).
			Return(expectedUsers, nil).
			Once()

		res, err := uc.GetUsersByEmail(ctx, &pb.GetUserByEmailRequest{
			Email: email,
		})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.True(t, res.Success)
		assert.Len(t, res.User, 2)
		assert.Equal(t, "gemini_1", res.User[0].Username)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository returns error", func(t *testing.T) {
		mockRepo := new(mock_repo.MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil, nil, nil, nil)

		mockRepo.
			On("GetUsersByEmail", ctx, email).
			Return(nil, errors.New("database error")).
			Once()

		res, err := uc.GetUsersByEmail(ctx, &pb.GetUserByEmailRequest{
			Email: email,
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "database error", err.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty result", func(t *testing.T) {
		mockRepo := new(mock_repo.MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil, nil, nil, nil)

		mockRepo.
			On("GetUsersByEmail", ctx, "unknown@example.com").
			Return([]*domain.User{}, nil).
			Once()

		res, err := uc.GetUsersByEmail(ctx, &pb.GetUserByEmailRequest{
			Email: "unknown@example.com",
		})

		assert.NoError(t, err)
		assert.True(t, res.Success)
		assert.Len(t, res.User, 0)

		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_UpdateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - update user", func(t *testing.T) {
		mockRepo := new(mock_repo.MockUserRepository)

		uc := usecase.NewUserUsecase(mockRepo, nil, nil, nil, nil)

		id := uuid.New()

		existingUser := &domain.User{
			ID:       id,
			Username: "old",
			Email:    "old@example.com",
		}

		mockRepo.
			On("GetUser", ctx, id).
			Return(existingUser, nil).
			Once()

		mockRepo.
			On("UpdateUser", ctx, mock.Anything, mock.Anything).
			Return(existingUser, nil).
			Once()

		req := &pb.UpdateUserRequest{
			Id: id.String(),
		}

		res, err := uc.UpdateUser(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.True(t, res.Success)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		mockRepo := new(mock_repo.MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil, nil, nil, nil)

		res, err := uc.UpdateUser(ctx, &pb.UpdateUserRequest{
			Id: "invalid-uuid",
		})

		assert.Error(t, err)
		assert.Nil(t, res)

		mockRepo.AssertNotCalled(t, "GetUser", mock.Anything, mock.Anything)
	})

	t.Run("GetUser returns error", func(t *testing.T) {
		mockRepo := new(mock_repo.MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil, nil, nil, nil)

		id := uuid.New()

		mockRepo.
			On("GetUser", ctx, id).
			Return(nil, errors.New("db error")).
			Once()

		res, err := uc.UpdateUser(ctx, &pb.UpdateUserRequest{
			Id: id.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, res)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateUser returns error", func(t *testing.T) {
		mockRepo := new(mock_repo.MockUserRepository)
		uc := usecase.NewUserUsecase(mockRepo, nil, nil, nil, nil)

		id := uuid.New()

		existingUser := &domain.User{ID: id}

		mockRepo.
			On("GetUser", ctx, id).
			Return(existingUser, nil).
			Once()

		mockRepo.
			On("UpdateUser", ctx, mock.Anything, mock.Anything).
			Return(nil, errors.New("update failed")).
			Once()

		res, err := uc.UpdateUser(ctx, &pb.UpdateUserRequest{
			Id: id.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, res)

		mockRepo.AssertExpectations(t)
	})
}
