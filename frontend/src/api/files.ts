// src/api/files.ts
export const API_BASE = import.meta.env.VITE_API_BASE ?? "http://localhost:8080";


import { getJSON, postJSON } from "./index";

export type FileMeta = {
  id: number;
  filename: string;
  size: number;          // in bytes
  uploaded_at: string;   // ISO timestamp from backend
};

// list files
export async function listFiles(): Promise<{ files: FileMeta[] }> {
  return getJSON("/api/files");
}

// upload file
export async function uploadFile(file: File): Promise<FileMeta> {
  const formData = new FormData();
  formData.append("file", file);

  const res = await fetch(`${API_BASE}/api/upload`, {
    method: "POST",
    credentials: "include", // so JWT cookie is sent
    body: formData,
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || "File upload failed");
  }

  return res.json();
}

// delete file
export async function deleteFile(fileId: number): Promise<{ success: boolean }> {
  return postJSON(`/api/files/${fileId}/delete`);
}