<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { getStatsOverview, getStatsTrend, getStatsTags, getStatsNotebooks } from '../../api/stats'
import { useGamificationStore } from '../../stores/gamification'
import * as echarts from 'echarts'
import { BarChart3, TrendingUp, Tag, BookOpen } from 'lucide-vue-next'
import StreakCounter from '../gamification/StreakCounter.vue'
import DailyGoalProgress from '../gamification/DailyGoalProgress.vue'
import AchievementGallery from '../gamification/AchievementGallery.vue'
import WritingReport from '../gamification/WritingReport.vue'

const { t } = useI18n()
const gamificationStore = useGamificationStore()

const overview = ref({})
const loading = ref(false)

// 图表实例
const trendChartRef = ref(null)
const tagChartRef = ref(null)
const notebookChartRef = ref(null)

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    // 加载游戏化数据
    await gamificationStore.fetchStatus()
    await gamificationStore.fetchAchievements()

    // 加载概览数据
    const overviewData = await getStatsOverview()
    overview.value = overviewData

    // 加载趋势数据并绘制图表
    const trendData = await getStatsTrend(7)
    renderTrendChart(trendData.list)

    // 加载标签统计并绘制图表
    const tagData = await getStatsTags(10)
    renderTagChart(tagData.list)

    // 加载笔记本统计并绘制图表
    const notebookData = await getStatsNotebooks()
    renderNotebookChart(notebookData.list)
  } catch (error) {
    console.error('Load stats failed:', error)
  } finally {
    loading.value = false
  }
}

// 绘制趋势图
const renderTrendChart = (data) => {
  if (!trendChartRef.value) return

  const chart = echarts.init(trendChartRef.value)
  const option = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      borderColor: '#000',
      borderWidth: 2,
      textStyle: { color: '#fff' }
    },
    xAxis: {
      type: 'category',
      data: data.map(d => d.date.substring(5)), // 只显示月-日
      axisLine: { lineStyle: { color: '#000', width: 2 } },
      axisLabel: { color: '#666', fontWeight: 'bold' }
    },
    yAxis: {
      type: 'value',
      axisLine: { lineStyle: { color: '#000', width: 2 } },
      axisLabel: { color: '#666', fontWeight: 'bold' },
      splitLine: { lineStyle: { type: 'dashed', color: '#e5e7eb' } }
    },
    series: [{
      data: data.map(d => d.count),
      type: 'line',
      smooth: true,
      lineStyle: { color: '#10b981', width: 3 },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: 'rgba(16, 185, 129, 0.3)' },
          { offset: 1, color: 'rgba(16, 185, 129, 0.05)' }
        ])
      },
      itemStyle: { color: '#10b981', borderColor: '#000', borderWidth: 2 }
    }]
  }
  chart.setOption(option)

  // 响应式调整
  window.addEventListener('resize', () => chart.resize())
}

// 绘制标签统计图
const renderTagChart = (data) => {
  if (!tagChartRef.value || !data.length) return

  const chart = echarts.init(tagChartRef.value)
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      borderColor: '#000',
      borderWidth: 2,
      textStyle: { color: '#fff' }
    },
    xAxis: {
      type: 'value',
      axisLine: { lineStyle: { color: '#000', width: 2 } },
      axisLabel: { color: '#666', fontWeight: 'bold' },
      splitLine: { lineStyle: { type: 'dashed', color: '#e5e7eb' } }
    },
    yAxis: {
      type: 'category',
      data: data.map(d => d.tag_name),
      axisLine: { lineStyle: { color: '#000', width: 2 } },
      axisLabel: { color: '#666', fontWeight: 'bold' }
    },
    series: [{
      data: data.map(d => d.count),
      type: 'bar',
      itemStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
          { offset: 0, color: '#8b5cf6' },
          { offset: 1, color: '#a78bfa' }
        ]),
        borderColor: '#000',
        borderWidth: 2
      },
      barWidth: '60%'
    }]
  }
  chart.setOption(option)
  window.addEventListener('resize', () => chart.resize())
}

