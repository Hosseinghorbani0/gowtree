# gowtree

<p align="center">
  🌍 <b>الreadme:</b>
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
  <b>أداة عرض شجرة المجلدات الحديثة لنظام Windows — ألوان، أيقونات، وتصدير JSON/Markdown/HTML.</b><br/>
  ملف تنفيذي واحد. بدون خوادم. بدون تعقيد زائد.
</p>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](../LICENSE)
[![Windows](https://img.shields.io/badge/platform-Windows-0078D6.svg)](https://github.com/hosseinghorbani0/gowtree/releases)

---

## ⚡ البدء السريع

1. حمّل **`gowtree-setup-1.4.0.exe`** من [Releases](https://github.com/hosseinghorbani0/gowtree/releases)
2. شغّل معالج التثبيت → فعّل **Add to PATH**
3. افتح **طرفية جديدة**:

```powershell
gowtree -a -s -L 2 --icons
```

أو من المصدر:

```powershell
git clone https://github.com/hosseinghorbani0/gowtree.git
cd gowtree
.\install.bat
```

---

## 🧐 لماذا gowtree؟

| الميزة | `tree` في Windows | **gowtree** |
|--------|:-----------------:|:-----------:|
| الألوان و Unicode | ❌ | ✅ |
| إخراج JSON | ❌ | ✅ |
| Markdown / HTML | ❌ | ✅ |
| شريط التقدم | ❌ | ✅ |
| أحجام المجلدات | ❌ | ✅ |
| تصفية Regex | ❌ | ✅ |
| أيقونات Nerd Font | ❌ | ✅ |
| الحافظة | ❌ | ✅ |
| إعدادات YAML | ❌ | ✅ |

---

## 📚 أمثلة

```powershell
gowtree -J > tree.json
gowtree --markdown --out TREE.md
gowtree -R "\.go$" --sort size -r --icons
gowtree --clip
```

---

## 📖 قصتنا

> **~20%** من المشروع بمساعدة الذكاء الاصطناعي (البنية، التوثيق، المثبّت).  
> **~80%** كُتب **يدوياً** — أثناء **انقطاع الإنترنت في إيران**، في ظل الصراع بين إيران وأمريكا وإسرائيل.

إن ألهمتك هذه القصة، ⭐ ضع نجمة على المستودع.

---

## 🔗 روابط

[GitHub](https://github.com/hosseinghorbani0/gowtree) · [Releases](https://github.com/hosseinghorbani0/gowtree/releases) · [Issues](https://github.com/hosseinghorbani0/gowtree/issues)

**المؤلف: [Hossein Ghorbani](https://github.com/hosseinghorbani0)** · رخصة MIT
