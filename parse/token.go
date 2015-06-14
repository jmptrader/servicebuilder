package parse

// Lexer token
//go:generate stringer -type=Token
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
	MODEL
	FIELDS
	PAGINATION
	ACTIONS

	// Types
	STRING
	INT
	DOUBLE
	DATE
	DATETIME
)
