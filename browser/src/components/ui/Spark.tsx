export function Spark({ values }: { values: number[] }) {
  const max = Math.max(...values, 1);
  return (
    <div className="flex items-end gap-0.5 h-6">
      {values.map((v, i) => (
        <div
          key={i}
          className="w-2 rounded-sm bg-muted-foreground/30"
          style={{ height: `${Math.round((v / max) * 100)}%`, minHeight: 2 }}
        />
      ))}
    </div>
  );
}
