package pms

import (
	"fmt"
	"time"

	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/message"
)

// Connection maintains connections to an MPD server. Two separate connections
// are made: one for IDLE events, and another as a control connection. The IDLE
// connection is kept open continuously, while the control connection is
// allowed to time out.
//
// This class is used by calling the Run method as a goroutine.
type Connection struct {
	Host       string
	Port       string
	Password   string
	Connected  chan struct{}
	IdleEvents chan string
	messages   chan message.Message
	mpdClient  *mpd.Client
	mpdIdle    *mpd.Watcher
}

// NewConnection returns Connection.
func NewConnection(messages chan message.Message) *Connection {
	return &Connection{
		messages:   messages,
		Connected:  make(chan struct{}, 16),
		IdleEvents: make(chan string, 16),
	}
}

// MpdClient pings the MPD server and returns the client object if both the
// IDLE connection and control connection are ready. Otherwise this function
// returns nil.
func (c *Connection) MpdClient() (*mpd.Client, error) {
	var err error

	if c.mpdIdle == nil {
		return nil, fmt.Errorf("MPD connection is not ready.")
	}

	addr := makeAddress(c.Host, c.Port)

	if c.mpdClient != nil {
		err = c.mpdClient.Ping()
		if err == nil {
			return c.mpdClient, nil
		}
		console.Log("MPD control connection timeout.")
	}

	console.Log("Establishing MPD control connection to %+v...", addr)

	c.mpdClient, err = mpd.DialAuthenticated(addr.network, addr.addr, c.Password)
	if err != nil {
		return nil, fmt.Errorf("MPD control connection error: %s", err)
	}

	console.Log("Established MPD control connection.")

	return c.mpdClient, nil
}

// Open sets the host, port, and password parameters, closes any existing
// connections, and asynchronously connects to MPD as long as Run() is called.
func (c *Connection) Open(host, port, password string) {
	c.Close()
	c.Host = host
	c.Port = port
	c.Password = password
}

// Close closes any MPD connections.
func (c *Connection) Close() {
	if c.mpdClient != nil {
		c.mpdClient.Close()
	}
	if c.mpdIdle != nil {
		c.mpdIdle.Close()
	}
	c.mpdClient = nil
	c.mpdIdle = nil
}

// Run is the main goroutine of Connection. This thread will maintain an IDLE
// connection to the MPD server, and reconnect if the connection has errors.
func (c *Connection) Run() {
	for {
		// Try to connect IDLE connection until successful.
		err := c.connectIdle()
		if err != nil {
			c.Error("Error connecting to MPD: %s", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Emit signal.
		c.Connected <- struct{}{}

		// Relay all IDLE wakeups on the IdleEvents channel.
		go c.relayIdle()

		// Wait until there is a connection error, and clean up.
		for err = range c.mpdIdle.Error {
			c.Error("Error in MPD IDLE connection: %s", err)
			c.mpdClient.Close()
			c.mpdIdle.Close()
		}
	}
}

// Message sends a message on the message bus.
func (c *Connection) Message(format string, a ...interface{}) {
	c.messages <- message.Format(format, a...)
}

// Error sends an error message on the message bus.
func (c *Connection) Error(format string, a ...interface{}) {
	c.messages <- message.Errorf(format, a...)
}

// connectIdle establishes the IDLE connection to MPD.
func (c *Connection) connectIdle() error {
	var err error

	c.mpdClient = nil
	c.mpdIdle = nil

	addr := makeAddress(c.Host, c.Port)

	c.Message("Establishing MPD IDLE connection to %+v...", addr)

	c.mpdIdle, err = mpd.NewWatcher(addr.network, addr.addr, c.Password)
	if err != nil {
		return fmt.Errorf("MPD connection error: %s", err)
	}

	c.Message("Connected to MPD server %s.", addr)

	return err
}

// relayIdle relays IDLE events. This function will exit when there the IDLE
// connection is closed.
func (c *Connection) relayIdle() {
	for subsystem := range c.mpdIdle.Event {
		c.IdleEvents <- subsystem
	}
}
