<div align="center">

# 🍀 LuckyPanel (幸运面板)

### 一体化开源服务器与网络管理平台
**All-in-One Modern Server & Network Management Dashboard**

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Vue](https://img.shields.io/badge/Vue-3.5+-4FC08D?style=flat&logo=vue.js)](https://vuejs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.9+-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org)
[![License](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com)

<p align="center">
  <b>融合 Nginx-UI + Lucky + ddns-go + 1Panel 核心精髓</b><br>
  集服务器总览、Web 反向代理、动态域名 DDNS、端口转发、防火墙管理、Docker 容器与 Compose 编排、Web 终端及 SSL 证书自动化于一体的轻量级综合控制面板。
</p>

[功能特性](#-功能特性) • [架构理念](#-架构理念) • [快速开始](#-快速开始) • [使用说明](#-使用说明) • [源码构建](#-源码构建) • [致谢与开源说明](#-致谢与开源说明)

</div>

---

## 📖 项目简介

**LuckyPanel** 旨在解决传统服务器管理工具分散、配置繁琐、资源占用高的问题。我们不从零制造轮子，而是精选并融合了社区中优秀且成熟的开源技术资产：

1. **Web 服务与反代**：汲取 **[0xJacky/nginx-ui](https://github.com/0xJacky/nginx-ui)** 优雅的可视化配置流、Monaco 代码编辑器、语法检测与平滑重载能力。
2. **网络中继与穿透**：移植 **[gdy666/lucky](https://github.com/gdy666/lucky)** 强大高性能的 TCP/UDP 端口转发池、实时连接监控、流量统计与 UPnP 路由器自动映射。
3. **动态域名解析**：融合 **[jeessy2/ddns-go](https://github.com/jeessy2/ddns-go)** 与 Lucky 的多厂商云解析驱动（Cloudflare、阿里云、腾讯云、华为云等）以及 IPv4/IPv6 双栈智能变动探测。
4. **系统防火墙管理**：吸收 **[1Panel](https://github.com/1Panel-dev/1Panel)** 的防火墙自适应抽象设计，自动兼容 UFW / Firewalld / iptables，内置 SSH 与管理端口防失联保护，并与端口转发一键联动放行。
5. **Docker 与 Compose 编排**：直接对接 Docker 官方 Moby SDK，实现容器全生命周期控制、实时日志流式查看、本地镜像管理与多容器 Compose 栈可视化管理。
6. **一体化单二进制交付**：将 Vue 3 前端静态资源完整嵌入 Go 后端二进制中，零复杂运行环境依赖，开箱即用。

---

## 🚀 功能特性

### 1. 🖥️ 服务器总览 (Dashboard)
- 实时 CPU 占用率、内存构成（Used / Cached / Buffers）、多磁盘挂载点容量与 IO 吞吐监控。
- 多网卡实时收发速率（Rx/Tx Speed）时序图表与系统负载（Load Average）展示。
- 关键状态总览卡片：快速掌握当前运行的 Nginx 站点数、活动转发连接、DDNS 状态及 Docker 容器分布。

### 2. 🛡️ 防火墙端口管理 (Firewall Management)
- **多底层自适应**：智能探测宿主机防火墙组件（自动兼容 `UFW`、`Firewalld`，提供 `iptables` 规则兜底）。
- **灵活放行规则**：支持单个端口（如 `8848`）或端口范围（如 `9000-9050`），支持 `TCP`、`UDP` 及 `TCP/UDP` 双协议放行。
- **安全白名单控制**：支持限定特定 IP / CIDR 白名单来源，强化内部数据库及私有服务的访问边界。
- **防锁死保护机制**：系统内置 SSH 默认端口（22）及当前面板 Web 端口的删除拦截，杜绝误操作导致远程失联。

### 3. 🌐 Web 服务与反向代理 (Nginx UI)
- **可视化建站向导**：快速创建 HTTP / HTTPS 反向代理站点、静态 SPA 托管与 TCP/UDP 四层流转发配置。
- **Docker 容器联动**：创建反向代理时可直接下拉选择本地正在运行的 Docker 容器及其映射端口。
- **专业 Monaco 编辑器**：支持 Nginx 语法高亮、自动补全、`nginx -t` 语法检测、版本 Diff 对比与配置历史一键回滚。

### 4. 🔀 端口转发与 NAT 穿透 (Port Forwarding)
- **高性能中继转发**：基于 Go 原生协程与缓冲池技术实现低延迟 TCP/UDP 全双工中继。
- **实时监控分析**：实时记录并刷新活跃连接数、当前传输速率、累积接收（Rx）与发送（Tx）字节大小。
- **UPnP 自动映射**：家庭宽带/内网环境下支持一键向兼容 UPnP 的路由器协商申请公网端口映射。
- **防火墙联动打通**：添加转发规则时勾选【自动在防火墙放行该端口】，平台即时在防火墙新增对应放行规则。

### 5. ⚡ 动态域名解析 (DDNS)
- **主流厂商全面覆盖**：内置 Cloudflare、阿里云 (AliDNS)、腾讯云 (DNSPod)、华为云及通用 Webhook 回调协议。
- **多途径 IP 探测**：支持外部公网探针轮询、物理网卡直读提取（IPv4/IPv6 正则过滤）、STUN NAT 探测。
- **智能增量更新**：后台定时轮询（周期可配置），仅在公网 IP 发生实质性漂移时才调用云厂商 API，附带完备日志与变动提醒。

### 6. 🐳 Docker 容器与 Compose 栈管理
- **容器运维**：启动、停止、重启、强制终止 (Kill)、删除容器；一键查看实时容器控制台日志流。
- **本地镜像**：镜像列表展示、占用体积统计与无用镜像清理。
- **Compose 编排中心**：在线编写 `docker-compose.yml`，提供多服务栈的一键部署 (`compose up`)、停止与重启更新。

### 7. 📜 统一日志中心 (Log Center)
- **全系统日志汇聚**：统一检索 Nginx Access/Error 日志、Docker 容器输出、Linux systemd journal / syslog 与面板操作审计日志。
- **实时流式追踪**：支持搜索关键词高亮、指定行数截取与终端风格的高性能滚动呈现。

### 8. 🔒 SSL 证书全自动生命周期 (ACME)
- 集成 `lego` ACME 客户端，支持申请 Let's Encrypt 与 ZeroSSL 证书。
- 支持 **HTTP-01**（无感 Nginx 验证）与 **DNS-01**（支持通配符泛域名，直接复用 DDNS 云厂商秘钥凭据）。
- 到期前 30 天静默自动续期，更新后自动触发 Nginx reload 平滑重载。

### 9. 💻 Web 终端 (Web Terminal)
- **本地 PTY Shell**：直接在浏览器中打开全功能 Linux Bash / Zsh 终端。
- **SSH 远程堡垒机**：内置凭据管理，支持多台远程主机 SSH 会话跳转与多 Tab 切换，基于 xterm.js 打造丝滑操作体验。

---

## 📦 架构理念

```text
┌─────────────────────────────────────────────────────────────┐
│                   LuckyPanel 统一管理门户                    │
│            (Vue 3 + Vite + TypeScript + Ant Design)         │
└──────────────────────────────┬──────────────────────────────┘
                               │ HTTP / WebSocket RESTful API
┌──────────────────────────────▼──────────────────────────────┐
│                  LuckyPanel 单二进制核心服务                │
│                 (Go 1.25+ / Gin / GORM / SQLite)            │
├──────────────┬──────────────┬──────────────┬────────────────┤
│ Nginx 站点   │ 防火墙控制   │ 端口转发池   │ DDNS 同步引擎  │
│ 配置/语法/重载 │ UFW/Firewall │ TCP/UDP 中继 │ CF/阿里/腾讯/HW │
├──────────────┼──────────────┼──────────────┼────────────────┤
│ Docker 管理  │ 统一日志中心 │ Web 终端     │ SSL ACME 自动化│
│ 容器/Compose │ Nginx/系统日志│ PTY / SSH 网关│ Let's Encrypt  │
└──────────────┴──────────────┴──────────────┴────────────────┘
```

---

## 🛠️ 快速开始

### 方式一：独立二进制直接运行 (推荐)

从 Releases 页面下载适合您系统架构的预编译单二进制包：

```bash
# 1. 赋予执行权限
chmod +x server-manager

# 2. 启动服务 (默认监听在 0.0.0.0:9000，使用当前目录下的 app.ini 配置)
./server-manager serve

# 3. 后台守护运行 (使用 nohup 或配置 systemd)
nohup ./server-manager serve > server.log 2>&1 &
```

打开浏览器访问 `http://<服务器IP>:9000` 即可进入面板。

---

### 方式二：Docker 一键部署

```bash
docker run -d \
  --name luckypanel \
  --restart always \
  --net host \
  -v /etc/nginx:/etc/nginx \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /data/luckypanel:/data \
  luckypanel/luckypanel:latest
```

> **挂载说明**：
> - `--net host`：便于面板直接感知宿主机网络网卡、执行端口转发与防火墙联动。
> - `/var/run/docker.sock`：使面板具备直接管理宿主机 Docker 容器与镜像的能力。
> - `/etc/nginx`：便于面板直接管理宿主机原生的 Nginx 站点配置。
> - `/data/luckypanel`：持久化存储面板的 SQLite 数据库 (`database.db`) 与 Compose 配置文件。

---

### 方式三：Docker Compose 快速编排

创建 `docker-compose.yml`：

```yaml
version: '3.8'

services:
  luckypanel:
    image: luckypanel/luckypanel:latest
    container_name: luckypanel
    restart: always
    network_mode: host
    volumes:
      - /etc/nginx:/etc/nginx
      - /var/run/docker.sock:/var/run/docker.sock
      - ./data:/data
    environment:
      - SERVER_PORT=9000
```

启动服务：
```bash
docker compose up -d
```

---

## 🔑 使用说明

### 1. 初始账号与密码管理
初次启动或需要重置管理员密码时，可执行自带的密码重置工具：

```bash
# 获取/重置 admin 初始密码
./server-manager reset-password --config app.ini
```
终端将输出新生成的安全密码，使用用户名 `admin` 和输出的密码登录控制台。

### 2. 配置文件说明 (`app.ini`)

系统主配置文件默认位于当前目录下的 `app.ini`（可参考 `app.example.ini`）：

```ini
[app]
PageSize  = 20
JwtSecret = "your-custom-jwt-secret"

[server]
Host    = "0.0.0.0"
Port    = 9000
RunMode = "release"

[database]
Name = "database" # 默认在当前目录或 /data 下生成 database.db
```

### 3. 防火墙端口管理演示
1. 进入左侧导航菜单 **【防火墙 (Firewall)】**。
2. 查看当前服务器识别到的防火墙底层状态（UFW / Firewalld / iptables）。
3. 点击 **+ 开放端口**，输入端口（例如 `8080` 或范围 `9000-9050`），选择协议（TCP/UDP/双协议），设定允许来源（支持 `any` 或 `192.168.1.0/24`）。
4. 确认后即时调用底层防火墙放行，并在面板清晰呈现所有放行规则。

### 4. 端口转发与防火墙联动演示
1. 进入左侧导航菜单 **【端口转发 (Port Forward)】**。
2. 点击 **+ 添加转发规则**，填写名称、监听端口（如 `18080`）、目标 IP（如 `192.168.1.100`）及目标端口（如 `80`）。
3. 勾选 **【自动在防火墙放行该端口】** 与 **【尝试 UPnP 路由器映射】**。
4. 保存后规则即时启动，外部流量即可平滑中继，并在列表实时查看到连接数与双向传输流量。

---

## 🔨 源码构建

如果您希望从源码自行编译打包整个项目：

### 前置要求
- Go 1.25+
- Node.js 20+ 及 pnpm

### 1. 构建前端
```bash
cd app
pnpm install
pnpm run build
cd ..
```
前端编译产物将生成在 `app/dist/` 中。

### 2. 编译嵌入后端的单二进制文件
```bash
# 采用 Go embed 将 app/dist 完整静态嵌入
go build -o server-manager main.go
```
执行完成后，根目录即可获得完整的单文件可执行程序 `server-manager`。

---

## 🤝 致谢与开源说明

**LuckyPanel** 的诞生离不开开源社区优秀先驱项目的启发与卓越贡献，特此向以下优秀项目致以最高敬意：

- **[0xJacky/nginx-ui](https://github.com/0xJacky/nginx-ui)**：提供了现代、专业的 Nginx 站点治理交互体系、配置解析与前端骨架。
- **[gdy666/lucky](https://github.com/gdy666/lucky)**：提供了极致轻量、高效的端口转发中继池与 DDNS 探测理念。
- **[jeessy2/ddns-go](https://github.com/jeessy2/ddns-go)**：提供了全面易用的多云厂商 DNS API 驱动实现。
- **[1Panel-dev/1Panel](https://github.com/1Panel-dev/1Panel)**：提供了健壮的 Linux 防火墙适配与容器运维设计思路。
- **[go-acme/lego](https://github.com/go-acme/lego)**：提供了工业级的自动化 ACME 证书签发与全 DNS 服务商驱动。
- **[shirou/gopsutil](https://github.com/shirou/gopsutil)**：提供了跨平台精准的系统硬件资源度量引擎。

本项目遵循 **GPL v3** 开源协议，欢迎广大开发者与运维爱好者提交 Issue、PR 共同完善 LuckyPanel！
