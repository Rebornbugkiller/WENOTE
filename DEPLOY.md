# 🚀 WeNote 服务器部署指南

## 📋 部署前准备

### 1. 服务器信息
- **公网 IP**: `47.97.78.32`（你的服务器 IP）
- **操作系统**: Ubuntu
- **配置**: 2 vCPU / 2 GiB / 40 GiB

### 2. 本地准备
- SSH 客户端（Windows 可用 PowerShell、PuTTY 或 Git Bash）
- 项目代码已准备好

---

## 🔧 第一步：设置服务器密码并连接

### 方式一：使用阿里云控制台远程连接（推荐新手）

1. 在阿里云控制台点击 **"远程连接"** 按钮
2. 选择 **"VNC 远程连接"** 或 **"Workbench 远程连接"**
3. 首次连接会提示设置密码，设置一个强密码（建议包含大小写字母、数字、特殊字符）

### 方式二：使用 SSH 客户端连接（推荐）

1. **设置/重置密码**：
   - 在控制台点击 **"设置密码"** 或 **"重置密码"**
   - 设置 root 用户密码（记住这个密码！）

2. **使用 SSH 连接**：
   ```bash
   # Windows PowerShell 或 Git Bash
   ssh root@47.97.78.32
   
   # 输入你设置的密码
   ```

---

## 📦 第二步：安装 Docker 和 Docker Compose

连接成功后，在服务器上执行以下命令：

```bash
# 1. 更新系统包
sudo apt update && sudo apt upgrade -y

# 2. 安装必要的工具
sudo apt install -y curl wget git

# 3. 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 4. 启动 Docker 服务
sudo systemctl start docker
sudo systemctl enable docker

# 5. 将当前用户添加到 docker 组（避免每次都用 sudo）
sudo usermod -aG docker $USER

# 6. 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 7. 验证安装
docker --version
docker-compose --version
```

**注意**：如果执行 `docker` 命令提示权限不足，需要重新登录：
```bash
# 退出 SSH 连接，然后重新连接
exit
# 重新连接
ssh root@47.97.78.32
```

---

## 📤 第三步：上传项目代码

### 方式一：使用 Git（推荐）

```bash
# 1. 安装 Git（如果还没安装）
sudo apt install -y git

# 2. 克隆项目
cd /root
git clone <你的项目仓库地址>
cd wenote
```

### 方式二：使用 SCP 上传（Windows）

在本地 PowerShell 或 Git Bash 中执行：

```bash
# 上传整个项目文件夹
scp -r e:\a\wenote root@47.97.78.32:/root/

# 或者只上传必要的文件
scp -r e:\a\wenote root@47.97.78.32:/root/wenote
```

### 方式三：使用 SFTP 工具
- Windows: WinSCP、FileZilla
- Mac: Cyberduck、FileZilla

---

## ⚙️ 第四步：配置环境变量

```bash
# 进入项目目录
cd /root/wenote

# 1. 复制环境变量模板
cp .env.example .env

# 2. 编辑 .env 文件
nano .env
# 或使用 vim
# vim .env
```

**编辑 `.env` 文件**，设置以下内容：

```bash
# 数据库密码（生产环境请使用强密码！）
MYSQL_ROOT_PASSWORD=你的强密码123!@#

# 智谱AI配置（可选，不配置不影响基础功能）
ZHIPU_API_KEY=your_zhipu_api_key_here

# JWT密钥（生产环境必须修改！至少32位随机字符串）
JWT_SECRET=你的随机密钥至少32位字符abcdefghijklmnopqrstuvwxyz123456
```

**保存文件**：
- nano: `Ctrl + O` 保存，`Ctrl + X` 退出
- vim: 按 `i` 进入编辑模式，编辑后按 `Esc`，输入 `:wq` 保存退出

---

## 🔐 第五步：配置后端配置文件

```bash
# 进入后端配置目录
cd /root/wenote/wenote-backend/config

# 复制配置模板
cp config.example.yaml config.yaml

# 编辑配置文件
nano config.yaml
```

**修改 `config.yaml` 中的关键配置**：

