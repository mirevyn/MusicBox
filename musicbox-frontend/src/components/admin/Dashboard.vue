<template>
  <div class="space-y-8">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-800 mb-1">仪表盘概览</h2>
        <p class="text-gray-500 text-sm">欢迎回来，查看最新的系统数据。</p>
      </div>
      <div class="flex gap-3">
        <button
          class="bg-primary text-white hover:bg-primary/90 transition-colors rounded-lg px-4 py-2 text-sm font-medium shadow-lg shadow-primary/20 disabled:opacity-60 disabled:cursor-not-allowed"
          :disabled="exporting || loading" @click="confirmExportReport">
          <div :class="exporting ? 'i-fa6-solid-spinner animate-spin' : 'i-fa6-solid-download'" class="mr-2"></div>
          {{ exporting ? '导出中...' : '导出报告' }}
        </button>
      </div>
    </div>

    <!-- 顶层统计卡片区域 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div
        class="bg-white rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition-all duration-300 p-6 relative overflow-hidden group">
        <div
          class="absolute right-0 top-0 w-32 h-32 bg-blue-500/5 rounded-full -mr-10 -mt-10 group-hover:scale-110 transition-transform duration-500">
        </div>
        <div class="relative z-10">
          <div class="flex items-center justify-between mb-4">
            <div class="w-12 h-12 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center text-xl">
              <div class="i-fa6-solid-music"></div>
            </div>
            <div class="text-xs font-bold px-2.5 py-1 rounded-full flex items-center gap-1 border"
              :class="trendClass(stats.songTrendPercent)">
              <div :class="trendIcon(stats.songTrendPercent)"></div>
              {{ loading ? '--' : `本周新增 ${stats.newSongs} (${formatTrend(stats.songTrendPercent)})` }}
            </div>
          </div>
          <div class="text-gray-500 text-sm font-medium mb-1">总歌曲数</div>
          <div class="text-3xl font-black text-gray-800 tracking-tight">
            {{ loading ? '--' : formatNumber(stats.totalSongs) }}
          </div>
        </div>
      </div>

      <div
        class="bg-white rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition-all duration-300 p-6 relative overflow-hidden group">
        <div
          class="absolute right-0 top-0 w-32 h-32 bg-purple-500/5 rounded-full -mr-10 -mt-10 group-hover:scale-110 transition-transform duration-500">
        </div>
        <div class="relative z-10">
          <div class="flex items-center justify-between mb-4">
            <div class="w-12 h-12 rounded-xl bg-purple-50 text-purple-600 flex items-center justify-center text-xl">
              <div class="i-fa6-solid-users"></div>
            </div>
            <div class="text-xs font-bold px-2.5 py-1 rounded-full flex items-center gap-1 border"
              :class="trendClass(stats.userTrendPercent)">
              <div :class="trendIcon(stats.userTrendPercent)"></div>
              {{ loading ? '--' : `本周新增 ${stats.newUsers} (${formatTrend(stats.userTrendPercent)})` }}
            </div>
          </div>
          <div class="text-gray-500 text-sm font-medium mb-1">注册用户</div>
          <div class="text-3xl font-black text-gray-800 tracking-tight">
            {{ loading ? '--' : formatNumber(stats.totalUsers) }}
          </div>
        </div>
      </div>

      <div
        class="bg-white rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition-all duration-300 p-6 relative overflow-hidden group">
        <div
          class="absolute right-0 top-0 w-32 h-32 bg-primary/5 rounded-full -mr-10 -mt-10 group-hover:scale-110 transition-transform duration-500">
        </div>
        <div class="relative z-10">
          <div class="flex items-center justify-between mb-4">
            <div class="w-12 h-12 rounded-xl bg-primary/10 text-primary flex items-center justify-center text-xl">
              <div class="i-fa6-solid-bolt"></div>
            </div>
            <div class="text-xs font-bold px-2.5 py-1 rounded-full flex items-center gap-1 border"
              :class="playTrendClass">
              <div :class="playTrendIcon"></div>
              {{ loading ? '--' : playTrendText }}
            </div>
          </div>
          <div class="text-gray-500 text-sm font-medium mb-1">今日播放</div>
          <div class="text-3xl font-black text-gray-800 tracking-tight">
            {{ loading ? '--' : formatPlayCount(stats.todayPlays) }}
          </div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 xl:grid-cols-5 gap-6">
      <div class="xl:col-span-3 bg-white rounded-2xl border border-gray-100 shadow-sm p-6">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h3 class="text-lg font-semibold text-gray-800">近 7 天趋势</h3>
            <p class="text-sm text-gray-400">播放、上传与注册数据走势</p>
          </div>
          <div class="text-xs font-semibold text-gray-400 bg-gray-50 px-3 py-1 rounded-full border border-gray-100">
            Last 7 Days
          </div>
        </div>
        <div v-if="loading" class="h-[360px] w-full flex flex-col items-center justify-center bg-gray-50/50 rounded-xl">
          <div class="i-fa6-solid-spinner animate-spin text-3xl text-primary/40 mb-3"></div>
          <div class="text-gray-400 text-sm">加载图表数据...</div>
        </div>
        <VChart v-else class="h-[360px] w-full" :option="trendOption" autoresize />
      </div>

      <div class="xl:col-span-2 bg-white rounded-2xl border border-gray-100 shadow-sm p-6">
        <div class="mb-4">
          <h3 class="text-lg font-semibold text-gray-800">播放来源分布</h3>
          <p class="text-sm text-gray-400">用户主要从哪些入口触发播放</p>
        </div>
        <div v-if="loading" class="h-[360px] w-full flex items-center justify-center bg-gray-50/50 rounded-xl"></div>
        <VChart v-else class="h-[360px] w-full" :option="playSourceOption" autoresize />
      </div>
    </div>

    <!-- 图表区域 2 -->
    <div class="grid grid-cols-1 xl:grid-cols-5 gap-6">
      <div class="xl:col-span-3 bg-white rounded-2xl border border-gray-100 shadow-sm p-6">
        <div class="mb-4">
          <h3 class="text-lg font-semibold text-gray-800">曲库歌手分布</h3>
          <p class="text-sm text-gray-400">当前收录歌曲最多的歌手</p>
        </div>
        <div v-if="loading" class="h-[300px] w-full flex items-center justify-center bg-gray-50/50 rounded-xl"></div>
        <VChart v-else class="h-[300px] w-full" :option="artistOption" autoresize />
      </div>

      <div class="xl:col-span-2 bg-white rounded-2xl border border-gray-100 shadow-sm p-6">
        <div class="mb-4">
          <h3 class="text-lg font-semibold text-gray-800">歌单审核状态</h3>
          <p class="text-sm text-gray-400">公开歌单当前审核分布</p>
        </div>
        <div v-if="loading" class="h-[300px] w-full flex items-center justify-center bg-gray-50/50 rounded-xl"></div>
        <VChart v-else class="h-[300px] w-full" :option="playlistStatusOption" autoresize />
      </div>
    </div>

    <!-- 导出报告确认弹窗 -->
    <ConfirmModal v-model="showExportConfirm" type="info" title="确认导出报告" content="是否确认导出包含最新统计数据的 CSV 报表？"
      :loading="exporting" @confirm="executeExport" />
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { BarChart, LineChart, PieChart } from 'echarts/charts'
import { GridComponent, LegendComponent, TooltipComponent, GraphicComponent } from 'echarts/components'
import type { EChartsOption } from 'echarts'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import { exportDashboardReport, getDashboardAnalytics, getDashboardStats } from '@/api/admin'
import { showToast } from '@/utils/common'
import { useNotificationStore } from '@/stores/notification'

