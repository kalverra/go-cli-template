package clients

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestLoggingRoundTripper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		responseBody string
		expectedBody string
	}{
		{
			name:         "plain text body",
			responseBody: "OK",
			expectedBody: `"body":"OK"`,
		},
		{
			name:         "json body",
			responseBody: `{"status":"ok"}`,
			expectedBody: `"body":{"status":"ok"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup buffer for logs
			buf := &bytes.Buffer{}
			logger := zerolog.New(buf).Level(zerolog.TraceLevel)

			// Mock server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer ts.Close()

			// Middleware setup
			rt := NewLoggingRoundTripper(http.DefaultTransport, logger, "test")
			client := &http.Client{Transport: rt}

			// Action
			resp, err := client.Get(ts.URL)
			require.NoError(t, err)
			defer func() { require.NoError(t, resp.Body.Close()) }()

			// Verification
			logOutput := buf.String()
			require.NotEmpty(t, logOutput, "expected log output, got none")

			// Check for expected keys
			expectedKeys := []string{"status", "method", "url", "client", "duration", "body"}
			for _, key := range expectedKeys {
				require.Contains(t, logOutput, key, "log output missing key: %s", key)
			}
			require.Contains(t, logOutput, tt.expectedBody)

			// Ensure body is still readable
			gotBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, tt.responseBody, string(gotBody))
		})
	}
}
