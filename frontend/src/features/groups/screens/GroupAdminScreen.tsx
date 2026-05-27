import { StatusBar } from 'expo-status-bar';
import { Alert, ScrollView, StyleSheet, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

import { BackButton } from '../../../shared/components/BackButton';
import { GroupAdminForm } from '../components/admin/GroupAdminForm';
import { GroupAdminHeader } from '../components/admin/GroupAdminHeader';
import { GroupAdminMembers } from '../components/admin/GroupAdminMembers';
import { GroupAdminRequests } from '../components/admin/GroupAdminRequests';
import { useGroupAdminScreen } from '../hooks/useGroupAdminScreen';
import type { Group } from '../services/groups';

type GroupAdminScreenProps = {
  group: Group;
  onBack: () => void;
  onGroupUpdated: (group: Group) => void;
};

export function GroupAdminScreen({ group, onBack, onGroupUpdated }: GroupAdminScreenProps) {
  const {
    approvingUserID,
    description,
    error,
    hasUnlimitedParticipants,
    isLoadingMembers,
    isLoadingRequests,
    isPrivate,
    isSaving,
    loadMembers,
    loadRequests,
    members,
    name,
    participantLimit,
    removingUserID,
    requests,
    setDescription,
    setHasUnlimitedParticipants,
    setIsPrivate,
    setName,
    setParticipantLimit,
    successMessage,
    handleApprove,
    handleRemoveMember,
    handleSaveGroup,
  } = useGroupAdminScreen(group, onGroupUpdated, onBack);

  function confirmRemoveMember(member: (typeof members)[number]) {
    const name = member.display_name || `Usuário ${member.user_id.slice(0, 8)}`;

    Alert.alert(
      'Remover participante',
      `Você tem certeza que deseja remover ${name} deste grupo?`,
      [
        { style: 'cancel', text: 'Cancelar' },
        { onPress: () => handleRemoveMember(member), style: 'destructive', text: 'Remover' },
      ],
    );
  }

  return (
    <SafeAreaView style={styles.safeArea}>
      <StatusBar style="dark" />
      <ScrollView contentContainerStyle={styles.container} showsVerticalScrollIndicator={false}>
        <View style={styles.backgroundMarkerTop} />
        <View style={styles.backgroundCircle} />

        <BackButton onPress={onBack} />

        <GroupAdminHeader groupName={group.name} error={error} successMessage={successMessage} />

        <GroupAdminForm
          description={description}
          hasUnlimitedParticipants={hasUnlimitedParticipants}
          isPrivate={isPrivate}
          isSaving={isSaving}
          name={name}
          participantLimit={participantLimit}
          onSave={handleSaveGroup}
          setDescription={setDescription}
          setHasUnlimitedParticipants={setHasUnlimitedParticipants}
          setIsPrivate={setIsPrivate}
          setName={setName}
          setParticipantLimit={setParticipantLimit}
        />

        <GroupAdminRequests
          approvingUserID={approvingUserID}
          isLoadingRequests={isLoadingRequests}
          loadRequests={loadRequests}
          onApprove={handleApprove}
          requests={requests}
        />

        <GroupAdminMembers
          isLoadingMembers={isLoadingMembers}
          loadMembers={loadMembers}
          members={members}
          onRemove={confirmRemoveMember}
          removingUserID={removingUserID}
        />
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
    backgroundColor: '#f5f8ef',
  },
  container: {
    backgroundColor: '#f5f8ef',
    flexGrow: 1,
    gap: 20,
    paddingHorizontal: 24,
    paddingVertical: 32,
  },
  backgroundMarkerTop: {
    borderColor: 'rgba(255, 255, 255, 0.68)',
    borderRadius: 8,
    borderWidth: 2,
    height: 116,
    left: 24,
    position: 'absolute',
    right: 24,
    top: -42,
  },
  backgroundCircle: {
    borderColor: 'rgba(32, 111, 67, 0.12)',
    borderRadius: 140,
    borderWidth: 2,
    height: 280,
    position: 'absolute',
    right: -128,
    top: 104,
    width: 280,
  },
});
