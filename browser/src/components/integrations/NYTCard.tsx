import { Button } from "@/components/ui/button";
import {
  CircleCheck,
  RefreshCw,
  Check,
  X,
  ArrowUpRightFromSquare,
} from "lucide-react";
import {
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from "../ui/tooltip";
import type { IntegrationPayload, NYTData } from "@/types";
import ProgressCircle from "../ui/ProgressCircle";
import { Badge } from "../ui/badge";

interface NYTCardProps {
  payload: IntegrationPayload;
  onSync: () => void;
  isSyncing: boolean;
}

export function NYTCard({
  payload,
  onSync,
  isSyncing,
}: NYTCardProps) {
  const data = (payload.data ?? {}) as NYTData;

  const mini = data.crossword_mini;
  const midi = data.crossword_midi;
  const wordle = data.wordle;
  const connections = data.connections;

  const hasMini =
    !!mini &&
    typeof mini.completionFraction === "number" &&
    typeof mini.firstSolve === "number" &&
    typeof mini.playTime === "number";

  const hasMidi =
    !!midi &&
    typeof midi.completionFraction === "number" &&
    typeof midi.firstSolve === "number" &&
    typeof midi.playTime === "number";

  const hasWordle =
    !!wordle &&
    !!wordle.game_data &&
    Array.isArray(wordle.game_data.boardState) &&
    typeof wordle.game_data.hardMode === "boolean" &&
    typeof wordle.game_data.status === "string" &&
    typeof wordle.game_data.currentRowIndex === "number" &&
    typeof wordle.solution === "string";

  const hasConnections =
    !!connections &&
    !!connections.game_data &&
    Array.isArray(connections.game_data.guesses) &&
    Array.isArray(connections.game_data.solvedCategories) &&
    typeof connections.game_data.mistakes === "number" &&
    typeof connections.game_data.puzzleComplete === "boolean";

  const isEmpty =
    !hasMini &&
    !hasMidi &&
    !hasWordle &&
    !hasConnections;

  const formatSeconds = (seconds: number): string => {
    if (!Number.isFinite(seconds) || seconds < 0) {
      return "0:00";
    }

    const m = Math.floor(seconds / 60);
    const s = String(Math.floor(seconds % 60)).padStart(2, "0");

    return `${m}:${s}`;
  };

  const getConnectionsBgColor = (level: number): string => {
    switch (level) {
      case 0:
        return "bg-[#f9df6d]";
      case 1:
        return "bg-[#a0c35a]";
      case 2:
        return "bg-[#b0c4ef]";
      case 3:
        return "bg-[#ba81c5]";
      default:
        return "bg-muted";
    }
  };

  const goToPuzzle = (
    puzzle: "mini" | "midi" | "wordle" | "connections",
  ) => {
    const publicationDate =
      puzzle === "mini"
        ? mini?.publicationDate
        : puzzle === "midi"
          ? midi?.publicationDate
          : undefined;

    if (puzzle === "mini" || puzzle === "midi") {
      if (!publicationDate) {
        return;
      }

      const datePath = publicationDate.split("-").join("/");

      window.open(
        `https://www.nytimes.com/crosswords/game/${puzzle}/${datePath}`,
        "_blank",
        "noopener,noreferrer",
      );

      return;
    }

    const date =
      publicationDate ??
      payload.fetched_at?.split("T")[0];

    if (!date) {
      return;
    }

    window.open(
      `https://www.nytimes.com/games/${puzzle}/${date}`,
      "_blank",
      "noopener,noreferrer",
    );
  };

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
              <RefreshCw
                size={10}
                className={isSyncing ? "animate-spin" : ""}
              />
              sync
            </Button>
          </TooltipTrigger>

          <TooltipContent side="left">
            <p className="text-xs">
              {payload.fetched_at
                ? `synced ${new Date(
                    payload.fetched_at,
                  ).toLocaleTimeString([], {
                    hour: "2-digit",
                    minute: "2-digit",
                  })}`
                : "never fetched"}
            </p>
          </TooltipContent>
        </Tooltip>
      </div>

      {!isEmpty && (
        <div className="flex flex-col gap-1">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">

            {/* MINI CROSSWORD */}
            {hasMini && mini && (
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
                    <CircleCheck
                      size={20}
                      className="text-blue-500"
                    />
                  ) : (
                    <ProgressCircle
                      size={20}
                      value={mini.completionFraction}
                    />
                  )}

                  {mini.completionFraction === 0 ? (
                    <div className="text-base font-medium text-muted-foreground leading-none my-1">
                      0:00
                    </div>
                  ) : (
                    <div className="text-base font-medium leading-none my-1">
                      {formatSeconds(
                        mini.completionFraction === 1
                          ? mini.firstSolve
                          : mini.playTime,
                      )}
                    </div>
                  )}
                </div>

                {Array.isArray(mini.constructors) &&
                  mini.constructors.length > 0 && (
                    <p className="text-[10px] text-muted-foreground">
                      {mini.constructors.join(", ")}
                    </p>
                  )}
              </div>
            )}

            {/* MIDI CROSSWORD */}
            {hasMidi && midi && (
              <div>
                <div className="flex items-center gap-1 text-muted-foreground">
                  <span className="text-[11px] font-medium uppercase">
                    midi crossword
                    {midi.title && ` - "${midi.title}"`}
                  </span>

                  <ArrowUpRightFromSquare
                    className="hover:text-foreground cursor-pointer"
                    size={12}
                    onClick={() => goToPuzzle("midi")}
                  />
                </div>

                <div className="flex items-center gap-1">
                  {midi.completionFraction === 1 ? (
                    <CircleCheck
                      size={20}
                      className="text-blue-500"
                    />
                  ) : (
                    <ProgressCircle
                      size={20}
                      value={midi.completionFraction}
                    />
                  )}

                  {midi.completionFraction === 0 ? (
                    <div className="text-base font-medium text-muted-foreground leading-none my-1">
                      0:00
                    </div>
                  ) : (
                    <div className="text-base font-medium leading-none my-1">
                      {formatSeconds(
                        midi.completionFraction === 1
                          ? midi.firstSolve
                          : midi.playTime,
                      )}
                    </div>
                  )}
                </div>

                {Array.isArray(midi.constructors) &&
                  midi.constructors.length > 0 && (
                    <p className="text-[10px] text-muted-foreground">
                      {midi.constructors.join(", ")}
                    </p>
                  )}
              </div>
            )}

            {/* WORDLE */}
            {hasWordle && wordle?.game_data && (
              <div>
                <div className="flex items-center gap-2">
                  <div className="flex items-center gap-1 text-muted-foreground">
                    <span className="text-[11px] font-medium uppercase">
                      wordle
                    </span>

                    <ArrowUpRightFromSquare
                      className="hover:text-foreground cursor-pointer"
                      size={12}
                      onClick={() => goToPuzzle("wordle")}
                    />
                  </div>

                  {wordle.game_data.hardMode && (
                    <Badge
                      className="text-[10px]"
                      variant="outline"
                    >
                      hard mode
                    </Badge>
                  )}

                  <Badge
                    className="text-[10px]"
                    variant={
                      wordle.game_data.status === "FAIL"
                        ? "destructive"
                        : "default"
                    }
                  >
                    {wordle.game_data.status === "WIN"
                      ? "won"
                      : wordle.game_data.status === "FAIL"
                        ? "failed"
                        : "in progress"}{" "}
                    {wordle.game_data.currentRowIndex} / 6
                  </Badge>
                </div>

                <div className="space-y-1 mt-2">
                  {wordle.game_data.boardState.map(
                    (guess, guessIdx) => {
                      if (typeof guess !== "string") {
                        return null;
                      }

                      const solution =
                        wordle.solution.toLowerCase();

                      /*
                       * Wordle guesses should be five characters.
                       * Normalize malformed data rather than letting
                       * it break rendering.
                       */
                      const normalizedGuess = guess
                        .slice(0, 5)
                        .toLowerCase();

                      const isCorrectAnswer =
                        normalizedGuess === solution;

                      const rowHeightClass = isCorrectAnswer
                        ? "h-12"
                        : "h-6 text-xs";

                      const lettersArray = normalizedGuess
                        .padEnd(5, " ")
                        .split("");

                      const solutionLetterCounts: Record<
                        string,
                        number
                      > = {};

                      for (const char of solution) {
                        solutionLetterCounts[char] =
                          (solutionLetterCounts[char] || 0) + 1;
                      }

                      const tileStatuses = Array(5).fill(
                        "empty",
                      ) as Array<
                        "empty" | "correct" | "present" | "absent"
                      >;

                      /*
                       * First pass: exact matches.
                       */
                      lettersArray.forEach(
                        (letter, idx) => {
                          if (
                            letter !== " " &&
                            solution[idx] === letter
                          ) {
                            tileStatuses[idx] = "correct";
                            solutionLetterCounts[letter]--;
                          }
                        },
                      );

                      /*
                       * Second pass: present/absent letters.
                       */
                      lettersArray.forEach(
                        (letter, idx) => {
                          if (
                            letter !== " " &&
                            tileStatuses[idx] !== "correct"
                          ) {
                            if (
                              (solutionLetterCounts[letter] ?? 0) >
                              0
                            ) {
                              tileStatuses[idx] = "present";
                              solutionLetterCounts[letter]--;
                            } else {
                              tileStatuses[idx] = "absent";
                            }
                          }
                        },
                      );

                      return (
                        <div
                          className="grid grid-cols-5 gap-1"
                          key={guessIdx}
                        >
                          {lettersArray.map(
                            (letter, letterIdx) => {
                              const status =
                                tileStatuses[letterIdx];

                              let bgColorClass = "bg-muted";

                              if (status === "correct") {
                                bgColorClass =
                                  "bg-[#538d4e]";
                              } else if (
                                status === "present"
                              ) {
                                bgColorClass =
                                  "bg-[#b59f3b]";
                              }

                              return (
                                <div
                                  className={`flex items-center justify-center rounded uppercase p-1 font-bold ${rowHeightClass} ${bgColorClass}`}
                                  key={letterIdx}
                                >
                                  {letter !== " " ? letter : ""}
                                </div>
                              );
                            },
                          )}
                        </div>
                      );
                    },
                  )}
                </div>
              </div>
            )}

            {/* CONNECTIONS */}
            {hasConnections && connections?.game_data && (
              <div>
                <div className="flex items-center gap-2">
                  <div className="flex items-center gap-1 text-muted-foreground">
                    <span className="text-[11px] font-medium uppercase">
                      connections
                    </span>

                    <ArrowUpRightFromSquare
                      className="hover:text-foreground cursor-pointer"
                      size={12}
                      onClick={() =>
                        goToPuzzle("connections")
                      }
                    />
                  </div>

                  {connections.game_data.puzzleWon ? (
                    <Badge
                      variant="default"
                      className="text-[10px]"
                    >
                      won
                    </Badge>
                  ) : (
                    <Badge
                      variant={
                        connections.game_data.puzzleComplete
                          ? "destructive"
                          : "outline"
                      }
                      className="text-[10px]"
                    >
                      {connections.game_data.puzzleComplete
                        ? "failed"
                        : `${
                            4 -
                            connections.game_data.mistakes
                          } mistakes left`}
                    </Badge>
                  )}
                </div>

                <div className="space-y-1 mt-2">
                  {connections.game_data.guesses.map(
                    (guess, guessIdx) => {
                      if (
                        !guess ||
                        !Array.isArray(guess.cards)
                      ) {
                        return null;
                      }

                      return (
                        <div
                          className="flex items-center gap-1"
                          key={guessIdx}
                        >
                          <div className="grid grid-cols-4 gap-1 flex-1 mr-2">
                            {guess.cards.map(
                              (card, cardIdx) => (
                                <div
                                  key={cardIdx}
                                  className={`h-6 rounded flex items-center justify-center transition-colors ${getConnectionsBgColor(
                                    card.level,
                                  )}`}
                                />
                              ),
                            )}
                          </div>

                          {guess.correct ? (
                            <Check
                              size={12}
                              className="text-green-500 shrink-0"
                            />
                          ) : (
                            <X
                              size={12}
                              className="text-muted-foreground/40 shrink-0"
                            />
                          )}
                        </div>
                      );
                    },
                  )}

                  <div className="mt-2 flex items-center gap-2">
                    {([0, 1, 2, 3] as const).map(
                      (level) => {
                        const solved =
                          connections.game_data!.solvedCategories.some(
                            (cat) => cat.level === level,
                          );

                        const colorClass =
                          getConnectionsBgColor(level);

                        const labels = [
                          "yellow",
                          "green",
                          "blue",
                          "purple",
                        ];

                        return (
                          <div
                            key={level}
                            className="flex items-center gap-1.5"
                          >
                            <div
                              className={`w-2.5 h-2.5 rounded-sm shrink-0 ${colorClass}`}
                            />

                            <span
                              className={`text-[10px] ${
                                solved
                                  ? "text-muted-foreground"
                                  : "text-muted-foreground/40 line-through"
                              }`}
                            >
                              {labels[level]}
                            </span>

                            {solved && (
                              <Check
                                size={9}
                                className="text-green-500 shrink-0"
                              />
                            )}
                          </div>
                        );
                      },
                    )}
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
