<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Star, Pin, Trash2, Sparkles } from 'lucide-vue-next'
import { getTagColor } from '../../utils/color'

const { t, locale } = useI18n()

const props = defineProps({
  note: { type: Object, required: true },
  isTrash: { type: Boolean, default: false },
  selected: { type: Boolean, default: false },
  index: { type: Number, default: 0 }
})

const emit = defineEmits(['click', 'delete', 'restore', 'toggle-status', 'toggle-select'])

// 3D 倾斜状态
const cardRef = ref(null)
const tiltStyle = ref({})
const glowStyle = ref({ opacity: 0 })

const handleMouseMove = (e) => {
  if (props.isTrash) return
  const card = cardRef.value
  if (!card) return
  const rect = card.getBoundingClientRect()
  const x = e.clientX - rect.left
  const y = e.clientY - rect.top
  const centerX = rect.width / 2
  const centerY = rect.height / 2
  const rotateX = (y - centerY) / 15
  const rotateY = (centerX - x) / 15
  tiltStyle.value = {
    transform: `perspective(1000px) rotateX(${rotateX}deg) rotateY(${rotateY}deg) scale(1.02)`
  }
  glowStyle.value = {
    opacity: 0.15,
    background: `radial-gradient(circle at ${x}px ${y}px, rgba(34,197,94,0.4) 0%, transparent 60%)`
  }
}

const handleMouseLeave = () => {
  tiltStyle.value = {}
  glowStyle.value = { opacity: 0 }
}

// 合并样式
const cardStyle = computed(() => ({
  ...tiltStyle.value,
  animationDelay: `${props.index * 0.05}s`
}))

const formatDate = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString(locale.value === 'zh-CN' ? 'zh-CN' : 'en-US')
}
</script>

<template>
  <div
    ref="cardRef"
    class="note-card relative group p-6 rounded-3xl border-4 border-black bg-white flex flex-col h-72 justify-between"
    :class="[
      !isTrash ? 'cursor-pointer hover:z-10' : 'cursor-default opacity-80',
      'shadow-[5px_5px_0px_0px_rgba(0,0,0,0.1)] hover:shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]',
      selected ? 'ring-4 ring-green-500' : ''
    ]"
    :style="cardStyle"
    @click="!isTrash && emit('click', note)"
    @mousemove="handleMouseMove"
    @mouseleave="handleMouseLeave"
  >
    <!-- 光标跟随高光 -->
    <div class="absolute inset-0 rounded-3xl pointer-events-none transition-opacity duration-300" :style="glowStyle" />

    <!-- 回收站复选框 -->
    <input
      v-if="isTrash"
      type="checkbox"
      :checked="selected"
      @click.stop
      @change="emit('toggle-select', note.id)"
      class="absolute top-3 left-3 w-5 h-5 accent-green-500 z-20 cursor-pointer"
    />
    <!-- Tape decoration -->
    <div class="absolute -top-3 left-1/2 -translate-x-1/2 w-16 h-6 bg-white/30 backdrop-blur-sm border-2 border-white/50 rotate-[-2deg] shadow-sm pointer-events-none" />

    <!-- Pin button (top right) -->
    <button
      v-if="!isTrash"
      @click.stop="emit('toggle-status', note.id, 'is_pinned')"
      class="absolute -right-3 -top-3 p-1.5 rounded-full border-2 border-black shadow-sm z-20 transition-all"
      :class="note.is_pinned ? 'bg-blue-500 text-white opacity-100' : 'bg-white text-slate-300 hover:text-blue-500 opacity-0 group-hover:opacity-100'"
      :title="note.is_pinned ? t('editor.unpin') : t('editor.pin')"
    >
      <Pin class="w-4 h-4" :fill="note.is_pinned ? 'currentColor' : 'none'" />
    </button>

    <!-- Pinned indicator for trash -->
    <div
      v-if="isTrash && note.is_pinned"
      class="absolute -right-3 -top-3 bg-blue-500 text-white p-1.5 rounded-full border-2 border-black shadow-sm z-20"
    >
      <Pin class="w-4 h-4 fill-current" />
    </div>

    <!-- Content -->
    <div>
      <div class="flex justify-between items-start mb-3">
        <h3 class="font-black text-xl line-clamp-2 leading-tight flex-1 mr-2">
          {{ note.title || t('noteCard.untitled') }}
        </h3>
        <button
          v-if="!isTrash"
          @click.stop="emit('toggle-status', note.id, 'is_starred')"
          class="p-1 rounded-full transition-colors hover:bg-slate-100"
          :class="note.is_starred ? 'text-yellow-400' : 'text-slate-300'"
        >
          <Star class="w-5 h-5" fill="currentColor" />
        </button>
      </div>

      <!-- AI Summary -->
      <div
        v-if="note.summary"
        class="mb-2 px-2 py-1 bg-green-50 border border-green-100 rounded text-[10px] font-bold text-green-700 flex items-start gap-1"
      >
        <Sparkles class="w-3 h-3 mt-0.5 shrink-0" />
        <span class="line-clamp-2">{{ note.summary }}</span>
      </div>

      <!-- Content preview -->
      <p v-if="!note.summary" class="text-slate-500 font-medium text-sm line-clamp-3 leading-relaxed">
        {{ note.content || t('noteCard.noContent') }}
      </p>
    </div>

    <!-- Footer -->
    <div class="mt-4">
      <!-- Tags -->
      <div class="flex flex-wrap gap-1 mb-3">
        <span
          v-for="tag in note.tags"
          :key="tag.id"
          class="text-[10px] font-bold px-1.5 py-0.5 rounded border border-black/10 text-slate-600"
          :style="{ backgroundColor: getTagColor(tag) + '40' }"
        >
          #{{ tag.name }}
        </span>
      </div>

      <!-- Meta info -->
      <div class="flex justify-between items-center border-t border-black/5 pt-3">
        <div class="flex items-center gap-2">
          <span class="text-[10px] font-bold text-slate-400">
            {{ formatDate(note.updated_at) }}
          </span>
        </div>

        <!-- Actions (非回收站) -->
        <button
          v-if="!isTrash"
          @click.stop="emit('delete', note.id)"
          class="p-1.5 rounded-lg bg-red-100 border-2 border-red-300 text-red-500 hover:bg-red-500 hover:text-white transition-all"
          :title="t('noteCard.moveToTrash')"
        >
          <Trash2 class="w-4 h-4" />
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.note-card {
  transition: transform 0.15s ease-out, box-shadow 0.2s ease;
  animation: fadeUp 0.4s ease-out both;
}

@keyframes fadeUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