```yaml
server:
  port: 8080
  mode: release  # 生产环境改为 release

database:
  host: mysql  # Docker Compose 中使用服务名
  port: 3306
  username: root
  password: 你的数据库密码  # 与 .env 中的 MYSQL_ROOT_PASSWORD 一致
  dbname: wenote

jwt:
  secret: 你的JWT密钥  # 与 .env 中的 JWT_SECRET 一致
  expire: 168

ai:
  zhipu:
    api_key: your-zhipu-api-key-here  # 可选
```

---

## 🚀 第六步：启动服务

```bash
# 回到项目根目录
cd /root/wenote

# 给启动脚本添加执行权限
chmod +x start.sh

# 启动服务
./start.sh
```

或者直接使用 Docker Compose：

```bash
# 构建并启动所有服务
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

---

## 🔥 第七步：配置防火墙和安全组

### 1. 阿里云安全组配置

在阿里云控制台 → 安全组 → 配置规则，开放以下端口：

- **80** (HTTP) - 前端访问
- **8080** (后端API) - 可选，如果不需要外部访问可不开放
- **22** (SSH) - 远程连接（默认已开放）

### 2. 服务器防火墙配置（Ubuntu UFW）

```bash
# 安装 UFW（如果未安装）
sudo apt install -y ufw

# 允许 SSH（重要！先允许 SSH，避免被锁在外面）
sudo ufw allow 22/tcp

# 允许 HTTP
sudo ufw allow 80/tcp

# 允许后端端口（可选）
sudo ufw allow 8080/tcp

# 启用防火墙
sudo ufw enable

# 查看防火墙状态
sudo ufw status
```

---

## ✅ 第八步：验证部署

### 1. 检查服务状态

```bash
# 查看容器运行状态
docker-compose ps

# 应该看到三个容器都在运行：
# - wenote-mysql
# - wenote-backend
# - wenote-frontend
```

### 2. 访问应用

在浏览器中访问：
- **前端**: `http://47.97.78.32`
- **后端API**: `http://47.97.78.32:8080`（如果开放了端口）

### 3. 查看日志（如有问题）

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f mysql
```

---

## 🛠️ 常用管理命令

### 服务管理

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 重新构建并启动
docker-compose up -d --build

# 查看日志
docker-compose logs -f

# 进入容器（调试用）
docker exec -it wenote-backend sh
docker exec -it wenote-mysql bash
```

### 数据备份

```bash
# 备份 MySQL 数据
docker exec wenote-mysql mysqldump -u root -p你的密码 wenote > backup_$(date +%Y%m%d).sql

# 恢复数据
docker exec -i wenote-mysql mysql -u root -p你的密码 wenote < backup_20240126.sql
```

---

## ⚠️ 常见问题

### 1. Docker 命令需要 sudo

**解决**：重新登录 SSH，或执行：
```bash
sudo usermod -aG docker $USER
newgrp docker
```

### 2. 端口被占用

**解决**：检查端口占用
```bash
sudo netstat -tlnp | grep :80
sudo netstat -tlnp | grep :8080
```

### 3. 无法访问前端

**检查**：
- 安全组是否开放 80 端口
- 防火墙是否允许 80 端口
- 容器是否正常运行：`docker-compose ps`

### 4. 数据库连接失败

**检查**：
- `.env` 和 `config.yaml` 中的数据库密码是否一致
- MySQL 容器是否健康：`docker-compose ps`
- 查看 MySQL 日志：`docker-compose logs mysql`

### 5. 服务启动失败

**排查**：
```bash
# 查看详细错误日志
docker-compose logs

# 检查配置文件语法
cat .env
cat wenote-backend/config/config.yaml
```

---

## 🔒 安全建议

1. **修改默认密码**：数据库密码、JWT 密钥都要使用强密码
2. **定期更新**：`sudo apt update && sudo apt upgrade`
3. **配置 SSL**：生产环境建议配置 HTTPS（使用 Nginx + Let's Encrypt）
4. **定期备份**：设置自动备份数据库
5. **监控日志**：定期检查应用日志，发现异常及时处理

---

## 📞 需要帮助？

如果遇到问题：
1. 查看日志：`docker-compose logs -f`
2. 检查容器状态：`docker-compose ps`
3. 查看项目 Issues：https://github.com/Rebornbugkiller/WENOTE/issues

---

**祝部署顺利！🎉**
