package middleware

import (
	"net/http"

	logger "github.com/agusespa/flogg"
	"github.com/golang-jwt/jwt"
)

func CorsMiddleware(next http.Handler, domain string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", domain)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type GatewayClaims struct {
	jwt.StandardClaims
}

func GatewayMiddleware(next http.Handler, allowedIPs []string, logg logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO uncomment
		// clientIP := helpers.GetIP(r)
		// isAllowedIP := false
		// for _, ip := range allowedIPs {
		// 	if strings.TrimSpace(ip) == clientIP {
		// 		isAllowedIP = true
		// 		break
		// 	}
		// }
		// if !isAllowedIP {
		// 	logg.LogError(fmt.Errorf("unauthorized access attempt from IP: %s", clientIP))
		// 	http.Error(w, "Unauthorized source", http.StatusForbidden)
		// 	return
		// }

		next.ServeHTTP(w, r)
	})
}

func ChainMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
