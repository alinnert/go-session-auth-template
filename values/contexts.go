package values

// ContextKey is a key type for request contexts
type ContextKey string

// DBContext database request context key
const DBContext ContextKey = "db"

// SessionContext session context key
const SessionContext ContextKey = "session"
