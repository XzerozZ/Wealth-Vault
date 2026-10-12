package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"time"
	"wealth-vault/auth-service/internal/domain"
	pb "wealth-vault/auth-service/pkg/pb/proto/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var (
	ErrLineRejected      = errors.New("LINE authentication rejected")
	ErrLineNotConfigured = errors.New("LINE channel is not configured")
	ErrLineUnavailable   = errors.New("LINE verification unavailable")
)

type VerifiedLineIdentity struct {
	Subject     string
	Email       string
	DisplayName string
}

type LineVerifier interface {
	Verify(ctx context.Context, rawIDToken, expectedChannelID string) (*VerifiedLineIdentity, error)
}

type LINEIDTokenVerifier struct{ client *http.Client }

const lineAccessTokenPrefix = "line_access_token:"

func NewLINEIDTokenVerifier() *LINEIDTokenVerifier {
	return &LINEIDTokenVerifier{client: &http.Client{Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
}

func (v *LINEIDTokenVerifier) Verify(ctx context.Context, rawIDToken, expectedChannelID string) (*VerifiedLineIdentity, error) {
	if strings.HasPrefix(rawIDToken, lineAccessTokenPrefix) {
		return v.verifyAccessToken(ctx, strings.TrimPrefix(rawIDToken, lineAccessTokenPrefix), expectedChannelID)
	}
	if len(rawIDToken) < 20 || len(rawIDToken) > 6000 || expectedChannelID == "" {
		return nil, ErrLineRejected
	}
	form := url.Values{"id_token": {rawIDToken}, "client_id": {expectedChannelID}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.line.me/oauth2/v2.1/verify", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, ErrLineUnavailable
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := v.client.Do(req)
	if err != nil {
		return nil, ErrLineUnavailable
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusBadRequest || response.StatusCode == http.StatusUnauthorized {
		return nil, ErrLineRejected
	}
	if response.StatusCode != http.StatusOK {
		return nil, ErrLineUnavailable
	}
	var claims struct {
		Issuer    string `json:"iss"`
		Audience  string `json:"aud"`
		Subject   string `json:"sub"`
		ExpiresAt int64  `json:"exp"`
		Email     string `json:"email"`
		Name      string `json:"name"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 16_384)).Decode(&claims); err != nil {
		return nil, ErrLineUnavailable
	}
	if claims.Issuer != "https://access.line.me" || claims.Audience != expectedChannelID || claims.Subject == "" || claims.ExpiresAt <= time.Now().Unix() {
		return nil, ErrLineRejected
	}
	return &VerifiedLineIdentity{Subject: claims.Subject, Email: strings.ToLower(strings.TrimSpace(claims.Email)), DisplayName: strings.TrimSpace(claims.Name)}, nil
}

// LIFF may retain an expired ID token when reopened from LINE's recent tabs.
// A current LINE access token is independently verified before reading userinfo.
func (v *LINEIDTokenVerifier) verifyAccessToken(ctx context.Context, accessToken, expectedChannelID string) (*VerifiedLineIdentity, error) {
	if len(accessToken) < 20 || len(accessToken) > 6000 || expectedChannelID == "" {
		return nil, ErrLineRejected
	}
	verifyURL := "https://api.line.me/oauth2/v2.1/verify?access_token=" + url.QueryEscape(accessToken)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, verifyURL, nil)
	if err != nil {
		return nil, ErrLineUnavailable
	}
	response, err := v.client.Do(request)
	if err != nil {
		return nil, ErrLineUnavailable
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusBadRequest || response.StatusCode == http.StatusUnauthorized {
		return nil, ErrLineRejected
	}
	if response.StatusCode != http.StatusOK {
		return nil, ErrLineUnavailable
	}
	var verified struct {
		ChannelID string `json:"client_id"`
		ExpiresIn int64  `json:"expires_in"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 16_384)).Decode(&verified); err != nil {
		return nil, ErrLineUnavailable
	}
	if verified.ChannelID != expectedChannelID || verified.ExpiresIn <= 0 {
		return nil, ErrLineRejected
	}

	profileRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.line.me/oauth2/v2.1/userinfo", nil)
	if err != nil {
		return nil, ErrLineUnavailable
	}
	profileRequest.Header.Set("Authorization", "Bearer "+accessToken)
	profileResponse, err := v.client.Do(profileRequest)
	if err != nil {
		return nil, ErrLineUnavailable
	}
	defer profileResponse.Body.Close()
	if profileResponse.StatusCode == http.StatusBadRequest || profileResponse.StatusCode == http.StatusUnauthorized || profileResponse.StatusCode == http.StatusForbidden {
		return nil, ErrLineRejected
	}
	if profileResponse.StatusCode != http.StatusOK {
		return nil, ErrLineUnavailable
	}
	var profile struct {
		Subject string `json:"sub"`
		Name    string `json:"name"`
	}
	if err := json.NewDecoder(io.LimitReader(profileResponse.Body, 16_384)).Decode(&profile); err != nil {
		return nil, ErrLineUnavailable
	}
	if profile.Subject == "" {
		return nil, ErrLineRejected
	}
	return &VerifiedLineIdentity{Subject: profile.Subject, DisplayName: strings.TrimSpace(profile.Name)}, nil
}

func (u *AuthUsecase) ConfigureLINE(channelID string, verifier LineVerifier) {
	u.lineChannelID = strings.TrimSpace(channelID)
	u.lineVerifier = verifier
}

func challengeDigest(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (u *AuthUsecase) ExchangeLine(ctx context.Context, rawIDToken, clientChannelID string) (*domain.LineAuthResult, error) {
	if u.lineChannelID == "" || u.lineVerifier == nil {
		return nil, ErrLineNotConfigured
	}
	if clientChannelID != "" && clientChannelID != u.lineChannelID {
		return nil, ErrLineRejected
	}
	verified, err := u.lineVerifier.Verify(ctx, rawIDToken, u.lineChannelID)
	if err != nil {
		return nil, err
	}
	identity, err := u.authRepo.FindLineIdentity(ctx, u.lineChannelID, verified.Subject)
	if err != nil {
		return nil, err
	}
	if identity != nil {
		return u.lineAuthenticated(ctx, identity.UserID)
	}

	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return nil, err
	}
	challengeID := hex.EncodeToString(secret)
	if err := u.authRepo.SaveLineChallenge(ctx, &domain.LineChallenge{
		Digest: challengeDigest(challengeID), ChannelID: u.lineChannelID, Subject: verified.Subject,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}); err != nil {
		return nil, err
	}
	status := "onboarding_required"
	if verified.Email != "" {
		account, err := u.authRepo.FindLocalAccountByEmail(ctx, verified.Email)
		if err != nil {
			return nil, err
		}
		if account != nil {
			status = "link_required"
		}
	}
	return &domain.LineAuthResult{Status: status, ChallengeID: challengeID}, nil
}

func (u *AuthUsecase) CompleteLineChallenge(ctx context.Context, input domain.LineChallengeInput) (*domain.LineAuthResult, error) {
	if u.lineChannelID == "" {
		return nil, ErrLineNotConfigured
	}
	if input.ChannelID != "" && input.ChannelID != u.lineChannelID {
		return nil, ErrLineRejected
	}
	if len(input.ChallengeID) != 64 || (input.Action != "link" && input.Action != "register") {
		return nil, ErrLineRejected
	}
	email := strings.ToLower(strings.TrimSpace(input.Email))
	parsedEmail, err := mail.ParseAddress(email)
	if err != nil || parsedEmail.Address != email || len(email) > 255 || len(input.Password) < 8 || len(input.Password) > 72 {
		return nil, ErrLineRejected
	}
	digest := challengeDigest(input.ChallengeID)
	challenge, err := u.authRepo.ClaimLineChallenge(ctx, digest)
	if err != nil {
		return nil, ErrLineRejected
	}
	if challenge.ChannelID != u.lineChannelID {
		return nil, ErrLineRejected
	}
	account, err := u.authRepo.FindLocalAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	var userID uuid.UUID
	if input.Action == "link" {
		if account == nil || bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(input.Password)) != nil {
			_ = u.authRepo.ReleaseLineChallenge(ctx, digest)
			return nil, ErrLineRejected
		}
		userID = account.UserID
	} else {
		if account != nil {
			_ = u.authRepo.ReleaseLineChallenge(ctx, digest)
			return &domain.LineAuthResult{Status: "link_required", ChallengeID: input.ChallengeID}, nil
		}
		username := strings.TrimSpace(input.Username)
		if len(username) < 3 || len(username) > 60 {
			_ = u.authRepo.ReleaseLineChallenge(ctx, digest)
			return nil, ErrLineRejected
		}
		hash, err := HashPassword(input.Password)
		if err != nil {
			_ = u.authRepo.ReleaseLineChallenge(ctx, digest)
			return nil, ErrLineRejected
		}
		created, err := u.userClient.CreateUser(ctx, &pb.CreateUserRequest{Email: email, Username: username})
		if err != nil {
			return nil, fmt.Errorf("create LINE user: %w", err)
		}
		if created.GetUser() == nil {
			return nil, errors.New("user service returned empty profile")
		}
		userID, err = uuid.Parse(created.User.Id)
		if err != nil {
			return nil, err
		}
		account = &domain.AuthAccount{UserID: userID, Provider: ProviderLocal, ProviderAccountID: email, Email: email, Password: hash, IsEmailVerified: false}
		if err := u.authRepo.Register(ctx, account); err != nil {
			return nil, fmt.Errorf("save LINE user: %w", err)
		}
		if input.FirstName != "" || input.LastName != "" {
			paths := []string{}
			if input.FirstName != "" {
				paths = append(paths, "first_name")
			}
			if input.LastName != "" {
				paths = append(paths, "last_name")
			}
			_, _ = u.userClient.UpdateUser(ctx, &pb.UpdateUserRequest{Id: userID.String(), Firstname: strings.TrimSpace(input.FirstName), Lastname: strings.TrimSpace(input.LastName), UpdateMask: &fieldmaskpb.FieldMask{Paths: paths}})
		}
	}

	existing, err := u.authRepo.FindLineIdentity(ctx, challenge.ChannelID, challenge.Subject)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		if existing.UserID != userID {
			return nil, ErrLineRejected
		}
	} else if err := u.authRepo.LinkLineIdentity(ctx, &domain.LineIdentity{ChannelID: challenge.ChannelID, Subject: challenge.Subject, UserID: userID}); err != nil {
		return nil, err
	}
	return u.lineAuthenticated(ctx, userID)
}

func (u *AuthUsecase) lineAuthenticated(ctx context.Context, userID uuid.UUID) (*domain.LineAuthResult, error) {
	account, err := u.authRepo.FindByUserIDAndProvider(ctx, userID, ProviderLocal)
	if err != nil {
		return nil, err
	}
	tokens, err := u.GenerateTokensAndSession(ctx, userID, account.Email)
	if err != nil {
		return nil, err
	}
	displayName := account.Email
	profile, err := u.userClient.GetUser(ctx, &pb.GetUserByIDRequest{Id: userID.String()})
	if err == nil && profile.GetUser() != nil {
		if name := strings.TrimSpace(profile.User.Firstname + " " + profile.User.Lastname); name != "" {
			displayName = name
		} else if profile.User.Username != "" {
			displayName = profile.User.Username
		}
	}
	return &domain.LineAuthResult{Status: "authenticated", UserID: userID.String(), DisplayName: displayName,
		AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken, AccessTokenExpiresIn: int32(AccessTokenTTL.Seconds())}, nil
}

func (u *AuthUsecase) RevokeRefreshSession(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return ErrLineRejected
	}
	return u.authRepo.RevokeSession(ctx, refreshToken)
}

func (u *AuthUsecase) DeleteExpiredLineChallenges(ctx context.Context) error {
	return u.authRepo.DeleteExpiredLineChallenges(ctx)
}
