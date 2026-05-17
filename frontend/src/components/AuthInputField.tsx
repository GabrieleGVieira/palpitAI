import { StyleSheet, Text, TextInput, View } from 'react-native';

type AuthInputFieldProps = {
  label: string;
  placeholder: string;
  secureTextEntry?: boolean;
  value: string;
  onChangeText: (text: string) => void;
  keyboardType?: 'default' | 'email-address';
  autoCapitalize?: 'none' | 'sentences' | 'words' | 'characters';
};

export function AuthInputField({
  label,
  placeholder,
  secureTextEntry = false,
  value,
  onChangeText,
  keyboardType = 'default',
  autoCapitalize = 'none',
}: AuthInputFieldProps) {
  return (
    <View style={styles.field}>
      <Text style={styles.label}>{label}</Text>
      <TextInput
        autoCapitalize={autoCapitalize}
        keyboardType={keyboardType}
        onChangeText={onChangeText}
        placeholder={placeholder}
        placeholderTextColor="#9ca6a0"
        secureTextEntry={secureTextEntry}
        selectionColor="#1c4c34"
        style={styles.input}
        value={value}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  field: {
    marginBottom: 20,
  },
  label: {
    color: '#1c4c34',
    fontSize: 14,
    fontWeight: '700',
    letterSpacing: 0.2,
    marginBottom: 8,
  },
  input: {
    backgroundColor: '#f7faf6',
    borderColor: '#c8d8cc',
    borderRadius: 20,
    borderWidth: 1,
    color: '#0e382a',
    fontSize: 16,
    minHeight: 50,
    paddingHorizontal: 18,
    paddingVertical: 14,
  },
});
