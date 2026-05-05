<template>
    <header class="h-20 bg-white/80 backdrop-blur-xl border-b border-gray-100/50 flex items-center justify-between px-4 md:px-8 sticky top-0 z-40 transition-all duration-300">
        <button @click="$emit('toggleSidebar')" class="lg:hidden p-2 -ml-2 text-gray-500 hover:text-primary transition-colors rounded-lg hover:bg-gray-50">
            <div class="i-fa6-solid-bars text-xl"></div>
        </button>

        <div class="flex items-center gap-4 ml-auto">
            <div class="relative" ref="notificationArea">
                <button
                    @click="notificationStore.togglePanel()"
                    class="relative w-10 h-10 rounded-full flex items-center justify-center text-gray-500 hover:bg-gray-50 hover:text-primary transition-all"
                >
                    <div class="i-fa6-solid-bell text-lg"></div>
                    <span
                        v-if="notificationStore.unreadCount > 0"
                        class="absolute right-1.5 top-1.5 min-w-4 h-4 rounded-full bg-red-500 px-1 text-[10px] font-bold text-white flex items-center justify-center"
                    >
                        {{ notificationStore.unreadCount > 99 ? '99+' : notificationStore.unreadCount }}
                    </span>
                    <span
                        class="absolute left-1.5 bottom-1.5 h-2 w-2 rounded-full"
                        :class="notificationStore.connected ? 'bg-emerald-500' : 'bg-gray-300'"
                    ></span>
                </button>

                <div
                    v-if="notificationStore.panelOpen"
                    class="absolute right-0 top-12 z-50 w-[360px] overflow-hidden rounded-2xl border border-gray-100 bg-white shadow-2xl shadow-slate-200/60"
                >
                    <div class="flex items-center justify-between border-b border-gray-100 px-4 py-3">
                        <div>
                            <div class="text-sm font-bold text-gray-800">实时通知</div>
                            <div class="text-[11px] text-gray-400">
                                {{ notificationStore.connected ? 'WebSocket 已连接' : '连接中断' }}
                            </div>
                        </div>
                        <button
                            @click="notificationStore.markAllRead()"
                            class="text-xs font-medium text-primary hover:text-primary/80 transition-colors"
                        >
                            全部已读
                        </button>
                    </div>

                    <div v-if="notificationStore.notifications.length === 0" class="px-4 py-10 text-center">
                        <div class="i-fa6-solid-bell-slash mx-auto mb-3 text-2xl text-gray-200"></div>
                        <div class="text-sm font-medium text-gray-500">暂无实时通知</div>
                        <div class="mt-1 text-xs text-gray-400">有新的待审核歌单时会立即出现在这里</div>
                    </div>

                    <div v-else class="max-h-[360px] overflow-y-auto">
                        <button
                            v-for="item in notificationStore.notifications"
                            :key="item.id"
                            @click="openNotification(item)"
                            class="w-full border-b border-gray-100 px-4 py-3 text-left transition-colors hover:bg-gray-50 last:border-b-0"
                        >
                            <div class="flex items-start gap-3">
                                <span
                                    class="mt-1 h-2.5 w-2.5 shrink-0 rounded-full"
                                    :class="item.read ? 'bg-gray-200' : 'bg-primary'"
                                ></span>
                                <div class="min-w-0 flex-1">
                                    <div class="flex items-center justify-between gap-3">
                                        <div class="truncate text-sm font-semibold text-gray-800">{{ item.title }}</div>
                                        <div class="shrink-0 text-[11px] text-gray-400">{{ formatRelativeTime(item.createdAt) }}</div>
                                    </div>
                                    <div class="mt-1 text-xs leading-5 text-gray-500">{{ item.content }}</div>
                                </div>
                            </div>
                        </button>
                    </div>
                </div>
            </div>
            
            <div class="h-8 w-px bg-gray-200"></div>

            <div class="flex items-center gap-3 pl-2 py-1.5 rounded-full hover:bg-gray-50 transition-colors cursor-pointer pr-2 md:pr-4 border border-transparent hover:border-gray-100">
                <div class="w-9 h-9 rounded-full bg-gradient-to-br from-primary to-blue-600 flex items-center justify-center text-white font-bold shadow-lg shadow-primary/20 ring-2 ring-white">
                    {{ (userStore.user?.username || 'A').charAt(0).toUpperCase() }}
                </div>
                <div class="hidden sm:flex flex-col">
                    <span class="text-sm font-bold text-gray-700 leading-none mb-1">{{ userStore.user?.username || 'Admin' }}</span>
                    <span class="text-xs text-gray-400 font-medium leading-none">Super Admin</span>
                </div>
            </div>
        </div>
    </header>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useNotificationStore, type AdminNotification } from '@/stores/notification'

const userStore = useUserStore()
const notificationStore = useNotificationStore()
const router = useRouter()

defineEmits<{
    (e: 'toggleSidebar'): void
}>()

// ---- 点击面板外部自动关闭 ----
const notificationArea = ref<HTMLElement | null>(null)

const onDocumentClick = (e: MouseEvent) => {
    if (!notificationStore.panelOpen) return
    if (notificationArea.value && !notificationArea.value.contains(e.target as Node)) {
        notificationStore.closePanel()
    }
}

onMounted(() => {
    document.addEventListener('click', onDocumentClick, true)
})
onBeforeUnmount(() => {
    document.removeEventListener('click', onDocumentClick, true)
})

const openNotification = (item: AdminNotification) => {
    notificationStore.markRead(item.id)
    notificationStore.closePanel()
    if (item.route) {
        router.push(item.route)
    }
}

const formatRelativeTime = (value: string) => {
    if (!value) return ''
    const target = new Date(value).getTime()
    const diff = Date.now() - target
    const minute = 60 * 1000
    const hour = 60 * minute
    if (diff < hour) {
        return `${Math.max(1, Math.floor(diff / minute))} 分钟前`
    }
    if (diff < 24 * hour) {
        return `${Math.max(1, Math.floor(diff / hour))} 小时前`
    }
    return new Date(value).toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}
</script>
