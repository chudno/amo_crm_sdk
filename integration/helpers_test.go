//go:build integration

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/chudno/amo_crm_sdk/client"
)

// testTimeout — таймаут для каждого запроса к API.
const testTimeout = 15 * time.Second

// setupClient создаёт клиент amoCRM из переменных окружения.
// Пропускает тест, если переменные не заданы.
func setupClient(t *testing.T) (*client.Client, context.Context) {
	t.Helper()

	baseURL := os.Getenv("AMO_BASE_URL")
	token := os.Getenv("AMO_ACCESS_TOKEN")

	if baseURL == "" || token == "" {
		t.Skip("AMO_BASE_URL и AMO_ACCESS_TOKEN не заданы, пропускаем интеграционный тест")
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	t.Cleanup(func() { cancel() })

	return client.NewClient(baseURL, token), ctx
}
