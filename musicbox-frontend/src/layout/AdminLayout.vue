<template>
    <div class="h-screen bg-[#f5f7fa] overflow-hidden flex font-sans">
        <!-- 侧边栏 -->
        <AdminSidebar :is-open="sidebarOpen" @close="sidebarOpen = false" />

        <!-- 主要内容区域 -->
        <div class="flex-1 flex flex-col lg:ml-72 h-screen min-h-0 overflow-hidden transition-all duration-300 ease-in-out">
            <!-- 顶部导航栏 -->
            <AdminHeader @toggle-sidebar="sidebarOpen = !sidebarOpen" />

            <!-- 业务内容展示区 -->
            <main class="flex-1 min-h-0 overflow-y-auto overflow-x-hidden px-4 py-4 md:px-6 md:py-5 lg:px-8 lg:py-6">
                <div class="max-w-7xl mx-auto w-full">
                    <RouterView v-slot="{ Component }">
                        <transition name="fade-slide" mode="out-in">
                            <component :is="Component" />
                        </transition>
                    </RouterView>
                </div>
            </main>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import AdminSidebar from '@/components/admin/AdminSidebar.vue'
import AdminHeader from '@/components/admin/AdminHeader.vue'
import { useUserStore } from '@/stores/user'
import { useNotificationStore } from '@/stores/notification'

const sidebarOpen = ref(false)
const userStore = useUserStore()
const notificationStore = useNotificationStore()

const syncNotificationConnection = () => {
    if (userStore.isLoggedIn && userStore.isAdmin) {
        notificationStore.connect()
        return
    }
    notificationStore.disconnect()
}

onMounted(() => {
    if (localStorage.getItem('token') && !userStore.user) {
        void userStore.fetchProfile().finally(() => {
            syncNotificationConnection()
        })
        return
    }
    syncNotificationConnection()
})

watch(() => [userStore.isLoggedIn, userStore.isAdmin], () => {
    syncNotificationConnection()
})

onBeforeUnmount(() => {
    notificationStore.disconnect()
})
</script>

<style scoped>
.fade-slide-enter-active,
.fade-slide-leave-active {
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-slide-enter-from {
    opacity: 0;
    transform: translateY(10px);
}

.fade-slide-leave-to {
    opacity: 0;
    transform: translateY(-10px);
}
</style>
