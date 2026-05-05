<template>
  <div class="p-4 md:p-8 max-w-5xl mx-auto space-y-6 animate-fade-in pb-10">
    <!-- 页面标题 -->
    <div class="flex items-center gap-4 mb-2">
      <div
        class="w-12 h-12 bg-gradient-to-br from-primary to-blue-600 rounded-2xl flex items-center justify-center text-white shadow-lg shadow-primary/20">
        <div class="i-fa6-solid-user-gear text-xl"></div>
      </div>
      <div>
        <h2 class="text-3xl font-black text-gray-800 tracking-tight">个人中心</h2>
        <div class="flex items-center gap-1.5 mt-0.5">
          <span class="w-1.5 h-1.5 bg-green-500 rounded-full animate-pulse"></span>
          <p class="text-gray-400 font-bold uppercase text-[9px] tracking-widest">账号设置与安全</p>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
      <!-- 左侧：基本资料 & 头像 -->
      <section
        class="lg:col-span-5 bg-white/70 backdrop-blur-xl p-6 rounded-[2.5rem] border border-white shadow-xl shadow-gray-200/50 space-y-6 relative overflow-hidden group/card">
        <div
          class="absolute -top-12 -right-12 w-32 h-32 bg-primary/5 rounded-full blur-2xl group-hover/card:bg-primary/10 transition-colors duration-700">
        </div>

        <div class="flex items-center justify-between relative z-10">
          <div class="flex items-center gap-2">
            <div class="p-1.5 bg-primary/10 rounded-lg text-primary">
              <div class="i-fa6-solid-address-card text-xs"></div>
            </div>
            <h3 class="text-lg font-bold text-gray-800">账户资料</h3>
          </div>
        </div>

        <!-- 头像上传区 -->
        <div class="flex flex-col items-center gap-4 relative z-10">
          <div class="relative group">
            <div v-if="previewUrl && !uploadingAvatar"
              class="absolute -inset-3 bg-primary/20 rounded-full animate-ping opacity-20"></div>

            <div
              class="w-32 h-32 rounded-full overflow-hidden border-4 border-white shadow-xl bg-gray-50 ring-4 ring-gray-50 transition-all group-hover:scale-105 duration-500 relative">
              <img v-if="(previewUrl || userStore.avatarUrl) && !avatarImgError"
                :src="previewUrl || userStore.avatarUrl" class="w-full h-full object-cover"
                :class="{ 'blur-[1px] opacity-70': uploadingAvatar }" @error="avatarImgError = true" />
              <div v-else
                class="w-full h-full bg-gradient-to-br from-blue-50 to-indigo-50 flex items-center justify-center text-primary text-4xl font-black">
                {{ (userStore.user?.username || 'U').charAt(0).toUpperCase() }}
              </div>

              <div v-if="uploadingAvatar"
                class="absolute inset-0 bg-black/30 backdrop-blur-sm flex flex-col items-center justify-center text-white">
                <div class="i-fa6-solid-spinner animate-spin text-2xl"></div>
              </div>

              <div v-if="previewUrl && !uploadingAvatar" class="absolute bottom-2 left-0 right-0 flex justify-center">
                <span
                  class="px-2 py-0.5 bg-primary text-white text-[8px] font-black rounded-full shadow-lg">PREVIEW</span>
              </div>
            </div>

            <label v-if="!previewUrl" for="avatar-upload"
              class="absolute inset-0 bg-black/40 flex flex-col items-center justify-center text-white opacity-0 group-hover:opacity-100 transition-all cursor-pointer rounded-full z-10">
              <div
                class="i-fa6-solid-camera text-2xl transform translate-y-1 group-hover:translate-y-0 transition-transform">
              </div>
            </label>

            <div v-if="!previewUrl"
              class="absolute bottom-0 right-0 w-8 h-8 bg-white rounded-full shadow-lg flex items-center justify-center text-primary border-2 border-gray-50 group-hover:bg-primary group-hover:text-white transition-all z-20 pointer-events-none">
              <div class="i-fa6-solid-camera text-xs"></div>
            </div>

            <input id="avatar-upload" type="file" accept="image/*" class="hidden" @change="handleFileSelect"
              :disabled="uploadingAvatar" />
          </div>

          <!-- 预览确认按钮组 -->
          <div v-if="previewUrl"
            class="flex items-center gap-2 animate-slide-up bg-gray-100/50 p-1 rounded-full border border-gray-100">
            <button @click="cancelSelection" :disabled="uploadingAvatar"
              class="px-4 py-1.5 text-gray-500 rounded-full text-xs font-bold hover:bg-white hover:text-gray-800 transition-all disabled:opacity-50">取消</button>
            <button @click="confirmUpload" :disabled="uploadingAvatar"
              class="px-5 py-1.5 bg-primary text-white rounded-full text-xs font-bold shadow-lg shadow-primary/20 hover:shadow-primary/30 transition-all disabled:opacity-50 flex items-center gap-1.5">
              <div v-if="uploadingAvatar" class="i-fa6-solid-spinner animate-spin"></div>
              保存头像
            </button>
          </div>

          <div v-else class="text-center">
            <div class="text-2xl font-black text-gray-800 tracking-tight">{{ userStore.user?.username }}</div>
            <div
              class="inline-flex items-center gap-1 px-2 py-0.5 bg-primary/5 text-primary text-[9px] font-black rounded-full uppercase mt-1 border border-primary/10">
              <div v-if="userStore.isAdmin" class="i-fa6-solid-crown text-[7px]"></div>
              {{ userStore.isAdmin ? 'Admin' : 'User' }}
            </div>
          </div>

          <!-- 详细信息网格 -->
          <div class="w-full pt-4 grid grid-cols-1 gap-2">
            <div
              class="flex items-center justify-between p-3 bg-gray-50/50 rounded-xl border border-gray-100 hover:bg-white hover:border-primary/20 transition-all group/info">
              <div class="flex items-center gap-3">
                <div
                  class="w-7 h-7 rounded-lg bg-gray-100 flex items-center justify-center text-gray-400 group-hover/info:text-primary transition-colors">
                  <div class="i-fa6-solid-fingerprint text-[10px]"></div>
                </div>
                <div class="text-[11px] font-mono font-bold text-gray-700">UID: {{ userStore.user?.id }}</div>
              </div>
              <button @click="copyUid" class="p-1.5 text-gray-300 hover:text-primary transition-colors">
                <div :class="uidCopied ? 'i-fa6-solid-check text-green-500' : 'i-fa6-solid-copy'"></div>
              </button>
            </div>

            <div
              class="flex items-center gap-3 p-3 bg-gray-50/50 rounded-xl border border-gray-100 hover:bg-white hover:border-primary/20 transition-all group/info">
              <div
                class="w-7 h-7 rounded-lg bg-gray-100 flex items-center justify-center text-gray-400 group-hover/info:text-primary transition-colors">
                <div class="i-fa6-solid-calendar-days text-[10px]"></div>
              </div>
              <div class="text-[11px] font-black text-gray-700">注册于 {{ formatDate(userStore.user?.createdAt || '') }}
              </div>
            </div>

            <button @click="handleLogout"
              class="flex items-center justify-between p-3 bg-rose-50/30 rounded-xl border border-rose-100/50 hover:bg-rose-50 hover:border-rose-200 transition-all group/logout">
              <div class="flex items-center gap-3">
                <div
                  class="w-7 h-7 rounded-lg bg-rose-100 text-rose-500 flex items-center justify-center group-hover/logout:bg-rose-500 group-hover/logout:text-white transition-all">
                  <div class="i-fa6-solid-right-from-bracket text-[10px]"></div>
                </div>
                <div class="text-[11px] font-black text-rose-600">退出当前账号</div>
              </div>
              <div
                class="i-fa6-solid-chevron-right text-[8px] text-rose-300 group-hover/logout:translate-x-0.5 transition-transform">
              </div>
            </button>
          </div>
        </div>
      </section>

      <!-- 右侧：安全设置 -->
      <section class="lg:col-span-7">
        <div
          class="bg-white/70 backdrop-blur-xl p-6 rounded-[2.5rem] border border-white shadow-xl shadow-gray-200/50 space-y-6">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <div class="p-1.5 bg-rose-50 rounded-lg text-rose-500">
                <div class="i-fa6-solid-shield-halved text-xs"></div>
              </div>
              <h3 class="text-lg font-bold text-gray-800">安全修改</h3>
            </div>
            <span
              class="text-[8px] font-black text-gray-400 border border-gray-100 px-1.5 py-0.5 rounded-md uppercase tracking-widest">Security</span>
          </div>

          <form @submit.prevent="handleUpdatePassword" class="space-y-4">
            <div class="space-y-1.5">
              <label class="text-[9px] font-black text-gray-400 uppercase tracking-widest px-2">旧密码验证</label>
              <div class="relative group">
                <input v-model="passwordForm.oldPassword" type="password" placeholder="当前密码" required
                  class="w-full pl-10 pr-4 py-3.5 rounded-xl bg-gray-50 border border-gray-100 focus:border-primary focus:bg-white focus:ring-4 focus:ring-primary/5 transition-all text-sm outline-none" />
                <div
                  class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400 group-focus-within:text-primary transition-colors">
                  <div class="i-fa6-solid-lock text-[10px]"></div>
                </div>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="space-y-1.5">
                <label class="text-[9px] font-black text-gray-400 uppercase tracking-widest px-2">新登录密码</label>
                <div class="relative group">
                  <input v-model="passwordForm.newPassword" type="password" placeholder="6位以上" required
                    class="w-full pl-10 pr-4 py-3.5 rounded-xl bg-gray-50 border border-gray-100 focus:border-primary focus:bg-white focus:ring-4 focus:ring-primary/5 transition-all text-sm outline-none" />
                  <div
                    class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400 group-focus-within:text-primary transition-colors">
                    <div class="i-fa6-solid-key text-[10px]"></div>
                  </div>
                </div>
              </div>

              <div class="space-y-1.5">
                <label class="text-[9px] font-black text-gray-400 uppercase tracking-widest px-2">再次确认</label>
                <div class="relative group">
                  <input v-model="passwordForm.confirmPassword" type="password" placeholder="重复输入" required
                    :class="['w-full pl-10 pr-10 py-3.5 rounded-xl bg-gray-50 border transition-all text-sm outline-none',
                      passwordMismatch ? 'border-rose-200 focus:ring-4 focus:ring-rose-50 focus:border-rose-400' : 'border-gray-100 focus:border-primary focus:bg-white focus:ring-4 focus:ring-primary/5']" />
                  <div class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">
                    <div class="i-fa6-solid-check-double text-[10px]"></div>
                  </div>
                  <div v-if="passwordMatch" class="absolute right-4 top-1/2 -translate-y-1/2 text-green-500">
                    <div class="i-fa6-solid-circle-check text-[10px]"></div>
                  </div>
                </div>
              </div>
            </div>

            <!--
              固定高度占位区：始终占据 h-5 的空间，错误文字用 v-show + transition 控制显隐，
            -->
            <div class="h-5 flex items-center px-2">
              <Transition enter-active-class="transition duration-200 ease-out"
                enter-from-class="opacity-0 -translate-y-1" enter-to-class="opacity-100 translate-y-0"
                leave-active-class="transition duration-150 ease-in" leave-from-class="opacity-100 translate-y-0"
                leave-to-class="opacity-0 -translate-y-1">
                <span v-if="passwordMismatch"
                  class="text-[9px] font-bold text-rose-500 uppercase flex items-center gap-1">
                  <div class="i-fa6-solid-triangle-exclamation"></div> 密码不一致
                </span>
              </Transition>
            </div>

            <button type="submit" :disabled="submitting || passwordMismatch"
              class="w-full py-4 bg-gray-900 text-white rounded-xl font-black shadow-lg shadow-gray-900/10 hover:shadow-gray-900/20 active:scale-95 transition-all disabled:opacity-30 flex items-center justify-center gap-2 relative group/btn overflow-hidden">
              <span>{{ submitting ? '提交中...' : '同步新密码' }}</span>
              <div v-if="submitting" class="i-fa6-solid-spinner animate-spin"></div>
              <div
                class="absolute top-0 -left-1/2 w-full h-full bg-white/5 skew-x-[-25deg] group-hover/btn:left-[150%] transition-all duration-700">
              </div>
            </button>
          </form>
        </div>

        <div class="mt-6 px-4 py-3 bg-gray-50/50 rounded-2xl border border-gray-100/50">
          <div class="flex gap-3">
            <div class="text-primary mt-0.5">
              <div class="i-fa6-solid-circle-info text-xs"></div>
            </div>
            <p class="text-[10px] text-gray-400 font-medium leading-relaxed">
              定期更换密码能有效保护您的账号安全。如果您发现头像更新不及时，可以尝试刷新页面或重新登录。
            </p>
          </div>
        </div>
      </section>
    </div>

    <ConfirmModal v-model="showLogoutConfirm" title="确认退出" content="确定要退出登录吗？" type="danger" confirm-text="确认退出"
      @confirm="confirmLogout" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { uploadAvatar, updatePassword } from '@/api/user'
