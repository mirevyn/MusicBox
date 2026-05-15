# MusicBox

MusicBox 是一个前后端分离的在线音乐盒项目，包含音乐播放、歌曲管理、歌单、收藏、每日推荐、后台管理和 AI 助手等功能。项目适合作为个人音乐库、私有化部署的音乐管理系统。

当前仓库包含：

- `musicbox-frontend`：Vue 3 + Vite + TypeScript 前端应用。
- `musicbox-backend`：Go + Gin + GORM 后端服务。
- `docker-compose.yml`：MySQL、后端、前端 Nginx 的一键部署编排。
- `musicbox-backend/musicbox.sql`：基础数据库结构和必要种子数据。

## 功能特性

- 用户注册、登录、JWT 鉴权、个人资料和密码修改。
- 音乐库浏览、搜索、播放队列、播放器控制。
- 歌曲收藏、我喜欢的音乐、每日推荐、播放历史记录。
- 歌单创建、编辑、公开推荐、搜索、添加和移除歌曲。
- 管理后台：仪表盘、歌曲管理、用户管理、歌单审核、系统设置、报表导出。
- 上传音乐、封面、歌词、头像等资源文件。
- AI 音乐助手：支持 Ollama 或 OpenAI 兼容接口，支持流式聊天。
- Swagger API 文档和健康检查接口。
- Docker Compose 一键部署，前端 Nginx 反向代理后端 API、上传资源和 WebSocket。

## 技术栈

| 模块 | 技术                                                               |
| ---- | ------------------------------------------------------------------ |
| 前端 | Vue 3, Vite, TypeScript, Pinia, Vue Router, UnoCSS, Axios, ECharts |
| 后端 | Go, Gin, GORM, MySQL, JWT, WebSocket, Swagger                      |
| 部署 | Docker, Docker Compose, Nginx                                      |
| AI   | Ollama 或 OpenAI 兼容 API                                          |

## 目录结构

```text
MusicBox
├── docker-compose.yml              # 标准部署编排
├── .env.example                    # Docker 环境变量示例
├── musicbox-backend
│   ├── api                         # HTTP / WebSocket 控制器
│   ├── config                      # 后端配置示例
│   ├── docs                        # Swagger 产物
│   ├── internal                    # 配置、模型、服务层
│   ├── middleware                  # 鉴权、CORS 等中间件
│   ├── router                      # 路由注册
│   ├── utils                       # 通用工具
│   └── musicbox.sql                # 基础数据库初始化脚本
└── musicbox-frontend
    ├── public
    ├── src
    │   ├── api                     # 前端 API 封装
    │   ├── components              # 通用、播放器、后台等组件
    │   ├── layout                  # 客户端和后台布局
    │   ├── router                  # 前端路由
    │   ├── stores                  # Pinia 状态
    │   └── views                   # 页面视图
    └── nginx.conf                  # 生产镜像 Nginx 配置
```

## 环境要求

推荐使用 Docker Compose 启动完整环境。

- Docker 和 Docker Compose。
- 本地开发后端：Go 1.25.x。
- 本地开发前端：Node.js 22.x、pnpm。
- 数据库：MySQL 8.0。
- 可选 AI 服务：Ollama，或任何 OpenAI 兼容接口。

## 快速启动：Docker Compose

1. 复制环境变量文件：

```bash
cp .env.example .env
```

2. 修改 `.env`，至少填写：

```env
MYSQL_ROOT_PASSWORD=your_secure_password_here
JWT_SECRET=your_random_jwt_secret_here
```

3. 启动服务：

```bash
docker compose up --build
```

4. 访问应用：

- 前端入口：`http://localhost`
- 默认前端端口来自 `.env` 中的 `FRONTEND_PORT`，未设置时为 `80`。

如果本机 80 端口被占用，可以在 `.env` 中增加：

```env
FRONTEND_PORT=8080
```

然后访问 `http://localhost:8080`。

## 初始化数据说明

- `musicbox.sql` 会创建基础表结构和必要种子数据。
- 为了让公开仓库保持轻量，项目不内置演示歌曲、封面、歌词或头像文件。
- 如需完整播放体验，请通过后台上传自己的合法音乐资源，或把自备文件放到 `musicbox-backend/uploads` 下的对应目录。
- MySQL 官方镜像只会在数据目录为空时执行初始化 SQL。如果 `mysql-data/` 已存在，重新启动不会再次导入基础数据。

基础脚本会创建管理员种子账号，管理员头像默认为空，前端会使用用户名首字母作为兜底头像。首次部署后请立即重置管理员密码，或在部署前替换 SQL 中的管理员密码哈希。

