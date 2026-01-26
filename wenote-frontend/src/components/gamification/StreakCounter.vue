<template>
  <div class="streak-counter bg-white border-4 border-black rounded-2xl p-4 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
    <div class="flex items-center gap-3">
      <!-- Fire Icon with animation -->
      <div class="relative">
        <span class="text-4xl" :class="{ 'animate-bounce': streak > 0 }">
          {{ streak > 0 ? '🔥' : '❄️' }}
        </span>
        <span v-if="atRisk" class="absolute -top-1 -right-1 w-3 h-3 bg-yellow-400 rounded-full animate-ping"></span>
      </div>

      <div>
        <div class="text-xs font-black text-slate-400 uppercase tracking-wider">连续天数</div>
        <div class="text-3xl font-black" :class="streak > 0 ? 'text-orange-500' : 'text-slate-300'">
          {{ streak }} <span class="text-lg">天</span>
        </div>
        <div class="text-xs text-slate-400">
          最长: {{ longest }} 天
        </div>
      </div>
    </div>

    <!-- Warning if streak at risk -->
    <div v-if="atRisk && streak > 0" class="mt-3 px-3 py-2 bg-yellow-100 border-2 border-yellow-400 rounded-xl text-xs font-bold text-yellow-700 flex items-center gap-2">
      <span>⚠️</span>
      <span>今天还没写笔记，连续天数有风险！</span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  streak: {
    type: Number,
    default: 0
  },
  longest: {
    type: Number,
    default: 0
  },
  atRisk: {
    type: Boolean,
    default: false
  }
})
</script>
