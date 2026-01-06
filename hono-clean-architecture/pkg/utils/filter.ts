import { SQL, sql, or, ilike } from "drizzle-orm";
import { PgColumn } from "drizzle-orm/pg-core";
import { Search, Pagination } from "../models";

/**
 * Apply search filter using SQL ILIKE (case-insensitive LIKE)
 * This replaces the vector embedding search from the Go version
 */
export function applySearch(search: Search, columns: PgColumn[]): SQL | undefined {
  if (!search.keyword || search.keyword.trim() === "") {
    return undefined;
  }

  const keyword = `%${search.keyword}%`;

  if (search.column && search.column.trim() !== "") {
    const columnNames = search.column.split(",").map((c) => c.trim());
    const targetColumns = columns.filter((col) => columnNames.includes(col.name));

    if (targetColumns.length === 0) {
      return undefined;
    }

    if (targetColumns.length === 1) {
      return ilike(targetColumns[0], keyword);
    }

    const conditions = targetColumns.map((col) => ilike(col, keyword));
    return or(...conditions);
  }

  if (columns.length === 0) {
    return undefined;
  }

  if (columns.length === 1) {
    return ilike(columns[0], keyword);
  }

  const conditions = columns.map((col) => ilike(col, keyword));
  return or(...conditions);
}

/**
 * Apply sorting to query
 */
export function applySort(orderBy?: string, order?: "asc" | "desc"): SQL | undefined {
  if (!orderBy || orderBy.trim() === "") {
    return undefined;
  }

  const direction = order === "desc" ? sql`DESC` : sql`ASC`;
  return sql`${sql.identifier(orderBy)} ${direction}`;
}

/**
 * Get offset for pagination
 */
export function getOffset(pagination: Pagination): number {
  const page = pagination.page >= 1 ? pagination.page : 1;
  const perPage = pagination.perPage >= 1 && pagination.perPage <= 1000 ? pagination.perPage : 10;
  return (page - 1) * perPage;
}

/**
 * Get limit for pagination
 */
export function getLimit(pagination: Pagination): number {
  return pagination.perPage >= 1 && pagination.perPage <= 1000 ? pagination.perPage : 10;
}
