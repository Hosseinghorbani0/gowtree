# Deploying gowtree to GitHub

## Prerequisites

- [Git](https://git-scm.com/download/win)
- [Go 1.21+](https://go.dev/dl/)
- GitHub account

## 1. Configure Git (once)

```powershell
git config --global user.email "your-email@example.com"
git config --global user.name "Your Name"
```

## 2. Create the repository on GitHub

1. Go to [github.com/new](https://github.com/new)
2. Name: `gowtree`
3. Description: `🌳 A modern tree command for Windows`
4. Public repository
5. Do **not** initialize with README (this repo already has one)

## 3. Push from your machine

```powershell
cd d:\golang
git init
git add .
git commit -m "feat: gowtree v1.4.0 — modular Windows tree command"
git branch -M main
git remote add origin https://github.com/hosseinghorbani0/gowtree.git
git push -u origin main
```

When prompted for credentials, use a **Personal Access Token** (not your password).

## 4. Create a release

```powershell
git tag v1.4.0
git push origin v1.4.0
```

CI will build `gowtree.exe` and `gowtree-setup-1.4.0.exe` automatically.

Or manually:

```powershell
.\scripts\build-installer.ps1
# Upload gowtree.exe and installer\output\gowtree-setup-*.exe on GitHub Releases
```

## Security note

Never commit tokens or passwords to the repository. Use GitHub's credential manager or `gh auth login`.

If a token is exposed, revoke it immediately at:
**Settings → Developer settings → Personal access tokens**

## Troubleshooting

| Error | Fix |
|-------|-----|
| `command not found: gowtree` | Open a new terminal after install |
| `Permission denied (push)` | Check token has `repo` scope |
| Inno Setup not found | Install from [jrsoftware.org](https://jrsoftware.org/isinfo.php) |

---

Need help? [Open an issue](https://github.com/hosseinghorbani0/gowtree/issues)
