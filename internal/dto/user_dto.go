package dto

// CreateUserDTO represents the data required to create a user
type CreateUserDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AuthUserRequestDTO represents the authentication request data
type AuthUserRequestDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AuthUserResponseDTO represents the response after successful authentication
type AuthUserResponseDTO struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}
