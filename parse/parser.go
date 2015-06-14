package parse

import (
	"errors"
	"fmt"
	"github.com/romanoff/servicebuilder/app"
	"io"
	"strconv"
)

func NewParser(name string, reader io.Reader) *Parser {
	parser := &Parser{lexer: NewLexer(name, reader), app: &app.Application{Models: make([]*app.Model, 0, 0)}}
	parser.scanner = parser.lexer.Scan()
	return parser
}

type Parser struct {
	lexer         *Lexer
	scanner       chan Lexeme
	lexeme        *Lexeme
	unscanEnabled bool
	index         int
	app           *app.Application
}

func (self *Parser) scan() *Lexeme {
	if !self.unscanEnabled {
		lexeme := <-self.scanner
		self.lexeme = &lexeme
	}
	self.unscanEnabled = false
	return self.lexeme
}

func (self *Parser) scanIgnoreWhitespace() *Lexeme {
	lexeme := self.scan()
	if lexeme.Token == WS {
		lexeme = self.scan()
	}
	return lexeme
}

func (self *Parser) unscan() {
	self.unscanEnabled = true
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
	model := &app.Model{Fields: make([]*app.Field, 0, 0)}
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
	var err error

modelSections:
	for {
		lexeme = self.scanIgnoreWhitespace()
		if lexeme.Token == EOF {
			return errors.New("Unexpected EOF")
		}
		switch lexeme.Token {
		case FIELDS:
			err = self.parseModelFields(model)
		case PAGINATION:
			err = self.parseModelPagination(model)
		case ACTIONS:
			err = self.parseModelActions(model)
		case RIGHTBRACE:
			self.unscan()
			break modelSections
		default:
			return errors.New("Unexpected EOF")
		}
		if err != nil {
			return err
		}
	}

	lexeme = self.scanIgnoreWhitespace()
	if lexeme.Token == EOF {
		return errors.New("Unexpected EOF")
	}
	if lexeme.Token != RIGHTBRACE {
		return fmt.Errorf("found %q, expected }", lexeme.Value)
	}
	self.app.Models = append(self.app.Models, model)
	return nil
}

func (self *Parser) parseModelFields(model *app.Model) error {
	lexeme := self.scanIgnoreWhitespace()
	if lexeme.Token == EOF {
		return errors.New("Unexpected EOF")
	}
	if lexeme.Token != LEFTBRACE {
		return fmt.Errorf("found %q, expected {", lexeme.Value)
	}

	for {
		if lexeme = self.scanIgnoreWhitespace(); lexeme.Token != IDENT {
			if lexeme.Token == RIGHTBRACE {
				self.unscan()
				break
			}
			return fmt.Errorf("found %q, expected field name", lexeme.Value)
		}
		field := &app.Field{}
		field.Name = lexeme.Value

		if lexeme = self.scanIgnoreWhitespace(); lexeme.Token != COLON {
			return fmt.Errorf("found %q, expected :", lexeme.Value)
		}

		lexeme = self.scanIgnoreWhitespace()
		switch lexeme.Token {
		case STRING:
			field.Type = app.STRING
		case INT:
			field.Type = app.INT
		case DOUBLE:
			field.Type = app.DOUBLE
		case DATE:
			field.Type = app.DATE
		case DATETIME:
			field.Type = app.DATETIME
		default:
			return fmt.Errorf("found %q, expected field type", lexeme.Value)
		}
		model.Fields = append(model.Fields, field)
	}

	lexeme = self.scanIgnoreWhitespace()
	if lexeme.Token == EOF {
		return errors.New("Unexpected EOF")
	}
	if lexeme.Token != RIGHTBRACE {
		return fmt.Errorf("found %q, expected }", lexeme.Value)
	}
	return nil
}

