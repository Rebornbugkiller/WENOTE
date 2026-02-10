<template>
  <div class="flex justify-center mb-4 md:mb-6 relative z-30">
    <!-- 对话气泡 - 手机端隐藏 -->
    <Transition name="bubble">
      <div v-if="message" class="absolute -top-16 md:-top-20 bg-white border-4 border-black p-2 md:p-3 rounded-2xl shadow-lg w-40 md:w-48 text-center hidden md:block">
        <p class="text-[10px] md:text-xs font-bold text-slate-800 leading-tight">{{ message }}</p>
        <div class="absolute -bottom-3 left-1/2 -translate-x-1/2 w-4 h-4 bg-white border-r-4 border-b-4 border-black transform rotate-45"></div>
      </div>
    </Transition>

    <!-- 贪吃蛇主体 -->
    <div
      ref="snakeRef"
      @click="$emit('click')"
      class="w-24 h-20 md:w-32 md:h-28 rounded-[50%] border-4 border-black relative flex items-center justify-center shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] cursor-pointer transition-all duration-300 hover:scale-105"
      :class="combo > 5 ? 'bg-emerald-400' : 'bg-green-500'"
    >
      <!-- 皇冠 -->
      <Transition name="crown">
        <svg v-if="combo > 5" class="absolute -top-8 md:-top-10 w-8 h-8 md:w-10 md:h-10 text-yellow-400 drop-shadow-[2px_2px_0px_rgba(0,0,0,1)] animate-bounce" fill="currentColor" viewBox="0 0 24 24">
          <path d="M5 16L3 5l5.5 5L12 4l3.5 6L21 5l-2 11H5zm14 3c0 .6-.4 1-1 1H6c-.6 0-1-.4-1-1v-1h14v1z"/>
        </svg>
      </Transition>

      <!-- 眼睛 -->
      <div class="flex gap-4 md:gap-6 mb-1 md:mb-2">
        <!-- 左眼 -->
        <div class="w-6 h-6 md:w-8 md:h-8 bg-white rounded-full border-2 border-black flex items-center justify-center overflow-hidden">
          <template v-if="isPasswordField">
            <!-- 闭眼 - 一条线 -->
            <div class="w-4 md:w-5 h-1 bg-black rounded-full"></div>
          </template>
          <template v-else>
            <!-- 睁眼 - 眼球跟随鼠标 -->
            <div
              class="w-3 h-3 md:w-4 md:h-4 bg-black rounded-full transition-transform duration-100"
              :style="{ transform: `translate(${eyeX * 0.3}px, ${eyeY * 0.3}px)` }"
            ></div>
          </template>
        </div>
        <!-- 右眼 -->
        <div class="w-6 h-6 md:w-8 md:h-8 bg-white rounded-full border-2 border-black flex items-center justify-center overflow-hidden">
          <template v-if="isPasswordField">
            <!-- 闭眼 -->
            <div class="w-4 md:w-5 h-1 bg-black rounded-full"></div>
          </template>
          <template v-else>
            <!-- 睁眼 -->
            <div
              class="w-3 h-3 md:w-4 md:h-4 bg-black rounded-full transition-transform duration-100"
              :style="{ transform: `translate(${eyeX * 0.3}px, ${eyeY * 0.3}px)` }"
            ></div>
          </template>
        </div>
      </div>

      <!-- 腮红 -->
      <div class="absolute bottom-4 md:bottom-6 left-3 md:left-4 w-3 h-1.5 md:w-4 md:h-2 bg-pink-300 rounded-full opacity-60"></div>
      <div class="absolute bottom-4 md:bottom-6 right-3 md:right-4 w-3 h-1.5 md:w-4 md:h-2 bg-pink-300 rounded-full opacity-60"></div>

      <!-- 舌头 -->
      <div class="absolute bottom-2 md:bottom-3 flex flex-col items-center">
        <div class="w-2 h-3 bg-red-400 rounded-b-full border-2 border-t-0 border-black"></div>
        <div class="flex gap-0.5 -mt-0.5">
          <div class="w-1 h-2 bg-red-400 rounded-b-full"></div>
          <div class="w-1 h-2 bg-red-400 rounded-b-full"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  focusField: String,
  combo: { type: Number, default: 0 },
  isLogin: { type: Boolean, default: true },
  message: String
})

defineEmits(['click'])

const snakeRef = ref(null)
const eyeX = ref(0)
const eyeY = ref(0)

const isPasswordField = computed(() =>
  props.focusField === 'password' || props.focusField === 'confirm'
)

const handleMouseMove = (e) => {
  if (!snakeRef.value || isPasswordField.value) return

  const rect = snakeRef.value.getBoundingClientRect()
  const centerX = rect.left + rect.width / 2
  const centerY = rect.top + rect.height / 2

  const deltaX = e.clientX - centerX
  const deltaY = e.clientY - centerY

  const maxOffset = 5
  const distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY)
  const scale = Math.min(distance / 100, 1)

  eyeX.value = (deltaX / distance) * maxOffset * scale || 0
  eyeY.value = (deltaY / distance) * maxOffset * scale || 0
}

onMounted(() => {
  window.addEventListener('mousemove', handleMouseMove)
})

onUnmounted(() => {
  window.removeEventListener('mousemove', handleMouseMove)
})
</script>

<style scoped>
.bubble-enter-active, .bubble-leave-active {
  transition: all 0.3s ease;
}
.bubble-enter-from, .bubble-leave-to {
  opacity: 0;
  transform: scale(0.8) translateY(10px);
}
.crown-enter-active, .crown-leave-active {
  transition: all 0.3s ease;
}
.crown-enter-from, .crown-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>
