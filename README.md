## 📦 Mint-File 项目说明

`Mint-File` 是一个支持多存储后端的文件上传下载服务，具备良好的可扩展性与模块化结构。

---

### ✨ 功能特性

* ✅ **多存储后端支持**：兼容火山引擎 TOS 和 MinIO
* ✅ **统一接口设计**：相同 API 屏蔽不同后端实现细节
* ✅ **对象管理功能**：支持上传、下载、公开访问
* ✅ **模块化架构**：便于功能扩展与替换

---

### 📁 项目结构

```bash
mint-file/
├─ main.go          # 启动入口，初始化服务
├─ upload.go        # 上传服务统一入口
├─ download.go      # 下载服务统一入口
├─ tosService/
│  ├─ upload/       # 各上传实现（minio.go, tos.go）
│  ├─ download/     # 各下载实现（minio.go, tos.go）
│  ├─ parser/       # 文件解析器（支持 CSV、DOCX 等）
│  ├─ listener/     # 存储事件监听（如 TOS 回调）
│  └─ public.go     # 公共访问权限处理
```

---

### ⚙️ 示例配置（YAML）

```yaml
file:
  switch: tos
  tos:
    tos_endpoint: your_tos_endpoint
    tos_access_key: your_tos_access_key
    tos_access_secret: your_tos_access_secret
    tos_region: your_tos_region
    tos_bucket_name: your_tos_bucket_name
    tos_location:
      picture: test/picture/
      file: test/file/
    tos_shard: 5242880 # 5*1024*1024
  minio:
    minio_endpoint: your_tos_endpoint
    minio_access_key: your_tos_access_key
    minio_access_secret: your_tos_access_secret
    minio_bucket_name: your_tos_bucket_name
    minio_location:
      picture: test/picture/
      file: test/file/
```
