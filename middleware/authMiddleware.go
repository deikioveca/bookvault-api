package middleware

import (
	"BookVault-API/helper"
	"BookVault-API/jwt"
	"net/http"
	"strings"
)


func AuthMiddleware(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                helper.WriteError(w, http.StatusUnauthorized, "missing authorization header")
                return
            }

            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                helper.WriteError(w, http.StatusUnauthorized, "invalid authorization header format")               
                return
            }

            tokenString := parts[1]
            claims, err := jwt.ParseToken(tokenString)
            if err != nil {
                helper.WriteError(w, http.StatusUnauthorized, "invalid token")               
                return
            }
            
            
            if len(allowedRoles) > 0 {
                role, ok := claims["role"].(string)
                if !ok {
                    helper.WriteError(w, http.StatusUnauthorized, "invalid token claims")                    
                    return
                }

                authorized := false
                for _, r := range allowedRoles {
                    if role == r {
                        authorized = true
                        break
                    }
                }

                if !authorized {
                    helper.WriteError(w, http.StatusForbidden, "forbidden: insufficient permissions")                   
                    return
                }
            }
            next(w, r)
        }
    }
}