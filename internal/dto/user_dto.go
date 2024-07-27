package dto

type CreateUserDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthUserRequestDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthUserResponseDTO struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}
