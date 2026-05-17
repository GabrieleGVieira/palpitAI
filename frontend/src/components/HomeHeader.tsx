import { Image, Pressable, StyleSheet, Text, View } from 'react-native';

type HomeHeaderProps = {
  userName?: string;
  onCreateGroup: () => void;
  onLogout: () => void;
  isSubmitting: boolean;
};

export function HomeHeader({ userName, onCreateGroup, onLogout, isSubmitting }: HomeHeaderProps) {
  return (
    <View style={styles.header}>
      <View style={styles.logoMark}>
        <Image
          accessibilityIgnoresInvertColors
          resizeMode="cover"
          source={require('../../assets/splash-palpitai.png')}
          style={styles.logoImage}
        />
      </View>
      <Text style={styles.title}>Olá, {userName || 'amigo'}</Text>
      <Text style={styles.subtitle}>
        Aqui você acompanha seus grupos e seus palpites com calma.
      </Text>

      <View style={styles.actionTabs}>
        <Pressable onPress={onCreateGroup} style={[styles.tabButton, styles.tabPrimary]}>
          <Text style={styles.tabButtonText}>Criar grupo</Text>
        </Pressable>
        <Pressable
          disabled={isSubmitting}
          onPress={onLogout}
          style={[styles.tabButton, styles.tabSecondary, isSubmitting && styles.buttonDisabled]}>
          <Text style={styles.tabSecondaryText}>{isSubmitting ? 'Saindo...' : 'Sair'}</Text>
        </Pressable>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  header: {
    paddingTop: 32,
  },
  logoMark: {
    alignItems: 'center',
    backgroundColor: '#ffffff',
    borderColor: '#d9e7d4',
    borderRadius: 32,
    borderWidth: 1,
    height: 64,
    justifyContent: 'center',
    marginBottom: 24,
    overflow: 'hidden',
    shadowColor: '#1e5c39',
    shadowOffset: { height: 8, width: 0 },
    shadowOpacity: 0.12,
    shadowRadius: 16,
    width: 64,
  },
  logoImage: {
    height: 76,
    transform: [{ scale: 1.18 }],
    width: 76,
  },
  title: {
    color: '#123d2a',
    fontSize: 38,
    fontWeight: '800',
    letterSpacing: 0,
  },
  subtitle: {
    color: '#486654',
    fontSize: 16,
    lineHeight: 24,
    marginTop: 12,
    maxWidth: 340,
  },
  actionTabs: {
    flexDirection: 'row',
    gap: 10,
    marginTop: 18,
  },
  tabButton: {
    flex: 1,
    alignItems: 'center',
    borderRadius: 999,
    justifyContent: 'center',
    minHeight: 48,
    paddingHorizontal: 12,
  },
  tabPrimary: {
    backgroundColor: '#1f7a4a',
  },
  tabSecondary: {
    backgroundColor: '#ffffff',
    borderColor: '#1f7a4a',
    borderWidth: 1,
  },
  tabButtonText: {
    color: '#ffffff',
    fontSize: 14,
    fontWeight: '800',
  },
  tabSecondaryText: {
    color: '#1f7a4a',
    fontSize: 14,
    fontWeight: '800',
  },
  buttonDisabled: {
    opacity: 0.72,
  },
});
