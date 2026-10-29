#Requires -Version 5.1
# Generates icon.ico and Inno Setup wizard images for gowtree
$ErrorActionPreference = "Stop"
Add-Type -AssemblyName System.Drawing

$Root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$Assets = Join-Path $Root "assets"
New-Item -ItemType Directory -Force -Path $Assets | Out-Null

function New-TreeBitmap([int]$Size) {
    $bmp = New-Object System.Drawing.Bitmap $Size, $Size
    $g = [System.Drawing.Graphics]::FromImage($bmp)
    $g.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::AntiAlias
    $g.Clear([System.Drawing.Color]::FromArgb(15, 23, 42))

    $bg = New-Object System.Drawing.Drawing2D.LinearGradientBrush (
        [System.Drawing.Rectangle]::new(0, 0, $Size, $Size),
        [System.Drawing.Color]::FromArgb(20, 83, 45),
        [System.Drawing.Color]::FromArgb(15, 23, 42),
        45
    )
    $g.FillRectangle($bg, 0, 0, $Size, $Size)
    $bg.Dispose()

    $trunkW = [Math]::Max(4, [int]($Size * 0.08))
    $trunkH = [int]($Size * 0.32)
    $cx = [int]($Size / 2)
    $trunk = New-Object System.Drawing.SolidBrush ([System.Drawing.Color]::FromArgb(133, 77, 14))
    $g.FillRectangle($trunk, $cx - $trunkW / 2, $Size - $trunkH - [int]($Size * 0.08), $trunkW, $trunkH)

    $leaf = New-Object System.Drawing.SolidBrush ([System.Drawing.Color]::FromArgb(74, 222, 128))
    $r = [int]($Size * 0.28)
    $g.FillEllipse($leaf, $cx - $r, [int]($Size * 0.18), $r * 2, $r * 2)
    $g.FillEllipse($leaf, $cx - $r - [int]($Size * 0.18), [int]($Size * 0.28), [int]($r * 1.2), [int]($r * 1.2))
    $g.FillEllipse($leaf, $cx + [int]($Size * 0.06), [int]($Size * 0.28), [int]($r * 1.2), [int]($r * 1.2))

    $g.Dispose()
    return $bmp
}

function Save-IconFile($icoPath) {
    $tree256 = New-TreeBitmap 256
    $hIcon = $tree256.GetHicon()
    $icon = [System.Drawing.Icon]::FromHandle($hIcon)
    $fs = [System.IO.File]::Create($icoPath)
    $icon.Save($fs)
    $fs.Close()
    $tree256.Dispose()
}

function Save-WizardBmp($outPath, [int]$W, [int]$H) {
    $bmp = New-Object System.Drawing.Bitmap $W, $H
    $g = [System.Drawing.Graphics]::FromImage($bmp)
    $g.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::AntiAlias
    $g.Clear([System.Drawing.Color]::FromArgb(15, 23, 42))
    $brush = New-Object System.Drawing.Drawing2D.LinearGradientBrush (
        [System.Drawing.Point]::new(0, 0), [System.Drawing.Point]::new($W, $H),
        [System.Drawing.Color]::FromArgb(20, 83, 45), [System.Drawing.Color]::FromArgb(15, 23, 42)
    )
    $g.FillRectangle($brush, 0, 0, $W, $H); $brush.Dispose()

    $iconSize = [Math]::Min($W - 16, [int]($H * 0.42))
    $tree = New-TreeBitmap $iconSize
    $x = [int](($W - $iconSize) / 2)
    $y = [int](($H - $iconSize) / 2) - 12
    $g.DrawImage($tree, $x, $y, $iconSize, $iconSize); $tree.Dispose()

    $sf = New-Object System.Drawing.StringFormat
    $sf.Alignment = [System.Drawing.StringAlignment]::Center
    $font = New-Object System.Drawing.Font("Segoe UI", [Math]::Max(10, $W / 14), [System.Drawing.FontStyle]::Bold)
    $g.DrawString("gowtree", $font, [System.Drawing.Brushes]::White,
        [System.Drawing.RectangleF]::new(0, $y + $iconSize + 6, $W, 36), $sf)
    $font2 = New-Object System.Drawing.Font("Segoe UI", [Math]::Max(7, $W / 22))
    $g.DrawString("Directory Tree", $font2, [System.Drawing.Brushes]::LightGray,
        [System.Drawing.RectangleF]::new(0, $y + $iconSize + 28, $W, 24), $sf)
    $g.Dispose()
    $bmp.Save($outPath, [System.Drawing.Imaging.ImageFormat]::Bmp)
    $bmp.Dispose()
}

Write-Host "Generating gowtree assets..." -ForegroundColor Cyan
Save-IconFile (Join-Path $Assets "icon.ico")
Save-WizardBmp (Join-Path $Assets "wizard-large.bmp") 164 314
Save-WizardBmp (Join-Path $Assets "wizard-small.bmp") 55 58
Write-Host "  assets/icon.ico" -ForegroundColor Green
Write-Host "  assets/wizard-large.bmp" -ForegroundColor Green
Write-Host "  assets/wizard-small.bmp" -ForegroundColor Green
