<template>
  <div class="leaderboard-container bg-black p-4 rounded-xl border-4 border-slate-700 shadow-xl w-64">
    <!-- Header -->
    <div class="flex items-center justify-between mb-4 border-b-2 border-slate-700 pb-2">
      <h3 class="text-yellow-400 font-black text-xs tracking-widest animate-pulse">{{ t('leaderboard.title') }}</h3>
      <div class="w-2 h-2 bg-red-500 rounded-full animate-ping"></div>
    </div>

    <!-- List -->
    <div class="space-y-3 relative overflow-hidden h-[300px]">
      <div class="animate-scroll-up">
        <!-- Mock Data -->
        <div v-for="(player, index) in players" :key="index" class="flex items-center gap-3 mb-3 font-mono text-xs">
          <div
            class="w-6 h-6 flex items-center justify-center font-bold rounded"
            :class="index < 3 ? 'bg-yellow-400 text-black' : 'bg-slate-800 text-slate-400'"
          >
            {{ index + 1 }}
          </div>
          <div class="flex-1">
            <div class="text-green-400 font-bold truncate">{{ player.name }}</div>
            <div class="text-slate-500 text-[10px]">{{ player.title }}</div>
          </div>
          <div class="text-pink-500 font-bold">{{ player.score }}</div>
        </div>
        <!-- Duplicate for seamless scroll -->
        <div v-for="(player, index) in players" :key="`dup-${index}`" class="flex items-center gap-3 mb-3 font-mono text-xs">
          <div
            class="w-6 h-6 flex items-center justify-center font-bold rounded"
            :class="index < 3 ? 'bg-yellow-400 text-black' : 'bg-slate-800 text-slate-400'"
          >
            {{ index + 1 }}
          </div>
          <div class="flex-1">
            <div class="text-green-400 font-bold truncate">{{ player.name }}</div>
            <div class="text-slate-500 text-[10px]">{{ player.title }}</div>
          </div>
          <div class="text-pink-500 font-bold">{{ player.score }}</div>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div class="mt-4 text-center">
      <div class="text-[10px] text-slate-500 font-mono">{{ t('leaderboard.rank', { rank: 999 }) }}</div>
      <div class="text-[10px] text-green-500 font-bold animate-pulse">{{ t('leaderboard.insertCoin') }}</div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const players = [
  { name: 'NEO_ONE', score: '99999', title: 'The Chosen One' },
  { name: 'TRINITY', score: '88450', title: 'Hacker Elite' },
  { name: 'MORPHEUS', score: '76200', title: 'Captain' },
  { name: 'AGENT_S', score: '54000', title: 'Program' },
  { name: 'CYPHER', score: '43200', title: 'Traitor' },
  { name: 'TANK', score: '32100', title: 'Operator' },
  { name: 'DOZER', score: '28900', title: 'Pilot' },
  { name: 'MOUSE', score: '15400', title: 'Coder' },
  { name: 'SWITCH', score: '12300', title: 'Guardian' },
  { name: 'APOC', score: '09800', title: 'Driver' }
]
</script>

<style scoped>
@keyframes scroll-up {
  0% { transform: translateY(0); }
  100% { transform: translateY(-50%); }
}

.animate-scroll-up {
  animation: scroll-up 8s linear infinite;
}

.animate-scroll-up:hover {
  animation-play-state: paused;
}
</style>