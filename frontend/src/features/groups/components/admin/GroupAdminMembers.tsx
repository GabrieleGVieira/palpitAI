import { Pressable, StyleSheet, Text, View } from 'react-native';

import { EmptyBox } from '../../../../shared/components/EmptyBox';
import { FinishButton } from '../../../../shared/components/FinishButton';
import { LoadingIndicator } from '../../../../shared/components/LoadingIndicator';
import type { GroupMember } from '../../services/groups';

type GroupAdminMembersProps = {
  isLoadingMembers: boolean;
  loadMembers: () => void;
  members: GroupMember[];
  onRemove: (member: GroupMember) => void;
  onTransferOwnership: (member: GroupMember) => void;
  removingUserID: string | null;
  transferringOwnerUserID: string | null;
};

export function GroupAdminMembers({
  isLoadingMembers,
  loadMembers,
  members,
  onRemove,
  onTransferOwnership,
  removingUserID,
  transferringOwnerUserID,
}: GroupAdminMembersProps) {
  return (
    <View style={styles.card}>
      <View style={styles.header}>
        <View>
          <Text style={styles.cardTitle}>Participantes</Text>
          <Text style={styles.cardSubtitle}>Membros ativos do grupo</Text>
        </View>
        <Pressable onPress={loadMembers} style={styles.refreshButton}>
          <Text style={styles.refreshButtonText}>Atualizar</Text>
        </Pressable>
      </View>

      {isLoadingMembers ? <LoadingIndicator text="Carregando..." /> : null}

      {!isLoadingMembers && members.length === 0 ? (
        <EmptyBox title="Nenhum participante." text="Nenhum membro ativo encontrado." />
      ) : null}

      {members.map((member) => {
        const isOwner = member.role === 'owner';

        return (
          <View key={member.user_id} style={styles.memberRow}>
            <View style={styles.memberInfo}>
              <Text style={styles.memberName}>
                {member.display_name || `Usuário ${member.user_id.slice(0, 8)}`}
              </Text>
              <Text style={styles.memberMeta}>{isOwner ? 'Dono do grupo' : 'Participante'}</Text>
            </View>

            {isOwner ? (
              <Text style={styles.ownerBadge}>Owner</Text>
            ) : (
              <View style={styles.memberActions}>
                <FinishButton
                  isLoading={transferringOwnerUserID === member.user_id}
                  loadingLabel="Transferindo..."
                  onPress={() => onTransferOwnership(member)}
                  waitingLabel="Tornar dono"
                />
                <FinishButton
                  isLoading={removingUserID === member.user_id}
                  loadingLabel="Removendo..."
                  onPress={() => onRemove(member)}
                  waitingLabel="Remover"
                />
              </View>
            )}
          </View>
        );
      })}
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: '#ffffff',
    borderColor: '#cfe0c9',
    borderRadius: 8,
    borderWidth: 1,
    gap: 16,
    padding: 16,
  },
  header: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  cardTitle: {
    color: '#123d2a',
    fontSize: 18,
    fontWeight: '800',
  },
  cardSubtitle: {
    color: '#486654',
    fontSize: 13,
    marginTop: 4,
  },
  refreshButton: {
    backgroundColor: '#f5f8ef',
    borderColor: '#cfe0c9',
    borderRadius: 8,
    borderWidth: 1,
    paddingHorizontal: 12,
    paddingVertical: 8,
  },
  refreshButtonText: {
    color: '#1f7a4a',
    fontSize: 13,
    fontWeight: '800',
  },
  memberRow: {
    alignItems: 'center',
    borderTopColor: '#edf3e8',
    borderTopWidth: 1,
    flexDirection: 'row',
    gap: 10,
    justifyContent: 'space-between',
    paddingTop: 12,
  },
  memberInfo: {
    flex: 1,
  },
  memberActions: {
    gap: 8,
  },
  memberName: {
    color: '#183f2d',
    fontSize: 14,
    fontWeight: '800',
  },
  memberMeta: {
    color: '#486654',
    fontSize: 12,
    marginTop: 3,
  },
  ownerBadge: {
    backgroundColor: '#edf3e8',
    borderRadius: 8,
    color: '#1f7a4a',
    fontSize: 12,
    fontWeight: '800',
    overflow: 'hidden',
    paddingHorizontal: 10,
    paddingVertical: 8,
  },
});
