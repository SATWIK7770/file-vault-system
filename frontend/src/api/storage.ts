// src/api/storage.ts

export const API_BASE = import.meta.env.VITE_API_BASE ?? "http://localhost:8080";
export interface StorageStats {
  total_storage: number;
  original_storage: number;
  savings: number;
}


export async function getStorageStats(): Promise<StorageStats> {
  const res = await fetch(`${API_BASE}/api/storage-stats`, {
    method: "GET",
    credentials: "include",
    headers: {
      "Authorization": `Bearer ${localStorage.getItem("token") || ""}`,
    },
  });

  if (!res.ok) {
    const data = await res.json().catch(() => null);
    throw new Error(data?.error || "Failed to fetch storage stats");
  }

  return res.json();
}

