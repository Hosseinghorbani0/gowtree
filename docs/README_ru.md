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
  <b>Современный просмотрщик дерева каталогов для Windows — цвета, иконки, экспорт JSON/Markdown/HTML.</b><br/>
  Один бинарник. Без серверов. Без лишней сложности.
</p>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](../LICENSE)
[![Windows](https://img.shields.io/badge/platform-Windows-0078D6.svg)](https://github.com/hosseinghorbani0/gowtree/releases)

---

## ⚡ Быстрый старт

1. Скачайте **`gowtree-setup-1.4.0.exe`** из [Releases](https://github.com/hosseinghorbani0/gowtree/releases)
2. Запустите мастер установки → включите **Add to PATH**
3. Откройте **новый** терминал:

```powershell
gowtree -a -s -L 2 --icons
```

Или из исходников:

```powershell
git clone https://github.com/hosseinghorbani0/gowtree.git
cd gowtree
.\install.bat
```

---

## 🧐 Зачем gowtree?

| Функция | Windows `tree` | **gowtree** |
|---------|:--------------:|:-----------:|
| Цвета и Unicode | ❌ | ✅ |
| JSON | ❌ | ✅ |
| Markdown / HTML | ❌ | ✅ |
| Прогресс-бар | ❌ | ✅ |
| Размер папок | ❌ | ✅ |
| Regex-фильтр | ❌ | ✅ |
| Иконки Nerd Font | ❌ | ✅ |
| Буфер обмена | ❌ | ✅ |
| YAML-конфиг | ❌ | ✅ |

---

## 📚 Примеры

```powershell
gowtree -J > tree.json
gowtree --markdown --out TREE.md
gowtree -R "\.go$" --sort size -r --icons
gowtree --clip
```

---

## 📖 Наша история

> **~20%** проекта создано с помощью AI (архитектура, документация, установщик).  
> **~80%** написано **вручную** — во время **отключений интернета в Иране**, на фоне конфликта между Ираном, США и Израилем.

Если эта история вам близка — поставьте ⭐ репозиторию.

---

## 🔗 Ссылки

[GitHub](https://github.com/hosseinghorbani0/gowtree) · [Releases](https://github.com/hosseinghorbani0/gowtree/releases) · [Issues](https://github.com/hosseinghorbani0/gowtree/issues)

**Автор: [Hossein Ghorbani](https://github.com/hosseinghorbani0)** · Лицензия MIT
