package parse

type Lexeme struct {
	Token Token
	Value []byte
	Pos   Position
}

// Position of lexeme
type Position struct {
	Line, Column int
}

// Lexical parser for servicebuilder
type Lexer struct {
	Name   String        // Input filename. For error messages
	Reader *bufio.Reader //  Content reader
	bus    chan Lexeme   // Lexemes bus is populated with lexemes as they are consumed
}

// Return new lexer
func NewLexer(name String, reader io.Reader) *Lexer {
	return &Lexer{
		Name:   name,
		Reader: bufio.NewReader(reader),
		bus:    make(chan Lexeme),
	}
}

// Read next rune
func (self *Lexer) read() rune {
	ch, _, err := self.Reader.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// Go to the previous rune
func (self *Lexer) unread() { _ = self.Reader.UnreadRune() }
