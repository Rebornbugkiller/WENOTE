<template>
  <div class="min-h-screen bg-green-50 flex items-center justify-center p-4 relative overflow-hidden select-none">
    <!-- CRT 扫描线 -->
    <div class="fixed inset-0 z-50 pointer-events-none opacity-5" style="background: linear-gradient(rgba(18,16,16,0) 50%, rgba(0,0,0,0.25) 50%), linear-gradient(90deg, rgba(255,0,0,0.06), rgba(0,255,0,0.02), rgba(0,0,255,0.06)); background-size: 100% 2px, 3px 100%;"></div>

    <!-- 网格背景 -->
    <div class="absolute inset-0 pointer-events-none opacity-20" style="background-image: linear-gradient(#22c55e 1px, transparent 1px), linear-gradient(90deg, #22c55e 1px, transparent 1px); background-size: 40px 40px;"></div>

    <!-- 浮动图标 -->
    <FloatingIcons />

    <!-- HUD 面板 -->
    <HudPanel />

    <!-- 音乐控制器 -->
    <MusicWidget :is-playing="isPlayingMusic" @toggle="toggleMusic" />

    <!-- Combo 显示 -->
    <Transition name="combo">
      <div v-if="combo > 0 && !showSettlement" class="fixed top-1/4 left-1/2 -translate-x-1/2 z-50 pointer-events-none">
        <div class="font-black text-6xl italic text-transparent bg-clip-text bg-gradient-to-b from-yellow-300 to-red-500 drop-shadow-[4px_4px_0_#000]" :style="{ transform: `scale(${1 + Math.min(combo * 0.1, 0.5)})` }">
          {{ combo }} COMBO!
        </div>
        <div class="text-center font-bold text-white text-sm bg-black px-2 rounded-full mt-2 border-2 border-white animate-bounce">SUPER HOT!!</div>
      </div>
    </Transition>

    <!-- Slogan -->
    <div class="absolute top-8 left-1/2 -translate-x-1/2 z-30">
      <h1 v-if="locale === 'en-US'" class="slogan-text text-4xl md:text-5xl font-black tracking-wider">
        <span class="slogan-word text-slate-800" style="--i:0">We</span>
        <span class="slogan-word text-green-500" style="--i:1">Note</span>
        <span class="slogan-word text-slate-800" style="--i:2">,</span>
        <span class="slogan-word text-slate-800" style="--i:3">We</span>
        <span class="slogan-word text-green-500" style="--i:4">Create</span>
        <span class="slogan-word text-slate-800" style="--i:5">.</span>
      </h1>
      <h1 v-else class="slogan-text text-4xl md:text-5xl font-black tracking-wider">
        <span class="slogan-word text-green-500" style="--i:0">记录</span>
        <span class="slogan-word text-slate-800" style="--i:1">此刻</span>
        <span class="slogan-word text-slate-800" style="--i:2">，</span>
        <span class="slogan-word text-green-500" style="--i:3">创造</span>
        <span class="slogan-word text-slate-800" style="--i:4">未来</span>
      </h1>
    </div>

    <!-- 主卡片容器 -->
    <div class="relative z-20 w-full max-w-7xl mx-auto flex items-center justify-center gap-8 lg:gap-16">
      
      <!-- Left Wing: AutoSnake Console -->
      <div class="hidden lg:block transform -rotate-6 hover:rotate-0 transition-transform duration-500 hover:scale-105">
        <AutoSnake />
      </div>

      <!-- Center Stage: Login Card -->
      <div class="w-full max-w-md relative group">
        <!-- Language & Fever Mode Toggle -->
        <div class="absolute -top-16 right-0 flex gap-2 z-30">
          <button
            @click="toggleLocale"
            class="bg-blue-600 text-white px-4 py-2 rounded-lg font-black transform hover:scale-110 active:scale-95 transition-all shadow-[4px_4px_0_#000] border-2 border-black"
          >
            {{ locale === 'zh-CN' ? 'EN' : '中' }}
          </button>
          <button
            @click="toggleFeverMode"
            class="bg-red-600 text-white px-4 py-2 rounded-lg font-black italic transform hover:scale-110 active:scale-95 transition-all shadow-[4px_4px_0_#000] border-2 border-black"
            :class="{ 'animate-pulse': isFeverMode }"
          >
            {{ isFeverMode ? '🔥 ' + t('login.feverOn') : '🚨 ' + t('login.feverMode') }}
          </button>
        </div>

        <div class="bg-white border-4 border-black rounded-[2rem] shadow-[16px_16px_0px_0px_rgba(0,0,0,1)] overflow-visible transition-all duration-100"
          :class="{ 'animate-shake border-red-500 shadow-red-500': isFeverMode }"
        >
          <!-- 卡片头部 -->
          <div class="h-6 bg-gradient-to-r from-green-400 to-green-500 border-b-4 border-black flex items-center justify-between px-4 rounded-t-[1.8rem]"
            :class="{ 'from-red-500 to-yellow-500': isFeverMode }"
          >
            <div class="flex gap-1">
              <div v-for="i in 4" :key="i" class="w-1.5 h-3 bg-black/20 rounded-full"></div>
            </div>
            <div class="text-[10px] font-black opacity-60 tracking-[0.2em]">{{ isLogin ? t('login.insertCoin') : t('login.newChallenger') }}</div>
            <div class="flex gap-1">
              <div v-for="i in 4" :key="i" class="w-1.5 h-3 bg-black/20 rounded-full"></div>
            </div>
          </div>
    
          <div class="p-8">
            <!-- 像素头像 -->
            <PixelAvatar
              :focus-field="focusField"
              :combo="combo"
              :is-login="isLogin"
              :message="avatarMessage"
              @click="handleAvatarClick"
            />
    
            <!-- 标题 & 模式切换 -->
            <div class="text-center mb-6">
              <div class="flex justify-center space-x-4 mb-2">
                <button
                  @click="isLogin && toggleMode()"
                  @mouseenter="playSFX('hover')"
                  class="text-xs font-black uppercase tracking-widest border-b-4 pb-1 transition-all"
                  :class="!isLogin ? 'text-black border-green-500' : 'text-slate-300 border-transparent hover:text-slate-500'"
                >{{ t('login.createHero') }}</button>
                <button
                  @click="!isLogin && toggleMode()"
                  @mouseenter="playSFX('hover')"
                  class="text-xs font-black uppercase tracking-widest border-b-4 pb-1 transition-all"
                  :class="isLogin ? 'text-black border-green-500' : 'text-slate-300 border-transparent hover:text-slate-500'"
                >{{ t('login.loadSave') }}</button>
              </div>
              <h1 class="text-4xl font-black text-slate-800 tracking-tighter drop-shadow-[2px_2px_0_rgba(0,0,0,1)]">
                {{ isLogin ? 'WE' : 'JOIN ' }}<span class="text-green-500" :class="{ 'text-red-500 animate-pulse': isFeverMode }">NOTE</span>
              </h1>
              <p class="text-xs font-bold text-slate-400 tracking-widest mt-1">{{ t('login.slogan') }}</p>
            </div>
    
            <!-- 表单 -->
            <form @submit.prevent="handleSubmit" class="space-y-4">
              <!-- 用户名 -->
              <div class="group relative">
                <label class="block text-xs font-black uppercase tracking-wider mb-1 ml-1">{{ isLogin ? t('login.playerId') : t('login.heroName') }}</label>
                <div class="relative transition-transform group-focus-within:-translate-y-1 group-focus-within:scale-[1.02]">
                  <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <svg class="w-6 h-6 text-black" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/></svg>
                  </div>
                  <input
                    v-model="form.username"
                    type="text"
                    class="block w-full pl-12 pr-4 py-3 border-4 border-black rounded-xl bg-slate-50 focus:bg-white focus:outline-none focus:shadow-[4px_4px_0px_0px_rgba(34,197,94,1)] transition-all font-black text-lg"
                    :class="{ 'focus:shadow-red-500': isFeverMode }"
                    :placeholder="isLogin ? 'USERNAME' : 'YOUR NAME'"
                    @focus="focusField = 'username'"
                    @blur="focusField = null"
                    @input="handleInput"
                  />
                </div>
              </div>
    
              <!-- 密码 -->
              <div class="group relative">
                <label class="block text-xs font-black uppercase tracking-wider mb-1 ml-1">{{ t('login.secretKey') }}</label>
                <div class="relative transition-transform group-focus-within:-translate-y-1 group-focus-within:scale-[1.02]">
                  <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <svg class="w-6 h-6 text-black" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/></svg>
                  </div>
                  <input
                    v-model="form.password"
                    :type="showPassword ? 'text' : 'password'"
                    class="block w-full pl-12 pr-12 py-3 border-4 border-black rounded-xl bg-slate-50 focus:bg-white focus:outline-none focus:shadow-[4px_4px_0px_0px_rgba(236,72,153,1)] transition-all font-black text-lg"
                    :class="{ 'focus:shadow-red-500': isFeverMode }"
                    placeholder="PASSWORD"
                    @focus="focusField = 'password'"
                    @blur="focusField = null"
                    @input="handleInput"
                  />
                  <button type="button" @click="showPassword = !showPassword" @mouseenter="playSFX('hover')" class="absolute inset-y-0 right-0 pr-3 flex items-center hover:text-green-500">
                    <svg v-if="showPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/></svg>
                    <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
                  </button>
                </div>
              </div>
    
              <!-- 确认密码（注册模式） -->
              <Transition name="slide">
                <div v-if="!isLogin" class="group relative">
                  <label class="block text-xs font-black uppercase tracking-wider mb-1 ml-1">{{ t('login.confirmKey') }}</label>
                  <div class="relative transition-transform group-focus-within:-translate-y-1 group-focus-within:scale-[1.02]">
                    <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                      <svg class="w-6 h-6 text-black" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/></svg>
                    </div>
                    <input
                      v-model="form.confirmPassword"
                      :type="showConfirmPassword ? 'text' : 'password'"
                      class="block w-full pl-12 pr-12 py-3 border-4 border-black rounded-xl bg-slate-50 focus:bg-white focus:outline-none focus:shadow-[4px_4px_0px_0px_rgba(234,179,8,1)] transition-all font-black text-lg"
                      :class="{ 'focus:shadow-red-500': isFeverMode }"
                      placeholder="REPEAT PASS"
                      @focus="focusField = 'confirm'"
                      @blur="focusField = null"
                      @input="handleInput"
                    />
                    <button type="button" @click="showConfirmPassword = !showConfirmPassword" @mouseenter="playSFX('hover')" class="absolute inset-y-0 right-0 pr-3 flex items-center hover:text-green-500">
                      <svg v-if="showConfirmPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/></svg>
                      <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
                    </button>
                  </div>
                </div>
              </Transition>
    
              <!-- 提交按钮 -->
              <button
                type="submit"
                :disabled="isLoading"
                @mouseenter="playSFX('hover')"
                class="w-full mt-6 bg-black text-white py-4 rounded-xl border-4 border-transparent flex items-center justify-center gap-3 hover:-translate-y-1 transition-all active:shadow-none active:translate-y-1 disabled:opacity-80"
                :class="isLogin ? 'shadow-[8px_8px_0px_0px_rgba(34,197,94,1)] hover:shadow-[12px_12px_0px_0px_rgba(34,197,94,1)]' : 'shadow-[8px_8px_0px_0px_rgba(59,130,246,1)] hover:shadow-[12px_12px_0px_0px_rgba(59,130,246,1)]'"
              >
                <template v-if="isLoading">
                  <div class="w-full h-8 bg-gray-800 rounded-full overflow-hidden border-2 border-white relative">
                    <div class="h-full bg-green-500 relative animate-loading-bar">
                      <div class="absolute right-0 top-1/2 -translate-y-1/2 translate-x-1/2 w-4 h-4 bg-yellow-400 rounded-full border-2 border-black z-10"></div>
                    </div>
                    <div class="absolute right-2 top-1/2 -translate-y-1/2 w-3 h-3 bg-red-500 rounded-full border border-white animate-pulse"></div>
                  </div>
                </template>
                <template v-else>
                  <span class="font-black text-2xl tracking-widest italic">{{ isLogin ? t('login.loadSave') : t('login.createHero') }}</span>
                  <svg v-if="isLogin" class="w-8 h-8 text-green-400 animate-pulse" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"/></svg>
                  <svg v-else class="w-8 h-8 text-blue-400 animate-bounce" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"/></svg>
                </template>
              </button>
            </form>
          </div>
        </div>
      </div>

      <!-- Right Wing: Leaderboard -->
      <div class="hidden lg:block transform rotate-6 hover:rotate-0 transition-transform duration-500 hover:scale-105">
        <Leaderboard />
      </div>
    </div>

    <!-- 底部跑马灯 -->
    <FooterMarquee />

    <!-- 结算弹窗 -->
    <SettlementModal :show="showSettlement" @next="goHome" />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { login, register } from '../api/auth'
