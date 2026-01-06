import { HTTPException } from "hono/http-exception";
import { ContentfulStatusCode } from "hono/utils/http-status";

export interface ResponseError {
  code: number;
  source: string;
  title: string;
  message: string;
}

export interface ResponseForm<T = unknown> {
  success: boolean;
  data?: T;
  errors?: ResponseError[];
  result?: Record<string, unknown>;
}

export function createError(
  code: number,
  source: string,
  title: string,
  message: string
): ResponseError {
  return { code, source, title, message };
}

export function whereAmI(): string {
  const err = new Error();
  const stack = err.stack?.split("\n")[2];
  if (stack) {
    const match = stack.match(/at\s+(.+?)\s+\((.+?):(\d+):(\d+)\)/);
    if (match) {
      return `${match[1]} (${match[2]}:${match[3]})`;
    }
    const simpleMatch = stack.match(/at\s+(.+?):(\d+):(\d+)/);
    if (simpleMatch) {
      return `${simpleMatch[1]}:${simpleMatch[2]}`;
    }
  }
  return "unknown";
}

export function createHttpException(code: ContentfulStatusCode, message: string): HTTPException {
  return new HTTPException(code, { message });
}

export function successResponse<T>(data: T, result?: Record<string, unknown>): ResponseForm<T> {
  return {
    success: true,
    data,
    result,
  };
}

export function errorResponse(errors: ResponseError[]): ResponseForm {
  return {
    success: false,
    errors,
  };
}
