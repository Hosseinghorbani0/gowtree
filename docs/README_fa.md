# gowtree

<p align="center">
  🌍 <b>راهنما:</b>
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
  <b>یک دستور tree مدرن و سریع برای ویندوز — با رنگ، آیکون و خروجی‌های JSON/Markdown/HTML.</b><br/>
  یک فایل اجرایی. بدون سرور. بدون پیچیدگی اضافه.
</p>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](../LICENSE)
[![Go 1.21+](https://img.shields.io/badge/go-1.21+-blue.svg)](https://go.dev/dl/)
[![Windows](https://img.shields.io/badge/platform-Windows-0078D6.svg)](https://github.com/hosseinghorbani0/gowtree/releases)

---

## ⚡ شروع سریع (۳۰ ثانیه)

1. **`gowtree-setup-1.4.0.exe`** را از [Releases](https://github.com/Hosseinghorbani0/gowtree/releases/latest) دانلود کنید
2. نصب‌کننده را اجرا کنید → Next → Next → **Add to PATH** → Install → Finish
3. یک **ترمینال جدید** باز کنید:

```powershell
gowtree -a -s -L 2 --icons
```

> **هشدار SmartScreen ویندوز:** بار اول ممکن است صفحه آبی «Windows protected your PC» بیاید. روی **More info** بزنید → **Run anyway**. این برای نرم‌افزار امضا‌نشده جدید طبیعی است.

---

## 🧐 چرا gowtree؟

ویندوز دستور `tree` دارد، اما برای توسعه‌دهندهٔ امروزی کافی نیست:

| قابلیت | `tree` ویندوز | **gowtree** |
|--------|:-------------:|:-----------:|
| رنگ و Unicode | ❌ | ✅ |
| خروجی JSON | ❌ | ✅ |
| Markdown / HTML | ❌ | ✅ |
| نوار پیشرفت | ❌ | ✅ |
| حجم پوشه‌ها | ❌ | ✅ |
| فیلتر Regex | ❌ | ✅ |
| آیکون Nerd Font | ❌ | ✅ |
| کپی در کلیپ‌بورد | ❌ | ✅ |
| فایل تنظیمات YAML | ❌ | ✅ |
| نصب‌کننده اختصاصی | ❌ | ✅ |

---

## 📚 مثال‌ها

```powershell
gowtree                          # درخت ساده
gowtree -a -s -L 2 --icons       # مخفی‌ها، حجم، عمق، آیکون
gowtree -R "\.go$" --sort size -r --icons
gowtree -J > tree.json           # JSON
gowtree --markdown --out TREE.md # Markdown
gowtree --html --du --out report.html
gowtree --clip                   # کپی در کلیپ‌بورد
```

---

## 🏗️ ساختار پروژه

```
cmd/gowtree/       ← نقطه ورود
internal/app/      ← منطق CLI
internal/tree/     ← موتور اصلی درخت
installer/         ← نصب‌کننده Inno Setup
scripts/           ← اسکریپت build
```

---

## 🛠️ ساخت نصب‌کننده

```powershell
.\scripts\build-installer.ps1
```

نیاز: Go 1.21+ و (اختیاری) [Inno Setup 6](https://jrsoftware.org/isinfo.php)

---

## 📖 داستان این پروژه

> **حدود ۲۰٪** این پروژه با کمک هوش مصنوعی شکل گرفته (معماری، مستندات، نصب‌کننده).  
> **حدود ۸۰٪** آن **دست‌نویس** است — در **قطعی اینترنت ایران**، در میانه جنگ میان ایران، آمریکا و اسرائیل.

وقتی آفلاین هستید، هنوز باید ساختار پروژه‌تان را ببینید. gowtree برای همین ساخته شده.

اگر این داستان برایتان معنا دارد، ⭐ به مخزن ستاره بدهید.

---

## 🐛 رفع مشکل

### هشدار آبی SmartScreen

بار اول ویندوز پیام **Windows protected your PC** نشان می‌دهد چون برنامه هنوز با گواهی code-signing امضا نشده. این برای نرم‌افزار متن‌باز جدید طبیعی است.

**برای نصب:** **More info** → **Run anyway**

آیکون درخت، نام سازنده و اطلاعات نسخه اضافه شده. برای حذف کامل هشدار به گواهی Authenticode (~۱۰۰+ دلار/سال) نیاز است.

### دستور پیدا نمی‌شود؟ ترمینال جدید باز کنید یا `%USERPROFILE%\bin` را به PATH اضافه کنید.

**رنگ نمی‌بینید؟** `gowtree --color always` — Windows Terminal را امتحان کنید.

---

## 🔗 لینک‌ها

- [GitHub](https://github.com/hosseinghorbani0/gowtree) · [Releases](https://github.com/hosseinghorbani0/gowtree/releases) · [Issues](https://github.com/hosseinghorbani0/gowtree/issues)
- **سازنده:** [Hossein Ghorbani](https://github.com/hosseinghorbani0)

---

## 📄 مجوز

MIT — جزئیات در [LICENSE](../LICENSE).

**ساخته شده با ❤️ توسط [Hossein Ghorbani](https://github.com/hosseinghorbani0)**