import { useUserStore } from '../stores/user'
import { AudioEngine } from '../components/login/AudioEngine'
import confetti from 'canvas-confetti'

const { t, locale } = useI18n()

const toggleLocale = () => {
  locale.value = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
  localStorage.setItem('locale', locale.value)
}

import AutoSnake from '../components/login/AutoSnake.vue'
import Leaderboard from '../components/login/Leaderboard.vue'
import PixelAvatar from '../components/login/PixelAvatar.vue'
import FloatingIcons from '../components/login/FloatingIcons.vue'
import HudPanel from '../components/login/HudPanel.vue'
import MusicWidget from '../components/login/MusicWidget.vue'
import FooterMarquee from '../components/login/FooterMarquee.vue'
import SettlementModal from '../components/login/SettlementModal.vue'

const router = useRouter()
const userStore = useUserStore()

// 状态
const isLogin = ref(true)
const isFeverMode = ref(false)
const focusField = ref(null)
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const isLoading = ref(false)
const showSettlement = ref(false)
const isPlayingMusic = ref(AudioEngine.getUserMusicPreference())
const avatarMessage = ref(t('login.insertCoinBubble'))

const form = reactive({
  username: '',
  password: '',
  confirmPassword: ''
})

// Combo 系统
const combo = ref(0)
let comboTimer = null

