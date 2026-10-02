package markdown

func isSpecial(r rune) bool {
	switch r {
	case '#', '*', '`', '[', '\n':
		return true
	default:
		return false
	}
}

type TokenBuffer struct {
	items []Token
}

func (q *TokenBuffer) Append(item Token) {
	q.items = append(q.items, item)
}

func (q *TokenBuffer) Clear() {
	q.items = q.items[:0]
}

func (q *TokenBuffer) Commit() []Token {
	tokens := q.items
	q.items = nil

	return tokens
}

type Lexer struct {
	source []rune
	pos    int

	buffer TokenBuffer
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		source: []rune(source),
	}
}

func (l *Lexer) Tokenize() []Token {
	var tokens []Token

	for l.pos < len(l.source) {
		token := l.nextTokens()

		if token != nil {
			tokens = append(tokens, token...)
		}
	}

	return tokens
}

func (l *Lexer) nextTokens() []Token {
	switch l.source[l.pos] {
	case '#':
		return l.lexHeading()
	case '*':
		return l.lexStar()
	case '`':
		return l.lexCode()
	case '[':
		return l.lexLink()
	case '\n':
		return l.lexNewline()
	default:
		return l.lexText()
	}
}

func (l *Lexer) lexHeading() []Token {
	startPos := l.pos
	l.pos++ // consume '#'

	l.buffer.Append(Token{
		Type:  TokenHeadingStart,
		Value: "#",
	})

	content, found := l.lexUntil('\n')
	if !found {
		return l.rejectSpeculation(startPos)
	}

	l.buffer.Append(Token{
		Type:  TokenText,
		Value: content,
	})

	l.buffer.Append(Token{
		Type: TokenHeadingEnd,
	})

	return l.buffer.Commit()
}

func (l *Lexer) lexStar() []Token {
	startpos := l.pos

	if l.pos+1 < len(l.source) && l.source[l.pos+1] == '*' {
		l.pos += 2 // consume '**'

		l.buffer.Append(Token{
			Type:  TokenBoldStart,
			Value: "**",
		})

		p, found := l.lexUntil('*')
		if !found ||
			l.pos+1 >= len(l.source) ||
			l.source[l.pos+1] != '*' {
			return l.rejectSpeculation(startpos)
		}

		l.buffer.Append(Token{
			Type:  TokenText,
			Value: p,
		})

		l.pos += 2 // consume '**'

		l.buffer.Append(Token{
			Type:  TokenBoldEnd,
			Value: "**",
		})

		return l.buffer.Commit()
	}

	l.pos++ // consume '*'

	l.buffer.Append(Token{
		Type:  TokenItalicStart,
		Value: "*",
	})

	p, found := l.lexUntil('*')
	if !found {
		return l.rejectSpeculation(startpos)
	}

	l.buffer.Append(Token{
		Type:  TokenText,
		Value: p,
	})

	l.pos++ // consume '*'

	l.buffer.Append(Token{
		Type:  TokenItalicEnd,
		Value: "*",
	})

	return l.buffer.Commit()
}

func (l *Lexer) lexCode() []Token {
	startPos := l.pos
	l.pos++ // consume '`'

	l.buffer.Append(Token{
		Type:  TokenCodeStart,
		Value: "`",
	})

	p, found := l.lexUntil('`')
	if !found {
		return l.rejectSpeculation(startPos)
	}

	l.buffer.Append(Token{
		Type:  TokenText,
		Value: p,
	})

	l.pos++ // consume '`'

	l.buffer.Append(Token{
		Type:  TokenCodeEnd,
		Value: "`",
	})

	return l.buffer.Commit()
}

func (l *Lexer) lexLink() []Token {
	startPos := l.pos
	l.pos++ // consume '['

	// enqueue the start of the link token
	l.buffer.Append(Token{
		Type:  TokenLinkStart,
		Value: "[",
	})

	display, found := l.lexUntil(']')
	if !found {
		return l.rejectSpeculation(startPos)
	}

	l.buffer.Append(Token{
		Type:  TokenLinkText,
		Value: display,
	})

	l.pos++ // consume ']'

	// expect '('
	if l.pos >= len(l.source) || l.source[l.pos] != '(' {
		return l.rejectSpeculation(startPos)
	}

	l.pos++ // consume '('

	url, found := l.lexUntil(')')
	if !found {
		return l.rejectSpeculation(startPos)
	}

	l.buffer.Append(Token{
		Type:  TokenLinkURL,
		Value: url,
	})

	l.pos++ // consume ')'

	l.buffer.Append(Token{
		Type: TokenLinkEnd,
	})

	return l.buffer.Commit()
}

func (l *Lexer) lexNewline() []Token {
	start := l.pos

	l.pos++

	if l.pos < len(l.source) && l.source[l.pos] == '\n' {
		for l.pos < len(l.source) && l.source[l.pos] == '\n' {
			l.pos++
		}

		return []Token{
			{
				Type:  TokenBlankline,
				Value: string(l.source[start:l.pos]),
			},
		}
	}

	return []Token{
		{
			Type:  TokenNewline,
			Value: "\n",
		},
	}
}

func (l *Lexer) lexText() []Token {
	start := l.pos
	for l.pos < len(l.source) && !isSpecial(l.source[l.pos]) {
		l.pos++
	}

	return []Token{
		{
			Type:  TokenText,
			Value: string(l.source[start:l.pos]),
		},
	}
}

func (l *Lexer) lexUntil(delimiter rune) (string, bool) {
	start := l.pos
	for l.pos < len(l.source) && l.source[l.pos] != delimiter {
		l.pos++
	}

	if l.pos == len(l.source) {
		return string(l.source[start:l.pos]), false
	}

	return string(l.source[start:l.pos]), true
}

func (l *Lexer) rejectSpeculation(startPos int) []Token {
	l.buffer = TokenBuffer{}
	l.pos = startPos

	char := l.source[l.pos]
	l.pos++

	return []Token{
		{
			Type:  TokenText,
			Value: string(char),
		},
	}
}
