import { useState } from "react";
import { cn, getQualityClass } from "@/lib/utils";

interface QualityTrackProps {
  value: number;
  editable: boolean;
  onChange?: (v: number) => void;
}

export function QualityTrack({ value, editable, onChange }: QualityTrackProps) {
  const [hoveredIndex, setHoveredIndex] = useState<number | null>(null);

  return (
    <div className="flex items-center gap-2 w-full">
      <div 
        className="grid grid-cols-10 gap-1 w-full"
        onMouseLeave={() => editable && setHoveredIndex(null)}
      >
        {Array.from({ length: 10 }, (_, i) => {
          const n = i + 1;
          const filled = n <= value;

          const isForwardHover =
            editable &&
            hoveredIndex !== null &&
            n > value && 
            n <= hoveredIndex;

          const isBackwardHover =
            editable &&
            hoveredIndex !== null &&
            n <= value &&
            n >= hoveredIndex;

          return (
            <button
              key={n}
              type="button"
              disabled={!editable}
              onClick={() => editable && onChange?.(n)}
              onMouseEnter={() => editable && setHoveredIndex(n)}
              className={cn(
                "h-8 rounded text-[11px] border transition-colors text-center font-bold",
                
                filled ? getQualityClass(value, "base") : "border-border text-muted-foreground",
                
                isForwardHover && "bg-muted",
                isBackwardHover && getQualityClass(value, "active"),
                
                editable && !filled && !isForwardHover && "hover:bg-muted",
              )}
            >
              {n}
            </button>
          );
        })}
      </div>
    </div>
  );
}