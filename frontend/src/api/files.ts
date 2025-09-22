// src/api/files.ts
export const API_BASE = import.meta.env.VITE_API_BASE ?? "http://localhost:8080";


// import { getJSON, postJSON } from "./index";

// export type FileMeta = {
//   id: number;        // user_files.id (used for delete)
//   file_id: number;   // files.id (used for download)
//   filename: string;  // user_files.file_name (shown to user)
// };

// // list files
// export async function listFiles(): Promise<{ files: FileMeta[] }> {
//   return getJSON("/api/files");
// }

// // upload file
// export async function uploadFile(file: File): Promise<FileMeta> {
//   const formData = new FormData();
//   formData.append("file", file);

//   const res = await fetch(`${API_BASE}/api/upload`, {
//     method: "POST",
//     credentials: "include", // so JWT cookie is sent
//     body: formData,
//   });

//   if (!res.ok) {
//     const text = await res.text();
//     throw new Error(text || "File upload failed");
//   }

//   return res.json();
// }

// // delete file
// export async function deleteFile(fileId: number): Promise<{ success: boolean }> {
//   return postJSON(`/api/files/${fileId}/delete`);
// }


export type FileMeta = {
  id: number;       // user_files.id
  file_id: number;  // reference to files.id
  filename: string;
  size?: number;
  uploader?: string;  
  uploadDate?: string;
  isPublic?: boolean;
  downloadCount?: number;
  dedupRefCount?: number;
  canDelete?: boolean;  
  tags?: string[];  
};

// List all files for the logged-in user
export async function listFiles(): Promise<{ files: FileMeta[] }> {
  const res = await fetch(`${API_BASE}/api/files`, {
    credentials: "include",
    headers: { "Accept": "application/json" },
  });
  if (!res.ok) throw new Error("Failed to fetch files");
  return res.json();
}

// Upload a single file
export async function uploadFile(file: File): Promise<FileMeta> {
  const formData = new FormData();
  formData.append("file", file);

  const res = await fetch(`${API_BASE}/api/upload`, {
    method: "POST",
    body: formData,
    credentials: "include",
  });

  if (!res.ok) {
    const data = await res.json().catch(() => null);
    throw new Error(data?.error || "Upload failed");
  }

  return res.json();
}

// Delete a file
export async function deleteFile(fileID: number): Promise<void> {
  const res = await fetch(`${API_BASE}/api/files/${fileID}`, {
    method: "DELETE",
    credentials: "include",
  });
  if (!res.ok) {
    const data = await res.json().catch(() => null);
    throw new Error(data?.error || "Failed to delete file");
  }
}

// Toggle visibility (public/private) for a file
export async function toggleVisibility(fileID: number, isPublic: boolean): Promise<FileMeta> {
  const res = await fetch(`${API_BASE}/api/files/${fileID}/visibility`, {
    method: "PATCH",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ isPublic }),
  });

  if (!res.ok) {
    const data = await res.json().catch(() => null);
    throw new Error(data?.error || "Failed to toggle visibility");
  }

  return res.json();
}

export type FileFilter = {
  filename?: string;
  mimeType?: string[];
  minSize?: number; // in bytes
  maxSize?: number; // in bytes
  startDate?: string; 
  endDate?: string;   
  uploader?: string;
  isPublic?: boolean;
  sizeRange?: [number, number];   
  dateRange?: [string, string]; 
};


// Filter files with query parameters
export async function filterFiles(filters: FileFilter): Promise<{ files: FileMeta[] }> {
  const query = new URLSearchParams();

  if (filters.filename) query.append("filename", filters.filename);
  if (filters.mimeType) query.append("mimeType", filters.mimeType.join(","));
  if (filters.minSize) query.append("minSize", filters.minSize.toString());
  if (filters.maxSize) query.append("maxSize", filters.maxSize.toString());
  if (filters.startDate) query.append("startDate", filters.startDate);
  if (filters.endDate) query.append("endDate", filters.endDate);
  if (filters.uploader) query.append("uploader", filters.uploader);
  if (filters.isPublic !== undefined) query.append("isPublic", filters.isPublic ? "true" : "false");

  const res = await fetch(`${API_BASE}/api/files/filter?${query.toString()}`, {
    credentials: "include",
    headers: { "Accept": "application/json" },
  });

  if (!res.ok) throw new Error("Failed to filter files");
  return res.json();
}
