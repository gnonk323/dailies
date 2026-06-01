import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import type { DailyEntry } from "@/types";
import { Spark } from "./ui/Spark";

export function InsightsStrip({ entries }: { entries: DailyEntry[] }) {
  const last30 = entries.slice(0, 30);
  const prior30 = entries.slice(30, 60);

  const avg = (arr: DailyEntry[]) =>
    arr.length ? arr.reduce((s, e) => s + e.day_quality, 0) / arr.length : 0;

  const avgNow = avg(last30);
  const avgPrior = avg(prior30);
  const delta = avgNow - avgPrior;

  const moodCounts: Record<string, number> = {};
  last30.forEach((e) => e.moods.forEach((m) => (moodCounts[m] = (moodCounts[m] ?? 0) + 1)));
  const topMoods = Object.entries(moodCounts)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 4);

  const qualityValues = last30.map((e) => e.day_quality).reverse();

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium">Insights</span>
      </div>
      <div className="grid grid-cols-3 gap-2.5">
        <div className="border bg-card rounded-lg p-3">
          <div className="text-xl font-medium">{avgNow.toFixed(1)}</div>
          <div className="text-[11px] text-muted-foreground mt-0.5">avg quality, 30d</div>
          <div className="my-2">
            <Spark values={qualityValues} />
          </div>
          {prior30.length > 0 && (
            <div className={cn("text-[11px] mt-1", delta >= 0 ? "text-emerald-600 dark:text-emerald-400" : "text-red-500")}>
              {delta >= 0 ? "↑" : "↓"} {Math.abs(delta).toFixed(1)} vs prior 30d
            </div>
          )}
        </div>

        <div className="border bg-card rounded-lg p-3">
          <div className="text-xl font-medium">
            {topMoods[0]?.[0] ?? "—"}
          </div>
          <div className="text-[11px] text-muted-foreground mt-0.5">top mood, 30d</div>
          <div className="flex flex-wrap gap-1 mt-2">
            {topMoods.slice(1).map(([m]) => (
              <Badge key={m} variant="secondary" className="text-[10px] font-normal">
                {m}
              </Badge>
            ))}
          </div>
        </div>

        <div className="border bg-card rounded-lg p-3 flex flex-col justify-between">
          <div>
            <div className="text-xl font-medium">{last30.length}</div>
            <div className="text-[11px] text-muted-foreground mt-0.5">entries, 30d</div>
          </div>
          <div className="text-[11px] text-muted-foreground mt-1">
            {entries.length} total
          </div>
        </div>
      </div>
    </div>
  );
}