use([CanvasRenderer, LineChart, PieChart, BarChart, GridComponent, TooltipComponent, LegendComponent, GraphicComponent])

interface ValueItem {
  name: string
  value: number
}

const loading = ref(true)
const exporting = ref(false)
const stats = reactive({
  totalSongs: 0,
  totalUsers: 0,
  todayPlays: 0,
  newSongs: 0,
  newUsers: 0,
  songTrendPercent: 0,
  userTrendPercent: 0,
  playTrendPercent: 0,
  playTrendDirection: 'flat'
})

const showExportConfirm = ref(false)

const analytics = reactive({
  trendSeries: {
    labels: [] as string[],
    plays: [] as number[],
    songs: [] as number[],
    users: [] as number[],
  },
  topArtists: [] as ValueItem[],
  playSourceDistribution: [] as ValueItem[],
  playlistStatus: [] as ValueItem[],
})

const formatNumber = (value: number) => value.toLocaleString('zh-CN')

const formatPlayCount = (value: number) => {
  if (value >= 10000) {
    const formatted = (value / 10000).toFixed(1)
    return `${formatted.replace(/\.0$/, '')}w`
  }
  return formatNumber(value)
}

const formatTrend = (value: number) => {
  const absValue = Math.abs(value)
  const formatted = absValue >= 100 ? Math.round(absValue) : absValue.toFixed(1).replace(/\.0$/, '')
  const prefix = value > 0 ? '+' : value < 0 ? '-' : ''
  return `${prefix}${formatted}%`
}

