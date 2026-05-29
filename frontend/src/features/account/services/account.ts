import { apiClient } from '../../../shared/services/apiClient';

export async function deleteAccount() {
  await apiClient<Record<string, string>>('/api/v1/me', {
    fallbackError: 'Não foi possível excluir sua conta agora.',
    method: 'DELETE',
  });
}
