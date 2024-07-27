package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SomchaiSPB/user-auth/internal/dto"
	"github.com/SomchaiSPB/user-auth/internal/entity"
	"github.com/SomchaiSPB/user-auth/internal/hash"
	"github.com/SomchaiSPB/user-auth/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"time"
)

const tokenExpTime = 15 * time.Minute

var validate *validator.Validate

var (
	ErrUserNameExists   = errors.New("user already exists error")
	ErrCreateUser       = errors.New("user create error")
	ErrWrongCredentials = errors.New("wrong credentials error")
	ErrGenerateToken    = errors.New("generate token error")
	ErrValidation       = errors.New("validation error")
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserSvc(ur repository.UserRepository) *UserService {
	return &UserService{userRepository: ur}
}

func (s UserService) Create(data []byte) ([]byte, error) {
	var userDto dto.CreateUserDTO

	if err := json.Unmarshal(data, &userDto); err != nil {
		return nil, err
	}

	err := validate.Struct(userDto)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
	}

	if exists := s.userRepository.Exists(userDto.Username); exists {
		return nil, fmt.Errorf("%s: %w", userDto.Username, ErrUserNameExists)
	}

	hashedPass, err := hash.NewHasher().HashPassword(userDto.Password)

	if err != nil {
		return nil, fmt.Errorf("password hash error: %v", err)
	}

	u := &entity.User{
		Name:     userDto.Username,
		Password: hashedPass,
	}

	createdUser, err := s.userRepository.Create(u)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreateUser, err)
	}

	createdUser.Password = ""

	return json.Marshal(createdUser)
}

func (s UserService) Authenticate(data []byte, jwtSecret []byte) ([]byte, error) {
	var authDto dto.AuthUserRequestDTO

	if err := json.Unmarshal(data, &authDto); err != nil {
		return nil, err
	}

	err := validate.Struct(authDto)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
	}

	u, err := s.userRepository.GetByName(authDto.Username)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWrongCredentials
		}
		return nil, err
	}

	if ok := hash.NewHasher().CheckPasswordHash(authDto.Password, u.Password); !ok {
		return nil, ErrWrongCredentials
	}

	exp := time.Now().Add(tokenExpTime).Unix()

	tkn, err := generateJWT(u.Name, jwtSecret, exp)

	if err != nil {
		return nil, ErrGenerateToken
	}

	response := dto.AuthUserResponseDTO{
		Token:     tkn,
		ExpiresAt: exp,
	}

	return json.Marshal(response)
}

func generateJWT(username string, jwtSecret []byte, exp int64) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}
