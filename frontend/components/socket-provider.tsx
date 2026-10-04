"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";

import { useAuth } from "@/components/auth-provider";
import { WS_URL } from "@/lib/api";
import type { WsEvent } from "@/lib/types";

type Handler = (event: WsEvent) => void;

type SocketState = {
  connected: boolean;
  subscribe: (topic: string, handler: Handler) => () => void;
};

const SocketContext = createContext<SocketState | null>(null);

const RECONNECT_DELAY_MS = 2000;

export function SocketProvider({ children }: { children: React.ReactNode }) {
  const { token } = useAuth();
  const [connected, setConnected] = useState(false);

  // One socket for the whole tab. Components register per topic handlers and
  // the provider multiplexes over that single connection.
  const socketRef = useRef<WebSocket | null>(null);
  const handlersRef = useRef<Map<string, Set<Handler>>>(new Map());

  const send = useCallback((action: "subscribe" | "unsubscribe", topic: string) => {
    const socket = socketRef.current;
    if (socket?.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({ action, topic }));
    }
  }, []);

  useEffect(() => {
    if (!token) {
      socketRef.current?.close();
      socketRef.current = null;
      return;
    }

    let closedByUs = false;
    let retry: ReturnType<typeof setTimeout> | undefined;

    const connect = () => {
      const socket = new WebSocket(`${WS_URL}?token=${encodeURIComponent(token)}`);
      socketRef.current = socket;

      socket.onopen = () => {
        setConnected(true);
        // Re-subscribe after a reconnect so live updates resume on their own.
        for (const topic of handlersRef.current.keys()) {
          socket.send(JSON.stringify({ action: "subscribe", topic }));
        }
      };

      socket.onmessage = (message) => {
        let event: WsEvent;
        try {
          event = JSON.parse(message.data);
        } catch {
          return;
        }
        handlersRef.current.get(event.topic)?.forEach((handler) => handler(event));
      };

      socket.onclose = () => {
        setConnected(false);
        if (!closedByUs) retry = setTimeout(connect, RECONNECT_DELAY_MS);
      };
    };

    connect();

    return () => {
      closedByUs = true;
      if (retry) clearTimeout(retry);
      socketRef.current?.close();
      socketRef.current = null;
    };
  }, [token]);

  const subscribe = useCallback(
    (topic: string, handler: Handler) => {
      let handlers = handlersRef.current.get(topic);
      if (!handlers) {
        handlers = new Set();
        handlersRef.current.set(topic, handlers);
        send("subscribe", topic);
      }
      handlers.add(handler);

      return () => {
        const current = handlersRef.current.get(topic);
        if (!current) return;
        current.delete(handler);
        if (current.size === 0) {
          handlersRef.current.delete(topic);
          send("unsubscribe", topic);
        }
      };
    },
    [send],
  );

  const value = useMemo(() => ({ connected, subscribe }), [connected, subscribe]);

  return <SocketContext.Provider value={value}>{children}</SocketContext.Provider>;
}

export function useSocket() {
  const ctx = useContext(SocketContext);
  if (!ctx) throw new Error("useSocket must be used inside SocketProvider");
  return ctx;
}

// useTopic subscribes for as long as the component is mounted. The handler is
// held in a ref so an inline callback does not re-open the subscription.
export function useTopic(topic: string | null, handler: Handler) {
  const { subscribe } = useSocket();
  const handlerRef = useRef(handler);

  useEffect(() => {
    handlerRef.current = handler;
  }, [handler]);

  useEffect(() => {
    if (!topic) return;
    return subscribe(topic, (event) => handlerRef.current(event));
  }, [topic, subscribe]);
}
