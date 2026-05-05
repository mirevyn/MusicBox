import request, {
  buildApiUrl,
  dispatchApiError,
  dispatchAuthError,
  extractApiErrorMessage,
  getStoredAuthToken,
} from "../utils/request";

export interface AIChatMessage {
  role: "user" | "assistant";
  content: string;
}

export interface AIStatusResponse {
  online: boolean;
  model: string;
  providerKey?: "ollama" | "openai" | "auto";
  provider?: string;
  models?: string[];
  baseUrl?: string;
}

// 获取本地 AI 服务状态
export function getAIStatus(baseUrl?: string, model?: string, apiKey?: string, provider?: string) {
  return request.post("/ai/status", { baseUrl, model, apiKey, provider });
}


// 调用本地 Ollama 流式聊天接口
export async function streamAIChat(
  messages: AIChatMessage[],
  baseUrl?: string,
  model?: string,
  apiKey?: string,
  provider?: string,
  signal?: AbortSignal,
) {
  try {
    const token = getStoredAuthToken();
    const response = await fetch(buildApiUrl("/ai/chat/stream"), {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
      body: JSON.stringify({ messages, baseUrl, model, apiKey, provider }),
      signal,
    });

    const contentType = response.headers.get("content-type") || "";
    if (contentType.includes("application/json")) {
      const result = await response.json().catch(() => ({}));
      const message = extractApiErrorMessage(result, "AI stream unavailable");
      if (response.status === 401) dispatchAuthError(message);
      else dispatchApiError(message, response.status);
      throw new Error(message);
    }

    if (!response.ok || !response.body) {
      const text = await response.text().catch(() => "");
      const compactText = text.trim().replace(/\s+/g, " ").slice(0, 120);
      const message = compactText || `AI stream unavailable (HTTP ${response.status})`;
      if (response.status === 401) dispatchAuthError(message);
      else dispatchApiError(message, response.status);
      throw new Error(message);
    }

    return response;
  } catch (error: any) {
    if (error?.name === "AbortError") {
      throw error;
    }

    if (error instanceof Error) {
      if (error.message === "Failed to fetch") {
        dispatchApiError("网络连接错误，请检查网络", "network");
        throw new Error("网络连接错误，请检查网络");
      }
      throw error;
    }

    dispatchApiError("AI stream unavailable", "network");
    throw new Error("AI stream unavailable");
  }
}
