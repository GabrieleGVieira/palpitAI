const SITE_URL = import.meta.env.VITE_PUBLIC_SITE_URL ?? 'https://palpitai.com.br';

export const LEGAL_URLS = {
  accountDeletion: `${SITE_URL}/account-deletion`,
  privacy: `${SITE_URL}/privacy`,
  terms: `${SITE_URL}/terms`,
} as const;
