package markdown

type Inline interface {
	inline()
}

type TextNode struct {
	Value string
}

func (TextNode) inline() {}

type Newline struct{}

func (Newline) inline() {}

type ItalicNode struct {
	Children []Inline
}

func (ItalicNode) inline() {}

type BoldNode struct {
	Children []Inline
}

func (BoldNode) inline() {}

type CodeNode struct {
	Value string
}

func (CodeNode) inline() {}

type LinkNode struct {
	Text []Inline
	URL  string
}

func (LinkNode) inline() {}

type ImageNode struct {
	AltText []Inline
	URL     string
}

func (ImageNode) inline() {}
