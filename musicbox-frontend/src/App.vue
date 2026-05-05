<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import Toast from '@/components/ui/Toast.vue'

const toastRef = ref<InstanceType<typeof Toast> | null>(null)

const handleGlobalToast = (event: CustomEvent) => {
    if (event.detail?.message) {
        toastRef.value?.show(event.detail.message, event.detail?.type)
    }
}

onMounted(() => {
    window.addEventListener('play-toast', handleGlobalToast as EventListener)
    window.addEventListener('add-playlist-toast', handleGlobalToast as EventListener)
})

onBeforeUnmount(() => {
    window.removeEventListener('play-toast', handleGlobalToast as EventListener)
    window.removeEventListener('add-playlist-toast', handleGlobalToast as EventListener)
})
</script>

<template>
  <Teleport to="body">
    <Toast ref="toastRef" />
  </Teleport>
  <router-view />
</template>

