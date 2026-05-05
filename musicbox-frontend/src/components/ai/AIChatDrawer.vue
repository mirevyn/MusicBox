<template>
    <div class="z-[200] relative">
        <Transition enter-active-class="transition-opacity duration-200 ease-linear" enter-from-class="opacity-0"
            enter-to-class="opacity-100" leave-active-class="transition-opacity duration-200 ease-linear"
            leave-from-class="opacity-100" leave-to-class="opacity-0">
            <div v-if="isOpen" @click="close" class="fixed inset-0 bg-black/20 z-[200]"></div>
        </Transition>

        <!-- 侧边栏面板 -->
        <Transition enter-active-class="transition-transform duration-200 ease-out" enter-from-class="translate-x-full"
            enter-to-class="translate-x-0" leave-active-class="transition-transform duration-150 ease-in"
            leave-from-class="translate-x-0" leave-to-class="translate-x-full">
            <div v-if="isOpen"
                class="fixed top-0 right-0 h-full w-full md:w-[450px] bg-white z-[201] shadow-2xl flex flex-col">
                <!-- 头部区域 -->
                <div
                    class="flex items-center justify-between p-4 border-b border-gray-100 bg-white/50 backdrop-blur-xl absolute top-0 left-0 right-0 z-10 transition-colors">
                    <div class="flex items-center gap-3">
                        <div
                            class="w-10 h-10 rounded-full bg-gradient-to-br from-primary to-blue-600 flex items-center justify-center text-white shadow-lg shadow-blue-500/30">
                            <div class="i-fa6-solid-robot text-lg"></div>
                        </div>
                        <div>
                            <h2 class="font-bold text-gray-800 text-lg">AI Music Guide</h2>
                            <div class="flex items-center gap-1.5 text-xs font-medium"
                                :class="aiOnline ? 'text-green-500' : 'text-amber-500'">
                                <div class="w-1.5 h-1.5 rounded-full"
                                    :class="aiOnline ? 'bg-green-500 animate-pulse' : 'bg-amber-500'"></div>
                                <span>{{ aiOnline ? `${providerName} · ${modelName} Online` : `${providerName} ·
                                    ${modelName} Offline` }}</span>
                            </div>
                        </div>
                    </div>
                    <div class="flex items-center gap-1">
                        <button @click="showSettings = true"
                            class="w-8 h-8 rounded-full hover:bg-gray-100 flex items-center justify-center text-gray-400 hover:text-gray-600 transition-colors">
                            <div class="i-fa6-solid-gear text-sm"></div>
                        </button>
                        <button @click="close"
                            class="w-8 h-8 rounded-full hover:bg-gray-100 flex items-center justify-center text-gray-400 hover:text-gray-600 transition-colors">
                            <div class="i-fa6-solid-xmark text-lg"></div>
                        </button>
                    </div>
                </div>

                <!-- AI 设置弹窗 -->
                <Transition enter-active-class="transition-all duration-200 ease-out"
                    enter-from-class="opacity-0 translate-y-2" enter-to-class="opacity-100 translate-y-0"
                    leave-active-class="transition-all duration-150 ease-in"
                    leave-from-class="opacity-100 translate-y-0" leave-to-class="opacity-0 translate-y-2">
                    <div v-if="showSettings"
                        class="absolute inset-0 z-[202] bg-white flex flex-col p-6 overflow-y-auto custom-scrollbar">
                        <div class="flex items-center justify-between mb-8">
                            <div>
                                <h3 class="font-black text-gray-900 text-xl tracking-tight">AI 配置中心</h3>
                                <p class="text-xs text-gray-400 mt-1">定制你的专属智能助手参数</p>
                            </div>
                            <button @click="showSettings = false"
                                class="w-10 h-10 rounded-full hover:bg-gray-200/50 flex items-center justify-center text-gray-400 hover:text-gray-600 transition-all active:scale-90">
                                <div class="i-fa6-solid-xmark text-lg"></div>
                            </button>
                        </div>

                        <div class="space-y-6 flex-1">
                            <!-- 服务配置卡片 -->
                            <div class="bg-white rounded-[1.5rem] p-5 shadow-sm border border-gray-100/50">
                                <div class="flex items-center gap-2 mb-4 text-gray-800">
                                    <div class="i-fa6-solid-server text-sm text-primary"></div>
                                    <h4 class="text-sm font-bold">服务连接</h4>
                                </div>
                                <div class="space-y-4">
                                    <div class="rounded-2xl bg-slate-50 px-4 py-3 border border-gray-100">
                                        <div class="flex items-center justify-between gap-3">
                                            <div>
                                                <div
                                                    class="text-[11px] font-bold uppercase tracking-wider text-gray-400">
                                                    当前生效连接</div>
                                                <div class="mt-1 text-sm font-semibold text-gray-800">{{
                                                    connectionSourceLabel }}</div>
                                            </div>
                                            <div class="text-[10px] font-bold uppercase tracking-wide"
                                                :class="aiOnline ? 'text-emerald-600' : 'text-amber-500'">
                                                {{ aiOnline ? 'online' : 'offline' }}
                                            </div>
                                        </div>
                                        <div class="mt-2 text-[11px] text-gray-500 break-all">{{ statusBaseUrl ||
                                            '跟随后端环境变量' }}</div>
                                        <div class="mt-1 text-[11px] text-gray-400">{{ providerName }} · {{ modelName }}
                                        </div>
                                        <div v-if="statusHint" class="mt-2 text-[11px]"
                                            :class="aiOnline ? 'text-emerald-600' : 'text-amber-600'">
                                            {{ statusHint }}
                                        </div>
                                    </div>
                                    <div>
                                        <label
                                            class="block text-[11px] font-bold text-gray-400 uppercase tracking-wider mb-1.5 ml-1">服务类型</label>
                                        <div class="grid grid-cols-3 gap-2">
                                            <button v-for="option in providerOptions" :key="option.value"
                                                @click="settingsProvider = option.value"
                                                class="rounded-xl px-3 py-2 text-xs font-bold transition-all" :class="settingsProvider === option.value
                                                    ? 'bg-primary text-white shadow-md shadow-blue-500/20'
                                                    : 'bg-slate-50 text-gray-500 hover:bg-blue-50 hover:text-primary'">
                                                {{ option.label }}
                                            </button>
                                        </div>
                                    </div>
                                    <div>
                                        <label
                                            class="block text-[11px] font-bold text-gray-400 uppercase tracking-wider mb-1.5 ml-1">AI
                                            地址 (Ollama/OpenAI)</label>
                                        <div class="flex gap-2 mb-4">
                                            <div class="relative flex-1 group">
                                                <div
                                                    class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-300 group-focus-within:text-primary transition-colors">
                                                    <div class="i-fa6-solid-link text-xs"></div>
                                                </div>
                                                <input v-model="settingsBaseUrl" type="text" placeholder="请输入共享 AI 服务地址"
                                                    class="w-full pl-9 pr-3 py-2.5 bg-slate-50 border-none rounded-xl focus:ring-2 focus:ring-primary/20 focus:bg-white text-sm transition-all" />
                                            </div>
                                            <button @click="testAndFetchModels"
                                                class="px-4 py-2.5 bg-white border border-gray-200 text-gray-600 rounded-xl hover:bg-slate-50 hover:border-primary/30 hover:text-primary transition-all text-xs font-bold flex items-center gap-2 shrink-0">
                                                <div class="i-fa6-solid-rotate"></div>
                                                <span>获取</span>
                                            </button>
                                        </div>
                                        <button @click="resetSettingsToEnv"
                                            class="mb-4 ml-1 text-[11px] font-bold text-primary hover:text-blue-700 transition-colors">
                                            使用后端默认配置
                                        </button>

                                        <label
                                            class="block text-[11px] font-bold text-gray-400 uppercase tracking-wider mb-1.5 ml-1">API
                                            密钥 (可选)</label>
                                        <div class="relative flex-1 group">
                                            <div
                                                class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-300 group-focus-within:text-primary transition-colors">
                                                <div class="i-fa6-solid-key text-xs"></div>
                                            </div>
                                            <input v-model="settingsApiKey" type="password" placeholder="sk-..."
                                                class="w-full pl-9 pr-3 py-2.5 bg-slate-50 border-none rounded-xl focus:ring-2 focus:ring-primary/20 focus:bg-white text-sm transition-all" />
                                        </div>
                                        <p class="mt-2 px-1 text-[10px] leading-normal text-gray-400">
                                            这里支持任意共享 Ollama 地址，也支持 OpenAI 兼容接口地址，例如 `https://api.openai.com/v1`。
                                            API Key 仅保存在当前浏览器会话中，不会写入已保存方案。
                                        </p>
                                    </div>

                                </div>
                            </div>

                            <div class="bg-white rounded-[1.5rem] p-5 shadow-sm border border-gray-100/50">
                                <div class="flex items-center gap-2 mb-4 text-gray-800">
                                    <div class="i-fa6-solid-layer-group text-sm text-emerald-500"></div>
                                    <h4 class="text-sm font-bold">连接方案</h4>
                                </div>
                                <div class="space-y-4">
                                    <div>
                                        <label
                                            class="block text-[11px] font-bold text-gray-400 uppercase tracking-wider mb-1.5 ml-1">方案名称</label>
                                        <input v-model="profileName" type="text" placeholder="例如：宿主机 Ollama / OpenAI"
                                            class="w-full px-3 py-2.5 bg-slate-50 border-none rounded-xl focus:ring-2 focus:ring-emerald-500/20 focus:bg-white text-sm transition-all" />
                                    </div>
                                    <div class="flex gap-2">
                                        <button @click="saveCurrentProfile"
                                            class="flex-1 px-4 py-2.5 bg-emerald-500 text-white rounded-xl hover:bg-emerald-600 transition-all text-xs font-bold">
                                            保存当前方案
                                        </button>
                                        <button @click="removeCurrentProfile" :disabled="!activeProfileId"
                                            class="px-4 py-2.5 bg-white border border-gray-200 text-gray-500 rounded-xl hover:bg-slate-50 transition-all text-xs font-bold disabled:opacity-50 disabled:cursor-not-allowed">
                                            删除
                                        </button>
                                    </div>
                                    <div v-if="savedProfiles.length" class="space-y-2">
                                        <button v-for="profile in savedProfiles" :key="profile.id"
                                            @click="applyProfile(profile)"
                                            class="w-full rounded-2xl border px-4 py-3 text-left transition-all"
                                            :class="activeProfileId === profile.id ? 'border-emerald-300 bg-emerald-50' : 'border-gray-100 bg-slate-50 hover:border-emerald-200 hover:bg-emerald-50/60'">
                                            <div class="flex items-center justify-between gap-3">
                                                <div class="min-w-0">
                                                    <div class="truncate text-sm font-semibold text-gray-800">{{
                                                        profile.name }}</div>
                                                    <div class="truncate text-[11px] text-gray-500">{{ profile.baseUrl
                                                        || '后端默认配置' }}</div>
                                                </div>
                                                <div class="shrink-0 text-[10px] font-bold uppercase tracking-wide"
                                                    :class="activeProfileId === profile.id ? 'text-emerald-600' : 'text-gray-400'">
                                                    {{ detectProviderLabel(profile.baseUrl, profile.provider) }}
                                                </div>
                                            </div>
                                            <div v-if="profile.model" class="mt-2 text-[11px] text-gray-500">模型：{{
                                                profile.model }}</div>
                                        </button>
                                    </div>
                                    <div v-else class="rounded-xl bg-slate-50 px-4 py-3 text-[11px] text-gray-400">
                                        还没有保存连接方案。你可以为不同的 Ollama 服务或 OpenAI 兼容服务各存一套配置。
                                    </div>
                                </div>
                            </div>

                            <!-- 模型配置卡片 -->
                            <div class="bg-white rounded-[1.5rem] p-5 shadow-sm border border-gray-100/50">
                                <div class="flex items-center gap-2 mb-4 text-gray-800">
                                    <div class="i-fa6-solid-microchip text-sm text-blue-500"></div>
                                    <h4 class="text-sm font-bold">模型参数</h4>
                                </div>
                                <div class="space-y-4">
                                    <div class="relative">
                                        <label
                                            class="block text-[11px] font-bold text-gray-400 uppercase tracking-wider mb-1.5 ml-1">运行模型</label>
                                        <Dropdown placement="bottom-left" width="w-full">
                                            <template #trigger>
                                                <div class="relative group/model">
                                                    <div
                                                        class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-300 group-focus-within:text-blue-500 transition-colors">
                                                        <div class="i-fa6-solid-cube text-xs"></div>
                                                    </div>
                                                    <input v-model="settingsModel" type="text" placeholder="qwen2.5:7b"
                                                        class="w-full pl-9 pr-10 py-2.5 bg-slate-50 border-none rounded-xl focus:ring-2 focus:ring-blue-500/20 focus:bg-white text-sm transition-all shadow-none" />
                                                    <div
                                                        class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 pointer-events-none transition-transform group-hover/model:translate-y-[-40%]">
                                                        <div class="i-fa6-solid-chevron-down text-[10px]"></div>
                                                    </div>
                                                </div>
                                            </template>

                                            <div class="max-h-60 overflow-y-auto py-2 custom-scrollbar p-1">
                                                <div v-if="availableModels.length === 0" class="px-4 py-6 text-center">
                                                    <div
                                                        class="i-fa6-solid-circle-exclamation text-2xl text-gray-200 mb-2 mx-auto">
                                                    </div>
                                                    <p class="text-xs text-gray-400">请先获取可用模型</p>
                                                </div>
                                                <div v-else class="space-y-1">
                                                    <button v-for="model in availableModels" :key="model"
                                                        @click="settingsModel = model"
                                                        class="w-full text-left px-4 py-2.5 text-sm text-gray-600 hover:bg-blue-50 hover:text-primary rounded-xl transition-all flex items-center justify-between group/item">
                                                        <div class="flex items-center gap-3 min-w-0">
                                                            <div class="w-2 h-2 rounded-full group-hover/item:bg-primary transition-colors"
                                                                :class="settingsModel === model ? 'bg-primary' : 'bg-gray-200'">
                                                            </div>
                                                            <span class="font-medium truncate"
                                                                :class="settingsModel === model ? 'text-gray-900' : ''">{{
                                                                    model
                                                                }}</span>
                                                        </div>
                                                        <div v-if="settingsModel === model"
                                                            class="i-fa6-solid-check text-primary text-xs shrink-0">
                                                        </div>
                                                    </button>
                                                </div>
                                            </div>
                                        </Dropdown>
                                    </div>
                                    <div
                                        class="flex items-start gap-3 p-3 bg-amber-50 rounded-xl border border-amber-100">
                                        <div class="i-fa6-solid-circle-info text-amber-500 mt-0.5 text-xs"></div>
                                        <p class="text-[10px] text-amber-700 leading-normal font-medium">
                                            提示：留空将回退至后端环境变量配置 `AI_BASE_URL`。更推荐把共享 AI 服务保存成连接方案长期复用。
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div class="pt-6 mt-auto">
                            <button @click="saveSettings"
                                class="w-full py-4 bg-gradient-to-r from-primary to-blue-600 text-white rounded-[1.2rem] hover:shadow-xl hover:shadow-blue-500/30 transition-all font-bold text-sm shadow-lg shadow-blue-500/20 active:scale-[0.98]">
                                应用并保存设置
                            </button>
                        </div>
                    </div>
                </Transition>


                <!-- 消息滚动区域 -->
                <div ref="scrollRef" class="flex-1 overflow-y-auto p-4 pt-20 pb-4 space-y-6 bg-slate-50 relative">
                    <!-- 背景装饰纹理 -->
                    <div class="absolute inset-0 opacity-[0.03] bg-noise pointer-events-none">
                    </div>

                    <!-- 欢迎消息（首条非历史消息） -->
                    <div class="flex gap-4">
                        <div
                            class="w-8 h-8 rounded-full bg-gradient-to-br from-primary to-blue-600 flex shrink-0 items-center justify-center text-white text-xs shadow-md">
                            <div class="i-fa6-solid-robot"></div>
                        </div>
                        <div class="flex flex-col gap-1 items-start max-w-[85%]">
                            <div
                                class="bg-white px-4 py-3 rounded-2xl rounded-tl-none shadow-sm border border-gray-100 text-sm text-gray-700 leading-relaxed shadow-slate-200">
                                <p>你好！我是你的专属 AI 音乐助手 🎵<br>我可以帮你寻找音乐、创建情绪歌单，或者聊聊音乐背后的故事。</p>
                            </div>
                            <!-- 快捷操作建议 -->
                            <div class="flex flex-wrap gap-2 mt-2">
                                <button v-for="action in quickActions" :key="action" @click="sendQuickAction(action)"
                                    class="px-3 py-1.5 bg-white border border-gray-200 rounded-full text-xs text-primary hover:bg-blue-50 hover:border-blue-200 transition-all shadow-sm">
                                    {{ action }}
                                </button>
                            </div>
                        </div>
                    </div>

                    <!-- 聊天记录列表 -->
                    <template v-for="(msg, index) in messages" :key="index">
                        <!-- 用户发送的消息 -->
                        <div v-if="msg.role === 'user'" class="flex gap-4 flex-row-reverse">
                            <div
                                class="w-8 h-8 rounded-full bg-gray-200 flex shrink-0 items-center justify-center text-gray-500 overflow-hidden">
                                <img v-if="userAvatarUrl" :src="userAvatarUrl" :alt="userAvatarAlt"
                                    class="h-full w-full object-cover" loading="lazy" />
                                <span v-else-if="userInitial" class="text-xs font-bold text-gray-600">{{ userInitial
                                    }}</span>
                                <div v-else class="i-fa6-solid-user text-sm"></div>
                            </div>
                            <div
                                class="bg-primary text-white px-4 py-3 rounded-2xl rounded-tr-none shadow-md shadow-blue-500/20 text-sm leading-relaxed whitespace-pre-wrap max-w-[85%]">
                                {{ msg.content }}
                            </div>
                        </div>

                        <!-- AI 回复的消息 -->
                        <div v-else class="flex gap-4">
                            <div
                                class="w-8 h-8 rounded-full bg-gradient-to-br from-primary to-blue-600 flex shrink-0 items-center justify-center text-white text-xs shadow-md">
                                <div class="i-fa6-solid-robot"></div>
                            </div>
                            <div
                                class="bg-white px-4 py-3 rounded-2xl rounded-tl-none shadow-sm border border-gray-100 text-sm text-gray-700 leading-relaxed shadow-slate-200 max-w-[85%] min-w-0 overflow-hidden">
                                <div class="whitespace-pre-wrap break-words mb-1">{{ msg.content }}</div>
                                <div v-if="msg.loading"
                                    class="inline-flex items-center gap-1.5 h-4 px-1 align-middle mt-1">
                                    <div class="typing-dot"></div>
                                    <div class="typing-dot"></div>
                                    <div class="typing-dot"></div>
                                </div>
                                <div v-if="isGenerating && index === messages.length - 1 && !msg.loading"
                                    class="inline-block w-1 h-3.5 ml-1 bg-primary/40 animate-pulse align-middle"></div>
                                <Transition enter-active-class="transition-all duration-150 ease-out"
                                    enter-from-class="opacity-0 translate-y-1"
                                    enter-to-class="opacity-100 translate-y-0">
                                    <div v-if="msg.showSongs && msg.songs?.length" class="mt-4 grid gap-2 w-full">
                                        <div v-for="song in msg.songs" :key="song.songId"
                                            @click="playRecommendedSong(song, msg.songs)"
                                            class="group flex items-center gap-3 rounded-2xl border border-gray-100 bg-slate-50/90 px-3 py-2 text-left transition-colors hover:border-primary/20 hover:bg-blue-50 cursor-pointer w-full min-w-0">
                                            <img :src="song.coverUrl || '/default.jpg'" :alt="song.title"
                                                class="h-12 w-12 shrink-0 rounded-xl object-cover bg-gray-100" />
                                            <div class="flex-1 min-w-0 pr-1">
                                                <div class="truncate text-sm font-semibold text-gray-800 w-full">{{
                                                    song.title }}</div>
                                                <div class="truncate text-xs text-gray-500 w-full">{{ song.artist
                                                }}<span v-if="song.album"> · {{ song.album }}</span></div>
                                            </div>
                                            <div
                                                class="flex items-center gap-1 shrink-0 opacity-80 transition-opacity group-hover:opacity-100 ml-auto">
                                                <button
                                                    class="flex h-8 w-8 items-center justify-center rounded-full bg-white text-primary shadow-sm transition-transform hover:scale-110 active:scale-95"
                                                    @click.stop="playRecommendedSong(song, msg.songs)">
                                                    <div class="i-fa6-solid-play text-xs ml-0.5"></div>
                                                </button>
                                                <button
                                                    class="flex h-8 w-8 items-center justify-center rounded-full bg-white text-gray-500 shadow-sm transition-transform hover:scale-110 active:scale-95 hover:text-primary"
                                                    @click.stop="addRecommendedToQueue(song)">
                                                    <div class="i-fa6-solid-plus text-xs"></div>
                                                </button>
                                            </div>
                                        </div>
                                    </div>
                                </Transition>
                            </div>
                        </div>
                    </template>
                </div>

                <!-- 底部输入区域 -->
                <div class="p-4 bg-white border-t border-gray-100 relative z-20">
                    <div
                        class="relative flex items-end gap-2 bg-gray-100/50 p-2 rounded-[1.5rem] border border-transparent focus-within:border-primary/20 focus-within:bg-white focus-within:ring-4 focus-within:ring-primary/5 transition-all">
                        <textarea v-model="input" rows="1" @keydown.enter.prevent="sendMessage" placeholder="告诉我要听什么..."
                            class="w-full bg-transparent border-none focus:ring-0 text-sm text-gray-700 placeholder-gray-400 resize-none py-3 px-3 max-h-32 min-h-[44px]"
                            style="field-sizing: content;"></textarea>

                        <button @click="sendMessage" :disabled="!isGenerating && !input.trim()"
                            class="w-10 h-10 rounded-full bg-primary text-white flex items-center justify-center shrink-0 mb-0.5 shadow-lg shadow-blue-500/30 hover:shadow-blue-500/40 hover:scale-105 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:shadow-none transition-all">
                            <div v-if="!isGenerating" class="i-fa6-solid-paper-plane text-sm ml-0.5"></div>
                            <div v-else class="i-fa6-solid-stop text-sm"></div>
                        </button>
                    </div>
                    <div class="text-center mt-2">
                        <span class="text-[10px] text-gray-400">Powered by {{ providerName }} {{ modelName }}</span>
                    </div>
                </div>
            </div>
        </Transition>
    </div>
