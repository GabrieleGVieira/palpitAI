import { useCallback, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import {
  approveJoinRequest,
  listGroupMembers,
  listJoinRequests,
  removeGroupMember,
  transferGroupOwnership,
  updateGroup,
  type Group,
  type GroupMember,
  type JoinRequest,
} from '../services/groups';

const emptyJoinRequests: JoinRequest[] = [];
const emptyMembers: GroupMember[] = [];

export function useGroupAdminScreen(
  group: Group,
  onGroupUpdated: (group: Group) => void,
  onBack: () => void,
) {
  const [name, setName] = useState(group.name);
  const [description, setDescription] = useState(group.description);
  const [isPrivate, setIsPrivate] = useState(group.is_private);
  const [hasUnlimitedParticipants, setHasUnlimitedParticipants] = useState(
    group.participant_limit === null,
  );
  const [participantLimit, setParticipantLimit] = useState(
    group.participant_limit ? String(group.participant_limit) : '20',
  );
  const [approvingUserID, setApprovingUserID] = useState<string | null>(null);
  const [removingUserID, setRemovingUserID] = useState<string | null>(null);
  const [transferringOwnerUserID, setTransferringOwnerUserID] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const queryClient = useQueryClient();
  const requestsQuery = useQuery({
    queryFn: () => listJoinRequests(group.id),
    queryKey: ['groups', group.id, 'join-requests'],
  });
  const membersQuery = useQuery({
    queryFn: () => listGroupMembers(group.id),
    queryKey: ['groups', group.id, 'members'],
  });
  const refetchRequests = requestsQuery.refetch;
  const refetchMembers = membersQuery.refetch;
  const updateMutation = useMutation({
    mutationFn: (payload: Parameters<typeof updateGroup>[1]) => updateGroup(group.id, payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['groups'] });
    },
  });
  const approveMutation = useMutation({
    mutationFn: (userID: string) => approveJoinRequest(group.id, userID),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['groups'] });
      await queryClient.invalidateQueries({ queryKey: ['groups', group.id, 'join-requests'] });
    },
  });
  const removeMemberMutation = useMutation({
    mutationFn: (userID: string) => removeGroupMember(group.id, userID),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['groups'] });
      await queryClient.invalidateQueries({ queryKey: ['groups', group.id, 'members'] });
      await queryClient.invalidateQueries({ queryKey: ['groups', group.id, 'ranking'] });
    },
  });
  const transferOwnershipMutation = useMutation({
    mutationFn: (userID: string) => transferGroupOwnership(group.id, userID),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['groups'] });
      await queryClient.invalidateQueries({ queryKey: ['groups', group.id, 'members'] });
    },
  });

  const loadRequests = useCallback(async () => {
    setError(null);
    const result = await refetchRequests();
    const nextError = queryErrorMessage(result.error, 'Não foi possível carregar solicitações.');
    if (nextError) {
      setError(nextError);
    }
  }, [refetchRequests]);

  const loadMembers = useCallback(async () => {
    setError(null);
    const result = await refetchMembers();
    const nextError = queryErrorMessage(result.error, 'Não foi possível carregar participantes.');
    if (nextError) {
      setError(nextError);
    }
  }, [refetchMembers]);

  async function handleSaveGroup() {
    setError(null);
    setSuccessMessage(null);

    if (!name.trim()) {
      setError('Informe o nome do grupo.');
      return;
    }

    if (!hasUnlimitedParticipants && Number(participantLimit) < 2) {
      setError('O limite precisa ser maior que 1.');
      return;
    }

    try {
      const updatedGroup = await updateMutation.mutateAsync({
        description,
        has_unlimited_participants: hasUnlimitedParticipants,
        is_private: isPrivate,
        name,
        participant_limit: hasUnlimitedParticipants ? null : Number(participantLimit),
      });

      onGroupUpdated({ ...group, ...updatedGroup });
      setSuccessMessage('Grupo atualizado.');
      onBack();
    } catch (saveError) {
      setError(queryErrorMessage(saveError, 'Não foi possível atualizar o grupo.'));
    }
  }

  async function handleApprove(request: JoinRequest) {
    setError(null);
    setSuccessMessage(null);
    setApprovingUserID(request.user_id);

    try {
      await approveMutation.mutateAsync(request.user_id);
      onGroupUpdated({
        ...group,
        member_count: group.member_count + 1,
        pending_requests_count: Math.max(group.pending_requests_count - 1, 0),
      });
      await queryClient.invalidateQueries({ queryKey: ['groups', group.id, 'members'] });
      setSuccessMessage('Solicitação aprovada.');
    } catch (approveError) {
      setError(queryErrorMessage(approveError, 'Não foi possível aprovar a solicitação.'));
    } finally {
      setApprovingUserID(null);
    }
  }

  async function handleRemoveMember(member: GroupMember) {
    setError(null);
    setSuccessMessage(null);
    setRemovingUserID(member.user_id);

    try {
      await removeMemberMutation.mutateAsync(member.user_id);
      onGroupUpdated({
        ...group,
        member_count: Math.max(group.member_count - 1, 1),
      });
      setSuccessMessage('Participante removido.');
    } catch (removeError) {
      setError(queryErrorMessage(removeError, 'Não foi possível remover o participante.'));
    } finally {
      setRemovingUserID(null);
    }
  }

  async function handleTransferOwnership(member: GroupMember) {
    setError(null);
    setSuccessMessage(null);
    setTransferringOwnerUserID(member.user_id);

    try {
      await transferOwnershipMutation.mutateAsync(member.user_id);
      onGroupUpdated({
        ...group,
        owner_id: member.user_id,
        role: 'member',
      });
      setSuccessMessage('Propriedade do grupo transferida.');
      onBack();
    } catch (transferError) {
      setError(queryErrorMessage(transferError, 'Não foi possível transferir a propriedade.'));
    } finally {
      setTransferringOwnerUserID(null);
    }
  }

  return {
    approvingUserID,
    description,
    error:
      error !== null
        ? error
        : queryErrorMessage(
            requestsQuery.isError ? requestsQuery.error : null,
            'Não foi possível carregar solicitações.',
          ),
    hasUnlimitedParticipants,
    isLoadingRequests: requestsQuery.isLoading,
    isLoadingMembers: membersQuery.isLoading,
    isPrivate,
    isSaving: updateMutation.isPending,
    loadRequests,
    loadMembers,
    members: Array.isArray(membersQuery.data) ? membersQuery.data : emptyMembers,
    name,
    participantLimit,
    removingUserID,
    requests: Array.isArray(requestsQuery.data) ? requestsQuery.data : emptyJoinRequests,
    setDescription,
    setHasUnlimitedParticipants,
    setIsPrivate,
    setName,
    setParticipantLimit,
    setSuccessMessage,
    successMessage,
    transferringOwnerUserID,
    handleApprove,
    handleRemoveMember,
    handleSaveGroup,
    handleTransferOwnership,
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
