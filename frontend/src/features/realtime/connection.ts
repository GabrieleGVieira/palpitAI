import { apiURL } from '../../shared/services/apiClient';
import { supabase } from '../../services/supabase';
import type { RealtimeOptions } from './types';
import { websocketURL } from './url';

const reconnectDelayMs = 2000;

export async function connectRealtime({ groupID, onEvent }: RealtimeOptions) {
  const {
    data: { session },
  } = await supabase.auth.getSession();

  if (!session?.access_token) {
    throw new Error('Sua sessao expirou. Entre novamente para receber atualizacoes ao vivo.');
  }

  const accessToken = session.access_token;
  let isClosed = false;
  let socket: WebSocket | null = null;
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null;

  function connect() {
    const url = new URL('/ws', websocketURL(apiURL));
    url.searchParams.set('token', accessToken);
    if (groupID) {
      url.searchParams.set('group_id', groupID);
    }

    socket = new WebSocket(url.toString());

    socket.onopen = () => {
      console.info('Realtime conectado.', { groupID });
    };

    socket.onmessage = (message) => {
      try {
        const event = JSON.parse(String(message.data));
        console.info('Realtime recebido.', {
          matchID: event?.payload?.match_id,
          name: event?.name,
          room: event?.room,
        });
        onEvent(event);
      } catch {
        // REST remains the source of truth if a realtime message is malformed.
      }
    };

    socket.onclose = (event) => {
      console.warn('Realtime desconectado.', { code: event.code, groupID, reason: event.reason });
      if (!isClosed) {
        reconnectTimer = setTimeout(connect, reconnectDelayMs);
      }
    };

    socket.onerror = () => {
      console.warn('Realtime encontrou um erro de conexão.', { groupID });
      socket?.close();
    };
  }

  connect();

  return () => {
    isClosed = true;
    if (reconnectTimer) {
      clearTimeout(reconnectTimer);
    }
    socket?.close();
  };
}

export { websocketURL };
