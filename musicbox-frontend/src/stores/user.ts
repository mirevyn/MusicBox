import { isAxiosError } from "axios";
import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { getUserInfo } from "@/api/user";
import { resolveUrl } from "@/utils/common";

import type { User } from "@/types/api";
import { usePlayerStore } from "./player";

export const useUserStore = defineStore("user", () => {
    // 状态定义数据区域
    const user = ref<User | null>(null);
    const loading = ref(false);

    // 计算属性定义区域
    const isLoggedIn = computed(() => !!user.value && !!localStorage.getItem("token"));

    const avatarUrl = computed(() => {
        if (!user.value) return "";
        const url = user.value.avatarUrl || user.value.avatar_url;
        return resolveUrl(url);
    });


    const isAdmin = computed(() => user.value?.role?.toLowerCase() === "admin");

    function clearAuthState() {
        localStorage.removeItem("token");
        user.value = null;
    }

    // 动作方法逻辑区域
    async function fetchProfile() {
        const token = localStorage.getItem("token");
        if (!token) {
            user.value = null;
            return;
        }

        loading.value = true;
        try {
            const res = await getUserInfo();
            if ([401, 403].includes(res.code)) {
                clearAuthState();
                return;
            }

            // 确保 data 存在且有效
            if (res.code === 200 && res.data && Object.keys(res.data).length > 0) {
                user.value = res.data;
            } else {
                // 避免把临时异常误判成登录失效
                console.warn("用户信息返回异常，保留登录态，等待后续重试", {
                    code: res.code,
                    msg: res.msg,
                });
                user.value = null;
            }
        } catch (error) {
            console.error("获取用户信息失败", error);
            if (isAxiosError(error) && [401, 403].includes(error.response?.status ?? 0)) {
                clearAuthState();
            }
        } finally {
            loading.value = false;
        }
    }

    function setUser(newUser: User | null) {
        user.value = newUser;
    }

    function logout() {
        clearAuthState();

        // 彻底清理：重置播放器状态
        try {
            const playerStore = usePlayerStore();
            playerStore.reset();
        } catch (e) {
            console.error("重置播放器失败:", e);
        }
    }

    return {
        // 状态属性区域
        user,
        loading,
        // 计算属性区域
        isLoggedIn,
        avatarUrl,
        isAdmin,
        // 动作方法区域
        fetchProfile,
        setUser,
        logout,
    };
});
