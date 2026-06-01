import { cn, formatSidebarDate, getQualityClass } from "@/lib/utils";
import type { DailyEntry } from "@/types";

interface SidebarRowProps {
  entry: DailyEntry;
  active: boolean;
  onClick: () => void;
}

export function SidebarRow({ entry, active, onClick }: SidebarRowProps) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        "w-full flex items-center gap-2.5 px-2.5 py-1 rounded-lg text-left transition-colors",
        active ? "bg-muted border border-border/60" : "border border-transparent hover:bg-muted/60"
      )}
    >
      <span className={cn("w-2 h-2 rounded-full shrink-0", getQualityClass(entry.day_quality, "base"))} />
      <div className="flex-1 min-w-0">
        <div className="text-sm">{formatSidebarDate(entry.date)}</div>
        {entry.moods.length > 0 && (
          <div className="text-[11px] text-muted-foreground mt-0.5 truncate leading-3 mb-1">
            {entry.moods.slice(0, 3).join(", ")}
          </div>
        )}
      </div>
      <span className="text-[11px] text-muted-foreground shrink-0">{entry.day_quality}/10</span>
    </button>
  );
}
