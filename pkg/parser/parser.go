package parser

import (
	"fmt"
	"io"
	"strconv"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/alexey-medvedchikov/lc3/pkg/bytecode"
)

type Program struct {
	Pos        lexer.Position
	Statements []*Statement `parser:"( @@ EOL* )*"`
}

type Statement struct {
	Pos       lexer.Position
	Labels    []*Label   `parser:"@@* Colon?" json:",omitempty"`
	Directive *Directive `parser:"( @@" json:",omitempty"`
	Op        *Op        `parser:"  | @@" json:",omitempty"`
	Trap      *Trap      `parser:"  | @@ )?" json:",omitempty"`
	Comment   *Comment   `parser:"@@?" json:",omitempty"`
}

type Label struct {
	Pos  lexer.Position
	Name *string `parser:"@Label" json:",omitempty"`
}

type Comment struct {
	Pos     lexer.Position
	Comment *string `parser:"@Comment" json:",omitempty"`
}

type Directive struct {
	Pos  lexer.Position
	Name *string         `parser:"@Directive" json:",omitempty"`
	Args []*DirectiveArg `parser:"@@*" json:",omitempty"`
}

type DirectiveArg struct {
	Pos    lexer.Position
	Label  *string `parser:"@Label" json:",omitempty"`
	Number *Number `parser:"| @Number" json:",omitempty"`
	String *String `parser:"| @String" json:",omitempty"`
}

type Op struct {
	Pos    lexer.Position
	OpCode *string   `parser:"@OpCode" json:",omitempty"`
	Args   []*OpArgs `parser:"( @@ ( ',' @@ )* )?" json:",omitempty"`
}

type OpArgs struct {
	Pos      lexer.Position
	Register *bytecode.Register `parser:"@Register" json:",omitempty"`
	Label    *string            `parser:"| @Label" json:",omitempty"`
	Number   *Number            `parser:"| @Number" json:",omitempty"`
}

type Trap struct {
	Name *string `parser:"@Trap" json:",omitempty"`
}

type Number int

func (n *Number) Capture(values []string) error {
	if len(values) != 1 {
		return fmt.Errorf("number can only capture single value: '%+v'", values)
	}

	v := values[0]
	base := -1
	switch v[0] {
	case 'x':
		base = 16
	case '#':
		base = 10
	}

	if base != -1 {
		i, err := strconv.ParseInt(v[1:], base, 64)
		if err != nil {
			return err
		}
		*n = Number(i)
		return nil
	}

	return fmt.Errorf("unknown number format: '%+v'", v)
}

type String string

func (s *String) Capture(values []string) error {
	if len(values) != 1 {
		return fmt.Errorf("string can only capture single value: '%+v'", values)
	}
	v := values[0]
	*s = String(v[1 : len(v)-1])

	return nil
}

var (
	asmLexer = lexer.MustSimple([]lexer.Rule{
		{Name: "EOL", Pattern: `[\r\n]+`},
		{Name: "Number", Pattern: `x-?[[:xdigit:]]+|#-?\d+`},
		{Name: "String", Pattern: `"[^"]*"`},
		{
			Name: "OpCode",
			Pattern: `(?i)\b(add|and|br|brnzp|brnz|brnp|brzp|brn|brz|bzp|jmp|jmpt|jsr` +
				`|jsrr|ld|ldi|ldr|lea|not|rti|st|sti|str|trap)\b`,
		},
		{Name: "Trap", Pattern: "(?i)\b(getc|in|out|puts|putsp|halt)\b"},
		{Name: "Directive", Pattern: `\.[[:alpha:]]\w*`},
		{Name: "Register", Pattern: `(?i)r\d`},
		{Name: "Label", Pattern: `[a-zA-Z0-9_]\w*`},
		{Name: "Comment", Pattern: `;.*`},
		{Name: "Comma", Pattern: `,`},
		{Name: "Colon", Pattern: `:`},
		{Name: "skip-whitespace", Pattern: `[[:blank:]]+`},
	})

	asmParser = participle.MustBuild(&Program{},
		participle.Lexer(asmLexer),
		participle.UseLookahead(1),
	)
)

func Parse(r io.Reader) (*Program, error) {
	var p Program

	if err := asmParser.Parse("", r, &p); err != nil {
		return nil, err
	}

	return &p, nil
}
