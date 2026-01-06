export interface Pagination {
  page: number;
  perPage: number;
  total: number;
  orderBy?: string;
  order?: "asc" | "desc";
}

export function createPagination(
  page?: number,
  perPage?: number,
  orderBy?: string,
  order?: string
): Pagination {
  return {
    page: page && page >= 1 ? page : 1,
    perPage: perPage && perPage >= 1 && perPage <= 1000 ? perPage : 10,
    total: 0,
    orderBy: orderBy || undefined,
    order: order === "desc" ? "desc" : "asc",
  };
}

export function getPaginationString(pagination: Pagination): string {
  return `page=${pagination.page}&per_page=${pagination.perPage}&order_by=${pagination.orderBy || ""}&order=${pagination.order || "asc"}`;
}
