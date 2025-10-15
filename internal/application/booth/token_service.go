package booth

import (
	"context"
	"errors"
	"fmt"

	domainBooth "go-ddd-clean/internal/domain/booth"

	"github.com/golang-jwt/jwt/v5"
)

const boothAccessTokenType = "booth_access"

var (
	ErrInvalidBoothToken = errors.New("invalid booth token")
	ErrTokenMismatch     = errors.New("booth token does not match current booth state")
)

type tokenClaims struct {
	BoothID      string `json:"booth_id"`
	BranchID     string `json:"branch_id"`
	TokenVersion int    `json:"token_version"`
	Type         string `json:"type"`
	jwt.RegisteredClaims
}

type ValidatedToken struct {
	BoothID      string
	BranchID     string
	TokenVersion int
	Raw          string
}

type TokenService struct {
	repo   domainBooth.Repository
	secret []byte
}

func NewTokenService(repo domainBooth.Repository, secret string) *TokenService {
	return &TokenService{
		repo:   repo,
		secret: []byte(secret),
	}
}

func (s *TokenService) Register(ctx context.Context, boothID string, branchID string) (string, error) {
	entity, err := s.repo.GetByID(ctx, boothID)
	if err != nil {
		return "", err
	}
	if branchID != "" && entity.BranchID != branchID {
		return "", ErrTokenMismatch
	}
	if entity.TokenVersion == 0 {
		entity.TokenVersion = 1
		if err := s.repo.UpdateTokenVersion(ctx, entity.ID, entity.TokenVersion); err != nil {
			return "", err
		}
	}
	return s.sign(entity)
}

func (s *TokenService) Regenerate(ctx context.Context, boothID string) (string, error) {
	entity, err := s.repo.GetByID(ctx, boothID)
	if err != nil {
		return "", err
	}
	entity.TokenVersion++
	if entity.TokenVersion <= 0 {
		entity.TokenVersion = 1
	}
	if err := s.repo.UpdateTokenVersion(ctx, entity.ID, entity.TokenVersion); err != nil {
		return "", err
	}
	return s.sign(entity)
}

func (s *TokenService) Validate(ctx context.Context, token string) (*ValidatedToken, error) {
	claims, err := s.parse(token)
	if err != nil {
		return nil, err
	}
	if claims.Type != boothAccessTokenType || claims.BoothID == "" || claims.BranchID == "" || claims.TokenVersion <= 0 {
		return nil, ErrInvalidBoothToken
	}
	entity, err := s.repo.GetByID(ctx, claims.BoothID)
	if err != nil {
		return nil, err
	}
	if entity.BranchID != claims.BranchID || entity.TokenVersion != claims.TokenVersion {
		return nil, ErrTokenMismatch
	}
	return &ValidatedToken{
		BoothID:      entity.ID,
		BranchID:     entity.BranchID,
		TokenVersion: entity.TokenVersion,
		Raw:          token,
	}, nil
}

func (s *TokenService) parse(token string) (*tokenClaims, error) {
	claims := &tokenClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, ErrInvalidBoothToken
	}
	if !parsedToken.Valid {
		return nil, ErrInvalidBoothToken
	}
	return claims, nil
}

func (s *TokenService) sign(entity *domainBooth.Booth) (string, error) {
	if entity == nil {
		return "", ErrInvalidBoothToken
	}
	claims := tokenClaims{
		BoothID:      entity.ID,
		BranchID:     entity.BranchID,
		TokenVersion: entity.TokenVersion,
		Type:         boothAccessTokenType,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}
