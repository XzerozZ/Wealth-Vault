package usecase

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
	"wealth-vault/auth-service/internal/domain"
	pb "wealth-vault/auth-service/pkg/pb/proto/user"
	tokenmock "wealth-vault/auth-service/test/mock/pkg/token"
	repomock "wealth-vault/auth-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestLINEIDTokenVerifier(t *testing.T) {
	var checkedRequest bool
	verifier := &LINEIDTokenVerifier{client: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.String() != "https://api.line.me/oauth2/v2.1/verify" {
			t.Fatal("unexpected verification URL")
		}
		if err := req.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if req.PostForm.Get("client_id") != "test-channel" || req.PostForm.Get("id_token") != strings.Repeat("x", 30) {
			t.Fatal("incorrect verification form")
		}
		checkedRequest = true
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"iss":"https://access.line.me","aud":"test-channel","sub":"line-subject","exp":` + intToString(time.Now().Add(time.Minute).Unix()) + `,"name":"Tester"}`)), Header: make(http.Header)}, nil
	})}}
	identity, err := verifier.Verify(context.Background(), strings.Repeat("x", 30), "test-channel")
	if err != nil || !checkedRequest || identity.Subject != "line-subject" {
		t.Fatalf("identity=%+v err=%v", identity, err)
	}
	if _, err := verifier.Verify(context.Background(), "short", "test-channel"); err != ErrLineRejected {
		t.Fatalf("short token err=%v", err)
	}
	badAudience := &LINEIDTokenVerifier{client: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"iss":"https://access.line.me","aud":"other-channel","sub":"line-subject","exp":` + intToString(time.Now().Add(time.Minute).Unix()) + `}`)), Header: make(http.Header)}, nil
	})}}
	if _, err := badAudience.Verify(context.Background(), strings.Repeat("x", 30), "test-channel"); err != ErrLineRejected {
		t.Fatalf("wrong audience err=%v", err)
	}
}

func TestLINEAccessTokenVerifier(t *testing.T) {
	accessToken := strings.Repeat("a", 40)
	profileCalled := false
	verifier := &LINEIDTokenVerifier{client: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/oauth2/v2.1/verify":
			if req.Method != http.MethodGet || req.URL.Query().Get("access_token") != accessToken {
				t.Fatal("incorrect access-token verification request")
			}
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"client_id":"test-channel","expires_in":3600}`)), Header: make(http.Header)}, nil
		case "/oauth2/v2.1/userinfo":
			profileCalled = true
			if req.Header.Get("Authorization") != "Bearer "+accessToken {
				t.Fatal("userinfo request did not use verified token")
			}
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"sub":"line-subject","name":"Tester"}`)), Header: make(http.Header)}, nil
		default:
			t.Fatal("unexpected LINE endpoint")
			return nil, nil
		}
	})}}
	identity, err := verifier.Verify(context.Background(), lineAccessTokenPrefix+accessToken, "test-channel")
	if err != nil || !profileCalled || identity.Subject != "line-subject" || identity.DisplayName != "Tester" {
		t.Fatalf("identity=%+v err=%v", identity, err)
	}

	for _, response := range []string{
		`{"client_id":"wrong-channel","expires_in":3600}`,
		`{"client_id":"test-channel","expires_in":0}`,
	} {
		profileCalled = false
		rejected := &LINEIDTokenVerifier{client: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Path == "/oauth2/v2.1/userinfo" {
				profileCalled = true
			}
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(response)), Header: make(http.Header)}, nil
		})}}
		if _, err := rejected.Verify(context.Background(), lineAccessTokenPrefix+accessToken, "test-channel"); err != ErrLineRejected || profileCalled {
			t.Fatalf("invalid access token metadata was accepted: err=%v profileCalled=%v", err, profileCalled)
		}
	}
}

func intToString(value int64) string { return strconv.FormatInt(value, 10) }

type lineVerifierFunc func(context.Context, string, string) (*VerifiedLineIdentity, error)

func (f lineVerifierFunc) Verify(ctx context.Context, token, channelID string) (*VerifiedLineIdentity, error) {
	return f(ctx, token, channelID)
}

type lineUserClient struct {
	pb.UserServiceClient
	userID uuid.UUID
}

func (c lineUserClient) CreateUser(_ context.Context, request *pb.CreateUserRequest, _ ...grpc.CallOption) (*pb.UserResponse, error) {
	return &pb.UserResponse{Success: true, User: &pb.User{Id: c.userID.String(), Email: request.Email, Username: request.Username}}, nil
}

func (c lineUserClient) GetUser(_ context.Context, _ *pb.GetUserByIDRequest, _ ...grpc.CallOption) (*pb.UserResponse, error) {
	return &pb.UserResponse{Success: true, User: &pb.User{Id: c.userID.String(), Username: "Local Tester"}}, nil
}

