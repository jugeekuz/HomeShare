package auth

import (
	"errors"
	"fmt"
	"time"
	"database/sql"

	"github.com/golang-jwt/jwt/v5"

	"file-server/config"
)

func GenerateTokens(accessParams, refreshParams *TokenParameters) (string, string, error) {
	cfg := config.LoadConfig()

	accessClaims := jwt.MapClaims{
		"user_id":   accessParams.UserId,
		"folder_id": accessParams.FolderId,
		"access":    accessParams.Access,
		"exp":       time.Now().Add(accessParams.ExpiryDuration).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(cfg.Secrets.Jwt.JwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.MapClaims{
		"user_id":   refreshParams.UserId,
		"folder_id": refreshParams.FolderId,
		"access":    refreshParams.Access,
		"exp":       time.Now().Add(refreshParams.ExpiryDuration).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(cfg.Secrets.Jwt.JwtSecret))
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


func Authenticate(db *sql.DB, creds Credentials) (*User, error) {

	user, err := getUserByUsername(db, creds.Username)
	if err != nil {
		return nil, err
	}

	hashedPassword := HashPassword(creds.Password, user.Salt)
	if hashedPassword != user.PasswordHash {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func HasAccess(claims jwt.MapClaims, folderID, requiredAccess string) (bool, error) {
	claimFolderID, ok := claims["folder_id"].(string)
	if !ok {
		return false, errors.New("folder_id claim missing or invalid")
	}

	claimAccess, ok := claims["access"].(string)
	if !ok {
		return false, errors.New("access claim missing or invalid")
	}

	if claimFolderID != folderID && claimFolderID != "/" {
		fmt.Printf("folderid : %s\nclaimFolderId: %s\n\n", folderID, claimFolderID)
		return false, nil
	}

	if requiredAccess == "r" && (claimAccess == "r" || claimAccess == "rw") {
		return true, nil
	}
	if requiredAccess == "w" && (claimAccess == "w" || claimAccess == "rw") {
		return true, nil
	}
	if requiredAccess == "rw" && claimAccess == "rw" {
		return true, nil
	}

	return false, nil
}
