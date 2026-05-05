<template>
    <div class="space-y-6">
        <div class="flex items-center justify-between">
            <div>
                <h2 class="text-2xl font-bold text-gray-800 dark:text-white mb-1">用户管理</h2>
                <p class="text-gray-500 text-sm">查看和管理注册用户</p>
            </div>
            
        </div>

        <!-- 用户表格卡片 -->
        <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden min-h-[620px] flex flex-col">
            <!-- 工具栏 -->
            <div class="p-4 border-b border-gray-100 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <!-- 左侧：搜索与筛选 -->
                <div class="flex items-center gap-3 flex-1">
                    <!-- 搜索 -->
                    <div class="relative w-64">
                        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400">
                            <div class="i-fa6-solid-magnifying-glass"></div>
                        </div>
                        <input v-model="searchKeyword" @keydown.enter="handleSearch" type="text" placeholder="搜索用户名..." 
                            class="w-full bg-gray-50 border-none rounded-lg py-2 pl-10 pr-4 text-sm focus:ring-2 focus:ring-primary/20 focus:bg-white transition-all">
                    </div>

                    <!-- 筛选器 -->
                    <Dropdown width="w-32">
                        <template #trigger>
                            <div class="w-32 px-4 py-2 bg-gray-50 rounded-lg text-sm font-medium text-gray-600 hover:bg-gray-100 transition-colors flex items-center justify-between cursor-pointer">
                                <span>{{ filterRoleLabel }}</span>
                                <div class="i-fa6-solid-chevron-down text-xs opacity-50"></div>
                            </div>
                        </template>
                        <div class="py-1">
                            <button @click="setRole('')" class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors" :class="{'text-primary font-medium': filterRole === ''}">
                                所有角色
                            </button>
                            <button @click="setRole('User')" class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors" :class="{'text-primary font-medium': filterRole === 'User'}">
                                普通用户
                            </button>
                            <button @click="setRole('Admin')" class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors" :class="{'text-primary font-medium': filterRole === 'Admin'}">
                                管理员
                            </button>
                        </div>
                    </Dropdown>

                    <button @click="handleSearch" class="px-4 py-2 bg-gray-50 rounded-lg text-sm font-medium text-gray-600 hover:bg-gray-100 transition-colors flex items-center gap-2 whitespace-nowrap">
                        搜索
                    </button>
                </div>
                
                <!-- 右侧：操作项 -->
                <div class="flex items-center gap-3">
                    <button @click="handleCreate" class="bg-primary text-white hover:bg-primary/90 transition-colors rounded-lg px-4 py-2 text-sm font-medium shadow-lg shadow-primary/20 flex items-center gap-2 whitespace-nowrap">
                        <div class="i-fa6-solid-user-plus"></div> 添加用户
                    </button>
                </div>
            </div>

            <!-- 内容区域 -->
            <div class="relative flex-1 flex flex-col">
                 <!-- 加载状态 -->
                <div v-if="loading && users.length === 0" class="absolute inset-0 flex items-center justify-center bg-white z-10">
                    <div class="flex flex-col items-center gap-3">
                        <div class="i-fa6-solid-spinner animate-spin text-2xl text-primary"></div>
                        <div class="text-sm text-gray-400">加载中...</div>
                    </div>
                </div>

                 <!-- 空状态 -->
                <div v-if="!loading && users.length === 0" class="absolute inset-0 flex items-center justify-center">
                    <div class="flex flex-col items-center gap-3 text-gray-400">
                        <div class="i-fa6-solid-users text-4xl opacity-20"></div>
                        <div>暂无用户数据</div>
                    </div>
                </div>

                <table class="w-full text-left border-collapse">
                    <thead>
                        <tr class="bg-gray-50/50 border-b border-gray-100 text-xs uppercase text-gray-500 font-semibold tracking-wider">
                            <th class="px-6 py-4 w-16">#</th>
                            <th class="px-6 py-4">用户</th>
                            <th class="px-6 py-4">角色</th>
                            <th class="px-6 py-4">状态</th>
                            <th class="px-6 py-4">注册时间</th>
                            <th class="px-6 py-4 text-right">操作</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-50">
                        <tr v-for="(user, index) in users" :key="user.id" 
                            class="hover:bg-gray-50/80 transition-colors group">
                             <td class="px-6 py-4 text-gray-400 font-medium text-sm">{{ (currentPage - 1) * pageSize + index + 1 }}</td>
                            <td class="px-6 py-4">
                                <div class="flex items-center gap-3">
                                    <template v-if="user.avatarUrl">
                                        <img :src="getAvatarUrl(user.avatarUrl)" :alt="user.username" class="w-10 h-10 rounded-full object-cover border border-gray-100 shadow-sm" />
                                    </template>
                                    <template v-else>
                                        <div class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-100 to-purple-100 flex items-center justify-center text-primary font-bold text-lg">
                                            {{ user.username.charAt(0).toUpperCase() }}
                                        </div>
                                    </template>
                                    <div class="text-sm font-bold text-gray-800">{{ user.username }}</div>
                                </div>
                            </td>
                            <td class="px-6 py-4">
                                <span v-if="user.role === 'Admin'" class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-purple-50 text-purple-600 ring-1 ring-inset ring-purple-500/10">管理员</span>
                                <span v-else class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-gray-50 text-gray-600 ring-1 ring-inset ring-gray-500/10">普通用户</span>
                            </td>
                            <td class="px-6 py-4">
                                <span v-if="user.status === 1" class="inline-flex items-center gap-1.5 px-2 py-1 rounded-full text-xs font-medium bg-green-50 text-green-600">
                                    <div class="w-1.5 h-1.5 rounded-full bg-green-500"></div> 正常
                                </span>
                                <span v-else class="inline-flex items-center gap-1.5 px-2 py-1 rounded-full text-xs font-medium bg-red-50 text-red-600">
                                    <div class="w-1.5 h-1.5 rounded-full bg-red-500"></div> 禁用
                                </span>
                            </td>
                            <td class="px-6 py-4 text-sm text-gray-500 font-mono">{{ formatDate(user.createdAt) }}</td>
                            <td class="px-6 py-4 text-right">
                                <div class="flex items-center justify-end gap-2">
                                    <button @click="handleEdit(user)" :disabled="isProtectedUser(user)"
                                        class="w-8 h-8 rounded-lg flex items-center justify-center transition-colors disabled:cursor-not-allowed"
                                        :class="isProtectedUser(user) ? 'text-gray-300 bg-gray-50' : 'text-gray-500 hover:bg-blue-50 hover:text-blue-600'"
                                        :title="isProtectedUser(user) ? protectedUserTitle(user) : '编辑'">
                                        <div :class="isProtectedUser(user) ? 'i-fa6-solid-lock' : 'i-fa6-solid-pen'"></div>
                                    </button>
                                    <button @click="toggleBan(user)" :disabled="isProtectedUser(user)"
                                        class="w-8 h-8 rounded-lg flex items-center justify-center transition-colors disabled:cursor-not-allowed"
                                        :class="isProtectedUser(user)
                                            ? 'text-gray-300 bg-gray-50'
                                            : (user.status === 1 ? 'text-gray-500 hover:bg-red-50 hover:text-red-500' : 'text-gray-500 hover:bg-green-50 hover:text-green-500')"
                                        :title="isProtectedUser(user) ? protectedUserTitle(user) : (user.status === 1 ? '封禁' : '解封')">
                                        <div :class="isProtectedUser(user) ? 'i-fa6-solid-lock' : (user.status === 1 ? 'i-fa6-solid-ban' : 'i-fa6-solid-check')"></div>
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
                
                <!-- 分页 -->
                <div class="mt-auto">
                    <Pagination v-if="total > 0"
                    :current="currentPage" 
                    :total="total" 
                    :page-size="pageSize"
                    @update:current="handlePageChange"
                />
                </div>
            </div>
        </div>

        <!-- 用户对话框 -->
        <UserModal v-model="showModal" :user="currentUser" :loading="modalLoading" @submit="handleSubmit" />

        <!-- 确认对话框 -->
        <ConfirmModal v-model="showConfirmModal" :type="confirmType" :title="confirmTitle" :content="confirmContent" 
            :loading="confirmLoading" @confirm="handleConfirmAction" />
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import Pagination from '@/components/ui/Pagination.vue'
import UserModal from './UserModal.vue'
import Dropdown from '@/components/ui/Dropdown.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import { getUsers, updateUser, createUser } from '@/api/admin'
import { resolveUrl, showToast } from '@/utils/common'
import { getApiErrorMessage } from '@/utils/request'
import { useUserStore } from '@/stores/user'


