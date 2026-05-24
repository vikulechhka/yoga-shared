package middleware

import (
    "net/http"

    "github.com/vikulechhka/yoga-shared/logger"
    "github.com/vikulechhka/yoga-shared/response" 
    "go.uber.org/zap"
)

func RecoverMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                logger.Error("Panic recovered", zap.Any("error", err))
                response.WriteError(w, http.StatusInternalServerError, "internal server error")
            }
        }()
        next.ServeHTTP(w, r)
    })
}
