package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

type (
	Entry struct {
		Type          EntryType
		Data          []byte
		StartPosition int
	}

	EntryType int
)

const (
	UNKNOWN = iota
	VARIABLE
	COMMENT
	STRING_CONSTANT
	CONSTANT
	OPPERAND
	OPPERAND_SYMBOL
	RESERVED
	FUNCTION
	NUMBER
	CLASS
)

var reservedStrings = []string{
	"abstract",
	"and",
	"array",
	"as",
	"break",
	"callable",
	"case",
	"catch",
	"class",
	"clone",
	"const",
	"continue",
	"declare",
	"default",
	"die",
	"do",
	"echo",
	"else",
	"elseif",
	"empty",
	"enddeclare",
	"endfor",
	"endforeach",
	"endif",
	"endswitch",
	"endwhile",
	"eval",
	"exit",
	"extends",
	"final",
	"finally",
	"for",
	"foreach",
	"function",
	"global",
	"goto",
	"if",
	"implements",
	"include",
	"include_once",
	"instanceof",
	"insteadof",
	"interface",
	"isset",
	"list",
	"namespace",
	"new",
	"or",
	"print",
	"private",
	"protected",
	"public",
	"require",
	"require_once",
	"return",
	"static",
	"switch",
	"throw",
	"trait",
	"try",
	"unset",
	"use",
	"var",
	"while",
}

var operators = []string{
	"[", "]",
	"**",
	"++",
	"--",
	"~", "(int)", "(float)", "(string)", "(array)", "(object)", "(bool)",
	"!",
	"*",
	"/",
	"%",
	"+", "-", ".",
	"<<",
	">>",
	"<",
	"<=",
	">", ">=",
	"==", "!=", "===", "!==", "<>", "<=>",
	"&",
	"^",
	"|",
	"&&",
	"||",
	"??",
	"?", ":",
	"=", "+=", "-=", "*=", "**=", "/=", ".=", "%=", "&=", "|=", "^=", "<<=", ">>=",
	"->",
	"and",
	"xor",
	"or",
}

var operatorSymbols = []string{
	"[",
	"]",
	"(",
	")",
	"{",
	"}",
	"~",
	"!",
	"*",
	//"/",
	"%",
	"+",
	"-",
	".",
	"<",
	">",
	"&",
	"^",
	"|",
	"?",
	":",
	"=",
	",",
	";",
}
var operatorSymbolsString = "[]()~!*%+-.<>&^|?:=,;"