const handleInput = () => {
  AudioEngine.playSFX('type')
  combo.value++
  
  // Power Mode: Screen Shake on input
  if (isFeverMode.value) {
    const intensity = Math.min(combo.value, 10)
    document.body.style.transform = `translate(${Math.random() * intensity - intensity/2}px, ${Math.random() * intensity - intensity/2}px)`
    setTimeout(() => document.body.style.transform = 'none', 50)
    
    // Extra particles in fever mode
    confetti({ 
      particleCount: 3, 
      spread: 20, 
      startVelocity: 10, 
      origin: { y: 0.8 }, 
      colors: ['#ff0000', '#ffff00'] 
    })
  }

  if (comboTimer) clearTimeout(comboTimer)
  comboTimer = setTimeout(() => { combo.value = 0 }, 1000)

  // 每 5 连击触发特效
  if (combo.value > 0 && combo.value % 5 === 0) {
    AudioEngine.playSFX('success')
    confetti({ particleCount: 20, spread: 30, startVelocity: 20, origin: { y: 0.6 }, colors: ['#22c55e'] })
  }
}

// 音效
const playSFX = (type) => AudioEngine.playSFX(type)

// 音乐控制
let musicToggling = false
const toggleMusic = () => {
  if (musicToggling) return
  musicToggling = true
  const playing = AudioEngine.toggleBGM()
  isPlayingMusic.value = playing
  AudioEngine.setUserMusicPreference(playing)
  setTimeout(() => { musicToggling = false }, 300)
}

