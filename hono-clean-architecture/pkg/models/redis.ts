export const CacheKeys = {
  PkgUserGetUser: "pkg:user:get",
  PkgUserGetUserByEmail: "pkg:user:email",
  PkgUserGetList: "pkg:user:list",

  RefreshToken: "refresh",

  InternalHandlerMiddlewareAuthPermission: "internal:handler:middleware:auth:permission",
} as const;

export type CacheKeyType = (typeof CacheKeys)[keyof typeof CacheKeys];
