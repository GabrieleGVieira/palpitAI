import { useCallback } from 'react';
import { useQuery } from '@tanstack/react-query';

import { getUserScore, listGroups, type Group } from '../services/groups';

const emptyGroups: Group[] = [];

export function useHomeData() {
  const groupsQuery = useQuery({
    queryFn: listGroups,
    queryKey: ['groups'],
  });
  const scoreQuery = useQuery({
    queryFn: getUserScore,
    queryKey: ['me', 'score'],
  });
  const refetchGroups = groupsQuery.refetch;
  const refetchScore = scoreQuery.refetch;

  const refreshHome = useCallback(async () => {
    await Promise.all([refetchGroups(), refetchScore()]);
  }, [refetchGroups, refetchScore]);

  return {
    groups: Array.isArray(groupsQuery.data) ? groupsQuery.data : emptyGroups,
    groupsError: queryErrorMessage(groupsQuery.isError ? groupsQuery.error : null),
    isLoadingGroups: groupsQuery.isLoading,
    isLoadingScore: scoreQuery.isLoading,
    refreshHome,
    scoreError: queryErrorMessage(scoreQuery.isError ? scoreQuery.error : null),
    totalPoints: scoreQuery.data?.total_points ?? 0,
  };
}

function queryErrorMessage(error: unknown) {
  if (error == null) {
    return null;
  }

  if (typeof error === 'string') {
    return error.trim() || 'Não foi possível carregar as informações.';
  }

  if (typeof error === 'object' && 'message' in error) {
    const message = (error as { message?: unknown }).message;
    if (typeof message === 'string' && message.trim()) {
      return message;
    }
  }

  return 'Não foi possível carregar as informações.';
}
