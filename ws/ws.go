package ws

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nekoite/go-napcat/config"
	"go.uber.org/zap"
)

type Client struct {
	logger    *zap.Logger
	setupFunc func() (*websocket.Conn, error)
	conn      *websocket.Conn
	closeCh   chan struct{}
	send      chan []byte
	stopped   atomic.Bool
	closeOnce sync.Once

	writeWait    time.Duration
	pongWait     time.Duration
	pingPeriod   time.Duration
	maxReconnect int

	onRecvMsg func([]byte)
}

func NewConn(logger *zap.Logger, cfg *config.BotConfig, onRecvMsg func([]byte)) (*Client, error) {
	logger = logger.Named("ws")
	setupFunc := func() (*websocket.Conn, error) {
		u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", cfg.Ws.Host, cfg.Ws.Port), Path: cfg.Ws.Endpoint}
		logger.Info("connecting to", zap.String("url", u.String()))
		header := http.Header{}
		header.Add("Authorization", fmt.Sprintf("Bearer %s", cfg.Ws.Token))
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
	conn, err := setupFunc()
	if err != nil {
		logger.Error("dial:", zap.Error(err))
		return nil, err
	}
	wsConn := &Client{
		logger:       logger,
		setupFunc:    setupFunc,
		conn:         conn,
		closeCh:      make(chan struct{}),
		send:         make(chan []byte, 256),
		stopped:      atomic.Bool{},
		writeWait:    time.Duration(cfg.Ws.Timeout) * time.Millisecond,
		pongWait:     time.Duration(cfg.Ws.PongTimeout) * time.Millisecond,
		pingPeriod:   time.Duration(cfg.Ws.PingPeriod) * time.Millisecond,
		maxReconnect: cfg.Ws.MaxReconnect,
		onRecvMsg:    onRecvMsg,
	}
	return wsConn, nil
}

func (c *Client) Start() error {
	if err := c.prepareConn(c.conn); err != nil {
		return err
	}
	go c.run()
	return nil
}

func (c *Client) prepareConn(conn *websocket.Conn) error {
	err := conn.SetReadDeadline(time.Now().Add(c.pongWait))
	if err != nil {
		c.logger.Error("set read deadline", zap.Error(err))
		return err
	}
	conn.SetPongHandler(func(string) error {
		c.logger.Debug("received pong")
		if err := conn.SetReadDeadline(time.Now().Add(c.pongWait)); err != nil {
			c.logger.Error("set read deadline", zap.Error(err))
			return err
		}
		return nil
	})
	return nil
}

func (c *Client) run() {
	defer c.stopped.Store(true)

	for {
		err := c.runConn(c.conn)
		if c.stopped.Load() || err == nil {
			return
		}

		attempt := 0
		for {
			if c.stopped.Load() {
				return
			}
			if c.maxReconnect >= 0 && attempt >= c.maxReconnect {
				c.logger.Warn("max reconnect attempts reached", zap.Int("maxReconnect", c.maxReconnect), zap.Error(err))
				return
			}

			attempt++
			c.logger.Warn("connection closed unexpectedly, reconnecting after 3 seconds", zap.Int("attempt", attempt), zap.Int("maxReconnect", c.maxReconnect), zap.Error(err))
			time.Sleep(3 * time.Second)

			nextConn, reconnectErr := c.setupFunc()
			if reconnectErr != nil {
				err = reconnectErr
				c.logger.Error("failed to reconnect", zap.Int("attempt", attempt), zap.Error(reconnectErr))
				continue
			}

			if prepareErr := c.prepareConn(nextConn); prepareErr != nil {
				err = prepareErr
				c.logger.Error("failed to prepare reconnected websocket", zap.Int("attempt", attempt), zap.Error(prepareErr))
				_ = nextConn.Close()
				continue
			}

			c.conn = nextConn
			break
		}
	}
}

func (c *Client) runConn(conn *websocket.Conn) error {
	pumpDone := make(chan error, 2)
	stopWrite := make(chan struct{})

	go c.readPump(conn, pumpDone)
	go c.writePump(conn, stopWrite, pumpDone)

	err := <-pumpDone
	close(stopWrite)
	_ = conn.Close()
	secondErr := <-pumpDone

	if c.stopped.Load() {
		return nil
	}
	if err != nil {
		return err
	}
	return secondErr
}

func (c *Client) readPump(conn *websocket.Conn, pumpDone chan<- error) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("wsrecv", zap.Error(err))
			}
			pumpDone <- err
			return
		}
		c.logger.Debug("wsrecv", zap.String("message", string(message)))
		if c.onRecvMsg != nil {
			go c.onRecvMsg(message)
		}
	}
}

func (c *Client) writePump(conn *websocket.Conn, stopWrite <-chan struct{}, pumpDone chan<- error) {
	ticker := time.NewTicker(c.pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-stopWrite:
			pumpDone <- nil
			return
		case <-ticker.C:
			if err := c.writeMessage(conn, websocket.PingMessage, nil); err != nil {
				c.logger.Error("wssend", zap.Error(err))
				pumpDone <- err
				return
			}
			c.logger.Debug("sent ping")
		case <-c.closeCh:
			c.logger.Info("ws connection close")
			err := c.writeMessage(conn, websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				c.logger.Error("close", zap.Error(err))
			}
			pumpDone <- nil
			return
		case message, ok := <-c.send:
			if !ok {
				c.logger.Error("send channel closed")
				_ = c.writeMessage(conn, websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				pumpDone <- nil
				return
			}
			c.logger.Debug("wssend", zap.String("message", string(message)))
			if err := c.writeMessage(conn, websocket.TextMessage, message); err != nil {
				c.logger.Error("wssend", zap.Error(err))
				pumpDone <- err
				return
			}
		}
	}
}

func (c *Client) writeMessage(conn *websocket.Conn, messageType int, data []byte) error {
	err := conn.SetWriteDeadline(time.Now().Add(c.writeWait))
	if err != nil {
		c.logger.Error("set write deadline", zap.Error(err))
		return err
	}
	return conn.WriteMessage(messageType, data)
}

func (c *Client) Close() {
	c.stopped.Store(true)
	c.closeOnce.Do(func() {
		close(c.closeCh)
	})
}

func (c *Client) Send(msg []byte) {
	c.send <- msg
}
