# 📒 WeNote - 你的智能笔记小助手

[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3.5-4FC08D?logo=vue.js)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

> 一个会思考的笔记本 - 用 AI 帮你整理思绪，让灵感不再乱飞 ✨

---

## 为什么选择 WeNote？

**厌倦了笔记越记越乱？** WeNote 来拯救你！

- 📝 **随心记录** - 想写就写，支持 Markdown，还有富文本编辑器
- 🤖 **AI 小助手** - 自动生成摘要和标签，再也不用纠结怎么分类
- 🗂️ **智能整理** - 多笔记本管理，标签系统，全文搜索，找笔记就像搜索引擎一样快
- 🔒 **安全可靠** - 你的笔记只属于你，密码加密，JWT 认证
- 🎨 **界面清爽** - 简洁美观，专注内容本身
- 🐳 **一键部署** - Docker Compose 搞定一切，5 分钟上线

---

## 快速开始

### 方式一：Docker 一键启动（推荐）

```bash
# 克隆项目
git clone https://github.com/Rebornbugkiller/WENOTE.git
cd wenote

# 配置环境变量（复制 .env.example 并修改）
cp .env.example .env

# 一键启动（Windows）
start.bat

# 一键启动（Linux/Mac）
./start.sh
```

打开浏览器访问 `http://localhost` 就能用了！

### 方式二：本地开发

**后端：**
```bash
cd wenote-backend
cp config/config.example.yaml config/config.yaml
# 编辑 config.yaml 填入数据库密码和 JWT 密钥
go run cmd/server/main.go
```

**前端：**
```bash
cd wenote-frontend
npm install
npm run dev
```

---

## 核心功能

### 基础功能
- ✅ 用户注册登录（安全第一）
- ✅ 创建、编辑、删除笔记
- ✅ 多笔记本管理
- ✅ 自定义标签和颜色
- ✅ 全文搜索（支持中文）
- ✅ 智能回收站（误删不怕）

### AI 智能功能
- 🤖 **自动摘要** - 长篇笔记一键生成摘要
- 🏷️ **智能标签** - AI 帮你推荐合适的标签
- ⚡ **快速整理** - 让 AI 帮你理清思路

### 高级特性
- 🛡️ **三级限流** - 防止接口被刷，保护服务稳定
- 🔌 **熔断降级** - AI 服务挂了也不影响正常使用
- 📊 **操作审计** - 所有操作都有记录

---

## 技术栈

**前端：** Vue 3 + Element Plus + Vite
**后端：** Go + Gin + GORM
**数据库：** MySQL 8.0
**AI：** 智谱 GLM-4
**部署：** Docker + Docker Compose

---

## 配置说明

### 必填配置

编辑 `wenote-backend/config/config.yaml`：

```yaml
database:
  password: YOUR_DB_PASSWORD      # 数据库密码

jwt:
  secret: YOUR_RANDOM_SECRET_KEY  # JWT 密钥（至少 32 位随机字符串）
```

### 可选配置（AI 功能）

如果想用 AI 功能，需要配置智谱 API Key：

```yaml
ai:
  zhipu:
    api_key: YOUR_ZHIPU_API_KEY   # 在 https://open.bigmodel.cn/ 获取
```

**获取 API Key：**
1. 访问 [智谱 AI 开放平台](https://open.bigmodel.cn/)
2. 注册并实名认证
3. 创建 API Key（新用户有免费额度）

---

## API 文档

### 认证接口
- `POST /api/v1/auth/register` - 注册
- `POST /api/v1/auth/login` - 登录

### 笔记接口
- `GET /api/v1/notes` - 获取笔记列表
- `POST /api/v1/notes` - 创建笔记
- `GET /api/v1/notes/:id` - 获取笔记详情
- `PATCH /api/v1/notes/:id` - 更新笔记
- `DELETE /api/v1/notes/:id` - 删除笔记
- `POST /api/v1/notes/:id/restore` - 恢复笔记
- `POST /api/v1/notes/:id/ai/generate` - AI 生成摘要和标签

### 笔记本接口
- `GET /api/v1/notebooks` - 获取笔记本列表
- `POST /api/v1/notebooks` - 创建笔记本

### 标签接口
- `GET /api/v1/tags` - 获取标签列表
- `POST /api/v1/tags` - 创建标签

---

## 项目结构

```
wenote/
├── wenote-backend/          # Go 后端
│   ├── cmd/                 # 程序入口
│   ├── config/              # 配置文件
│   ├── internal/            # 核心业务逻辑
│   │   ├── handler/         # HTTP 处理器
│   │   ├── service/         # 业务服务
│   │   ├── repo/            # 数据访问
│   │   └── model/           # 数据模型
│   └── pkg/                 # 公共工具包
│       ├── ai/              # AI 客户端
│       ├── jwt/             # JWT 认证
│       └── logger/          # 日志
├── wenote-frontend/         # Vue 前端
│   ├── src/
│   │   ├── views/           # 页面
│   │   ├── components/      # 组件
│   │   ├── api/             # API 调用
│   │   └── composables/     # 组合式函数
│   └── public/              # 静态资源
├── docker-compose.yaml      # Docker 编排
├── start.bat                # Windows 启动脚本
└── start.sh                 # Linux/Mac 启动脚本
```

---

## 常见问题

**Q: AI 功能不可用怎么办？**
A: 检查 `config.yaml` 中的 `ai.zhipu.api_key` 是否配置正确。如果不需要 AI 功能，可以不配置，不影响其他功能使用。

**Q: 如何修改端口？**
A: 编辑 `docker-compose.yaml` 中的端口映射，或者修改 `.env` 文件中的端口配置。

**Q: 数据存储在哪里？**
A: 使用 Docker 部署时，数据存储在 Docker Volume 中。如果需要备份，可以导出 MySQL 数据。

**Q: 支持多用户吗？**
A: 支持！每个用户的数据完全隔离，互不影响。

---

## 开发计划

- [ ] 支持笔记分享
- [ ] 支持协作编辑
- [ ] 移动端适配
- [ ] 更多 AI 功能（智能问答、内容扩写等）
- [ ] 支持更多导出格式（PDF、Word 等）

---

## 贡献指南

欢迎提交 Issue 和 Pull Request！

如果你有好的想法或建议，欢迎在 [Issues](https://github.com/Rebornbugkiller/WENOTE/issues) 中讨论。

---

## 致谢

感谢导师 **王申跃** 的悉心指导。

感谢以下开源项目：
- [Gin](https://github.com/gin-gonic/gin) - 高性能 HTTP 框架
- [GORM](https://gorm.io/) - 优雅的 ORM 库
- [Vue.js](https://vuejs.org/) - 渐进式前端框架
- [Element Plus](https://element-plus.org/) - Vue 3 组件库
- [智谱 AI](https://open.bigmodel.cn/) - AI 能力支持

---

## License

MIT License - 随便用，记得 Star 就行 ⭐

---

**如果这个项目对你有帮助，欢迎 Star！**

**有问题？** 提 [Issue](https://github.com/Rebornbugkiller/WENOTE/issues)
**想贡献？** 发 [Pull Request](https://github.com/Rebornbugkiller/WENOTE/pulls)