import { showToast, formatDate } from '@/utils/common'
import { getApiErrorMessage } from '@/utils/request'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'

const router = useRouter()
const userStore = useUserStore()

watch(() => userStore.avatarUrl, () => { avatarImgError.value = false })

const submitting = ref(false)
const uploadingAvatar = ref(false)
const uidCopied = ref(false)
const avatarImgError = ref(false)
const showLogoutConfirm = ref(false)
const previewUrl = ref('')
const selectedFile = ref<File | null>(null)

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const passwordMismatch = computed(() =>
  !!(passwordForm.confirmPassword && passwordForm.newPassword !== passwordForm.confirmPassword)
)
const passwordMatch = computed(() =>
  !!(passwordForm.confirmPassword && passwordForm.newPassword === passwordForm.confirmPassword)
)

const handleFileSelect = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return
  if (file.size > 2 * 1024 * 1024) {
    showToast('头像文件不能超过 2MB', 'error')
    return
  }
  if (previewUrl.value) URL.revokeObjectURL(previewUrl.value)
  selectedFile.value = file
  previewUrl.value = URL.createObjectURL(file)
    ; (event.target as HTMLInputElement).value = ''
}

const cancelSelection = () => {
  if (previewUrl.value) URL.revokeObjectURL(previewUrl.value)
  previewUrl.value = ''
  selectedFile.value = null
}

