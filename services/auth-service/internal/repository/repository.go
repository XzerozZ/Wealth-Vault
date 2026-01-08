package repository

import (
	"context"
	"wealth-vault/auth-service/internal/domain"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Register(ctx context.Context, auth *domain.AuthAccount) error {
	if err := r.db.WithContext(ctx).Create(&auth).Error; err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*domain.AuthAccount, error) {
	var auth domain.AuthAccount
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&auth).Error; err != nil {
		return nil, err
	}
	return &auth, nil
}

func (r *AuthRepository) CreateSession(ctx context.Context, session *domain.AuthSession) error {
	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*domain.AuthSession, error) {
	var session domain.AuthSession
	err := r.db.WithContext(ctx).
		Where("refresh_token = ? AND revoked = ?", refreshToken, false).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *AuthRepository) RevokeSession(ctx context.Context, refreshToken string) error {
	err := r.db.WithContext(ctx).
		Model(&domain.AuthSession{}).
		Where("refresh_token = ?", refreshToken).
		Update("revoked", true).Error
	if err != nil {
		return err
	}
	return nil
}
