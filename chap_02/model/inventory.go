package model

// inventory model
type Inventory struct {
	ID     string
	Data   string
	Cached int64
}

// interact with database
func NewInventory(id, data string, cached int64) Inventory {
	return Inventory{
		ID:     id,
		Data:   data,
		Cached: cached,
	}
}