const trendClass = (value: number) => {
  if (value > 0) return 'text-green-600 bg-green-50 border-green-100'
  if (value < 0) return 'text-red-500 bg-red-50 border-red-100'
  return 'text-gray-500 bg-gray-50 border-gray-100'
}

const trendIcon = (value: number) => {
  if (value > 0) return 'i-fa6-solid-arrow-trend-up'
  if (value < 0) return 'i-fa6-solid-arrow-trend-down'
  return 'i-fa6-solid-minus'
}

const playTrendClass = computed(() => {
  if (stats.playTrendDirection === 'up') return 'text-green-600 bg-green-50 border-green-100'
  if (stats.playTrendDirection === 'down') return 'text-red-500 bg-red-50 border-red-100'
  return 'text-primary bg-primary/10 border-primary/20'
})

const playTrendIcon = computed(() => {
  if (stats.playTrendDirection === 'up') return 'i-fa6-solid-arrow-trend-up'
  if (stats.playTrendDirection === 'down') return 'i-fa6-solid-arrow-trend-down'
  return 'i-fa6-solid-wave-square'
})

const playTrendText = computed(() => {
  if (stats.playTrendDirection === 'flat') return '较昨日持平'
  return `${formatTrend(stats.playTrendPercent)} vs 昨日`
})

const trendOption = computed<EChartsOption>(() => ({
  color: ['#335bf1', '#10b981', '#8b5cf6'],
  tooltip: {
    trigger: 'axis',
    backgroundColor: 'rgba(15, 23, 42, 0.92)',
    borderWidth: 0,
    textStyle: { color: '#fff' },
  },
  legend: {
    bottom: 4,
    itemWidth: 10,
    itemHeight: 10,
    itemGap: 24,
    textStyle: { color: '#64748b', fontSize: 12 },
  },
  grid: {
    left: 45,
    right: 15,
    top: 20,
    bottom: 58,
  },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: analytics.trendSeries.labels,
    axisLine: { lineStyle: { color: '#e2e8f0' } },
    axisLabel: { color: '#94a3b8' },
  },
  yAxis: {
    type: 'value',
    splitLine: { lineStyle: { color: '#eef2ff' } },
    axisLabel: { color: '#94a3b8' },
  },
  series: [
    {
      name: '播放',
      type: 'line',
      smooth: true,
      symbol: 'circle',
      symbolSize: 8,
      data: analytics.trendSeries.plays,
      areaStyle: {
        color: 'rgba(51, 91, 241, 0.12)',
      },
      lineStyle: { width: 3 },
    },
    {
      name: '上传',
      type: 'line',
      smooth: true,
      symbol: 'circle',
      symbolSize: 7,
      data: analytics.trendSeries.songs,
      areaStyle: {
        color: 'rgba(16, 185, 129, 0.10)',
      },
      lineStyle: { width: 3 },
    },
    {
      name: '注册',
      type: 'line',
      smooth: true,
      symbol: 'circle',
      symbolSize: 7,
      data: analytics.trendSeries.users,
      areaStyle: {
        color: 'rgba(139, 92, 246, 0.10)',
      },
      lineStyle: { width: 3 },
    },
  ],
}))

