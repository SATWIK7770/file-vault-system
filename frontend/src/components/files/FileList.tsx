import React from "react";
import { deleteFile, API_BASE } from "../../api/files";
import type { FileMeta } from "../../api/files";

type Props = {
  files: FileMeta[];
  onDeleted: (id: number) => void;
};

export const FileList: React.FC<Props> = ({ files, onDeleted }) => {
  const handleDelete = async (id: number) => {
    try {
      await deleteFile(id);
      onDeleted(id);
    } catch (err) {
      console.error("Failed to delete file:", err);
      alert("Failed to delete file");
    }
  };

  if (!files || files.length === 0) {
    return <p>No files uploaded yet.</p>;
  }

  return (
    <ul>
      {files.map((f, idx) => (
        <li key={`${f.id}-${idx}`} style={{ marginBottom: "0.5rem" }}>
          <span>{f.filename}</span>
          <a
            href={`${API_BASE}/api/files/${f.file_id}/download`}
            style={{ marginLeft: "1rem" }}
          >
            Download
          </a>

          <button
            style={{ marginLeft: "1rem" }}
            onClick={() => handleDelete(f.id)}
          >
            Delete
          </button>
        </li>
      ))}
    </ul>
  );
};
