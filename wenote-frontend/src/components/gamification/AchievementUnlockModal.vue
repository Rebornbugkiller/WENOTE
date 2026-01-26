<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="show && achievement" class="fixed inset-0 z-[200] flex items-center justify-center bg-black/80 backdrop-blur-sm">
        <div :class="rarityGlow" class="p-1 rounded-[2.5rem]">
          <div class="bg-white border-4 border-black rounded-[2.3rem] p-8 w-[90vw] max-w-md text-center">
            <!-- Achievement Icon -->
            <div class="text-7xl mb-4 animate-bounce">{{ achievement.icon }}</div>

            <!-- Title -->
            <div class="font-black text-2xl text-slate-800 mb-2">
              成就解锁！
            </div>

            <!-- Achievement Name -->
            <div class="font-black text-xl" :class="rarityTextColor">
              {{ achievement.name_zh || achievement.name }}
            </div>

            <!-- Description -->
            <div class="text-slate-500 text-sm mt-2 mb-6">
              {{ achievement.description_zh || achievement.description }}
            </div>

            <!-- Rarity Badge -->
            <div class="inline-block px-4 py-1 rounded-full text-xs font-black uppercase tracking-wider" :class="rarityBadgeClass">
              {{ rarityText }}
            </div>

            <!-- Close Button -->
            <button
              @click="handleClose"
              class="mt-6 w-full py-3 bg-black text-white font-black rounded-xl border-2 border-transparent shadow-[4px_4px_0px_0px_rgba(34,197,94,1)] hover:shadow-[2px_2px_0px_0px_rgba(34,197,94,1)] hover:translate-x-[2px] hover:translate-y-[2px] transition-all"
            >
              太棒了！
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  achievement: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close'])

const rarityGlow = computed(() => {
  const glows = {
    common: 'bg-gradient-to-r from-slate-300 to-slate-400',
    rare: 'bg-gradient-to-r from-blue-400 to-blue-600',
    epic: 'bg-gradient-to-r from-purple-400 to-purple-600',
    legendary: 'bg-gradient-to-r from-yellow-400 to-orange-500 animate-pulse'
  }
  return glows[props.achievement?.rarity] || glows.common
})

const rarityTextColor = computed(() => {
  const colors = {
    common: 'text-slate-600',
    rare: 'text-blue-600',
    epic: 'text-purple-600',
    legendary: 'text-yellow-600'
  }
  return colors[props.achievement?.rarity] || colors.common
})

const rarityBadgeClass = computed(() => {
  const classes = {
    common: 'bg-slate-200 text-slate-600',
    rare: 'bg-blue-100 text-blue-600',
    epic: 'bg-purple-100 text-purple-600',
    legendary: 'bg-yellow-100 text-yellow-600'
  }
  return classes[props.achievement?.rarity] || classes.common
})

const rarityText = computed(() => {
  const texts = {
    common: '普通',
    rare: '稀有',
    epic: '史诗',
    legendary: '传说'
  }
  return texts[props.achievement?.rarity] || '普通'
})

const handleClose = () => {
  emit('close', props.achievement?.id)
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: all 0.4s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from > div,
.modal-leave-to > div {
  transform: scale(0.8) rotate(-5deg);
}
</style>
