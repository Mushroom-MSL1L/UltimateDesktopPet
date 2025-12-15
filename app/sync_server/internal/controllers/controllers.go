package controllers

import (
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/sync_server/internal/models"
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/sync_server/internal/services"

	"github.com/gin-gonic/gin"
)

// @Summary		Login with username and password
// @Description	Authenticate user and return a token
// @Description	Only "admin" with password "1234" is valid
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			loginRequest	body		controllers.LoginRequest	true	"Login information"
// @Success		200				{object}	controllers.LoginResponse	"Authentication token"
// @Failure		400				{object}	controllers.ResponseMessage	"Invalid request body | Invalid username or password"
// @Failure		500				{object}	controllers.ResponseMessage	"Could not generate token"
// @Router			/login [post]
func Login(c *gin.Context) {
	var loginReq LoginRequest
	if err := c.BindJSON(&loginReq); err != nil {
		c.JSON(400, ResponseMessage{Message: "Invalid request body"})
		return
	}

	if !models.ValidateAccount(loginReq.Account) {
		c.JSON(400, ResponseMessage{Message: "Invalid username or password"})
		return
	}

	token, err := services.GenerateToken(loginReq.Username)
	if err != nil {
		c.JSON(500, ResponseMessage{Message: "Could not generate token"})
		return
	}

	c.JSON(200, LoginResponse{Token: token})
}

// @Summary		Get user profile
// @Description	Using token and return user profile information
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string						true	"<your_token_here>"
// @Success		200				{object}	controllers.ProfileResponse	"User profile information"
// @Failure		400				{object}	controllers.ResponseMessage	"Invalid request body"
// @Failure		401				{object}	controllers.ResponseMessage	"Unauthorized"
// @Router			/profile [get]
func Profile(c *gin.Context) {
	usename, isExist := c.Get("username")
	if !isExist {
		c.JSON(400, ResponseMessage{Message: "How do you get here?"})
	}
	c.JSON(200, ProfileResponse{Username: usename.(string)})
}

// @Summary		Get server health status
// @Description	Return the current health status of the server
// @Tags			health
// @Produce		json
// @Success		200	{object}	controllers.HealthResponse	"Current health status"
// @Router			/health [get]
func Health(c *gin.Context) {
	c.JSON(200, HealthResponse{Status: "ok"})
}
