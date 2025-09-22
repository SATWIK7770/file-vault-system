// import React, { useEffect, useState } from "react";
// import { useAuth } from "../contexts/AuthContext";
// import { listFiles } from "../api/files";
// import type { FileMeta, FileFilter } from "../api/files";
// import { FileUpload } from "../components/files/FileUpload";
// import { FileList } from "../components/files/FileList";
// import { FileFilters } from "../components/files/FileFilters";

// const Dashboard: React.FC = () => {
//   const { user, logout } = useAuth();
//   const [files, setFiles] = useState<FileMeta[]>([]);
//   const [filters, setFilters] = useState<FileFilter>({});
//   const [loading, setLoading] = useState(true);

//   useEffect(() => {
//     async function fetchFiles() {
//       try {
//         const res = await listFiles();
//         setFiles(res.files);
//       } catch (err) {
//         console.error("Failed to fetch files:", err);
//       } finally {
//         setLoading(false);
//       }
//     }
//     fetchFiles();
//   }, []);

//   if (!user) return <p>Not authorized</p>;

//   // Apply filters
//   const filteredFiles = files.filter((f) => {
//     if (filters.filename && !f.filename.includes(filters.filename)) return false;
//     if (filters.uploader && f.uploader !== filters.uploader) return false;
//     if (filters.minSize && (f.size === undefined || f.size < filters.minSize)) return false;
//     if (filters.maxSize && (f.size === undefined || f.size > filters.maxSize)) return false;
//     if (filters.startDate && f.uploadDate && new Date(f.uploadDate) < new Date(filters.startDate)) return false;
//     if (filters.endDate && f.uploadDate && new Date(f.uploadDate) > new Date(filters.endDate)) return false;
//     if (filters.tags && !filters.tags.every((t) => f.tags?.includes(t))) return false;
//     if (filters.isPublic !== undefined && f.isPublic !== filters.isPublic) return false;
//     return true;
//   });

//   // Compute storage stats
//   const totalOriginal = files.reduce((a, f) => a + (f.size || 0), 0);
//   const totalDeduped = files.reduce(
//     (a, f) => a + (f.size || 0) / (f.dedupRefCount && f.dedupRefCount > 0 ? f.dedupRefCount : 1),
//     0
//   );
//   const storageSaved = totalOriginal - totalDeduped;
//   const storageSavedPercent = totalOriginal ? (storageSaved / totalOriginal) * 100 : 0;

//   return (
//     <div style={{ padding: "2rem" }}>
//       <h1>Welcome, {user.username} 👋</h1>
//       <button onClick={logout} className="mb-4 px-3 py-1 bg-gray-300 rounded">
//         Logout
//       </button>

//       <div className="mb-4 p-2 border rounded">
//         <p>Total Storage (Deduped): {totalDeduped.toFixed(2)} bytes</p>
//         <p>Total Original Storage: {totalOriginal} bytes</p>
//         <p>Saved: {storageSaved.toFixed(2)} bytes ({storageSavedPercent.toFixed(2)}%)</p>
//       </div>

//       <h2>Filters</h2>
//       <FileFilters onChange={setFilters} />

//       <h2>Your Files</h2>
//       {loading ? (
//         <p>Loading files...</p>
//       ) : (
//         <FileList files={filteredFiles} onDeleted={(id) => setFiles(prev => prev.filter(f => f.id !== id))} onUpdated={(updated) =>
//         setFiles(prev => prev.map(f => f.id === updated.id ? updated : f))}/>
//       )}
//       <hr className="my-4" />

//       <h2>Upload New Files</h2>
//       <FileUpload onUploaded={(file) => setFiles((prev) => [...prev, file])} />
//     </div>
//   );
// };

// export default Dashboard;


import React, { useEffect, useState } from "react";
import { useAuth } from "../contexts/AuthContext";
import { listFiles } from "../api/files";
import type { FileMeta, FileFilter } from "../api/files";
import { FileUpload } from "../components/files/FileUpload";
import { FileList } from "../components/files/FileList";
import { FileFilters } from "../components/files/FileFilters";
import "./dashboard.css"; 

const Dashboard: React.FC = () => {
  const { user, logout } = useAuth();
  const [files, setFiles] = useState<FileMeta[]>([]);
  const [filters, setFilters] = useState<FileFilter>({});
  const [loading, setLoading] = useState(true);

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

  // Storage stats
  const totalOriginal = files.reduce((a, f) => a + (f.size || 0), 0);
  const totalDeduped = files.reduce(
    (a, f) => a + (f.size || 0) / (f.dedupRefCount && f.dedupRefCount > 0 ? f.dedupRefCount : 1),
    0
  );
  const storageSaved = totalOriginal - totalDeduped;
  const storageSavedPercent = totalOriginal ? (storageSaved / totalOriginal) * 100 : 0;

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <h1>Welcome, {user.username} 👋</h1>
        <button onClick={logout} className="logout-btn">Logout</button>
      </header>

      <section className="stats">
        <p><strong>Total Storage (Deduped):</strong> {totalDeduped} bytes</p>
        <p><strong>Total Original Storage:</strong> {totalOriginal} bytes</p>
        <p>
          <strong>Saved:</strong> {storageSaved} bytes 
          ({storageSavedPercent.toFixed(2)}%)
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
              // Refresh entire file list after any operation
              onDeleted={async () => await fetchFiles()}
              onUpdated={async () => await fetchFiles()}
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

