package storage

import (
	"encoding/base64"
	"encoding/json"
	"sync"

	"github.com/gin-gonic/gin"
)

// workerVars stores environment variables forwarded from the Cloudflare Worker
// via the X-Worker-Vars request header. Cloudflare Containers don't expose
// [vars] as OS env vars, so the Worker relays them in this header.
var (
	workerVars   = make(map[string]string)
	workerVarsMu sync.RWMutex
)

// SetWorkerVar stores a variable forwarded from the Worker.
func SetWorkerVar(key, value string) {
	workerVarsMu.Lock()
	workerVars[key] = value
	workerVarsMu.Unlock()
}

// GetWorkerVar retrieves a variable forwarded from the Worker.
func GetWorkerVar(key string) string {
	workerVarsMu.RLock()
	defer workerVarsMu.RUnlock()
	return workerVars[key]
}

// WorkerVarMiddleware reads the X-Worker-Vars header and stores the decoded
// key-value pairs globally for use by storage providers and other components.
func WorkerVarMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("X-Worker-Vars")
		if raw == "" {
			c.Next()
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(raw)
		if err != nil {
			c.Next()
			return
		}

		var vars map[string]string
		if err := json.Unmarshal(decoded, &vars); err != nil {
			c.Next()
			return
		}

		workerVarsMu.Lock()
		for k, v := range vars {
			if _, exists := workerVars[k]; !exists {
				workerVars[k] = v
			}
		}
		workerVarsMu.Unlock()

		c.Next()
	}
}