</template>

<script setup lang="ts">
import { computed, ref, nextTick, onMounted, onUnmounted } from 'vue'
import { getAIStatus, streamAIChat } from '@/api/ai'
import { formatSongs, showToast } from '@/utils/common'
import { usePlayerStore, type Song as PlayerSong } from '@/stores/player'
import { useUserStore } from '@/stores/user'
import type { SongDTO } from '@/types/api'
import Dropdown from '@/components/ui/Dropdown.vue'

const isOpen = ref(false)
const input = ref('')
const isGenerating = ref(false)
const scrollRef = ref<HTMLElement | null>(null)
const aiOnline = ref(false)
const modelName = ref('qwen2.5:7b')
const providerName = ref('AI')
let abortController: AbortController | null = null
const player = usePlayerStore()
const userStore = useUserStore()

const userAvatarUrl = computed(() => userStore.avatarUrl)
const userInitial = computed(() => (userStore.user?.username?.trim().charAt(0) || '').toUpperCase())
const userAvatarAlt = computed(() => userStore.user?.username ? `${userStore.user.username} Avatar` : 'User Avatar')

const showSettings = ref(false)
const availableModels = ref<string[]>([])
type AIProfile = {
    id: string
    name: string
    baseUrl: string
    model: string
    provider: AIProvider
}

