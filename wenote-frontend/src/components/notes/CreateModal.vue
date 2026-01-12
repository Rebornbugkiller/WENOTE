<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { X } from 'lucide-vue-next'

const { t } = useI18n()

const props = defineProps({
  show: Boolean,
  type: { type: String, default: 'notebook' } // 'notebook' | 'tag'
})
const emit = defineEmits(['close', 'create'])

const inputName = ref('')

watch(() => props.show, (val) => {
  if (val) inputName.value = ''
})

const handleCreate = () => {
  if (!inputName.value.trim()) return
  emit('create', inputName.value.trim())
  inputName.value = ''
}

const handleKeydown = (e) => {
  if (e.key === 'Enter') handleCreate()
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
                {{ type === 'notebook' ? t('createModal.createNotebook') : t('createModal.createTag') }}
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
              v-model="inputName"
              autofocus
              type="text"
              class="w-full border-2 border-slate-200 rounded-xl px-4 py-3 font-bold text-slate-800 focus:outline-none focus:border-black mb-6 transition-colors bg-slate-50 focus:bg-white"
              :placeholder="type === 'notebook' ? t('createModal.notebookPlaceholder') : t('createModal.tagPlaceholder')"
              @keydown="handleKeydown"
            />

            <!-- Actions -->
            <div class="flex justify-end gap-2">
              <button
                @click="$emit('close')"
                class="px-4 py-2 font-bold text-slate-500 hover:bg-slate-100 rounded-lg transition-colors"
              >
                {{ t('common.cancel') }}
              </button>
              <button
                @click="handleCreate"
                class="px-6 py-2 bg-black text-white border-2 border-transparent hover:border-black hover:bg-slate-800 font-bold rounded-xl shadow-lg transition-all active:scale-95"
              >
                {{ t('common.create') }}
              </button>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>
