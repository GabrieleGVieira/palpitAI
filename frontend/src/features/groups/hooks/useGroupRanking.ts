import { useCallback } from 'react';
import { useQuery } from '@tanstack/react-query';

import { listGroupRanking, type RankingEntry } from '../services/groups';

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
    // Garante que se não for um array válido, retorna um array vazio imutável
    ranking: rankingQuery.data || ([] as RankingEntry[]),
    // Uma checagem mais segura que evita falsos positivos no primeiro render
    rankingError: rankingQuery.error
      ? (rankingQuery.error as Error).message || 'Erro ao carregar ranking'
      : null,
  };
}
