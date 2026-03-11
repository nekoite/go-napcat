package ws

import (
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nekoite/go-napcat/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testUpgrader = websocket.Upgrader{}

func newTestClient(t *testing.T, serverURL string, maxReconnect int) *Client {
	t.Helper()

	parsedURL, err := url.Parse(serverURL)
	require.NoError(t, err)

	host, portText, err := net.SplitHostPort(parsedURL.Host)
	require.NoError(t, err)

	port, err := strconv.Atoi(portText)
	require.NoError(t, err)

	cfg := config.DefaultBotConfig(123456, "token").WithWs(host, port, "/")
	cfg.Ws.PingPeriod = 60000
	cfg.Ws.PongTimeout = 60000
	cfg.Ws.Timeout = 1000
	cfg.Ws.MaxReconnect = maxReconnect

	client, err := NewConn(zap.NewNop(), cfg, nil)
	require.NoError(t, err)

	return client
}

func TestClientStopsAfterMaxReconnectAttempts(t *testing.T) {
	acceptedConn := make(chan *websocket.Conn, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade failed: %v", err)
			return
		}
		acceptedConn <- conn
	}))
	defer server.Close()

	client := newTestClient(t, server.URL, 2)
	originalSetup := client.setupFunc
	var reconnectAttempts atomic.Int32
	client.setupFunc = func() (*websocket.Conn, error) {
		reconnectAttempts.Add(1)
		return originalSetup()
	}

	require.NoError(t, client.Start())

	firstConn := <-acceptedConn
	require.NoError(t, firstConn.Close())
	server.Close()

	require.Eventually(t, func() bool {
		return client.stopped.Load()
	}, 20*time.Second, 300*time.Millisecond)
	require.EqualValues(t, 2, reconnectAttempts.Load())
}

func TestClientResetsReconnectAttemptsAfterSuccessfulReconnect(t *testing.T) {
	acceptedConn := make(chan *websocket.Conn, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade failed: %v", err)
			return
		}
		acceptedConn <- conn
	}))
	defer server.Close()

	client := newTestClient(t, server.URL, 2)
	originalSetup := client.setupFunc
	var reconnectAttempts atomic.Int32
	client.setupFunc = func() (*websocket.Conn, error) {
		attempt := reconnectAttempts.Add(1)
		if attempt > 1 {
			return nil, errors.New("reconnect failed")
		}
		return originalSetup()
	}

	require.NoError(t, client.Start())

	firstConn := <-acceptedConn
	require.NoError(t, firstConn.Close())

	secondConn := <-acceptedConn
	require.NoError(t, secondConn.Close())

	require.Eventually(t, func() bool {
		return client.stopped.Load()
	}, 20*time.Second, 300*time.Millisecond)
	require.EqualValues(t, 3, reconnectAttempts.Load())
}