const confirmUpload = async () => {
  if (!selectedFile.value) return
  const formData = new FormData()
  formData.append('file', selectedFile.value)
  uploadingAvatar.value = true
  try {
    const res = await uploadAvatar(formData)
    if (res.code === 200) {
      showToast('头像更新成功')
      cancelSelection()
      await userStore.fetchProfile()
    } else {
      showToast(res.msg || '上传失败', 'error')
    }
  } catch (error) {
    showToast(getApiErrorMessage(error, '上传失败，请稍后重试'), 'error')
  } finally {
    uploadingAvatar.value = false
  }
}

const handleUpdatePassword = async () => {
  if (passwordMismatch.value) return
  if (passwordForm.newPassword.length < 6) {
    showToast('新密码至少为 6 位', 'error')
    return
  }
  submitting.value = true
  try {
    const res = await updatePassword({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword,
    })
    if (res.code === 200) {
      showToast('密码修改成功')
      passwordForm.oldPassword = ''
      passwordForm.newPassword = ''
      passwordForm.confirmPassword = ''
    } else {
      showToast(res.msg || '修改失败', 'error')
    }
  } catch (error) {
    showToast(getApiErrorMessage(error, '修改失败，原始密码可能错误'), 'error')
  } finally {
    submitting.value = false
  }
}

const handleLogout = () => { showLogoutConfirm.value = true }

const confirmLogout = () => {
  userStore.logout()
  showLogoutConfirm.value = false
  showToast('您已安全退出登录')
  router.push('/auth')
}

const copyUid = () => {
  if (!userStore.user?.id) return
  navigator.clipboard.writeText(String(userStore.user.id)).then(() => {
    uidCopied.value = true
    showToast('UID 已复制')
    setTimeout(() => { uidCopied.value = false }, 2000)
  })
}

onMounted(() => { if (!userStore.user) userStore.fetchProfile() })
onBeforeUnmount(() => { if (previewUrl.value) URL.revokeObjectURL(previewUrl.value) })
</script>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.8s cubic-bezier(0.16, 1, 0.3, 1);
}

.animate-slide-up {
  animation: slideUp 0.5s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(10px) scale(0.95);
  }

  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}
</style>