// 模式切换
const toggleMode = () => {
  AudioEngine.playSFX('switch')
  isLogin.value = !isLogin.value
  avatarMessage.value = isLogin.value ? t('login.welcomeBack') : t('login.newHero')
  combo.value = 0
  form.confirmPassword = ''
}

// 狂热模式切换
const toggleFeverMode = () => {
  isFeverMode.value = !isFeverMode.value
  AudioEngine.playSFX(isFeverMode.value ? 'win' : 'switch')
  
  if (isFeverMode.value) {
    if (isPlayingMusic.value) {
      AudioEngine.tempo = 200 // Speed up music
    }
  } else {
    AudioEngine.tempo = 150
  }
}

// 头像点击
const handleAvatarClick = () => {
  AudioEngine.playSFX('hover')
  const messages = [
    t('login.justAPixel'),
    t('login.keepTyping'),
    t('login.youGotThis'),
    t('login.wenote'),
    t('login.snakeMsg1'),
    t('login.snakeMsg2'),
    t('login.snakeMsg3'),
    t('login.snakeMsg4'),
    t('login.snakeMsg5'),
    t('login.snakeMsg6'),
    t('login.snakeMsg7'),
    t('login.snakeMsg8'),
    t('login.snakeMsg9'),
    t('login.snakeMsg10'),
    t('login.snakeMsg11'),
    t('login.snakeMsg12'),
    t('login.snakeMsg13'),
    t('login.snakeMsg14'),
    t('login.snakeMsg15'),
    t('login.snakeMsg16'),
    t('login.snakeMsg17'),
    t('login.snakeMsg18'),
    t('login.snakeMsg19'),
    t('login.snakeMsg20'),
    t('login.snakeMsg21'),
    t('login.snakeMsg22'),
    t('login.snakeMsg23'),
    t('login.snakeMsg24'),
    t('login.snakeMsg25'),
    t('login.snakeMsg26'),
    t('login.snakeMsg27'),
    t('login.snakeMsg28'),
    t('login.snakeMsg29'),
    t('login.snakeMsg30'),
    t('login.snakeMsg31'),
    t('login.snakeMsg32'),
    t('login.snakeMsg33'),
    t('login.snakeMsg34'),
    t('login.snakeMsg35'),
    t('login.snakeMsg36'),
    t('login.snakeMsg37'),
    t('login.snakeMsg38'),
    t('login.snakeMsg39'),
    t('login.snakeMsg40'),
    t('login.snakeMsg41'),
    t('login.snakeMsg42'),
    t('login.snakeMsg43'),
    t('login.snakeMsg44'),
    t('login.snakeMsg45'),
    t('login.snakeMsg46')
  ]
  avatarMessage.value = messages[Math.floor(Math.random() * messages.length)]
}

