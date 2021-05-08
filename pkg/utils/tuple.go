package utils

// Tuple wrapped two elements
type Tuple struct {
	First, Second interface{}
}

func NewTuple(a, b interface{}) *Tuple {
	return &Tuple{First: a, Second: b}
}