type AIProvider = 'auto' | 'ollama' | 'openai'

const providerOptions: Array<{ value: AIProvider, label: string }> = [
    { value: 'auto', label: '自动识别' },
    { value: 'ollama', label: 'Ollama' },
    { value: 'openai', label: 'OpenAI 兼容' },
]

const loadProfiles = (): AIProfile[] => {
    try {
        const raw = localStorage.getItem('ai_profiles')
        if (!raw) return []
        const parsed = JSON.parse(raw)
        if (!Array.isArray(parsed)) return []
        return parsed
            .map((item) => ({
                id: String(item?.id || ''),
                name: String(item?.name || ''),
                baseUrl: String(item?.baseUrl || ''),
                model: String(item?.model || ''),
                provider: normalizeProviderValue(item?.provider),
            }))
            .filter((item) => item.id)
    } catch {
        return []
    }
}

const getActiveStoredProfile = () => {
    const activeId = localStorage.getItem('ai_profile_id') || ''
    if (!activeId) return null
    return loadProfiles().find((item) => item.id === activeId) || null
}

const loadProfileApiKeys = (): Record<string, string> => {
    try {
        const raw = sessionStorage.getItem('ai_profile_api_keys')
        if (!raw) return {}
        const parsed = JSON.parse(raw)
        if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) return {}
        const entries = Object.entries(parsed)
            .map(([key, value]) => [String(key), String(value || '')] as const)
            .filter(([, value]) => value.trim())
        return Object.fromEntries(entries)
    } catch {
        return {}
    }
}

