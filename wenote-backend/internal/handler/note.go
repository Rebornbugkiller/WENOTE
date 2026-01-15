package handler

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

import (
	"wenote-backend/internal/model"
	"wenote-backend/internal/repo"
	"wenote-backend/internal/service"
	"wenote-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// NoteHandler 笔记处理器
type NoteHandler struct {
	noteService *service.NoteService
	auditRepo   *repo.AuditRepo
}

// NewNoteHandler 创建笔记处理器实例
func NewNoteHandler() *NoteHandler {
	return &NoteHandler{
		noteService: service.NewNoteService(),
		auditRepo:   repo.NewAuditRepo(),
	}
}

// Create 创建笔记
func (h *NoteHandler) Create(c *gin.Context) {
	userID := c.GetUint64("userID")

	var req model.NoteCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	note, err := h.noteService.Create(userID, &req)
	if err != nil {
		if err == service.ErrNotebookNotFound {
			response.BadRequest(c, "笔记本不存在")
			return
		}
		response.InternalError(c, "创建笔记失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", note)
}

// GetByID 获取笔记详情
func (h *NoteHandler) GetByID(c *gin.Context) {
	userID := c.GetUint64("userID")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	note, err := h.noteService.GetByID(userID, noteID)
	if err != nil {
		if err == service.ErrNoteNotFound {
			response.NotFound(c, "笔记不存在")
			return
		}
		response.InternalError(c, "获取笔记失败")
		return
	}

	response.Success(c, note)
}

// Update 更新笔记
func (h *NoteHandler) Update(c *gin.Context) {
	userID := c.GetUint64("userID")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	var req model.NoteUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	note, err := h.noteService.Update(userID, noteID, &req)
	if err != nil {
		if err == service.ErrNoteNotFound {
			response.NotFound(c, "笔记不存在")
			return
		}
		if err == service.ErrNotebookNotFound {
			response.BadRequest(c, "笔记本不存在")
			return
		}
		response.InternalError(c, "更新笔记失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", note)
}

// Delete 删除笔记（软删除）
func (h *NoteHandler) Delete(c *gin.Context) {
	userID := c.GetUint64("userID")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	if err := h.noteService.Delete(userID, noteID); err != nil {
		if err == service.ErrNoteNotFound {
			response.NotFound(c, "笔记不存在")
			return
		}
		response.InternalError(c, "删除笔记失败")
		return
	}

	// 记录审计日志
	_ = h.auditRepo.Create(&model.AuditLog{
		UserID:       userID,
		Action:       "delete",
		ResourceType: "note",
		ResourceID:   noteID,
		IPAddress:    c.ClientIP(),
	})

	response.SuccessWithMessage(c, "删除成功", nil)
}

// Restore 恢复已删除的笔记
func (h *NoteHandler) Restore(c *gin.Context) {
	userID := c.GetUint64("userID")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	note, err := h.noteService.Restore(userID, noteID)
	if err != nil {
		if err == service.ErrNoteNotFound {
			response.NotFound(c, "笔记不存在")
			return
		}
		response.InternalError(c, "恢复笔记失败")
		return
	}

	// 记录审计日志
	_ = h.auditRepo.Create(&model.AuditLog{
		UserID:       userID,
		Action:       "restore",
		ResourceType: "note",
		ResourceID:   noteID,
		IPAddress:    c.ClientIP(),
	})

	response.SuccessWithMessage(c, "恢复成功", note)
}

// List 获取笔记列表
func (h *NoteHandler) List(c *gin.Context) {
	userID := c.GetUint64("userID")

	var req model.NoteListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resp, err := h.noteService.List(userID, &req)
	if err != nil {
		response.InternalError(c, "获取笔记列表失败")
		return
	}

	response.Success(c, resp)
}

// UpdateTags 更新笔记标签
func (h *NoteHandler) UpdateTags(c *gin.Context) {
	userID := c.GetUint64("userID")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	var req model.NoteTagsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	note, err := h.noteService.UpdateTags(userID, noteID, req.TagIDs)
	if err != nil {
		if err == service.ErrNoteNotFound {
			response.NotFound(c, "笔记不存在")
			return
		}
		response.InternalError(c, "更新标签失败")
		return
	}

	response.SuccessWithMessage(c, "更新成功", note)
}

// ApplySuggestedTags 应用AI建议的标签
func (h *NoteHandler) ApplySuggestedTags(c *gin.Context) {
	userID := c.GetUint64("userID")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	if err := h.noteService.ApplySuggestedTags(userID, noteID); err != nil {
		if err == service.ErrNoteNotFound {
			response.NotFound(c, "笔记不存在")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "标签应用成功", nil)
}

// ListDeleted 获取回收站笔记列表
func (h *NoteHandler) ListDeleted(c *gin.Context) {
	userID := c.GetUint64("userID")

	// 获取分页参数
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if val, err := strconv.Atoi(ps); err == nil && val > 0 {
			pageSize = val
		}
	}

	resp, err := h.noteService.ListDeleted(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, "获取回收站列表失败")
		return
	}

	response.Success(c, resp)
}

// BatchDelete 批量删除笔记
func (h *NoteHandler) BatchDelete(c *gin.Context) {
	userID := c.GetUint64("userID")

	var req model.BatchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	count, err := h.noteService.BatchHardDelete(req.NoteIDs, userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 记录审计日志
	_ = h.auditRepo.Create(&model.AuditLog{
		UserID:       userID,
		Action:       "batch_delete",
		ResourceType: "note",
		Details: map[string]interface{}{
			"note_ids":      req.NoteIDs,
			"deleted_count": count,
		},
		IPAddress: c.ClientIP(),
	})

	response.SuccessWithMessage(c, "批量删除成功", map[string]interface{}{
		"deleted_count": count,
	})
}

// BatchRestore 批量恢复笔记
func (h *NoteHandler) BatchRestore(c *gin.Context) {
	userID := c.GetUint64("userID")

	var req model.BatchRestoreReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	count, err := h.noteService.BatchRestore(req.NoteIDs, userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 记录审计日志
	_ = h.auditRepo.Create(&model.AuditLog{
		UserID:       userID,
		Action:       "batch_restore",
		ResourceType: "note",
		Details: map[string]interface{}{
			"note_ids":       req.NoteIDs,
			"restored_count": count,
		},
		IPAddress: c.ClientIP(),
	})

	response.SuccessWithMessage(c, "批量恢复成功", map[string]interface{}{
		"restored_count": count,
	})
}

// BatchMove 批量移动笔记
func (h *NoteHandler) BatchMove(c *gin.Context) {
	userID := c.GetUint64("userID")

	var req model.BatchMoveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	count, err := h.noteService.BatchMove(req.NoteIDs, req.NotebookID, userID)
	if err != nil {
		if err == service.ErrNotebookNotFound {
			response.BadRequest(c, "笔记本不存在")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	// 记录审计日志
	_ = h.auditRepo.Create(&model.AuditLog{
		UserID:       userID,
		Action:       "batch_move",
		ResourceType: "note",
		Details: map[string]interface{}{
			"note_ids":    req.NoteIDs,
			"notebook_id": req.NotebookID,
			"moved_count": count,
		},
		IPAddress: c.ClientIP(),
	})

	response.SuccessWithMessage(c, "批量移动成功", map[string]interface{}{
		"moved_count": count,
	})
}

// GenerateSummaryAndTags 同步生成摘要和标签
func (h *NoteHandler) GenerateSummaryAndTags(c *gin.Context) {
	userID := c.GetUint64("userID")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	result, err := h.noteService.GenerateSummaryAndTagsSync(userID, noteID)
	if err != nil {
		if err == service.ErrNoteNotFound {
			response.NotFound(c, "笔记不存在")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"summary": result.Summary,
		"tags":    result.Tags,
	})
}

// AIAssist AI写作助手
// POST /api/v1/notes/ai/assist
func (h *NoteHandler) AIAssist(c *gin.Context) {
	var req model.AIAssistReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 创建超时上下文
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	result, err := h.noteService.AIAssist(ctx, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// ExportNote 导出单个笔记
// GET /api/v1/notes/:id/export
func (h *NoteHandler) ExportNote(c *gin.Context) {
	userID := c.GetUint64("userID")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	note, err := h.noteService.GetByID(userID, noteID)
	if err != nil {
		response.NotFound(c, "笔记不存在")
		return
	}

	// 生成Markdown内容（包含YAML元数据）
	markdown := h.generateMarkdown(note)
	filename := fmt.Sprintf("%s_%s.md", sanitizeFilename(note.Title), time.Now().Format("20060102"))

	c.Header("Content-Type", "text/markdown; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.String(200, markdown)
}

// ExportAllNotes 导出所有笔记为ZIP
// GET /api/v1/notes/export
func (h *NoteHandler) ExportAllNotes(c *gin.Context) {
	userID := c.GetUint64("userID")

	// 获取所有笔记
	req := &model.NoteListReq{Page: 1, PageSize: 10000}
	resp, err := h.noteService.List(userID, req)
	if err != nil {
		response.InternalError(c, "获取笔记失败")
		return
	}

	// 创建ZIP文件
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, note := range resp.List {
		markdown := h.generateMarkdown(note)
		filename := fmt.Sprintf("%s_%d.md", sanitizeFilename(note.Title), note.ID)
		
		writer, err := zipWriter.Create(filename)
		if err != nil {
			continue
		}
		writer.Write([]byte(markdown))
	}

	zipWriter.Close()

	filename := fmt.Sprintf("wenote_backup_%s.zip", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(200, "application/zip", buf.Bytes())
}

// ImportNotes 导入笔记
// POST /api/v1/notes/import
func (h *NoteHandler) ImportNotes(c *gin.Context) {
	userID := c.GetUint64("userID")

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "未找到上传文件")
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		response.InternalError(c, "打开文件失败")
		return
	}
	defer src.Close()

	var imported int
	var failed int

	// 根据文件类型处理
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == ".zip" {
		// ZIP文件 - 批量导入
		imported, failed, err = h.importFromZip(userID, src, file.Size)
	} else if ext == ".md" {
		// 单个Markdown文件
		content, _ := io.ReadAll(src)
		if h.importMarkdownFile(userID, string(content)) {
			imported = 1
		} else {
			failed = 1
		}
	} else {
		response.BadRequest(c, "不支持的文件类型（仅支持.md和.zip）")
		return
	}

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"imported": imported,
		"failed":   failed,
	})
}

// generateMarkdown 生成带YAML元数据的Markdown
func (h *NoteHandler) generateMarkdown(note *model.Note) string {
	// 准备元数据
	metadata := map[string]interface{}{
		"title":      note.Title,
		"created_at": note.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 添加标签
	if len(note.Tags) > 0 {
		tagNames := make([]string, len(note.Tags))
		for i, tag := range note.Tags {
			tagNames[i] = tag.Name
		}
		metadata["tags"] = tagNames
	}

	// 序列化为YAML
	yamlData, _ := yaml.Marshal(metadata)

	// 组合YAML前置和Markdown内容
	return fmt.Sprintf("---\n%s---\n\n%s", string(yamlData), note.Content)
}

// importFromZip 从ZIP导入
func (h *NoteHandler) importFromZip(userID uint64, reader io.Reader, size int64) (int, int, error) {
	// 读取ZIP内容到内存
	buf := new(bytes.Buffer)
	io.Copy(buf, reader)

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), size)
	if err != nil {
		return 0, 0, fmt.Errorf("解析ZIP文件失败")
	}

	imported := 0
	failed := 0

	for _, file := range zipReader.File {
		if !strings.HasSuffix(strings.ToLower(file.Name), ".md") {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			failed++
			continue
		}

		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			failed++
			continue
		}

		if h.importMarkdownFile(userID, string(content)) {
			imported++
		} else {
			failed++
		}
	}

	return imported, failed, nil
}

// importMarkdownFile 导入单个Markdown文件
func (h *NoteHandler) importMarkdownFile(userID uint64, content string) bool {
	// 解析YAML前置和Markdown内容
	title, tags, mdContent := parseMarkdown(content)

	// 获取默认笔记本
	notebooks, _ := repo.NewNotebookRepo().ListByUserID(userID)
	var notebookID uint64
	for _, nb := range notebooks {
		if nb.IsDefault {
			notebookID = nb.ID
			break
		}
	}
	if notebookID == 0 && len(notebooks) > 0 {
		notebookID = notebooks[0].ID
	}

	// 创建笔记
	req := &model.NoteCreateReq{
		NotebookID: notebookID,
		Title:      title,
		Content:    mdContent,
	}

	note, err := h.noteService.Create(userID, req)
	if err != nil {
		return false
	}

	// 如果有标签，创建并关联
	if len(tags) > 0 {
		tagRepo := repo.NewTagRepo()
		var tagIDs []uint64
		for _, tagName := range tags {
			// 尝试查找已存在的标签
			existingTags, _ := tagRepo.ListByUserID(userID)
			var tagID uint64
			for _, t := range existingTags {
				if t.Name == tagName {
					tagID = t.ID
					break
				}
			}
			
			// 如果不存在则创建
			if tagID == 0 {
				newTag := &model.Tag{
					UserID: userID,
					Name:   tagName,
					Color:  "#6B7280",
				}
				if err := tagRepo.Create(newTag); err == nil {
					tagID = newTag.ID
				}
			}
			
			if tagID > 0 {
				tagIDs = append(tagIDs, tagID)
			}
		}
		
		if len(tagIDs) > 0 {
			repo.NewNoteRepo().ReplaceNoteTags(note.ID, tagIDs)
		}
	}

	return true
}

// parseMarkdown 解析Markdown（提取YAML前置）
func parseMarkdown(content string) (string, []string, string) {
	title := "未命名笔记"
	var tags []string
	mdContent := content

	// 检查是否有YAML前置
	if strings.HasPrefix(content, "---\n") {
		parts := strings.SplitN(content, "---\n", 3)
		if len(parts) >= 3 {
			yamlContent := parts[1]
			mdContent = strings.TrimSpace(parts[2])

			// 解析YAML
			var metadata map[string]interface{}
			if err := yaml.Unmarshal([]byte(yamlContent), &metadata); err == nil {
				if t, ok := metadata["title"].(string); ok && t != "" {
					title = t
				}
				if tagList, ok := metadata["tags"].([]interface{}); ok {
					for _, tag := range tagList {
						if tagStr, ok := tag.(string); ok {
							tags = append(tags, tagStr)
						}
					}
				}
			}
		}
	}

	// 如果没有标题，从内容第一行提取
	if title == "未命名笔记" && mdContent != "" {
		lines := strings.Split(mdContent, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				// 移除Markdown标题标记
				title = strings.TrimPrefix(line, "#")
				title = strings.TrimSpace(title)
				if len(title) > 50 {
					title = title[:50] + "..."
				}
				break
			}
		}
	}

	return title, tags, mdContent
}

// sanitizeFilename 清理文件名
func sanitizeFilename(name string) string {
	if name == "" {
		return "untitled"
	}
	// 移除非法字符
	name = strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' {
			return '_'
		}
		return r
	}, name)
	
	// 限制长度
	if len(name) > 50 {
		name = name[:50]
	}
	return name
}
