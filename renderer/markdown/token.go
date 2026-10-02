package markdown

type TokenType int

const (
	TokenText TokenType = iota
	TokenHeadingStart
	TokenHeadingEnd

	TokenNewline
	TokenBlankline

	TokenBoldStart
	TokenBoldEnd
	TokenItalicStart
	TokenItalicEnd

	TokenCodeStart
	TokenCodeEnd

	TokenCodeBlockStart
	TokenCodeBlockEnd

	TokenLinkStart
	TokenLinkText
	TokenLinkEnd
	TokenLinkURL
)

type Token struct {
	Type  TokenType
	Value string
}
