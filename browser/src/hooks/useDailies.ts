import { useState, useEffect, useCallback } from 'react';
import type { DailyEntry, DailiesConfig } from '@/types';
import { api } from '../lib/api';

export function useDailies() {
  const [entries, setEntries] = useState<DailyEntry[]>([]);
  const [config, setConfig] = useState<DailiesConfig | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refreshData = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const [fetchedEntries, fetchedConfig] = await Promise.all([
        api.getEntries(),
        api.getConfig()
      ]);
      setEntries(fetchedEntries);
      setConfig(fetchedConfig);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An unknown error occurred');
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    let isMounted = true;

    async function loadInitialData() {
      try {
        const [fetchedEntries, fetchedConfig] = await Promise.all([
          api.getEntries(),
          api.getConfig()
        ]);
        
        if (isMounted) {
          setEntries(fetchedEntries);
          setConfig(fetchedConfig);
          setIsLoading(false);
        }
      } catch (err) {
        if (isMounted) {
          setError(err instanceof Error ? err.message : 'An unknown error occurred');
          setIsLoading(false);
        }
      }
    }

    loadInitialData();

    return () => {
      isMounted = false;
    };
  }, []);

  const saveEntry = async (entry: DailyEntry) => {
    try {
      await api.saveEntry(entry);
      await refreshData(); 
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save entry');
      throw err; 
    }
  };

  const runIntegration = async (date: string, integrationName: string) => {
    try {
      await api.triggerIntegration(date, integrationName);
      await refreshData(); 
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Integration sync failed');
      throw err;
    }
  };

  return {
    entries,
    config,
    isLoading,
    error,
    saveEntry,
    runIntegration,
    refreshData,
  };
}
