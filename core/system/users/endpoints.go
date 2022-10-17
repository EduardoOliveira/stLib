package users

import (
	"log"
	"net/http"

	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func login(c echo.Context) error {
	login := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := c.Bind(&login)
	if err != nil {
		log.Println("Error binding login struct", err)
		return c.NoContent(http.StatusBadRequest)
	}
	log.Println(login)

	if login.Username == "" || login.Password == "" {
		log.Println("Username or password is empty")
		return c.NoContent(http.StatusBadRequest)
	}
	log.Println(users)
	user, ok := users[login.Username]
	if !ok {
		log.Println("User not found")
		return c.NoContent(http.StatusUnauthorized)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		log.Println("Password does not match", err)
		return c.NoContent(http.StatusUnauthorized)
	}

	claims := jwt.MapClaims{
		"sub":      user.Uuid,
		"username": user.Username,
	}

	for _, permission := range user.Permissions {
		claims[permission] = true
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(runtime.Cfg.JwtSecret))

	if err != nil {
		log.Println("Error signing token", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": tokenString,
	})
}
