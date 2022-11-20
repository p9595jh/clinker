package dto

type Transformable interface {
	Transform() (map[string]any, error)
}

// sustain its type and only change value by reference (pointer)
type ReferProcessable interface {
	Process() error
}