## 本地开发

### 1. 启动 MySQL

可以使用本地 MySQL，也可以只启动 Docker 中的 MySQL：

```bash
docker compose up mysql
```

如果使用本地 MySQL，请手动创建数据库并导入基础脚本：

```bash
mysql -uroot -p -e "CREATE DATABASE IF NOT EXISTS music_box DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
mysql -uroot -p music_box < musicbox-backend/musicbox.sql
```

### 2. 配置后端

复制配置文件：

```bash
cd musicbox-backend
cp config/config.example.yaml config/config.local.yaml
```

修改 `config/config.local.yaml` 中的数据库 DSN、JWT Secret、上传路径、CORS 和 AI 配置。

也可以使用环境变量覆盖配置，例如：

```bash
MUSICBOX_MYSQL_DSN="root:password@tcp(127.0.0.1:3306)/music_box?charset=utf8mb4&parseTime=True&loc=Local"
MUSICBOX_JWT_SECRET="replace_with_a_long_random_secret"
MUSICBOX_HTTP_ADDR=":8000"
```

启动后端：

```bash
go run main.go
```

后端默认监听 `http://localhost:8000`。

常用入口：

- 健康检查：`http://localhost:8000/healthz`
- Swagger 文档：`http://localhost:8000/swagger/index.html`
- API 前缀：`http://localhost:8000/api`

### 3. 启动前端

```bash
cd musicbox-frontend
pnpm install
```

本地开发时，前端和后端通常不在同一个端口，需要设置后端地址。可以创建 `musicbox-frontend/.env.local`：

```env
VITE_BACKEND_URL=http://localhost:8000
```

启动前端开发服务器：

```bash
pnpm run dev
```

访问：

```text
http://localhost:5173
```

## 环境变量说明

Docker Compose 使用仓库根目录的 `.env`。

| 变量                         | 必填 | 说明                                               |
| ---------------------------- | ---- | -------------------------------------------------- |
| `MYSQL_ROOT_PASSWORD`        | 是   | MySQL root 密码，也是后端容器连接 MySQL 使用的密码 |
| `JWT_SECRET`                 | 是   | JWT 签名密钥，请使用足够长的随机字符串             |
| `MUSICBOX_HTTP_ADDR`         | 否   | 后端监听地址，默认 `:8000`                         |
| `GIN_MODE`                   | 否   | Gin 运行模式，生产建议 `release`                   |
| `ALLOWED_ORIGINS`            | 否   | CORS 允许来源，多个来源用英文逗号分隔              |
| `FRONTEND_PORT`              | 否   | 前端容器映射到宿主机的端口，默认 `80`              |
| `AI_BASE_URL`                | 否   | AI 服务地址，例如 Ollama 或 OpenAI 兼容 API        |
| `AI_MODEL`                   | 否   | AI 模型名，默认可使用 `qwen2.5:7b`                 |
| `AI_PROVIDER`                | 否   | AI 服务类型，建议填写 `ollama` 或 `openai`         |
| `AI_API_KEY`                 | 否   | OpenAI 兼容服务 API Key，Ollama 通常不需要         |
| `AI_ALLOW_PRIVATE_BASE_URLS` | 否   | 是否允许前端配置内网 AI 地址，默认 `false`         |
| `AI_REQUEST_TIMEOUT_SECONDS` | 否   | AI 请求超时时间，默认 `300` 秒                     |

后端还支持以 `MUSICBOX_` 前缀的环境变量覆盖 YAML 配置，例如：

- `MUSICBOX_CONFIG_PATH`
- `MUSICBOX_MYSQL_DSN`
- `MUSICBOX_JWT_SECRET`
- `MUSICBOX_JWT_EXPIRES_HOURS`
- `MUSICBOX_UPLOAD_SONG_PATH`
- `MUSICBOX_UPLOAD_COVER_PATH`
- `MUSICBOX_UPLOAD_LYRIC_PATH`
- `MUSICBOX_UPLOAD_AVATAR_PATH`
- `MUSICBOX_CORS_ALLOWED_ORIGINS`

前端构建可使用：

- `VITE_BACKEND_URL`：后端基础地址。Docker 同源部署时通常留空；本地开发时建议设置为 `http://localhost:8000`。

## AI 配置

AI 助手是可选能力。未配置 `AI_BASE_URL` 时，主应用仍可正常运行，AI 状态会提示不可用或未配置。

### Ollama 示例

```env
AI_BASE_URL=http://host.docker.internal:11434
AI_PROVIDER=ollama
AI_MODEL=qwen2.5:7b
AI_API_KEY=
```

