import { Button } from "@/components/ui/button";
import { RefreshCw } from "lucide-react";
import { Tooltip, TooltipTrigger, TooltipContent } from "../ui/tooltip";
import type { IntegrationPayload } from "@/types";

interface GenericIntegrationCardProps {
  name: string;
  payload: IntegrationPayload;
  onSync: () => void;
  isSyncing: boolean;
}

export function GenericIntegrationCard({ name, payload, onSync, isSyncing }: GenericIntegrationCardProps) {
  const entries = Object.entries(payload.data ?? {});

  return (
    <div className="rounded-lg border border-border p-3 flex flex-col gap-2.5">
      <div className="flex items-center justify-between">
        <span className="text-xs font-medium capitalize">{name}</span>
        <Tooltip>
          <TooltipTrigger asChild>
            <Button
              variant="ghost"
              size="xs"
              className="text-muted-foreground"
              onClick={onSync}
              disabled={isSyncing}
            >
              <RefreshCw size={10} className={isSyncing ? "animate-spin" : ""} />
              sync
            </Button>
          </TooltipTrigger>
          <TooltipContent side="left">
            <p className="text-xs">
              {payload.fetched_at ? `synced ${new Date(payload.fetched_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}` : "never fetched"}
            </p>
          </TooltipContent>
        </Tooltip>
      </div>
      <div className="flex flex-col gap-1">
        {entries.map(([key, val]) => (
          <div key={key} className="flex items-start justify-between gap-2 text-[11px]">
            <span className="text-muted-foreground capitalize">{key.replace(/_/g, " ")}</span>
            <span className="text-foreground text-right break-all max-w-[60%]">
              {typeof val === "object" ? JSON.stringify(val) : String(val)}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}
