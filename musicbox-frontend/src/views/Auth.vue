<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { login, register } from '@/api/user'
import { showToast } from '@/utils/common'
import { getApiErrorMessage } from '@/utils/request'

interface LoginForm { username: string; password: string }
interface RegisterForm { username: string; password: string; confirmPassword: string }

const isRegister = ref(false)
const loading = ref(false)
const router = useRouter()

const loginForm = reactive<LoginForm>({ username: '', password: '' })
const registerForm = reactive<RegisterForm>({ username: '', password: '', confirmPassword: '' })

const toggleMode = () => {
  isRegister.value = !isRegister.value
  loginForm.username = ''
  loginForm.password = ''
  registerForm.username = ''
  registerForm.password = ''
  registerForm.confirmPassword = ''
}

const handleLogin = async () => {
  if (!loginForm.username || !loginForm.password) {
    showToast('请输入账号和密码', 'error')
    return
  }
  loading.value = true
  try {
    const res: any = await login({ username: loginForm.username, password: loginForm.password })
    if (res.data?.token) {
      localStorage.setItem('token', res.data.token)
      showToast('欢迎回来！正在进入音乐世界...', 'success')
      setTimeout(() => router.push('/'), 1500)
    } else {
      showToast(res.msg || '登录失败，响应格式不正确', 'error')
    }
  } catch (error: any) {
    const msg = getApiErrorMessage(error, '登录失败，请检查您的凭据')
    showToast(msg.includes('账号已被禁用') ? '您的账号已被封禁，请联系管理员' : msg, 'error')
  } finally {
    loading.value = false
  }
}

