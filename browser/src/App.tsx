import { useState, useMemo } from "react";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { getQualityClass, getLocalDateString } from "./lib/utils";
import { useDailies } from "@/hooks/useDailies";
import type { DailyEntry, IntegrationPayload } from "@/types";

import { SidebarRow } from "@/components/SidebarRow";
import { EntryPanel } from "@/components/EntryPanel";
import { InsightsStrip } from "@/components/InsightsStrip";

import {
  ArrowLeft,
  Plus,
  ChevronLeft,
  ChevronRight,
  RotateCw,
  Menu,
  X,
} from "lucide-react";

function getStartOfWeek(dateStr: string): Date {
  const [year, month, day] = dateStr.split("-").map(Number);
  const date = new Date(year, month - 1, day);
  const dayOfWeek = date.getDay();
  const diff = date.getDate() - (dayOfWeek === 0 ? 6 : dayOfWeek - 1);

  return new Date(date.setDate(diff));
}

export default function DailiesPage() {
  const TODAY = getLocalDateString();

  const {
    entries,
    config,
    isLoading,
    error,
    saveEntry,
    runIntegration,
    refreshData,
  } = useDailies();

  const [selectedDate, setSelectedDate] = useState<string>(TODAY);
  const [syncingIntegration, setSyncingIntegration] = useState<string | null>(
    null
  );
  const [draft, setDraft] = useState<DailyEntry | null>(null);
  const [mode, setMode] = useState<"view" | "edit" | "create">("view");
  const [weekOffset, setWeekOffset] = useState<number>(0);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  const sortedEntries = useMemo(
    () => [...entries].sort((a, b) => b.date.localeCompare(a.date)),
    [entries]
  );

  const groupedEntries = useMemo(() => {
    const groups: { label: string; items: DailyEntry[] }[] = [];
    const map: Record<string, DailyEntry[]> = {};

    sortedEntries.forEach((entry) => {
      const monday = getStartOfWeek(entry.date);
      const sunday = new Date(monday);
      sunday.setDate(monday.getDate() + 6);

      const label = `${monday.toLocaleDateString([], {
        month: "short",
        day: "numeric",
      })} - ${sunday.toLocaleDateString([], {
        month: "short",
        day: "numeric",
      })}`;

      if (!map[label]) {
        map[label] = [];
        groups.push({ label, items: map[label] });
      }

      map[label].push(entry);
    });

    return groups;
  }, [sortedEntries]);

  const selectedEntry = useMemo(
    () => sortedEntries.find((e) => e.date === selectedDate) ?? null,
    [sortedEntries, selectedDate]
  );

  const isToday = selectedDate === TODAY;
  const entry = (draft ?? selectedEntry) as DailyEntry | null;

  const weekData = useMemo(() => {
    const list = [];

    for (let i = 6; i >= 0; i--) {
      const d = new Date();
      d.setDate(d.getDate() - i - weekOffset * 7);

      const dateStr = getLocalDateString(d);
      const existing = sortedEntries.find((e) => e.date === dateStr);

      list.push(existing ?? { date: dateStr, day_quality: 0 });
    }

    return list;
  }, [sortedEntries, weekOffset]);

  const dateRangeLabel = useMemo(() => {
    if (weekData.length === 0) return "";

    const start = new Date(weekData[0].date).toLocaleDateString([], {
      month: "short",
      day: "numeric",
    });

    const end = new Date(
      weekData[weekData.length - 1].date
    ).toLocaleDateString([], {
      month: "short",
      day: "numeric",
    });

    return `${start} - ${end}`;
  }, [weekData]);

  function selectDate(date: string) {
    setSelectedDate(date);
    setDraft(null);
    setMode("view");
    setMobileMenuOpen(false);
  }

  function startCreate() {
    let emptyIntegrations: Record<string, IntegrationPayload> = {};

    if (config?.integrations) {
      emptyIntegrations = Object.entries(config.integrations).reduce<
        Record<string, IntegrationPayload>
      >((acc, [key, integration]) => {
        if (integration?.enabled) {
          acc[key] = {
            data: {},
            fetched_at: null,
          };
        }

        return acc;
      }, {});
    }

    const empty = {
      date: TODAY,
      day_quality: 0,
      moods: [],
      context_tags: [],
      journal: "",
      integrations: emptyIntegrations,
    };

    setDraft(empty);
    setMode("create");
    setMobileMenuOpen(false);
  }

  function startEdit() {
    if (!selectedEntry) return;

    setDraft(selectedEntry);
    setMode("edit");
  }

  function cancelEdit() {
    setDraft(null);
    setMode("view");
  }

  async function handleSave() {
    if (!draft) return;

    await saveEntry(draft);
    setDraft(null);
    setMode("view");
    refreshData();
  }

  async function handleSync(date: string, integration: string) {
    setSyncingIntegration(integration);

    try {
      await runIntegration(date, integration);
    } finally {
      setSyncingIntegration(null);
    }
  }

  if (isLoading) {
    return (
      <div className="h-screen flex items-center justify-center">
        loading…
      </div>
    );
  }

  if (error) {
    return (
      <div className="h-screen flex items-center justify-center text-red-500">
        {error}
      </div>
    );
  }

  return (
    <div className="flex flex-col h-screen bg-background">
      {/* HEADER */}
      <header className="flex items-center justify-between px-4 sm:px-6 h-12 border-b shrink-0">
        <div className="flex items-center gap-2">
          {/* Mobile menu button */}
          <Button
            variant="ghost"
            size="icon"
            className="md:hidden h-8 w-8"
            onClick={() => setMobileMenuOpen(true)}
            aria-label="Open menu"
          >
            <Menu className="h-4 w-4" />
          </Button>

          <span className="text-sm font-medium">dailies</span>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8"
            onClick={refreshData}
            aria-label="Refresh"
          >
            <RotateCw className="h-4 w-4" />
          </Button>
        </div>
      </header>

      <div className="flex flex-1 overflow-hidden">
        {/* DESKTOP SIDEBAR */}
        <aside className="hidden md:flex w-64 border-r flex-col shrink-0">
          <ScrollArea className="flex-1 px-2">
            <div className="p-3 border-b flex flex-col gap-2">
              <div className="flex items-center justify-between">
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-6 w-6"
                  onClick={() =>
                    setWeekOffset((prev) => prev + 1)
                  }
                >
                  <ChevronLeft className="h-3.5 w-3.5" />
                </Button>

                <span className="text-[11px] font-medium uppercase text-muted-foreground">
                  week
                </span>

                <Button
                  variant="ghost"
                  size="icon"
                  className="h-6 w-6"
                  disabled={weekOffset === 0}
                  onClick={() =>
                    setWeekOffset((prev) =>
                      Math.max(0, prev - 1)
                    )
                  }
                >
                  <ChevronRight className="h-3.5 w-3.5" />
                </Button>
              </div>

              <div className="flex gap-1.5 h-12 items-end px-1 pt-1">
                {weekData.map((e) => {
                  const rating = e.day_quality ?? 0;
                  const percentHeight =
                    rating > 0
                      ? `${(rating / 10) * 100}%`
                      : "8%";

                  const isActive =
                    e.date === selectedDate;

                  const isUnrated = rating === 0;

                  const baseBg = isUnrated
                    ? "bg-muted/40"
                    : getQualityClass(rating, "base");

                  const hoverBg = isUnrated
                    ? "hover:bg-muted/70"
                    : getQualityClass(rating, "hover");

                  return (
                    <button
                      key={e.date}
                      title={`${e.date}: ${rating}/10`}
                      onClick={() => selectDate(e.date)}
                      className={`flex-1 rounded-sm transition-colors relative outline-none ${baseBg} ${hoverBg} ${
                        isActive
                          ? "ring-2 ring-primary ring-offset-1 ring-offset-background scale-105 z-10"
                          : ""
                      }`}
                      style={{ height: percentHeight }}
                    />
                  );
                })}
              </div>

              <div className="text-[10px] text-center text-muted-foreground font-medium mt-0.5">
                {dateRangeLabel}
              </div>
            </div>

            <div className="space-y-2">
              {groupedEntries.map((group) => (
                <div key={group.label} className="space-y-1">
                  <span className="text-[11px] font-semibold uppercase text-muted-foreground block text-left px-2 pt-2">
                    {group.label}
                  </span>

                  {group.items.map((e) => (
                    <SidebarRow
                      key={e.date}
                      entry={e}
                      active={e.date === selectedDate}
                      onClick={() => selectDate(e.date)}
                    />
                  ))}
                </div>
              ))}
            </div>
          </ScrollArea>
        </aside>

        {/* MOBILE SIDEBAR DRAWER */}
        {mobileMenuOpen && (
          <div
            className="fixed inset-0 z-50 md:hidden"
            onClick={() => setMobileMenuOpen(false)}
          >
            {/* Backdrop */}
            <div className="absolute inset-0 bg-black/40" />

            {/* Drawer */}
            <aside
              className="absolute inset-y-0 left-0 w-[min(20rem,85vw)] bg-background border-r shadow-xl flex flex-col"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="h-12 shrink-0 border-b flex items-center justify-between px-3">
                <span className="text-sm font-medium">
                  dailies
                </span>

                <Button
                  variant="ghost"
                  size="icon"
                  className="h-8 w-8"
                  onClick={() => setMobileMenuOpen(false)}
                  aria-label="Close menu"
                >
                  <X className="h-4 w-4" />
                </Button>
              </div>

              <ScrollArea className="flex-1 px-2">
                <div className="p-3 border-b flex flex-col gap-2">
                  <div className="flex items-center justify-between">
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-6 w-6"
                      onClick={() =>
                        setWeekOffset((prev) => prev + 1)
                      }
                    >
                      <ChevronLeft className="h-3.5 w-3.5" />
                    </Button>

                    <span className="text-[11px] font-medium uppercase text-muted-foreground">
                      week
                    </span>

                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-6 w-6"
                      disabled={weekOffset === 0}
                      onClick={() =>
                        setWeekOffset((prev) =>
                          Math.max(0, prev - 1)
                        )
                      }
                    >
                      <ChevronRight className="h-3.5 w-3.5" />
                    </Button>
                  </div>

                  <div className="flex gap-1.5 h-12 items-end px-1 pt-1">
                    {weekData.map((e) => {
                      const rating = e.day_quality ?? 0;

                      const percentHeight =
                        rating > 0
                          ? `${(rating / 10) * 100}%`
                          : "8%";

                      const isActive =
                        e.date === selectedDate;

                      const isUnrated = rating === 0;

                      const baseBg = isUnrated
                        ? "bg-muted/40"
                        : getQualityClass(
                            rating,
                            "base"
                          );

                      const hoverBg = isUnrated
                        ? "hover:bg-muted/70"
                        : getQualityClass(
                            rating,
                            "hover"
                          );

                      return (
                        <button
                          key={e.date}
                          title={`${e.date}: ${rating}/10`}
                          onClick={() => selectDate(e.date)}
                          className={`flex-1 rounded-sm transition-colors relative outline-none ${baseBg} ${hoverBg} ${
                            isActive
                              ? "ring-2 ring-primary ring-offset-1 ring-offset-background scale-105 z-10"
                              : ""
                          }`}
                          style={{
                            height: percentHeight,
                          }}
                        />
                      );
                    })}
                  </div>

                  <div className="text-[10px] text-center text-muted-foreground font-medium mt-0.5">
                    {dateRangeLabel}
                  </div>
                </div>

                <div className="space-y-2">
                  {groupedEntries.map((group) => (
                    <div
                      key={group.label}
                      className="space-y-1"
                    >
                      <span className="text-[11px] font-semibold uppercase text-muted-foreground block text-left px-2 pt-2">
                        {group.label}
                      </span>

                      {group.items.map((e) => (
                        <SidebarRow
                          key={e.date}
                          entry={e}
                          active={e.date === selectedDate}
                          onClick={() => selectDate(e.date)}
                        />
                      ))}
                    </div>
                  ))}
                </div>
              </ScrollArea>
            </aside>
          </div>
        )}

        {/* MAIN CONTENT */}
        <main className="flex-1 overflow-hidden min-w-0">
          <ScrollArea className="h-full">
            <div className="max-w-3xl mx-auto px-4 sm:px-6 py-5 flex flex-col gap-6">
              {!isToday && (
                <div className="flex">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() =>
                      setSelectedDate(TODAY)
                    }
                  >
                    <ArrowLeft className="w-3 h-3" />
                    back to today
                  </Button>
                </div>
              )}

              {!entry && isToday && (
                <Button
                  onClick={startCreate}
                  className="w-full"
                >
                  <Plus className="w-4 h-4" />
                  create
                </Button>
              )}

              {entry && (
                <EntryPanel
                  entry={entry}
                  mode={mode}
                  isToday={isToday}
                  onEdit={startEdit}
                  onCancel={cancelEdit}
                  onSave={handleSave}
                  onChange={setDraft}
                  onSync={handleSync}
                  syncingIntegration={syncingIntegration}
                />
              )}

              <InsightsStrip entries={sortedEntries} />
            </div>
          </ScrollArea>
        </main>
      </div>
    </div>
  );
}