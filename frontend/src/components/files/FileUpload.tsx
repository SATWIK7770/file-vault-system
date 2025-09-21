import React, { useState, useCallback, useRef } from "react";
import { useDropzone } from "react-dropzone";
import { uploadFile } from "../../api/files";
import type { FileMeta } from "../../api/files";

type Props = {
  onUploaded: (file: FileMeta) => void;
};

export const FileUpload: React.FC<Props> = ({ onUploaded }) => {
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
  const [loading, setLoading] = useState(false);
  const [err, setErr] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const addFiles = (files: File[]) => {
    setSelectedFiles((prev) => [...prev, ...files]);
  };

  const removeFile = (index: number) => {
    setSelectedFiles((prev) => prev.filter((_, i) => i !== index));
  };

  // Trigger actual upload on button click
  const handleUpload = async () => {
    if (selectedFiles.length === 0) return;

    setErr(null);
    setLoading(true);
    try {
      for (const file of selectedFiles) {
        const uploaded = await uploadFile(file);
        onUploaded(uploaded);
      }
      // Reset after upload
      setSelectedFiles([]);
      if (fileInputRef.current) {
        fileInputRef.current.value = ""; // ✅ clears file input
      }
    } catch (e: any) {
      setErr(e.message || "Upload failed");
    } finally {
      setLoading(false);
    }
  };

  // Drag & Drop
  const onDrop = useCallback((acceptedFiles: File[]) => {
    addFiles(acceptedFiles);
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({ onDrop });

  // Classic file picker
  const onFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      addFiles(Array.from(e.target.files));
      e.target.value = ""; // ✅ allows reselecting same file
    }
  };

  return (
    <div>
      {err && <div style={{ color: "red" }}>{err}</div>}

      <div
        {...getRootProps()}
        className="border-2 border-dashed border-gray-400 p-12 min-h-[200px] flex items-center justify-center text-lg text-gray-600 rounded-xl cursor-pointer mb-4"
      >
        <input {...getInputProps()} />
        {isDragActive ? <p>Drop files here...</p> : <p>Drag & drop files here</p>}
      </div>

      {/* Classic input */}
      <input ref={fileInputRef} type="file" multiple onChange={onFileChange} />

      {/* Show selected files */}
      {selectedFiles.length > 0 && (
        <ul className="mt-2 space-y-1">
          {selectedFiles.map((f, idx) => (
            <li
              key={idx}
              className="flex items-center justify-between bg-gray-100 px-2 py-1 rounded"
            >
              <span className="truncate">{f.name}</span>
              <button
                onClick={() => removeFile(idx)}
                className="ml-2 text-red-500 hover:text-red-700"
              >
                ✕
              </button>
            </li>
          ))}
        </ul>
      )}

      {/* Upload button */}
      <button
        onClick={handleUpload}
        disabled={loading || selectedFiles.length === 0}
        className="mt-3 px-4 py-2 bg-blue-500 text-white rounded disabled:opacity-50"
      >
        {loading ? "Uploading..." : "Upload Selected"}
      </button>
    </div>
  );
};