const initialProfile = getActiveStoredProfile()
const profileApiKeys = ref<Record<string, string>>(loadProfileApiKeys())
const settingsBaseUrl = ref(initialProfile?.baseUrl || localStorage.getItem('ai_base_url') || '')
const settingsModel = ref(initialProfile?.model || localStorage.getItem('ai_model') || '')
const settingsApiKey = ref(
    (initialProfile?.id ? profileApiKeys.value[initialProfile.id] : '') ||
    sessionStorage.getItem('ai_api_key') ||
    '',
)
const settingsProvider = ref<AIProvider>(normalizeProviderValue(initialProfile?.provider || localStorage.getItem('ai_provider')))
const savedProfiles = ref<AIProfile[]>(loadProfiles())
const activeProfileId = ref(initialProfile?.id || '')
const profileName = ref(initialProfile?.name || '')
const statusBaseUrl = ref('')
const statusHint = ref('')

const persistProfiles = () => {
    localStorage.setItem('ai_profiles', JSON.stringify(savedProfiles.value))
}

const persistProfileApiKeys = () => {
    sessionStorage.setItem('ai_profile_api_keys', JSON.stringify(profileApiKeys.value))
}

const getOptionalSetting = (value: string) => {
    const trimmed = value.trim()
    return trimmed ? trimmed : undefined
}

