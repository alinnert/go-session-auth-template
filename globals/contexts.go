package globals

// ContextKey is a key type for request contexts
type ContextKey int

const (
	// DBContext database request context key
	DBContext ContextKey = iota
	// SessionContext session context key
	SessionContext
	// ValidatorContext validator context key
	ValidatorContext
)
