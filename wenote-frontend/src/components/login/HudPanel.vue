<template>
  <div class="absolute top-6 left-6 hidden md:flex flex-col gap-3 z-30 pointer-events-none select-none">
    <!-- HP 血条 -->
    <div class="flex items-center gap-3 bg-white/90 p-2 rounded-lg border-2 border-black shadow-[4px_4px_0_rgba(0,0,0,0.2)] backdrop-blur-sm -rotate-1 hover:rotate-0 transition-transform">
      <svg class="w-6 h-6 text-red-500 fill-red-500 animate-pulse" viewBox="0 0 24 24"><path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/></svg>
      <div class="w-40 h-5 border-2 border-black rounded-full bg-slate-200 relative overflow-hidden">
        <div class="h-full bg-red-500 relative animate-hp" style="width: 85%">
          <div class="absolute top-0 right-0 bottom-0 w-1 bg-white/30"></div>
          <div class="absolute top-1 right-2 w-full h-1 bg-white/20 rounded-full"></div>
        </div>
      </div>
      <span class="font-black text-xs text-slate-700">{{ t('hud.level') }}.99</span>
    </div>

    <!-- MP 能量条 -->
    <div class="flex items-center gap-3 bg-white/90 p-2 rounded-lg border-2 border-black shadow-[4px_4px_0_rgba(0,0,0,0.2)] backdrop-blur-sm rotate-1 hover:rotate-0 transition-transform">
      <svg class="w-6 h-6 text-blue-500 fill-blue-500" viewBox="0 0 24 24"><path d="M7 2v11h3v9l7-12h-4l4-8z"/></svg>
      <div class="w-32 h-5 border-2 border-black rounded-full bg-slate-200 relative overflow-hidden">
        <div class="h-full bg-blue-400 w-3/4 relative">
          <div class="absolute top-1 right-2 w-full h-1 bg-white/20 rounded-full"></div>
        </div>
      </div>
      <span class="font-black text-xs text-slate-700">{{ t('hud.max') }}</span>
    </div>

    <!-- 排行榜 -->
    <div class="hidden lg:block bg-white/90 backdrop-blur border-4 border-black p-4 rounded-2xl shadow-[8px_8px_0px_0px_rgba(0,0,0,0.2)] w-64 -rotate-2 hover:rotate-0 transition-transform pointer-events-auto">
      <div class="flex items-center gap-2 mb-4 border-b-2 border-black/10 pb-2">
        <svg class="w-5 h-5 text-yellow-500 fill-yellow-500" viewBox="0 0 24 24"><path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"/></svg>
        <h3 class="font-black text-slate-800 italic tracking-tighter">{{ t('hud.leaderboard') }}</h3>
      </div>
      <div class="space-y-3 font-mono text-sm font-bold">
        <div v-for="(player, i) in leaderboard" :key="i" class="flex justify-between items-center group cursor-pointer hover:bg-black/5 p-1 rounded transition-colors">
          <div class="flex items-center gap-2">
            <span class="w-5 h-5 rounded flex items-center justify-center text-[10px] text-white" :class="i === 0 ? 'bg-yellow-400 border border-black' : 'bg-slate-400'">{{ i + 1 }}</span>
            <span class="text-slate-700 group-hover:text-black">{{ player.name }}</span>
          </div>
          <span :class="player.color">{{ player.score }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const leaderboard = [
  { name: 'SNAKE_GOD', score: '99,999', color: 'text-red-500' },
  { name: 'CODE_MASTER', score: '88,400', color: 'text-orange-500' },
  { name: 'WENOTE_DEV', score: '72,000', color: 'text-blue-500' }
]
</script>

<style scoped>
@keyframes hp-pulse {
  0%, 100% { width: 85%; }
  25% { width: 100%; }
  50% { width: 90%; }
  75% { width: 95%; }
}
.animate-hp {
  animation: hp-pulse 5s ease-in-out infinite;
}
</style>
