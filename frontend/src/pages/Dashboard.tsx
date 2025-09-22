// Dashboard.tsx
import React, { useEffect, useState, useCallback } from "react";
import { useAuth } from "../contexts/AuthContext";
import type { FileMeta, FileFilter } from "../api/files";
import { FileUpload } from "../components/files/FileUpload";
import { FileList } from "../components/files/FileList";
import { FileFilters } from "../components/files/FileFilters";
import { filterFiles } from "../api/files"; // backend API call
import "./dashboard.css"; 
import { useStorageStats } from "../hooks/userStorageStats";

const Dashboard: React.FC = () => {
  const { user, logout } = useAuth();
  const [files, setFiles] = useState<FileMeta[]>([]);
  const [filters, setFilters] = useState<FileFilter>({});
  const [loading, setLoading] = useState(true);
  const { stats: storageStats, refresh: refreshStorage } = useStorageStats();

  // Fetch filtered files from backend
  const fetchFilteredFiles = useCallback(async (filters: FileFilter) => {
    setLoading(true);
    try {
      const res = await filterFiles(filters);
      setFiles(res.files);
    } catch (err) {
      console.error("Failed to fetch files:", err);
      setFiles([]);
    } finally {
      setLoading(false);
    }
  }, []);

  // Refresh both files and storage stats
  const refreshData = async () => {
    await fetchFilteredFiles(filters);
    await refreshStorage();
  };

  // Fetch initially
  useEffect(() => {
    refreshData();
  }, []);

  // Refetch whenever filters change
  useEffect(() => {
    const timeout = setTimeout(() => fetchFilteredFiles(filters), 300); // debounce
    return () => clearTimeout(timeout);
  }, [filters, fetchFilteredFiles]);

  if (!user) return <p>Not authorized</p>;

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <h1>Welcome, {user.username} 👋</h1>
        <button onClick={logout} className="logout-btn">Logout</button>
      </header>

      <section className="stats">
        <p><strong>Total Storage (Deduped):</strong> {storageStats.total_storage} bytes</p>
        <p><strong>Total Original Storage:</strong> {storageStats.original_storage} bytes</p>
        <p><strong>Saved:</strong> {storageStats.savings} bytes</p>
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
              files={files}
              onDeleted={refreshData}
              onUpdated={refreshData}
            />
          </div>
        )}
      </section>

      <section>
        <h2>Upload New Files</h2>
        <FileUpload onUploaded={refreshData} />
      </section>
    </div>
  );
};

export default Dashboard;