// 绘制笔记本统计图
const renderNotebookChart = (data) => {
  if (!notebookChartRef.value || !data.length) return

  const chart = echarts.init(notebookChartRef.value)
  const colors = ['#f59e0b', '#10b981', '#3b82f6', '#ef4444', '#8b5cf6', '#ec4899']

  // 处理笔记本名称国际化（默认笔记本"未分类"需要翻译）
  const getNotebookDisplayName = (name, isDefault) => {
    if (isDefault || name === '未分类') {
      return t('sidebar.uncategorized')
    }
    return name
  }

  const option = {
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      borderColor: '#000',
      borderWidth: 2,
      textStyle: { color: '#fff' },
      formatter: '{b}: {c} ({d}%)'
    },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: true,
      itemStyle: {
        borderColor: '#000',
        borderWidth: 2
      },
      label: {
        color: '#000',
        fontWeight: 'bold'
      },
      data: data.map((d, i) => ({
        name: getNotebookDisplayName(d.notebook_name, d.is_default),
        value: d.count,
        itemStyle: { color: colors[i % colors.length] }
      }))
    }]
  }
  chart.setOption(option)
  window.addEventListener('resize', () => chart.resize())
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="stats-dashboard p-6 space-y-6">
    <!-- 标题 -->
    <div class="flex items-center gap-3">
      <BarChart3 class="w-8 h-8" />
      <h2 class="text-3xl font-black">{{ t('stats.dataStats') }}</h2>
    </div>

    <!-- 游戏化组件 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <StreakCounter
        :streak="gamificationStore.streakDisplay.current"
        :longest="gamificationStore.streakDisplay.longest"
        :at-risk="gamificationStore.streakDisplay.atRisk"
      />
      <DailyGoalProgress
        :current="gamificationStore.goalProgress.current"
        :target="gamificationStore.goalProgress.target"
        :percent="gamificationStore.goalProgress.percent"
        :completed="gamificationStore.goalProgress.completed"
        @update-goal="gamificationStore.setDailyGoal"
      />
      <WritingReport />
    </div>

    <!-- 成就展示 -->
    <AchievementGallery :achievements="gamificationStore.achievements" />

    <!-- 概览卡片 -->
    <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
      <div class="bg-white border-4 border-black rounded-2xl p-4 shadow-[6px_6px_0px_0px_rgba(0,0,0,1)]">
        <div class="text-sm font-bold text-slate-400 uppercase mb-1">{{ t('stats.totalNotes') }}</div>
        <div class="text-4xl font-black text-slate-800">{{ overview.total_notes || 0 }}</div>
        <div class="text-xs text-green-600 font-bold mt-1">{{ t('stats.thisWeek') }} +{{ overview.this_week_notes || 0 }}</div>
      </div>

      <div class="bg-white border-4 border-black rounded-2xl p-4 shadow-[6px_6px_0px_0px_rgba(0,0,0,1)]">
        <div class="text-sm font-bold text-slate-400 uppercase mb-1">{{ t('stats.notebooks') }}</div>
        <div class="text-4xl font-black text-slate-800">{{ overview.total_notebooks || 0 }}</div>
      </div>

      <div class="bg-white border-4 border-black rounded-2xl p-4 shadow-[6px_6px_0px_0px_rgba(0,0,0,1)]">
        <div class="text-sm font-bold text-slate-400 uppercase mb-1">{{ t('stats.tags') }}</div>
        <div class="text-4xl font-black text-slate-800">{{ overview.total_tags || 0 }}</div>
      </div>

      <div class="bg-white border-4 border-black rounded-2xl p-4 shadow-[6px_6px_0px_0px_rgba(0,0,0,1)]">
        <div class="text-sm font-bold text-slate-400 uppercase mb-1">{{ t('stats.totalWords') }}</div>
        <div class="text-4xl font-black text-slate-800">{{ (overview.total_words || 0).toLocaleString() }}</div>
        <div class="text-xs text-blue-600 font-bold mt-1">{{ t('stats.thisWeek') }} +{{ (overview.this_week_words || 0).toLocaleString() }}</div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- 趋势图 -->
      <div class="bg-white border-4 border-black rounded-2xl p-6 shadow-[6px_6px_0px_0px_rgba(0,0,0,1)]">
        <div class="flex items-center gap-2 mb-4">
          <TrendingUp class="w-5 h-5 text-green-600" />
          <h3 class="text-xl font-black">{{ t('stats.last7Days') }}</h3>
        </div>
        <div ref="trendChartRef" class="w-full h-64"></div>
      </div>

      <!-- 标签统计 -->
      <div class="bg-white border-4 border-black rounded-2xl p-6 shadow-[6px_6px_0px_0px_rgba(0,0,0,1)]">
        <div class="flex items-center gap-2 mb-4">
          <Tag class="w-5 h-5 text-purple-600" />
          <h3 class="text-xl font-black">{{ t('stats.top10Tags') }}</h3>
        </div>
        <div ref="tagChartRef" class="w-full h-64"></div>
      </div>

      <!-- 笔记本分布 -->
      <div class="bg-white border-4 border-black rounded-2xl p-6 shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] lg:col-span-2">
        <div class="flex items-center gap-2 mb-4">
          <BookOpen class="w-5 h-5 text-blue-600" />
          <h3 class="text-xl font-black">{{ t('stats.notebookDistribution') }}</h3>
        </div>
        <div ref="notebookChartRef" class="w-full h-80"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.stats-dashboard {
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>



