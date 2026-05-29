import { StyleSheet, Text } from 'react-native';

import { LEGAL_URLS } from '../constants/legal';
import { colors } from '../theme';
import { openExternalUrl } from '../utils/openExternalUrl';

export function LegalConsentText() {
  return (
    <Text style={styles.text}>
      Ao continuar, você concorda com os{' '}
      <Text style={styles.link} onPress={() => openExternalUrl(LEGAL_URLS.terms)}>
        Termos de Uso
      </Text>{' '}
      e a{' '}
      <Text style={styles.link} onPress={() => openExternalUrl(LEGAL_URLS.privacy)}>
        Política de Privacidade
      </Text>
      .
    </Text>
  );
}

const styles = StyleSheet.create({
  text: {
    color: colors.mutedText,
    fontSize: 12,
    lineHeight: 18,
    marginTop: 18,
    textAlign: 'center',
  },
  link: {
    color: colors.primary,
    fontWeight: '800',
    textDecorationLine: 'underline',
  },
});
