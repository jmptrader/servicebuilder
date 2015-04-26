package parse

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
