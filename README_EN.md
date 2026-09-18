<div align="center">

# 🍀 LuckyPanel

### All-in-One Modern Server & Network Management Dashboard

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Vue](https://img.shields.io/badge/Vue-3.5+-4FC08D?style=flat&logo=vue.js)](https://vuejs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.9+-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org)
[![License](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com)

<p align="center">
  <b>Combining the best features of Nginx-UI, Lucky, ddns-go, and 1Panel</b><br>
  A lightweight, full-featured server dashboard combining hardware metrics, Nginx reverse proxy, dynamic DNS (DDNS), TCP/UDP port forwarding, Linux firewall management, Docker container & Compose stack control, Web Terminal, and automated SSL certificate lifecycle.
</p>

[Features](#-features) • [Quick Start](#-quick-start) • [Configuration](#-configuration) • [Build from Source](#-build-from-source) • [Acknowledgements](#-acknowledgements)

</div>

---

## 🚀 Features

- 🖥️ **Server Dashboard**: Real-time CPU, RAM, Disk I/O, Network card speed charts, and system load.
- 🛡️ **Firewall Management**: Adaptive support for UFW, Firewalld, and iptables. Single and range port opening with CIDR IP whitelist support. Built-in protection preventing accidental SSH (22) and panel port closure.
- 🌐 **Nginx Reverse Proxy**: Visual virtual host wizards, Docker container auto-detection, Monaco configuration editor with syntax highlight and `nginx -t` validation.
- 🔀 **TCP/UDP Port Forwarding**: High-performance stream relay pool inspired by Lucky, active connection tracker, realtime bytes transferred (Rx/Tx), UPnP router auto-mapping, and automatic firewall opening.
- ⚡ **Dynamic DNS (DDNS)**: Multi-provider support (Cloudflare, Aliyun, Tencent Cloud, Huawei Cloud, Webhooks), IPv4/IPv6 dual-stack, periodic change detection.
- 🐳 **Docker & Compose**: Manage container lifecycle, stream live logs, inspect images, and edit/deploy multi-container Compose stacks online.
- 📜 **Log Center**: Unified search and filter across Nginx access/error logs, Docker container streams, and Linux syslog/journal.
- 🔒 **Automated SSL (ACME)**: Automatic Let's Encrypt / ZeroSSL certificate issuance (HTTP-01 & DNS-01) with silent auto-renewal.
- 💻 **Web Terminal**: In-browser local PTY Shell and multi-host SSH bastion client based on xterm.js.
- 📦 **Single Standalone Binary**: Complete Vue 3 frontend assets embedded inside the Go executable.

---

## 🛠️ Quick Start

### 1. Run Prebuilt Binary

```bash
chmod +x server-manager
./server-manager serve
```

Visit `http://<YOUR-SERVER-IP>:9000` in your web browser.

### 2. Reset or Retrieve Initial Password

```bash
./server-manager reset-password --config app.ini
```

### 3. Docker One-Click Run

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

---

## 🔨 Build from Source

### Prerequisites
- Go 1.25+
- Node.js 20+ & pnpm

```bash
# 1. Build Vue 3 Frontend
cd app
pnpm install
pnpm run build
cd ..

# 2. Build Single Executable with Embedded Frontend
go build -o server-manager main.go
```

---

## 🤝 Acknowledgements

LuckyPanel is built upon and inspired by these great open-source projects:
- **[0xJacky/nginx-ui](https://github.com/0xJacky/nginx-ui)**
- **[gdy666/lucky](https://github.com/gdy666/lucky)**
- **[jeessy2/ddns-go](https://github.com/jeessy2/ddns-go)**
- **[1Panel-dev/1Panel](https://github.com/1Panel-dev/1Panel)**
- **[go-acme/lego](https://github.com/go-acme/lego)**
- **[shirou/gopsutil](https://github.com/shirou/gopsutil)**
