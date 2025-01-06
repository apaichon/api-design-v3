package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"api/config"
	"api/internal/handler"

	"github.com/dgrijalva/jwt-go"
)

func cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}
}

func generateSalt() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	cors(w, r)
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userRepo := NewUserRepo()
	config := config.NewConfig()

	userdb, ok := userRepo.GetUserByName(user.Username) // users[user.Username]
	if ok != nil {
		http.Error(w, "Invalid username", http.StatusUnauthorized)
		return
	}
	// fmt.Println("User")
	hash := HashString(user.Username + user.Password + userdb.Salt)
	// fmt.Println("UserDbPassword:", userdb.Password, "Hash:",hash, "UserNameDb:", userdb.Username, "username", user.Username )
	if userdb.Password != hash {
		http.Error(w, "Invalid username", http.StatusUnauthorized)
		return
	}

	// fmt.Println("Duration", time.Duration(config.TokenAge))
	expirationTime := time.Now().Add(time.Duration(config.TokenAge) * time.Minute)

	claims := &JwtClaims{
		UserId:   userdb.UserId,
		Username: userdb.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	response := JwtToken{
		Token:     tokenString,
		ExpiredAt: expirationTime.Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AuthenticationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the JWT token is present in the request
		tokenString := getTokenFromRequest(r)
		// fmt.Println("tokenString:" + tokenString)
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate the JWT token
		token, err := validateToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, tokenString)
		r = r.WithContext(ctx)

		// Pass the request to the next handler if the token is valid
		next.ServeHTTP(w, r)
	})
}

// DecodeJWTToken decodes a JWT and verifies its signature with a secret key
func DecodeJWTToken(tokenString, secretKey string) (*JwtClaims, error) {
	config := config.NewConfig()
	// Parse the token and validate the signature
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is what we expect (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}

		// Return the secret key for signature verification
		return []byte(config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid and contains our expected claims
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token or claims")
}

func getTokenFromRequest(r *http.Request) string {
	// Check if the token is present in the request header
	token := r.Header.Get("Authorization")
	if token != "" {
		return token
	}

	// Check if the token is present in the request cookies
	cookie, err := r.Cookie("token")
	if err == nil {
		return cookie.Value
	}

	return ""
}

func validateToken(tokenString string) (*jwt.Token, error) {
	config := config.NewConfig()
	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil // Replace "your-secret-key" with your actual secret key
	})

	// fmt.Println("Token", token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		resp := handler.NewErrorResponse(
			http.StatusBadRequest,
			"Bad Request",
			"INVALID_REQUEST",
			"Invalid request body",
			r.Header.Get("X-Request-ID"),
		)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	if user.Username == "" || user.Password == "" {
		resp := handler.NewErrorResponse(
			http.StatusBadRequest,
			"Bad Request",
			"MISSING_FIELDS",
			"Username and password are required",
			r.Header.Get("X-Request-ID"),
		)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	userRepo := NewUserRepo()
	exists, err := userRepo.ExistsUserByName(user.Username)
	if err != nil {
		resp := handler.NewErrorResponse(
			http.StatusInternalServerError,
			"Internal Server Error",
			"DATABASE_ERROR",
			"Database operation failed",
			r.Header.Get("X-Request-ID"),
		)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	if exists {
		resp := handler.NewErrorResponse(
			http.StatusConflict,
			"Conflict",
			"USER_EXISTS",
			"Username already exists",
			r.Header.Get("X-Request-ID"),
		)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	salt := generateSalt()
	hashedPassword := HashString(user.Username + user.Password + salt)

	newUser := User{
		Username:  user.Username,
		Password:  hashedPassword,
		Salt:      salt,
		CreatedAt: time.Now(),
		CreatedBy: user.Username,
		StatusID:  1,
	}

	err = userRepo.CreateUser(&newUser)
	if err != nil {
		resp := handler.NewErrorResponse(
			http.StatusInternalServerError,
			"Internal Server Error",
			"REGISTRATION_FAILED",
			"Failed to create user",
			r.Header.Get("X-Request-ID"),
		)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	responseData := map[string]string{
		"username": user.Username,
		"message":  "User registered successfully",
	}

	resp := handler.NewResponse(http.StatusCreated, "Created", responseData, r.Header.Get("X-Request-ID"))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