const playSourceOption = computed<EChartsOption>(() => ({
  color: ['#335bf1', '#10b981', '#f59e0b', '#8b5cf6', '#ef4444'],
  tooltip: {
    trigger: 'item',
    confine: true,
    backgroundColor: 'rgba(15, 23, 42, 0.92)',
    borderWidth: 0,
    textStyle: { color: '#fff' },
    formatter: (params: any) => `${params.name}<br />播放次数：${params.value}<br />占比：${params.percent}%`,
  },
  legend: {
    bottom: 0,
    left: 'center',
    icon: 'circle',
    itemWidth: 10,
    itemHeight: 10,
    textStyle: { color: '#64748b', fontSize: 12 },
  },
  series: [
    {
      type: 'pie',
      radius: ['48%', '72%'],
      center: ['50%', '44%'],
      padAngle: 2,
      label: { show: false },
      emphasis: {
        scale: true,
        label: { show: false },
      },
      data: analytics.playSourceDistribution,
    },
  ],
}))

const artistOption = computed<EChartsOption>(() => ({
  color: ['#335bf1'],
  tooltip: {
    trigger: 'axis',
    axisPointer: { type: 'shadow' },
    backgroundColor: 'rgba(15, 23, 42, 0.92)',
    borderWidth: 0,
    textStyle: { color: '#fff' },
  },
  grid: {
    left: 95,
    right: 50,
    top: 8,
    bottom: 20,
  },
  xAxis: {
    type: 'value',
    splitLine: { lineStyle: { color: '#eef2ff' } },
    axisLabel: { color: '#94a3b8' },
  },
  yAxis: {
    type: 'category',
    data: analytics.topArtists.map((item) => item.name),
    inverse: true,
    axisTick: { show: false },
    axisLine: { show: false },
    axisLabel: {
      color: '#64748b',
      fontWeight: 600,
      width: 88,
      overflow: 'truncate',
    },
  },
  series: [
    {
      type: 'bar',
      data: analytics.topArtists.map((item) => item.value),
      barWidth: 14,
      showBackground: true,
      backgroundStyle: {
        color: '#f8fafc',
        borderRadius: [0, 99, 99, 0],
      },
      itemStyle: {
        borderRadius: [0, 99, 99, 0],
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 1, y2: 0,
          colorStops: [
            { offset: 0, color: '#3b82f6' },
            { offset: 1, color: '#335bf1' }
          ]
        } as any,
        shadowColor: 'rgba(51, 91, 241, 0.2)',
        shadowBlur: 10
      },
      emphasis: {
        itemStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 1, y2: 0,
            colorStops: [
              { offset: 0, color: '#60a5fa' },
              { offset: 1, color: '#2563eb' }
            ]
          } as any,
          shadowBlur: 20,
          shadowColor: 'rgba(59, 130, 246, 0.4)'
        }
      },
      label: {
        show: true,
        position: 'right',
        color: '#64748b',
        fontWeight: 600,
        distance: 10
      },
    } as any,
  ],
}))

