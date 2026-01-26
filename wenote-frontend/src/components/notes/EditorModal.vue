<script setup>
import { ref, watch, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { X, Pin, Star, Save, Sparkles, CheckCircle2, Book, Bot } from 'lucide-vue-next'
import { applySuggestedTags, generateSummaryAndTags, getNote, uploadImage, aiAssist } from '../../api/note'
import { ElMessage, ElLoading } from 'element-plus'
import Vditor from 'vditor'
import 'vditor/dist/index.css'

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

// Vditor instance
const vditor = ref(null)
const editorContainer = ref(null)

// Watch for note changes
watch(() => props.note, (newNote, oldNote) => {
  console.log('🔍 EditorModal watch triggered')
  console.log('📦 newNote.id:', newNote?.id)
  console.log('📦 oldNote.id:', oldNote?.id)
  console.log('🔍 vditor.value:', vditor.value)

  if (!newNote) {
    console.log('⚠️ newNote is null, skipping')
    return
  }

  // 只有当笔记 ID 改变时才更新（避免 deep: true 导致的循环）
  if (oldNote && newNote.id === oldNote.id) {
    console.log('⚠️ Same note, skipping update to avoid overwriting user input')
    return
  }

  formData.value = { ...newNote }
  console.log('✅ formData updated:', formData.value.content?.substring(0, 50) || '(空)')

  // 延迟设置 Vditor 内容，确保实例已准备好
  if (vditor.value) {
    setTimeout(() => {
      if (vditor.value) {
        console.log('📝 Setting Vditor content:', newNote.content?.substring(0, 50) || '(空)')
        vditor.value.setValue(newNote.content || '')
        console.log('✅ Vditor content set successfully')
      }
    }, 50)
  } else {
    console.log('⚠️ vditor not ready yet')
  }
}, { immediate: true })

// Initialize Vditor editor
onMounted(async () => {
  await nextTick()
  
  // 确保formData有初始值
  formData.value = { ...props.note }
  
  console.log('🔧 Vditor初始化开始')
  console.log('📝 Note内容:', props.note.content?.substring(0, 100) || '(空)')
  console.log('📦 Container:', editorContainer.value)
  
  if (editorContainer.value) {
    console.log('🚀 开始初始化 Vditor')
    vditor.value = new Vditor(editorContainer.value, {
      height: '100%',
      width: '100%',
      placeholder: t('editor.contentPlaceholder'),
      theme: 'classic',
      mode: 'ir', // 即时渲染模式
      cdn: '/vditor', // 使用本地打包的资源（不加 dist）
      lang: 'zh_CN', // 设置语言
      className: 'vditor-content',
      toolbar: [
        'emoji',
        'headings',
        'bold',
        'italic',
        'strike',
        '|',
        'line',
        'quote',
        'list',
        'ordered-list',
        'check',
        '|',
        'code',
        'inline-code',
        'link',
        'table',
        '|',
        'upload', // 图片上传按钮
        '|',
        'undo',
        'redo',
        '|',
        'preview',
        'fullscreen',
        {
          name: 'more',
          toolbar: [
            'both',
            'code-theme',
            'content-theme',
            'export',
            'outline',
            'preview-theme',
          ],
        }
      ],
      upload: {
        accept: 'image/*',
        max: 5 * 1024 * 1024, // 5MB
        handler: async (files) => {
          if (!formData.value.id) {
            ElMessage.warning(t('messages.saveFirst'))
            return null
          }
          
          const results = []
          for (const file of files) {
            try {
              const response = await uploadImage(formData.value.id, file)
              results.push({
                url: response.url
              })
            } catch (error) {
              ElMessage.error(`图片上传失败: ${file.name}`)
              console.error('Upload error:', error)
            }
          }
          return results.length > 0 ? JSON.stringify(results) : null
        }
      },
      cache: {
        enable: false
      },
      after: () => {
        // 编辑器初始化完成后设置内容
        console.log('✅ Vditor初始化完成')
        console.log('📄 props.note.content:', props.note.content?.substring(0, 100) || '(空)')
        console.log('📄 formData.value.content:', formData.value.content?.substring(0, 100) || '(空)')
        console.log('🔍 vditor实例:', vditor.value)

        if (vditor.value) {
          // 优先使用 formData.value（因为 watch 可能已经更新了）
          const content = formData.value.content || props.note.content || ''
          console.log('🔍 实际内容长度:', content.length)
          console.log('🔍 实际内容:', content.substring(0, 200))
          try {
            vditor.value.setValue(content)
            console.log('✅ setValue调用成功')
            // 尝试再次获取值验证
            const getValue = vditor.value.getValue()
            console.log('🔍 getValue返回长度:', getValue.length)
            console.log('🔍 getValue返回内容:', getValue.substring(0, 200))
          } catch (error) {
            console.error('❌ setValue失败:', error)
          }
        } else {
          console.error('❌ vditor.value为null!')
        }
      },
      input: (value) => {
        // 同步内容到formData
        formData.value.content = value
      },
      blur: (value) => {
        formData.value.content = value
      }
    })
  }
})

// Cleanup on unmount
onBeforeUnmount(() => {
  if (vditor.value) {
    vditor.value.destroy()
    vditor.value = null
  }
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
  // 确保获取最新的编辑器内容
  if (vditor.value) {
    formData.value.content = vditor.value.getValue()
  }
  emit('save', formData.value)
}

// AI助手功能
const handleAIAssist = async (action) => {
  if (!vditor.value) return

  let loading
  try {
    const selectedText = vditor.value.getSelection()
    const fullText = vditor.value.getValue()
    
    // 根据不同操作准备参数
    let context = ''
    let text = ''
    let language = ''

    switch (action) {
      case 'continue':
        // 续写：需要上下文
        context = fullText || ''
        if (!context) {
          ElMessage.warning('请先输入一些内容作为上下文')
          return
        }
        break
      case 'rewrite':
      case 'expand':
        // 改写/扩写：需要选中文本
        text = selectedText
        if (!text) {
          ElMessage.warning('请先选中要处理的文本')
          return
        }
        break
      case 'translate':
        // 翻译：需要选中文本，自动检测语言
        text = selectedText
        if (!text) {
          ElMessage.warning('请先选中要翻译的文本')
          return
        }
        // 简单判断：包含中文则翻译为英文，否则翻译为中文
        language = /[\u4e00-\u9fa5]/.test(text) ? 'en' : 'zh'
        break
    }

    loading = ElLoading.service({
      lock: true,
      text: 'AI 正在处理中...',
      background: 'rgba(0, 0, 0, 0.7)',
    })

    const response = await aiAssist(action, context, text, language)
    const result = response.result

    // 根据操作类型插入或替换文本
    if (action === 'continue') {
      // 续写：追加到末尾
      vditor.value.insertValue('\n\n' + result)
    } else {
      // 改写/扩写/翻译：替换选中文本
      vditor.value.insertValue(result)
    }

    ElMessage.success('AI 处理完成')
  } catch (error) {
    console.error('AI assist error:', error)
    ElMessage.error(error.response?.data?.message || 'AI 处理失败')
  } finally {
    if (loading) {
      loading.close()
    }
  }
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
            <div class="flex-1 flex flex-col overflow-hidden">
              <!-- Title Input -->
              <div class="px-8 pt-6 pb-2">
                <input
                  v-model="formData.title"
                  type="text"
                  class="w-full bg-transparent text-4xl font-black text-slate-800 placeholder-slate-400 focus:outline-none"
                  :placeholder="t('editor.titlePlaceholder')"
                />
              </div>
              
              <!-- Vditor Editor Container -->
              <div class="flex-1 px-8 pb-8 overflow-hidden">
                <div ref="editorContainer" class="h-full vditor-wrapper"></div>
              </div>
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

                  <!-- AI写作助手快捷按钮 -->
                  <div class="mb-4 pt-3 border-t border-green-100">
                    <p class="text-[10px] font-bold text-slate-400 uppercase mb-2">写作助手</p>
                    <div class="grid grid-cols-2 gap-2">
                      <button
                        @click="handleAIAssist('continue')"
                        class="px-2 py-1.5 bg-blue-50 hover:bg-blue-100 border border-blue-200 text-blue-700 text-xs font-bold rounded transition-colors"
                      >
                        ✍️ AI续写
                      </button>
                      <button
                        @click="handleAIAssist('rewrite')"
                        class="px-2 py-1.5 bg-purple-50 hover:bg-purple-100 border border-purple-200 text-purple-700 text-xs font-bold rounded transition-colors"
                      >
                        ✨ AI改写
                      </button>
                      <button
                        @click="handleAIAssist('expand')"
                        class="px-2 py-1.5 bg-orange-50 hover:bg-orange-100 border border-orange-200 text-orange-700 text-xs font-bold rounded transition-colors"
                      >
                        📝 AI扩写
                      </button>
                      <button
                        @click="handleAIAssist('translate')"
                        class="px-2 py-1.5 bg-teal-50 hover:bg-teal-100 border border-teal-200 text-teal-700 text-xs font-bold rounded transition-colors"
                      >
                        🌐 AI翻译
                      </button>
                    </div>
                    <p class="text-[9px] text-slate-400 mt-2 italic">
                      提示：改写/扩写/翻译需先选中文本
                    </p>
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

/* Vditor样式调整 */
.vditor-wrapper :deep(.vditor) {
  border: none;
  height: 100%;
}

.vditor-wrapper :deep(.vditor-toolbar) {
  background-color: #f8fafc;
  border-bottom: 2px solid #e2e8f0;
  padding: 8px;
  position: relative;
}

/* 确保 Vditor 的 tooltip 显示在最上层 */
.vditor-wrapper :deep(.vditor-hint) {
  z-index: 9999 !important;
}

.vditor-wrapper :deep(.vditor-panel) {
  z-index: 9999 !important;
}

.vditor-wrapper :deep(.vditor-tip) {
  z-index: 9999 !important;
}

.vditor-wrapper :deep(.vditor-content) {
  background-color: white;
  font-size: 16px;
  line-height: 1.8;
}

.vditor-wrapper :deep(.vditor-ir) {
  padding: 20px;
}

.vditor-wrapper :deep(.vditor-ir pre.vditor-reset) {
  font-family: 'Courier New', monospace;
}
</style>
