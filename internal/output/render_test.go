package output

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/hosseinghorbani0/gowtree/internal/style"
	"github.com/hosseinghorbani0/gowtree/internal/tree"
)

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) {
	return 0, w.err
}

type shortWriter struct{}

func (shortWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	return len(p) - 1, io.ErrShortWrite
}

func testRenderNode() *tree.Node {
	return &tree.Node{
		Name: "root",
		Type: "directory",
		Children: []*tree.Node{
			{Name: "unsafe<&.txt", Type: "file"},
		},
	}
}

func TestRenderMarkdownPropagatesWriterError(t *testing.T) {
	errWant := errors.New("disk full")
	err := RenderMarkdown(failingWriter{err: errWant}, testRenderNode(), "", true, style.Theme{}, false)
	if !errors.Is(err, errWant) {
		t.Fatalf("RenderMarkdown() error = %v, want %v", err, errWant)
	}
}

func TestRenderHTMLPropagatesWriterError(t *testing.T) {
	errWant := errors.New("broken pipe")
	err := RenderHTML(failingWriter{err: errWant}, testRenderNode(), "", true, style.Theme{}, false)
	if !errors.Is(err, errWant) {
		t.Fatalf("RenderHTML() error = %v, want %v", err, errWant)
	}
}

func TestRenderMarkdownRejectsPartialWriter(t *testing.T) {
	err := RenderMarkdown(shortWriter{}, testRenderNode(), "", true, style.Theme{}, false)
	if !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("RenderMarkdown() error = %v, want %v", err, io.ErrShortWrite)
	}
}

func TestRenderHTMLEscapesNames(t *testing.T) {
	var out strings.Builder
	if err := RenderHTML(&out, testRenderNode(), "", true, style.Theme{}, false); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out.String(), "unsafe<&.txt") {
		t.Fatal("HTML output contains unescaped filename")
	}
	if !strings.Contains(out.String(), "unsafe&lt;&amp;.txt") {
		t.Fatalf("HTML output does not contain escaped filename: %s", out.String())
	}
}
