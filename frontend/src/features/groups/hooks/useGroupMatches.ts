import { useCallback, useEffect, useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';

import { listGroupMatches, type GroupMatch } from '../services/groups';
import type { RealtimeEvent } from '../../realtime/types';
import type { ScoreDraft } from '../types';

export function useGroupMatches(groupID: string) {
  const [drafts, setDrafts] = useState<Record<string, ScoreDraft>>({});
  const [error, setError] = useState<string | null>(null);
  const queryClient = useQueryClient();
  const matchesQuery = useQuery({
    queryFn: () => listGroupMatches(groupID),
    queryKey: ['groups', groupID, 'matches'],
  });
  const refetchMatches = matchesQuery.refetch;

  const matches = matchesQuery.data ?? [];

  useEffect(() => {
    if (matchesQuery.data) {
      setDrafts(buildDrafts(matchesQuery.data));
    }
  }, [matchesQuery.data]);

  const loadMatches = useCallback(
    async (showLoading = true) => {
      setError(null);
      const result = await refetchMatches({ cancelRefetch: showLoading });
      if (result.error) {
        setError(errorMessage(result.error, 'Não foi possível carregar jogos.'));
      }
    },
    [refetchMatches],
  );

  const updateDraft = useCallback((matchID: string, key: keyof ScoreDraft, value: string) => {
    setDrafts((currentDrafts) => ({
      ...currentDrafts,
      [matchID]: {
        ...(currentDrafts[matchID] ?? { awayScore: '', homeScore: '' }),
        [key]: value.replace(/\D/g, '').slice(0, 2),
      },
    }));
  }, []);

  const updateMatchPrediction = useCallback(
    (matchID: string, prediction: GroupMatch['my_prediction']) => {
      queryClient.setQueryData<GroupMatch[]>(['groups', groupID, 'matches'], (currentMatches) =>
        (currentMatches ?? []).map((currentMatch) =>
          currentMatch.id === matchID
            ? {
                ...currentMatch,
                my_prediction: prediction,
              }
            : currentMatch,
        ),
      );
    },
    [groupID, queryClient],
  );

  const updateMatchFromRealtime = useCallback(
    (event: RealtimeEvent) => {
      const matchID = stringValue(event.payload.match_id);
      if (!matchID) {
        return false;
      }

      let wasUpdated = false;
      queryClient.setQueryData<GroupMatch[]>(['groups', groupID, 'matches'], (currentMatches) =>
        (currentMatches ?? []).map((currentMatch) => {
          if (currentMatch.id !== matchID) {
            return currentMatch;
          }

          wasUpdated = true;
          const status = matchStatusValue(event.payload.status) ?? currentMatch.status;
          const homeScore =
            numberValue(event.payload.final_home_score) ??
            numberValue(event.payload.home_score) ??
            currentMatch.final_home_score;
          const awayScore =
            numberValue(event.payload.final_away_score) ??
            numberValue(event.payload.away_score) ??
            currentMatch.final_away_score;

          return {
            ...currentMatch,
            final_away_score: awayScore,
            final_home_score: homeScore,
            finished_at:
              status === 'finished'
                ? (stringValue(event.payload.finished_at) ??
                  currentMatch.finished_at ??
                  new Date().toISOString())
                : currentMatch.finished_at,
            status,
          };
        }),
      );

      return wasUpdated;
    },
    [groupID, queryClient],
  );

  return {
    drafts,
    error:
      error ||
      (matchesQuery.error
        ? errorMessage(matchesQuery.error, 'Não foi possível carregar jogos.')
        : null),
    isLoading: matchesQuery.isLoading,
    loadMatches,
    matches,
    setError,
    updateDraft,
    updateMatchFromRealtime,
    updateMatchPrediction,
  };
}

function errorMessage(error: unknown, fallback: string) {
  return error instanceof Error ? error.message : fallback;
}

function buildDrafts(matches: GroupMatch[]) {
  return Object.fromEntries(
    matches.map((match) => [
      match.id,
      {
        awayScore: match.my_prediction ? String(match.my_prediction.away_score) : '',
        homeScore: match.my_prediction ? String(match.my_prediction.home_score) : '',
      },
    ]),
  ) as Record<string, ScoreDraft>;
}

function stringValue(value: unknown) {
  return typeof value === 'string' && value.length > 0 ? value : null;
}

function numberValue(value: unknown) {
  if (typeof value === 'number') {
    return value;
  }

  return null;
}

function matchStatusValue(value: unknown): GroupMatch['status'] | null {
  if (
    value === 'scheduled' ||
    value === 'live' ||
    value === 'finished' ||
    value === 'postponed' ||
    value === 'cancelled'
  ) {
    return value;
  }

  return null;
}
