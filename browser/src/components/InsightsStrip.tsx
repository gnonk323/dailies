import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import type { DailyEntry, NYTSummary } from "@/types";
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

  // Find the most recent entry that has NYT summary data
  const nytSummary: NYTSummary | null =
    entries.find((e) => e.integrations?.nyt?.data?.summary)
      ?.integrations.nyt.data.summary ?? null;

  const wordle = nytSummary?.wordle;
  const connections = nytSummary?.connections;

  const wordleTotal = wordle?.totalStats;
  const wordleStreak = wordle?.calculatedStats;

  const wordleWinRate = wordleTotal
    ? Math.round((wordleTotal.gamesWon / wordleTotal.gamesPlayed) * 100)
    : null;

  const guessDistribution = wordleTotal
    ? ([1, 2, 3, 4, 5, 6] as const).map((n) => ({
        n,
        count: wordleTotal.guesses[String(n) as keyof typeof wordleTotal.guesses] ?? 0,
      }))
    : [];
  const maxGuessCount = Math.max(...guessDistribution.map((g) => g.count), 1);

  const connWinRate = connections
    ? Math.round((connections.puzzles_won / connections.puzzles_completed) * 100)
    : null;

  const mistakeDist = connections
    ? ([0, 1, 2, 3, 4] as const).map((n) => ({
        n,
        count: connections.mistakes[String(n) as keyof typeof connections.mistakes] ?? 0,
      }))
    : [];
  const maxMistakeCount = Math.max(...mistakeDist.map((m) => m.count), 1);

  return (
    <div className="space-y-3 text-left">
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium">Insights</span>
      </div>
      <div className="grid grid-cols-3 gap-2.5">
        <div className="border bg-card rounded-lg p-3">
          <div className="text-2xl font-medium">{avgNow.toFixed(1)}</div>
          <div className="text-[11px] text-muted-foreground">avg quality, 30d</div>
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
          <div className="text-2xl font-medium">
            {topMoods[0]?.[0] ?? "—"}
          </div>
          <div className="text-[11px] text-muted-foreground">top mood, 30d</div>
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
            <div className="text-2xl font-medium">{last30.length}</div>
            <div className="text-[11px] text-muted-foreground">entries, 30d</div>
          </div>
          <div className="text-[11px] text-muted-foreground mt-1">
            {entries.length} total
          </div>
        </div>
      </div>

      {nytSummary && (
        <>
          <span className="block text-[11px] font-semibold uppercase text-muted-foreground">NYT Games - all time</span>
          <div className="grid grid-cols-2 gap-2.5">

            <div className="border bg-card rounded-lg p-3 space-y-2.5">
              <span className="text-xs font-medium text-muted-foreground">wordle</span>

              <div className="grid grid-cols-4 gap-3 mt-2 text-center">
                <div>
                  <div className="text-base font-medium leading-none">{wordleTotal?.gamesPlayed ?? "—"}</div>
                  <div className="text-[10px] text-muted-foreground mt-0.5">played</div>
                </div>
                <div>
                  <div className="text-base font-medium leading-none">{wordleWinRate != null ? `${wordleWinRate}%` : "—"}</div>
                  <div className="text-[10px] text-muted-foreground mt-0.5">win rate</div>
                </div>
                <div>
                  <div className="text-base font-medium leading-none">{wordleStreak?.currentStreak ?? "—"}</div>
                  <div className="text-[10px] text-muted-foreground mt-0.5">streak</div>
                </div>
                <div>
                  <div className="text-base font-medium leading-none">{wordleStreak?.maxStreak ?? "—"}</div>
                  <div className="text-[10px] text-muted-foreground mt-0.5">best</div>
                </div>
              </div>

              <div className="space-y-1">
                {guessDistribution.map(({ n, count }) => (
                  <div key={n} className="flex items-center gap-1.5">
                    <span className="text-[10px] text-muted-foreground w-2.5 shrink-0">{n}</span>
                    <div className="flex-1 h-3 bg-muted rounded-sm overflow-hidden">
                      <div
                        className="h-full bg-[#538d4e] rounded-sm"
                        style={{ width: `${(count / maxGuessCount) * 100}%` }}
                      />
                    </div>
                    <span className="text-[10px] text-muted-foreground w-4 text-right shrink-0">{count}</span>
                  </div>
                ))}
              </div>
            </div>

            <div className="border bg-card rounded-lg p-3 space-y-2.5">
              <span className="text-xs font-medium text-muted-foreground">connections</span>

              <div className="grid grid-cols-4 gap-3 mt-2 text-center">
                <div>
                  <div className="text-base font-medium leading-none">{connections?.puzzles_completed ?? "—"}</div>
                  <div className="text-[10px] text-muted-foreground mt-0.5">played</div>
                </div>
                <div>
                  <div className="text-base font-medium leading-none">{connWinRate != null ? `${connWinRate}%` : "—"}</div>
                  <div className="text-[10px] text-muted-foreground mt-0.5">win rate</div>
                </div>
                <div>
                  <div className={cn("text-base font-medium leading-none", connections?.current_streak === 0 && "text-muted-foreground")}>
                    {connections?.current_streak ?? "—"}
                  </div>
                  <div className="text-[10px] text-muted-foreground mt-0.5">streak</div>
                </div>
                <div>
                  <div className="text-base font-medium leading-none">{connections?.max_streak ?? "—"}</div>
                  <div className="text-[10px] text-muted-foreground mt-0.5">best</div>
                </div>
              </div>

              <div className="space-y-1">
                {mistakeDist.map(({ n, count }) => (
                  <div key={n} className="flex items-center gap-1.5">
                    <span className="text-[10px] text-muted-foreground w-2.5 shrink-0">{n}</span>
                    <div className="flex-1 h-3 bg-muted rounded-sm overflow-hidden">
                      <div
                        className="h-full rounded-sm"
                        style={{
                          width: `${(count / maxMistakeCount) * 100}%`,
                          // fewer mistakes = greener, more mistakes = more muted
                          backgroundColor: n === 0 ? "#a0c35a" : n === 4 ? "#e06060" : `color-mix(in srgb, #a0c35a ${100 - n * 20}%, #e06060)`,
                        }}
                      />
                    </div>
                    <span className="text-[10px] text-muted-foreground w-4 text-right shrink-0">{count}</span>
                  </div>
                ))}
              </div>
            </div>

          </div>
        </>
      )}
    </div>
  );
}
