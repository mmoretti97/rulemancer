//go:generate re2go $INPUT -o $OUTPUT --api simple
package rulemancer

import (
	"errors"
	"fmt"
)

const (
	ScopeMain = 0 + iota
	ScopeComment
	ScopeBlock
	ScopeDefTemplate
	ScopeSlot
	ScopeMultiSlot
)

type scopeData struct {
	blockDepth int
	name       string
}

type scopeLevel struct {
	*ProtocolData
	*scopeData
	scopeId int
	prev    *scopeLevel
	next    *scopeLevel
}

func newScopeLevel(e *ProtocolData, scope int) (*scopeLevel, int) {
	return &scopeLevel{ProtocolData: e, scopeData: new(scopeData), scopeId: scope, prev: nil, next: nil}, scope
}

func (s *scopeLevel) push(scope int) (*scopeLevel, int) {
	sl, _ := newScopeLevel(s.ProtocolData, scope)
	sl.prev = s
	s.next = sl
	if s.Debug {
		fmt.Println(purple("scope:"), sl)
	}
	return sl, sl.scopeId
}

func (s *scopeLevel) pop() (*scopeLevel, int) {
	if s.prev == nil {
		return s, ScopeMain
	}
	prev := s.prev
	s.prev = nil
	prev.next = nil
	if s.Debug {
		fmt.Println(purple("scope:"), prev)
	}
	return prev, prev.scopeId
}

func (s *scopeLevel) descend(scope int) (*scopeLevel, int) {
	for {
		if s.scopeId == scope {
			return s, s.scopeId
		}

		if s.prev == nil {
			break
		}
		newDesc := s.prev
		s.prev = nil
		s.next = nil
		newDesc.next = nil
		s = newDesc
	}
	return s, ScopeMain
}

func (s *scopeLevel) isInsideScope(scope int) *scopeLevel {
	for desc := s; desc != nil; desc = desc.prev {
		if desc.scopeId == scope {
			return desc
		}
	}
	return nil
}

func (s *scopeLevel) isCurrentScope(scope int) bool {
	return s.scopeId == scope
}

func (s *scopeLevel) currentScope() int {
	return s.scopeId
}

func (s *scopeLevel) String() string {
	result := ""
	desc := s
	for {
		if desc == nil {
			break
		}
		result = cyan(scopeName(desc.scopeId)) + result
		if desc.prev != nil {
			result = ", " + result
		}
		desc = desc.prev
	}
	result = "[" + result + "]"
	return result
}

func scopeName(scope int) string {
	switch scope {
	case ScopeMain:
		return "main"
	case ScopeComment:
		return "comment"
	case ScopeBlock:
		return "block"
	case ScopeDefTemplate:
		return "deftemplate"
	case ScopeSlot:
		return "slot"
	case ScopeMultiSlot:
		return "multislot"
	}
	return "unknown"
}

// Returns "fake" terminating null if cursor has reached limit.
func peek(str string, cur int) byte {
	if cur >= len(str) {
		return 0 // fake null
	} else {
		return str[cur]
	}
}