const playlistStatusOption = computed<EChartsOption>(() => ({
  color: ['#f59e0b', '#10b981', '#ef4444'],
  tooltip: {
    trigger: 'axis',
    axisPointer: { type: 'shadow' },
    backgroundColor: 'rgba(15, 23, 42, 0.92)',
    borderWidth: 0,
    textStyle: { color: '#fff' },
  },
  grid: {
    left: 65,
    right: 15,
    top: 10,
    bottom: 20,
  },
  xAxis: {
    type: 'value',
    splitLine: { lineStyle: { color: '#f1f5f9' } },
    axisLabel: { color: '#94a3b8' },
  },
  yAxis: {
    type: 'category',
    data: analytics.playlistStatus.map((item) => item.name),
    axisTick: { show: false },
    axisLine: { show: false },
    axisLabel: { color: '#64748b', fontWeight: 600 },
  },
  series: [
    {
      type: 'bar',
      data: analytics.playlistStatus.map((item) => item.value),
      barWidth: 16,
      showBackground: true,
      backgroundStyle: {
        color: '#f8fafc',
        borderRadius: [0, 99, 99, 0],
      },
      itemStyle: {
        borderRadius: [0, 99, 99, 0],
        color: (params: any) => {
          const colors = [
            ['#fbbf24', '#f59e0b'], // 待审核
            ['#34d399', '#10b981'], // 已通过
            ['#f87171', '#ef4444']  // 已拒绝
          ]
          const colorPair = (colors[params.dataIndex % colors.length] || colors[0]) as string[]
          return {
            type: 'linear',
            x: 0, y: 0, x2: 1, y2: 0,
            colorStops: [
              { offset: 0, color: colorPair[0] || '#fbbf24' },
              { offset: 1, color: colorPair[1] || '#f59e0b' }
            ]
          } as any
        },
        shadowBlur: 8,
        shadowColor: 'rgba(0, 0, 0, 0.05)'
      },
      emphasis: {
        itemStyle: {
          shadowBlur: 15,
          shadowColor: 'rgba(0, 0, 0, 0.1)'
        }
      },
      label: {
        show: true,
        position: 'right',
        color: '#64748b',
        fontWeight: 600,
        distance: 10
      },
    } as any,
  ],
}))

const loadDashboardData = async () => {
  try {
    const [statsRes, analyticsRes] = await Promise.all([
      getDashboardStats() as any,
      getDashboardAnalytics() as any,
    ])

    const statsData = statsRes?.data || {}
    stats.totalSongs = Number(statsData.totalSongs || 0)
    stats.totalUsers = Number(statsData.totalUsers || 0)
    stats.todayPlays = Number(statsData.todayPlays || 0)
    stats.newSongs = Number(statsData.newSongs || 0)
    stats.newUsers = Number(statsData.newUsers || 0)
    stats.songTrendPercent = Number(statsData.songTrendPercent || 0)
    stats.userTrendPercent = Number(statsData.userTrendPercent || 0)
    stats.playTrendPercent = Number(statsData.playTrendPercent || 0)
    stats.playTrendDirection = statsData.playTrendDirection || 'flat'

    const analyticsData = analyticsRes?.data || {}
    analytics.trendSeries.labels = analyticsData.trendSeries?.labels || []
    analytics.trendSeries.plays = analyticsData.trendSeries?.plays || []
    analytics.trendSeries.songs = analyticsData.trendSeries?.songs || []
    analytics.trendSeries.users = analyticsData.trendSeries?.users || []
    analytics.topArtists = analyticsData.topArtists || []
    analytics.playSourceDistribution = analyticsData.playSourceDistribution || []
    analytics.playlistStatus = analyticsData.playlistStatus || []
  } finally {
    loading.value = false
  }
}

const confirmExportReport = () => {
  showExportConfirm.value = true
}

const executeExport = async () => {
  if (exporting.value) return

  exporting.value = true
  try {
    const blob = await exportDashboardReport() as unknown as Blob
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    const timestamp = new Date().toISOString().slice(0, 19).replace(/[T:]/g, '-')

    link.href = url
    link.download = `dashboard-report-${timestamp}.csv`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    showToast('报告导出成功')
    showExportConfirm.value = false
  } catch {
    showToast('导出报告失败', 'error')
  } finally {
    exporting.value = false
  }
}

const notificationStore = useNotificationStore()

const handleDashboardRefresh = () => {
  loadDashboardData()
}

const POLL_INTERVAL_MS = 30_000
let pollTimer: number | null = null

onMounted(() => {
  loadDashboardData()
  notificationStore.onDashboardRefresh(handleDashboardRefresh)
  pollTimer = window.setInterval(() => loadDashboardData(), POLL_INTERVAL_MS)
})

onBeforeUnmount(() => {
  notificationStore.offDashboardRefresh(handleDashboardRefresh)
  if (pollTimer) {
    window.clearInterval(pollTimer)
    pollTimer = null
  }
})
</script>
