import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function getLocalDateString(date: Date = new Date()): string {
  const offset = date.getTimezoneOffset();
  const localDate = new Date(date.getTime() - (offset * 60 * 1000));
  return localDate.toISOString().slice(0, 10);
}

export function formatDisplayDate(dateStr: string): string {
  const d = new Date(dateStr + "T00:00:00");
  return d.toLocaleDateString("en-US", {
    weekday: "short",
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

export function formatSidebarDate(dateStr: string): string {
  const TODAY = getLocalDateString();
  const d = new Date(dateStr + "T00:00:00");
  const now = new Date();
  const yesterday = new Date(now);
  yesterday.setDate(now.getDate() - 1);

  if (dateStr === TODAY) return "Today";
  if (dateStr === yesterday.toISOString().slice(0, 10)) return "Yesterday";
  return d.toLocaleDateString("en-US", { weekday: "short", month: "short", day: "numeric" });
}

export function qualityColor(q: number): string {
  if (q >= 7) return "bg-emerald-500";
  if (q >= 4) return "bg-amber-400";
  return "bg-red-400";
}

const QUALITY_THEMES = {
  exceptional: {
    base: "bg-purple-500 border-purple-500 text-white",
    active: "bg-purple-600 border-purple-600 text-white",
    hover: "hover:bg-purple-600 hover:border-purple-600",
  },
  high: {
    base: "bg-emerald-500 border-emerald-500 text-white",
    active: "bg-emerald-600 border-emerald-600 text-white",
    hover: "hover:bg-emerald-600 hover:border-emerald-600",
  },
  medium: {
    base: "bg-amber-400 border-amber-400 text-amber-950",
    active: "bg-amber-500 border-amber-500 text-amber-950",
    hover: "hover:bg-amber-500 hover:border-amber-500",
  },
  low: {
    base: "bg-red-400 border-red-400 text-white",
    active: "bg-red-500 border-red-500 text-white",
    hover: "hover:bg-red-500 hover:border-red-500",
  },
};

function getQualityKey(q: number): keyof typeof QUALITY_THEMES {
  if (q === 10) return "exceptional";
  if (q >= 7) return "high";
  if (q >= 4) return "medium";
  return "low";
}

export function getQualityClass(q: number, variant: "base" | "active" | "hover"): string {
  return QUALITY_THEMES[getQualityKey(q)][variant];
}
