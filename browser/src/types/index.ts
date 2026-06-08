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

export interface NYTData {
  user_id: string;
  connections: NYTConnections;
  crossword_midi: NYTCrossword;
  crossword_mini: NYTCrosswordMini;
  wordle: NYTWordle;
  summary: NYTSummary;
}

export interface NYTConnections {
  puzzle_id: number;
  game_data: {
    isPlayingArchive: boolean;
    mistakes: number;
    puzzleComplete: boolean;
    puzzleWon: boolean;
    guesses: ConnectionsGuess[];
    solvedCategories: ConnectionsSolvedCategory[];
  };
}

export interface ConnectionsCard {
  level: number;
  position: number;
}

export interface ConnectionsGuess {
  cards: ConnectionsCard[];
  correct: boolean;
}

export interface ConnectionsSolvedCategory {
  cards: ConnectionsCard[];
  level: number;
  orderSolved: number;
}

export interface CrosswordDimensions {
  height: number;
  width: number;
}

export interface NYTCrossword {
  id: number;
  publicationDate: string; // YYYY-MM-DD
  title: string;
  editor: string;
  constructors: string[];
  dimensions: CrosswordDimensions;
  completionFraction: number;
  firstSolve: number;
  playTime: number;
}

export interface NYTCrosswordMini {
  id: number;
  publicationDate: string;
  constructors: string[];
  dimensions: CrosswordDimensions;
  completionFraction: number;
  firstSolve: number;
  playTime: number;
}

export interface NYTWordle {
  puzzle_id: number;
  solution: string;
  game_data: {
    boardState: string[];
    currentRowIndex: number;
    hardMode: boolean;
    isPlayingArchive: boolean;
    status: 'IN_PROGRESS' | 'WIN' | 'FAIL'; // Narrowed based on typical Wordle state
  };
}

export interface NYTSummary {
  connections: {
    current_streak: number;
    last_played_print_date: string;
    max_streak: number;
    mistakes: Record<string, number>; // Dynamic map, e.g., {"0": 67, "1": 36...}
    puzzles_completed: number;
    puzzles_won: number;
  };
  crossword_midi: {
    generation: number;
    puzzlesSolved: number;
    puzzlesStarted: number;
    solveRate: number;
    streaks: {
      current: number;
      longest: number;
      startDate: string;
    };
  };
  crossword_mini: {
    bestDate: string;
    bestTimeSeconds: number;
    generation: number;
  };
  wordle: {
    calculatedStats: {
      currentStreak: number;
      hasPlayed: boolean;
      lastCompletedPrintDate: string;
      lastWonPrintDate: string;
      maxStreak: number;
    };
    legacyStats: {
      autoOptInTimestamp: number;
      currentStreak: number;
      gamesPlayed: number;
      gamesWon: number;
      guesses: Record<string, number>; // Maps guess counts "1"-"6" and "fail"
      hasMadeStatsChoice: boolean;
      hasPlayed: boolean;
      lastWonDayOffset: number;
      maxStreak: number;
      timestamp: number;
    };
    seaOfGreensStats: {
      sea_of_greens_wins: number;
    };
    totalStats: {
      gamesPlayed: number;
      gamesWon: number;
      guesses: Record<string, number>;
      hasPlayed: boolean;
      hasPlayedArchive: boolean;
    };
  };
}