interface User {
    id: number
    username: string
    role: string
    status: number
    avatarUrl?: string
    createdAt: string
}

const users = ref<User[]>([])
const userStore = useUserStore()
const currentPage = ref(1)
const pageSize = ref(6)
const total = ref(0)
const loading = ref(false)
const searchKeyword = ref('')
const filterRole = ref('')

// 获取头像 URL 的工具函数
const getAvatarUrl = (path: string) => {
    return resolveUrl(path)
}


// 下拉菜单辅助逻辑状态管理
const filterRoleLabel = computed(() => {
    switch (filterRole.value) {
        case 'User': return '普通用户'
        case 'Admin': return '管理员'
        default: return '所有角色'
    }
})

const currentAdminID = computed(() => Number(userStore.user?.id || 0))

const isDefaultAdminUser = (user: User) => {
    return user.id === 1 || user.username.trim().toLowerCase() === 'admin'
}

const isCurrentAdminUser = (user: User) => {
    return currentAdminID.value > 0 && user.id === currentAdminID.value
}

const isProtectedUser = (user: User) => {
    return isDefaultAdminUser(user) || isCurrentAdminUser(user)
}

const protectedUserTitle = (user: User) => {
    if (isDefaultAdminUser(user)) return '默认管理员不可在后台用户管理中操作'
    if (isCurrentAdminUser(user)) return '不能在后台用户管理中操作自己'
    return '受保护用户'
}

