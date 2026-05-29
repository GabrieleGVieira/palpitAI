import { Pressable, StyleSheet, Text, View } from 'react-native';

import { colors } from '../../../shared/theme';

type AccountSettingsCardProps = {
  onDeleteAccount: () => void;
};

export function AccountSettingsCard({ onDeleteAccount }: AccountSettingsCardProps) {
  return (
    <View style={styles.card}>
      <Text style={styles.title}>Configurações</Text>
      <Pressable onPress={onDeleteAccount} style={styles.row}>
        <View>
          <Text style={styles.deleteText}>Excluir conta</Text>
          <Text style={styles.description}>Remover sua conta e dados associados.</Text>
        </View>
      </Pressable>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.surface,
    borderColor: colors.border,
    borderRadius: 8,
    borderWidth: 1,
    paddingHorizontal: 16,
    paddingVertical: 14,
  },
  title: {
    color: colors.primaryText,
    fontSize: 14,
    fontWeight: '800',
    marginBottom: 8,
  },
  row: {
    minHeight: 48,
    justifyContent: 'center',
  },
  deleteText: {
    color: colors.dangerStrong,
    fontSize: 14,
    fontWeight: '800',
  },
  description: {
    color: colors.mutedText,
    fontSize: 12,
    lineHeight: 17,
    marginTop: 2,
  },
});
