# Mint-File

`mint-file` 是一个基于 [Hertz](https://www.cloudwego.io/docs/hertz/) 的轻量级服务，封装了与 [火山引擎 TOS](https://www.volcengine.com/product/tos) 的交互，支持上传、下载等操作，并通过简单的 Web 接口对外提供服务。

---

## ✨ 功能特性

- ✅ 封装 TOS SDK，简化调用
- ✅ 支持上传、下载、文件访问公开化
- ✅ Web API 接口基于 Hertz 框架实现
- ✅ 模块结构清晰，可按需拓展（如监听器、异步处理等）

---

## 📁 目录结构

```bash
mint-file/
├── service/
│   ├── download/
│   │   └── tos.go            # 下载服务封装
│   ├── listener/
│   │   └── tos.go            # 对 TOS 操作的监听或回调处理
│   └── upload/
│       ├── tos.go            # 上传服务封装
│       ├── public.go         # 设置对象访问权限为公开
│       └── download.go       # 上传后下载访问相关逻辑
├── go.mod                    # Go 模块定义
├── main.go                   # Hertz 启动入口
├── README.md                 # 项目说明文档
└── upload.go                 # 通用上传入口接口（或路由注册）
