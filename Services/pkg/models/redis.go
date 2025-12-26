package models

var RedisKey string

const (
	InternalHandlerMiddlewareAuthPermission = "internal:handler:middleware:auth:permission"
	PkgUserGetUsers                         = "pkg:user:get-users"
	PkgUserGetUser                          = "pkg:user:get-user"
	PkgUserGetUserGmail                     = "pkg:user:get-user-gmail"
	PkgRoleGetRoles                         = "pkg:role:get-roles"
	PkgRoleGetRole                          = "pkg:role:get-role"
	PkgRoleGetRoleByName                    = "pkg:role:get-role-by-name"
	PkgPermissionGetPermissions             = "pkg:permission:get-permissions"
	PkgPermissionGetPermission              = "pkg:permission:get-permission"
	PkgRolePermissionGetByRoleID            = "pkg:role-permission:get-by-role-id"
	PkgUserRoleGetByUserID                  = "pkg:user-role:get-by-user-id"
)
