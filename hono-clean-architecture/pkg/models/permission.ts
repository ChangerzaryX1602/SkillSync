import { pgTable, serial, varchar, timestamp, integer } from "drizzle-orm/pg-core";
import { InferSelectModel, InferInsertModel } from "drizzle-orm";
import { users } from "./user";

export const permissions = pgTable("permissions", {
  id: serial("id").primaryKey(),
  group: varchar("pkg", { length: 255 }).notNull(),
  name: varchar("name", { length: 255 }).notNull(),
  createdAt: timestamp("created_at").defaultNow().notNull(),
  updatedAt: timestamp("updated_at").defaultNow().notNull(),
});

export type Permission = InferSelectModel<typeof permissions>;
export type NewPermission = InferInsertModel<typeof permissions>;

export const roles = pgTable("roles", {
  id: serial("id").primaryKey(),
  name: varchar("name", { length: 255 }).notNull().unique(),
  createdAt: timestamp("created_at").defaultNow().notNull(),
  updatedAt: timestamp("updated_at").defaultNow().notNull(),
});

export type Role = InferSelectModel<typeof roles>;
export type NewRole = InferInsertModel<typeof roles>;

export const rolePermissions = pgTable("role_permissions", {
  id: serial("id").primaryKey(),
  roleId: integer("role_id")
    .notNull()
    .references(() => roles.id, { onDelete: "cascade", onUpdate: "cascade" }),
  permissionId: integer("permission_id")
    .notNull()
    .references(() => permissions.id, {
      onDelete: "cascade",
      onUpdate: "cascade",
    }),
  createdAt: timestamp("created_at").defaultNow().notNull(),
  updatedAt: timestamp("updated_at").defaultNow().notNull(),
});

export type RolePermission = InferSelectModel<typeof rolePermissions>;
export type NewRolePermission = InferInsertModel<typeof rolePermissions>;

export const userRoles = pgTable("user_roles", {
  id: serial("id").primaryKey(),
  userId: integer("user_id")
    .notNull()
    .references(() => users.id, { onDelete: "cascade", onUpdate: "cascade" }),
  roleId: integer("role_id")
    .notNull()
    .references(() => roles.id, { onDelete: "cascade", onUpdate: "cascade" }),
  createdAt: timestamp("created_at").defaultNow().notNull(),
  updatedAt: timestamp("updated_at").defaultNow().notNull(),
});

export type UserRole = InferSelectModel<typeof userRoles>;
export type NewUserRole = InferInsertModel<typeof userRoles>;

export const PermissionName = {
  Create: "create",
  Read: "read",
  Update: "update",
  Delete: "delete",
  List: "list",
  Me: "me",
} as const;

export type PermissionNameType = (typeof PermissionName)[keyof typeof PermissionName];

export const PermissionGroupType = {
  User: "user",
  Role: "role",
  Permission: "permission",
  RolePermission: "role_permission",
  UserRole: "user_role",
} as const;

export type PermissionGroupTypeValue =
  (typeof PermissionGroupType)[keyof typeof PermissionGroupType];

export function permissionGroupName(
  group: PermissionGroupTypeValue,
  name: PermissionNameType
): string {
  return `${group}:${name}`;
}
