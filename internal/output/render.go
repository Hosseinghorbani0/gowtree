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

func RenderMarkdown(w io.Writer, node *tree.Node, prefix string, isRoot bool, theme style.Theme, showDU bool) {
	if node.Type == "directory" {
		name := node.Name
		if theme.UseIcons {
			name = style.DirIcon() + " " + name
		}
		if isRoot {
			fmt.Fprintf(w, "# %s\n", name)
		} else {
			fmt.Fprintf(w, "%s- 📂 **%s**", prefix, name)
			if node.Size > 0 && showDU {
				fmt.Fprintf(w, " (%s)", FormatSize(node.Size))
			}
			fmt.Fprintln(w)
		}
		for _, child := range node.Children {
			newPrefix := prefix + "    "
			if isRoot {
				newPrefix = ""
			}
			RenderMarkdown(w, child, newPrefix, false, theme, showDU)
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
		fmt.Fprintf(w, "%s- 📄 %s", prefix, name)
		if node.Size > 0 {
			fmt.Fprintf(w, " (%s)", FormatSize(node.Size))
		}
		if node.ModTime != "" {
			fmt.Fprintf(w, " - %s", node.ModTime)
		}
		fmt.Fprintln(w)
	}
}

func RenderHTML(w io.Writer, node *tree.Node, prefix string, isRoot bool, theme style.Theme, showDU bool) {
	if node.Type == "directory" {
		name := html.EscapeString(node.Name)
		if theme.UseIcons {
			name = style.DirIcon() + " " + name
		}
		if isRoot {
			fmt.Fprintf(w, "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<meta charset=\"utf-8\">\n<title>%s</title>\n<style>body{font-family:Segoe UI,system-ui,sans-serif;margin:2rem;line-height:1.6}ul{list-style:none;padding-left:1.2rem}li{margin:.2rem 0}</style>\n</head>\n<body>\n<h1>%s</h1>\n<ul>\n", name, name)
			for _, child := range node.Children {
				RenderHTML(w, child, "  ", false, theme, showDU)
			}
			fmt.Fprintf(w, "</ul>\n</body>\n</html>\n")
		} else {
			fmt.Fprintf(w, "%s<li><strong>%s</strong>", prefix, name)
			if node.Size > 0 && showDU {
				fmt.Fprintf(w, " (%s)", FormatSize(node.Size))
			}
			fmt.Fprintln(w)
			if len(node.Children) > 0 {
				fmt.Fprintf(w, "%s<ul>\n", prefix+"  ")
				for _, child := range node.Children {
					RenderHTML(w, child, prefix+"    ", false, theme, showDU)
				}
				fmt.Fprintf(w, "%s</ul>\n", prefix+"  ")
			}
			fmt.Fprintf(w, "%s</li>\n", prefix)
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
		fmt.Fprintf(w, "%s<li>%s", prefix, name)
		if node.Size > 0 {
			fmt.Fprintf(w, " (%s)", FormatSize(node.Size))
		}
		if node.ModTime != "" {
			fmt.Fprintf(w, " - %s", node.ModTime)
		}
		fmt.Fprintln(w, "</li>")
	}
}
