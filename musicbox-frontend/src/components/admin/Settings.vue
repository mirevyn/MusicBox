<template>
    <div class="space-y-6">
        <div>
            <h2 class="text-2xl font-bold text-gray-800 dark:text-white mb-1">系统设置</h2>
            <p class="text-gray-500 text-sm">配置网站全局参数</p>
        </div>

        <!-- 安全设置 -->
        <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden relative">
            <div v-if="loading"
                class="absolute inset-0 bg-white/60 z-10 flex items-center justify-center backdrop-blur-[1px]">
                <div class="i-fa6-solid-spinner animate-spin text-primary text-xl"></div>
            </div>

            <div class="px-6 py-4 border-b border-gray-100 bg-gray-50/30">
                <h3 class="font-bold text-gray-800 flex items-center gap-2">
                    <div class="i-fa6-solid-shield-halved"></div> 安全与访问
                </h3>
            </div>
            <div class="p-6 space-y-4">
                <div class="flex items-center justify-between p-4 rounded-xl bg-gray-50">
                    <div>
                        <div class="font-bold text-gray-800 text-sm">开放注册</div>
                        <div class="text-xs text-gray-500">允许新用户注册账号</div>
                    </div>
                    <label class="relative inline-flex items-center cursor-pointer">
                        <input type="checkbox" v-model="settings.allowRegister" class="sr-only peer">
                        <div
                            class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/20 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary">
                        </div>
                    </label>
                </div>
            </div>
        </div>

        <div class="flex justify-end pt-4">
            <button @click="saveSettings" :disabled="saving"
                class="bg-primary text-white hover:bg-primary/90 disabled:opacity-50 transition-colors rounded-lg px-8 py-2.5 text-sm font-bold shadow-lg shadow-primary/20 flex items-center gap-2">
                <div v-if="saving" class="i-fa6-solid-spinner animate-spin"></div>
                <div v-else class="i-fa6-solid-floppy-disk"></div>
                {{ saving ? '正在保存...' : '保存更改' }}
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { fetchSystemSettings, updateSystemSettings } from '@/api/settings'
import { showToast } from '@/utils/common'

const loading = ref(false)
const saving = ref(false)

const settings = reactive({
    allowRegister: true
})

const loadSettings = async () => {
    loading.value = true
    try {
        const res = await fetchSystemSettings()
        if (res.code === 200 && res.data) {
            settings.allowRegister = res.data.allowRegister
        }
    } catch (error) {
        console.error('Failed to load settings:', error)
        showToast('加载设置失败', 'error')
    } finally {
        loading.value = false
    }
}

const saveSettings = async () => {
    saving.value = true
    try {
        const res = await updateSystemSettings(settings)
        if (res.code === 200) {
            showToast('设置已保存')
        } else {
            showToast(res.msg || '保存失败', 'error')
        }
    } catch (error) {
        console.error('Failed to save settings:', error)
        showToast('保存设置失败', 'error')
    } finally {
        saving.value = false
    }
}

onMounted(() => {
    loadSettings()
})
</script>
