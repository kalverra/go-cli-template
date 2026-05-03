package clients

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

type mockRoundTripper struct {
	resp *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.resp, nil
}

func BenchmarkRoundTripper_Baseline(b *testing.B) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("OK")),
	}
	rt := &mockRoundTripper{resp: resp}
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	for b.Loop() {
		r, err := rt.RoundTrip(req)
		require.NoError(b, err)
		err = r.Body.Close()
		require.NoError(b, err)
		// Reset body for next iteration
		resp.Body = io.NopCloser(bytes.NewBufferString("OK"))
	}
}

func BenchmarkLoggingRoundTripper_Simple(b *testing.B) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("OK")),
	}
	mock := &mockRoundTripper{resp: resp}
	logger := zerolog.New(io.Discard).Level(zerolog.TraceLevel)
	rt := NewLoggingRoundTripper(mock, logger, "test")
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	for b.Loop() {
		r, err := rt.RoundTrip(req)
		require.NoError(b, err)
		err = r.Body.Close()
		require.NoError(b, err)
		// Reset body for next iteration
		resp.Body = io.NopCloser(bytes.NewBufferString("OK"))
	}
}

func runBodyBenchmark(b *testing.B, size int, isJSON bool) {
	body := make([]byte, size)
	for i := range body {
		body[i] = 'a'
	}
	if isJSON {
		body = []byte(`{"data":"` + string(body) + `"}`)
	}

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(body)),
	}
	mock := &mockRoundTripper{resp: resp}
	logger := zerolog.New(io.Discard).Level(zerolog.TraceLevel)
	rt := NewLoggingRoundTripper(mock, logger, "test")
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	for b.Loop() {
		r, err := rt.RoundTrip(req)
		require.NoError(b, err)
		err = r.Body.Close()
		require.NoError(b, err)
		// Reset body for next iteration
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
	}
}

func BenchmarkLoggingRoundTripper_SmallBody(b *testing.B) {
	runBodyBenchmark(b, 128, false)
}

func BenchmarkLoggingRoundTripper_LargeBody(b *testing.B) {
	runBodyBenchmark(b, 1024*1024, false)
}

func BenchmarkLoggingRoundTripper_SmallJSON(b *testing.B) {
	runBodyBenchmark(b, 128, true)
}

func BenchmarkLoggingRoundTripper_LargeJSON(b *testing.B) {
	runBodyBenchmark(b, 1024*1024, true)
}
