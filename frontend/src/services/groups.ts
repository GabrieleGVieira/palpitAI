import { supabase } from './supabase';

const apiURL = process.env.EXPO_PUBLIC_API_URL ?? 'http://localhost:3000';

export type CreateGroupPayload = {
  description: string;
  has_unlimited_participants: boolean;
  is_private: boolean;
  match_scope: 'all' | 'selected';
  name: string;
  participant_limit: number | null;
  selected_teams: string[];
};

export type Group = {
  created_at: string;
  description: string;
  id: string;
  invite_code: string;
  is_private: boolean;
  match_scope: string;
  name: string;
  owner_id: string;
  participant_limit: number | null;
  selected_teams: string[];
};

type APIError = {
  error?: string;
};

async function readJSONResponse(response: Response) {
  const responseText = await response.text();

  if (!responseText) {
    return {};
  }

  try {
    return JSON.parse(responseText) as Group | APIError;
  } catch {
    throw new Error(`A API respondeu em formato inesperado: ${responseText}`);
  }
}

export async function createGroup(payload: CreateGroupPayload) {
  const {
    data: { session },
  } = await supabase.auth.getSession();

  if (!session?.access_token) {
    throw new Error('Sua sessao expirou. Entre novamente para criar um grupo.');
  }

  let response: Response;

  try {
    response = await fetch(`${apiURL}/api/v1/groups`, {
      body: JSON.stringify(payload),
      headers: {
        Authorization: `Bearer ${session.access_token}`,
        'Content-Type': 'application/json',
      },
      method: 'POST',
    });
  } catch {
    throw new Error(
      `Nao foi possivel acessar a API em ${apiURL}. Verifique se o backend esta rodando e se o celular esta na mesma rede.`,
    );
  }

  const data = await readJSONResponse(response);

  if (!response.ok) {
    throw new Error('error' in data && data.error ? data.error : 'Nao foi possivel criar o grupo.');
  }

  return data as Group;
}
