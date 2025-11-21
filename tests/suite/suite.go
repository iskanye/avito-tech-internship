package suite

import (
	"context"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/iskanye/avito-tech-internship/internal/config"
	"github.com/iskanye/avito-tech-internship/pkg/api"
	"github.com/stretchr/testify/require"
)

type Suite struct {
	c *api.ClientWithResponses
}

func New(t *testing.T) (*Suite, context.Context) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath(configPath())
	cfg.LoadEnv()

	t.Cleanup(func() {
		t.Helper()
	})

	hc := http.Client{}
	c, err := api.NewClientWithResponses(
		net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		api.WithHTTPClient(&hc),
	)
	require.NoError(t, err)

	return &Suite{
		c: c,
	}, context.TODO()
}

func configPath() string {
	const key = "CONFIG_PATH"

	if v := os.Getenv(key); v != "" {
		return v
	}

	return "../config/tests.yaml"
}
