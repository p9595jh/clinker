package reposh

type Raw []interface{}

type EntityParts[E any] struct {
	Entities   *[]E
	Entity     *E
	EntitiesFn func(*[]E)
	EntityFn   func(*E)
	Raw        Raw
}

type OrderBy struct {
	Column string
	Desc   bool
}

type FindOption[E any] struct {
	Select  []string
	Where   EntityParts[E]
	Order   OrderBy
	Limit   int
	Offset  int
	Preload []string
}

// to attach preloading fields
func (o *FindOption[E]) Preloads(s ...string) *FindOption[E] {
	o.Preload = append(o.Preload, s...)
	return o
}

// to find 'or' conditions
func (o *FindOption[E]) Entities(e *[]E) *FindOption[E] {
	o.Where.Entities = e
	return o
}

// to find 'and' conditions
func (o *FindOption[E]) Entity(e *E) *FindOption[E] {
	o.Where.Entity = e
	return o
}

// to find 'or' conditions && for embedded fields
func (o *FindOption[E]) EntitiesFn(f func(*[]E)) *FindOption[E] {
	o.Where.EntitiesFn = f
	return o
}

// to find 'and' conditions && for embedded fields
func (o *FindOption[E]) EntityFn(f func(*E)) *FindOption[E] {
	o.Where.EntityFn = f
	return o
}

func (o *FindOption[E]) Raw(raw *Raw) *FindOption[E] {
	o.Where.Raw = *raw
	return o
}