const getOptionalProvider = (provider: AIProvider) => provider === 'auto' ? undefined : provider

const hasAuthToken = () => Boolean(localStorage.getItem('token'))

function normalizeProviderValue(provider?: string | null): AIProvider {
    switch ((provider || '').trim().toLowerCase()) {
        case 'ollama':
            return 'ollama'
        case 'openai':
            return 'openai'
        default:
            return 'auto'
    }
}

const detectProviderLabel = (baseUrl?: string, provider: AIProvider = 'auto') => {
    if (provider === 'ollama') return 'Ollama'
    if (provider === 'openai') return 'OpenAI'
    const normalized = (baseUrl || '').trim().toLowerCase()
    if (!normalized) return '自动识别'
    if (normalized.includes('/v1') || normalized.includes('openai')) return 'OpenAI'
    return '自动识别'
}

const buildProfileName = () => {
    return profileName.value.trim() ||
        `${detectProviderLabel(settingsBaseUrl.value, settingsProvider.value)} ${settingsModel.value.trim() || 'Default'}`
}

const connectionSourceLabel = computed(() => {
    if (activeProfileId.value) {
        const activeProfile = savedProfiles.value.find((item) => item.id === activeProfileId.value)
        if (activeProfile?.name) return activeProfile.name
    }
    if (settingsBaseUrl.value.trim() || settingsModel.value.trim() || settingsApiKey.value.trim()) {
        return '当前手动配置'
    }
    return '后端默认配置'
})

