package parse

func (self *Lexer) scanPagination() {
	self.nextFunc = nil
	originalParenDepth := self.parenDepth
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
			self.bus <- self.scanIdentifier()
			continue
		} else if isDigit(ch) {
			self.unread()
			self.bus <- self.scanNumber()
		}

		switch ch {
		case eof:
			self.unread()
			break ScanLoop
		case '{':
			self.parenDepth++
			self.bus <- Lexeme{LEFTBRACE, string(ch), *self.Position}
			continue
		case ':':
			self.bus <- Lexeme{COLON, string(ch), *self.Position}
			continue
		case '}':
			self.parenDepth--
			self.bus <- Lexeme{RIGHTBRACE, string(ch), *self.Position}
			if self.parenDepth == originalParenDepth {
				break ScanLoop
			}
			continue
		}
	}
}
