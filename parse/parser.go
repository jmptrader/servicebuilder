package parse

import (
	"errors"
	"fmt"
	"github.com/romanoff/servicebuilder/app"
	"io"
)

func NewParser(name string, reader io.Reader) *Parser {
	parser := &Parser{lexer: NewLexer(name, reader), buffer: make([]Lexeme, 0, 0), app: &app.Application{}}
	parser.scanner = parser.lexer.Scan()
	return parser
}

type Parser struct {
	lexer   *Lexer
	scanner chan Lexeme
	buffer  []Lexeme
	index   int
	app     *app.Application
}

func (self *Parser) scan() *Lexeme {
	if self.index == len(self.buffer) {
		lexeme := <-self.scanner
		self.buffer = append(self.buffer, lexeme)
		self.index++
		return &lexeme
	}
	lexeme := self.buffer[self.index]
	self.index++
	return &lexeme
}

func (self *Parser) scanIgnoreWhitespace() *Lexeme {
	lexeme := self.scan()
	if lexeme.Token == WS {
		lexeme = self.scan()
	}
	return lexeme
}

func (self *Parser) unscan() *Lexeme {
	if self.index == 0 {
		return nil
	}
	self.index--
	lexeme := self.buffer[self.index]
	return &lexeme
}
func (self *Parser) Parse() (app.Application, error) {
	for {
		lexeme := self.scanIgnoreWhitespace()
		if lexeme.Token == EOF {
			break
		}
		if lexeme.Token == MODEL {
			err := self.parseModel()
			if err != nil {
				return *self.app, err
			}
		}
	}
	return *self.app, nil
}

func (self *Parser) parseModel() error {
	model := &app.Model{}
	lexeme := self.scanIgnoreWhitespace()
	if lexeme.Token == EOF {
		return errors.New("Unexpected EOF")
	}
	if lexeme.Token != IDENT {
		return fmt.Errorf("found %q, expected model identifier", lexeme.Token)
	}
	model.Name = lexeme.Value
	lexeme = self.scanIgnoreWhitespace()
	if lexeme.Token == EOF {
		return errors.New("Unexpected EOF")
	}
	if lexeme.Token != LEFTBRACE {
		return fmt.Errorf("found %q, expected {", lexeme.Value)
	}
	// if err != self.parseModelFields(model) {
	// 	return err
	// }
	lexeme = self.scanIgnoreWhitespace()
	if lexeme.Token == EOF {
		return errors.New("Unexpected EOF")
	}
	if lexeme.Token != RIGHTBRACE {
		return fmt.Errorf("found %q, expected }", lexeme.Value)
	}
	return nil
}

func (self *Parser) parseModelFields(*app.Model) error {
	return nil
}
