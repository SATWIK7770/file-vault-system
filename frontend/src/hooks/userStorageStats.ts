// src/hooks/useStorageStats.ts
import { useEffect, useState } from "react";
import { getStorageStats, type StorageStats } from "../api/storage";

export const useStorageStats = (pollInterval = 0) => {
  const [stats, setStats] = useState<StorageStats>({
    total_storage: 0,
    original_storage: 0,
    savings: 0,
  });

  const fetchStats = async () => {
    try {
      const data = await getStorageStats();
      setStats(data);
    } catch (err) {
      console.error("Failed to fetch storage stats:", err);
    }
  };

  useEffect(() => {
    fetchStats();
    if (pollInterval > 0) {
      const interval = setInterval(fetchStats, pollInterval);
      return () => clearInterval(interval);
    }
  }, []);

  return { stats, refresh: fetchStats };
};
