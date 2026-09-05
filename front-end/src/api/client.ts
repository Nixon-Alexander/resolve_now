import type { ErrorResponse, SuccessResponse } from "../types";

// Configure via .env: VITE_API_BASE_URL=http://localhost:8080/api/v1
export const API_BASE_URL: string =
  import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8000/api/v1";

export class ApiError extends Error {
  code?: number;
  constructor(message: string, code?: number) {
    super(message);
    this.name = "ApiError";
    this.code = code;
  }
}

async function unwrap<T>(res: Response): Promise<T> {
  let body: unknown = null;
  try {
    body = await res.json();
  } catch {
    // no body / not JSON
  }

  if (!res.ok) {
    const err = body as ErrorResponse | null;
    throw new ApiError(err?.message ?? `Request failed with status ${res.status}`, err?.error_code);
  }

  const success = body as SuccessResponse<T> | T;
  if (success && typeof success === "object" && "data" in (success as Record<string, unknown>)) {
    return (success as SuccessResponse<T>).data;
  }
  return success as T;
}

export async function apiGet<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    method: "GET",
    headers: { Accept: "application/json" },
  });
  return unwrap<T>(res);
}

export async function apiPostJson<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json", Accept: "application/json" },
    body: JSON.stringify(body),
  });
  return unwrap<T>(res);
}

export async function apiPostForm<T>(path: string, form: FormData): Promise<T> {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    method: "POST",
    headers: { Accept: "application/json" },
    body: form,
  });
  return unwrap<T>(res);
}
