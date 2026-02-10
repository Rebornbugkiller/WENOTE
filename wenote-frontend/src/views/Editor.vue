<script setup>
import { ref, onMounted, computed, onBeforeUnmount, watch, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ArrowLeft, Pin, Star, Save, Sparkles, CheckCircle2, Book, Bot } from 'lucide-vue-next'
import { getNote, createNote, updateNote } from '../api/note'
import { getNotebooks, getDefaultNotebook } from '../api/notebook'
import { getTags } from '../api/tag'
import { ElMessage, ElMessageBox } from 'element-plus'
import Vditor from 'vditor'
import 'vditor/dist/index.css'

const router = useRouter()
const route = useRoute()
const { t, locale } = useI18n()

// Get notebook display name (handle default notebook translation)
const getNotebookDisplayName = (notebook) => {
  if (notebook.is_default) {
    return t('sidebar.uncategorized')
  }
  return notebook.name
}

// 判断是否为新建模式
const isNewMode = computed(() => route.name === 'EditorNew' || route.path === '/editor/new')

// Loading state
const isLoading = ref(true)
const editorReady = ref(false)
const isSaved = ref(false) // 标记是否已保存过

// Note data
const noteId = computed(() => route.params.id)
const note = ref(null)
const notebooks = ref([])
const tags = ref([])

// Form data
const formData = ref({
  id: null,
  title: '',
  content: '',
  notebook_id: null,
  is_pinned: false,
  is_starred: false,
  tags: [],
  summary: '',
  suggested_tags: [],
  ai_status: 'pending'
})

// AI loading state
const aiLoading = ref(false)
const showTagSelect = ref(false)

// Vditor instance
const vditor = ref(null)
const editorContainer = ref(null)

// 检查是否有内容
const hasContent = computed(() => {
  const title = formData.value.title?.trim() || ''
  const content = formData.value.content?.trim() || ''
  return title.length > 0 || content.length > 0
})

// Load note data (编辑模式)
const loadNote = async () => {
  isLoading.value = true
  try {
    const [noteData, notebooksData, tagsData] = await Promise.all([
      getNote(noteId.value),
      getNotebooks(),
      getTags()
    ])

    note.value = noteData
    notebooks.value = notebooksData.list || []
    tags.value = tagsData.list || []

    // Set form data
    formData.value = {
      id: noteData.id,
      title: noteData.title || '',
      content: noteData.content || '',
      notebook_id: noteData.notebook_id,
      is_pinned: noteData.is_pinned || false,
      is_starred: noteData.is_starred || false,
      tags: noteData.tags || [],
      summary: noteData.summary || '',
      suggested_tags: noteData.suggested_tags || [],
      ai_status: noteData.ai_status || 'pending'
    }
    isSaved.value = true // 编辑模式下已有笔记
  } catch (err) {
    console.error('Failed to load note:', err)
    ElMessage.error(t('editor.loadNoteFailed'))
    router.push('/')
  } finally {
    isLoading.value = false
  }
}

// Load initial data (新建模式)
const loadInitialData = async () => {
  isLoading.value = true
  try {
    const [notebooksData, tagsData, defaultNotebook] = await Promise.all([
      getNotebooks(),
      getTags(),
      getDefaultNotebook()
    ])

    notebooks.value = notebooksData.list || []
    tags.value = tagsData.list || []

    // 设置默认笔记本
    formData.value.notebook_id = defaultNotebook.id
  } catch (err) {
    console.error('Failed to load initial data:', err)
    ElMessage.error(t('editor.loadDataFailed'))
  } finally {
    isLoading.value = false
  }
}

// Initialize Vditor editor
const initVditor = async () => {
  if (!editorContainer.value) {
    console.error('Editor container not found')
    return
  }

  try {
    vditor.value = new Vditor(editorContainer.value, {
      height: 'calc(100vh - 120px)',
      mode: 'ir',
      cdn: '/vditor',
      placeholder: t('editor.contentPlaceholder'),
      toolbar: [
        'emoji', 'headings', 'bold', 'italic', 'strike', '|',
        'line', 'quote', 'list', 'ordered-list', 'check', '|',
        'code', 'inline-code', 'link', 'table', '|',
        'upload', '|', 'undo', 'redo', '|',
        'preview', 'fullscreen'
      ],
      upload: {
        accept: 'image/*',
        handler: () => {
          ElMessage.warning(t('editor.useImageLink'))
          return null
        }
      },
      cache: { enable: false },
      after: () => {
        editorReady.value = true
        const content = formData.value.content || ''
        if (vditor.value && vditor.value.setValue) {
          vditor.value.setValue(content)
        }
      },
      input: (value) => {
        formData.value.content = value
      }
    })
  } catch (error) {
    console.error('❌ Failed to initialize Vditor:', error)
    vditor.value = null
  }
}