func main() {
	//c := color.New(color.FgCyan)
	//c.Println("Prints cyan text")
	//
	////c.DisableColor()
	//fmt.Println("This is printed without any color")
	//
	////c.EnableColor()
	//c.Println("This prints again cyan...")
	//
	//return

	var Entries []Entry

	f, err := os.Open("request.php")
	if err != nil {
		log.Error(err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Error(err)
	}
	//
	//for _, line := range bytes.Split(data, []byte("\n")) {
	//	for _, subLine := range bytes.Split(line, []byte(" ")) {
	//		for _, byteEntry := range bytes.Split(subLine, []byte("\t")) {
	//			Entries = append(Entries, Entry{
	//				Type: UNKNOWN,
	//				Data: byteEntry,
	//			}, Entry{
	//				Type: SPLITTER,
	//			})
	//		}
	//	}
	//	Entries[len(Entries)-1].Type = LINE_SPLITTER
	//}

	var it = Iterator{
		StateMachine: buildStateMachine(operatorSymbols),
		data:         data,
	}

	for {
		e, err := it.Next()

		if err == io.EOF {
			break
		}

		if e != nil {
			Entries = append(Entries, *e)
		}
	}

	Entries = ProceedEntries(Entries, data)

	var (
		lPos   int
		colors = map[EntryType]*color.Color{
			UNKNOWN:         color.New(color.FgBlack, color.BgWhite),
			COMMENT:         color.New(color.FgWhite),
			STRING_CONSTANT: color.New(color.FgHiRed),
			CONSTANT:        color.New(color.FgWhite),
			OPPERAND:        color.New(color.FgHiBlue),
			OPPERAND_SYMBOL: color.New(color.FgRed),
			RESERVED:        color.New(color.FgGreen, color.Bold),
			NUMBER:          color.New(color.FgHiYellow),
			VARIABLE:        color.New(color.FgHiMagenta),
			FUNCTION:        color.New(color.FgBlue),
			CLASS:           color.New(color.FgBlue),
		}
	)
	for _, e := range Entries {
		if e.StartPosition > lPos {
			fmt.Print(string(data[lPos:e.StartPosition]))
		}

		colors[e.Type].Print(string(e.Data))
		lPos = e.StartPosition + len(e.Data)

		//fmt.Println(e.Type, string(e.Data))
	}
	fmt.Print(string(data[lPos:]))
}

type State struct {
	Next map[byte]int
	Zero int
	Else int
}

type StateMachine []State

func buildStateMachine(operatorSymbols []string) StateMachine {
	var m = make(StateMachine, 15)

	s := State{
		Next: map[byte]int{
			'"':  1,
			'\'': 12,
			'/':  5,
			'\n': 10,
			'\t': 10,
			' ':  10,
		},
		Else: 0,
		Zero: -1,
	}
	// operators
	for _, sc := range operatorSymbols {
		s.Next[sc[0]] = 11
	}
	//s.Next['*'] = 11
	m[0] = s

	// [CHAIN] String constant
	m[1] = State{
		Next: map[byte]int{
			'\\': 2,
			'"':  4,
		},
		Else: 1,
		Zero: -1,
	}

	m[2] = State{
		Next: map[byte]int{
			'\\': 3,
		},
		Else: 1,
		Zero: -1,
	}

	m[3] = State{
		Next: map[byte]int{
			'\\': 2,
			'"':  4,
		},
		Else: 1,
		Zero: -1,
	}

	m[4] = State{
		Next: map[byte]int{},
		Else: -1,
		Zero: 0,
	}

	// [CHAIN] String constant
	m[12] = State{
		Next: map[byte]int{
			'\\': 13,
			'\'': 4,
		},
		Else: 12,
		Zero: -1,
	}

	m[13] = State{
		Next: map[byte]int{
			'\\': 14,
		},
		Else: 12,
		Zero: -1,
	}

	m[14] = State{
		Next: map[byte]int{
			'\\': 13,
			'\'': 4,
		},
		Else: 12,
		Zero: -1,
	}

	// [CHAIN] Comment
	m[5] = State{
		Next: map[byte]int{
			'/': 6,
			'*': 8,
		},
		Else: -1,
		Zero: 0,
	}

	m[6] = State{
		Next: map[byte]int{
			'\n': 7,
		},
		Else: 6,
		Zero: -1,
	}

	m[7] = State{
		Next: map[byte]int{},
		Else: -1,
		Zero: 0,
	}

	m[8] = State{
		Next: map[byte]int{
			'*': 9,
		},
		Else: 8,
		Zero: -1,
	}

	m[9] = State{
		Next: map[byte]int{
			'*': 9,
			'/': 7,
		},
		Else: 8,
		Zero: -1,
	}

	// [CHAIN] Split
	m[10] = State{
		Next: map[byte]int{},
		Else: -1,
		Zero: 0,
	}

	// [CHAIN] Operators
	m[11] = State{
		Next: map[byte]int{},
		Else: -1,
		Zero: 0,
	}

	return m
}

type Iterator struct {
	StateMachine
	data  []byte
	state int
	pos   int
	lPos  int
}

func (it *Iterator) Next() (*Entry, error) {
	if it.pos >= len(it.data) {
		return nil, io.EOF
	}
	c := it.data[it.pos]
	//fmt.Println("\t", it.state, string(c), it.pos)
	it.pos++
	st := it.StateMachine[it.state]
	next, ok := st.Next[c]
	if !ok {
		if st.Else != -1 {
			next = st.Else
		} else {
			it.state = st.Zero
			it.pos--
			return it.Next()
		}
	}

	var entry *Entry
	if it.state < next {
		switch next {
		case 11:
			entry = &Entry{
				Data:          it.data[it.lPos:it.pos],
				Type:          OPPERAND_SYMBOL,
				StartPosition: it.lPos,
			}
			it.lPos = it.pos
		case 12:
			entry = &Entry{
				Data:          it.data[it.lPos : it.pos-1],
				Type:          UNKNOWN,
				StartPosition: it.lPos,
			}
			it.lPos = it.pos - 1
		}
	}

	it.state = next

	switch it.state {
	case 4:
		entry = &Entry{
			Data:          it.data[it.lPos:it.pos],
			Type:          STRING_CONSTANT,
			StartPosition: it.lPos,
		}
		it.lPos = it.pos
	case 7:
		entry = &Entry{
			Data:          it.data[it.lPos:it.pos],
			Type:          COMMENT,
			StartPosition: it.lPos,
		}
		it.lPos = it.pos
	case 10:
		entry = &Entry{
			Data:          it.data[it.lPos : it.pos-1],
			Type:          UNKNOWN,
			StartPosition: it.lPos,
		}
		it.lPos = it.pos
	case 5:
		if it.data[it.pos] != '/' && it.data[it.pos] != '*' {
			entry = &Entry{
				Data:          it.data[it.lPos:it.pos],
				Type:          OPPERAND_SYMBOL,
				StartPosition: it.lPos,
			}
			it.lPos = it.pos
		}
	}

	if it.state == 0 && entry == nil {
		if strings.Contains(operatorSymbolsString, string(it.data[it.pos:it.pos+1])) {
			entry = &Entry{
				Data:          it.data[it.lPos:it.pos],
				Type:          UNKNOWN,
				StartPosition: it.lPos,
			}
			it.lPos = it.pos
		}
	}

	if entry != nil && len(entry.Data) == 0 {
		entry = nil
	}

	return entry, nil
}

func ProceedEntries(entries []Entry, data []byte) []Entry {
	var (
		n           int
		k           = len(entries)
		numberRegex = regexp.MustCompile(`^([1-9][0-9]*|0|0[xX][0-9a-fA-F]+|0[0-7]+|0b[01]+|[0-9]*[\.][0-9]+|[0-9]+[\.][0-9]*|([0-9]+|[0-9]*[\.][0-9]+|[0-9]+[\.][0-9]*)[eE][+-]?[0-9]+)$`)
	)

	var opperands = map[string]bool{}

	for _, op := range operators {
		opperands[string(op)] = true
	}

	var reserved = map[string]bool{}

	for _, res := range reservedStrings {
		reserved[res] = true
	}

	for i := 0; i < k; i++ {
		//fmt.Println("\t\t", string(entries[i].Data))
		var lvl = 1
		if entries[i].Type == OPPERAND_SYMBOL {
			if i+1 < k && entries[i+1].Type == OPPERAND_SYMBOL && entries[i].StartPosition+1 == entries[i+1].StartPosition {
				if i+2 < k && entries[i+2].Type == OPPERAND_SYMBOL && entries[i+1].StartPosition+1 == entries[i+2].StartPosition {
					if _, ok := opperands[string(data[entries[i].StartPosition:entries[i].StartPosition+3])]; ok {
						lvl = 3
						entries[i].Data = data[entries[i].StartPosition : entries[i].StartPosition+3]
						entries[i].Type = OPPERAND
					}
				}

				if _, ok := opperands[string(data[entries[i].StartPosition:entries[i].StartPosition+2])]; lvl < 3 && ok {
					lvl = 2
					entries[i].Data = data[entries[i].StartPosition : entries[i].StartPosition+2]
					entries[i].Type = OPPERAND
				}
			}

			if _, ok := opperands[string(data[entries[i].StartPosition:entries[i].StartPosition+1])]; lvl < 2 && ok {
				entries[i].Type = OPPERAND
			}
			//fmt.Println(string(entries[i].Data), " ", lvl-1)
		}

		if entries[i].Type == UNKNOWN {
			if entries[i].Data[0] == '$' {
				entries[i].Type = VARIABLE
			} else if _, ok := reserved[string(entries[i].Data)]; ok {
				entries[i].Type = RESERVED
			} else if numberRegex.Match(entries[i].Data) {
				entries[i].Type = NUMBER
			} else if i+1 < k && entries[i+1].Type == OPPERAND_SYMBOL && string(entries[i+1].Data) == "(" {
				entries[i].Type = FUNCTION
			} else if i > 0 && entries[n-1].Type == OPPERAND && string(entries[n-1].Data) == "->" {
				entries[i].Type = FUNCTION
			} else if i > 0 && entries[n-1].Type == RESERVED && string(entries[n-1].Data) == "new" {
				entries[i].Type = CLASS
			}
		}

		entries[n] = entries[i]
		n++
		i += lvl - 1
	}

	return entries[:n]
}
