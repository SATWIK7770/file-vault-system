import React, { useEffect, useState } from "react";
import { useAuth } from "../contexts/AuthContext";
import { listFiles } from "../api/files";
import type { FileMeta } from "../api/files";
import { FileUpload } from "../components/files/FileUpload";
import { FileList } from "../components/files/FileList";

const Dashboard: React.FC = () => {
const { user, logout } = useAuth();
const [files, setFiles] = useState<FileMeta[]>([]);
const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchFiles() {
      try {
        const res = await listFiles();
        setFiles(res.files);
      } catch (err) {
        console.error("Failed to fetch files:", err);
      } finally {
        setLoading(false);
      }
    }
    fetchFiles();
  }, []);

  if (!user) return <p>Not authorized</p>;

  return (
    <div style={{ padding: "2rem" }}>
      <h1>Welcome, {user.username} 👋</h1>
      <button onClick={logout}>Logout</button>

      <hr />

      <h2>Your Files</h2>
      {loading ? (
        <p>Loading files...</p>
      ) : (
        <FileList
          files={files}
          onDeleted={(id) =>
            setFiles((prev) => prev.filter((file) => file.id !== id))
          }
        />
      )}

      <hr />

      <h2>Upload a new file</h2>
      <FileUpload onUploaded={(file) => setFiles((prev) => [...prev, file])} />
    </div>
  );
};

export default Dashboard;