const setRole = (role: string) => {
    filterRole.value = role
    handleSearch()
}

const showModal = ref(false)
const currentUser = ref<User | null>(null)
const modalLoading = ref(false)

// 确认对话框状态变量
const showConfirmModal = ref(false)
const confirmLoading = ref(false)
const confirmType = ref<'info' | 'danger'>('danger')
const confirmTitle = ref('')
const confirmContent = ref('')
const confirmAction = ref<(() => Promise<void>) | null>(null)

// 格式化日期
const formatDate = (dateStr: string) => {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    return date.toLocaleDateString()
}

// 获取用户列表
const fetchUsers = async () => {
    loading.value = true
    try {
        const res: any = await getUsers({
            pageIndex: currentPage.value,
            pageSize: pageSize.value,
            keyword: searchKeyword.value,
            role: filterRole.value
        })
        if (res.data) {
            users.value = res.data.users || []
            total.value = res.data.total || 0
        }
    } catch (error) {
        console.error("Fetch users err:", error)
        showToast('获取用户列表失败', 'error')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => {
    currentPage.value = 1
    fetchUsers()
}

const handlePageChange = (page: number) => {
    currentPage.value = page
    fetchUsers()
}

// 切换禁用状态
const toggleBan = (user: User) => {
    if (isProtectedUser(user)) {
        showToast(protectedUserTitle(user), 'warning')
        return
    }

    const newStatus = user.status === 1 ? 0 : 1
    const actionName = user.status === 1 ? '封禁' : '解封'
    
    confirmTitle.value = `确认${actionName}`
    confirmContent.value = `确定要${actionName}用户 "${user.username}" 吗？`
    confirmType.value = user.status === 1 ? 'danger' : 'info'
    
    confirmAction.value = async () => {
        try {
            await updateUser(user.id, { status: newStatus })
            showToast(`${actionName}成功`, 'success')
            user.status = newStatus
            showConfirmModal.value = false
        } catch (error) {
            showToast(`${actionName}失败`, 'error')
        }
    }
    
    showConfirmModal.value = true
}

const handleConfirmAction = async () => {
    if (confirmAction.value) {
        confirmLoading.value = true
        try {
            await confirmAction.value()
        } finally {
            confirmLoading.value = false
        }
    }
}

// 弹窗对话框处理逻辑
const handleCreate = () => {
    currentUser.value = null
    showModal.value = true
}

const handleEdit = (user: User) => {
    if (isProtectedUser(user)) {
        showToast(protectedUserTitle(user), 'warning')
        return
    }

    currentUser.value = { ...user } // 避免直接修改原始对象
    showModal.value = true
}

const handleSubmit = async (formData: any) => {
    modalLoading.value = true
    try {
        let res: any
        if (currentUser.value) {
            // 更新用户信息
            const updates: any = {
                role: formData.role,
                status: formData.status
            }
            // 仅在用户名发生变化时发送（虽然后端也会校验唯一性）
            if (formData.username !== currentUser.value.username) {
                updates.username = formData.username
            }
            // 仅当提供了密码时才发送密码更新字段
            if (formData.password) {
                updates.password = formData.password
            }

            res = await updateUser(currentUser.value.id, updates)
            if (res.code !== 200) throw new Error(res.msg || '更新失败')
            showToast('更新成功', 'success')
        } else {
            // 创建新用户
            if (!formData.password) {
                 showToast('创建用户必须填写密码', 'error')
                 modalLoading.value = false
                 return
            }
            res = await createUser(formData)
            if (res.code !== 200 && res.code !== 201) throw new Error(res.msg || '创建失败')
            showToast('创建成功', 'success')
        }
        showModal.value = false
        fetchUsers() // 刷新列表数据
    } catch (error) {
        let msg = getApiErrorMessage(error, '操作失败')
        
        // 前端友好错误映射，防止后端未及时更新报错信息时的兜底提示
        if (msg === '用户名已存在') {
            msg = '用户已存在，请换用户名'
        }
        
        showToast(msg, 'error')
        // 对话框保持打开状态，因为在 catch 块中我们没有显式将其设为 false，允许用户修正后再次提交
    } finally {
        modalLoading.value = false
    }
}

onMounted(() => {
    fetchUsers()
})
</script>
