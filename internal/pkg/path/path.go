package path

import "fmt"

type Type string

const (
	Trace      Type = "trace"
	Generation Type = "generation"
	Span       Type = "span"
)

type Path []Element

type Element struct {
	Type Type
	ID   string
}

func (p *Path) Push(t Type, id string) {
	*p = append(*p, Element{Type: t, ID: id})
}

func (p *Path) Pop() *Element {
	if len(*p) == 0 {
		return nil
	}

	lastElement := (*p)[len(*p)-1]
	*p = (*p)[:len(*p)-1]
	return &lastElement
}

func (p *Path) PopIf(t Type) (*Element, error) {
	if len(*p) == 0 {
		return nil, fmt.Errorf("path is empty")
	}

	lastElement := (*p)[len(*p)-1]
	if lastElement.Type != t {
		return nil, fmt.Errorf("expected %s, got %s", t, lastElement.Type)
	}

	*p = (*p)[:len(*p)-1]
	return &lastElement, nil
}

func (p *Path) At(i int) *Element {
	if i < 0 || i >= len(*p) {
		return nil
	}
	return &(*p)[i]
}

func (p *Path) Last() Element {
	if len(*p) == 0 {
		return Element{}
	}
	return (*p)[len(*p)-1]
}