const handleRegister = async () => {
  const username = registerForm.username.trim()
  const password = registerForm.password
  if (!username || !password || !registerForm.confirmPassword) return showToast('请完善注册信息', 'error')
  if (username.length < 3) return showToast('用户名至少需要 3 个字符', 'error')
  if (username.length > 30) return showToast('用户名长度不能超过 30 个字符', 'error')
  if (password !== registerForm.confirmPassword) return showToast('两次密码输入不一致', 'error')
  if (password.length < 6) return showToast('密码长度至少需要 6 位', 'error')
  loading.value = true
  try {
    const res: any = await register({ username, password })
    if (res.code !== 201) throw new Error(res.msg || '注册失败')
    showToast(res.msg || '注册成功！请登录', 'success')
    setTimeout(() => {
      toggleMode()
      loginForm.username = registerForm.username
      loginForm.password = ''
    }, 1500)
  } catch (error: any) {
    showToast(getApiErrorMessage(error, '注册失败，用户名可能已被占用'), 'error')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-root">

    <!-- 主卡片 -->
    <div class="auth-card" :class="{ 'is-register': isRegister }">

      <!-- ── 左侧：视觉面板 ── -->
      <div class="visual-panel">
        <!-- 背景网格纹理 -->
        <div class="grid-texture" />

        <!-- 动态光晕 -->
        <div class="glow glow-1" />
        <div class="glow glow-2" />

        <!-- 内容 -->
        <div class="visual-content" :class="{ shifted: isRegister }">

          <!-- 旋转黑胶 -->
          <div class="vinyl-wrap">
            <div class="vinyl-glow" />
            <div class="vinyl">
              <div class="vinyl-grooves" />
              <div class="vinyl-label">
                <div class="vinyl-hole" />
              </div>
            </div>
            <!-- 唱针 -->
            <div class="tonearm" :class="{ lifted: isRegister }" />
          </div>

          <!-- 音乐律动条 -->
          <div class="bars">
            <span v-for="i in 7" :key="i" class="bar" :style="`--d:${(i - 1) * 80}ms`" />
          </div>

          <!-- 文字切换 -->
          <Transition name="slide-text" mode="out-in">
            <div v-if="!isRegister" class="visual-text">
              <h2>欢迎回来</h2>
              <p>登录以继续播放您的收藏</p>
              <button class="panel-btn" @click="toggleMode">没有账号？去注册</button>
            </div>
            <div v-else class="visual-text">
              <h2>开启旅程</h2>
              <p>注册即享高品质无损音乐体验</p>
              <button class="panel-btn" @click="toggleMode">已有账号？去登录</button>
            </div>
          </Transition>
        </div>
      </div>

      <!-- ── 右侧：表单区 ── -->
      <div class="form-panel">

        <!-- 登录表单 -->
        <Transition name="form-swap" mode="out-in">
          <form v-if="!isRegister" key="login" @submit.prevent="handleLogin" class="auth-form">
            <div class="form-header">
              <div class="logo-dot" />
              <span class="brand">MusicBox</span>
            </div>

            <h1 class="form-title">登录</h1>
            <p class="form-sub">用音乐填满每一刻</p>

            <div class="fields">
              <label class="field">
                <div class="field-icon i-mdi-account-music-outline" />
                <input v-model="loginForm.username" type="text" placeholder="用户名" autocomplete="username" />
              </label>
              <label class="field">
                <div class="field-icon i-mdi-lock-outline" />
                <input v-model="loginForm.password" type="password" placeholder="密码" autocomplete="current-password" />
              </label>
            </div>

            <button type="submit" class="submit-btn" :disabled="loading">
              <div v-if="loading" class="i-svg-spinners-bars-fade" />
              <span>{{ loading ? '登录中...' : '登 录' }}</span>
            </button>

            <div class="mobile-switch">
              没有账号？<button type="button" @click="toggleMode">去注册</button>
            </div>
          </form>

          <!-- 注册表单 -->
          <form v-else key="register" @submit.prevent="handleRegister" class="auth-form">
            <div class="form-header">
              <div class="logo-dot" />
              <span class="brand">MusicBox</span>
            </div>

            <h1 class="form-title">注册</h1>
            <p class="form-sub">加入，定制您的专属频率</p>

            <div class="fields">
              <label class="field">
                <div class="field-icon i-mdi-account-outline" />
                <input v-model="registerForm.username" type="text" placeholder="设置用户名" autocomplete="username" />
              </label>
              <label class="field">
                <div class="field-icon i-mdi-lock-outline" />
                <input v-model="registerForm.password" type="password" placeholder="设置密码" autocomplete="new-password" />
              </label>
              <label class="field">
                <div class="field-icon i-mdi-lock-check-outline" />
                <input v-model="registerForm.confirmPassword" type="password" placeholder="确认密码"
                  autocomplete="new-password" />
              </label>
            </div>

            <button type="submit" class="submit-btn" :disabled="loading">
              <div v-if="loading" class="i-svg-spinners-bars-fade" />
              <span>{{ loading ? '注册中...' : '立即注册' }}</span>
            </button>

            <div class="mobile-switch">
              已有账号？<button type="button" @click="toggleMode">去登录</button>
            </div>
          </form>
        </Transition>

      </div>
    </div>
  </div>
</template>

<style scoped>
/* ── 根布局 ── */
.auth-root {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100svh;
  background: #f0f4ff;
  background-image:
    radial-gradient(ellipse 80% 60% at 10% 0%, theme('colors.primary')14 0%, transparent 60%),
    radial-gradient(ellipse 50% 50% at 90% 100%, theme('colors.primary')0f 0%, transparent 60%);
  padding: 1rem;
  font-family: 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
}

/* ── 主卡片 ── */
.auth-card {
  display: flex;
  width: 100%;
  max-width: 960px;
  min-height: 580px;
  background: #fff;
  border-radius: 28px;
  overflow: hidden;
  box-shadow:
    0 0 0 1px theme('colors.primary')14,
    0 24px 64px rgba(13, 17, 23, 0.1),
    0 4px 16px theme('colors.primary')0f;
}

/* ── 视觉面板（左） ── */
.visual-panel {
  position: relative;
  width: 44%;
  background: #0d1117;
  overflow: hidden;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 网格纹理 */
.grid-texture {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(theme('colors.primary')1a 1px, transparent 1px),
    linear-gradient(90deg, theme('colors.primary')1a 1px, transparent 1px);
  background-size: 32px 32px;
}

/* 光晕 */
.glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  pointer-events: none;
}

.glow-1 {
  width: 320px;
  height: 320px;
  background: theme('colors.primary')47;
  top: -80px;
  left: -60px;
}

.glow-2 {
  width: 200px;
  height: 200px;
  background: rgba(107, 143, 245, 0.2);
  bottom: -40px;
  right: -40px;
}

/* 视觉内容 */
.visual-content {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  padding: 40px 32px;
  text-align: center;
  transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 黑胶唱片 */
.vinyl-wrap {
  position: relative;
  width: 160px;
  height: 160px;
}

.vinyl-glow {
  position: absolute;
  inset: -12px;
  background: radial-gradient(circle, theme('colors.primary')59 0%, transparent 70%);
  border-radius: 50%;
  animation: pulse-glow 3s ease-in-out infinite;
}

.vinyl {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: radial-gradient(circle at 50% 50%, #1e2535 0%, #0d1117 100%);
  box-shadow:
    0 0 0 2px theme('colors.primary')40,
    inset 0 0 30px rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  animation: spin 12s linear infinite;
  position: relative;
}

.vinyl-grooves {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  background: repeating-radial-gradient(circle at 50%,
      transparent 0px,
      transparent 6px,
      rgba(255, 255, 255, 0.025) 6px,
      rgba(255, 255, 255, 0.025) 7px);
}

.vinyl-label {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: linear-gradient(135deg, theme('colors.primary') 0%, #5b7ef5 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 12px theme('colors.primary')99;
  z-index: 1;
}

.vinyl-hole {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #0d1117;
}

/* 唱针 */
.tonearm {
  position: absolute;
  top: -10px;
  right: -22px;
  width: 52px;
  height: 4px;
  background: linear-gradient(90deg, rgba(255, 255, 255, 0.6), rgba(255, 255, 255, 0.2));
  border-radius: 2px;
  transform-origin: right center;
  transform: rotate(-30deg);
  transition: transform 0.8s cubic-bezier(0.4, 0, 0.2, 1);
}

.tonearm.lifted {
  transform: rotate(-60deg);
}

/* 音乐律动条 */
.bars {
  display: flex;
  align-items: flex-end;
  gap: 5px;
  height: 28px;
}

.bar {
  width: 3px;
  background: theme('colors.primary');
  border-radius: 2px;
  animation: bar-dance 1.2s ease-in-out infinite;
  animation-delay: var(--d);
  opacity: 0.7;
}

.bar:nth-child(odd) {
  background: #6b8ff5;
  opacity: 0.9;
}

/* 视觉文字 */
.visual-text h2 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #fff;
  margin: 0 0 6px;
  letter-spacing: -0.01em;
}

.visual-text p {
  font-size: 0.8rem;
  color: rgba(255, 255, 255, 0.5);
  margin: 0 0 20px;
  line-height: 1.6;
}

.panel-btn {
  display: inline-block;
  padding: 8px 22px;
  border-radius: 100px;
  border: 1px solid theme('colors.primary')80;
  background: theme('colors.primary')26;
  color: #6b8ff5;
  font-size: 0.78rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.panel-btn:hover {
  background: rgba(51, 91, 241, 0.3);
  border-color: #6b8ff5;
  color: #fff;
}

/* ── 表单面板（右） ── */
.form-panel {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 40px;
  background: #fff;
}

/* ── 表单 ── */
.auth-form {
  width: 100%;
  max-width: 320px;
  display: flex;
  flex-direction: column;
  gap: 0;
}

.form-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 32px;
}

.logo-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: theme('colors.primary');
  box-shadow: 0 0 0 3px theme('colors.primary')33;
}

.brand {
  font-size: 0.85rem;
  font-weight: 600;
  color: #94a3b8;
  letter-spacing: 0.04em;
}

.form-title {
  font-size: 2rem;
  font-weight: 700;
  color: #0f172a;
  margin: 0 0 6px;
  letter-spacing: -0.03em;
}

.form-sub {
  font-size: 0.84rem;
  color: #94a3b8;
  margin: 0 0 28px;
}

/* 输入组 */
.fields {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 24px;
}

.field {
  position: relative;
  display: flex;
  align-items: center;
}

.field-icon {
  position: absolute;
  left: 14px;
  font-size: 16px;
  color: #cbd5e1;
  transition: color 0.2s;
  pointer-events: none;
}

.field input {
  width: 100%;
  padding: 13px 16px 13px 40px;
  background: #f8fafc;
  border: 1.5px solid #e8edf5;
  border-radius: 12px;
  font-size: 0.875rem;
  color: #0f172a;
  outline: none;
  transition: all 0.2s;
  box-sizing: border-box;
}

.field input::placeholder {
  color: #c0ccda;
}

.field input:focus {
  background: #fff;
  border-color: theme('colors.primary');
  box-shadow: 0 0 0 3px theme('colors.primary')1a;
}

.field:focus-within .field-icon {
  color: theme('colors.primary');
}

/* 提交按钮 */
.submit-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 14px;
  border-radius: 12px;
  background: theme('colors.primary');
  color: #fff;
  font-size: 0.9rem;
  font-weight: 600;
  letter-spacing: 0.06em;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
  box-shadow: 0 6px 20px theme('colors.primary')59;
  position: relative;
  overflow: hidden;
}

.submit-btn::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.12) 0%, transparent 60%);
  pointer-events: none;
}