func (self *Parser) parseModelPagination(model *app.Model) error {
	lexeme := self.scanIgnoreWhitespace()
	if lexeme.Token == EOF {
		return errors.New("Unexpected EOF")
	}
	if lexeme.Token != LEFTBRACE {
		return fmt.Errorf("found %q, expected {", lexeme.Value)
	}

	for {
		lexeme = self.scanIgnoreWhitespace()
		if lexeme.Token == RIGHTBRACE {
			self.unscan()
			break
		}
		if lexeme.Token != IDENT || (lexeme.Value != "per_page" && lexeme.Value != "max_per_page") {
			return fmt.Errorf("expected per_page or max_per_page, but got %q", lexeme.Value)
		}
		attribute := lexeme.Value
		if lexeme = self.scanIgnoreWhitespace(); lexeme.Token != COLON {
			return fmt.Errorf("found %q, expected :", lexeme.Value)
		}
		if lexeme = self.scanIgnoreWhitespace(); lexeme.Token != NUMERIC {
			return fmt.Errorf("found %q, expected numeric value", lexeme.Value)
		}
		paginationValue, err := strconv.Atoi(lexeme.Value)
		if err != nil {
			return err
		}
		if model.Pagination == nil {
			model.Pagination = &app.Pagination{}
		}
		switch attribute {
		case "per_page":
			model.Pagination.PerPage = paginationValue
		case "max_per_page":
			model.Pagination.MaxPerPage = paginationValue
		}
	}

	lexeme = self.scanIgnoreWhitespace()
	if lexeme.Token == EOF {
		return errors.New("Unexpected EOF")
	}
	if lexeme.Token != RIGHTBRACE {
		return fmt.Errorf("found %q, expected }", lexeme.Value)
	}
	return nil
}

func (self *Parser) parseModelActions(model *app.Model) error {
	lexeme := self.scanIgnoreWhitespace()
	if lexeme.Token == EOF {
		return errors.New("Unexpected EOF")
	}
	if lexeme.Token != LEFTBRACE {
		return fmt.Errorf("found %q, expected {", lexeme.Value)
	}

	lexeme = self.scanIgnoreWhitespace()
	if lexeme.Token != IDENT || lexeme.Value != "rest_actions" {
		self.unscan()
	} else {
		lexeme = self.scanIgnoreWhitespace()
		if lexeme.Token == EOF {
			return errors.New("Unexpected EOF")
		}
		if lexeme.Token != COLON {
			return fmt.Errorf("found %q, expected :", lexeme.Value)
		}
		lexeme = self.scanIgnoreWhitespace()
		if lexeme.Token == EOF {
			return errors.New("Unexpected EOF")
		}
		if lexeme.Token != LEFTSQBRACE {
			return fmt.Errorf("found %q, expected [", lexeme.Value)
		}
		model.Actions = &app.RestfulActions{}
		for {
			lexeme = self.scanIgnoreWhitespace()
			if lexeme.Token == RIGHTSQBRACE {
				self.unscan()
				break
			}
			if lexeme.Token == EOF {
				return errors.New("Unexpected EOF")
			}
			if lexeme.Token != IDENT {
				return fmt.Errorf("found %q, expected restful action [index, create, show, update, destroy]", lexeme.Value)
			}
			switch lexeme.Value {
			case "index":
				model.Actions.Index = true
			case "create":
				model.Actions.Create = true
			case "show":
				model.Actions.Show = true
			case "update":
				model.Actions.Update = true
			case "destroy":
				model.Actions.Destroy = true
			default:
				return fmt.Errorf("unexpected restful action. Expected one of [index, create, show, update, destroy], but got %q", lexeme.Value)
			}
			lexeme = self.scanIgnoreWhitespace()
			if lexeme.Token == EOF {
				return errors.New("Unexpected EOF")
			}
			if lexeme.Token == COMMA {
				continue
			} else {
				self.unscan()
			}
		}

		lexeme = self.scanIgnoreWhitespace()
		if lexeme.Token == EOF {
			return errors.New("Unexpected EOF")
		}
		if lexeme.Token != RIGHTSQBRACE {
			return fmt.Errorf("found %q, expected ]", lexeme.Value)
		}
	}

	lexeme = self.scanIgnoreWhitespace()
	if lexeme.Token == EOF {
		return errors.New("Unexpected EOF")
	}
	if lexeme.Token != RIGHTBRACE {
		return fmt.Errorf("found %q, expected }", lexeme.Value)
	}
	return nil
}
