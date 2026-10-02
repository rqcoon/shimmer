package markdown

import (
	"strings"
)

type MDRenderer struct{}

func NewMDRenderer() *MDRenderer {
	return &MDRenderer{}
}

func (r *MDRenderer) Render(source string) (string, error) {
	lexer := NewLexer(source)
	tokens := lexer.Tokenize()

	parser := NewParser(tokens)
	document := parser.Parse()

	html := r.renderDocument(document)

	return html, nil
}

func (r *MDRenderer) renderDocument(doc Document) string {
	var output strings.Builder

	for _, block := range doc.Children {
		r.renderBlock(&output, block)
	}

	return output.String()
}

func (r *MDRenderer) renderBlock(output *strings.Builder, block Block) {
	switch b := block.(type) {
	case Heading:
		output.WriteString("<h1>")
		for _, inline := range b.Children {
			r.renderInline(output, inline)
		}
		output.WriteString("</h1>\n")
	case Paragraph:
		output.WriteString("<p>")
		for _, inline := range b.Children {
			r.renderInline(output, inline)
		}
		output.WriteString("</p>\n")
	case CodeBlock:
		output.WriteString("<pre><code>")
		output.WriteString(b.Content)
		output.WriteString("</code></pre>\n")
	default:
		output.WriteString("<!-- md.go: Unknown block type -->\n")
	}
}

func (r *MDRenderer) renderInline(output *strings.Builder, inline Inline) {
	switch node := inline.(type) {
	case TextNode:
		output.WriteString(node.Value)
	case LinkNode:
		output.WriteString("<a href=\"")
		output.WriteString(node.URL)
		output.WriteString("\">")
		for _, child := range node.Text {
			r.renderInline(output, child)
		}
		output.WriteString("</a>")
	case ItalicNode:
		output.WriteString("<em>")
		for _, child := range node.Children {
			r.renderInline(output, child)
		}
		output.WriteString("</em>")
	case BoldNode:
		output.WriteString("<strong>")
		for _, child := range node.Children {
			r.renderInline(output, child)
		}
		output.WriteString("</strong>")
	case CodeNode:
		output.WriteString("<code>")
		output.WriteString(node.Value)
		output.WriteString("</code>")
	default:
		output.WriteString("<!-- md.go: Unknown inline type -->")
	}
}
