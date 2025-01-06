package auth

import (
	"time"
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	UserId    int       `json:"user_id"`
	Username  string    `json:"user_name"`
	Password  string    `json:"password"`
	Salt      string    `json:"salt"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	StatusID  int       `json:"status_id"`
}


type JwtClaims struct {
	UserId int `json:"user_id"`
	Username string `json:"user_name"`
	jwt.StandardClaims
}


type JwtToken struct {
	Token string `json:"token"`
	ExpiredAt int64 `json:"expiredAt"`
	
}
