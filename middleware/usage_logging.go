package middleware

import (
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request information
		//fmt.Printf("[%s] %s from IP: %s at %s\n", r.Method, r.URL.Path, getClientIP(r), time.Now().Format(time.RFC3339))
		//we wil put this in log files
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// getClientIP extracts the original IP address from X-Forwarded-For header
// func getClientIP(r *http.Request) string {
// 	ip := r.Header.Get("X-Forwarded-For")
// 	if ip == "" {
// 		// If X-Forwarded-For is empty, use RemoteAddr
// 		ip = r.RemoteAddr
// 	} else {
// 		// X-Forwarded-For may contain a list of IPs, take the first one
// 		ips := strings.Split(ip, ",")
// 		ip = strings.TrimSpace(ips[0])
// 	}
// 	return ip
// }
