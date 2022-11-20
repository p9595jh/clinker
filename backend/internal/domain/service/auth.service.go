package service

import (
	"clinker-backend/common/logger"
	"clinker-backend/internal/infrastructure/database/entity"
	"clinker-backend/internal/infrastructure/database/repository"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	privateKey     *rsa.PrivateKey
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) *AuthService {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return &AuthService{
		privateKey:     privateKey,
		userRepository: userRepository,
	}
}

func (s *AuthService) PK() *rsa.PrivateKey {
	return s.privateKey
}

func (s *AuthService) compare(input, password string) bool {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("AuthService", "compare", "GenerateFromPassword").E(err).D("input", input).D("password", password).W()
		return false
	} else {
		return bcrypt.CompareHashAndPassword(hash, []byte(input)) == nil
	}
}

func (s *AuthService) PublishToken(id string) (string, error) {
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *AuthService) Validate(id, password string) (bool, error) {
	admin, err := s.userRepository.FindById(id)
	if err != nil {
		return false, err
	} else if admin == nil {
		return false, nil
	} else {
		return s.compare(password, admin.Password), nil
	}
}

func (s *AuthService) Insert(user *entity.User) (*entity.User, error) {
	fmt.Println("register:", user)
	if pw, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
		return nil, err
	} else {
		user.Password = string(pw)
	}

	newUser, err := s.userRepository.Save(user)
	if err != nil {
		return nil, err
	} else {
		fmt.Println("registered:", newUser)
		newUser.Password = ""
		return newUser, nil
	}
}
