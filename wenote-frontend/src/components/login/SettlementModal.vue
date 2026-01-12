<template>
  <Transition name="modal">
    <div v-if="show" class="fixed inset-0 z-[100] flex items-center justify-center bg-black/80 backdrop-blur-sm">
      <div class="bg-gradient-to-b from-yellow-300 to-yellow-500 p-1 rounded-[2.5rem] shadow-[0_0_50px_rgba(234,179,8,0.5)]">
        <div class="bg-white border-4 border-black rounded-[2.3rem] p-8 w-[90vw] max-w-lg text-center overflow-hidden relative">
          <!-- 背景光效 -->
          <div class="absolute inset-0 bg-[radial-gradient(circle_at_center,rgba(255,255,0,0.2)_0%,transparent_70%)] animate-pulse pointer-events-none"></div>

          <!-- 标题 -->
          <div class="font-black text-5xl italic text-transparent bg-clip-text bg-gradient-to-b from-yellow-400 to-orange-500 drop-shadow-[4px_4px_0_#000] mb-4 animate-bounce-in">
            {{ t('settlement.stageClear') }}
          </div>

          <!-- 欢迎信息 -->
          <div class="text-slate-500 font-bold mb-6">{{ t('login.welcomeBack') }}</div>

          <!-- 统计 -->
          <div class="flex justify-center gap-4 mb-6">
            <div class="bg-slate-100 p-4 rounded-xl border-2 border-black w-1/3">
              <div class="text-xs font-black text-slate-400 mb-1">{{ t('settlement.score') }}</div>
              <div class="text-3xl font-black text-slate-800">{{ score }}</div>
            </div>
            <div class="bg-black p-4 rounded-xl border-2 border-black w-1/3 text-white relative overflow-hidden">
              <div class="text-xs font-black text-slate-400 mb-1">{{ t('settlement.rank') }}</div>
              <div class="text-4xl font-black text-yellow-400">{{ rank }}</div>
            </div>
          </div>

          <!-- 跳动头像 -->
          <div class="w-20 h-20 mx-auto bg-green-500 rounded-full border-4 border-black flex items-center justify-center mb-6 relative shadow-lg animate-bounce">
            <div class="w-14 h-7 bg-white rounded-full border-2 border-black flex items-center justify-center">
              <div class="w-8 h-3 bg-red-400 rounded-full"></div>
            </div>
            <svg class="absolute -top-6 -right-3 w-8 h-8 text-yellow-400 fill-yellow-400 drop-shadow-md animate-bounce" viewBox="0 0 24 24">
              <path d="M5 16L3 5l5.5 5L12 4l3.5 6L21 5l-2 11H5zm14 3c0 .6-.4 1-1 1H6c-.6 0-1-.4-1-1v-1h14v1z"/>
            </svg>
          </div>

          <!-- 自动跳转提示 -->
          <div class="text-slate-400 text-sm font-bold">
            {{ countdown }}s 后自动进入...
          </div>

          <!-- 进度条 -->
          <div class="mt-4 h-2 bg-slate-200 rounded-full overflow-hidden">
            <div class="h-full bg-green-500 transition-all duration-[1500ms] ease-linear" :style="{ width: progressWidth }"></div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { computed, ref, watch, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { AudioEngine } from './AudioEngine.js'

const { t } = useI18n()

const props = defineProps({
  show: Boolean
})

const emit = defineEmits(['next'])

const countdown = ref(1)
const progress = ref(0)
let timer = null

const progressWidth = computed(() => `${progress.value}%`)

watch(() => props.show, (val) => {
  if (val) {
    countdown.value = 1
    progress.value = 0

    // 进度条动画
    setTimeout(() => progress.value = 100, 100)

    // 倒计时结束时播放音效并跳转
    timer = setTimeout(() => {
      countdown.value = 0
      AudioEngine.playSFX('firework')
      // 等音效播放完再跳转
      setTimeout(() => emit('next'), 600)
    }, 1500)
  } else {
    clearTimeout(timer)
  }
})

onUnmounted(() => {
  clearTimeout(timer)
})

const score = computed(() => 99999)

const rank = computed(() => 'SSS')
</script>

<style scoped>
.modal-enter-active, .modal-leave-active {
  transition: all 0.3s ease;
}
.modal-enter-from, .modal-leave-to {
  opacity: 0;
}
.modal-enter-from > div, .modal-leave-to > div {
  transform: scale(0.8);
}

@keyframes bounce-in {
  0% { transform: scale(2); opacity: 0; }
  50% { transform: scale(1.1); }
  100% { transform: scale(1); opacity: 1; }
}
.animate-bounce-in {
  animation: bounce-in 0.5s ease-out;
}
</style>
