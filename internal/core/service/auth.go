package service

import (
	"context"
	user2 "github.com/GoBootCamp-Group1/Task-Management/internal/core/domain/user"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/jwt"
	"time"

	jwt2 "github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userOps                *user2.Ops
	secret                 []byte
	tokenExpiration        uint
	refreshTokenExpiration uint
}

func NewAuthService(userOps *user2.Ops, secret []byte,
	tokenExpiration uint, refreshTokenExpiration uint) *AuthService {
	return &AuthService{
		userOps:                userOps,
		secret:                 secret,
		tokenExpiration:        tokenExpiration,
		refreshTokenExpiration: refreshTokenExpiration,
	}
}

type UserToken struct {
	AuthorizationToken string
	RefreshToken       string
	ExpiresAt          int64
}

func (s *AuthService) Login(ctx context.Context, email, pass string) (*UserToken, error) {
	user, err := s.userOps.GetUserByEmailAndPassword(ctx, email, pass)
	if err != nil {
		return nil, err
	}

	// calc expiration time values
	var (
		authExp    = time.Now().Add(time.Minute * time.Duration(s.tokenExpiration))
		refreshExp = time.Now().Add(time.Minute * time.Duration(s.refreshTokenExpiration))
	)

	authToken, err := jwt.CreateToken(s.secret, s.userClaims(user, authExp))
	if err != nil {
		return nil, err // todo
	}

	refreshToken, err := jwt.CreateToken(s.secret, s.userClaims(user, refreshExp))
	if err != nil {
		return nil, err // todo
	}

	return &UserToken{
		AuthorizationToken: authToken,
		RefreshToken:       refreshToken,
		ExpiresAt:          authExp.Unix(),
	}, nil
}

func (s *AuthService) RefreshAuth(ctx context.Context, refreshToken string) (*UserToken, error) {
	claim, err := jwt.ParseToken(refreshToken, s.secret)
	if err != nil {
		return nil, err
	}

	u, err := s.userOps.GetUserByID(ctx, claim.UserID)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, user2.ErrUserNotFound
	}

	// calc expiration time values
	var (
		authExp = time.Now().Add(time.Minute * time.Duration(s.tokenExpiration))
	)

	authToken, err := jwt.CreateToken(s.secret, s.userClaims(u, authExp))
	if err != nil {
		return nil, err // todo
	}

	return &UserToken{
		AuthorizationToken: authToken,
		RefreshToken:       refreshToken,
		ExpiresAt:          authExp.UnixMilli(),
	}, nil
}

func (s *AuthService) userClaims(user *user2.User, exp time.Time) *jwt.UserClaims {
	return &jwt.UserClaims{
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: &jwt2.NumericDate{
				Time: exp,
			},
		},
		UserID: user.ID,
		Role:   user.Role.String(),
	}
}
