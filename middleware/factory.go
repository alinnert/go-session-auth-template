package middleware

import "net/http"

// Middleware is a type that functions can return.
// This type is what needs to be passed to `router.Use`.
// If you write a function that returns this,
// you can make a middleware configurabel with an options parameter.
// Example: `router.Use(myFactory(options))`
type Middleware func(http.Handler) http.Handler
