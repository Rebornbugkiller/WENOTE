# 📒 WeNote - 你的智能笔记小助手

[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3.5-4FC08D?logo=vue.js)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

> 一个会思考的笔记本 - 用 AI 帮你整理思绪，让灵感不再乱飞 ✨

## 🌐 在线体验

**👉 [http://47.97.78.32](http://47.97.78.32)**

---

<img width="1910" height="827" alt="image" src="https://github.com/user-attachments/assets/5cfa4c14-b83c-47ff-9760-5d4b4c9ec620" />
<img width="1918" height="826" alt="image" src="https://github.com/user-attachments/assets/03f73025-05a6-49b8-9658-84cf30208ae6" />

<img width="1467" height="710" alt="image" src="https://github.com/user-attachments/assets/be2dcd65-d8f8-4ceb-ad22-a8523397f564" />

## 为什么选择 WeNote？

**厌倦了笔记越记越乱？** WeNote 来拯救你！

这不只是一个笔记应用，更是一个让你**爱上记录**的游戏化写作平台。

| 特性 | 描述 |
|------|------|
| 📝 **随心记录** | Markdown + 富文本编辑器，想怎么写就怎么写 |
| 🤖 **AI 加持** | 一键生成摘要和标签，告别分类焦虑 |
| 🎮 **游戏化体验** | 连续打卡、成就系统、写作报告，让记笔记像打游戏一样上瘾 |
| 📊 **数据可视化** | ECharts 图表展示你的写作趋势和习惯 |
| 🔒 **安全可靠** | 密码加密 + JWT 认证，你的笔记只属于你 |
| 🌐 **双语支持** | 中英文无缝切换，国际化体验 |
| 🐳 **一键部署** | Docker Compose 搞定一切，5 分钟上线 |

---

## 亮点功能

### 🎮 游戏化写作系统

让记笔记变得有趣！

- **连续打卡** - 每天写笔记，点燃你的写作火焰 🔥
- **成就徽章** - 解锁各种成就，收集专属徽章 🏆
- **每日目标** - 设定字数目标，追踪完成进度
- **写作报告** - 周报/月报统计，见证你的成长轨迹
- **Combo 连击** - 登录页面的趣味互动，打出超级连击！

### 🎨 复古像素风登录页

不只是登录，更是一场视觉盛宴：

- **GameBoy 贪吃蛇** - 等待时来一局经典游戏
- **像素风头像** - 可爱的互动小蛇陪你登录（50+ 条趣味台词！）
- **CRT 扫描线** - 复古显示器效果，满满的怀旧感
- **Fever 模式** - 开启狂热模式，体验不一样的视觉冲击
- **背景音乐** - 8-bit 风格 BGM，支持一键开关，记住你的偏好
- **记住密码** - 登录成功后自动记住账号密码，下次无需重复输入

### 📊 数据统计仪表盘

用数据见证你的写作之旅：

- **趋势图表** - 最近 7 天的写作趋势一目了然
- **标签云** - TOP10 热门标签可视化
- **笔记本分布** - 饼图展示各笔记本占比
- **字数统计** - 总字数、本周新增，量化你的努力

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

## 功能清单

### 基础功能
- ✅ 用户注册登录（安全第一）
- ✅ 创建、编辑、删除笔记
- ✅ 多笔记本管理
- ✅ 自定义标签和颜色
- ✅ 全文搜索（支持中文）
- ✅ 智能回收站（误删不怕）
- ✅ 图片上传支持
- ✅ 个人设置中心

### AI 智能功能
- 🤖 **自动摘要** - 长篇笔记一键生成摘要
- 🏷️ **智能标签** - AI 帮你推荐合适的标签
- 💾 **自动保存** - AI 生成前自动保存，不丢失任何内容

### 游戏化功能
- 🔥 **连续打卡** - 追踪你的写作连续天数
- 🏆 **成就系统** - 多种成就等你解锁（笔记、连续、字数、目标）
- 📈 **每日目标** - 自定义字数目标，培养写作习惯
- 📋 **写作报告** - 周报/月报，对比上期数据

### 高级特性
- 🛡️ **三级限流** - 防止接口被刷，保护服务稳定
- 🔌 **熔断降级** - AI 服务挂了也不影响正常使用
- 📊 **操作审计** - 所有操作都有记录
- 🌍 **国际化** - 中英文双语支持

---

## 技术栈

| 层级 | 技术 |
|------|------|
| **前端** | Vue 3 + Vite + Element Plus + Pinia + TailwindCSS + ECharts |
| **后端** | Go 1.23 + Gin + GORM |
| **数据库** | MySQL 8.0 |
| **AI** | 智谱 GLM-4 |
| **国际化** | Vue I18n（中/英双语） |
| **部署** | Docker + Nginx + Systemd |

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
│   ├── cmd/server/          # 程序入口
│   ├── config/              # 配置文件
│   ├── internal/            # 核心业务逻辑
│   │   ├── handler/         # HTTP 处理器
│   │   ├── service/         # 业务服务层
│   │   ├── repo/            # 数据访问层
│   │   ├── model/           # 数据模型
│   │   ├── router/          # 路由定义
│   │   └── middleware/      # 中间件（认证、限流、日志）
│   ├── pkg/                 # 公共工具包
│   │   ├── ai/              # AI 客户端（智谱 GLM）
│   │   ├── jwt/             # JWT 认证
│   │   ├── logger/          # 日志工具
│   │   ├── hash/            # 密码加密
│   │   └── response/        # 统一响应格式
│   └── scripts/             # 数据库初始化脚本
├── wenote-frontend/         # Vue 前端
│   ├── src/
│   │   ├── views/           # 页面组件（Home, Login, Editor, Settings）
│   │   ├── components/      # 通用组件
│   │   │   ├── login/       # 登录页组件（贪吃蛇、像素头像、HUD等）
│   │   │   ├── notes/       # 笔记相关组件
│   │   │   ├── gamification/# 游戏化组件（成就、连续打卡、写作报告）
│   │   │   ├── stats/       # 数据统计仪表盘
│   │   │   └── settings/    # 设置页组件
│   │   ├── api/             # API 调用封装
│   │   ├── stores/          # Pinia 状态管理
│   │   ├── router/          # 路由配置
│   │   ├── composables/     # 组合式函数
│   │   ├── i18n/            # 国际化配置
│   │   ├── styles/          # 样式文件
│   │   └── utils/           # 工具函数
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

**已完成：**
- [x] 游戏化写作系统（连续打卡、成就、每日目标）
- [x] 数据统计仪表盘
- [x] 复古像素风登录页
- [x] 国际化支持（中/英）
- [x] 个人设置中心
- [x] 图片上传功能
- [x] 背景音乐控制（记住用户偏好）
- [x] 记住登录账号密码
- [x] 小蛇 50+ 条趣味互动台词

**进行中：**
- [ ] AI 写作助手（续写、改写、扩写、翻译）
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



感谢以下开源项目：
- [Gin](https://github.com/gin-gonic/gin) - 高性能 HTTP 框架
- [GORM](https://gorm.io/) - 优雅的 ORM 库
- [Vue.js](https://vuejs.org/) - 渐进式前端框架
- [Element Plus](https://element-plus.org/) - Vue 3 组件库
- [ECharts](https://echarts.apache.org/) - 强大的可视化图表库
- [智谱 AI](https://open.bigmodel.cn/) - AI 能力支持

---

## License

MIT License - 随便用，记得 Star 就行 ⭐

---

<div align="center">

### 让记笔记变成一件有趣的事

**如果这个项目对你有帮助，请给个 Star 支持一下！**

[![Star History Chart](https://img.shields.io/github/stars/Rebornbugkiller/WENOTE?style=social)](https://github.com/Rebornbugkiller/WENOTE)

**有问题？** 提 [Issue](https://github.com/Rebornbugkiller/WENOTE/issues) | **想贡献？** 发 [Pull Request](https://github.com/Rebornbugkiller/WENOTE/pulls)

</div>
