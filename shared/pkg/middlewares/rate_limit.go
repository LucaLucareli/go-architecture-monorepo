package middlewares

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips   map[string]*rate.Limiter
	mu    *sync.RWMutex
	r     rate.Limit
	b     int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips:   make(map[string]*rate.Limiter),
		mu:    &sync.RWMutex{},
		r:     r,
		b:     b,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.ips[ip]
	i.mu.RUnlock()

	if !exists {
		i.mu.Lock()
		defer i.mu.Unlock()
		
		if limiter, exists = i.ips[ip]; exists {
			return limiter
		}

		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

// RateLimitMiddleware protege a API contra excesso de requisições usando Token Bucket por IP.
func RateLimitMiddleware(rps float64, burst int) echo.MiddlewareFunc {
	limiterManager := NewIPRateLimiter(rate.Limit(rps), burst)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			limiter := limiterManager.GetLimiter(ip)

			if !limiter.Allow() {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"message": "Muitas requisições. Tente novamente mais tarde.",
				})
			}

			return next(c)
		}
	}
}
