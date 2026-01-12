-- WENOTE 数据库初始化脚本
-- 创建数据库
CREATE DATABASE IF NOT EXISTS wenote DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE wenote;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE COMMENT '用户名',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 笔记本表
CREATE TABLE IF NOT EXISTS notebooks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    name VARCHAR(255) NOT NULL COMMENT '笔记本名称',
    is_default TINYINT(1) DEFAULT 0 COMMENT '是否为默认笔记本',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='笔记本表';

-- 笔记表
CREATE TABLE IF NOT EXISTS notes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    notebook_id BIGINT UNSIGNED NOT NULL COMMENT '笔记本ID',
    title VARCHAR(255) COMMENT '标题',
    content LONGTEXT COMMENT '内容',
    -- AI 相关字段
    summary TEXT COMMENT 'AI生成的摘要',
    summary_len INT DEFAULT 200 COMMENT '摘要长度控制',
    suggested_tags JSON COMMENT 'AI生成的标签建议',
    ai_status ENUM('pending','running','done','failed') DEFAULT 'pending' COMMENT 'AI任务状态',
    ai_error TEXT COMMENT 'AI任务错误信息',
    ai_enabled TINYINT(1) DEFAULT 0 COMMENT '是否启用AI',
    -- 状态字段
    is_pinned TINYINT(1) DEFAULT 0 COMMENT '是否置顶',
    is_starred TINYINT(1) DEFAULT 0 COMMENT '是否星标',
    -- 软删除
    deleted_at DATETIME DEFAULT NULL COMMENT '删除时间（NULL表示未删除）',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_updated (user_id, updated_at DESC),
    INDEX idx_notebook_id (notebook_id),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_ai_status (ai_status),
    FULLTEXT INDEX ft_title_content (title, content) WITH PARSER ngram,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (notebook_id) REFERENCES notebooks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='笔记表';

-- 标签表
CREATE TABLE IF NOT EXISTS tags (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    name VARCHAR(100) NOT NULL COMMENT '标签名称',
    color VARCHAR(20) DEFAULT '#6B7280' COMMENT '标签颜色',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE INDEX idx_user_name (user_id, name),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='标签表';

-- 笔记标签关联表
CREATE TABLE IF NOT EXISTS note_tags (
    note_id BIGINT UNSIGNED NOT NULL COMMENT '笔记ID',
    tag_id BIGINT UNSIGNED NOT NULL COMMENT '标签ID',
    PRIMARY KEY (note_id, tag_id),
    INDEX idx_tag_id (tag_id),
    FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='笔记标签关联表';

-- 审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '操作用户ID',
    action VARCHAR(50) NOT NULL COMMENT '操作类型',
    resource_type VARCHAR(50) NOT NULL COMMENT '资源类型',
    resource_id BIGINT UNSIGNED COMMENT '资源ID',
    details JSON COMMENT '详细信息',
    ip_address VARCHAR(50) COMMENT '操作IP',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
    INDEX idx_user_action (user_id, action),
    INDEX idx_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审计日志表';
