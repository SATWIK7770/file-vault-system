import React, { useEffect, useState } from "react";
import { useAuth } from "../contexts/AuthContext";
import { listFiles } from "../api/files";
import type { FileMeta, FileFilter } from "../api/files";
import { FileUpload } from "../components/files/FileUpload";
import { FileList } from "../components/files/FileList";
import { FileFilters } from "../components/files/FileFilters";
import "./dashboard.css"; 
import { useStorageStats } from "../hooks/userStorageStats";

const Dashboard: React.FC = () => {
  const { user, logout } = useAuth();
  const [files, setFiles] = useState<FileMeta[]>([]);
  const [filters, setFilters] = useState<FileFilter>({});
  const [loading, setLoading] = useState(true);
  const { stats: storageStats, refresh: refreshStorage } = useStorageStats();

    const refreshData = async () => {
    await fetchFiles();       // existing file list refresh
    await refreshStorage();   // refresh storage stats
  };

  // Centralized fetch function
  const fetchFiles = async () => {
    setLoading(true);
    try {
      const res = await listFiles();
      setFiles(res.files);
    } catch (err) {
      console.error("Failed to fetch files:", err);
    } finally {
      setLoading(false);
    }
  };

  // Initial fetch
  useEffect(() => {
    fetchFiles();
    refreshStorage();
  }, []);

  if (!user) return <p>Not authorized</p>;

  // Apply filters
  const filteredFiles = files.filter((f) => {
    if (filters.filename && !f.filename.includes(filters.filename)) return false;
    if (filters.uploader && f.uploader !== filters.uploader) return false;
    if (filters.minSize && (f.size === undefined || f.size < filters.minSize)) return false;
    if (filters.maxSize && (f.size === undefined || f.size > filters.maxSize)) return false;
    if (filters.startDate && f.uploadDate && new Date(f.uploadDate) < new Date(filters.startDate)) return false;
    if (filters.endDate && f.uploadDate && new Date(f.uploadDate) > new Date(filters.endDate)) return false;
    if (filters.isPublic !== undefined && f.isPublic !== filters.isPublic) return false;
    return true;
  });

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <h1>Welcome, {user.username} 👋</h1>
        <button onClick={logout} className="logout-btn">Logout</button>
      </header>

      <section className="stats">
        <p><strong>Total Storage (Deduped):</strong> {storageStats.total_storage} bytes</p>
        <p><strong>Total Original Storage:</strong> {storageStats.original_storage} bytes</p>
        <p>
          <strong>Saved:</strong> {storageStats.savings} bytes
        </p>
      </section>


      <section>
        <h2>Filters</h2>
        <FileFilters onChange={setFilters} />
      </section>

      <section>
        <h2>Your Files</h2>
        {loading ? (
          <p>Loading files...</p>
        ) : (
          <div className="overflow-x-auto w-full">
          <FileList
            files={filteredFiles}
            onDeleted={refreshData}
            onUpdated={refreshData}
          />

          </div>
        )}
      </section>

      <section>
        <h2>Upload New Files</h2>
        <FileUpload onUploaded={async () => await fetchFiles()} />
      </section>
    </div>
  );
};

export default Dashboard;