如果 Ollama 和后端都在同一台宿主机，但后端运行在 Docker 容器内，通常不能在容器里使用 `127.0.0.1:11434` 访问宿主机服务。Windows 和 macOS Docker Desktop 可优先尝试 `host.docker.internal`。

### OpenAI 兼容接口示例

```env
AI_BASE_URL=https://your-openai-compatible.example.com/v1
AI_PROVIDER=openai
AI_MODEL=your-model-name
AI_API_KEY=your_api_key
```

## API 文档

后端集成 Swagger 文档：

```text
http://localhost:8000/swagger/index.html
```

Docker Compose 默认只暴露前端 Nginx 端口，后端端口主要在容器网络内部使用。如果需要在 Docker 部署中直接访问 Swagger，可以临时给 `backend` 服务添加端口映射，例如：

```yaml
ports:
  - "8000:8000"
```

## 测试和构建

后端测试：

```bash
cd musicbox-backend
go test ./...
```

前端类型检查和生产构建：

```bash
cd musicbox-frontend
pnpm install
pnpm run build
```

Docker 构建验证：

```bash
docker compose up --build
```

建议发布前至少确认：

- 前端首页能打开。
- 登录和注册流程可用。
- 后端 `/healthz` 返回 `{"status":"ok"}`。
- 后端 Swagger 页面可访问。
- 上传目录挂载正常。
- AI 功能在未配置时能正常提示，在配置后能正常连接。

## GitHub 发布前检查

本项目适合公开仓库，但提交前请确认不要把本地密钥、数据库和上传资源一起推上去。

推荐检查：

```bash
git status --short --branch
git ls-files .env musicbox-backend/config/config.local.yaml musicbox-backend/config/config.prod.yaml mysql-data musicbox-backend/uploads
git diff --stat
```

这些内容不应该出现在 Git 跟踪列表中：

- `.env`
- `musicbox-backend/config/config.local.yaml`
- `musicbox-backend/config/config.prod.yaml`
- `mysql-data/`
- `musicbox-backend/uploads/`
- `musicbox-frontend/node_modules/`
- `musicbox-frontend/dist/`

首次推送示例：

```bash
git branch -M main
git remote add origin https://github.com/<your-name>/MusicBox.git
git add -A
git commit -m "docs: prepare GitHub release"
git push -u origin main
```

如果你希望保留当前 `master` 分支，也可以不执行 `git branch -M main`。

## 部署建议

- 生产环境务必修改 `MYSQL_ROOT_PASSWORD` 和 `JWT_SECRET`。
- 首次部署后立即修改管理员密码，或在初始化 SQL 中替换管理员密码哈希。
- 不要把真实音乐文件、封面、歌词、头像或数据库目录提交到公开仓库。
- 给 `musicbox-backend/uploads` 做持久化和备份。
- 如果部署到公网，建议放在 HTTPS 反向代理后面。
- CORS 只开放真实前端域名，不要在生产环境随意放开来源。
- 上传音乐资源前，请确认你拥有合法使用和传播权限。

## 常见问题

### 修改 SQL 后为什么没有重新初始化数据库？

MySQL 容器只会在数据目录为空时执行 `/docker-entrypoint-initdb.d` 中的 SQL。已经存在 `mysql-data/` 时，修改 SQL 不会自动生效。开发环境可以在备份后清空数据目录再重建；生产环境请使用迁移脚本或手动执行 SQL。

### Docker 启动后为什么 Swagger 打不开？

默认 Docker Compose 只对外暴露前端 Nginx，后端 `8000` 端口未映射到宿主机。可以本地开发时直接运行后端访问 Swagger，或临时给 `backend` 服务添加端口映射。

### 前端开发环境请求后端失败怎么办？

检查 `musicbox-frontend/.env.local`：

```env
VITE_BACKEND_URL=http://localhost:8000
```

同时确认后端 `config.local.yaml` 中的 CORS 允许来源包含：

```yaml
cors:
  allowed_origins:
    - "http://localhost:5173"
```

### 为什么初始化后没有歌曲或封面？

公开仓库不内置音乐媒体资源。请上传自己的音乐资源，或把自备文件放到 `musicbox-backend/uploads/music`、`uploads/covers`、`uploads/lyrics` 等目录下。

### Git 提示 `.git/index.lock` 怎么办？

确认没有正在运行的 Git 进程后，可以删除 `.git/index.lock`，再重新执行 `git status` 或提交命令。不要在 Git 命令仍在运行时删除这个文件。

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
