package parse

// Lexical parser for servicebuilder
type Lexer struct {
	Name  String      // Input filename. For error messages
	Input []byte      // Input file content
	bus   chan Lexeme // Lexemes bus is populated with lexemes as they are consumed
}

type Lexeme struct {
	Token Token
	Value []byte
	Pos   Position
}

type Position struct {
	Line, Column int
}

// Lexer token
type Token int

const (
	ILLEGAL Token = iota
	EOF           // end of file
	WS            // whitespace

	IDENT        // identifier
	LEFTBRACE    // {
	RIGHTBRACE   // }
	LEFTSQBRACE  // ]
	RIGHTSQBRACE // ]
	COLON        // :
	COMMA        // ,
	NUMERIC      // 1234

	// modes
	FIELDS
	PAGINATION
)
