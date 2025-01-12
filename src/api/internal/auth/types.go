package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	UserID    int       `json:"user_id"`
	Username  string    `json:"user_name"`
	Password  string    `json:"password"`
	Salt      string    `json:"salt"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	StatusID  int       `json:"status_id"`
}

type Role struct {
	RoleID       int       `json:"role_id"`
	RoleName     string    `json:"role_name"`
	RoleDesc     string    `json:"role_desc,omitempty"`
	IsSuperAdmin bool      `json:"is_super_admin"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	StatusID     int       `json:"status_id"`
}

type UserRoles struct {
	UserRoleID int       `json:"user_role_id"`
	RoleID     int       `json:"role_id"`
	UserID     string    `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
	StatusID   int       `json:"status_id"`
}

type RolePermissions struct {
	RolePermissionID   int       `json:"role_permission_id"`
	RolePermissionDesc string    `json:"role_permission_desc"`
	ResourceTypeID     int       `json:"resource_type_id"`
	ResourceName       string    `json:"resource_name"`
	CanExecute         bool      `json:"can_execute"`
	CanRead            bool      `json:"can_read"`
	CanWrite           bool      `json:"can_write"`
	CanDelete          bool      `json:"can_delete"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          string    `json:"updated_by"`
	StatusID           int       `json:"status_id"`
}

type PermissionView struct {
	UserID       string
	RoleID       int
	RoleName     string
	ResourceName string
	CanExecute   bool
	CanRead      bool
	CanWrite     bool
	CanDelete    bool
}

type JwtClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"user_name"`
	jwt.StandardClaims
}

type JwtToken struct {
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expiredAt"`
}

type UserPermissionView struct {
	UserID           string
	RoleID           int
	RoleName         string
	IsSuperAdmin     bool
	RolePermissionID int
	ResourceTypeID   int
	ResourceName     string
	CanExecute       bool
	CanRead          bool
	CanWrite         bool
	CanDelete        bool
}

type IsSuperAdmin struct {
	UserID       int
	IsSuperAdmin int
}
