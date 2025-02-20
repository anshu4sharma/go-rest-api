package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/anshu4sharma/go-rest-api/internal/models"
	"github.com/anshu4sharma/go-rest-api/internal/repository"
	"golang.org/x/oauth2"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type GoogleAuthService struct {
	config   *oauth2.Config
	userRepo *repository.UserRepository
	authSvc  *AuthService
}

func NewGoogleAuthService(config *oauth2.Config, userRepo *repository.UserRepository, authSvc *AuthService) *GoogleAuthService {
	return &GoogleAuthService{
		config:   config,
		userRepo: userRepo,
		authSvc:  authSvc,
	}
}

func (s *GoogleAuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state)
}

func (s *GoogleAuthService) HandleCallback(code string) (*models.TokenResponse, error) {
	token, err := s.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.getUserInfo(token.AccessToken)
	if err != nil {
		return nil, err
	}

	// Check if user exists
	user, err := s.userRepo.FindByGoogleID(userInfo.ID)
	if err != nil {
		// Create new user
		now := time.Now()
		user = &models.User{
			Email:       userInfo.Email,
			Name:        userInfo.Name,
			GoogleID:    &userInfo.ID,
			Role:        "user",
			LastLoginAt: &now,
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, err
		}
	} else {
		// Update last login time
		now := time.Now()
		user.LastLoginAt = &now
		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}
	}

	// Generate JWT tokens
	accessToken, expiresIn, err := s.authSvc.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.authSvc.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *GoogleAuthService) getUserInfo(accessToken string) (*GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info from Google")
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
