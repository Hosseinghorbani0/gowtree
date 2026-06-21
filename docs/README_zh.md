# gowtree

<p align="center">
  🌍 <b>Readme:</b>
  <a href="../README.md"><img src="https://flagcdn.com/20x15/us.png" alt="English"> English</a> ·
  <a href="README_fa.md"><img src="https://flagcdn.com/20x15/ir.png" alt="Persian"> فارسی</a> ·
  <a href="README_zh.md"><img src="https://flagcdn.com/20x15/cn.png" alt="Chinese"> 中文</a> ·
  <a href="README_tr.md"><img src="https://flagcdn.com/20x15/tr.png" alt="Turkish"> Türkçe</a> ·
  <a href="README_ar.md"><img src="https://flagcdn.com/20x15/sa.png" alt="Arabic"> العربية</a> ·
  <a href="README_ru.md"><img src="https://flagcdn.com/20x15/ru.png" alt="Russian"> Русский</a>
</p>

<p align="center">
  <img src="../assets/banner.svg" alt="gowtree banner" width="100%">
</p>

<p align="center">
  <b>为 Windows 打造的现代目录树工具 — 支持颜色、图标与 JSON/Markdown/HTML 导出。</b><br/>
  单个可执行文件。无需服务器。轻量高效。
</p>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](../LICENSE)
[![Windows](https://img.shields.io/badge/platform-Windows-0078D6.svg)](https://github.com/hosseinghorbani0/gowtree/releases)

---

## ⚡ 快速开始

1. 从 [Releases](https://github.com/hosseinghorbani0/gowtree/releases) 下载 **`gowtree-setup-1.4.0.exe`**
2. 运行安装向导，勾选 **Add to PATH**
3. 打开**新**终端：

```powershell
gowtree -a -s -L 2 --icons
```

或从源码安装：

```powershell
git clone https://github.com/hosseinghorbani0/gowtree.git
cd gowtree
.\install.bat
```

---

## 🧐 为什么选择 gowtree？

| 功能 | Windows `tree` | **gowtree** |
|------|:--------------:|:-----------:|
| 颜色与 Unicode | ❌ | ✅ |
| JSON 输出 | ❌ | ✅ |
| Markdown / HTML | ❌ | ✅ |
| 进度条 | ❌ | ✅ |
| 目录大小 | ❌ | ✅ |
| 正则过滤 | ❌ | ✅ |
| Nerd Font 图标 | ❌ | ✅ |
| 剪贴板 | ❌ | ✅ |
| YAML 配置 | ❌ | ✅ |

---

## 📚 示例

```powershell
gowtree -J > tree.json
gowtree --markdown --out TREE.md
gowtree -R "\.go$" --sort size -r --icons
gowtree --clip
```

---

## 📖 项目故事

> 约 **20%** 由 AI 协助完成（架构、文档、安装程序）。  
> 约 **80%** 为**手工编写** — 在**伊朗断网**期间，于伊朗、美国与以色列冲突背景下完成。

若您认同这份坚持，欢迎 ⭐ Star 本仓库。

---

## 🔗 链接

[GitHub](https://github.com/hosseinghorbani0/gowtree) · [Releases](https://github.com/hosseinghorbani0/gowtree/releases) · [Issues](https://github.com/hosseinghorbani0/gowtree/issues)

**作者：[Hossein Ghorbani](https://github.com/hosseinghorbani0)** · MIT 许可证
