<script setup>
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { X, Pin, Star, Save, Sparkles, CheckCircle2, Book, Bot } from 'lucide-vue-next'
import { applySuggestedTags, generateSummaryAndTags, getNote } from '../../api/note'
import { ElMessage } from 'element-plus'

const { t } = useI18n()

const props = defineProps({
  note: { type: Object, required: true },
  notebooks: { type: Array, default: () => [] },
  allTags: { type: Array, default: () => [] }
})

const emit = defineEmits(['close', 'save', 'update:note', 'refresh'])

// Form data
const formData = ref({ ...props.note })
const aiLoading = ref(false)
const showTagSelect = ref(false)

// Watch for note changes - only reset when switching to a different note
watch(() => props.note.id, () => {
  formData.value = { ...props.note }
})

// Available tags (not already added)
const availableTags = computed(() => {
  const currentTagIds = formData.value.tags?.map(t => t.id) || []
  return props.allTags.filter(t => !currentTagIds.includes(t.id))
})

// AI 生成摘要和标签
const handleGenerateAI = async () => {
  if (!formData.value.id) {
    ElMessage.warning(t('messages.saveFirst'))
    return
  }
  if (!formData.value.content) {
    ElMessage.warning(t('messages.contentEmpty'))
    return
  }

  aiLoading.value = true
  try {
    const res = await generateSummaryAndTags(formData.value.id)
    formData.value.summary = res.summary
    formData.value.suggested_tags = res.tags
    formData.value.ai_status = 'done'
    emit('refresh') // 只刷新列表，不关闭弹窗
    ElMessage.success(t('messages.aiGenerateSuccess'))
  } catch (err) {
    ElMessage.error(err.response?.data?.message || t('messages.aiGenerateFailed'))
  } finally {
    aiLoading.value = false
  }
}

// Apply AI suggested tags
const handleApplySuggestedTags = async () => {
  if (!formData.value.suggested_tags?.length) return

  try {
    await applySuggestedTags(formData.value.id)
    // 重新获取笔记数据，更新标签列表
    const updatedNote = await getNote(formData.value.id)
    formData.value.tags = updatedNote.tags || []
    formData.value.suggested_tags = []
    ElMessage.success(t('messages.tagsApplied'))
    // 触发父组件刷新
    emit('save', formData.value)
  } catch (err) {
    ElMessage.error(t('messages.applyTagsFailed'))
  }
}

// Add tag to note
const addTagToNote = (tag) => {
  if (!formData.value.tags) {
    formData.value.tags = []
  }
  if (!formData.value.tags.find(t => t.id === tag.id)) {
    formData.value.tags.push(tag)
  }
  showTagSelect.value = false
}

// Remove tag from note
const removeTag = (tagId) => {
  formData.value.tags = formData.value.tags.filter(t => t.id !== tagId)
}

// Save handler
const handleSave = () => {
  emit('save', formData.value)
}
</script>

