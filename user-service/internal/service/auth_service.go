package service

import (
	"errors"
	"log"

	"user-service/internal/events"
	"user-service/internal/model"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	userRepository *repository.UserRepository
	tokenService   *TokenService
	publisher      *events.Publisher
}

func NewAuthService(userRepository *repository.UserRepository, tokenService *TokenService, publisher *events.Publisher) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		tokenService:   tokenService,
		publisher:      publisher,
	}
}

func (s *AuthService) Register(email, password string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.userRepository.Create(user); err != nil {
		return nil, err
	}

	if s.publisher != nil {
		if err := s.publisher.Publish("user.created", map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		}); err != nil {
			log.Printf("failed to publish user.created event: %v", err)
		}
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (*model.User, string, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrInvalidCredentials
		}

		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.tokenService.Generate(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) GetByID(id uint) (*model.User, error) {
	return s.userRepository.FindByID(id)
}
