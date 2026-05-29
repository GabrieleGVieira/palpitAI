import { Pressable, StyleSheet, Text, View } from 'react-native';

import { LEGAL_URLS } from '../constants/legal';
import { colors } from '../theme';
import { openExternalUrl } from '../utils/openExternalUrl';

export function LegalLinksCard() {
  return (
    <View style={styles.card}>
      <Text style={styles.title}>Informações legais</Text>
      <Pressable onPress={() => openExternalUrl(LEGAL_URLS.terms)} style={styles.row}>
        <Text style={styles.linkText}>Termos de Uso</Text>
      </Pressable>
      <Pressable onPress={() => openExternalUrl(LEGAL_URLS.privacy)} style={styles.row}>
        <Text style={styles.linkText}>Política de Privacidade</Text>
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
    minHeight: 38,
    justifyContent: 'center',
  },
  linkText: {
    color: colors.primary,
    fontSize: 14,
    fontWeight: '700',
    textDecorationLine: 'underline',
  },
});
