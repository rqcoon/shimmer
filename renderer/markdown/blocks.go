package markdown

type Block interface {
	block()
}

type Document struct {
	Children []Block
}

func (Document) block() {}

type Heading struct {
	level    int
	Children []Inline
}

func (Heading) block() {}

type Paragraph struct {
	Children []Inline
}

func (Paragraph) block() {}

type CodeBlock struct {
	Language string
	Content  string
}

func (CodeBlock) block() {}
