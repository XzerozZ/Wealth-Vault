package repository

import (
	"context"
	"errors"
	"time"
	"wealth-vault/auth-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Register(ctx context.Context, auth *domain.AuthAccount) error {
	if err := r.db.WithContext(ctx).Create(auth).Error; err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) FindByEmailAndProvider(ctx context.Context, email string, provider string) (*domain.AuthAccount, error) {
	var account domain.AuthAccount
	if err := r.db.WithContext(ctx).Where("email = ? AND provider = ?", email, provider).First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AuthRepository) FindByID(ctx context.Context, userid string) (*domain.AuthAccount, error) {
	var auth domain.AuthAccount
	if err := r.db.WithContext(ctx).Where("user_id = ?", userid).First(&auth).Error; err != nil {
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

func (r *AuthRepository) RotateSession(ctx context.Context, oldRefreshToken string, next *domain.AuthSession) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&domain.AuthSession{}).
			Where("refresh_token = ? AND revoked = ? AND refresh_expires_at > ?", oldRefreshToken, false, time.Now()).
			Update("revoked", true)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return errors.New("refresh token was already used")
		}
		return tx.Create(next).Error
	})
}

func (r *AuthRepository) FindLocalAccountByEmail(ctx context.Context, email string) (*domain.AuthAccount, error) {
	var account domain.AuthAccount
	err := r.db.WithContext(ctx).Where("LOWER(email) = ? AND provider = ?", email, "local").First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AuthRepository) FindLineIdentity(ctx context.Context, channelID, subject string) (*domain.LineIdentity, error) {
	var identity domain.LineIdentity
	err := r.db.WithContext(ctx).Where("channel_id = ? AND subject = ?", channelID, subject).First(&identity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *AuthRepository) SaveLineChallenge(ctx context.Context, challenge *domain.LineChallenge) error {
	return r.db.WithContext(ctx).Create(challenge).Error
}

func (r *AuthRepository) ClaimLineChallenge(ctx context.Context, digest string) (*domain.LineChallenge, error) {
	var challenge domain.LineChallenge
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("digest = ?", digest).First(&challenge).Error; err != nil {
			return err
		}
		if challenge.ConsumedAt != nil || !challenge.ExpiresAt.After(time.Now()) || challenge.Attempts >= 5 {
			return errors.New("challenge expired or already used")
		}
		now := time.Now()
		return tx.Model(&challenge).Updates(map[string]interface{}{"consumed_at": now, "attempts": challenge.Attempts + 1}).Error
	})
	if err != nil {
		return nil, err
	}
	return &challenge, nil
}

func (r *AuthRepository) ReleaseLineChallenge(ctx context.Context, digest string) error {
	return r.db.WithContext(ctx).Model(&domain.LineChallenge{}).
		Where("digest = ? AND attempts < ? AND expires_at > ?", digest, 5, time.Now()).
		Update("consumed_at", nil).Error
}

func (r *AuthRepository) LinkLineIdentity(ctx context.Context, identity *domain.LineIdentity) error {
	return r.db.WithContext(ctx).Create(identity).Error
}

func (r *AuthRepository) DeleteExpiredLineChallenges(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&domain.LineChallenge{}).Error
}

func (r *AuthRepository) DeleteExpiredSessions(ctx context.Context) error {
	fourteenDaysAgo := time.Now().AddDate(0, 0, -14)
	err := r.db.WithContext(ctx).
		Where("refresh_expires_at < ?", time.Now()).
		Or("revoked = ? AND updated_at < ?", true, fourteenDaysAgo).
		Delete(&domain.AuthSession{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) DeleteExpiredOTPs(ctx context.Context) error {
	err := r.db.WithContext(ctx).Unscoped().Where("expired_at < ?", time.Now()).Delete(&domain.AuthOTP{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) SaveOTP(ctx context.Context, otp *domain.AuthOTP) error {
	if err := r.db.WithContext(ctx).Create(otp).Error; err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) GetValidOTP(ctx context.Context, userID uuid.UUID, code string) (*domain.AuthOTP, error) {
	var otp domain.AuthOTP
	err := r.db.WithContext(ctx).Where("user_id = ? AND otp = ? AND expired_at > ?", userID, code, time.Now()).First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *AuthRepository) DeleteOTP(ctx context.Context, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).Where("user_id = ?", userID.String()).Delete(&domain.AuthOTP{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, newHash string) error {
	err := r.db.WithContext(ctx).Model(&domain.AuthAccount{}).Where("user_id = ?", userID).Update("password", newHash).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) FindByUserIDAndProvider(ctx context.Context, userID uuid.UUID, provider string) (*domain.AuthAccount, error) {
	var account domain.AuthAccount
	if err := r.db.WithContext(ctx).Where("user_id = ? AND provider = ?", userID, provider).First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AuthRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.AuthAccount, error) {
	var accounts []domain.AuthAccount

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}
