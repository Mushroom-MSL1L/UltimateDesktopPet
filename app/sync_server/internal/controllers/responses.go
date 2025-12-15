package controllers

import "github.com/Mushroom-MSL1L/UltimateDesktopPet/app/sync_server/internal/models"

type LoginRequest struct {
	models.Account
}

type LoginResponse struct {
	Token string `json:"token" example:"<Header>.<Payload>.<Signature>"`
}

type ProfileRequest struct {
	LoginResponse
}

type ProfileResponse struct {
	Username string `json:"username" example:"admin"`
}

type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

type ResponseMessage struct {
	Message string `json:"message" example:"error message"`
}
