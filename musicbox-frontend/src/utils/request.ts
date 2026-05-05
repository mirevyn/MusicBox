import axios, { type AxiosError, type AxiosRequestConfig } from "axios";

// 基础 URL
const BASE_URL = import.meta.env.VITE_BACKEND_URL || "";
export type ApiErrorStatus = number | "timeout" | "network";

// 扩展 AxiosRequestConfig 类型，支持重试
declare module "axios" {
  export interface AxiosRequestConfig {
    retries?: number; // 最大重试次数
    retryDelay?: number; // 重试间隔(ms)
    retryCount?: number; // 当前已重试次数
  }
}

// 定义后端通用的响应结构
export interface ApiResponse<T = any> {
  code: number;
  msg: string;
  data?: T;
}

// 定义后端错误返回类型
interface ApiErrorData {
  code?: number;
  msg?: string;
  data?: any;
}

export function buildApiUrl(path: string) {
  return `${BASE_URL}/api${path}`.replace(/([^:])\/\/+?/g, "$1/");
}

export function getStoredAuthToken() {
  return localStorage.getItem("token");
}

export function extractApiErrorMessage(data: unknown, fallback = "请求失败") {
  const payload = data as ApiErrorData | undefined;
  return payload?.msg || fallback;
}

export function getApiErrorMessage(error: unknown, fallback = "请求失败") {
  if (axios.isAxiosError<ApiErrorData>(error)) {
    return extractApiErrorMessage(error.response?.data, fallback);
  }
  if (error instanceof Error && error.message) {
    return error.message;
  }
  return fallback;
}

export async function getApiBlobErrorMessage(
  error: unknown,
  fallback = "请求失败",
) {
  if (axios.isAxiosError(error) && error.response?.data instanceof Blob) {
    try {
      const text = await error.response.data.text();
      const payload = JSON.parse(text) as ApiErrorData;
      return payload?.msg || fallback;
    } catch {
      return fallback;
    }
  }
  return getApiErrorMessage(error, fallback);
}

export function dispatchAuthError(message: string) {
  localStorage.removeItem("token");
  window.dispatchEvent(
    new CustomEvent("auth-error", { detail: { message, status: 401 } }),
  );
}

export function dispatchApiError(message: string, status: ApiErrorStatus) {
  window.dispatchEvent(
    new CustomEvent("api-error", {
      detail: { message, status },
    }),
  );
}

// 全局错误处理
const handleErrorEvent = (error: AxiosError<unknown>) => {
  const message = extractApiErrorMessage(error.response?.data, "请求失败");
  const status = error.response?.status;

  // 优化日志显示：
  // 4xx (客户端错误/业务逻辑错误) 使用 warn，避免控制台满屏红字
  // 5xx (服务器错误) 或 网络错误 使用 error
  if (status && status >= 400 && status < 500) {
    console.warn(`API业务提示 [${status}]:`, message);
  } else {
    console.error("API请求系统错误:", error);
    console.log("详细错误信息:", message);
  }

  if (status === 401) {
    dispatchAuthError(message);
  } else if (error.response) {
    dispatchApiError(message, status ?? 500);
  } else if (error.code === "ECONNABORTED") {
    dispatchApiError("请求超时，请检查网络", "timeout");
  } else {
    dispatchApiError("网络连接错误，请检查网络", "network");
  }
};

// 创建 axios 实例
const service = axios.create({
  // 拼接 /api
  baseURL: buildApiUrl(""),
  timeout: 60000,
});

// 请求拦截器：自动携带 token
service.interceptors.request.use(
  (config) => {
    const token = getStoredAuthToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// 响应拦截器：统一处理错误 + 重试逻辑
service.interceptors.response.use(
  (response) => {
    return response.data;
  },
  async (error: AxiosError & { config?: AxiosRequestConfig }) => {
    const config = error.config;

    // 如果没有配置重试，直接处理错误
    if (!config || !config.retries) {
      handleErrorEvent(error);
      return Promise.reject(error);
    }

    // 初始化 retryCount
    config.retryCount ??= 0;

    // 超过最大重试次数，处理错误
    if (config.retryCount >= config.retries) {
      handleErrorEvent(error);
      return Promise.reject(error);
    }

    // 增加重试计数
    config.retryCount++;

    // 延迟重试
    await new Promise((resolve) =>
      setTimeout(resolve, config.retryDelay ?? 1000),
    );
    console.log(
      `正在重试请求 (${config.retryCount}/${config.retries}): ${config.url}`,
    );

    // 重试请求
    return service(config);
  },
);

export default service;
