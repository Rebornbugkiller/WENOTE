<template>
  <!-- Mobile: Simple icon button -->
  <div class="absolute top-4 right-4 md:hidden z-30">
    <button
      @click="$emit('toggle')"
      class="bg-black text-white p-2.5 rounded-xl shadow-[2px_2px_0px_0px_rgba(34,197,94,1)] border-2 border-white/20 cursor-pointer"
      :class="{ 'animate-pulse': isPlaying }"
    >
      <svg v-if="isPlaying" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3"/>
      </svg>
      <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z"/>
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2"/>
      </svg>
    </button>
  </div>

  <!-- Desktop: Full widget -->
  <div class="absolute top-6 right-6 hidden md:block z-30">
    <div
      @click="$emit('toggle')"
      class="bg-black text-white p-3 rounded-xl shadow-[4px_4px_0px_0px_rgba(34,197,94,1)] flex items-center gap-3 border-2 border-white/20 rotate-2 hover:rotate-0 transition-transform cursor-pointer"
      :class="{ 'animate-pulse': isPlaying }"
    >
      <div class="w-10 h-10 rounded-lg flex items-center justify-center border-2 border-white/20 relative overflow-hidden" :class="isPlaying ? 'bg-indigo-500' : 'bg-slate-700'">
        <svg v-if="isPlaying" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3"/>
        </svg>
        <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z"/>
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2"/>
        </svg>
      </div>
      <div class="flex flex-col">
        <div class="flex items-center gap-1">
          <span class="w-2 h-2 rounded-full" :class="isPlaying ? 'bg-green-500 animate-pulse' : 'bg-red-500'"></span>
          <span class="text-[10px] font-bold text-green-400 tracking-wider">{{ isPlaying ? 'NOW PLAYING' : 'MUSIC OFF' }}</span>
        </div>
        <div class="relative w-24 h-4 overflow-hidden">
          <div v-if="isPlaying" class="absolute whitespace-nowrap text-xs font-mono animate-marquee">
            8-BIT ADVENTURE.mp3 • WENOTE OST •
          </div>
          <span v-else class="text-xs font-mono text-slate-500">CLICK TO PLAY</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  isPlaying: Boolean
})
defineEmits(['toggle'])
</script>

<style scoped>
@keyframes marquee {
  from { transform: translateX(0); }
  to { transform: translateX(-100%); }
}
.animate-marquee {
  animation: marquee 5s linear infinite;
}
</style>
