<template>
  <div
    class="achievement-badge p-3 rounded-xl border-4 transition-all cursor-pointer hover:scale-105"
    :class="badgeClasses"
    @click="$emit('click', achievement)"
  >
    <div class="text-center">
      <div class="text-3xl mb-1" :class="{ 'grayscale opacity-50': !unlocked }">
        {{ achievement.icon }}
      </div>
      <div class="text-xs font-black truncate" :class="unlocked ? 'text-slate-800' : 'text-slate-400'">
        {{ achievementName }}
      </div>
      <div v-if="unlocked && achievement.unlocked_at" class="text-[10px] text-slate-400 mt-1">
        {{ formatDate(achievement.unlocked_at) }}
      </div>
      <div v-if="!unlocked" class="text-[10px] text-slate-400 mt-1">
        {{ getThresholdText() }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()

const props = defineProps({
  achievement: {
    type: Object,
    required: true
  },
  unlocked: {
    type: Boolean,
    default: false
  }
})

defineEmits(['click'])

const badgeClasses = computed(() => {
  if (!props.unlocked) {
    return 'border-slate-200 bg-slate-50'
  }
  const rarityColors = {
    common: 'border-slate-400 bg-white shadow-[2px_2px_0px_0px_rgba(0,0,0,0.2)]',
    rare: 'border-blue-400 bg-blue-50 shadow-[2px_2px_0px_0px_rgba(59,130,246,0.3)]',
    epic: 'border-purple-400 bg-purple-50 shadow-[2px_2px_0px_0px_rgba(168,85,247,0.3)]',
    legendary: 'border-yellow-400 bg-yellow-50 shadow-[0_0_20px_rgba(234,179,8,0.3)]'
  }
  return rarityColors[props.achievement.rarity] || rarityColors.common
})

const achievementName = computed(() => {
  if (locale.value === 'zh-CN') {
    return props.achievement.name_zh || props.achievement.name
  }
  return props.achievement.name
})

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

const getThresholdText = () => {
  const { category, threshold } = props.achievement
  switch (category) {
    case 'notes':
      return t('achievementThreshold.notes', { threshold })
    case 'streak':
      return t('achievementThreshold.streak', { threshold })
    case 'words':
      if (threshold >= 1000) {
        return t('achievementThreshold.words', { threshold: (threshold / 1000) + 'k' })
      }
      return t('achievementThreshold.words', { threshold })
    case 'goals':
      return t('achievementThreshold.goals', { threshold })
    default:
      return ''
  }
}
</script>
