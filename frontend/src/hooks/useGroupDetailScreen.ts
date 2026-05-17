import { useCallback, useEffect, useState } from 'react';

import { connectRealtime } from '../services/realtime';
import {
  listGroupMatches,
  listGroupRanking,
  savePrediction,
  type Group,
  type GroupMatch,
  type RankingEntry,
} from '../services/groups';
import { notificationMessageFromEvent } from '../utils/realtimeNotifications';

export type ScoreDraft = {
  awayScore: string;
  homeScore: string;
};

export type GroupDetailTab = 'matches' | 'ranking';

export function useGroupDetailScreen(group: Group) {
  const [matches, setMatches] = useState<GroupMatch[]>([]);
  const [ranking, setRanking] = useState<RankingEntry[]>([]);
  const [drafts, setDrafts] = useState<Record<string, ScoreDraft>>({});
  const [activeTab, setActiveTab] = useState<GroupDetailTab>('matches');
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingRanking, setIsLoadingRanking] = useState(false);
  const [savingMatchID, setSavingMatchID] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [notificationMessage, setNotificationMessage] = useState<string | null>(null);
  const [rankingError, setRankingError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const loadMatches = useCallback(
    async (showLoading = true) => {
      setError(null);
      if (showLoading) {
        setIsLoading(true);
      }

      try {
        const nextMatches = await listGroupMatches(group.id);
        setMatches(nextMatches);
        setDrafts(buildDrafts(nextMatches));
      } catch (loadError) {
        setError(
          loadError instanceof Error ? loadError.message : 'Não foi possível carregar jogos.',
        );
      } finally {
        if (showLoading) {
          setIsLoading(false);
        }
      }
    },
    [group.id],
  );

  const loadRanking = useCallback(
    async (showLoading = true) => {
      setRankingError(null);
      if (showLoading) {
        setIsLoadingRanking(true);
      }

      try {
        const nextRanking = await listGroupRanking(group.id);
        setRanking(nextRanking);
      } catch (loadError) {
        setRankingError(
          loadError instanceof Error ? loadError.message : 'Não foi possível carregar o ranking.',
        );
      } finally {
        if (showLoading) {
          setIsLoadingRanking(false);
        }
      }
    },
    [group.id],
  );

  useEffect(() => {
    void loadMatches();
  }, [loadMatches]);

  useEffect(() => {
    if (activeTab === 'ranking') {
      void loadRanking();
    }
  }, [activeTab, loadRanking]);

  useEffect(() => {
    let cleanup: (() => void) | undefined;
    let isMounted = true;

    connectRealtime({
      groupID: group.id,
      onEvent: (event) => {
        const nextNotification = notificationMessageFromEvent(event, group.name);

        if (nextNotification) {
          setNotificationMessage(nextNotification);
        }

        if (
          event.name === 'match.updated' ||
          event.name === 'match.finished' ||
          event.name === 'match.goal'
        ) {
          void loadMatches(false);
        }

        if (event.name === 'ranking.updated' || event.name === 'match.finished') {
          void loadRanking(false);
        }
      },
    })
      .then((nextCleanup) => {
        if (isMounted) {
          cleanup = nextCleanup;
        } else {
          nextCleanup();
        }
      })
      .catch(() => {
        // Realtime falha silenciosamente; REST ainda funcionará.
      });

    return () => {
      isMounted = false;
      cleanup?.();
    };
  }, [group.id, group.name, loadMatches, loadRanking]);

  useEffect(() => {
    if (!notificationMessage) {
      return;
    }

    const timer = setTimeout(() => setNotificationMessage(null), 5000);
    return () => clearTimeout(timer);
  }, [notificationMessage]);

  function updateDraft(matchID: string, key: keyof ScoreDraft, value: string) {
    setDrafts((currentDrafts) => ({
      ...currentDrafts,
      [matchID]: {
        ...(currentDrafts[matchID] ?? { awayScore: '', homeScore: '' }),
        [key]: value.replace(/\D/g, '').slice(0, 2),
      },
    }));
  }

  async function handleSavePrediction(match: GroupMatch) {
    const draft = drafts[match.id];
    setError(null);
    setSuccessMessage(null);

    if (!draft?.homeScore || !draft.awayScore) {
      setError('Informe os dois placares para salvar o palpite.');
      return;
    }

    setSavingMatchID(match.id);

    try {
      const prediction = await savePrediction(group.id, match.id, {
        away_score: Number(draft.awayScore),
        home_score: Number(draft.homeScore),
      });

      setMatches((currentMatches) =>
        currentMatches.map((currentMatch) =>
          currentMatch.id === match.id
            ? {
                ...currentMatch,
                my_prediction: prediction,
              }
            : currentMatch,
        ),
      );
      await loadRanking();
      setSuccessMessage('Palpite salvo.');
    } catch (saveError) {
      setError(
        saveError instanceof Error ? saveError.message : 'Não foi possível salvar o palpite.',
      );
    } finally {
      setSavingMatchID(null);
    }
  }

  return {
    activeTab,
    drafts,
    error,
    isLoading,
    isLoadingRanking,
    loadMatches,
    loadRanking,
    matches,
    notificationMessage,
    ranking,
    rankingError,
    savePrediction: handleSavePrediction,
    setActiveTab,
    savingMatchID,
    successMessage,
    updateDraft,
  };
}

export function buildDrafts(matches: GroupMatch[]) {
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

export function formatDate(value: string) {
  return new Intl.DateTimeFormat('pt-BR', {
    dateStyle: 'short',
    timeStyle: 'short',
  }).format(new Date(value));
}

export function formatUserID(userID: string) {
  if (userID.length <= 12) {
    return userID;
  }

  return `${userID.slice(0, 8)}...${userID.slice(-4)}`;
}

export function formatMatchStatus(status: GroupMatch['status']) {
  const statusLabels: Record<GroupMatch['status'], string> = {
    cancelled: 'Cancelado',
    finished: 'Encerrado',
    live: 'Ao vivo',
    postponed: 'Adiado',
    scheduled: 'Agendado',
  };

  return statusLabels[status];
}

export function formatMatchStage(stage: GroupMatch['stage']) {
  const stageLabels: Record<GroupMatch['stage'], string> = {
    GROUP_STAGE: 'Fase de grupos',
    LAST_32: 'Mata-mata inicial',
    LAST_16: 'Oitavas de final',
    QUARTER_FINALS: 'Quartas de final',
    SEMI_FINALS: 'Semi-finais',
    FINAL: 'Final',
  };

  return stageLabels[stage];
}
