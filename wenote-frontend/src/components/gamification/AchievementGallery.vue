<template>
  <div class="achievement-gallery bg-white border-4 border-black rounded-2xl p-4 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
    <div class="flex justify-between items-center mb-4">
      <div class="text-xs font-black text-slate-400 uppercase tracking-wider">成就</div>
      <div class="text-sm font-bold text-slate-600">
        {{ unlockedCount }} / {{ totalCount }}
      </div>
    </div>

    <!-- Category Tabs -->
    <div class="flex gap-2 mb-4 overflow-x-auto pb-2">
      <button
        v-for="cat in categories"
        :key="cat.id"
        @click="selectedCategory = cat.id"
        class="px-3 py-1 rounded-lg text-xs font-bold whitespace-nowrap transition-all border-2"
        :class="selectedCategory === cat.id
          ? 'bg-black text-white border-black'
          : 'bg-white text-slate-600 border-slate-200 hover:border-slate-400'"
      >
        {{ cat.label }}
      </button>
    </div>

    <!-- Achievement Grid -->
    <div class="grid grid-cols-4 gap-2">
      <AchievementBadge
        v-for="achievement in filteredAchievements"
        :key="achievement.id"
        :achievement="achievement"
        :unlocked="achievement.unlocked"
        @click="selectedAchievement = achievement"
      />
    </div>

    <!-- Empty State -->
    <div v-if="filteredAchievements.length === 0" class="text-center py-8 text-slate-400">
      <div class="text-4xl mb-2">🏆</div>
      <div class="text-sm">暂无成就</div>
    </div>

    <!-- Achievement Detail Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="selectedAchievement" class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 backdrop-blur-sm" @click.self="selectedAchievement = null">
          <div class="bg-white border-4 border-black rounded-2xl p-6 w-[90vw] max-w-sm shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
            <div class="text-center">
              <div class="text-6xl mb-4" :class="{ 'grayscale opacity-50': !selectedAchievement.unlocked }">
                {{ selectedAchievement.icon }}
              </div>
              <div class="font-black text-xl mb-2">
                {{ selectedAchievement.name_zh || selectedAchievement.name }}
              </div>
              <div class="text-slate-500 text-sm mb-4">
                {{ selectedAchievement.description_zh || selectedAchievement.description }}
              </div>
              <div v-if="selectedAchievement.unlocked" class="text-green-600 font-bold text-sm">
                ✓ 已解锁
              </div>
              <div v-else class="text-slate-400 text-sm">
                🔒 未解锁
              </div>
            </div>
            <button
              @click="selectedAchievement = null"
              class="mt-4 w-full py-2 border-2 border-black rounded-xl font-bold hover:bg-slate-100 transition-colors"
            >
              关闭
            </button>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import AchievementBadge from './AchievementBadge.vue'

const props = defineProps({
  achievements: {
    type: Array,
    default: () => []
  }
})

const selectedCategory = ref('all')
const selectedAchievement = ref(null)

const categories = [
  { id: 'all', label: '全部' },
  { id: 'notes', label: '笔记' },
  { id: 'streak', label: '连续' },
  { id: 'words', label: '字数' },
  { id: 'goals', label: '目标' }
]

const filteredAchievements = computed(() => {
  if (selectedCategory.value === 'all') {
    return props.achievements
  }
  return props.achievements.filter(a => a.category === selectedCategory.value)
})

const unlockedCount = computed(() => props.achievements.filter(a => a.unlocked).length)
const totalCount = computed(() => props.achievements.length)
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
</style>
