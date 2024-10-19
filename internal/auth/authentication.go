package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), err
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiriesIn time.Duration) (string, error) {

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiriesIn).UTC()),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(tokenSecret))

	return ss, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims := token.Claims.(*jwt.RegisteredClaims)
	if claims.ExpiresAt.Time.UTC().Before(time.Now().UTC()) {
		return uuid.Nil, fmt.Errorf("token expired at %v", claims.ExpiresAt)
	}

	return uuid.MustParse(claims.Subject), nil
}

func GetBearerToken(header http.Header) (string, error) {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("there is no authorization header")
	}
	tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
	if !ok {
		return "", fmt.Errorf("there is no bearer token")
	}
	return tokenString, nil
}

func MakeRefreshToken() (string, error) {
	randString := make([]byte, 32)
	_, err := rand.Read(randString)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randString), nil
}
