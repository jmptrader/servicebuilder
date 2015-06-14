package parse

import (
	"strings"
)

func (self *Lexer) scanModelIdentifier() Lexeme {
	lexeme := self.scanIdentifier()
	switch strings.ToLower(lexeme.Value) {
	case "fields":
		lexeme.Token = FIELDS
		self.nextFunc = self.scanFields
	case "pagination":
		lexeme.Token = PAGINATION
		self.nextFunc = self.scanPagination
	case "actions":
		lexeme.Token = ACTIONS
		self.nextFunc = self.scanActions
	}
	return lexeme
}

func (self *Lexer) scanModel() {
	self.nextFunc = nil
ScanLoop:
	for {
		if self.nextFunc != nil {
			self.nextFunc()
		}
		ch := self.read()
		if isWhitespace(ch) {
			self.unread()
			self.bus <- self.scanWhitespace()
			continue
		} else if isLetter(ch) {
			self.unread()
			self.bus <- self.scanModelIdentifier()
			continue
		}
		switch ch {
		case eof:
			self.unread()
			break ScanLoop
		case '{':
			self.parenDepth++
			self.bus <- Lexeme{LEFTBRACE, string(ch), *self.Position}
			continue
		case '}':
			self.parenDepth--
			self.bus <- Lexeme{RIGHTBRACE, string(ch), *self.Position}
			continue
		}
		self.bus <- Lexeme{ILLEGAL, string(ch), *self.Position}
	}
}
