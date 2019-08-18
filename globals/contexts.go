package globals

// ContextKey is a key type for request contexts
type ContextKey int

const (
	// DBContext provides access to the database
	DBContext ContextKey = iota
	// SessionContext provides access to the sessionManager
	SessionContext
	// ValidatorContext provides access to the validator
	ValidatorContext
	// UserContext represents a user struct
	UserContext
)
