# 🚀 LuckyPanel 上传至 GitHub 完整操作指引

本项目已经完成所有功能代码整合、前端构建、后端嵌入与本地 Git 初始化。请按照以下步骤将其发布至您的 GitHub 仓库。

---

## 步骤 1：在 GitHub 上创建新的空仓库

1. 登录您的 GitHub 账号：https://github.com
2. 点击右上角的 **+** 按钮，选择 **New repository**（新建仓库）。
3. 填写仓库信息：
   - **Repository name**: `luckypanel`（或您自定义的仓库名称）
   - **Description**: `All-in-One Modern Server & Network Management Dashboard (Nginx UI + Lucky + ddns-go + 1Panel)`
   - **Public / Private**: 根据需求选择公开（Public）或私有（Private）
   - **重要**：**不要** 勾选 "Add a README file"、"Add .gitignore" 或 "Choose a license"（因为本地已经包含这些文件）。
4. 点击绿色按钮 **Create repository**。

---

## 步骤 2：在本地终端执行推送命令

在当前项目根目录下打开终端，执行以下两条命令即可完成上传：

```bash
# 1. 关联远程 GitHub 仓库地址 (请将 YOUR_USERNAME 替换为您的 GitHub 用户名或使用您创建好的仓库地址)
git remote add origin https://github.com/YOUR_USERNAME/luckypanel.git

# 2. 推送代码至 GitHub main 分支
git push -u origin main
```

> **提示**：如果您的 GitHub 账号配置了 SSH 密钥，可以使用 SSH 地址推送：
> ```bash
> git remote add origin git@github.com:YOUR_USERNAME/luckypanel.git
> git push -u origin main
> ```

---

## 步骤 3：使用 GitHub CLI（可选快速方式）

如果您安装了官方 GitHub CLI 工具 (`gh`)，可直接一条命令完成自动创建仓库与推送：

```bash
gh repo create luckypanel --public --source=. --remote=origin --push
```

---

## 验证与后续建议

- 推送成功后，刷新 GitHub 仓库页面即可看到精美排版的 `README.md` 与全部模块源码。
- 项目二进制文件 `server-manager` 及临时数据库已配置在 `.gitignore` 中，保证不会触发 GitHub 100MB 单文件限制。
- 发布首个版本时，建议在 GitHub 仓库右侧点击 **Releases** -> **Draft a new release** 打上 `v1.0.0` 标签。
