package handler

import (
	"net/http"
	"user-service/config"
	"user-service/dto"
	"user-service/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(name, email, password string) (int, error) {
	user := models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with the provided details.
// @Tags Users
// @Accept json
// @Produce json
// @Param register body entity.RegisterUser true "Register User"
// @Success 200 {object} map[string]interface{} "Success message and user details"
// @Failure 400 {object} map[string]interface{} "Invalid Request Parameters"
// @Failure 500 {object} map[string]interface{} "Register Failed"
// @Router /users/register [post]
func Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid Request Parameters"})
	}

	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Name is required"})
	}
	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Email is required"})
	}
	if req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Password is required"})
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Invalid generate password"})
	}

	userID, err := CreateUser(req.Name, req.Email, string(hashPassword))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Register Failed"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":    userID,
		"name":  req.Name,
		"email": req.Email,
	})
}
