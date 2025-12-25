package models

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;index"`
	Group     string    `json:"group" gorm:"column:pkg;not null;index"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
type Role struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;index"`
	Name      string    `json:"name" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
type RolePermission struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement;index"`
	RoleID       uint       `json:"role_id" gorm:"not null;index"`
	PermissionID uint       `json:"permission_id" gorm:"not null;index"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	Role         Role       `json:"-" gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Permission   Permission `json:"-" gorm:"foreignKey:PermissionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
type UserRole struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;index"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	RoleID    uint      `json:"role_id" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Role      Role      `json:"-" gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User      User      `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func MigrateTablePermissions(db *gorm.DB) error {
	err := db.AutoMigrate(&Permission{}, &Role{}, &RolePermission{}, &UserRole{})
	if err != nil {
		return err
	}
	return err
}

type PermissionGroup_ string
type PermissionName_ string

const (
	Create PermissionName_ = "create"
	Read   PermissionName_ = "read"
	Update PermissionName_ = "update"
	Delete PermissionName_ = "delete"
	List   PermissionName_ = "list"
	Me     PermissionName_ = "me"
)
const (
	UserGroup           PermissionGroup_ = "user"
	RoleGroup           PermissionGroup_ = "role"
	PermissionGroup     PermissionGroup_ = "permission"
	RolePermissionGroup PermissionGroup_ = "role_permission"
	UserRoleGroup       PermissionGroup_ = "user_role"
)

func PermissionGroupName(group PermissionGroup_, name PermissionName_) string {
	return string(group) + ":" + string(name)
}
