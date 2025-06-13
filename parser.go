package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Program struct {
	Repo   string
	Branch string
	Body   []Stmt
}

type Parser struct {
	tokens  []Token
	current int
}

type ForStmt struct {
	Iterator Token
	Iterable Expr
	Body     []Stmt
}

type IfStatment struct {
	Condition Expr
	ThenBody  []Stmt
	ElseBody  []Stmt
}

type VarStmt struct {
	Name  Token
	Value Expr
}

type Stmt interface {
	isStmt()
}

type Expr interface {
	isExpr()
}

type CommandStmt struct {
	Name Token
	Args []Expr
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

type LiteralExpr struct {
	Value interface{}
}

type VariableExpr struct {
	Name Token
}

func (VariableExpr) isExpr() {}

func (UnaryExpr) isExpr() {}

func (LiteralExpr) isExpr() {}

func (CommandStmt) isStmt() {}

func (ForStmt) isStmt() {}

func (IfStatment) isStmt() {}

func (BinaryExpr) isExpr() {}

func (VarStmt) isStmt() {}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(ttype TokenType, msg string) (Token, error) {
	if p.check(ttype) {
		return p.advance(), nil
	}
	return Token{}, fmt.Errorf("%s at line %d", msg, p.peek().Line)
}

func (p *Parser) check(ttype TokenType) bool {
	return !p.isAtEnd() && p.peek().Type == ttype
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) peek() Token {
	if p.isAtEnd() {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens) || p.tokens[p.current].Type == TokenEOF
}

func (p *Parser) parseVarStmt() (Stmt, error) {
	name, err := p.consume(TokenIdent, "Expected variable name after 'var'")
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(TokenEqal, "Expected '-' after varaible name"); err != nil {
		return nil, err
	}

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return VarStmt{
		Name:  name,
		Value: value,
	}, nil
}

func (p *Parser) parseForStatement() (Stmt, error) {
	iterator, err := p.consume(TokenIdent, "Expected loop variable after 'for'")
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(TokenIn, "Expected 'in' after loop variable"); err != nil {
		return nil, err
	}

	iterable, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(TokenLBrace, "Expected '{' at start of for loop"); err != nil {
		return nil, err
	}

	var body []Stmt
	for !p.check(TokenRBrace) && !p.isAtEnd() {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		body = append(body, stmt)
	}

	if _, err := p.consume(TokenRBrace, "Expected '}' after for loop body"); err != nil {
		return nil, err
	}

	return ForStmt{
		Iterator: iterator,
		Iterable: iterable,
		Body:     body,
	}, nil

}

func (p *Parser) parseEquality() (Expr, error) {
	expr, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.match(TokenEqality, TokenBang) {
		operator := p.previous()
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) parseComparison() (Expr, error) {
	expr, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.match(TokenGreater, TokenGreaterEq, TokenLess, TokenLessEq) {
		operator := p.previous()
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) parseTerm() (Expr, error) {
	expr, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.match(TokenPlus, TokenMinus) {
		operator := p.previous()
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) parseFactor() (Expr, error) {
	expr, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.match(TokenStar, TokenFWDSlash) {
		operator := p.previous()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) parseUnary() (Expr, error) {
	if p.match(TokenBang, TokenMinus) {
		operator := p.previous()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return UnaryExpr{
			Operator: operator,
			Right:    right,
		}, nil
	}
	return p.parsePrimary()
}

func (p *Parser) parseBlock() ([]Stmt, error) {
	if _, err := p.consume(TokenIf, "Expected '{' after if conditon"); err != nil {
		return nil, err
	}
	var statements []Stmt
	for !p.check(TokenRBrace) && !p.isAtEnd() {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	if _, err := p.consume(TokenRBrace, "Expected '}' after if block"); err != nil {
		return nil, err
	}
	return statements, nil

}

func (p *Parser) parseIfStatement() (Stmt, error) {
	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	thenBody, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	var elseBody []Stmt
	if p.match(TokenElse) {
		if p.match(TokenIf) {
			elseStmt, err := p.parseIfStatement()
			if err != nil {
				return nil, err
			}
			elseBody = []Stmt{elseStmt}
		} else {
			elseBody, err = p.parseBlock()
			if err != nil {
				return nil, err
			}
		}
	}

	return IfStatment{
		Condition: condition,
		ThenBody:  thenBody,
		ElseBody:  elseBody,
	}, nil
}

func (p *Parser) parseExpression() (Expr, error) {
	return p.parseEquality()
}

func (p *Parser) parseCommand() (Stmt, error) {
	name, err := p.consume(TokenIdent, "Expected command name")
	if err != nil {
		return nil, err
	}

	var args []Expr
	for !p.check(TokenRBrace) && !p.check(TokenEOF) {
		if p.check(TokenVar) || p.check(TokenIf) || p.check(TokenFor) {
			break
		}
		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return CommandStmt{
		Name: name,
		Args: args,
	}, nil
}

func (p *Parser) parseStatement() (Stmt, error) {
	switch {
	case p.match(TokenVar):
		return p.parseVarStmt()
	case p.match(TokenIf):
		return p.parseIfStatement()
	case p.match(TokenFor):
		return p.parseForStatement()
	case p.check(TokenIdent):
		return p.parseCommand()
	}

	return nil, nil
}

func (p *Parser) parsePrimary() (Expr, error) {
	switch {
	case p.match(TokenFalse):
		return LiteralExpr{Value: false}, nil
	case p.match(TokenTrue):
		return LiteralExpr{Value: true}, nil
	case p.match(TokenNil):
		return LiteralExpr{Value: nil}, nil
	case p.match(TokenNumber):
		val, err := strconv.ParseFloat(p.previous().Lexeme, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number '%s'", p.previous().Lexeme)
		}
		return LiteralExpr{Value: val}, nil
	case p.match(TokenString):
		lex := p.previous().Lexeme
		unquoted := strings.Trim(lex, `"`)
		return LiteralExpr{Value: unquoted}, nil
	case p.match(TokenIdent):
		return VariableExpr{Name: p.previous()}, nil
	case p.match(TokenLParen):
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(TokenRParen, "Expected ')' after expression"); err != nil {
			return nil, err
		}
		return expr, nil
	default:
		return nil, fmt.Errorf("unexpected token '%s' at line %d", p.peek().Lexeme, p.peek().Line)
	}
}

func (p *Parser) Parse() ([]Stmt, error) {
	var statements []Stmt
	for !p.isAtEnd() {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	return statements, nil
}

// func (p *Parser) parseStatement() {
// 	switch{
// 	case p.match(TokenFor):
// 		return p.parse
// 	}
// }
