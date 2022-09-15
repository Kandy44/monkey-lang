package lexer

import (
	"monkey_lang/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

type FType int

const (
	FInitial         FType = iota
	FRationalBegin   FType = iota
	FRational        FType = iota
	FFractionalBegin FType = iota
	FFractional      FType = iota
)

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhiteSpace()
	switch l.ch {
	case '=':

		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		if l.peekChar() == '*' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EXP, Literal: literal}
		} else {
			tok = newToken(token.ASTERISK, l.ch)
		}

	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)

	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isNumber(l.ch) {
			tok := l.readNum()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func isNumber(ch byte) bool {
	return isDigit(ch) || isSign(ch) || isExponent(ch)
}

func (l *Lexer) readNum() token.Token {
	position := l.position
	var t token.Token
	state := FInitial

	for {
		l.readChar()
		// fmt.Println(l.ch)
		// fmt.Println(state)
		// fmt.Println(l.ch == ';')

		if l.position == len(l.input) || (!isDigit(l.ch) && !isSign(l.ch) && !isExponent(l.ch) && l.ch != '.') {
			break
		}

		switch state {
		case FInitial:
			if isSign(l.ch) {
				state = FRationalBegin
			} else if isDigit(l.ch) {
				state = FRational
			} else if l.ch == '.' {
				state = FFractionalBegin
			} else {
				return newToken(token.ILLEGAL, l.ch)
			}
		case FRationalBegin:
			if isDigit(l.ch) {
				state = FRational
			} else if l.ch == '.' {
				state = FFractionalBegin
			} else {
				return newToken(token.ILLEGAL, l.ch)
			}
		case FRational:
			if isDigit(l.ch) {

			} else if l.ch == '.' {
				state = FFractional
			} else {
				return newToken(token.ILLEGAL, l.ch)
			}
		case FFractionalBegin:
			if isDigit(l.ch) {
				state = FFractional
			} else {
				return newToken(token.ILLEGAL, l.ch)
			}
		case FFractional:
			if isDigit(l.ch) {
				state = FFractional
			} else {
				return newToken(token.ILLEGAL, l.ch)
			}
		}
		// l.readChar()
	}

	// fmt.Printf("Cur num: %s\n", l.input[position:l.position])
	// fmt.Println(state)

	if state == FInitial || state == FRational || state == FFractional {
		t.Literal = l.input[position:l.position]
		t.Type = token.INT

		if !isDecimal(t.Literal) {
			t.Type = token.FLOAT
		}
		return t
	}

	return newToken(token.ILLEGAL, l.ch)
}

// else if isFloat(l.ch) {
// 	tok.Type = token.FLOAT
// 	tok.Literal = l.readFloat()
// 	return tok
// } else if isDigit(l.ch) {
// 	tok.Type = token.INT
// 	tok.Literal = l.readNumber()
// 	return tok
// }

// func isFloat(ch byte) bool {
// 	return isDigit(ch) || (ch == '+' || ch == '-')
// }

func isDecimal(str string) bool {
	for i := 0; i < len(str); i++ {
		if !isDigit(str[i]) {
			return false
		}
	}
	return true
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isSign(ch byte) bool {
	return ch == '+' || ch == '-'
}

func isExponent(ch byte) bool {
	return ch == 'e' || ch == 'E'
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// func (l *Lexer) readNumber() string {
// 	position := l.position
// 	for isDigit(l.ch) {
// 		l.readChar()
// 	}
// 	return l.input[position:l.position]
// }

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}
