package entity

import "strings"

type Pair [2]string

func NewPair(base, quote string) Pair {
	return [2]string{strings.ToUpper(base), strings.ToUpper(quote)}
}

// Uppercase return uppercase string of pair joined by sep char.
func (p Pair) Uppercase(sep ...string) string {
	if len(sep) != 0 && len(sep[0]) != 0 {
		return p[0] + sep[0] + p[1]
	}

	return p[0] + p[1]
}

// Lowercase return lowercase string of pair joined by sep char.
func (p Pair) Lowercase(sep ...string) string {
	return strings.ToLower(p.Uppercase(sep...))
}

func (p *Pair) SetBase(s string) {
	p[0] = strings.ToUpper(s)
}

func (p *Pair) SetQuote(s string) {
	p[1] = strings.ToUpper(s)
}

func (p Pair) Base() string {
	return p[0]
}

func (p Pair) Quote() string {
	return p[1]
}
