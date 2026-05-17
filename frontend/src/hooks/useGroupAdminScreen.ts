import { useCallback, useEffect, useState } from 'react';

import {
  approveJoinRequest,
  listJoinRequests,
  updateGroup,
  type Group,
  type JoinRequest,
} from '../services/groups';

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
  const [requests, setRequests] = useState<JoinRequest[]>([]);
  const [isLoadingRequests, setIsLoadingRequests] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [approvingUserID, setApprovingUserID] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const loadRequests = useCallback(async () => {
    setError(null);
    setIsLoadingRequests(true);

    try {
      const nextRequests = await listJoinRequests(group.id);
      setRequests(nextRequests);
    } catch (loadError) {
      setError(
        loadError instanceof Error ? loadError.message : 'Não foi possível carregar solicitações.',
      );
    } finally {
      setIsLoadingRequests(false);
    }
  }, [group.id]);

  useEffect(() => {
    void loadRequests();
  }, [loadRequests]);

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

    setIsSaving(true);

    try {
      const updatedGroup = await updateGroup(group.id, {
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
      setError(
        saveError instanceof Error ? saveError.message : 'Não foi possível atualizar o grupo.',
      );
    } finally {
      setIsSaving(false);
    }
  }

  async function handleApprove(request: JoinRequest) {
    setError(null);
    setSuccessMessage(null);
    setApprovingUserID(request.user_id);

    try {
      await approveJoinRequest(group.id, request.user_id);
      setRequests((currentRequests) =>
        currentRequests.filter((currentRequest) => currentRequest.user_id !== request.user_id),
      );
      onGroupUpdated({
        ...group,
        member_count: group.member_count + 1,
        pending_requests_count: Math.max(group.pending_requests_count - 1, 0),
      });
      setSuccessMessage('Solicitação aprovada.');
    } catch (approveError) {
      setError(
        approveError instanceof Error
          ? approveError.message
          : 'Não foi possível aprovar a solicitação.',
      );
    } finally {
      setApprovingUserID(null);
    }
  }

  return {
    approvingUserID,
    description,
    error,
    hasUnlimitedParticipants,
    isLoadingRequests,
    isPrivate,
    isSaving,
    loadRequests,
    name,
    participantLimit,
    requests,
    setDescription,
    setHasUnlimitedParticipants,
    setIsPrivate,
    setName,
    setParticipantLimit,
    setRequests,
    setSuccessMessage,
    successMessage,
    handleApprove,
    handleSaveGroup,
  };
}
