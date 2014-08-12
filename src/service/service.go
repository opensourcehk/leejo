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

// pointer to key
type KeyPtr interface{}

// pointer to parent key
type ParentKeyPtr interface{}

// pointer to condition
type CondPtr interface{}

// service interface
type Service interface {
	Create(EntityPtr) error
	List(ParentKeyPtr, ListCond, EntityListPtr) error
	Retrieve(KeyPtr, EntityListPtr) error
	Update(KeyPtr, EntityPtr) error
	Delete(KeyPtr) error
}