// 烟花特效
const fireFireworks = () => {
  const end = Date.now() + 2000
  const frame = () => {
    confetti({ particleCount: 5, angle: 60, spread: 55, origin: { x: 0 }, colors: ['#22c55e', '#eab308', '#ec4899'] })
    confetti({ particleCount: 5, angle: 120, spread: 55, origin: { x: 1 }, colors: ['#22c55e', '#eab308', '#ec4899'] })
    if (Date.now() < end) requestAnimationFrame(frame)
  }
  frame()
}

// 提交表单
const handleSubmit = async () => {
  if (!form.username || !form.password) {
    AudioEngine.playSFX('error')
    avatarMessage.value = t('login.fillAllFields')
    return
  }

  if (!isLogin.value && form.password !== form.confirmPassword) {
    AudioEngine.playSFX('error')
    avatarMessage.value = t('login.passwordsDontMatch')
    return
  }

  AudioEngine.playSFX('start')
  isLoading.value = true
  const startTime = Date.now()

  try {
    if (isLogin.value) {
      const data = await login({ username: form.username, password: form.password })
      userStore.setToken(data.token)
      userStore.setUser(data.user)
      // 保存凭据
      localStorage.setItem('wenote_saved_username', form.username)
      localStorage.setItem('wenote_saved_password', form.password)
      // 等待进度条动画完成（0.8 秒）
      const elapsed = Date.now() - startTime
      if (elapsed < 800) await new Promise(r => setTimeout(r, 800 - elapsed))
      isLoading.value = false
      AudioEngine.playSFX('win')
      showSettlement.value = true
      fireFireworks()
    } else {
      await register({ username: form.username, password: form.password })
      const elapsed = Date.now() - startTime
      if (elapsed < 800) await new Promise(r => setTimeout(r, 800 - elapsed))
      isLoading.value = false
      AudioEngine.playSFX('win')
      ElMessage.success(t('login.registerSuccess'))
      isLogin.value = true
      form.password = ''
      form.confirmPassword = ''
    }
  } catch (e) {
    isLoading.value = false
    AudioEngine.playSFX('error')
    avatarMessage.value = t('login.tryAgain')
  }
}

// 跳转首页
const goHome = () => {
  router.push('/')
}

// 重置表单
const resetForm = () => {
  showSettlement.value = false
  form.username = ''
  form.password = ''
  form.confirmPassword = ''
  combo.value = 0
}

onUnmounted(() => {
  if (comboTimer) clearTimeout(comboTimer)
})

// Auto-play BGM on first user interaction
onMounted(() => {
  // 读取保存的凭据
  const savedUsername = localStorage.getItem('wenote_saved_username')
  const savedPassword = localStorage.getItem('wenote_saved_password')
  if (savedUsername) form.username = savedUsername
  if (savedPassword) form.password = savedPassword

  window.addEventListener('click', () => {
    AudioEngine.init()
    if (isPlayingMusic.value && !AudioEngine.isPlaying) {
      AudioEngine.startBGM()
    }
  }, { once: true })
})
</script>

<style scoped>
/* Slogan Animation */
.slogan-text {
  -webkit-text-stroke: 2px white;
}
.slogan-word {
  display: inline-block;
  animation: bounce-letter 2s ease-in-out infinite;
  animation-delay: calc(var(--i) * 0.15s);
}
@keyframes bounce-letter {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

.combo-enter-active, .combo-leave-active {
  transition: all 0.3s ease;
}
.combo-enter-from, .combo-leave-to {
  opacity: 0;
  transform: translateX(-50%) scale(0);
}

.slide-enter-active, .slide-leave-active {
  transition: all 0.3s ease;
}
.slide-enter-from, .slide-leave-to {
  opacity: 0;
  height: 0;
  overflow: hidden;
}

@keyframes loading-bar {
  from { width: 0%; }
  to { width: 100%; }
}
.animate-loading-bar {
  animation: loading-bar 0.8s linear;
}

/* Shake Animation */
@keyframes shake {
  0%, 100% { transform: translate(0, 0) rotate(0deg); }
  25% { transform: translate(-2px, 2px) rotate(-1deg); }
  50% { transform: translate(2px, -2px) rotate(1deg); }
  75% { transform: translate(-2px, -2px) rotate(-1deg); }
}

.animate-shake {
  animation: shake 0.2s ease-in-out infinite;
}
</style>
