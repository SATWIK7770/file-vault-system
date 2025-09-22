import React from "react";
import { deleteFile, toggleVisibility, API_BASE } from "../../api/files";
import type { FileMeta } from "../../api/files";
import "./filelist.css";

type Props = {
  files: FileMeta[];
  onDeleted: (id: number) => void;
  onUpdated: (updatedFile: FileMeta) => void;
  canDelete?: (file: FileMeta) => boolean;
};

export const FileList: React.FC<Props> = ({ files, onDeleted, onUpdated, canDelete }) => {
  const handleDelete = async (file: FileMeta) => {
    if (canDelete && !canDelete(file)) return alert("Cannot delete this file");
    try {
      await deleteFile(file.id);
      onDeleted(file.id);
    } catch (err) {
      console.error(err);
      alert("Delete failed");
    }
  };

  const handleToggleVisibility = async (file: FileMeta) => {
    try {
      const updated = await toggleVisibility(file.id, !file.isPublic);
      onUpdated(updated);
    } catch (err) {
      console.error(err);
      alert("Failed to update visibility");
    }
  };

  if (!files.length) return <p>No files found.</p>;

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm border-collapse border border-gray-300">
        <thead>
          <tr className="bg-gray-100">
            <th className="border p-2 w-1/5">Filename</th>
            <th className="border p-2 w-1/10">Size</th>
            <th className="border p-2 w-1/10">Uploader</th>
            <th className="border p-2 w-1/10">Upload Date</th>
            <th className="border p-2 w-1/12">Public</th>
            <th className="border p-2 w-1/12">Downloads</th>
            <th className="border p-2 w-1/12">Dedup Refs</th>
            <th className="border p-2 w-1/6">Tags</th>
            <th className="border p-2 w-1/6">Actions</th>
          </tr>
        </thead>
        <tbody>
          {files.map((f) => (
            <tr key={f.id} className="hover:bg-gray-50">
              <td className="border p-2">{f.filename}</td>
              <td className="border p-2">{f.size ? (f.size / 1024).toFixed(2) + " KB" : "-"}</td>
              <td className="border p-2">{f.uploader || "-"}</td>
              <td className="border p-2">{f.uploadDate ? new Date(f.uploadDate).toLocaleDateString() : "-"}</td>
              <td className="border p-2">{f.isPublic ? "Yes" : "No"}</td>
              <td className="border p-2">{f.downloadCount ?? "-"}</td>
              <td className="border p-2">{f.dedupRefCount ?? "-"}</td>
              <td className="border p-2">{f.tags?.join(", ") || "-"}</td>
              <td className="border p-2 space-x-2">
                <a href={`${API_BASE}/api/files/${f.file_id}/download`} className="text-blue-600 text-xs px-2 py-1 border rounded">
                  Download
                </a>
                <button onClick={() => handleToggleVisibility(f)} className="text-white bg-blue-600 text-xs px-2 py-1 border rounded">
                {f.isPublic ? "Make Private" : "Make Public"} </button>
                {f.canDelete && (
                  <button
                    onClick={() => handleDelete(f)}
                    className="text-red-600 text-xs px-2 py-1 border rounded"
                  >
                    Delete
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};