// Load note and init editor on mount
onMounted(async () => {
  if (isNewMode.value) {
    // 新建模式：只加载笔记本和标签
    const dataPromise = loadInitialData()
    await nextTick()
    initVditor()
    await dataPromise
  } else {
    // 编辑模式：加载笔记数据
    const dataPromise = loadNote()
    await nextTick()
    initVditor()
    await dataPromise
    if (vditor.value && editorReady.value) {
      vditor.value.setValue(formData.value.content || '')
    }
  }
})

// Watch for data loaded - update editor content
watch([isLoading, editorReady], ([loading, ready]) => {
  if (!loading && ready && vditor.value && !isNewMode.value) {
    vditor.value.setValue(formData.value.content || '')
  }
})

// Cleanup on unmount
onBeforeUnmount(() => {
  if (vditor.value) {
    vditor.value.destroy()
    vditor.value = null
  }
})

// Available tags
const availableTags = computed(() => {
  const currentTagIds = formData.value.tags?.map(t => t.id) || []
  return tags.value.filter(t => !currentTagIds.includes(t.id))
})

// Navigate back with unsaved check
const goBack = async () => {
  // 新建模式下，如果有内容但未保存，提示用户
  if (isNewMode.value && hasContent.value && !isSaved.value) {
    try {
      await ElMessageBox.confirm(
        t('editor.unsavedConfirm'),
        t('common.hint'),
        {
          confirmButtonText: t('editor.discard'),
          cancelButtonText: t('editor.continueEdit'),
          type: 'warning'
        }
      )
      // 用户确认放弃，直接返回
      router.push('/')
    } catch {
      // 用户取消，继续编辑
    }
  } else {
    router.push('/')
  }
}

// Save note
const handleSave = async () => {
  if (vditor.value) {
    formData.value.content = vditor.value.getValue()
  }

  // 验证：标题和内容至少有一个
  if (!formData.value.title?.trim() && !formData.value.content?.trim()) {
    ElMessage.warning(t('editor.enterTitleOrContent'))
    return
  }

  try {
    if (isNewMode.value && !isSaved.value) {
      // 新建模式：创建笔记
      const newNote = await createNote({
        notebook_id: formData.value.notebook_id,
        title: formData.value.title,
        content: formData.value.content
      })

      formData.value.id = newNote.id
      isSaved.value = true
      ElMessage.success(t('editor.createSuccess'))

      // 替换URL为编辑模式（不产生历史记录）
      router.replace(`/editor/${newNote.id}`)
    } else {
      // 编辑模式：更新笔记
      await updateNote(formData.value.id, {
        title: formData.value.title,
        content: formData.value.content,
        notebook_id: formData.value.notebook_id,
        is_starred: formData.value.is_starred,
        is_pinned: formData.value.is_pinned,
        tag_ids: formData.value.tags?.map(t => t.id) || []
      })
      ElMessage.success(t('editor.saveSuccess'))
    }
  } catch (err) {
    console.error('Save failed:', err)
    ElMessage.error(t('editor.saveFailed'))
  }
}

