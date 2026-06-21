# gowtree

<p align="center">
  🌍 <b>Readme:</b>
  <a href="README.md"><img src="https://flagcdn.com/20x15/us.png" alt="English"> English</a> ·
  <a href="docs/README_fa.md"><img src="https://flagcdn.com/20x15/ir.png" alt="Persian"> فارسی</a> ·
  <a href="docs/README_zh.md"><img src="https://flagcdn.com/20x15/cn.png" alt="Chinese"> 中文</a> ·
  <a href="docs/README_tr.md"><img src="https://flagcdn.com/20x15/tr.png" alt="Turkish"> Türkçe</a> ·
  <a href="docs/README_ar.md"><img src="https://flagcdn.com/20x15/sa.png" alt="Arabic"> العربية</a> ·
  <a href="docs/README_ru.md"><img src="https://flagcdn.com/20x15/ru.png" alt="Russian"> Русский</a>
</p>

<p align="center">
  <img src="assets/banner.svg" alt="gowtree banner" width="100%">
</p>

<p align="center">
  <b>A modern, fast directory tree viewer built for Windows — with colors, icons, and rich export formats.</b><br/>
  One binary. No servers. No bloat.
</p>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go 1.21+](https://img.shields.io/badge/go-1.21+-blue.svg)](https://go.dev/dl/)
[![Windows](https://img.shields.io/badge/platform-Windows-0078D6.svg)](https://github.com/hosseinghorbani0/gowtree/releases)
[![GitHub Stars](https://img.shields.io/github/stars/hosseinghorbani0/gowtree?style=social)](https://github.com/hosseinghorbani0/gowtree)

---

## ⚡ Quick Start (30 seconds)

### Option A — Graphical Installer (Recommended)

1. Download from **[Releases v1.4.0](https://github.com/Hosseinghorbani0/gowtree/releases/tag/v1.4.0)**:
   - **`gowtree-setup-1.4.0.exe`** — Next/Next wizard installer (recommended)
   - **`gowtree.exe`** — portable, no install needed
   - **`gowtree-portable-1.4.0.zip`** — exe + install scripts
2. Run `gowtree-setup-1.4.0.exe` → check **Add to PATH** → Finish
3. Open a **new** terminal:

```powershell
gowtree
gowtree -a -s -L 2 --icons
```

### Option B — One-line script install

```powershell
git clone https://github.com/hosseinghorbani0/gowtree.git
cd gowtree
.\install.bat
```

### Option C — Build from source

```powershell
go build -o gowtree.exe ./cmd/gowtree
.\install.ps1
```

---

## 🧐 Why gowtree?

Windows ships with `tree`, but it feels stuck in 1995. **gowtree** gives you what modern developers expect:

| Feature | Windows `tree` | **gowtree** |
|---------|:--------------:|:-----------:|
| Colors & Unicode | ❌ | ✅ auto-detect |
| JSON output | ❌ | ✅ `-J` |
| Markdown / HTML | ❌ | ✅ `--markdown` / `--html` |
| Progress bar | ❌ | ✅ `-p` |
| Directory sizes | ❌ | ✅ `--du` |
| Regex filter | ❌ | ✅ `-R` |
| Nerd Font icons | ❌ | ✅ `--icons` |
| Clipboard copy | ❌ | ✅ `--clip` |
| YAML config | ❌ | ✅ `~/.gowtree.yaml` |
| Custom installer | ❌ | ✅ Inno Setup wizard |

---

## 📚 Examples

```powershell
# Basic tree
gowtree

# Hidden files, sizes, depth limit, icons
gowtree -a -s -L 2 --icons

# Only Go files, sorted by size (largest last)
gowtree -R "\.go$" --sort size -r --icons

# JSON for scripting
gowtree -J > tree.json

# Markdown for docs
gowtree --markdown --out TREE.md

# HTML report with folder sizes
gowtree --html --du --out report.html

# Copy to clipboard
gowtree --clip
```

---

## ⚙️ Flags

| Flag | Description |
|------|-------------|
| `-a` | Show hidden files (dot-files + Windows hidden attribute) |
| `-L <n>` | Max depth (0 = unlimited) |
| `-d` / `-f` | Directories only / files only |
| `-s` | Show file sizes |
| `-I <glob>` | Ignore pattern (repeatable) |
| `-R <regex>` | Filter by regular expression |
| `--sort` | `name`, `size`, `time`, or `none` |
| `-r` | Reverse sort |
| `--time` | Show time: `mod`, `change`, `create` |
| `--files-first` | List files before directories |
| `-l` | Follow symbolic links |
| `--icons` | Nerd Font icons |
| `-J` | JSON output |
| `--markdown` / `--html` | Export formats |
| `--du` | Cumulative directory sizes |
| `-p` | Progress bar |
| `--out <file>` | Write to file |
| `--clip` | Copy to clipboard |
| `--color` | `auto`, `always`, `never` |
| `--charset` | `utf8` or `ascii` |
| `--config` | Path to YAML config |
| `--version` / `--author` | Info |
| `--no-credit` | Hide signature line |

---

## 🏗️ Project Structure

```
gowtree/
├── cmd/gowtree/          # Entry point
├── internal/
│   ├── app/              # CLI orchestration
│   ├── config/           # YAML config loader
│   ├── format/           # Size formatting
│   ├── output/           # Markdown, HTML, clipboard
│   ├── platform/         # Windows-specific helpers
│   ├── style/            # Colors, icons, theme
│   └── tree/             # Core tree engine
├── installer/            # Inno Setup wizard (Windows)
├── scripts/              # Build scripts
└── docs/                 # Translated READMEs
```

---

## 🛠️ Build the Windows Installer

Requirements: [Go 1.21+](https://go.dev/dl/) and optionally [Inno Setup 6](https://jrsoftware.org/isinfo.php)

```powershell
.\scripts\build-installer.ps1
# Output: gowtree.exe + installer\output\gowtree-setup-1.4.0.exe
```

---

## 📋 Config File

Save as `%USERPROFILE%\.gowtree.yaml`:

```yaml
color: auto
icons: true
sort: name
max_depth: 3
ignore:
  - node_modules
  - .git
```

CLI flags always override config values.

---

## 🐛 Troubleshooting

### `gowtree` is not recognized

Open a **new** terminal after install. Or manually add `%USERPROFILE%\bin` to PATH:

```powershell
[Environment]::SetEnvironmentVariable('Path', $env:Path + ';%USERPROFILE%\bin', 'User')
```

### Colors not showing

```powershell
gowtree --color always
```

Use Windows Terminal for best Unicode/icon support.

### Installer won't build

Install [Inno Setup 6](https://jrsoftware.org/isinfo.php), or use `.\install.bat` instead.

---

## 📖 Our Story

> **~20%** of this project was shaped with AI assistance (architecture, docs, installer polish).  
> **~80%** was written by hand — during **internet blackouts in Iran**, amid the conflict between Iran, the United States, and Israel.

This tool exists because when you're offline, you still need to understand your codebase. A fast, beautiful `tree` on Windows shouldn't depend on the network — or on outdated built-ins.

If that story resonates with you, ⭐ star the repo. It helps more than you think.

---

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Ideas welcome: more icons, performance, translations.

---

## 🔗 Links

- **GitHub**: [hosseinghorbani0/gowtree](https://github.com/hosseinghorbani0/gowtree)
- **Releases**: [Download installer](https://github.com/hosseinghorbani0/gowtree/releases)
- **Issues**: [Report a bug](https://github.com/hosseinghorbani0/gowtree/issues)
- **Author**: [Hossein Ghorbani](https://github.com/hosseinghorbani0)

---

## 📄 License

MIT — see [LICENSE](LICENSE).

**Made with ❤️ by [Hossein Ghorbani](https://github.com/hosseinghorbani0)**
