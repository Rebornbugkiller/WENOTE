import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { getNotes, getTrashNotes, createNote, updateNote, deleteNote, restoreNote, batchDelete, batchRestore, emptyTrash } from '../api/note'
import { getNotebooks, getDefaultNotebook, createNotebook, updateNotebook, deleteNotebook } from '../api/notebook'
import { getTags, createTag, updateTag, deleteTag } from '../api/tag'

export function useNotes() {
  const { t } = useI18n()

  // State
  const notes = ref([])
  const notebooks = ref([])
  const tags = ref([])
  const isLoading = ref(false)
  const currentView = ref('active') // 'active' | 'trash' | 'starred' | number (notebook_id)
  const filterTagId = ref(null) // 标签筛选
  const searchQuery = ref('')

  // 内部分页状态（不导出）
  const page = ref(1)
  const pageSize = ref(20)
  const total = ref(0) // 总记录数

  // Fetch notebooks and tags
  const fetchInitialData = async () => {
    try {
      const [nbRes, tagRes] = await Promise.all([getNotebooks(), getTags()])
      notebooks.value = nbRes.list || []
      tags.value = tagRes.list || []
    } catch (err) {
      console.error('Failed to fetch initial data:', err)
    }
  }

  // Fetch notes based on current view
  const fetchNotes = async () => {
    isLoading.value = true
    try {
      const params = { page: page.value, page_size: pageSize.value }

      if (searchQuery.value) {
        params.keyword = searchQuery.value
      }

      if (currentView.value === 'starred') {
        params.is_starred = true
      } else if (typeof currentView.value === 'number') {
        params.notebook_id = currentView.value
      }

      if (filterTagId.value) {
        params.tag_id = filterTagId.value
      }

      let res
      if (currentView.value === 'trash') {
        res = await getTrashNotes(params)
      } else {
        res = await getNotes(params)
      }

      notes.value = res.list || []
      total.value = res.total || 0
    } catch (err) {
      console.error('Failed to fetch notes:', err)
      ElMessage.error(t('messages.fetchNotesFailed'))
    } finally {
      isLoading.value = false
    }
  }

  // Create a new note (优化：不等待刷新列表)
  const handleCreateNote = async (notebookId) => {
    try {
      // 如果没有指定笔记本，使用默认笔记本
      let targetNotebookId = notebookId
      if (!targetNotebookId) {
        const defaultNotebook = await getDefaultNotebook()
        targetNotebookId = defaultNotebook.id
      }

      const data = await createNote({
        notebook_id: targetNotebookId,
        title: '',
        content: ''
      })
      // 后台刷新列表，不阻塞返回
      Promise.all([fetchNotes(), fetchInitialData()])
      return data
    } catch (err) {
      console.error('Failed to create note:', err)
      ElMessage.error(t('messages.createNoteFailed'))
      throw err
    }
  }

  // Update a note
  const handleUpdateNote = async (noteData) => {
    try {
      await updateNote(noteData.id, {
        title: noteData.title,
        content: noteData.content,
        notebook_id: noteData.notebook_id,
        is_starred: noteData.is_starred,
        is_pinned: noteData.is_pinned,
        tag_ids: noteData.tags?.map(t => t.id) || []
      })
      ElMessage.success(t('editor.saveSuccess'))
      await Promise.all([fetchNotes(), fetchInitialData()])
    } catch (err) {
      console.error('Failed to update note:', err)
      ElMessage.error(t('editor.saveFailed'))
      throw err
    }
  }

  // Delete a note (soft delete or permanent)
  const handleDeleteNote = async (id) => {
    const isTrash = currentView.value === 'trash'
    const confirmMsg = isTrash ? t('messages.deleteNoteConfirm') : t('messages.moveToTrashConfirm')

    try {
      await ElMessageBox.confirm(confirmMsg, t('common.hint'), {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: isTrash ? 'warning' : 'info'
      })
    } catch {
      return // 用户取消
    }

    try {
      await deleteNote(id)
      ElMessage.success(isTrash ? t('messages.permanentlyDeleted') : t('messages.movedToTrash'))
      await Promise.all([fetchNotes(), fetchInitialData()])
    } catch (err) {
      console.error('Failed to delete note:', err)
      ElMessage.error(t('messages.deleteFailed'))
    }
  }

  // Restore a note from trash
  const handleRestoreNote = async (id) => {
    try {
      await restoreNote(id)
      ElMessage.success(t('messages.restored'))
      await Promise.all([fetchNotes(), fetchInitialData()])
    } catch (err) {
      console.error('Failed to restore note:', err)
      ElMessage.error(t('messages.restoreFailed'))
    }
  }

  // Toggle note status (starred/pinned)
  const handleToggleStatus = async (id, field) => {
    const note = notes.value.find(n => n.id === id)
    if (!note) return

    try {
      await updateNote(id, { [field]: !note[field] })
      await fetchNotes()
    } catch (err) {
      console.error('Failed to toggle status:', err)
      ElMessage.error(t('messages.operationFailed'))
    }
  }

  // Create notebook
  const handleCreateNotebook = async (name) => {
    try {
      await createNotebook({ name })
      ElMessage.success(t('messages.notebookCreated'))
      await fetchInitialData()
    } catch (err) {
      console.error('Failed to create notebook:', err)
      ElMessage.error(t('messages.deleteFailed'))
      throw err
    }
  }

  // Delete notebook
  const handleDeleteNotebook = async (id) => {
    try {
      await deleteNotebook(id)
      ElMessage.success(t('messages.notebookDeleted'))
      if (currentView.value === id) {
        currentView.value = 'active'
      }
      await fetchInitialData()
      await fetchNotes()
    } catch (err) {
      console.error('Failed to delete notebook:', err)
      ElMessage.error(t('messages.deleteFailed'))
    }
  }

  // Update notebook
  const handleUpdateNotebook = async (id, name) => {
    try {
      await updateNotebook(id, { name })
      ElMessage.success(t('messages.notebookUpdated'))
      await fetchInitialData()
    } catch (err) {
      console.error('Failed to update notebook:', err)
      ElMessage.error(t('messages.updateFailed'))
      throw err
    }
  }

  // Create tag
  const handleCreateTag = async (name, color) => {
    try {
      await createTag({ name, color })
      ElMessage.success(t('messages.tagCreated'))
      await fetchInitialData()
    } catch (err) {
      console.error('Failed to create tag:', err)
      ElMessage.error(t('messages.deleteFailed'))
      throw err
    }
  }

  // Delete tag
  const handleDeleteTag = async (id) => {
    try {
      await deleteTag(id)
      ElMessage.success(t('messages.tagDeleted'))
      await fetchInitialData()
    } catch (err) {
      console.error('Failed to delete tag:', err)
      ElMessage.error(t('messages.deleteFailed'))
    }
  }

  // Update tag
  const handleUpdateTag = async (id, name, color) => {
    try {
      await updateTag(id, { name, color })
      ElMessage.success(t('messages.tagUpdated'))
      await fetchInitialData()
      await fetchNotes()
    } catch (err) {
      console.error('Failed to update tag:', err)
      ElMessage.error(t('messages.updateFailed'))
      throw err
    }
  }

  // Change view
  const setView = (view) => {
    currentView.value = view
    page.value = 1
    fetchNotes()
  }

  // Search
  const setSearch = (query) => {
    searchQuery.value = query
    page.value = 1
    fetchNotes()
  }

  // Filter by tag
  const setFilterTag = (tagId) => {
    filterTagId.value = filterTagId.value === tagId ? null : tagId // 点击同一个取消筛选
    page.value = 1
    fetchNotes()
  }

  // Set page for pagination
  const setPage = (newPage) => {
    page.value = newPage
    fetchNotes()
  }

  // Batch delete notes permanently
  const handleBatchDelete = async (ids) => {
    try {
      await ElMessageBox.confirm(t('messages.batchDeleteConfirm', { count: ids.length }), t('common.hint'), {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      })
    } catch {
      return false
    }
    try {
      await batchDelete(ids)
      ElMessage.success(t('messages.permanentlyDeleted'))
      await Promise.all([fetchNotes(), fetchInitialData()])
      return true
    } catch (err) {
      console.error('Failed to batch delete:', err)
      ElMessage.error(t('messages.deleteFailed'))
      return false
    }
  }

  // Batch restore notes
  const handleBatchRestore = async (ids) => {
    try {
      await batchRestore(ids)
      ElMessage.success(t('messages.restored'))
      await Promise.all([fetchNotes(), fetchInitialData()])
      return true
    } catch (err) {
      console.error('Failed to batch restore:', err)
      ElMessage.error(t('messages.restoreFailed'))
      return false
    }
  }

  // Empty trash - delete all notes in trash
  const handleEmptyTrash = async () => {
    try {
      await ElMessageBox.confirm(t('messages.emptyTrashConfirm'), t('messages.emptyTrashTitle'), {
        confirmButtonText: t('messages.confirmEmpty'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      })
    } catch {
      return false
    }
    try {
      await emptyTrash()
      ElMessage.success(t('messages.permanentlyDeleted'))
      await Promise.all([fetchNotes(), fetchInitialData()])
      return true
    } catch (err) {
      console.error('Failed to empty trash:', err)
      ElMessage.error(t('messages.emptyTrashFailed'))
      return false
    }
  }

  return {
    // State
    notes,
    notebooks,
    tags,
    isLoading,
    currentView,
    filterTagId,
    searchQuery,
    // Pagination
    page,
    pageSize,
    total,

    // Methods
    fetchInitialData,
    fetchNotes,
    handleCreateNote,
    handleUpdateNote,
    handleDeleteNote,
    handleRestoreNote,
    handleToggleStatus,
    handleCreateNotebook,
    handleUpdateNotebook,
    handleDeleteNotebook,
    handleCreateTag,
    handleUpdateTag,
    handleDeleteTag,
    handleBatchDelete,
    handleBatchRestore,
    handleEmptyTrash,
    setView,
    setSearch,
    setFilterTag,
    setPage
  }
}
