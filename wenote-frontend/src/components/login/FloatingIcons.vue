<template>
  <div class="absolute inset-0 overflow-hidden pointer-events-none hidden md:block">
    <!-- 云朵 -->
    <div
      v-for="cloud in clouds"
      :key="cloud.id"
      class="absolute text-white opacity-40 drop-shadow-md"
      :style="{ top: cloud.y + 'vh', animationDuration: cloud.duration + 's', animationDelay: cloud.delay + 's' }"
    >
      <svg class="w-16 h-16 animate-cloud" :style="{ transform: `scale(${cloud.scale})` }" fill="currentColor" viewBox="0 0 24 24">
        <path d="M19.35 10.04C18.67 6.59 15.64 4 12 4 9.11 4 6.6 5.64 5.35 8.04 2.34 8.36 0 10.91 0 14c0 3.31 2.69 6 6 6h13c2.76 0 5-2.24 5-5 0-2.64-2.05-4.78-4.65-4.96z"/>
      </svg>
    </div>

    <!-- 游戏图标 -->
    <div
      v-for="item in floatingItems"
      :key="item.id"
      class="absolute opacity-30"
      :class="item.color"
      :style="{
        left: item.x + 'vw',
        top: item.y + 'vh',
        animationDuration: item.duration + 's',
        animationDelay: item.delay + 's'
      }"
    >
      <component :is="item.icon" class="w-12 h-12 animate-float-rotate" :style="{ '--rotate-dir': item.rotateDir }" />
    </div>
  </div>
</template>

<script setup>
import { h } from 'vue'

// SVG 图标组件
const Gamepad = { render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24', class: 'w-12 h-12' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z' })]) }
const Star = { render: () => h('svg', { fill: 'currentColor', viewBox: '0 0 24 24', class: 'w-12 h-12' }, [h('path', { d: 'M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z' })]) }
const Ghost = { render: () => h('svg', { fill: 'currentColor', viewBox: '0 0 24 24', class: 'w-12 h-12' }, [h('path', { d: 'M12 2C7.58 2 4 5.58 4 10v10l2-2 2 2 2-2 2 2 2-2 2 2 2-2 2 2V10c0-4.42-3.58-8-8-8zm-2 9a2 2 0 110-4 2 2 0 010 4zm4 0a2 2 0 110-4 2 2 0 010 4z' })]) }
const Zap = { render: () => h('svg', { fill: 'currentColor', viewBox: '0 0 24 24', class: 'w-12 h-12' }, [h('path', { d: 'M7 2v11h3v9l7-12h-4l4-8z' })]) }
const Trophy = { render: () => h('svg', { fill: 'currentColor', viewBox: '0 0 24 24', class: 'w-12 h-12' }, [h('path', { d: 'M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z' })]) }
const Skull = { render: () => h('svg', { fill: 'currentColor', viewBox: '0 0 24 24', class: 'w-12 h-12' }, [h('path', { d: 'M12 2C6.48 2 2 6.48 2 12v8h4v-2h2v2h4v-2h2v2h4v-2h2v-6c0-5.52-4.48-10-10-10zm-3 12a2 2 0 110-4 2 2 0 010 4zm6 0a2 2 0 110-4 2 2 0 010 4z' })]) }

const icons = [
  { icon: Gamepad, color: 'text-green-400' },
  { icon: Star, color: 'text-yellow-400' },
  { icon: Ghost, color: 'text-pink-400' },
  { icon: Zap, color: 'text-blue-400' },
  { icon: Trophy, color: 'text-orange-400' },
  { icon: Skull, color: 'text-red-400' }
]

const clouds = Array.from({ length: 3 }, (_, i) => ({
  id: `cloud-${i}`,
  y: 15 + i * 30,
  scale: 1 + Math.random() * 0.5,
  duration: 60 + i * 20,
  delay: i * 5
}))

const floatingItems = Array.from({ length: 15 }, (_, i) => ({
  id: i,
  ...icons[i % icons.length],
  x: Math.random() * 100,
  y: Math.random() * 100,
  duration: 15 + Math.random() * 20,
  delay: Math.random() * 5,
  rotateDir: Math.random() > 0.5 ? 1 : -1
}))
</script>

<style scoped>
@keyframes cloud-move {
  from { transform: translateX(-200px); }
  to { transform: translateX(120vw); }
}
.animate-cloud {
  animation: cloud-move var(--duration, 60s) linear infinite;
}

@keyframes float-rotate {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(calc(var(--rotate-dir, 1) * 180deg)); }
}
.animate-float-rotate {
  animation: float-rotate var(--duration, 15s) ease-in-out infinite;
}
</style>
