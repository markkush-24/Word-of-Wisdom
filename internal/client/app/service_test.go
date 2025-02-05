package app

import (
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	core "word_of_wisdom/config"
)

// dummyLogger implements a minimal logger interface.
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

// fakeConn implements the net.Conn interface for testing.
type fakeConn struct{}

func (fc *fakeConn) Close() error                       { return nil }
func (fc *fakeConn) Read(b []byte) (int, error)         { return 0, nil }
func (fc *fakeConn) Write(b []byte) (int, error)        { return len(b), nil }
func (fc *fakeConn) LocalAddr() net.Addr                { return nil }
func (fc *fakeConn) RemoteAddr() net.Addr               { return nil }
func (fc *fakeConn) SetDeadline(t time.Time) error      { return nil }
func (fc *fakeConn) SetReadDeadline(t time.Time) error  { return nil }
func (fc *fakeConn) SetWriteDeadline(t time.Time) error { return nil }

// mockConnectionService – a minimal mock for network.ConnectionService.
type mockConnectionService struct {
	mock.Mock
}

func (m *mockConnectionService) Connect(address string) (net.Conn, error) {
	args := m.Called(address)
	connVal := args.Get(0)
	if connVal == nil {
		return nil, args.Error(1)
	}
	conn, ok := connVal.(net.Conn)
	if !ok {
		return nil, fmt.Errorf("expected net.Conn, got %T", connVal)
	}
	return conn, args.Error(1)
}

func (m *mockConnectionService) ReadMessage(conn net.Conn) (string, error) {
	args := m.Called(conn)
	return args.String(0), args.Error(1)
}

func (m *mockConnectionService) SendMessage(conn net.Conn, message string) error {
	args := m.Called(conn, message)
	return args.Error(0)
}

// mockPowSolver – a minimal mock for pow.PowSolver with the updated signature.
type mockPowSolver struct {
	mock.Mock
}

// The new SolvePoW method now accepts a challenge parameter.
func (m *mockPowSolver) SolvePoW(challenge string) string {
	args := m.Called(challenge)
	return args.String(0)
}

// TestStartSuccess verifies that the business logic executes successfully.
func TestStartSuccess(t *testing.T) { //nolint:paralleltest
	logger := &dummyLogger{}
	connService := new(mockConnectionService)
	powSolver := new(mockPowSolver)
	conn := &fakeConn{}

	// Expect a successful connection
	connService.
		On("Connect", core.ServerAddressDocker).
		Return(conn, nil)
	// First call to ReadMessage – receiving the unique challenge.
	connService.
		On("ReadMessage", conn).
		Return("challenge-message", nil).Once()
	// Expect that powSolver will be called with the parameter "challenge-message".
	powSolver.
		On("SolvePoW", "challenge-message").
		Return("solution")
	connService.
		On("SendMessage", conn, "solution").
		Return(nil)
	// Second call to ReadMessage – receiving the quote.
	connService.
		On("ReadMessage", conn).
		Return("quote-message", nil).Once()

	service := NewClientService(core.ServerAddressDocker, powSolver, connService, logger)
	err := service.Start()
	require.NoError(t, err)

	connService.AssertExpectations(t)
	powSolver.AssertExpectations(t)
}

// TestStart_ConnectionFailure verifies an error when establishing the connection.
func TestStart_ConnectionFailure(t *testing.T) { //nolint:paralleltest
	logger := &dummyLogger{}
	connService := new(mockConnectionService)
	powSolver := new(mockPowSolver)

	connErr := errors.New("connection error")
	connService.
		On("Connect", core.ServerAddressDocker).
		Return((net.Conn)(nil), connErr)

	service := NewClientService(core.ServerAddressDocker, powSolver, connService, logger)
	err := service.Start()
	require.Error(t, err)
	require.Contains(t, err.Error(), "connection error")

	connService.AssertExpectations(t)
}

// TestStart_ReadMessageFailure verifies an error when reading the challenge.
func TestStart_ReadMessageFailure(t *testing.T) { //nolint:paralleltest
	logger := &dummyLogger{}
	connService := new(mockConnectionService)
	powSolver := new(mockPowSolver)
	conn := &fakeConn{}

	connService.
		On("Connect", core.ServerAddressDocker).
		Return(conn, nil)
	readErr := errors.New("read error")
	connService.
		On("ReadMessage", conn).
		Return("", readErr).Once()

	service := NewClientService(core.ServerAddressDocker, powSolver, connService, logger)
	err := service.Start()
	require.Error(t, err)
	require.Contains(t, err.Error(), "read error")

	connService.AssertExpectations(t)
}

// TestStart_SendMessageFailure verifies an error when sending the solution.
func TestStart_SendMessageFailure(t *testing.T) { //nolint:paralleltest
	logger := &dummyLogger{}
	connService := new(mockConnectionService)
	powSolver := new(mockPowSolver)
	conn := &fakeConn{}

	connService.
		On("Connect", core.ServerAddressDocker).
		Return(conn, nil)
	connService.
		On("ReadMessage", conn).
		Return("challenge-message", nil).Once()
	powSolver.
		On("SolvePoW", "challenge-message").
		Return("solution")
	sendErr := errors.New("send error")
	connService.
		On("SendMessage", conn, "solution").
		Return(sendErr)

	service := NewClientService(core.ServerAddressDocker, powSolver, connService, logger)
	err := service.Start()
	require.Error(t, err)
	require.Contains(t, err.Error(), "send error")

	connService.AssertExpectations(t)
	powSolver.AssertExpectations(t)
}

// TestStart_QuoteReadFailure verifies an error when reading the quote after sending the solution.
func TestStart_QuoteReadFailure(t *testing.T) { //nolint:paralleltest
	logger := &dummyLogger{}
	connService := new(mockConnectionService)
	powSolver := new(mockPowSolver)
	conn := &fakeConn{}

	connService.
		On("Connect", core.ServerAddressDocker).
		Return(conn, nil)
	connService.
		On("ReadMessage", conn).
		Return("challenge-message", nil).Once()
	powSolver.
		On("SolvePoW", "challenge-message").
		Return("solution")
	connService.
		On("SendMessage", conn, "solution").
		Return(nil)
	quoteErr := errors.New("quote error")
	connService.
		On("ReadMessage", conn).
		Return("", quoteErr).Once()

	service := NewClientService(core.ServerAddressDocker, powSolver, connService, logger)
	err := service.Start()
	require.Error(t, err)
	require.Contains(t, err.Error(), "quote error")

	connService.AssertExpectations(t)
	powSolver.AssertExpectations(t)
}
