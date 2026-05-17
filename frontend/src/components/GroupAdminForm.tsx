import { Pressable, StyleSheet, Text, TextInput, View } from 'react-native';

import { AuthInputField } from './global/AuthInputField';
import { SwitchBox } from './global/SwitchBox';
import { ParticipantsCard } from './global/ParticipantsCard';

type GroupAdminFormProps = {
  description: string;
  hasUnlimitedParticipants: boolean;
  isPrivate: boolean;
  isSaving: boolean;
  name: string;
  participantLimit: string;
  onSave: () => void;
  setDescription: (value: string) => void;
  setHasUnlimitedParticipants: (value: boolean) => void;
  setIsPrivate: (value: boolean) => void;
  setName: (value: string) => void;
  setParticipantLimit: (value: string) => void;
};

export function GroupAdminForm({
  description,
  hasUnlimitedParticipants,
  isPrivate,
  isSaving,
  name,
  participantLimit,
  onSave,
  setDescription,
  setHasUnlimitedParticipants,
  setIsPrivate,
  setName,
  setParticipantLimit,
}: GroupAdminFormProps) {
  return (
    <View style={styles.card}>
      <Text style={styles.cardTitle}>Informações</Text>

      <AuthInputField
        autoCapitalize="words"
        keyboardType="default"
        label="Nome"
        onChangeText={setName}
        placeholder="Nome do grupo"
        value={name}
      />

      <View style={styles.fieldGroup}>
        <Text style={styles.label}>Descrição</Text>
        <TextInput
          multiline
          onChangeText={setDescription}
          placeholder="Descrição do grupo"
          placeholderTextColor="#7c8898"
          style={[styles.input, styles.textArea]}
          textAlignVertical="top"
          value={description}
        />
      </View>

      <ParticipantsCard
        hasUnlimitedParticipants={hasUnlimitedParticipants}
        participantLimit={participantLimit}
        setHasUnlimitedParticipants={setHasUnlimitedParticipants}
        setParticipantLimit={setParticipantLimit}
      />

      <SwitchBox
        title="Privado"
        subtitle="Novos membros precisam de aprovação"
        value={isPrivate}
        onPress={setIsPrivate}
      />

      <Pressable
        disabled={isSaving}
        onPress={onSave}
        style={[styles.primaryButton, isSaving && styles.buttonDisabled]}>
        <Text style={styles.primaryButtonText}>{isSaving ? 'Salvando...' : 'Salvar'}</Text>
      </Pressable>
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
  cardTitle: {
    color: '#123d2a',
    fontSize: 18,
    fontWeight: '800',
  },
  fieldGroup: {
    gap: 8,
  },
  label: {
    color: '#183f2d',
    fontSize: 14,
    fontWeight: '700',
  },
  input: {
    backgroundColor: '#f5f8ef',
    borderColor: '#cfe0c9',
    borderRadius: 8,
    borderWidth: 1,
    color: '#183f2d',
    fontSize: 16,
    minHeight: 52,
    paddingHorizontal: 14,
  },
  inputDisabled: {
    backgroundColor: '#edf3e8',
    color: '#7c8898',
  },
  textArea: {
    minHeight: 96,
    paddingTop: 12,
  },
  row: {
    flexDirection: 'row',
    gap: 12,
  },
  limitField: {
    flex: 1,
  },
  primaryButton: {
    alignItems: 'center',
    backgroundColor: '#1f7a4a',
    borderRadius: 8,
    justifyContent: 'center',
    minHeight: 52,
  },
  primaryButtonText: {
    color: '#ffffff',
    fontSize: 15,
    fontWeight: '800',
  },
  buttonDisabled: {
    opacity: 0.72,
  },
});