func TestExchangeLineCreatesShortLivedOpaqueChallenge(t *testing.T) {
	ctx := context.Background()
	repo := new(repomock.MockAuthRepository)
	repo.On("FindLineIdentity", ctx, "test-channel", "line-subject").Return((*domain.LineIdentity)(nil), nil).Once()
	repo.On("SaveLineChallenge", ctx, mock.MatchedBy(func(challenge *domain.LineChallenge) bool {
		return len(challenge.Digest) == 64 && challenge.ChannelID == "test-channel" && challenge.Subject == "line-subject" &&
			challenge.ExpiresAt.After(time.Now().Add(4*time.Minute)) && challenge.ExpiresAt.Before(time.Now().Add(6*time.Minute))
	})).Return(nil).Once()
	repo.On("FindLocalAccountByEmail", ctx, "tester@example.invalid").Return((*domain.AuthAccount)(nil), nil).Once()

	uc := NewAuthUsecase(repo, nil, nil, nil, nil)
	uc.ConfigureLINE("test-channel", lineVerifierFunc(func(_ context.Context, token, channelID string) (*VerifiedLineIdentity, error) {
		if token != "raw-id-token" || channelID != "test-channel" {
			t.Fatal("unexpected verifier input")
		}
		return &VerifiedLineIdentity{Subject: "line-subject", Email: "tester@example.invalid"}, nil
	}))
	result, err := uc.ExchangeLine(ctx, "raw-id-token", "test-channel")
	if err != nil || result.Status != "onboarding_required" || len(result.ChallengeID) != 64 {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	repo.AssertExpectations(t)
}

func TestCompleteLineChallengeRegistersAndLinksNewUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	challengeID := strings.Repeat("a", 64)
	repo := new(repomock.MockAuthRepository)
	repo.On("ClaimLineChallenge", ctx, challengeDigest(challengeID)).Return(&domain.LineChallenge{ChannelID: "test-channel", Subject: "line-subject"}, nil).Once()
	repo.On("FindLocalAccountByEmail", ctx, "tester@example.invalid").Return((*domain.AuthAccount)(nil), nil).Once()
	repo.On("Register", ctx, mock.MatchedBy(func(account *domain.AuthAccount) bool {
		return account.UserID == userID && account.Email == "tester@example.invalid" && account.Provider == ProviderLocal && account.Password != "Password123!"
	})).Return(nil).Once()
	repo.On("FindLineIdentity", ctx, "test-channel", "line-subject").Return((*domain.LineIdentity)(nil), nil).Once()
	repo.On("LinkLineIdentity", ctx, mock.MatchedBy(func(identity *domain.LineIdentity) bool {
		return identity.ChannelID == "test-channel" && identity.Subject == "line-subject" && identity.UserID == userID
	})).Return(nil).Once()
	repo.On("FindByUserIDAndProvider", ctx, userID, ProviderLocal).Return(&domain.AuthAccount{UserID: userID, Email: "tester@example.invalid"}, nil).Once()
	repo.On("CreateSession", ctx, mock.AnythingOfType("*domain.AuthSession")).Return(nil).Once()
	tokens := new(tokenmock.MockTokenMaker)
	tokens.On("CreateToken", userID.String(), "tester@example.invalid", TokenTypeAccess, AccessTokenTTL).Return("access-token", nil).Once()
	tokens.On("CreateToken", userID.String(), "tester@example.invalid", TokenTypeRefresh, RefreshTokenTTL).Return("refresh-token", nil).Once()

	uc := NewAuthUsecase(repo, lineUserClient{userID: userID}, tokens, nil, nil)
	uc.ConfigureLINE("test-channel", nil)
	result, err := uc.CompleteLineChallenge(ctx, domain.LineChallengeInput{
		ChallengeID: challengeID, ChannelID: "test-channel", Action: "register", Email: "tester@example.invalid", Password: "Password123!", Username: "local-tester",
	})
	if err != nil || result.Status != "authenticated" || result.UserID != userID.String() || result.RefreshToken != "refresh-token" {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	repo.AssertExpectations(t)
	tokens.AssertExpectations(t)
}

func TestCompleteLineChallengeLinksExistingAccountByPassword(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	challengeID := strings.Repeat("b", 64)
	hash, err := HashPassword("Password123!")
	if err != nil {
		t.Fatal(err)
	}
	account := &domain.AuthAccount{UserID: userID, Email: "existing@example.invalid", Provider: ProviderLocal, Password: hash}
	repo := new(repomock.MockAuthRepository)
	repo.On("ClaimLineChallenge", ctx, challengeDigest(challengeID)).Return(&domain.LineChallenge{ChannelID: "test-channel", Subject: "line-subject"}, nil).Once()
	repo.On("FindLocalAccountByEmail", ctx, account.Email).Return(account, nil).Once()
	repo.On("FindLineIdentity", ctx, "test-channel", "line-subject").Return((*domain.LineIdentity)(nil), nil).Once()
	repo.On("LinkLineIdentity", ctx, mock.MatchedBy(func(identity *domain.LineIdentity) bool { return identity.UserID == userID })).Return(nil).Once()
	repo.On("FindByUserIDAndProvider", ctx, userID, ProviderLocal).Return(account, nil).Once()
	repo.On("CreateSession", ctx, mock.AnythingOfType("*domain.AuthSession")).Return(nil).Once()
	tokens := new(tokenmock.MockTokenMaker)
	tokens.On("CreateToken", userID.String(), account.Email, TokenTypeAccess, AccessTokenTTL).Return("access-token", nil).Once()
	tokens.On("CreateToken", userID.String(), account.Email, TokenTypeRefresh, RefreshTokenTTL).Return("refresh-token", nil).Once()

	uc := NewAuthUsecase(repo, lineUserClient{userID: userID}, tokens, nil, nil)
	uc.ConfigureLINE("test-channel", nil)
	result, err := uc.CompleteLineChallenge(ctx, domain.LineChallengeInput{
		ChallengeID: challengeID, ChannelID: "test-channel", Action: "link", Email: account.Email, Password: "Password123!",
	})
	if err != nil || result.Status != "authenticated" || result.UserID != userID.String() {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	repo.AssertExpectations(t)
	tokens.AssertExpectations(t)
}
