<template>
    <nav
        class="fixed top-0 left-0 right-0 h-[calc(3.5rem+env(safe-area-inset-top))] md:h-[calc(4rem+env(safe-area-inset-top))] pt-[env(safe-area-inset-top)] bg-white/70 backdrop-blur-lg border-b border-gray-100 z-[100] grid grid-cols-[auto_1fr_auto] lg:grid-cols-[1fr_auto_1fr] items-center px-3 md:px-8 transition-all duration-300">

        <!-- 左侧：Logo -->
        <div class="flex items-center">
            <div class="flex items-center gap-2.5 cursor-pointer group/logo" @click="$router.push('/')">
                <img src="/logo.svg" alt="MusicBox Logo"
                    class="w-8 h-8 md:w-10 md:h-10 rounded-[0.6rem] md:rounded-xl shadow-lg shadow-primary/20 group-hover/logo:scale-105 transition-transform object-cover" />
                <h1 class="text-lg md:text-xl font-black tracking-tight text-gray-800 hidden sm:block">
                    MusicBox
                </h1>
            </div>
        </div>

        <!-- 中间：桌面端导航 (绝对居中) -->
        <div class="flex justify-center h-full">
            <div class="flex items-center gap-5 sm:gap-8 h-full">
                <a v-for="item in navItems" :key="item.name" href="#"
                    class="text-sm font-bold text-gray-500 lg:text-gray-400 hover:text-primary transition-all relative group py-1 h-full flex items-center"
                    :class="{ '!text-primary': activeItem === item.name }" @click.prevent="handleNavClick(item.path)">
                    {{ item.label }}
                    <span
                        class="absolute bottom-0 left-0 w-0 h-0.5 bg-primary transition-all duration-300 group-hover:w-full"
                        :class="{ '!w-full': activeItem === item.name }"></span>
                </a>
            </div>
        </div>

        <!-- 右侧：搜索与用户管理 -->
        <div class="flex items-center gap-3 md:gap-4 justify-end">
            <!-- 搜索框 -->
            <div class="hidden md:block w-full max-w-[240px] xl:max-w-[280px] transition-all">
                <SearchInput />
            </div>

            <!-- 用户区域 -->
            <div class="flex items-center gap-1.5 md:gap-2 shrink-0">
                <!-- 移动端搜索触发 (占位符) -->
                <button
                    class="md:hidden w-9 h-9 flex items-center justify-center text-gray-500 hover:text-primary rounded-full hover:bg-gray-100 transition-all"
                    @click="toggleMobileSearch">
                    <div class="i-fa6-solid-magnifying-glass text-sm"></div>
                </button>

                <!-- 分隔线 (桌面端) -->
                <div class="hidden md:block w-px h-6 bg-gray-100 mx-1"></div>

                <!-- 用户下拉菜单 -->
                <Dropdown v-if="userStore.user" placement="bottom-right" width="w-56">
                    <template #trigger>
                        <div
                            class="flex items-center gap-2 p-1 md:gap-2.5 md:p-1 md:pl-1 md:pr-2 rounded-full hover:bg-gray-50 border border-transparent hover:border-gray-100 transition-all group/user">
                            <div class="relative">
                                <img v-if="userStore.isLoggedIn && userStore.avatarUrl" :src="userStore.avatarUrl"
                                    class="
                                    w-8 h-8 rounded-full object-cover
                                    shadow-sm bg-gray-100
                                    transition-transform group-hover/user:scale-105
                                " loading="lazy" alt="User Avatar" />
                                <div v-else
                                    class="w-8 h-8 rounded-full bg-blue-50 flex items-center justify-center text-primary border border-blue-100 group-hover/user:scale-105 transition-all">
                                    <div class="i-fa6-solid-user text-xs"></div>
                                </div>
                            </div>
                            <span class="text-xs font-bold text-gray-700 hidden xl:block truncate max-w-[80px]">{{
                                userStore.user.username }}</span>
                            <div
                                class="i-fa6-solid-chevron-down text-[8px] text-gray-300 group-hover/user:text-primary transition-colors hidden md:block">
                            </div>
                        </div>
                    </template>
                    <div class="p-1.5">
                        <div class="px-4 py-3 bg-gray-50/50 rounded-lg mb-1.5 flex items-center justify-between gap-4">
                            <div class="text-[10px] font-bold text-gray-400 uppercase tracking-widest">当前账户</div>
                            <div class="text-sm font-black text-gray-800 truncate">{{ userStore.user.username }}</div>
                        </div>

                        <!-- 管理员入口 -->
                        <button v-if="userStore.isAdmin" @click="handleAdminClick"
                            class="w-full text-left px-4 py-2.5 text-sm font-bold text-blue-600 hover:bg-blue-50 rounded-lg transition-all flex items-center gap-3 active:scale-95 mb-1">
                            <div class="i-fa6-solid-gauge-high text-xs opacity-70"></div> 管理后台
                        </button>

                        <button @click="router.push('/profile')"
                            class="w-full text-left px-4 py-2.5 text-sm font-bold text-gray-600 hover:bg-blue-50 hover:text-primary rounded-lg transition-all flex items-center gap-3 active:scale-95">
                            <div class="i-fa6-solid-user text-xs opacity-50"></div> 个人中心
                        </button>
                        <button
                            class="w-full text-left px-4 py-2.5 text-sm font-bold text-rose-500 hover:bg-rose-50 rounded-lg transition-all flex items-center gap-3 active:scale-95"
                            @click="handleLogout">
                            <div class="i-fa6-solid-right-from-bracket text-xs opacity-50"></div> 退出登录
                        </button>
                    </div>
                </Dropdown>

                <!-- 未登录状态 -->
                <button v-else
                    class="px-5 md:px-6 h-8 md:h-10 bg-primary text-white text-xs md:text-sm font-black rounded-full shadow-lg shadow-primary/25 hover:shadow-primary/40 hover:-translate-y-0.5 active:scale-95 transition-all"
                    @click="router.push('/auth')">
                    立即登录
                </button>
            </div>
        </div>
    </nav>

    <!-- 移动端搜索面板 (根据缩小后的Navbar高度自适应top位置) -->
    <div v-if="showMobileSearch" class="fixed inset-0 z-[120] md:hidden">
        <div class="absolute inset-0 bg-black/20" @click="closeMobileSearch"></div>
        <div class="absolute top-[calc(3.5rem+env(safe-area-inset-top))] left-0 right-0 px-4 pb-4">
            <div
                class="bg-white/90 backdrop-blur-xl rounded-2xl shadow-xl border border-gray-100 p-3 flex items-center gap-2">
                <div class="flex-1">
                    <SearchInput />
                </div>
                <button
                    class="w-9 h-9 flex items-center justify-center rounded-full text-gray-500 hover:text-primary hover:bg-gray-100 transition-all"
                    @click="closeMobileSearch">
                    <div class="i-fa6-solid-xmark text-xs"></div>
                </button>
            </div>
        </div>
    </div>

    <!-- 退出登录确认弹窗 -->
    <ConfirmModal v-model="showLogoutConfirm" title="确认退出" content="确定要退出登录吗？" type="danger" confirm-text="确认退出"
        @confirm="confirmLogout" />
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, computed, ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { usePlayerStore } from '@/stores/player'
import Dropdown from '@/components/ui/Dropdown.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import SearchInput from '@/components/layout/SearchInput.vue'
import { showToast } from '@/utils/common'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const playerStore = usePlayerStore()

