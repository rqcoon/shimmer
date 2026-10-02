package markdown

import "strings"

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() Document {
	var blocks []Block

	for p.pos < len(p.tokens) {
		block := p.parseBlock()
		if block != nil {
			blocks = append(blocks, block)
		}
	}

	return Document{
		Children: blocks,
	}
}

func (p *Parser) parseBlock() Block {
	switch p.tokens[p.pos].Type {
	case TokenHeadingStart:
		return p.parseHeading()
	case TokenCodeBlockStart:
		return p.parseCodeBlock()
	default:
		return p.parseParagraph()
	}
}

func (p *Parser) parseHeading() Block {
	p.pos++ // consume starting token

	var children []Inline
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type != TokenHeadingEnd {
		inline := p.parseInline()
		if inline != nil {
			children = append(children, inline)
		}
	}

	if p.pos < len(p.tokens) {
		p.pos++
	}

	return Heading{
		level:    1, // todo: heading level lexing
		Children: children,
	}
}

func (p *Parser) parseCodeBlock() Block {
	lang := p.tokens[p.pos].Value // Language is stored in Value of the opening token

	if len(lang) == 0 {
		// lang = DEFAULT // todo: infer default from frontmatter
		lang = "Go"
	}

	p.pos++

	var children strings.Builder
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type != TokenCodeBlockEnd {
		children.WriteString(p.tokens[p.pos].Value)
		p.pos++
	}

	if p.pos < len(p.tokens) {
		p.pos++
	}

	return CodeBlock{
		Language: lang,
		Content:  children.String(),
	}
}

func (p *Parser) parseParagraph() Block {
	node := Paragraph{}

	for p.pos < len(p.tokens) {
		switch p.tokens[p.pos].Type {
		case TokenHeadingStart, TokenCodeBlockStart:
			if len(node.Children) > 0 {
				return node
			}

			return nil

		case TokenBlankline:
			p.pos++
			return node
		}

		inline := p.parseInline()

		if inline != nil {
			node.Children = append(node.Children, inline)
		}
	}

	if len(node.Children) == 0 {
		return nil
	}

	return node
}

func (p *Parser) parseInline() Inline {
	switch p.tokens[p.pos].Type {
	case TokenText:
		value := p.tokens[p.pos].Value
		p.pos++

		return TextNode{Value: value}

	case TokenNewline:
		p.pos++
		return Newline{}

	case TokenBoldStart:
		return p.parseBold()

	case TokenItalicStart:
		return p.parseItalic()

	case TokenCodeStart:
		return p.parseCode()

	case TokenLinkStart:
		return p.parseLink()

	default:
		p.pos++
		return nil
	}
}

func (p *Parser) parseBold() Inline {
	p.pos++ // consume starting token
	var children []Inline
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type != TokenBoldEnd {
		inline := p.parseInline()
		if inline != nil {
			children = append(children, inline)
		}
	}

	if p.pos < len(p.tokens) {
		p.pos++
	}

	return BoldNode{Children: children}
}

func (p *Parser) parseItalic() Inline {
	p.pos++
	var children []Inline
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type != TokenItalicEnd {
		inline := p.parseInline()
		if inline != nil {
			children = append(children, inline)
		}
	}

	if p.pos < len(p.tokens) {
		p.pos++
	}

	return ItalicNode{Children: children}
}

func (p *Parser) parseCode() Inline {
	p.pos++ // consume TokenCodeStart

	var children strings.Builder

	for p.pos < len(p.tokens) && p.tokens[p.pos].Type != TokenCodeEnd {
		children.WriteString(p.tokens[p.pos].Value)
		p.pos++
	}

	if p.pos < len(p.tokens) {
		p.pos++
	}

	return CodeNode{
		Value: children.String(),
	}
}

func (p *Parser) parseLink() Inline {
	p.pos++

	var children []Inline

	if p.pos < len(p.tokens) && p.tokens[p.pos].Type == TokenLinkText {
		children = append(children, TextNode{Value: p.tokens[p.pos].Value})
		p.pos++
	}

	var url string

	if p.pos < len(p.tokens) && p.tokens[p.pos].Type == TokenLinkURL {
		url = p.tokens[p.pos].Value
		p.pos++
	}

	if p.pos < len(p.tokens) && p.tokens[p.pos].Type == TokenLinkEnd {
		p.pos++
	}

	return LinkNode{
		URL:  url,
		Text: children,
	}
}
