import { Badge } from "@/components/ui/badge";
import { Plus, X } from "lucide-react";

interface ChipFieldProps {
  values: string[];
  editable: boolean;
  onAdd: (value: string) => void
  onRemove: (value: string) => void
  placeholder?: string;
}

export function ChipField({ values, editable, onRemove, onAdd, placeholder }: ChipFieldProps) {
  if (!editable) {
    return (
      <div className="flex flex-wrap gap-1 min-h-[6]">
        {!values || values.length === 0 ? (
          <span className="text-xs text-muted-foreground">none</span>
        ) : (
          values.map((v) => (
            <Badge key={v} variant="secondary" className="text-[11px] font-normal">
              {v}
            </Badge>
          ))
        )}
      </div>
    );
  }

  return (
    <div className="flex flex-wrap gap-1 items-center min-h-[9] px-2 py-1.5 rounded-md border border-input bg-background">
      {values.map((v) => (
        <Badge key={v} variant="secondary" className="text-[11px] font-normal pr-1 gap-1">
          {v}
          <button
            type="button"
            onClick={() => onRemove?.(v)}
            className="hover:text-foreground text-muted-foreground transition-colors"
          >
            <X className="w-2.5 h-2.5" />
          </button>
        </Badge>
      ))}
      <button
        type="button"
        onClick={() => {
          const value = window.prompt("Add item");
          if (value) onAdd(value);
        }}
        className="text-[11px] text-muted-foreground hover:text-foreground flex items-center gap-0.5 px-1 transition-colors"
      >
        <Plus className="w-3 h-3" />
        <span>{placeholder ?? "add"}</span>
      </button>
    </div>
  );
}
