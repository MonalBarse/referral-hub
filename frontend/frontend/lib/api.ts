import type { Comment, Job, Referral, User } from "@/lib/types";

export const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8090/api/v1"\;

export const WS_URL =
  process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8090/ws";

const TOKEN_KEY = "referhub.token";

export const tokenStore = {
  get(): string | null {
    if (typeof window === "undefined") return null;
    try {
      return window.localStorage.getItem(TOKEN_KEY);
    } catch {
      return null;
    }
  },
  set(token: string) {
    try {
      window.localStorage.setItem(TOKEN_KEY, token);
    } catch {
      // Private browsing can block storage. The session still works in memory.
    }
  },
  clear() {
    try {
      window.localStorage.removeItem(TOKEN_KEY);
    } catch {
      // Nothing to do.
    }
  },
};

// ApiError carries the backend's status and per-field validation messages so
// forms can mark the offending inputs.
export class ApiError extends Error {
  status: number;
  fields?: Record<string, string>;

  constructor(status: number, message: string, fields?: Record<string, string>) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.fields = fields;
  }
}

async function request<T>(
  path: string,
  init: RequestInit = {},
  token?: string | null,
): Promise<T> {
  const headers = new Headers(init.headers);
  headers.set("Content-Type", "application/json");
  if (token) headers.set("Authorization", `Bearer ${token}`);

  let res: Response;
  try {
    res = await fetch(`${API_URL}${path}`, { ...init, headers });
  } catch {
    throw new ApiError(0, "Cannot reach the API. Is the backend running?");
  }

  if (res.status === 204) return undefined as T;

  const body = await res.json().catch(() => null);

  if (!res.ok) {
    const message =
      (body && typeof body.error === "string" && body.error) ||
      `Request failed with status ${res.status}`;
    throw new ApiError(res.status, message, body?.fields);
  }

  return body as T;
}

export const api = {
  login: (email: string, name: string) =>
    request<{ token: string; user: User }>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, name }),
    }),

  me: (token: string) => request<User>("/auth/me", {}, token),

  listJobs: () => request<{ jobs: Job[] }>("/jobs"),

  getJob: (jobId: string) =>
    request<{ job: Job; referrals: Referral[] }>(`/jobs/${jobId}`),

  createJob: (
    token: string,
    input: { title: string; company: string; description: string; status?: string },
  ) =>
    request<Job>("/jobs", { method: "POST", body: JSON.stringify(input) }, token),

  createReferral: (
    token: string,
    jobId: string,
    input: { candidateName: string; candidateEmail: string; message: string },
  ) =>
    request<Referral>(
      `/jobs/${jobId}/referrals`,
      { method: "POST", body: JSON.stringify(input) },
      token,
    ),

  listComments: (referralId: string) =>
    request<{ comments: Comment[] }>(`/referrals/${referralId}/comments`),

  createComment: (token: string, referralId: string, body: string) =>
    request<Comment>(
      `/referrals/${referralId}/comments`,
      { method: "POST", body: JSON.stringify({ body }) },
      token,
    ),
};
