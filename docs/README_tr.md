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
  <b>Windows için modern bir dizin ağacı aracı — renkler, simgeler ve JSON/Markdown/HTML dışa aktarma.</b><br/>
  Tek binary. Sunucu yok. Gereksiz karmaşa yok.
</p>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](../LICENSE)
[![Windows](https://img.shields.io/badge/platform-Windows-0078D6.svg)](https://github.com/hosseinghorbani0/gowtree/releases)

---

## ⚡ Hızlı Başlangıç

1. [Releases](https://github.com/hosseinghorbani0/gowtree/releases) sayfasından **`gowtree-setup-1.4.0.exe`** indirin
2. Sihirbazı çalıştırın → **Add to PATH** seçeneğini işaretleyin
3. **Yeni** bir terminal açın:

```powershell
gowtree -a -s -L 2 --icons
```

Kaynak koddan kurulum:

```powershell
git clone https://github.com/hosseinghorbani0/gowtree.git
cd gowtree
.\install.bat
```

---

## 🧐 Neden gowtree?

| Özellik | Windows `tree` | **gowtree** |
|---------|:--------------:|:-----------:|
| Renk & Unicode | ❌ | ✅ |
| JSON çıktı | ❌ | ✅ |
| Markdown / HTML | ❌ | ✅ |
| İlerleme çubuğu | ❌ | ✅ |
| Klasör boyutları | ❌ | ✅ |
| Regex filtresi | ❌ | ✅ |
| Nerd Font simgeleri | ❌ | ✅ |
| Pano kopyalama | ❌ | ✅ |
| YAML yapılandırma | ❌ | ✅ |

---

## 📚 Örnekler

```powershell
gowtree -J > tree.json
gowtree --markdown --out TREE.md
gowtree -R "\.go$" --sort size -r --icons
gowtree --clip
```

---

## 📖 Hikayemiz

> Projenin **~%20'si** AI desteğiyle şekillendi (mimari, dokümantasyon, kurulum).  
> **~%80'i** elle yazıldı — **İran'daki internet kesintileri** sırasında, İran, ABD ve İsrail arasındaki çatışma döneminde.

Bu hikâye size dokunuyorsa ⭐ yıldız verin.

---

## 🔗 Bağlantılar

[GitHub](https://github.com/hosseinghorbani0/gowtree) · [Releases](https://github.com/hosseinghorbani0/gowtree/releases) · [Issues](https://github.com/hosseinghorbani0/gowtree/issues)

**Yazar: [Hossein Ghorbani](https://github.com/hosseinghorbani0)** · MIT Lisansı