const usingBackendDefaults = computed(() => {
    return !activeProfileId.value
        && !settingsBaseUrl.value.trim()
        && !settingsModel.value.trim()
        && !settingsApiKey.value.trim()
        && settingsProvider.value === 'auto'
})

const syncActiveProfile = () => {
    if (!activeProfileId.value) return
    const index = savedProfiles.value.findIndex((item) => item.id === activeProfileId.value)
    if (index === -1) return
    const currentProfile = savedProfiles.value[index]
    if (!currentProfile) return
    savedProfiles.value[index] = {
        ...currentProfile,
        name: buildProfileName(),
        baseUrl: settingsBaseUrl.value.trim(),
        model: settingsModel.value.trim(),
        provider: settingsProvider.value,
    }
    const apiKey = settingsApiKey.value.trim()
    if (apiKey) profileApiKeys.value[activeProfileId.value] = apiKey
    else delete profileApiKeys.value[activeProfileId.value]
    persistProfileApiKeys()
    profileName.value = savedProfiles.value[index]?.name || ''
    persistProfiles()
}


const testAndFetchModels = async () => {
    if (!hasAuthToken()) {
        aiOnline.value = false
        availableModels.value = []
        statusHint.value = '登录后可检测和使用 AI 服务'
        showToast('请先登录后再配置 AI 助手', 'warning')
        return
    }

    try {
        const res = await getAIStatus(
            getOptionalSetting(settingsBaseUrl.value),
            getOptionalSetting(settingsModel.value),
            getOptionalSetting(settingsApiKey.value),
            getOptionalProvider(settingsProvider.value)
        ) as any
        const data = res?.data || {}

        aiOnline.value = Boolean(data.online)
        availableModels.value = Array.isArray(data.models) ? data.models : []
        providerName.value = data.provider || detectProviderLabel(settingsBaseUrl.value, settingsProvider.value)
        modelName.value = data.model || settingsModel.value || modelName.value
        statusBaseUrl.value = data.baseUrl || settingsBaseUrl.value.trim()
        statusHint.value = data.online
            ? (usingBackendDefaults.value
                ? `当前使用后端默认配置，检测到 ${availableModels.value.length || 0} 个可用模型`
                : `已检测到 ${availableModels.value.length || 0} 个可用模型`)
            : (data.error || (data.models?.length ? '服务在线，但当前模型不可用' : '服务在线，但未检测到模型'))
        if (data.online) {
            showToast(usingBackendDefaults.value
                ? `已切换到后端默认配置，当前服务类型：${providerName.value}`
                : `连接成功，当前服务类型：${providerName.value}`)
        } else {
            showToast(data.error || '连接失败，请检查服务地址、模型和网络连通性', 'warning')
        }
    } catch {
        aiOnline.value = false
        availableModels.value = []
        statusHint.value = '连接检测失败，请检查服务地址、模型和网络连通性'
        showToast('连接失败，请检查地址', 'error')
    }
}

const saveSettings = () => {
    const baseUrl = settingsBaseUrl.value.trim()
    const model = settingsModel.value.trim()
    const apiKey = settingsApiKey.value.trim()

    if (baseUrl) localStorage.setItem('ai_base_url', baseUrl)
    else localStorage.removeItem('ai_base_url')

    if (model) localStorage.setItem('ai_model', model)
    else localStorage.removeItem('ai_model')

    if (apiKey) sessionStorage.setItem('ai_api_key', apiKey)
    else sessionStorage.removeItem('ai_api_key')

    localStorage.setItem('ai_provider', settingsProvider.value)
    syncActiveProfile()
    showSettings.value = false
    void fetchAIStatus()
}

