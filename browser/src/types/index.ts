export interface WordBank {
  promotion_threshold: number;
  max_words: number;
  moods: string[];
  context_tags: string[];
}

export interface DailiesConfig {
  word_bank: WordBank;
  integrations: Record<string, Record<string, boolean>>;
}

export interface IntegrationPayload {
  // integration payloads are intentionally generic to reduce friction
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  data: Record<string, any>;
  fetched_at: string | null;
}

export interface DailyEntry {
  date: string;         // YYYY-MM-DD
  day_quality: number;  // e.g., 1-5
  moods: string[];
  context_tags: string[];
  journal: string;
  integrations: Record<string, IntegrationPayload>;
}
