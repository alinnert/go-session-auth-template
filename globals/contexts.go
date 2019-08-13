package globals

// ContextKey is a key type for request contexts
type ContextKey string

// DBContext database request context key
const DBContext ContextKey = "db"

// SessionContext session context key
const SessionContext ContextKey = "session"

// ValidateRequest validator context key
const ValidatorContext ContextKey = "validate_request"
