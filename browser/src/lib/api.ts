import type { DailyEntry, DailiesConfig } from '@/types';

const API_BASE = '/api';

export const api = {
  async getConfig(): Promise<DailiesConfig> {
    const res = await fetch(`${API_BASE}/config`);
    if (!res.ok) throw new Error('Failed to fetch config');
    return res.json();
  },

  async getEntries(): Promise<DailyEntry[]> {
    const res = await fetch(`${API_BASE}/entries`);
    if (!res.ok) throw new Error('Failed to fetch entries');
    return res.json();
  },

  async saveEntry(entry: DailyEntry): Promise<void> {
    const res = await fetch(`${API_BASE}/entries`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(entry),
    });
    if (!res.ok) throw new Error('Failed to save entry');
  },

  async triggerIntegration(date: string, integration: string): Promise<void> {
    const res = await fetch(`${API_BASE}/entries/${date}/fetch/${integration}`, {
      method: 'POST',
    });
    if (!res.ok) throw new Error(`Failed to fetch integration: ${integration}`);
  }
};
