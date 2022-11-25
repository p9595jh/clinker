package service

import (
	"clinker-backend/common/logger"
	"clinker-backend/internal/domain/model/res"
	"clinker-backend/internal/infrastructure/database/entity"
	"clinker-backend/internal/infrastructure/database/repository"
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/gofiber/fiber/v2"
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

func (s *AuthService) PublishToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"id":          user.Id,
		"authority":   user.Authority,
		"confirmed":   user.Confirmed,
		"availableAt": user.StopUntil.Unix(),
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *AuthService) Validate(id, password string) (bool, *entity.User, error) {
	user, err := s.userRepository.FindById(id)
	if err != nil {
		return false, nil, err
	} else if user == nil {
		return false, nil, nil
	} else {
		return s.compare(password, user.Password), user, nil
	}
}

func (s *AuthService) Insert(user *entity.User) (*entity.User, error) {
	if pw, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
		return nil, err
	} else {
		user.Password = string(pw)
	}

	newUser, err := s.userRepository.Save(user)
	if err != nil {
		return nil, err
	} else {
		newUser.Password = ""
		return newUser, nil
	}
}

func (s *AuthService) Available(ctx *fiber.Ctx) *res.ErrorClientRes {
	if !ctx.Locals("confirmed").(bool) {
		return res.NewErrorClientRes(ctx, "not confirmed yet")
	}

	i := ctx.Locals("availableAt").(int64)
	if time.Now().Unix() > i {
		return nil
	}
	return res.NewErrorClientRes(ctx, "available at %d", i)
}
