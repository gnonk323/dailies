import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import type { IntegrationPayload } from "@/types";
import { useState } from "react";
import {
  RefreshCw,
  GitPullRequestArrow,
  GitPullRequestCreate,
  FolderRoot,
  Plus,
  Minus,
  GitCommitHorizontal,
  ExternalLink,
} from "lucide-react";

interface GitHubCardProps {
  payload: IntegrationPayload;
  onSync: () => void;
  isSyncing: boolean;
}

export function GitHubCard({ payload, onSync, isSyncing }: GitHubCardProps) {
  const data = (payload.data ?? {}) as {
    commits?: Record<string, string>;
    commits_count?: number;
    prs_merged?: Record<string, string>;
    prs_opened?: Record<string, string>;
    repos?: Record<string, string>;
  };

  const commits = Object.entries(data.commits ?? {});
  const prsMerged = Object.keys(data.prs_merged ?? {});
  const prsOpened = Object.keys(data.prs_opened ?? {});
  const repos = Object.entries(data.repos ?? {});
  const mergedPrTitles = new Set(prsMerged);
  const VISIBLE = 5;
  const extra = commits.length - VISIBLE;

  const [expanded, setExpanded] = useState(false);
  const visibleCommits = expanded ? commits : commits.slice(0, VISIBLE);

  const toggleExpanded = () => {
    setExpanded(!expanded);
  }

  return (
    <div className="rounded-lg border border-border p-3 flex flex-col gap-2.5">
      <div className="flex items-center justify-between">
        <span className="text-xs font-medium flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" className="bi bi-github" viewBox="0 0 16 16">
            <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27s1.36.09 2 .27c1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.01 8.01 0 0 0 16 8c0-4.42-3.58-8-8-8"/>
          </svg>
          github
        </span>
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

      <div className="flex gap-4">
        {[
          { label: "commits", val: data.commits_count ?? commits.length, icon: <GitCommitHorizontal /> },
          { label: "prs opened", val: prsOpened.length, icon: <GitPullRequestCreate /> },
          { label: "prs merged", val: prsMerged.length, icon: <GitPullRequestArrow /> },
          { label: "repos", val: repos.length, icon: <FolderRoot /> },
        ].map(({ label, val, icon }) => (
          <div key={label} className="text-center text-muted-foreground">
            <div className="flex items-center justify-center gap-2">
              {icon}
              <div className="font-medium leading-none text-foreground text-xl">{val}</div>
            </div>
            <div className="text-[10px] mt-0.5 text-center px-3">{label}</div>
          </div>
        ))}
      </div>

      {visibleCommits.length > 0 && (
        <div className="flex flex-col gap-1">
          {visibleCommits.map(([msg, sha]) => (
            <div key={sha} className="flex items-center gap-1.5 text-[11px] text-muted-foreground overflow-hidden">
              <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 shrink-0" />
              <span className="truncate flex-1">{msg}</span>
              {mergedPrTitles.has(msg) && (
                <Badge className="bg-[#8250df] text-white font-semibold text-[11px]">
                  PR
                </Badge>
              )}
            </div>
          ))}
          {extra > 0 && (
            <div>
              <Button
                variant="ghost"
                size="xs"
                className="text-[11px] text-muted-foreground"
                onClick={toggleExpanded}
              >
                {!expanded ? (
                  <>
                  <Plus />
                  <span>see more</span>
                  </>
                ) : (
                  <>
                  <Minus />
                  <span>see less</span>
                  </>
                )}
              </Button>
            </div>
          )}
        </div>
      )}

      {repos.length > 0 && (
        <div className="flex flex-wrap gap-1">
          {repos.map(([name, url]) => (
            <a
              key={name}
              href={url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-[10px] px-1.5 py-0.5 flex items-center gap-1 rounded border border-border bg-muted text-muted-foreground hover:text-foreground transition-colors"
            >
              {name}
              <ExternalLink size={10}/>
            </a>
          ))}
        </div>
      )}
    </div>
  );
}
