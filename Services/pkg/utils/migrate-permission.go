package utils

import (
	"github.com/ChangerzaryX1602/SkillSync/pkg/models"

	"gorm.io/gorm"
)

func MigratePermission(db *gorm.DB) error {
	permissions := []models.Permission{
		// User permissions
		{Group: string(models.UserGroup), Name: string(models.Create)},
		{Group: string(models.UserGroup), Name: string(models.Read)},
		{Group: string(models.UserGroup), Name: string(models.Update)},
		{Group: string(models.UserGroup), Name: string(models.Delete)},
		{Group: string(models.UserGroup), Name: string(models.List)},
		// Role permissions
		{Group: string(models.RoleGroup), Name: string(models.Create)},
		{Group: string(models.RoleGroup), Name: string(models.Read)},
		{Group: string(models.RoleGroup), Name: string(models.Update)},
		{Group: string(models.RoleGroup), Name: string(models.Delete)},
		{Group: string(models.RoleGroup), Name: string(models.List)},
		// Permission permissions
		{Group: string(models.PermissionGroup), Name: string(models.Create)},
		{Group: string(models.PermissionGroup), Name: string(models.Read)},
		{Group: string(models.PermissionGroup), Name: string(models.Update)},
		{Group: string(models.PermissionGroup), Name: string(models.Delete)},
		{Group: string(models.PermissionGroup), Name: string(models.List)},
		// RolePermission permissions
		{Group: string(models.RolePermissionGroup), Name: string(models.Create)},
		{Group: string(models.RolePermissionGroup), Name: string(models.Read)},
		{Group: string(models.RolePermissionGroup), Name: string(models.Update)},
		{Group: string(models.RolePermissionGroup), Name: string(models.Delete)},
		{Group: string(models.RolePermissionGroup), Name: string(models.List)},
		// UserRole permissions
		{Group: string(models.UserRoleGroup), Name: string(models.Create)},
		{Group: string(models.UserRoleGroup), Name: string(models.Read)},
		{Group: string(models.UserRoleGroup), Name: string(models.Update)},
		{Group: string(models.UserRoleGroup), Name: string(models.Delete)},
		{Group: string(models.UserRoleGroup), Name: string(models.List)},
	}

	for _, perm := range permissions {
		var existingPerm models.Permission
		err := db.Where("pkg = ? AND name = ?", perm.Group, perm.Name).First(&existingPerm).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&perm).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	err := MigrateRoles(db, []RolePermission{
		// Admin - Full access to everything
		{
			Role: models.Role{Name: string(Admin)},
			Permissions: []models.Permission{
				// User permissions
				{Group: string(models.UserGroup), Name: string(models.Create)},
				{Group: string(models.UserGroup), Name: string(models.Read)},
				{Group: string(models.UserGroup), Name: string(models.Update)},
				{Group: string(models.UserGroup), Name: string(models.Delete)},
				{Group: string(models.UserGroup), Name: string(models.List)},
				// Role permissions
				{Group: string(models.RoleGroup), Name: string(models.Create)},
				{Group: string(models.RoleGroup), Name: string(models.Read)},
				{Group: string(models.RoleGroup), Name: string(models.Update)},
				{Group: string(models.RoleGroup), Name: string(models.Delete)},
				{Group: string(models.RoleGroup), Name: string(models.List)},
				// Permission permissions
				{Group: string(models.PermissionGroup), Name: string(models.Create)},
				{Group: string(models.PermissionGroup), Name: string(models.Read)},
				{Group: string(models.PermissionGroup), Name: string(models.Update)},
				{Group: string(models.PermissionGroup), Name: string(models.Delete)},
				{Group: string(models.PermissionGroup), Name: string(models.List)},
				// RolePermission permissions
				{Group: string(models.RolePermissionGroup), Name: string(models.Create)},
				{Group: string(models.RolePermissionGroup), Name: string(models.Read)},
				{Group: string(models.RolePermissionGroup), Name: string(models.Update)},
				{Group: string(models.RolePermissionGroup), Name: string(models.Delete)},
				{Group: string(models.RolePermissionGroup), Name: string(models.List)},
				// UserRole permissions
				{Group: string(models.UserRoleGroup), Name: string(models.Create)},
				{Group: string(models.UserRoleGroup), Name: string(models.Read)},
				{Group: string(models.UserRoleGroup), Name: string(models.Update)},
				{Group: string(models.UserRoleGroup), Name: string(models.Delete)},
				{Group: string(models.UserRoleGroup), Name: string(models.List)},
			},
		},
		// User - Basic read access only
		{
			Role: models.Role{Name: string(User)},
			Permissions: []models.Permission{
				// UserGroup
				{Group: string(models.UserGroup), Name: string(models.Read)},
				{Group: string(models.UserGroup), Name: string(models.List)},
				{Group: string(models.UserGroup), Name: string(models.Me)},

				// PermissionGroup
				{Group: string(models.PermissionGroup), Name: string(models.Me)},
			},
		},
	})

	return err
}

func MigrateRoles(db *gorm.DB, rolePermissions []RolePermission) error {
	for _, rp := range rolePermissions {
		var role models.Role

		err := db.Where("name = ?", rp.Role.Name).First(&role).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				role = rp.Role
				if err := db.Create(&role).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		for _, pDef := range rp.Permissions {
			var permission models.Permission

			if err := db.Where("pkg = ? AND name = ?", pDef.Group, pDef.Name).First(&permission).Error; err != nil {

				continue
			}

			var existingRP models.RolePermission
			err := db.Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).First(&existingRP).Error

			if err == gorm.ErrRecordNotFound {

				newRP := models.RolePermission{
					RoleID:       role.ID,
					PermissionID: permission.ID,
				}
				if err := db.Create(&newRP).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

type RolePermission struct {
	Role        models.Role
	Permissions []models.Permission
}
type Role string

const (
	Admin Role = "Admin"
	User  Role = "User"
)
