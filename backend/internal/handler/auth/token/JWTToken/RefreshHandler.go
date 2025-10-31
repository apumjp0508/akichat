package JWTHandler;

import (
	"time"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func RefreshHandler(c *gin.Context) {
    // Cookie から refresh_token を取得
    refresh_token, err := c.Cookie("refresh_token")
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not provided"})
        return
    }

    // トークンをパース
    token, err := jwt.ParseWithClaims(refresh_token, &RefreshClaims{}, func(t *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
        return
    }

    rc := token.Claims.(*RefreshClaims)
    // 追加チェック：ユーザー存在、リフレッシュトークンが有効か（DBチェックなど）

    // 新しいトークンペアを発行
    newAccess, newRefresh, err := GenerateTokens(rc.UserID, "") 
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate tokens"})
        return
    }

    // 新しい refresh_token を Cookie にセットし直す（トークンローテーション）
    c.SetCookie("refresh_token", newRefresh, int(7*24*time.Hour.Seconds()), "/", "", false, true)

    c.JSON(http.StatusOK, gin.H{
        "access_token": newAccess,
    })
}