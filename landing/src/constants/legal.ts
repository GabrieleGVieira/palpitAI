const SITE_URL = import.meta.env.VITE_PUBLIC_SITE_URL ?? 'https://palpitai.app';

export const LEGAL_URLS = {
  privacy: `${SITE_URL}/privacy`,
  terms: `${SITE_URL}/terms`,
} as const;
