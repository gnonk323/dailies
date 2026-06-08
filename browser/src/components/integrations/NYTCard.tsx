import { Button } from "@/components/ui/button";
import { CircleCheck, RefreshCw, Check, X, ArrowUpRightFromSquare } from "lucide-react";
import { Tooltip, TooltipTrigger, TooltipContent } from "../ui/tooltip";
import type { IntegrationPayload, NYTData } from "@/types";
import ProgressCircle from "../ui/ProgressCircle";
import { Badge } from "../ui/badge";

interface NYTCardProps {
  payload: IntegrationPayload;
  onSync: () => void;
  isSyncing: boolean;
}

export function NYTCard({ payload, onSync, isSyncing }: NYTCardProps) {
  const data = (payload.data ?? {}) as NYTData;

  const mini = data.crossword_mini;
  const midi = data.crossword_midi;
  const wordle = data.wordle;
  const connections = data.connections;

  const isEmpty = !mini && !midi && !wordle && !connections;

  const formatSeconds = (seconds: number): string => {
    const m = Math.floor(seconds / 60);
    const s = String(seconds % 60).padStart(2, '0')
    return `${m}:${s}`
  }

  const getConnectionsBgColor = (level: number) => {
    switch (level) {
      case 0: return "bg-[#f9df6d]";
      case 1: return "bg-[#a0c35a]";
      case 2: return "bg-[#b0c4ef]";
      case 3: return "bg-[#ba81c5]";
      default: return "bg-muted";
    }
  };

  const goToPuzzle = (puzzle: string) => {
    if (puzzle === "mini" || puzzle === "midi") {
      window.open(`https://www.nytimes.com/crosswords/game/${puzzle}/${mini.publicationDate.split("-").join("/")}`)
    } else {
      window.open(`https://www.nytimes.com/games/${puzzle}/${mini.publicationDate}`)
    }
  }

  return (
    <div className="rounded-lg border border-border p-3 flex flex-col gap-2.5">
      <div className="flex items-center justify-between">
        <span className="text-xs font-medium flex items-center gap-2">
          <img src="nyt-icon.svg" alt="NYT Icon" />
          nyt games
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
      {!isEmpty && (
        <div className="flex flex-col gap-1">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <div className="flex items-center gap-1 text-muted-foreground">
                <span className="text-[11px] font-medium uppercase">
                  mini crossword
                </span>
                <ArrowUpRightFromSquare
                  className="hover:text-foreground cursor-pointer"
                  size={12}
                  onClick={() => goToPuzzle("mini")}
                />
              </div>
              <div className="flex items-center gap-1">
                {mini.completionFraction === 1 ? (
                  <CircleCheck size={20} className="text-blue-500" />
                ) : (
                  <ProgressCircle size={20} value={mini.completionFraction} />
                )}
                {mini.completionFraction === 0 ? (
                  <div className="text-base font-medium text-muted-foreground leading-none my-1">0:00</div>
                ) : (
                  <div className="text-base font-medium leading-none my-1">{mini.completionFraction === 1 ? formatSeconds(mini.firstSolve) : formatSeconds(mini.playTime)}</div>
                )}
              </div>
              <p className="text-[10px] text-muted-foreground">{mini.constructors.join(", ")}</p>
            </div>
            <div>
              <div className="flex items-center gap-1 text-muted-foreground">
                <span className="text-[11px] font-medium uppercase">midi crossword{midi.title && ` - "${midi.title}"`}</span>
                <ArrowUpRightFromSquare
                  className="hover:text-foreground cursor-pointer"
                  size={12}
                  onClick={() => goToPuzzle("midi")}
                />
              </div>
              <div className="flex items-center gap-1">
                {midi.completionFraction === 1 ? (
                  <CircleCheck size={20} className="text-blue-500" />
                ) : (
                  <ProgressCircle size={20} value={midi.completionFraction} />
                )}
                {midi.completionFraction === 0 ? (
                  <div className="text-base font-medium text-muted-foreground leading-none my-1">0:00</div>
                ) : (
                  <div className="text-base font-medium leading-none my-1">{midi.completionFraction === 1 ? formatSeconds(midi.firstSolve) : formatSeconds(midi.playTime)}</div>
                )}
              </div>
              <p className="text-[10px] text-muted-foreground">{midi.constructors.join(", ")}</p>
            </div>
            <div>
              <div className="flex items-center gap-2">
                <div className="flex items-center gap-1 text-muted-foreground">
                  <span className="text-[11px] font-medium uppercase">wordle</span>
                  <ArrowUpRightFromSquare
                    className="hover:text-foreground cursor-pointer"
                    size={12}
                    onClick={() => goToPuzzle("wordle")}
                  />
                </div>
                {wordle.game_data.hardMode && <Badge className="text-[10px]" variant='outline'>hard mode</Badge>}
                <Badge className="text-[10px]" variant={wordle.game_data.status === "FAIL" ? 'destructive' : 'default'}>{wordle.game_data.status === "WIN" ? "won" : wordle.game_data.status === "FAIL" ? "failed" : "in progress"} {wordle.game_data.currentRowIndex} / 6</Badge>
              </div>
              <div className="space-y-1 mt-2">
                {wordle.game_data.boardState.map((guess, guessIdx) => {
                  const solution = wordle.solution.toLowerCase();
                  const isCorrectAnswer = guess === solution;
                  const rowHeightClass = (isCorrectAnswer) ? "h-12" : "h-6 text-xs";

                  const lettersArray = guess.padEnd(5, " ").split("");

                  const solutionLetterCounts: Record<string, number> = {};
                  for (const char of solution) {
                    solutionLetterCounts[char] = (solutionLetterCounts[char] || 0) + 1;
                  }

                  const tileStatuses = Array(5).fill("empty");

                  lettersArray.forEach((letter, idx) => {
                    if (letter !== " ") {
                      if (solution[idx] === letter) {
                        tileStatuses[idx] = "correct";
                        solutionLetterCounts[letter]--;
                      }
                    }
                  });

                  lettersArray.forEach((letter, idx) => {
                    if (letter !== " " && tileStatuses[idx] !== "correct") {
                      if (solutionLetterCounts[letter] > 0) {
                        tileStatuses[idx] = "present";
                        solutionLetterCounts[letter]--;
                      } else {
                        tileStatuses[idx] = "absent";
                      }
                    }
                  });

                  return (
                    <div className="grid grid-cols-5 gap-1" key={guessIdx}>
                      {lettersArray.map((letter, letterIdx) => {
                        const status = tileStatuses[letterIdx];
                        
                        let bgColorClass = "bg-muted";

                        if (status === "correct") {
                          bgColorClass = "bg-[#538d4e]";
                        } else if (status === "present") {
                          bgColorClass = "bg-[#b59f3b]";
                        }

                        return (
                          <div
                            className={`flex items-center justify-center rounded uppercase p-1 font-bold ${rowHeightClass} ${bgColorClass}`}
                            key={letterIdx}
                          >
                            {letter !== " " ? letter : ""}
                          </div>
                        );
                      })}
                    </div>
                  );
                })}
              </div>
            </div>
            
            <div>
              <div className="flex items-center gap-2">
                <div className="flex items-center gap-1 text-muted-foreground">
                  <span className="text-[11px] font-medium uppercase">connections</span>
                  <ArrowUpRightFromSquare
                    className="hover:text-foreground cursor-pointer"
                    size={12}
                    onClick={() => goToPuzzle("connections")}
                  />
                </div>
                {connections.game_data.puzzleWon ? (
                  <Badge variant="default" className="text-[10px]">won</Badge>
                ) : (
                  <Badge variant={connections.game_data.puzzleComplete ? 'destructive' : 'outline'} className="text-[10px]">
                    {connections.game_data.puzzleComplete ? "failed" : `${4 - connections.game_data.mistakes} mistakes left`}
                  </Badge>
                )}
              </div>
              <div className="space-y-1 mt-2">
                {connections.game_data.guesses.map((guess, guessIdx) => (
                  <div className="flex items-center gap-1" key={guessIdx}>
                    <div className="grid grid-cols-4 gap-1 flex-1 mr-2">
                      {guess.cards.map((card, cardIdx) => (
                        <div
                          key={cardIdx}
                          className={`h-6 rounded flex items-center justify-center transition-colors ${getConnectionsBgColor(card.level)}`}
                        />
                      ))}
                    </div>
                    {guess.correct
                      ? <Check size={12} className="text-green-500 shrink-0" />
                      : <X size={12} className="text-muted-foreground/40 shrink-0" />
                    }
                  </div>
                ))}
                <div className="mt-2 flex items-center gap-2">
                  {([0, 1, 2, 3] as const).map((level) => {
                    const solved = connections.game_data.solvedCategories.some(
                      (cat) => cat.level === level
                    );
                    const colorClass = getConnectionsBgColor(level);
                    const labels = ["yellow", "green", "blue", "purple"];
                    return (
                      <div key={level} className="flex items-center gap-1.5">
                        <div className={`w-2.5 h-2.5 rounded-sm shrink-0 ${colorClass}`} />
                        <span
                          className={`text-[10px] ${
                            solved
                              ? "text-muted-foreground"
                              : "text-muted-foreground/40 line-through"
                          }`}
                        >
                          {labels[level]}
                        </span>
                        {solved && <Check size={9} className="text-green-500 shrink-0" />}
                      </div>
                    );
                  })}
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}