func (e *ProtocolData) Compile(yyinput string) error {
	yycursor := 0
	yytext := 0
	yymarker := 0
	prev := 0
	sS, scope := newScopeLevel(e, ScopeMain)

	for {
		if e.Debug && e.DebugLevel >= debugLevelMax {
			fmt.Println("----")
			fmt.Println("yycursor:", yycursor, "yytext:", yytext, "yymarker:", yymarker)
			fmt.Println("scopeStack:", sS.String())
			fmt.Println("scope:", scopeName(scope))
			fmt.Println(">>>>")
		}

		switch scope {
		case ScopeMain:
		 	/*!re2c
			re2c:yyfill:enable = 0;
			re2c:YYCTYPE = byte;
			re2c:YYPEEK = "peek(str, cur)";
			re2c:YYSKIP = "cur += 1";

			comment = ";";
			varname = [a-z][A-Za-z0-9-_]*;
			blockstart = "(";
			blockend = ")";
			deftemplate = "deftemplate";
			slot = "slot";
			multislot = "multislot";
			w = [ \t]+;

			*      { return errors.New("Unexpected input: "+string(yyinput[prev:yycursor])) }
			[\x00] { return nil }
			comment {
					sS, scope = sS.push(ScopeComment)
					prev = yycursor
					continue
				}
			blockstart {
					sS, scope = sS.push(ScopeBlock)
					sS.blockDepth = sS.prev.blockDepth + 1
					prev = yycursor
					continue
				}
			w 	{
					//fmt.Println("space:", yyinput[prev:yycursor])
					prev = yycursor
					continue
				}
			[\n]+ {
					//fmt.Println("newline:", yyinput[prev:yycursor])
					prev = yycursor
					continue
				}
			*/
		case ScopeComment:
			/*!re2c
			re2c:yyfill:enable = 0;
			re2c:YYCTYPE = byte;
			re2c:YYPEEK = "peek(str, cur)";
			re2c:YYSKIP = "cur += 1";
			
			[\x00] { return errors.New("Unexpected end of input") }
			"\n"   {
					sS, scope = sS.pop()
					prev = yycursor
					continue
				}
			*      {
					prev = yycursor
					continue
				}
			*/
		case ScopeBlock:
			/*!re2c
			re2c:yyfill:enable = 0;
			re2c:YYCTYPE = byte;
			re2c:YYPEEK = "peek(str, cur)";
			re2c:YYSKIP = "cur += 1";
			
			[\x00] { return errors.New("Unexpected end of input") }
			comment {
					sS, scope = sS.push(ScopeComment)
					prev = yycursor
					continue
				}
			blockstart {
					sS, scope = sS.push(ScopeBlock)
					sS.blockDepth = sS.prev.blockDepth + 1
					prev = yycursor
					continue
				}
			blockend {
					sS, scope = sS.pop()
					prev = yycursor
					continue
				}
			deftemplate {
					sS, scope = sS.push(ScopeDefTemplate)
					prev = yycursor
					continue
				}
			slot {
					sS, scope = sS.push(ScopeSlot)
					prev = yycursor
					continue
				}
			multislot {
					sS, scope = sS.push(ScopeMultiSlot)
					prev = yycursor
					continue
				}
			*      {
					prev = yycursor
					continue
				}
			*/
		case ScopeDefTemplate:
			/*!re2c
			re2c:yyfill:enable = 0;
			re2c:YYCTYPE = byte;
			re2c:YYPEEK = "peek(str, cur)";
			re2c:YYSKIP = "cur += 1";
			
			[\x00] { return errors.New("Unexpected end of input") }
			comment {
					sS, scope = sS.push(ScopeComment)
					prev = yycursor
					continue
				}
			blockstart {
					sS, scope = sS.push(ScopeBlock)
					sS.blockDepth = sS.prev.blockDepth + 1
					prev = yycursor
					continue
				}
			blockend {	
					sS, scope = sS.descend(ScopeMain)
					prev = yycursor
					continue
				}
			varname {
					//fmt.Println("deftemplate varname:", yyinput[prev:yycursor])
					relation:= yyinput[prev:yycursor]
					if _, exists := sS.Slots[relation]; exists {
						sS.name = relation
					}
					prev = yycursor
					continue
				}
			*      {
					prev = yycursor
					continue
				}
			*/
		case ScopeMultiSlot:
			/*!re2c
			re2c:yyfill:enable = 0;
			re2c:YYCTYPE = byte;
			re2c:YYPEEK = "peek(str, cur)";
			re2c:YYSKIP = "cur += 1";
			
			[\x00] { return errors.New("Unexpected end of input") }
			comment {
					sS, scope = sS.push(ScopeComment)
					prev = yycursor
					continue
				}
			blockstart {
					sS, scope = sS.push(ScopeBlock)
					sS.blockDepth = sS.prev.blockDepth + 1
					prev = yycursor
					continue
				}
			blockend {	
					sS, scope = sS.descend(ScopeDefTemplate)
					prev = yycursor
					continue
				}
			*      {
					prev = yycursor
					continue
				}
			*/
		case ScopeSlot:
			/*!re2c
			re2c:yyfill:enable = 0;
			re2c:YYCTYPE = byte;
			re2c:YYPEEK = "peek(str, cur)";
			re2c:YYSKIP = "cur += 1";
			
			[\x00] { return errors.New("Unexpected end of input") }
			comment {
					sS, scope = sS.push(ScopeComment)
					prev = yycursor
					continue
				}
			blockstart {
					sS, scope = sS.push(ScopeBlock)
					sS.blockDepth = sS.prev.blockDepth + 1
					prev = yycursor
					continue
				}
			blockend {	
					sS, scope = sS.descend(ScopeDefTemplate)
					prev = yycursor
					continue
				}
			varname {
					//fmt.Println("slot varname:", yyinput[prev:yycursor])
					slotName := yyinput[prev:yycursor]
					deftempl:= sS.isInsideScope(ScopeDefTemplate)
					if deftempl != nil && deftempl.name != "" {
						if slotList, exists := sS.Slots[deftempl.name]; exists {
							sS.Slots[deftempl.name] = append(slotList, slotName)
						}
					}
					prev = yycursor
					continue
				}
			*      {
					prev = yycursor
					continue
				}
			*/
		}
	}
}
