<template>
  <div class="daily-goal bg-white border-4 border-black rounded-2xl p-4 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
    <div class="flex justify-between items-center mb-2">
      <div class="text-xs font-black text-slate-400 uppercase tracking-wider">{{ t('dailyGoal.title') }}</div>
      <button @click="showSettings = true" class="text-slate-400 hover:text-slate-600 transition-colors">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
        </svg>
      </button>
    </div>

    <div class="flex items-end gap-2 mb-3">
      <span class="text-3xl font-black" :class="completed ? 'text-green-500' : 'text-slate-800'">
        {{ current.toLocaleString() }}
      </span>
      <span class="text-slate-400 font-bold text-sm pb-1">/ {{ target.toLocaleString() }} {{ t('gamification.chars') }}</span>
    </div>

    <!-- Progress bar -->
    <div class="h-4 bg-slate-200 rounded-full overflow-hidden border-2 border-black">
      <div
        class="h-full transition-all duration-500 ease-out"
        :class="completed ? 'bg-green-500' : 'bg-blue-500'"
        :style="{ width: `${Math.min(percent, 100)}%` }"
      >
        <div v-if="completed" class="h-full bg-gradient-to-r from-green-400 to-green-600 animate-pulse"></div>
      </div>
    </div>

    <!-- Completion celebration -->
    <div v-if="completed" class="mt-3 text-center">
      <span class="text-2xl animate-bounce inline-block">🎉</span>
      <span class="text-green-600 font-black text-sm ml-2">{{ t('dailyGoal.goalAchieved') }}</span>
    </div>

    <!-- Settings Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showSettings" class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 backdrop-blur-sm" @click.self="showSettings = false">
          <div class="bg-white border-4 border-black rounded-2xl p-6 w-[90vw] max-w-sm shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
            <h3 class="text-xl font-black mb-4">{{ t('dailyGoal.setGoal') }}</h3>
            <div class="mb-4">
              <label class="block text-sm font-bold text-slate-600 mb-2">{{ t('dailyGoal.dailyCharGoal') }}</label>
              <input
                v-model.number="newGoal"
                type="number"
                min="100"
                max="10000"
                step="100"
                class="w-full px-4 py-3 border-2 border-black rounded-xl font-bold text-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <p class="text-xs text-slate-400 mt-1">{{ t('dailyGoal.range') }}</p>
            </div>
            <div class="flex gap-3">
              <button
                @click="showSettings = false"
                class="flex-1 py-3 border-2 border-black rounded-xl font-bold hover:bg-slate-100 transition-colors"
              >
                {{ t('common.cancel') }}
              </button>
              <button
                @click="saveGoal"
                class="flex-1 py-3 bg-black text-white border-2 border-black rounded-xl font-bold hover:bg-slate-800 transition-colors"
              >
                {{ t('common.save') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  current: {
    type: Number,
    default: 0
  },
  target: {
    type: Number,
    default: 500
  },
  percent: {
    type: Number,
    default: 0
  },
  completed: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update-goal'])

const showSettings = ref(false)
const newGoal = ref(props.target)

watch(() => props.target, (val) => {
  newGoal.value = val
})

const saveGoal = () => {
  if (newGoal.value >= 100 && newGoal.value <= 10000) {
    emit('update-goal', newGoal.value)
    showSettings.value = false
  }
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .bg-white,
.modal-leave-to .bg-white {
  transform: scale(0.9);
}
</style>
