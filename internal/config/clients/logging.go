// Package clients provides logging for HTTP requests and responses.
package clients

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type loggingRoundTripper struct {
	next   http.RoundTripper
	logger zerolog.Logger
	name   string
}

// NewLoggingRoundTripper returns a new http.RoundTripper that logs responses.
func NewLoggingRoundTripper(next http.RoundTripper, logger zerolog.Logger, name string) http.RoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	return &loggingRoundTripper{
		next:   next,
		logger: logger,
		name:   name,
	}
}

func (l *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if !l.logger.Trace().Enabled() { //nolint:zerologlint
		return l.next.RoundTrip(req)
	}

	start := time.Now()

	resp, err := l.next.RoundTrip(req)

	duration := time.Since(start)

	event := l.logger.Trace().
		Str("client", l.name).
		Str("method", req.Method).
		Str("url", req.URL.String()).
		Str("duration", duration.String())

	if err != nil {
		event.Err(err).Msg("HTTP request failed")
		return resp, err
	}

	// Get response body and log it
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		event.Err(err).Msg("Failed to read response body")
		return resp, err
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if json.Valid(bodyBytes) {
		event.RawJSON("body", bodyBytes)
	} else {
		event.Str("body", string(bodyBytes))
	}

	event.Int("status", resp.StatusCode).
		Msg("HTTP response")

	return resp, nil
}