const showLogoutConfirm = ref(false)
const showMobileSearch = ref(false)

// 导航菜单配置
const navItems = [
    { name: 'home', label: '首页', path: '/' },
    { name: 'library', label: '我的', path: '/library' }
]

// 自动根据当前路由计算激活的菜单项
const activeItem = computed(() => {
    const currentPath = route.path
    if (currentPath === '/') return 'home'
    const found = navItems.find(item => item.path !== '/' && currentPath.startsWith(item.path))
    return found ? found.name : ''
})

// 点击导航跳转
const handleNavClick = (path: string) => {
    router.push(path)
}

const toggleMobileSearch = () => {
    showMobileSearch.value = !showMobileSearch.value
}

const closeMobileSearch = () => {
    showMobileSearch.value = false
}

const handleAdminClick = () => {
    if (playerStore.playing) {
        playerStore.playing = false
    }
    router.push('/admin')
}

const handleLogout = () => {
    showLogoutConfirm.value = true
}

const confirmLogout = () => {
    userStore.logout()
    showLogoutConfirm.value = false
    showToast('已退出登录', 'success')
    router.push('/auth')
}

// 监听全局 Auth 错误事件
const handleAuthError = () => {
    userStore.logout()
    showToast('登录已过期，请重新登录', 'warning')
    router.push('/auth')
}

onMounted(() => {
    userStore.fetchProfile()
    window.addEventListener('auth-error', handleAuthError)
})

onBeforeUnmount(() => {
    window.removeEventListener('auth-error', handleAuthError)
})

watch(() => route.fullPath, () => {
    showMobileSearch.value = false
})
</script>