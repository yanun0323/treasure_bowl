package model

type Pair [2]string

func NewPair(base, quote string) Pair {
	return [2]string{base, quote}
}

func (p Pair) String() string {
	return p[0] + "/" + p[1]
}

func (p *Pair) SetBase(base string) {
	p[0] = base
}

func (p *Pair) SetQuote(quote string) {
	p[1] = quote
}

func (p Pair) Base() string {
	return p[0]
}

func (p Pair) Quote() string {
	return p[1]
}
