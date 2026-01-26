<template>
  <div class="writing-report bg-white border-4 border-black rounded-2xl p-4 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
    <div class="flex justify-between items-center mb-4">
      <div class="text-xs font-black text-slate-400 uppercase tracking-wider">写作报告</div>
      <div class="flex gap-2">
        <button
          v-for="p in periods"
          :key="p.id"
          @click="changePeriod(p.id)"
          class="px-3 py-1 rounded-lg text-xs font-bold transition-all border-2"
          :class="period === p.id
            ? 'bg-black text-white border-black'
            : 'bg-white text-slate-600 border-slate-200 hover:border-slate-400'"
        >
          {{ p.label }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="text-center py-8">
      <div class="animate-spin text-2xl">⏳</div>
    </div>

    <div v-else-if="report">
      <!-- Period Info -->
      <div class="text-sm text-slate-500 mb-4">
        {{ report.start_date }} ~ {{ report.end_date }}
      </div>

      <!-- Stats Grid -->
      <div class="grid grid-cols-2 gap-3 mb-4">
        <div class="bg-slate-50 rounded-xl p-3 border-2 border-slate-200">
          <div class="text-xs text-slate-400 font-bold">新增笔记</div>
          <div class="text-2xl font-black text-slate-800">{{ report.notes_created }}</div>
          <div class="text-xs" :class="report.notes_delta >= 0 ? 'text-green-500' : 'text-red-500'">
            {{ report.notes_delta >= 0 ? '+' : '' }}{{ report.notes_delta }} 较上期
          </div>
        </div>
        <div class="bg-slate-50 rounded-xl p-3 border-2 border-slate-200">
          <div class="text-xs text-slate-400 font-bold">写作字数</div>
          <div class="text-2xl font-black text-slate-800">{{ formatNumber(report.chars_written) }}</div>
          <div class="text-xs" :class="report.chars_delta >= 0 ? 'text-green-500' : 'text-red-500'">
            {{ report.chars_delta >= 0 ? '+' : '' }}{{ formatNumber(report.chars_delta) }} 较上期
          </div>
        </div>
        <div class="bg-slate-50 rounded-xl p-3 border-2 border-slate-200">
          <div class="text-xs text-slate-400 font-bold">活跃天数</div>
          <div class="text-2xl font-black text-slate-800">{{ report.active_days }}</div>
          <div class="text-xs" :class="report.active_days_delta >= 0 ? 'text-green-500' : 'text-red-500'">
            {{ report.active_days_delta >= 0 ? '+' : '' }}{{ report.active_days_delta }} 较上期
          </div>
        </div>
        <div class="bg-slate-50 rounded-xl p-3 border-2 border-slate-200">
          <div class="text-xs text-slate-400 font-bold">目标达成</div>
          <div class="text-2xl font-black text-slate-800">{{ report.goals_met }}</div>
          <div class="text-xs text-slate-400">次</div>
        </div>
      </div>

      <!-- Achievements Earned -->
      <div v-if="report.achievements_earned?.length > 0" class="mb-4">
        <div class="text-xs font-black text-slate-400 mb-2">本期解锁成就</div>
        <div class="flex gap-2 flex-wrap">
          <div
            v-for="a in report.achievements_earned"
            :key="a.id"
            class="flex items-center gap-1 px-2 py-1 bg-yellow-100 rounded-lg border-2 border-yellow-400"
          >
            <span>{{ a.icon }}</span>
            <span class="text-xs font-bold text-yellow-700">{{ a.name_zh || a.name }}</span>
          </div>
        </div>
      </div>

      <!-- Daily Chart Placeholder -->
      <div v-if="report.daily_stats?.length > 0" class="mt-4">
        <div class="text-xs font-black text-slate-400 mb-2">每日活动</div>
        <div class="flex items-end gap-1 h-20">
          <div
            v-for="(day, index) in report.daily_stats"
            :key="index"
            class="flex-1 rounded-t transition-all"
            :class="day.goal_met ? 'bg-green-500' : 'bg-blue-400'"
            :style="{ height: `${Math.max(10, (day.chars / maxChars) * 100)}%` }"
            :title="`${day.date}: ${day.chars} 字`"
          ></div>
        </div>
        <div class="flex justify-between text-[10px] text-slate-400 mt-1">
          <span>{{ report.daily_stats[0]?.date?.slice(5) }}</span>
          <span>{{ report.daily_stats[report.daily_stats.length - 1]?.date?.slice(5) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useGamificationStore } from '../../stores/gamification'

const gamificationStore = useGamificationStore()

const period = ref('week')
const report = ref(null)
const loading = ref(false)

const periods = [
  { id: 'week', label: '本周' },
  { id: 'month', label: '本月' }
]

const maxChars = computed(() => {
  if (!report.value?.daily_stats) return 1
  return Math.max(...report.value.daily_stats.map(d => d.chars), 1)
})

const formatNumber = (num) => {
  if (num >= 10000) return (num / 10000).toFixed(1) + '万'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'k'
  return num
}

const changePeriod = async (p) => {
  period.value = p
  await fetchReport()
}

const fetchReport = async () => {
  loading.value = true
  try {
    report.value = await gamificationStore.fetchReport(period.value)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchReport()
})
</script>
