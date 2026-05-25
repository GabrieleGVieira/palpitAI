import { useCallback } from 'react';
import { useQuery } from '@tanstack/react-query';

import { listGroupRanking, type RankingEntry } from '../services/groups';

const emptyRanking: RankingEntry[] = [];

export function useGroupRanking(groupID: string) {
  const rankingQuery = useQuery({
    enabled: false,
    queryFn: () => listGroupRanking(groupID),
    queryKey: ['groups', groupID, 'ranking'],
  });
  const refetchRanking = rankingQuery.refetch;

  const loadRanking = useCallback(
    async (showLoading = true) => {
      await refetchRanking({ cancelRefetch: showLoading });
    },
    [refetchRanking],
  );

  return {
    isLoadingRanking: rankingQuery.isFetching,
    loadRanking,
    ranking: Array.isArray(rankingQuery.data) ? rankingQuery.data : emptyRanking,
    rankingError: queryErrorMessage(
      rankingQuery.isError ? rankingQuery.error : null,
      'Erro ao carregar ranking',
    ),
  };
}

function queryErrorMessage(error: unknown, fallback: string) {
  if (error == null) {
    return null;
  }

  if (typeof error === 'string') {
    return error.trim() || fallback;
  }

  if (typeof error === 'object' && 'message' in error) {
    const message = (error as { message?: unknown }).message;
    if (typeof message === 'string' && message.trim()) {
      return message;
    }
  }

  return fallback;
}
