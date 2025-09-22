import React, { useState, useEffect } from "react";
import type { FileFilter } from "../../api/files";
import "./filefilters.css";

type Props = {
  onChange: (filters: FileFilter) => void;
};

export const FileFilters: React.FC<Props> = ({ onChange }) => {
  const [filename, setFilename] = useState("");
  const [mimeType, setMimeType] = useState("");
  const [uploader, setUploader] = useState("");
  const [minSize, setMinSize] = useState("");
  const [maxSize, setMaxSize] = useState("");
  const [fromDate, setFromDate] = useState("");
  const [toDate, setToDate] = useState("");

  useEffect(() => {
    const filter: FileFilter = {};

    if (filename) filter.filename = filename;
    if (mimeType) filter.mimeType = mimeType.split(",").map((m) => m.trim());
    if (uploader) filter.uploader = uploader;

    if (minSize || maxSize) {
      const min = minSize ? parseInt(minSize) * 1024 : 0;
      const max = maxSize ? parseInt(maxSize) * 1024 : Infinity;
      filter.sizeRange = [min, max];
      filter.minSize = min;
      filter.maxSize = max;
    }

    if (fromDate || toDate) {
      const start = fromDate || "1970-01-01";
      const end = toDate || new Date().toISOString().split("T")[0];
      filter.dateRange = [start, end];
      filter.startDate = start;
      filter.endDate = end;
    }

    // for smoother filtering
    const timeout = setTimeout(() => onChange(filter), 300);
    return () => clearTimeout(timeout);
  }, [filename, mimeType, uploader, minSize, maxSize, fromDate, toDate, onChange]);

  return (
<div className="filters-container">
  <input
    type="text"
    placeholder="Search filename"
    value={filename}
    onChange={(e) => setFilename(e.target.value)}
  />
<label> MIME Type:
  <select value={mimeType} onChange={(e) => setMimeType(e.target.value)}>
    <option value="">All</option>
    <option value="image/jpeg">Image (JPEG)</option>
    <option value="image/png">Image (PNG)</option>
    <option value="image/gif">Image (GIF)</option>
    <option value="application/pdf">PDF</option>
    <option value="video/mp4">Video (MP4)</option>
    <option value="video/webm">Video (WebM)</option>
    <option value="audio/mpeg">Audio (MP3)</option>
    <option value="audio/wav">Audio (WAV)</option>
    <option value="application/msword">Word Document (.doc)</option>
    <option value="application/vnd.openxmlformats-officedocument.wordprocessingml.document">
      Word Document (.docx)
    </option>
    </select>
    </label>
  <input
    type="text"
    placeholder="Uploader"
    value={uploader}
    onChange={(e) => setUploader(e.target.value)}
  />

  <div className="range-group">
    <input
      type="number"
      placeholder="Min KB"
      value={minSize}
      onChange={(e) => setMinSize(e.target.value)}
    />
    <span>—</span>
    <input
      type="number"
      placeholder="Max KB"
      value={maxSize}
      onChange={(e) => setMaxSize(e.target.value)}
    />
  </div>

  {/* Date Range */}
  <div className="range-group">
    <input
      type="date"
      value={fromDate}
      onChange={(e) => setFromDate(e.target.value)}
    />
    <span>—</span>
    <input
      type="date"
      value={toDate}
      onChange={(e) => setToDate(e.target.value)}
    />
  </div>
</div>

);
};