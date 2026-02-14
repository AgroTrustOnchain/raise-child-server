package middleware

import (
	"raise-child/util"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type rateLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.RWMutex
	r   rate.Limit
	b   int
}

func InitializeRateLimiter(r rate.Limit, burst int) *rateLimiter {
	return &rateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   r,
		b:   burst,
	}
}

func (i *rateLimiter) getLimiter(address string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	res, isExist := i.ips[address]
	if !isExist {
		res = rate.NewLimiter(i.r, i.b)
		i.ips[address] = res
	}

	return res
}

func RateLimitMiddleware(limiter *rateLimiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		identifier, ok := ctx.Value("address").(string)
		if !ok || identifier == "" {
			identifier = ctx.ClientIP()
		}

		if !limiter.getLimiter(identifier).Allow() {
			util.ProcessResponse(util.GetRateLimitRequestBodyResponse(ctx))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
