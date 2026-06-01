import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Textarea } from "@/components/ui/textarea";
import { Separator } from "@/components/ui/separator";
import { Pencil, Check } from "lucide-react";
import { formatDisplayDate } from "@/lib/utils";
import type { DailyEntry } from "@/types";

import { QualityTrack } from "./ui/QualityTrack";
import { ChipField } from "./ui/ChipField";
import { IntegrationCard } from "./integrations/IntegrationCard";

type Mode = "view" | "edit" | "create";

interface Props {
  entry: DailyEntry;
  mode: Mode;
  isToday: boolean;

  onEdit: () => void;
  onCancel: () => void;
  onSave: () => void;
  onChange: (e: DailyEntry) => void;

  onSync: (date: string, integration: string) => Promise<void>;
  syncingIntegration: string | null;
}

export function EntryPanel({
  entry,
  mode,
  isToday,
  onEdit,
  onCancel,
  onSave,
  onChange,
  onSync,
  syncingIntegration,
}: Props) {

  const isEditing = mode !== "view";

  const update = (patch: Partial<DailyEntry>) => {
    onChange({ ...entry, ...patch });
  };

  const integrationEntries = Object.entries(entry.integrations ?? {});

  return (
    <div className="rounded-xl border bg-card text-left">

      <div className="flex justify-between px-5 py-3 border-b">
        <div className="flex gap-2 items-center">
          <span className="text-sm font-medium">
            {formatDisplayDate(entry.date)}
          </span>

          {isToday && (
            <Badge variant="outline">today</Badge>
          )}
        </div>

        <div className="flex gap-2">
          {!isEditing && (
            <Button size="xs" variant="outline" onClick={onEdit}>
              <Pencil />
              edit
            </Button>
          )}

          {isEditing && (
            <>
              <Button size="xs" variant="ghost" onClick={onCancel}>
                cancel
              </Button>
              <Button size="xs" onClick={onSave}>
                <Check />
                save
              </Button>
            </>
          )}
        </div>
      </div>

      <div className="p-5 flex flex-col gap-5">
        
        <div className="flex flex-col gap-2">
          <span className="text-[11px] font-medium uppercase text-muted-foreground">day quality</span>
          <QualityTrack
            value={entry.day_quality}
            editable={isEditing}
            onChange={(v) => update({ day_quality: v })}
          />
        </div>

        <div className="grid grid-cols-2 gap-4">

          <div className="flex flex-col gap-2">
            <span className="text-[11px] font-medium uppercase text-muted-foreground">moods</span>
            <ChipField
              values={entry.moods}
              editable={isEditing}
              onAdd={(v) => update({ moods: [...entry.moods, v] })}
              onRemove={(v) =>
                update({ moods: entry.moods.filter(m => m !== v) })
              }
            />
          </div>

          <div className="flex flex-col gap-2">
            <span className="text-[11px] font-medium uppercase text-muted-foreground">context tags</span>
            <ChipField
              values={entry.context_tags}
              editable={isEditing}
              onAdd={(v) => update({ context_tags: [...entry.context_tags, v] })}
              onRemove={(v) =>
                update({ context_tags: entry.context_tags.filter(t => t !== v) })
              }
            />
          </div>
        </div>

        <Separator />

        <div className="flex flex-col gap-2">
          <span className="text-[11px] font-medium uppercase text-muted-foreground">journal</span>
          {isEditing ? (
            <Textarea
              value={entry.journal}
              placeholder="notes, reflection, whatever..."
              onChange={(e) => update({ journal: e.target.value })}
              className="resize-y"
            />
          ) : (
            <p className="text-sm text-muted-foreground">
              {entry.journal || "no journal entry"}
            </p>
          )}
        </div>

        {integrationEntries.length > 0 && (
          <>
            <Separator />
            {integrationEntries.map(([name, payload]) => (
              <IntegrationCard
                key={name}
                name={name}
                payload={payload}
                onSync={() => onSync(entry.date, name)}
                isSyncing={syncingIntegration === name}
              />
            ))}
          </>
        )}
      </div>
    </div>
  );
}