// AI generate summary and tags
const handleGenerateAI = async () => {
  if (!formData.value.id) {
    ElMessage.warning(t('messages.saveFirst'))
    return
  }

  // Get latest content from editor
  if (vditor.value) {
    formData.value.content = vditor.value.getValue()
  }

  if (!formData.value.content || !formData.value.content.trim()) {
    ElMessage.warning(t('messages.contentEmpty'))
    return
  }

  // Save note first before generating AI
  try {
    await updateNote(formData.value.id, {
      title: formData.value.title,
      content: formData.value.content,
      notebook_id: formData.value.notebook_id,
      is_starred: formData.value.is_starred,
      is_pinned: formData.value.is_pinned,
      tag_ids: formData.value.tags?.map(t => t.id) || []
    })
  } catch (err) {
    console.error('Save failed before AI generation:', err)
    ElMessage.error(t('editor.saveFailedAI'))
    return
  }

  aiLoading.value = true
  try {
    const { generateSummaryAndTags } = await import('../api/note')
    const res = await generateSummaryAndTags(formData.value.id)
    formData.value.summary = res.summary
    formData.value.suggested_tags = res.tags
    formData.value.ai_status = 'done'
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
    const { applySuggestedTags } = await import('../api/note')
    await applySuggestedTags(formData.value.id)

    // Fetch updated note data to get new tags
    const noteData = await getNote(formData.value.id)

    // Only update tags, don't reload everything
    formData.value.tags = noteData.tags || []
    formData.value.suggested_tags = []

    ElMessage.success(t('messages.tagsApplied'))
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
</script>

<template>
  <div class="min-h-screen bg-slate-50">
    <!-- Header -->
    <header class="bg-white border-b-4 border-black sticky top-0 z-10">
      <div class="max-w-6xl mx-auto px-4 md:px-6 py-3 md:py-4 flex items-center justify-between">
        <div class="flex items-center gap-2 md:gap-4">
          <button
            @click="goBack"
            class="px-3 md:px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white font-bold rounded-lg border-2 border-black shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] hover:shadow-[3px_3px_0px_0px_rgba(0,0,0,1)] hover:-translate-y-0.5 transition-all flex items-center gap-2"
            :title="t('editor.back')"
          >
            <ArrowLeft class="w-5 h-5" />
            <span class="hidden md:inline">{{ t('editor.back') }}</span>
          </button>
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
          @click="handleSave"
          class="px-4 md:px-6 py-2 bg-green-500 text-white border-2 border-black rounded-xl font-bold shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] hover:shadow-[3px_3px_0px_0px_rgba(0,0,0,1)] hover:-translate-y-0.5 transition-all active:translate-y-0 active:shadow-none text-sm md:text-base"
        >
          <Save class="w-4 h-4 inline mr-1 md:mr-2" />
          <span class="hidden md:inline">{{ isNewMode && !isSaved ? t('editor.createNote') : t('editor.saveChanges') }}</span>
          <span class="md:hidden">{{ t('common.save') }}</span>
        </button>
      </div>
    </header>

    <!-- Content -->
    <div class="max-w-screen-2xl mx-auto px-4 md:px-6 py-4 md:py-8 flex flex-col lg:flex-row gap-4 lg:gap-8">
      <!-- Editor Area -->
      <div class="flex-1 min-w-0">
        <!-- Title -->
        <div class="mb-4">
          <input
            v-model="formData.title"
            type="text"
            :disabled="isLoading"
            class="w-full bg-transparent text-2xl md:text-3xl lg:text-4xl font-black text-slate-800 placeholder-slate-400 focus:outline-none disabled:opacity-50"
            :placeholder="isLoading ? t('common.loading') : t('editor.titlePlaceholder')"
          />
        </div>

        <!-- Vditor -->
        <div class="bg-white border-2 border-black rounded-xl shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] relative vditor-container" style="overflow: visible;">
          <div ref="editorContainer" class="vditor-wrapper"></div>
          <!-- Loading overlay -->
          <div v-if="isLoading" class="absolute inset-0 bg-white/80 flex items-center justify-center rounded-xl">
            <div class="text-slate-400 font-bold">{{ isNewMode ? t('editor.preparing') : t('editor.loadingNote') }}</div>
          </div>
        </div>
      </div>

      <!-- Sidebar -->
      <div class="w-full lg:w-72 lg:flex-shrink-0 space-y-4">
        <!-- AI Panel (只在编辑模式显示) -->
        <div v-if="!isNewMode || isSaved" class="bg-white p-4 rounded-xl border-2 border-green-200 shadow-sm">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <Bot class="w-4 h-4 text-green-600" />
              <span class="text-xs font-black text-green-700 uppercase">{{ t('editor.aiAssistant') }}</span>
            </div>
            <button
              v-if="formData.ai_status !== 'done'"
              @click="handleGenerateAI"
              :disabled="aiLoading || !formData.content"
              class="flex items-center gap-1 px-3 py-1.5 bg-gradient-to-r from-green-500 to-emerald-500 text-white border-2 border-black rounded-lg font-bold text-xs shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] disabled:opacity-50 disabled:cursor-not-allowed"
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

        <!-- New Mode Hint -->
        <div v-if="isNewMode && !isSaved" class="bg-yellow-50 p-4 rounded-xl border-2 border-yellow-200 shadow-sm">
          <p class="text-sm text-yellow-800 font-bold">
            📝 {{ t('editor.newNoteMode') }}
          </p>
          <p class="text-xs text-yellow-600 mt-1">
            {{ t('editor.newNoteHint') }}
          </p>
        </div>

        <!-- Notebook -->
        <div class="bg-white p-4 rounded-xl border-2 border-slate-200 shadow-sm">
          <label class="text-xs font-black text-slate-400 uppercase mb-2 block">{{ t('editor.notebook') }}</label>
          <div class="relative">
            <select
              v-model="formData.notebook_id"
              class="w-full appearance-none bg-slate-50 border-2 border-slate-200 rounded-xl py-3 pl-4 pr-10 font-bold text-sm text-slate-700 focus:outline-none focus:border-black transition-colors"
            >
              <option v-for="nb in notebooks" :key="nb.id" :value="nb.id">
                {{ getNotebookDisplayName(nb) }}
              </option>
            </select>
            <Book class="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
          </div>
        </div>

        <!-- Tags (只在编辑模式显示) -->
        <div v-if="!isNewMode || isSaved" class="bg-white p-4 rounded-xl border-2 border-slate-200 shadow-sm">
          <label class="text-xs font-black text-slate-400 uppercase mb-2 block">{{ t('editor.tags') }}</label>
          <div class="flex flex-wrap gap-2 p-3 bg-slate-50 border-2 border-slate-200 rounded-xl min-h-[80px] content-start">
            <span
              v-for="tag in formData.tags"
              :key="tag.id"
              class="px-2 py-1 bg-white border border-slate-200 text-xs font-bold rounded-md flex items-center gap-1 group shadow-sm"
            >
              {{ tag.name }}
              <button
                @click="removeTag(tag.id)"
                class="w-3 h-3 cursor-pointer text-slate-400 hover:text-red-500"
              >
                ×
              </button>
            </span>

            <!-- Add Tag Button -->
            <div class="relative">
              <button
                @click="showTagSelect = !showTagSelect"
                class="px-2 py-1 bg-slate-100 border border-dashed border-slate-300 text-xs font-bold rounded-md text-slate-400 hover:text-green-500 hover:border-green-500 transition-colors"
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
    </div>
  </div>
</template>

<style scoped>
/* 移动端编辑器高度 */
.vditor-container :deep(.vditor) {
  height: 50vh !important;
}

@media (min-width: 768px) {
  .vditor-container :deep(.vditor) {
    height: calc(100vh - 120px) !important;
  }
}

.vditor-wrapper :deep(.vditor) {
  border: none;
  height: 100%;
}

.vditor-wrapper :deep(.vditor-toolbar) {
  background-color: #f8fafc;
  border-bottom: 2px solid #e2e8f0;
  padding: 8px;
  position: relative;
  z-index: 1000;
}

/* 确保工具栏按钮的 tooltip 显示在最上层 */
.vditor-wrapper :deep(.vditor-toolbar .vditor-tooltipped) {
  position: relative;
  z-index: 1001;
}

.vditor-wrapper :deep(.vditor-toolbar .vditor-tooltipped::after),
.vditor-wrapper :deep(.vditor-toolbar .vditor-tooltipped::before) {
  z-index: 10000 !important;
}

.vditor-wrapper :deep(.vditor-hint) {
  z-index: 10000 !important;
}

.vditor-wrapper :deep(.vditor-panel) {
  z-index: 10000 !important;
}

.vditor-wrapper :deep(.vditor-tip) {
  z-index: 10000 !important;
}

/* 确保内容区域不会盖住 tooltip */
.vditor-wrapper :deep(.vditor-content) {
  background-color: white;
  font-size: 16px;
  line-height: 1.8;
  position: relative;
  z-index: 1;
}

.vditor-wrapper :deep(.vditor-ir) {
  padding: 20px;
  position: relative;
  z-index: 1;
}
</style>
