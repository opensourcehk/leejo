package service

// pointer to entiity
type EntityPtr interface{}

// pointer to list of entiity
type EntityListPtr interface{}

// pointer to entity list
type ListPtr interface{}

// interface listing conditions
type ListCond interface {
	GetLimit() uint64
	GetOffset() uint64
}

// basic implementation of listing conditions
type BasicListCond struct {
	Limit  uint64
	Offset uint64
}

func (c *BasicListCond) GetLimit() uint64 {
	return c.Limit
}
func (c *BasicListCond) GetOffset() uint64 {
	return c.Offset
}

// pointer to condition
type CondPtr interface{}

// struct to store context keys
type Context interface {
	Get(string) string
	GetCond() ListCond
}

// basic implementation of context
type BasicContext struct {
	Values Values
	Cond   ListCond
}

func (c *BasicContext) Get(key string) string {
	return c.Values.Get(key)
}
func (c *BasicContext) GetCond() ListCond {
	return c.Cond
}

// values in a context
type Values interface {
	Get(string) string
}

// service interface
type Service interface {
	Create(Context, EntityPtr) error
	List(Context, EntityListPtr) error
	Retrieve(Context, EntityListPtr) error
	Update(Context, EntityPtr) error
	Delete(Context) error
}
