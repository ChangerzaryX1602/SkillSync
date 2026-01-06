import { NodePgDatabase } from "drizzle-orm/node-postgres";
import { eq, and } from "drizzle-orm";
import {
  permissions,
  roles,
  rolePermissions,
  Role,
  PermissionName,
  PermissionGroupType,
} from "../models";

export const RoleType = {
  Admin: "admin",
  User: "user",
} as const;

export type RoleTypeValue = (typeof RoleType)[keyof typeof RoleType];

interface RolePermissionConfig {
  role: { name: string };
  permissions: { group: string; name: string }[];
}

export async function migratePermissions(db: NodePgDatabase): Promise<void> {
  const allPermissions = [
    { group: PermissionGroupType.User, name: PermissionName.Create },
    { group: PermissionGroupType.User, name: PermissionName.Read },
    { group: PermissionGroupType.User, name: PermissionName.Update },
    { group: PermissionGroupType.User, name: PermissionName.Delete },
    { group: PermissionGroupType.User, name: PermissionName.List },
    { group: PermissionGroupType.User, name: PermissionName.Me },

    { group: PermissionGroupType.Role, name: PermissionName.Create },
    { group: PermissionGroupType.Role, name: PermissionName.Read },
    { group: PermissionGroupType.Role, name: PermissionName.Update },
    { group: PermissionGroupType.Role, name: PermissionName.Delete },
    { group: PermissionGroupType.Role, name: PermissionName.List },

    { group: PermissionGroupType.Permission, name: PermissionName.Create },
    { group: PermissionGroupType.Permission, name: PermissionName.Read },
    { group: PermissionGroupType.Permission, name: PermissionName.Update },
    { group: PermissionGroupType.Permission, name: PermissionName.Delete },
    { group: PermissionGroupType.Permission, name: PermissionName.List },

    { group: PermissionGroupType.RolePermission, name: PermissionName.Create },
    { group: PermissionGroupType.RolePermission, name: PermissionName.Read },
    { group: PermissionGroupType.RolePermission, name: PermissionName.Update },
    { group: PermissionGroupType.RolePermission, name: PermissionName.Delete },
    { group: PermissionGroupType.RolePermission, name: PermissionName.List },

    { group: PermissionGroupType.UserRole, name: PermissionName.Create },
    { group: PermissionGroupType.UserRole, name: PermissionName.Read },
    { group: PermissionGroupType.UserRole, name: PermissionName.Update },
    { group: PermissionGroupType.UserRole, name: PermissionName.Delete },
    { group: PermissionGroupType.UserRole, name: PermissionName.List },
  ];

  for (const perm of allPermissions) {
    const existing = await db
      .select()
      .from(permissions)
      .where(and(eq(permissions.group, perm.group), eq(permissions.name, perm.name)))
      .limit(1);

    if (existing.length === 0) {
      await db.insert(permissions).values(perm);
    }
  }

  const roleConfigs: RolePermissionConfig[] = [
    {
      role: { name: RoleType.Admin },
      permissions: allPermissions,
    },
    {
      role: { name: RoleType.User },
      permissions: [
        { group: PermissionGroupType.User, name: PermissionName.Read },
        { group: PermissionGroupType.User, name: PermissionName.Me },
      ],
    },
  ];

  for (const config of roleConfigs) {
    let roleResult = await db.select().from(roles).where(eq(roles.name, config.role.name)).limit(1);

    let role: Role;
    if (roleResult.length === 0) {
      const insertedRoles = await db.insert(roles).values({ name: config.role.name }).returning();
      role = insertedRoles[0];
    } else {
      role = roleResult[0];
    }

    for (const perm of config.permissions) {
      const permResult = await db
        .select()
        .from(permissions)
        .where(and(eq(permissions.group, perm.group), eq(permissions.name, perm.name)))
        .limit(1);

      if (permResult.length > 0) {
        const permission = permResult[0];

        const existing = await db
          .select()
          .from(rolePermissions)
          .where(
            and(
              eq(rolePermissions.roleId, role.id),
              eq(rolePermissions.permissionId, permission.id)
            )
          )
          .limit(1);

        if (existing.length === 0) {
          await db.insert(rolePermissions).values({
            roleId: role.id,
            permissionId: permission.id,
          });
        }
      }
    }
  }
}
