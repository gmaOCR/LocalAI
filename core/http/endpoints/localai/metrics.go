package localai

import (
	"time"

	"github.com/go-skynet/LocalAI/pkg/schema"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MetricsHandler() fiber.Handler {
	return adaptor.HTTPHandler(promhttp.Handler())
}

type apiMiddlewareConfig struct {
	Filter  func(c *fiber.Ctx) bool
	metrics *schema.LocalAIMetrics
}

func MetricsAPIMiddleware(metrics *schema.LocalAIMetrics) fiber.Handler {
	cfg := apiMiddlewareConfig{
		metrics: metrics,
		Filter: func(c *fiber.Ctx) bool {
			return c.Path() == "/metrics"
		},
	}

	return func(c *fiber.Ctx) error {
		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}
		path := c.Path()
		method := c.Method()

		start := time.Now()
		err := c.Next()
		elapsed := float64(time.Since(start)) / float64(time.Second)
		cfg.metrics.ObserveAPICall(method, path, elapsed)
		return err
	}
}
