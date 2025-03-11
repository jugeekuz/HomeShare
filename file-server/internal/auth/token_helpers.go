package auth


import (
	"fmt"
	"os"
	"log"
	"errors"
	"time"
	"database/sql"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(accessParams, refreshParams *TokenParameters) (string, string, error) {
	accessClaims := jwt.MapClaims{
		"user_id":   accessParams.UserId,
		"folder_id": accessParams.FolderId,
		"access":    accessParams.Access,
		"exp":       time.Now().Add(accessParams.ExpiryDuration * time.Hour).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.MapClaims{
		"user_id":   refreshParams.UserId,
		"folder_id": refreshParams.FolderId,
		"access":    refreshParams.Access,
		"exp":       time.Now().Add(refreshParams.ExpiryDuration * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func DecodeToken(tokenStr string, secret string) (TokenParameters, error) {
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return TokenParameters{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return TokenParameters{}, errors.New("invalid token")
	}

	var params TokenParameters

	if uid, ok := claims["user_id"].(string); ok {
		params.UserId = uid
	}

	if fid, ok := claims["folder_id"].(string); ok {
		params.FolderId = fid
	}

	if access, ok := claims["access"].(string); ok {
		params.Access = access
	}

	if exp, ok := claims["exp"].(float64); ok {
		expTime := time.Unix(int64(exp), 0)
		params.ExpiryDuration = time.Until(expTime)
	}

	return params, nil
}

var (
	infoLog  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func logAsync(logger *log.Logger, message string) {
	go func() {
		logger.Println(message)
	}()
}

func Authenticate(db *sql.DB, creds Credentials) (*User, error) {
	logAsync(infoLog, "Starting authentication process")

	user, err := getUserByUsername(db, creds.Username)
	if err != nil {
		logAsync(errorLog, "Failed to retrieve user: "+err.Error())
		return nil, err
	}

	logAsync(infoLog, "User retrieved successfully: "+user.Username)

	hashedPassword := HashPassword(creds.Password, user.Salt)
	if hashedPassword != user.PasswordHash {
		logAsync(errorLog, "Invalid credentials for user: "+user.Username)
		return nil, errors.New("invalid credentials")
	}

	logAsync(infoLog, "Authentication successful for user: "+user.Username)
	return user, nil
}