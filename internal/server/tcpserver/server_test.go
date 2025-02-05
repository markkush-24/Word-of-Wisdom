package tcpserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"
)

// Constants for tests.
const (
	// serverStartupDelay provides a delay that allows the server to start.
	serverStartupDelay = 100 * time.Millisecond
	// testTimeout defines the maximum waiting time in tests.
	testTimeout = 1 * time.Second
)

// dummyTCPHandler is a test implementation of the TCPHandler interface.
// When its Handle method is called, it sends the accepted connection to a channel
// and then attempts to close the connection, logging any error that occurs.
type dummyTCPHandler struct {
	handled chan net.Conn
}

func (d *dummyTCPHandler) Handle(ctx context.Context, conn net.Conn) {
	// Signal that the connection is accepted.
	d.handled <- conn
	// After handling, close the connection and log any error.
	if err := conn.Close(); err != nil {
		fmt.Printf("Error closing connection: %v\n", err) //nolint:forbidigo
	}
}

// dummyLogger is a simple implementation of the logger.Logger interface that does nothing.
type dummyLogger struct{}

func (d *dummyLogger) Debug(args ...interface{})                    {}
func (d *dummyLogger) Debugf(template string, args ...interface{})  {}
func (d *dummyLogger) Info(args ...interface{})                     {}
func (d *dummyLogger) Infof(template string, args ...interface{})   {}
func (d *dummyLogger) Warn(args ...interface{})                     {}
func (d *dummyLogger) Warnf(template string, args ...interface{})   {}
func (d *dummyLogger) Error(args ...interface{})                    {}
func (d *dummyLogger) Errorf(template string, args ...interface{})  {}
func (d *dummyLogger) DPanic(args ...interface{})                   {}
func (d *dummyLogger) DPanicf(template string, args ...interface{}) {}
func (d *dummyLogger) Panic(args ...interface{})                    {}
func (d *dummyLogger) Panicf(template string, args ...interface{})  {}
func (d *dummyLogger) Fatal(args ...interface{})                    {}
func (d *dummyLogger) Fatalf(template string, args ...interface{})  {}

// getFreePort returns a free port for testing. If obtaining a port fails,
// the test is terminated fatally.
func getFreePort(t *testing.T) int {
	t.Helper()
	// Open a temporary listener on port 0 so that the OS assigns a free port.
	listener, err := net.Listen("tcp", ":0") //nolint:gosec
	if err != nil {
		t.Fatalf("Failed to obtain a free port: %v", err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			t.Logf("Failed to close listener: %v", err)
		}
	}()

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("Failed to cast Addr to *net.TCPAddr")
	}
	return addr.Port
}

// TestStartServerListenError verifies that starting two servers on the same port results in an error.
func TestStartServerListenError(t *testing.T) { //nolint:paralleltest
	port := getFreePort(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := &dummyTCPHandler{handled: make(chan net.Conn, 1)}
	logger := &dummyLogger{}

	// Start the first server.
	serverErrChan := make(chan error, 1)
	go func() {
		serverErrChan <- StartServer(ctx, port, handler, logger)
	}()

	// Allow time for the first server to start.
	time.Sleep(serverStartupDelay)

	// Attempting to start a second server on the same port should return an error.
	err := StartServer(ctx, port, handler, logger)
	if err == nil {
		t.Error("Expected error when listening on an already used port")
	}

	// Cancel the context to shut down the first server.
	cancel()
	select {
	case err := <-serverErrChan:
		if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Expected context cancellation error, got: %v", err)
		}
	case <-time.After(testTimeout):
		t.Error("Timeout waiting for server shutdown")
	}
}

// TestServerHandlesConnection verifies that the server correctly accepts incoming connections
// and passes them to the handler.
func TestServerHandlesConnection(t *testing.T) { //nolint:paralleltest
	port := getFreePort(t)
	ctx, cancel := context.WithCancel(context.Background())
	// Ensure the context is cancelled at the end of the test.
	defer cancel()

	handler := &dummyTCPHandler{handled: make(chan net.Conn, 1)}
	logger := &dummyLogger{}
	serverErrChan := make(chan error, 1)

	// Start the server.
	go func() {
		serverErrChan <- StartServer(ctx, port, handler, logger)
	}()

	// Wait for the server to start.
	time.Sleep(serverStartupDelay)

	// Connect to the server.
	conn, err := net.Dial("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		t.Fatalf("Failed to connect to the server: %v", err)
	}
	// Ensure the connection is closed, logging any error.
	defer func(conn net.Conn) {
		if err := conn.Close(); err != nil {
			t.Logf("Error closing connection: %v", err)
		}
	}(conn)

	// Wait for the handler to receive the connection.
	select {
	case handledConn := <-handler.handled:
		if handledConn == nil {
			t.Error("Handler received nil instead of a valid connection")
		}
	case <-time.After(testTimeout):
		t.Error("Timeout waiting for connection to be handled")
	}

	// Cancel the context to shut down the server.
	cancel()
	select {
	case err := <-serverErrChan:
		if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Expected context cancellation error, got: %v", err)
		}
	case <-time.After(testTimeout):
		t.Error("Timeout waiting for server shutdown")
	}
}

// TestServerShutdownWithoutConnections verifies that the server correctly shuts down
// when no incoming connections occur and the context is cancelled.
func TestServerShutdownWithoutConnections(t *testing.T) { //nolint:paralleltest
	port := getFreePort(t)
	ctx, cancel := context.WithCancel(context.Background())
	// Cancel the context at the end of the test.
	defer cancel()

	handler := &dummyTCPHandler{handled: make(chan net.Conn, 1)}
	logger := &dummyLogger{}
	serverErrChan := make(chan error, 1)

	go func() {
		serverErrChan <- StartServer(ctx, port, handler, logger)
	}()

	// Allow time for the server to start.
	time.Sleep(serverStartupDelay)

	// Cancel the context, initiating server shutdown.
	cancel()
	select {
	case err := <-serverErrChan:
		if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Expected context cancellation error, got: %v", err)
		}
	case <-time.After(testTimeout):
		t.Error("Timeout waiting for server shutdown")
	}
}
