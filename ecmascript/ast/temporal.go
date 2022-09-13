package ast

type TemporalEmptyArrowHead struct {
	BaseNode
}

func (t TemporalEmptyArrowHead) ESTree() interface{} {
	panic("TemporalEmptyArrowHead should not appear inside of ESTree.")
}

func (t TemporalEmptyArrowHead) ContainsTemporalNodes() bool {
	return true
}

type TemporalArrayRestElement struct {
	BaseNode
	BindingPattern
}

func (t TemporalArrayRestElement) ESTree() interface{} {
	panic("TemporalArrayRestElement should not appear inside of ESTree.")
}

func (t TemporalArrayRestElement) ContainsTemporalNodes() bool {
	return true
}

type TemporalObjectRestElement struct {
	BaseNode
	Identifier string
}

func (t TemporalObjectRestElement) ESTree() interface{} {
	panic("TemporalObjectRestElement should not appear inside of ESTree.")
}

func (t TemporalObjectRestElement) ContainsTemporalNodes() bool {
	return true
}

type TemporalFloatingRestElement struct {
	BaseNode
	Identifier string
}

func (t TemporalFloatingRestElement) ESTree() interface{} {
	panic("TemporalFloatingRestElement should not appear inside of ESTree.")
}

func (t TemporalFloatingRestElement) ContainsTemporalNodes() bool {
	return true
}
