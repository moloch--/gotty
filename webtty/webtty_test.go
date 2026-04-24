package webtty

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"testing"
)

type pipePair struct {
	*io.PipeReader
	*io.PipeWriter
}

type testSlave struct {
	reader io.Reader
	writer io.Writer
}

func (s *testSlave) Read(p []byte) (int, error) {
	return s.reader.Read(p)
}

func (s *testSlave) Write(p []byte) (int, error) {
	return s.writer.Write(p)
}

func (s *testSlave) WindowTitleVariables() map[string]interface{} {
	return nil
}

func (s *testSlave) ResizeTerminal(columns int, rows int) error {
	return nil
}

func TestWriteFromPTY(t *testing.T) {
	connInPipeReader, connInPipeWriter := io.Pipe() // in to conn
	connOutPipeReader, _ := io.Pipe()               // out from conn
	slaveOutPipeReader, slaveOutPipeWriter := io.Pipe()

	conn := pipePair{
		connOutPipeReader,
		connInPipeWriter,
	}
	slave := &testSlave{
		reader: slaveOutPipeReader,
		writer: io.Discard,
	}
	dt, err := New(conn, slave)
	if err != nil {
		t.Fatalf("Unexpected error from New(): %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	errs := make(chan error, 1)
	go func() {
		errs <- dt.Run(ctx)
	}()

	buf := make([]byte, 1024)
	n, err := connInPipeReader.Read(buf)
	if err != nil {
		t.Fatalf("Unexpected error from Read(): %s", err)
	}
	if n != 1 || buf[0] != SetWindowTitle {
		t.Fatalf("Unexpected initialize message `%s`", buf[:n])
	}

	message := []byte("foobar")
	n, err = slaveOutPipeWriter.Write(message)
	if err != nil {
		t.Fatalf("Unexpected error from Write(): %s", err)
	}
	if n != len(message) {
		t.Fatalf("Write() accepted `%d` for message `%s`", n, message)
	}

	n, err = connInPipeReader.Read(buf)
	if err != nil {
		t.Fatalf("Unexpected error from Read(): %s", err)
	}
	if buf[0] != Output {
		t.Fatalf("Unexpected message type `%c`", buf[0])
	}
	decoded := make([]byte, 1024)
	n, err = base64.StdEncoding.Decode(decoded, buf[1:n])
	if err != nil {
		t.Fatalf("Unexpected error from Decode(): %s", err)
	}
	if !bytes.Equal(decoded[:n], message) {
		t.Fatalf("Unexpected message received: `%s`", decoded[:n])
	}

	cancel()
	if err := <-errs; err != context.Canceled {
		t.Fatalf("Unexpected error from Run(): %s", err)
	}
}

func TestWriteFromConn(t *testing.T) {
	connInPipeReader, connInPipeWriter := io.Pipe()   // in to conn
	connOutPipeReader, connOutPipeWriter := io.Pipe() // out from conn
	slaveOutPipeReader, _ := io.Pipe()
	slaveInPipeReader, slaveInPipeWriter := io.Pipe()

	conn := pipePair{
		connOutPipeReader,
		connInPipeWriter,
	}
	slave := &testSlave{
		reader: slaveOutPipeReader,
		writer: slaveInPipeWriter,
	}

	dt, err := New(conn, slave, WithPermitWrite())
	if err != nil {
		t.Fatalf("Unexpected error from New(): %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	errs := make(chan error, 1)
	go func() {
		errs <- dt.Run(ctx)
	}()

	var (
		message []byte
		n       int
	)
	readBuf := make([]byte, 1024)
	n, err = connInPipeReader.Read(readBuf)
	if err != nil {
		t.Fatalf("Unexpected error from Read(): %s", err)
	}
	if n != 1 || readBuf[0] != SetWindowTitle {
		t.Fatalf("Unexpected initialize message `%s`", readBuf[:n])
	}

	// input
	message = append([]byte{Input}, []byte("hello\n")...) // line buffered canonical mode
	n, err = connOutPipeWriter.Write(message)
	if err != nil {
		t.Fatalf("Unexpected error from Write(): %s", err)
	}
	if n != len(message) {
		t.Fatalf("Write() accepted `%d` for message `%s`", n, message)
	}

	n, err = slaveInPipeReader.Read(readBuf)
	if err != nil {
		t.Fatalf("Unexpected error from Write(): %s", err)
	}
	if !bytes.Equal(readBuf[:n], message[1:]) {
		t.Fatalf("Unexpected message received: `%s`", readBuf[:n])
	}

	// ping
	message = []byte{Ping}
	n, err = connOutPipeWriter.Write(message)
	if err != nil {
		t.Fatalf("Unexpected error from Write(): %s", err)
	}
	if n != len(message) {
		t.Fatalf("Write() accepted `%d` for message `%s`", n, message)
	}

	n, err = connInPipeReader.Read(readBuf)
	if err != nil {
		t.Fatalf("Unexpected error from Read(): %s", err)
	}
	if !bytes.Equal(readBuf[:n], []byte{Pong}) {
		t.Fatalf("Unexpected message received: `%s`", readBuf[:n])
	}

	// TODO: resize

	cancel()
	if err := <-errs; err != context.Canceled {
		t.Fatalf("Unexpected error from Run(): %s", err)
	}
}
