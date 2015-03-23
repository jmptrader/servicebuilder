package servicebuilder

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

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
	FIELDS
	PAGINATION
)

type ModeType int

const (
	M_MODEL ModeType = iota
	M_FIELDS
	M_PAGINATION
)

var eof = rune(0)

// Checks if rune is whitespace
func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

// Checks if rune is letter
func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// Checks if rune is a digit
func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// Scanner represents a lexical scanner.
type Scanner struct {
	r     *bufio.Reader
	modes []*Mode
}

type Mode struct {
	mode         ModeType
	braceCounter int
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	} else if isDigit(ch) {
		s.unread()
		return s.scanNumber()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case ':':
		return COLON, string(ch)
	case ',':
		return COMMA, string(ch)
	case '{':
		if len(s.modes) == 0 {
			s.modes = append(s.modes, &Mode{mode: M_MODEL})
		}
		s.modes[0].braceCounter++
		return LEFTBRACE, string(ch)
	case '}':
		if len(s.modes) == 0 {
			s.modes = append(s.modes, &Mode{mode: M_MODEL})
		}
		s.modes[0].braceCounter--
		if s.modes[0].braceCounter == 0 {
			copy(s.modes[0:], s.modes[1:])
			s.modes[len(s.modes)-1] = nil
			s.modes = s.modes[:len(s.modes)-1]
		}
		return RIGHTBRACE, string(ch)
	case '[':
		return LEFTSQBRACE, string(ch)
	case ']':
		return RIGHTSQBRACE, string(ch)
	}

	return ILLEGAL, string(ch)
}

// scanNumber consumes the current rune and all contiguous number runes
func (s *Scanner) scanNumber() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isDigit(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}
	return NUMERIC, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	if len(s.modes) > 0 && s.modes[0].mode == M_MODEL {
		switch strings.ToUpper(buf.String()) {
		case "FIELDS":
			s.modes = append(s.modes, nil)
			copy(s.modes[1:], s.modes[0:])
			s.modes[0] = &Mode{mode: M_FIELDS}
			return FIELDS, buf.String()
		case "PAGINATION":
			s.modes = append(s.modes, nil)
			copy(s.modes[1:], s.modes[0:])
			s.modes[0] = &Mode{mode: M_PAGINATION}
			return PAGINATION, buf.String()
		}
	}

	// Otherwise return as a regular identifier.
	return IDENT, buf.String()
}
