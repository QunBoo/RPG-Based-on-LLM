package utils

import (
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var jwtKey = []byte("your_secret_key") // 保持安全，不要公开

func GenerateJWT(userID string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

func AuthenticateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, BearerSchema) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed JWT"})
			return
		}

		tokenString := authHeader[len(BearerSchema):]
		claims := &jwt.StandardClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT Token"})
			return
		}

		// Token is valid
		c.Set("userID", claims.Subject)
		c.Next()
	}
}

func GeneratePasswordHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	hexString := fmt.Sprintf("%x", hash)
	return hexString
	//passwordHash := sha256.Sum256([]byte(password))
	//hashStr := string(passwordHash[:])
	//return hashStr
}