.submit-btn:hover:not(:disabled) {
  background: theme('colors.primary');
  opacity: 0.9;
  box-shadow: 0 8px 28px theme('colors.primary')73;
  transform: translateY(-1px);
}

.submit-btn:active:not(:disabled) {
  transform: translateY(0);
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* 移动端切换入口 */
.mobile-switch {
  display: none;
  margin-top: 20px;
  text-align: center;
  font-size: 0.82rem;
  color: #94a3b8;
}

.mobile-switch button {
  color: theme('colors.primary');
  font-weight: 600;
  background: none;
  border: none;
  cursor: pointer;
  margin-left: 4px;
}

/* ── 过渡动画 ── */
.form-swap-enter-active,
.form-swap-leave-active {
  transition: all 0.28s cubic-bezier(0.4, 0, 0.2, 1);
}

.form-swap-enter-from {
  opacity: 0;
  transform: translateX(16px);
}

.form-swap-leave-to {
  opacity: 0;
  transform: translateX(-16px);
}

.slide-text-enter-active,
.slide-text-leave-active {
  transition: all 0.3s ease;
}

.slide-text-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.slide-text-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* ── 关键帧 ── */
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@keyframes bar-dance {

  0%,
  100% {
    height: 8px;
    opacity: 0.5;
  }

  50% {
    height: 24px;
    opacity: 1;
  }
}

@keyframes pulse-glow {

  0%,
  100% {
    opacity: 0.6;
    transform: scale(1);
  }

  50% {
    opacity: 1;
    transform: scale(1.08);
  }
}

/* ── 移动端响应式 ── */
@media (max-width: 720px) {
  .auth-card {
    flex-direction: column;
    min-height: unset;
    border-radius: 20px;
  }

  .visual-panel {
    width: 100%;
    padding: 28px 0 20px;
  }

  .visual-content {
    gap: 12px;
    padding: 0 24px;
  }

  .vinyl-wrap {
    width: 110px;
    height: 110px;
  }

  .panel-btn {
    display: none;
  }

  .form-panel {
    padding: 28px 24px 36px;
  }

  .form-header {
    margin-bottom: 20px;
  }

  .form-title {
    font-size: 1.6rem;
  }

  .mobile-switch {
    display: block;
  }
}
</style>
