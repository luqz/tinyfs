# tinyfs

一个用 curl 就能上传/下载文件的轻量 HTTP 文件服务。

## 快速开始

### 直接运行

```bash
go build -o tinyfs .
./tinyfs
```

### Docker Compose

```bash
docker compose up -d
```

## 用法

| 操作 | 命令 |
|------|------|
| 上传 (multipart) | `curl -F "file=@/path/to/file" http://localhost:8082/upload` |
| 上传 (raw body) | `curl --data-binary @/path/to/file http://localhost:8082/upload/filename` |
| 下载 | `curl -O http://localhost:8082/download/filename` |
| 查看文件列表 | `curl http://localhost:8082/` |

## 命令行参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-listen` | `:8082` | 监听地址 (e.g. `:9090`, `127.0.0.1:8080`) |
| `-dir` | `uploads` | 文件存储目录 |
