import api from './index'

export const getNotes = (params) => api.get('/notes', { params })
export const getTrashNotes = (params) => api.get('/notes/trash', { params })
export const getNote = (id) => api.get(`/notes/${id}`)
export const createNote = (data) => api.post('/notes', data)
export const updateNote = (id, data) => api.patch(`/notes/${id}`, data)
export const deleteNote = (id) => api.delete(`/notes/${id}`)
export const restoreNote = (id) => api.post(`/notes/${id}/restore`)
export const updateNoteTags = (id, tagIds) => api.put(`/notes/${id}/tags`, { tag_ids: tagIds })
export const applySuggestedTags = (id) => api.put(`/notes/${id}/tags/apply-suggestions`)
export const generateSummaryAndTags = (id) => api.post(`/notes/${id}/ai/generate`)
export const batchDelete = (noteIds) => api.post('/notes/batch/delete', { note_ids: noteIds })
export const batchRestore = (noteIds) => api.post('/notes/batch/restore', { note_ids: noteIds })
export const batchMove = (noteIds, notebookId) => api.post('/notes/batch/move', { note_ids: noteIds, notebook_id: notebookId })

// 附件相关
export const uploadImage = (noteId, file) => {
  const formData = new FormData()
  formData.append('file', file)
  return api.post(`/notes/${noteId}/attachments`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
export const getAttachments = (noteId) => api.get(`/notes/${noteId}/attachments`)
export const deleteAttachment = (attachmentId) => api.delete(`/attachments/${attachmentId}`)

// AI写作助手
export const aiAssist = (action, context, selectedText, language) => {
  return api.post('/notes/ai/assist', {
    action,
    context,
    selected_text: selectedText,
    language
  })
}

// 导入导出
export const exportAllNotes = () => {
  return api.get('/notes/export', { responseType: 'blob' })
}

export const exportNote = (noteId) => {
  return api.get(`/notes/${noteId}/export`, { responseType: 'blob' })
}

export const importNotes = (file) => {
  const formData = new FormData()
  formData.append('file', file)
  return api.post('/notes/import', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