const resetSettingsToEnv = () => {
    settingsBaseUrl.value = ''
    settingsModel.value = ''
    settingsApiKey.value = ''
    settingsProvider.value = 'auto'
    localStorage.removeItem('ai_base_url')
    localStorage.removeItem('ai_model')
    localStorage.removeItem('ai_provider')
    localStorage.removeItem('ai_profile_id')
    sessionStorage.removeItem('ai_api_key')
    activeProfileId.value = ''
    profileName.value = ''
    availableModels.value = []
    statusBaseUrl.value = ''
    statusHint.value = '已切换为后端默认配置；本地开发时会读取后端 config.local.yaml 或环境变量'
    showToast('已切换为后端默认配置')
    void fetchAIStatus()
}

const applyProfile = (profile: AIProfile) => {
    activeProfileId.value = profile.id
    profileName.value = profile.name
    settingsBaseUrl.value = profile.baseUrl
    settingsModel.value = profile.model
    settingsApiKey.value = profileApiKeys.value[profile.id] || ''
    settingsProvider.value = normalizeProviderValue(profile.provider)
    localStorage.setItem('ai_profile_id', profile.id)
    if (settingsApiKey.value.trim()) sessionStorage.setItem('ai_api_key', settingsApiKey.value.trim())
    else sessionStorage.removeItem('ai_api_key')
    void fetchAIStatus()
}

const saveCurrentProfile = () => {
    const id = activeProfileId.value || `${Date.now()}`
    const profile: AIProfile = {
        id,
        name: buildProfileName(),
        baseUrl: settingsBaseUrl.value.trim(),
        model: settingsModel.value.trim(),
        provider: settingsProvider.value,
    }
    const index = savedProfiles.value.findIndex((item) => item.id === id)
    if (index >= 0) savedProfiles.value[index] = profile
    else savedProfiles.value.unshift(profile)
    activeProfileId.value = id
    profileName.value = profile.name
    localStorage.setItem('ai_profile_id', id)
    persistProfiles()
    saveSettings()
    showToast('连接方案已保存')
}

const removeCurrentProfile = () => {
    if (!activeProfileId.value) return
    savedProfiles.value = savedProfiles.value.filter((item) => item.id !== activeProfileId.value)
    persistProfiles()
    resetSettingsToEnv()
    showToast('连接方案已删除')
}


interface Message {
    role: 'user' | 'assistant'
    content: string
    loading?: boolean
    songs?: PlayerSong[]
    showSongs?: boolean
}

const messages = ref<Message[]>([])

const quickActions = [
    '推荐几首周杰伦的歌',
    '适合学习听的音乐',
    '今天心情不好求安慰'
]

const scrollToBottom = async () => {
    await nextTick()
    if (scrollRef.value) {
        scrollRef.value.scrollTop = scrollRef.value.scrollHeight
    }
}

const close = () => {
    isOpen.value = false
}

// 全局事件监听：用于外部触发面板开关
const handleToggle = () => {
    isOpen.value = !isOpen.value
}

onMounted(() => {
    localStorage.removeItem('ai_api_key')
    persistProfiles()
    window.addEventListener('toggle-ai-chat', handleToggle)
    void fetchAIStatus()
})

onUnmounted(() => {
    window.removeEventListener('toggle-ai-chat', handleToggle)
    abortController?.abort()
})

const sendQuickAction = (text: string) => {
    input.value = text
    sendMessage()
}

const fetchAIStatus = async () => {
    if (!hasAuthToken()) {
        aiOnline.value = false
        availableModels.value = []
        statusBaseUrl.value = settingsBaseUrl.value.trim()
        statusHint.value = '登录后可使用 AI 助手'
        return
    }

    try {
        const res = await getAIStatus(
            getOptionalSetting(settingsBaseUrl.value),
            getOptionalSetting(settingsModel.value),
            getOptionalSetting(settingsApiKey.value),
            getOptionalProvider(settingsProvider.value)
        ) as any
        const data = res?.data || {}
        aiOnline.value = Boolean(data.online)
        providerName.value = data.provider || detectProviderLabel(settingsBaseUrl.value, settingsProvider.value)
        modelName.value = data.model || settingsModel.value || 'qwen2.5:7b'
        availableModels.value = Array.isArray(data.models) ? data.models : []
        statusBaseUrl.value = data.baseUrl || settingsBaseUrl.value.trim()
        statusHint.value = data.online
            ? (usingBackendDefaults.value
                ? `当前使用后端默认配置：${data.baseUrl || '由后端配置决定'}`
                : `已连接到 ${data.baseUrl || '后端默认地址'}`)
            : (data.error || 'AI 服务暂时不可用')

    } catch {
        aiOnline.value = false
        availableModels.value = []
        statusBaseUrl.value = settingsBaseUrl.value.trim()
        statusHint.value = '无法获取当前连接状态'
    }
}

const stopGeneration = () => {
    abortController?.abort()
    abortController = null
    isGenerating.value = false

    const lastMessage = messages.value[messages.value.length - 1]
    if (lastMessage?.role === 'assistant' && lastMessage.loading) {
        messages.value.pop()
    }
}

