package rpn

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	valueSplitter  *regexp.Regexp
	variableRegexp *regexp.Regexp
	funcRegexp     *regexp.Regexp
	stringSanitize *regexp.Regexp
)

func init() {
	valueSplitter = regexp.MustCompile(`[\(\)\,]`)
	variableRegexp = regexp.MustCompile(`^\$[0-9]+$`)
	funcRegexp = regexp.MustCompile(`^[A-Z][a-zA-Z0-9_]+$`)
	stringSanitize = regexp.MustCompile(`".+"`)
}

type Parser struct {
}

type RawRpnNode struct {
	ValueOk bool
	Name    string
}

func NewParser() *Parser {
	return new(Parser)
}

func (p *Parser) Parse(query string) ([]RawRpnNode, error) {
	query, santValues, err := p.sanitizeString(query)
	if err != nil {
		return nil, err
	}

	var nodes []RawRpnNode
	opStack := NewStack()
	tokens := p.split(query)

	for i, token := range tokens {
		if p.isValue(token) {
			nodes = append(nodes, RawRpnNode{
				ValueOk: true,
				Name:    token,
			})
			continue
		}

		if p.isFunction(tokens, i) || token == "(" {
			opStack.Push(token)

			if p.nextEquals(tokens, i+1, ",") {
				return nil, fmt.Errorf("misplaced ','")
			}

			if token != "(" && !funcRegexp.MatchString(token) {
				return nil, fmt.Errorf("invalid function name '%s'", token)
			}

			continue
		}

		if token == "," {
			poppedNodes, err := p.popToLeftPar(opStack)
			if err != nil {
				return nil, err
			}

			nodes = append(nodes, poppedNodes...)
			funcNode, ok := p.popFunc(opStack)
			if !ok {
				return nil, fmt.Errorf("misplaced ','")
			}

			if p.nextEquals(tokens, i+1, ")") {
				return nil, fmt.Errorf("misplaced ','")
			}

			opStack.Push(funcNode.Name)
			opStack.Push("(")
		}

		if token == ")" {
			poppedNodes, err := p.popToLeftPar(opStack)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, poppedNodes...)

			funcNode, ok := p.popFunc(opStack)
			if ok {
				nodes = append(nodes, funcNode)
			}
		}
	}

	if value, ok := opStack.Pop(); ok {
		return nil, fmt.Errorf("misplaced %s", value)
	}

	return p.putStringsBack(nodes, santValues), nil
}

func (p *Parser) sanitizeString(query string) (string, []string, error) {
	if index := strings.IndexAny(query, "^"); index >= 0 {
		return "", nil, fmt.Errorf("Invalid character '^'")
	}

	var values []string
	santQuery := stringSanitize.ReplaceAllStringFunc(query, func(value string) string {
		values = append(values, p.stripQuotes(value))
		return fmt.Sprintf("^%d", len(values)-1)
	})

	if strings.Index(santQuery, `"`) >= 0 {
		return "", nil, fmt.Errorf(`Non-matching '"'`)
	}

	return santQuery, values, nil
}

func (p *Parser) stripQuotes(value string) string {
	if len(value) < 2 {
		return value
	}

	return value[1 : len(value)-1]
}

func (p *Parser) putStringsBack(nodes []RawRpnNode, values []string) []RawRpnNode {
	for i, node := range nodes {
		if !node.ValueOk || node.Name[0] != '^' {
			continue
		}

		index, err := strconv.Atoi(node.Name[1:])
		if err != nil || index >= len(values) {
			panic(fmt.Sprintf("Unexpected index: %s", node.Name[1:]))
		}
		nodes[i] = RawRpnNode{
			ValueOk: true,
			Name:    values[index],
		}
	}
	return nodes
}

func (p *Parser) popToLeftPar(opStack *Stack) ([]RawRpnNode, error) {
	var nodes []RawRpnNode
	for {
		value, ok := opStack.Pop()
		if !ok {
			return nil, fmt.Errorf("invalid ')'")
		}

		if value == "(" {
			return nodes, nil
		}

		nodes = append(nodes, RawRpnNode{
			ValueOk: false,
			Name:    value,
		})
	}

	return nodes, nil
}

func (p *Parser) popFunc(opStack *Stack) (RawRpnNode, bool) {
	value, ok := opStack.Pop()
	if !ok {
		return RawRpnNode{}, false
	}

	if value == "(" {
		opStack.Push("(")
		return RawRpnNode{}, false
	}

	return RawRpnNode{
		ValueOk: false,
		Name:    value,
	}, true
}

func (p *Parser) split(query string) []string {
	var start int
	var tokens []string
	for _, index := range valueSplitter.FindAllStringIndex(query, -1) {
		tokens = p.appendNonEmpty(tokens, query[start:index[0]])
		tokens = p.appendNonEmpty(tokens, query[index[0]:index[1]])
		start = index[1]
	}

	return tokens
}

func (p *Parser) appendNonEmpty(tokens []string, value string) []string {
	value = strings.TrimSpace(value)
	if len(value) > 0 {
		tokens = append(tokens, value)
	}

	return tokens
}

func (p *Parser) nextEquals(tokens []string, index int, value string) bool {
	if len(tokens) <= index {
		return false
	}

	return tokens[index] == value
}

func (p *Parser) isValue(token string) bool {
	if token[0] == '^' {
		return true
	}

	_, err := strconv.ParseBool(token)
	if err == nil {
		return true
	}

	_, err = strconv.ParseFloat(token, 64)
	return err == nil || variableRegexp.MatchString(token)
}

func (p *Parser) isFunction(tokens []string, index int) bool {
	if index+1 >= len(tokens) {
		return false
	}

	return tokens[index+1] == "("
}
