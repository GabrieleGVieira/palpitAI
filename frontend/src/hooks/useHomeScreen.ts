import { useCallback, useEffect, useState } from 'react';
import { getUserScore, joinGroup, listGroups, type Group } from '../services/groups';
import { connectRealtime } from '../services/realtime';
import { notificationMessageFromEvent } from '../utils/realtimeNotifications';

export function useHomeScreen() {
  const [groups, setGroups] = useState<Group[]>([]);
  const [totalPoints, setTotalPoints] = useState(0);
  const [isLoadingGroups, setIsLoadingGroups] = useState(true);
  const [isLoadingScore, setIsLoadingScore] = useState(true);
  const [groupsError, setGroupsError] = useState<string | null>(null);
  const [scoreError, setScoreError] = useState<string | null>(null);
  const [inviteCode, setInviteCode] = useState('');
  const [joinError, setJoinError] = useState<string | null>(null);
  const [joinSuccess, setJoinSuccess] = useState<string | null>(null);
  const [isJoiningGroup, setIsJoiningGroup] = useState(false);
  const [notificationMessage, setNotificationMessage] = useState<string | null>(null);

  const loadGroups = useCallback(async () => {
    setGroupsError(null);
    setIsLoadingGroups(true);

    try {
      const nextGroups = await listGroups();
      setGroups(nextGroups);
    } catch (error) {
      setGroupsError(
        error instanceof Error ? error.message : 'Não foi possível carregar seus grupos.',
      );
    } finally {
      setIsLoadingGroups(false);
    }
  }, []);

  const loadScore = useCallback(async () => {
    setScoreError(null);
    setIsLoadingScore(true);

    try {
      const score = await getUserScore();
      setTotalPoints(score.total_points);
    } catch (error) {
      setScoreError(
        error instanceof Error ? error.message : 'Não foi possível carregar sua pontuação.',
      );
    } finally {
      setIsLoadingScore(false);
    }
  }, []);

  const refreshHome = useCallback(async () => {
    await Promise.all([loadGroups(), loadScore()]);
  }, [loadGroups, loadScore]);

  useEffect(() => {
    void refreshHome();
  }, [refreshHome]);

  useEffect(() => {
    let cleanup: (() => void) | undefined;
    let isMounted = true;

    connectRealtime({
      onEvent: (event) => {
        if (event.name === 'ranking.updated' || event.name === 'match.finished') {
          setNotificationMessage(notificationMessageFromEvent(event));
          void refreshHome();
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
        // Home remains usable through REST even when realtime is unavailable.
      });

    return () => {
      isMounted = false;
      cleanup?.();
    };
  }, [refreshHome]);

  useEffect(() => {
    if (!notificationMessage) {
      return;
    }

    const timer = setTimeout(() => setNotificationMessage(null), 5000);
    return () => clearTimeout(timer);
  }, [notificationMessage]);

  const handleJoinGroup = useCallback(async () => {
    setJoinError(null);
    setJoinSuccess(null);

    if (!inviteCode.trim()) {
      setJoinError('Informe o codigo do grupo.');
      return;
    }

    setIsJoiningGroup(true);

    try {
      const response = await joinGroup(inviteCode);
      setInviteCode('');
      setJoinSuccess(
        response.membership_status === 'pending'
          ? 'Solicitação enviada. Aguarde a aprovação do dono do grupo.'
          : 'Você entrou no grupo.',
      );
      await refreshHome();
    } catch (error) {
      setJoinError(error instanceof Error ? error.message : 'Não foi possível entrar no grupo.');
    } finally {
      setIsJoiningGroup(false);
    }
  }, [inviteCode, refreshHome]);

  return {
    groups,
    totalPoints,
    isLoadingGroups,
    isLoadingScore,
    groupsError,
    scoreError,
    inviteCode,
    setInviteCode,
    joinError,
    joinSuccess,
    isJoiningGroup,
    notificationMessage,
    refreshHome,
    handleJoinGroup,
  };
}