const playRecommendedSong = (song: PlayerSong, contextSongs?: PlayerSong[]) => {
    if (contextSongs && contextSongs.length > 0) {
        player.setPlayList(contextSongs, 'AI 推荐')
    }
    player.playSong(song)
    player.playlistSource = 'AI 推荐'
    showToast(`开始播放《${song.title}》`)
}

const addRecommendedToQueue = (song: PlayerSong) => {
    player.addToQueue(song)
}

const applyStreamEvent = (eventName: string, payload: any, target: Message) => {
    if (eventName === 'meta') {
        modelName.value = payload?.model || modelName.value
        providerName.value = payload?.provider || providerName.value
        const songs = Array.isArray(payload?.songs)
            ? formatSongs(payload.songs as SongDTO[])
            : []
        target.songs = songs
        target.showSongs = false
        return
    }

    if (eventName === 'chunk') {
        const content = payload?.content || ''
        if (content) {
            target.loading = false
            target.content += content
        }
        return
    }

    if (eventName === 'error') {
        target.loading = false
        target.showSongs = false
        target.content = payload?.message || 'AI 回复失败，请稍后重试。'
        return
    }

    if (eventName === 'done') {
        target.loading = false
        if (target.songs?.length) {
            window.setTimeout(() => {
                target.showSongs = true
            }, 120)
        }
    }
}

const consumeSSEStream = async (reader: ReadableStreamDefaultReader<Uint8Array>, target: Message) => {
    const decoder = new TextDecoder('utf-8')
    let buffer = ''

    while (true) {
        const { value, done } = await reader.read()
        if (done) break

        buffer += decoder.decode(value, { stream: true })
        const frames = buffer.split('\n\n')
        buffer = frames.pop() || ''

        for (const frame of frames) {
            const lines = frame.split('\n').map((line) => line.trim()).filter(Boolean)
            const eventLine = lines.find((line) => line.startsWith('event:'))
            const dataLine = lines.find((line) => line.startsWith('data:'))
            if (!eventLine || !dataLine) continue

            const eventName = eventLine.replace(/^event:\s*/, '')
            const rawData = dataLine.replace(/^data:\s*/, '')
            try {
                applyStreamEvent(eventName, JSON.parse(rawData), target)
            } catch {
                continue
            }
        }

        scrollToBottom()
    }
}

const sendMessage = async () => {
    if (isGenerating.value) {
        stopGeneration()
        return
    }

    const text = input.value.trim()
    if (!text) return
    if (!hasAuthToken()) {
        showToast('请先登录后再使用 AI 助手', 'warning')
        return
    }

    // 添加用户消息并更新界面状态
    messages.value.push({ role: 'user', content: text })
    input.value = ''
    scrollToBottom()

    isGenerating.value = true
    const aiMsgIndex = messages.value.push({ role: 'assistant', content: '', loading: true }) - 1
    scrollToBottom()

    abortController = new AbortController()

    try {
        const history = messages.value
            .filter((item) => !item.loading)
            .slice(-12)
            .map((item) => ({
                role: item.role,
                content: item.content,
            }))

        const aiMsg = messages.value[aiMsgIndex]
        const response = await streamAIChat(
            history,
            getOptionalSetting(settingsBaseUrl.value),
            getOptionalSetting(settingsModel.value),
            getOptionalSetting(settingsApiKey.value),
            getOptionalProvider(settingsProvider.value),
            abortController.signal
        )
        const reader = response.body?.getReader()


        aiOnline.value = true

        if (!reader || !aiMsg) {
            throw new Error('AI stream unavailable')
        }

        await consumeSSEStream(reader, aiMsg)

        if (!aiMsg.content.trim()) {
            aiMsg.loading = false
            aiMsg.content = '暂时没有生成有效内容。'
        }
    } catch (error: any) {
        if (error?.name === 'CanceledError' || error?.name === 'AbortError') {
            return
        }

        const aiMsg = messages.value[aiMsgIndex]
        if (aiMsg) {
            aiMsg.loading = false
            aiMsg.content = error?.message || 'AI 服务暂时不可用，请检查当前连接方案与模型配置。'
        }

        aiOnline.value = false
        showToast(error?.message || 'AI 助手暂时不可用', 'error')
    } finally {
        abortController = null
        isGenerating.value = false
        scrollToBottom()
        void fetchAIStatus()
    }
}
</script>

<style scoped>
/* 滚动条样式 */
::-webkit-scrollbar {
    width: 6px;
}

::-webkit-scrollbar-track {
    background: transparent;
}

::-webkit-scrollbar-thumb {
    background: #e2e8f0;
    border-radius: 10px;
}

::-webkit-scrollbar-thumb:hover {
    background: #cbd5e1;
}

/* AI 输入动画 */
.typing-dot {
    width: 6px;
    height: 6px;
    background-color: #94a3b8;
    border-radius: 50%;
    animation: typing-bounce 1.4s infinite ease-in-out both;
}

.typing-dot:nth-child(1) {
    animation-delay: -0.32s;
}

.typing-dot:nth-child(2) {
    animation-delay: -0.16s;
}

@keyframes typing-bounce {

    0%,
    80%,
    100% {
        transform: scale(0.6);
        opacity: 0.4;
    }

    40% {
        transform: scale(1);
        opacity: 1;
    }
}
</style>
