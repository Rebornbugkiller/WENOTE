<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { X } from 'lucide-vue-next'
import { tagColors } from '@/utils/color'

const { t } = useI18n()

const props = defineProps({
  show: Boolean,
  type: { type: String, default: 'notebook' }, // 'notebook' | 'tag'
  initialName: { type: String, default: '' },
  initialColor: { type: String, default: '#6B7280' }
})
const emit = defineEmits(['close', 'save'])

const editName = ref('')
const editColor = ref('#6B7280')

watch(() => props.show, (val) => {
  if (val) {
    editName.value = props.initialName
    editColor.value = props.initialColor || '#6B7280'
  }
})

const handleSave = () => {
  if (!editName.value.trim()) return
  emit('save', {
    name: editName.value.trim(),
    color: props.type === 'tag' ? editColor.value : undefined
  })
}

const handleKeydown = (e) => {
  if (e.key === 'Enter') handleSave()
  if (e.key === 'Escape') emit('close')
}
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="$emit('close')" />

        <!-- Modal -->
        <Transition name="scale">
          <div
            v-if="show"
            class="relative bg-white p-6 rounded-2xl border-4 border-black shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] w-full max-w-sm"
          >
            <!-- Header -->
            <div class="flex justify-between items-center mb-4">
              <h3 class="font-black text-xl">
                {{ type === 'notebook' ? t('createModal.editNotebook') : t('createModal.editTag') }}
              </h3>
              <button
                @click="$emit('close')"
                class="p-1 text-slate-400 hover:text-red-500 transition-colors"
              >
                <X class="w-5 h-5" />
              </button>
            </div>

            <!-- Input -->
            <input
              v-model="editName"
              autofocus
              type="text"
              class="w-full border-2 border-slate-200 rounded-xl px-4 py-3 font-bold text-slate-800 focus:outline-none focus:border-black mb-4 transition-colors bg-slate-50 focus:bg-white"
              :placeholder="t('createModal.namePlaceholder')"
              @keydown="handleKeydown"
            />

            <!-- Color Picker for Tags -->
            <div v-if="type === 'tag'" class="mb-6">
              <label class="block text-sm font-bold text-slate-600 mb-2">
                {{ t('createModal.selectColor') }}
              </label>
              <div class="flex gap-2 flex-wrap">
                <button
                  v-for="color in tagColors"
                  :key="color"
                  @click="editColor = color"
                  class="w-10 h-10 rounded-lg border-2 transition-all hover:scale-110"
                  :class="editColor === color ? 'border-black scale-110' : 'border-slate-300'"
                  :style="{ backgroundColor: color }"
                  :title="color"
                />
              </div>
            </div>

            <!-- Actions -->
            <div class="flex justify-end gap-2">
              <button
                @click="$emit('close')"
                class="px-4 py-2 font-bold text-slate-500 hover:bg-slate-100 rounded-lg transition-colors"
              >
                {{ t('common.cancel') }}
              </button>
              <button
                @click="handleSave"
                :disabled="!editName.trim()"
                class="px-6 py-2 bg-black text-white border-2 border-transparent hover:border-black hover:bg-slate-800 font-bold rounded-xl shadow-lg transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ t('common.save') }}
              </button>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>
