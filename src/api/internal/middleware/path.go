package middleware

import (
    "net/http"
    "strings"
)

// PathPrefixMiddleware applies middleware only to paths not in the exclude list
func PathPrefixMiddleware(excludePaths []string, middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        // Chain the middlewares for non-excluded paths
        handler := next
        for i := len(middlewares) - 1; i >= 0; i-- {
            handler = middlewares[i](handler)
        }

        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Check if path should be excluded from middleware
            for _, prefix := range excludePaths {
                if strings.HasPrefix(r.URL.Path, prefix) {
                    next.ServeHTTP(w, r)
                    return
                }
            }
            
            // Apply middleware chain
            handler.ServeHTTP(w, r)
        })
    }
} 