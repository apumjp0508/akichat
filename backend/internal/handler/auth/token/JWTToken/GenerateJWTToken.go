package JWTHandler;

import (
    "time"
    "akichat/backend/internal/config"

    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(config.Load().JwtSecret)

type AccessClaims struct {
    UserID uint   `json:"userID"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}
type RefreshClaims struct {
    UserID uint `json:"userID"`
    jwt.RegisteredClaims
}

func GenerateTokens(userID uint, email string) (accessToken string, refreshToken string, err error) {
    accessExp := time.Now().Add(15 * time.Minute)
    accessClaims := AccessClaims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(accessExp),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Subject:   "access",
        },
    }
    at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessToken, err = at.SignedString(jwtSecret)
    if err != nil {
        return "", "", err
    }

    refreshExp := time.Now().Add(7 * 24 * time.Hour)
    rc := RefreshClaims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(refreshExp),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Subject:   "refresh",
        },
    }
    rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)
    refreshToken, err = rt.SignedString(jwtSecret)
    if err != nil {
        return "", "", err
    }

    return accessToken, refreshToken, nil
}