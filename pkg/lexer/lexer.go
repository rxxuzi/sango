package lexer

import (
	"strings"
	"unicode"
)

// Lexer tokenizes Sango source code
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int
	column       int
}

// New creates a new Lexer
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	// Skip comments
	if l.ch == '/' {
		if l.peekChar() == '/' {
			l.skipLineComment()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.skipBlockComment()
			return l.NextToken()
		}
	}

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = NewToken(EQ, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = NewToken(DARROW, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else {
			tok = NewToken(ASSIGN, string(l.ch), tok.Line, tok.Column)
		}
	case '+':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = NewToken(PLUSASSIGN, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else {
			tok = NewToken(PLUS, string(l.ch), tok.Line, tok.Column)
		}
	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = NewToken(ARROW, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = NewToken(MINUSASSIGN, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else {
			tok = NewToken(MINUS, string(l.ch), tok.Line, tok.Column)
		}
	case '*':
		if l.peekChar() == '*' {
			ch := l.ch
			l.readChar()
			tok = NewToken(POWER, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = NewToken(ASTERISKASSIGN, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else {
			tok = NewToken(ASTERISK, string(l.ch), tok.Line, tok.Column)
		}
	case '/':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = NewToken(SLASHASSIGN, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else {
			tok = NewToken(SLASH, string(l.ch), tok.Line, tok.Column)
		}
	case '%':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = NewToken(PERCENTASSIGN, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else {
			tok = NewToken(PERCENT, string(l.ch), tok.Line, tok.Column)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = NewToken(AND, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = NewToken(AMPERSANDASSIGN, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else {
			tok = NewToken(AMPERSAND, string(l.ch), tok.Line, tok.Column)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = NewToken(OR, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = NewToken(PIPEASSIGN, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else {
			tok = NewToken(PIPE, string(l.ch), tok.Line, tok.Column)
		}
	case '^':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = NewToken(CARETASSIGN, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else {
			tok = NewToken(CARET, string(l.ch), tok.Line, tok.Column)
		}
	case '~':
		tok = NewToken(TILDE, string(l.ch), tok.Line, tok.Column)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = NewToken(NEQ, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else {
			tok = NewToken(NOT, string(l.ch), tok.Line, tok.Column)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = NewToken(LEQ, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = NewToken(LARROW, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else if l.peekChar() == '<' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = NewToken(LSHIFTASSIGN, "<<"+string(l.ch), tok.Line, tok.Column)
			} else {
				tok = NewToken(LSHIFT, string(ch)+string(l.ch), tok.Line, tok.Column)
			}
		} else {
			tok = NewToken(LT, string(l.ch), tok.Line, tok.Column)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = NewToken(GEQ, string(ch)+string(l.ch), tok.Line, tok.Column)
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = NewToken(RSHIFTASSIGN, ">>"+string(l.ch), tok.Line, tok.Column)
			} else {
				tok = NewToken(RSHIFT, string(ch)+string(l.ch), tok.Line, tok.Column)
			}
		} else {
			tok = NewToken(GT, string(l.ch), tok.Line, tok.Column)
		}
	case ';':
		tok = NewToken(SEMICOLON, string(l.ch), tok.Line, tok.Column)
	case ',':
		tok = NewToken(COMMA, string(l.ch), tok.Line, tok.Column)
	case '(':
		tok = NewToken(LPAREN, string(l.ch), tok.Line, tok.Column)
	case ')':
		tok = NewToken(RPAREN, string(l.ch), tok.Line, tok.Column)
	case '{':
		tok = NewToken(LBRACE, string(l.ch), tok.Line, tok.Column)
	case '}':
		tok = NewToken(RBRACE, string(l.ch), tok.Line, tok.Column)
	case '[':
		tok = NewToken(LBRACKET, string(l.ch), tok.Line, tok.Column)
	case ']':
		tok = NewToken(RBRACKET, string(l.ch), tok.Line, tok.Column)
	case ':':
		tok = NewToken(COLON, string(l.ch), tok.Line, tok.Column)
	case '.':
		if l.peekChar() == '.' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = NewToken(DOTDOTEQ, ".."+string(l.ch), tok.Line, tok.Column)
			} else {
				tok = NewToken(DOTDOT, string(ch)+string(l.ch), tok.Line, tok.Column)
			}
		} else {
			tok = NewToken(DOT, string(l.ch), tok.Line, tok.Column)
		}
	case '@':
		tok = NewToken(AT, string(l.ch), tok.Line, tok.Column)
	case '_':
		tok = NewToken(UNDERSCORE, string(l.ch), tok.Line, tok.Column)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			literal, isFloat := l.readNumber()
			tok.Literal = literal
			if isFloat {
				tok.Type = FLOAT
			} else {
				tok.Type = INT
			}
			return tok
		} else {
			tok = NewToken(ILLEGAL, string(l.ch), tok.Line, tok.Column)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
		if l.ch == '\n' {
			l.line++
			l.column = 0
		} else {
			l.column++
		}
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() (string, bool) {
	position := l.position
	isFloat := false

	for isDigit(l.ch) {
		l.readChar()
	}

	// Check for decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		isFloat = true
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	// Check for scientific notation
	if l.ch == 'e' || l.ch == 'E' {
		isFloat = true
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position], isFloat
}

func (l *Lexer) readString() string {
	var result strings.Builder
	l.readChar() // skip opening quote

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				result.WriteByte('\n')
			case 't':
				result.WriteByte('\t')
			case 'r':
				result.WriteByte('\r')
			case '\\':
				result.WriteByte('\\')
			case '"':
				result.WriteByte('"')
			default:
				result.WriteByte(l.ch)
			}
		} else {
			result.WriteByte(l.ch)
		}
		l.readChar()
	}

	return result.String()
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipLineComment() {
	// Skip //
	l.readChar()
	l.readChar()

	// Skip until end of line
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipBlockComment() {
	// Skip /*
	l.readChar()
	l.readChar()

	// Skip until */
	for l.ch != 0 {
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // skip *
			l.readChar() // skip /
			break
		}
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}