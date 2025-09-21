// src/api/index.ts
export const API_BASE = import.meta.env.VITE_API_BASE ?? "http://localhost:8080";

async function parseJSON(response: Response) {
  const text = await response.text();
  try {
    return text ? JSON.parse(text) : null;
  } catch {
    return text;
  }
}

export async function postJSON(path: string, body?: unknown) {
  const opts: RequestInit = {
    method: "POST",
    credentials: "include", 
  };

  // ✅ Only add headers/body if body is provided
  if (body !== undefined) {
    opts.headers = { "Content-Type": "application/json" };
    opts.body = JSON.stringify(body);
  }

  const res = await fetch(`${API_BASE}${path}`, opts);

  const data = await parseJSON(res);
  if (!res.ok) {
    const msg = (data && (data.error || data.message)) || res.statusText;
    const err: any = new Error(msg);
    err.status = res.status;
    throw err;
  }
  return data;
}


export async function getJSON(path: string) {
  const res = await fetch(`${API_BASE}${path}`, {
    credentials: "include",
    headers: { "Accept": "application/json" },
  });
  const data = await parseJSON(res);
  if (!res.ok) {
    const msg = (data && (data.error || data.message)) || res.statusText;
    const err: any = new Error(msg);
    err.status = res.status;
    throw err;
  }
  return data;
}
