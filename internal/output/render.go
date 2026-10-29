package output

import (
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/hosseinghorbani0/gowtree/internal/format"
	"github.com/hosseinghorbani0/gowtree/internal/style"
	"github.com/hosseinghorbani0/gowtree/internal/tree"
)

func FormatSize(bytes int64) string { return format.Size(bytes) }

func RenderMarkdown(w io.Writer, node *tree.Node, prefix string, isRoot bool, theme style.Theme, showDU bool) error {
	if node.Type == "directory" {
		name := node.Name
		if theme.UseIcons {
			name = style.DirIcon() + " " + name
		}
		if isRoot {
			if _, err := fmt.Fprintf(w, "# %s\n", name); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(w, "%s- 📂 **%s**", prefix, name); err != nil {
				return err
			}
			if node.Size > 0 && showDU {
				if _, err := fmt.Fprintf(w, " (%s)", FormatSize(node.Size)); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		for _, child := range node.Children {
			newPrefix := prefix + "    "
			if isRoot {
				newPrefix = ""
			}
			if err := RenderMarkdown(w, child, newPrefix, false, theme, showDU); err != nil {
				return err
			}
		}
	} else if node.Type == "file" {
		name := node.Name
		if theme.UseIcons {
			ext := ""
			if dot := strings.LastIndex(node.Name, "."); dot >= 0 {
				ext = node.Name[dot:]
			}
			name = style.FileIcon(ext) + " " + name
		}
		if _, err := fmt.Fprintf(w, "%s- 📄 %s", prefix, name); err != nil {
			return err
		}
		if node.Size > 0 {
			if _, err := fmt.Fprintf(w, " (%s)", FormatSize(node.Size)); err != nil {
				return err
			}
		}
		if node.ModTime != "" {
			if _, err := fmt.Fprintf(w, " - %s", node.ModTime); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	return nil
}

func RenderHTML(w io.Writer, node *tree.Node, prefix string, isRoot bool, theme style.Theme, showDU bool) error {
	if node.Type == "directory" {
		name := html.EscapeString(node.Name)
		if theme.UseIcons {
			name = style.DirIcon() + " " + name
		}
		if isRoot {
			if _, err := fmt.Fprintf(w, "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<meta charset=\"utf-8\">\n<title>%s</title>\n<style>body{font-family:Segoe UI,system-ui,sans-serif;margin:2rem;line-height:1.6}ul{list-style:none;padding-left:1.2rem}li{margin:.2rem 0}</style>\n</head>\n<body>\n<h1>%s</h1>\n<ul>\n", name, name); err != nil {
				return err
			}
			for _, child := range node.Children {
				if err := RenderHTML(w, child, "  ", false, theme, showDU); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintf(w, "</ul>\n</body>\n</html>\n"); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(w, "%s<li><strong>%s</strong>", prefix, name); err != nil {
				return err
			}
			if node.Size > 0 && showDU {
				if _, err := fmt.Fprintf(w, " (%s)", FormatSize(node.Size)); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
			if len(node.Children) > 0 {
				if _, err := fmt.Fprintf(w, "%s<ul>\n", prefix+"  "); err != nil {
					return err
				}
				for _, child := range node.Children {
					if err := RenderHTML(w, child, prefix+"    ", false, theme, showDU); err != nil {
						return err
					}
				}
				if _, err := fmt.Fprintf(w, "%s</ul>\n", prefix+"  "); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintf(w, "%s</li>\n", prefix); err != nil {
				return err
			}
		}
	} else if node.Type == "file" {
		name := html.EscapeString(node.Name)
		if theme.UseIcons {
			ext := ""
			if dot := strings.LastIndex(node.Name, "."); dot >= 0 {
				ext = node.Name[dot:]
			}
			name = style.FileIcon(ext) + " " + name
		}
		if _, err := fmt.Fprintf(w, "%s<li>%s", prefix, name); err != nil {
			return err
		}
		if node.Size > 0 {
			if _, err := fmt.Fprintf(w, " (%s)", FormatSize(node.Size)); err != nil {
				return err
			}
		}
		if node.ModTime != "" {
			if _, err := fmt.Fprintf(w, " - %s", node.ModTime); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w, "</li>"); err != nil {
			return err
		}
	}
	return nil
}
