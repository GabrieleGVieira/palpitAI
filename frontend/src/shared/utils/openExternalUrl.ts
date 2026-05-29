import { Alert, Linking } from 'react-native';

export async function openExternalUrl(url: string) {
  try {
    await Linking.openURL(url);
  } catch {
    Alert.alert('Não foi possível abrir o link', 'Tente novamente em alguns instantes.');
  }
}