<template>
  <Teleport to="body">
    <div v-if="note" class="fixed inset-0 z-50 flex items-center justify-center p-4 md:p-10">
      <!-- Backdrop -->
      <Transition name="fade">
        <div v-if="note" class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="$emit('close')" />
      </Transition>

      <!-- Modal -->
      <Transition name="scale" appear>
        <div
          v-if="note"
          class="relative w-full max-w-5xl h-[85vh] bg-white border-4 border-black rounded-3xl shadow-[12px_12px_0px_0px_rgba(0,0,0,1)] flex flex-col overflow-hidden"
        >
          <!-- Header -->
          <div class="p-4 border-b-4 border-black bg-slate-50 flex justify-between items-center select-none">
            <div class="flex items-center gap-4">
              <span class="font-black bg-black text-white px-2 py-1 rounded text-xs uppercase">{{ t('editor.title') }}</span>
              <div class="flex gap-2">
                <button
                  @click="formData.is_pinned = !formData.is_pinned"
                  class="p-1.5 rounded-lg border-2 border-transparent hover:bg-slate-200 transition-colors"
                  :class="formData.is_pinned ? 'text-green-600 bg-green-100 border-green-200' : 'text-slate-400'"
                  :title="formData.is_pinned ? t('editor.unpin') : t('editor.pin')"
                >
                  <Pin class="w-4 h-4" :fill="formData.is_pinned ? 'currentColor' : 'none'" />
                </button>
                <button
                  @click="formData.is_starred = !formData.is_starred"
                  class="p-1.5 rounded-lg border-2 border-transparent hover:bg-slate-200 transition-colors"
                  :class="formData.is_starred ? 'text-yellow-500 bg-yellow-50 border-yellow-200' : 'text-slate-400'"
                  :title="t('editor.star')"
                >
                  <Star class="w-4 h-4" :fill="formData.is_starred ? 'currentColor' : 'none'" />
                </button>
              </div>
            </div>
            <button
              @click="$emit('close')"
              class="p-2 bg-red-500 text-white border-2 border-black hover:bg-red-600 rounded-lg transition-colors"
            >
              <X class="w-5 h-5" />
            </button>
          </div>

          <!-- Content -->
          <div class="flex-1 flex flex-col md:flex-row overflow-hidden">
            <!-- Editor Area -->
            <div class="flex-1 p-8 overflow-y-auto flex flex-col">
              <input
                v-model="formData.title"
                type="text"
                class="w-full bg-transparent text-4xl font-black text-slate-800 placeholder-slate-400 focus:outline-none mb-6"
                :placeholder="t('editor.titlePlaceholder')"
              />
              <textarea
                v-model="formData.content"
                class="w-full flex-1 min-h-[400px] bg-transparent text-xl font-medium text-slate-600 placeholder-slate-400 focus:outline-none resize-none leading-relaxed"
                :placeholder="t('editor.contentPlaceholder')"
              />
            </div>

            <!-- Sidebar -->
            <div class="w-full md:w-80 bg-slate-50 border-l-4 border-black flex flex-col overflow-hidden">
              <div class="flex-1 overflow-y-auto p-6 space-y-6">
                <!-- AI Insight -->
                <div class="bg-white p-4 rounded-xl border-2 border-green-200 shadow-sm">
                  <div class="flex items-center justify-between mb-3">
                    <div class="flex items-center gap-2">
                      <Bot class="w-4 h-4 text-green-600" />
                      <span class="text-xs font-black text-green-700 uppercase">{{ t('editor.aiAssistant') }}</span>
                    </div>
                    <button
                      v-if="formData.ai_status !== 'done'"
                      @click="handleGenerateAI"
                      :disabled="aiLoading || !formData.content"
                      class="flex items-center gap-1 px-3 py-1.5 bg-gradient-to-r from-green-500 to-emerald-500 text-white border-2 border-black rounded-lg font-bold text-xs shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] hover:shadow-[3px_3px_0px_0px_rgba(0,0,0,1)] hover:-translate-y-0.5 transition-all disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:translate-y-0 disabled:hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,1)]"
                    >
                      <Sparkles class="w-3 h-3" :class="aiLoading ? 'animate-spin' : ''" />
                      {{ aiLoading ? t('editor.generating') : t('editor.generate') }}
                    </button>
                    <span v-else class="text-xs text-green-600 font-bold flex items-center gap-1">
                      <CheckCircle2 class="w-3 h-3" /> {{ t('editor.generated') }}
                    </span>
                  </div>

                  <!-- Summary -->
                  <div class="mb-4">
                    <p class="text-[10px] font-bold text-slate-400 uppercase mb-1">{{ t('editor.summary') }}</p>
                    <p v-if="formData.summary" class="text-sm text-slate-700 bg-green-50/50 p-2 rounded border border-green-100 leading-snug">
                      {{ formData.summary }}
                    </p>
                    <p v-else class="text-sm text-slate-400 italic text-center py-2">
                      {{ t('editor.summaryHint') }}
                    </p>
                  </div>

                  <!-- Suggested Tags -->
                  <div v-if="formData.suggested_tags?.length">
                    <p class="text-[10px] font-bold text-slate-400 uppercase mb-1">{{ t('editor.suggestedTags') }}</p>
                    <div class="flex flex-wrap gap-1 mb-2">
                      <span
                        v-for="tag in formData.suggested_tags"
                        :key="tag"
                        class="px-2 py-1 bg-yellow-100 border border-yellow-200 text-yellow-800 text-xs font-bold rounded-md"
                      >
                        {{ tag }}
                      </span>
                    </div>
                    <button
                      @click="handleApplySuggestedTags"
                      class="w-full py-1.5 border-2 border-black bg-black text-white text-xs font-bold rounded hover:bg-slate-800 transition-colors flex items-center justify-center gap-1"
                    >
                      <CheckCircle2 class="w-3 h-3" /> {{ t('editor.applySuggestedTags') }}
                    </button>
                  </div>
                </div>

                <!-- Notebook -->
                <div>
                  <label class="text-xs font-black text-slate-400 uppercase mb-2 block">{{ t('editor.notebook') }}</label>
                  <div class="relative">
                    <select
                      v-model="formData.notebook_id"
                      class="w-full appearance-none bg-white border-2 border-slate-200 rounded-xl py-3 pl-4 pr-10 font-bold text-sm text-slate-700 focus:outline-none focus:border-black transition-colors"
                    >
                      <option v-for="nb in notebooks" :key="nb.id" :value="nb.id">
                        {{ nb.name }}
                      </option>
                    </select>
                    <Book class="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                  </div>
                </div>

                <!-- Tags -->
                <div>
                  <label class="text-xs font-black text-slate-400 uppercase mb-2 block">{{ t('editor.tags') }}</label>
                  <div class="flex flex-wrap gap-2 p-3 bg-white border-2 border-slate-200 rounded-xl min-h-[80px] content-start">
                    <span
                      v-for="tag in formData.tags"
                      :key="tag.id"
                      class="px-2 py-1 bg-slate-100 border border-slate-200 text-xs font-bold rounded-md flex items-center gap-1 group"
                    >
                      {{ tag.name }}
                      <X
                        class="w-3 h-3 cursor-pointer text-slate-400 hover:text-red-500"
                        @click="removeTag(tag.id)"
                      />
                    </span>

                    <!-- Add Tag Button -->
                    <div class="relative">
                      <button
                        @click="showTagSelect = !showTagSelect"
                        class="px-2 py-1 bg-slate-50 border border-dashed border-slate-300 text-xs font-bold rounded-md text-slate-400 hover:text-green-500 hover:border-green-500 transition-colors"
                      >
                        {{ t('editor.addTag') }}
                      </button>

                      <!-- Tag Dropdown -->
                      <div
                        v-if="showTagSelect"
                        class="absolute top-full left-0 mt-2 w-48 bg-white border-2 border-black rounded-lg shadow-lg z-50 max-h-40 overflow-y-auto"
                      >
                        <template v-if="availableTags.length">
                          <button
                            v-for="tag in availableTags"
                            :key="tag.id"
                            @click="addTagToNote(tag)"
                            class="w-full text-left px-3 py-2 text-xs font-bold hover:bg-slate-100 flex items-center gap-2"
                          >
                            <span
                              class="w-2 h-2 rounded-full"
                              :style="{ backgroundColor: tag.color || '#ccc' }"
                            />
                            {{ tag.name }}
                          </button>
                        </template>
                        <div v-else class="px-3 py-2 text-xs text-slate-400 italic">
                          {{ t('editor.noMoreTags') }}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Save Button -->
              <div class="p-6 bg-white border-t-4 border-black">
                <button
                  @click="handleSave"
                  class="w-full flex items-center justify-center gap-2 px-6 py-3 bg-green-500 text-white border-2 border-black rounded-xl font-black text-lg shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] hover:-translate-y-1 transition-all"
                >
                  <Save class="w-5 h-5" />
                  {{ t('editor.saveChanges') }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </div>
  </Teleport>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.4s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.scale-enter-active { transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1); }
.scale-leave-active { transition: all 0.25s ease-in; }
.scale-enter-from { opacity: 0; transform: scale(0.5); }
.scale-leave-to { opacity: 0; transform: scale(0.9); }
</style>
