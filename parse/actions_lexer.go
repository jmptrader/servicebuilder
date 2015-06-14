package parse

func (self *Lexer) scanActions() {
	self.nextFunc = nil
	originalParentDepth := self.parenDepth
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
			lexeme := self.scanIdentifier()
			self.bus <- lexeme
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
			if self.parenDepth == originalParentDepth {
				break ScanLoop
			}
			continue
		case '[':
			self.bus <- Lexeme{LEFTSQBRACE, string(ch), *self.Position}
		case ']':
			self.bus <- Lexeme{RIGHTSQBRACE, string(ch), *self.Position}
		case ',':
			self.bus <- Lexeme{COMMA, string(ch), *self.Position}
		case ':':
			self.bus <- Lexeme{COLON, string(ch), *self.Position}
		}
		// Illegal token
	}
}
