import { computed, ref } from "vue";
import { defineStore } from "pinia";
import { showToast } from "@/utils/common";

export interface AdminNotification {
  id: string;
  type: string;
  title: string;
  content: string;
  targetId?: number;
  route?: string;
  createdAt: string;
  read: boolean;
}

const HEARTBEAT_MS = 25000;
const RECONNECT_MS = 3000;

export const useNotificationStore = defineStore("notification", () => {
  const socket = ref<WebSocket | null>(null);
  const notifications = ref<AdminNotification[]>([]);
  const connected = ref(false);
  const panelOpen = ref(false);

  let heartbeatTimer: number | null = null;
  let reconnectTimer: number | null = null;
  let shouldReconnect = false;

  // ---- 事件回调：dashboard_refresh 信号 ----
  const dashboardRefreshCallbacks = new Set<() => void>();

  const unreadCount = computed(
    () => notifications.value.filter((item) => !item.read).length
  );

  const connect = () => {
    if (typeof window === "undefined") return;
    if (
      socket.value &&
      (socket.value.readyState === WebSocket.OPEN ||
        socket.value.readyState === WebSocket.CONNECTING)
    ) {
      return;
    }

    const token = localStorage.getItem("token");
    if (!token) return;

    const url = buildSocketUrl();
    if (!url) return;

    shouldReconnect = true;
    clearReconnect();

    const ws = new WebSocket(url, ["musicbox.jwt", token]);
    socket.value = ws;

    ws.onopen = () => {
      connected.value = true;
      startHeartbeat();
    };

    ws.onmessage = (event) => {
      try {
        const payload = JSON.parse(event.data) as Omit<AdminNotification, "read">;

        // dashboard_refresh 是静默信号，不进通知列表，直接触发回调
        if (payload.type === "dashboard_refresh") {
          dashboardRefreshCallbacks.forEach((cb) => cb());
          return;
        }

        // 按 id 去重
        if (notifications.value.some((n) => n.id === payload.id)) return;
        notifications.value.unshift({ ...payload, read: false });
        if (notifications.value.length > 20) {
          notifications.value = notifications.value.slice(0, 20);
        }
        showToast(payload.title, "info");
      } catch (error) {
        console.warn("通知解析失败", error);
      }
    };

    ws.onclose = () => {
      connected.value = false;
      stopHeartbeat();
      socket.value = null;
      if (shouldReconnect) {
        reconnectTimer = window.setTimeout(() => connect(), RECONNECT_MS);
      }
    };

    ws.onerror = () => {
      connected.value = false;
    };
  };

  const disconnect = () => {
    shouldReconnect = false;
    clearReconnect();
    stopHeartbeat();
    connected.value = false;
    socket.value?.close();
    socket.value = null;
  };

  const togglePanel = () => {
    panelOpen.value = !panelOpen.value;
  };

  const closePanel = () => {
    panelOpen.value = false;
  };

  const markAllRead = () => {
    notifications.value = notifications.value.map((item) => ({
      ...item,
      read: true,
    }));
  };

  const markRead = (id: string) => {
    notifications.value = notifications.value.map((item) =>
      item.id === id ? { ...item, read: true } : item
    );
  };

  const startHeartbeat = () => {
    stopHeartbeat();
    heartbeatTimer = window.setInterval(() => {
      if (socket.value?.readyState === WebSocket.OPEN) {
        socket.value.send("ping");
      }
    }, HEARTBEAT_MS);
  };

  const stopHeartbeat = () => {
    if (heartbeatTimer) {
      window.clearInterval(heartbeatTimer);
      heartbeatTimer = null;
    }
  };

  const clearReconnect = () => {
    if (reconnectTimer) {
      window.clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }
  };

  const onDashboardRefresh = (cb: () => void) => {
    dashboardRefreshCallbacks.add(cb);
  };
  const offDashboardRefresh = (cb: () => void) => {
    dashboardRefreshCallbacks.delete(cb);
  };

  return {
    notifications,
    connected,
    panelOpen,
    unreadCount,
    connect,
    disconnect,
    togglePanel,
    closePanel,
    markAllRead,
    markRead,
    onDashboardRefresh,
    offDashboardRefresh,
  };
});

const buildSocketUrl = () => {
  const baseUrl = (import.meta.env.VITE_BACKEND_URL || "").trim();

  if (baseUrl) {
    const normalized = baseUrl.replace(/\/+$/, "");
    const wsBase = normalized.replace(/^http/i, "ws");
    return `${wsBase}/ws/admin/notifications`;
  }

  if (typeof window === "undefined") return "";
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  return `${protocol}//${window.location.host}/ws/admin/notifications`;
};
