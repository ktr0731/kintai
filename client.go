package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/jason0x43/go-toggl"
	"github.com/skratchdot/open-golang/open"
)

type Client struct {
	ssid string

	duration  time.Duration
	logger    *log.Logger
	toggl     Toggler
	started   bool
	timeEntry *toggl.TimeEntry
}

func NewClient(logger *log.Logger, ssid, apiToken string, duration time.Duration) *Client {
	return &Client{
		ssid:     ssid,
		logger:   logger,
		toggl:    NewTogglClient(apiToken),
		duration: duration,
	}
}

func (c *Client) Start() error {
	c.logger.Println("start pouring...")

	exitCh := createSigCh()

	ch := time.Tick(c.duration)
	for {
		select {
		case <-exitCh:
			c.logger.Println("pouring finished")
			return nil
		case <-ch:
			ssid, err := GetSSID()
			if err != nil {
				c.logger.Printf("failed to get SSID: %s", err)
				continue
			}
			switch {
			case c.isStartNewTimeEntry(ssid):
				c.started = true
				if err := c.startTimeEntry(); err != nil {
					return err
				}
			case c.isEndOfTimeEntry(ssid):
				ctx, cancel := context.WithCancel(context.Background())

				closeCh := make(chan *Report)
				srv, err := NewServer(c.logger, closeCh)
				if err != nil {
					return err
				}
				srv.Start(ctx)

				open.Run("http://127.0.0.1:8080/report")

				// stop server when report submitted
				c.logger.Println("waiting for submitting report...")
				report := <-closeCh
				cancel()

				c.started = false
				if err := c.stopTimeEntry(report); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (c *Client) isStartNewTimeEntry(ssid string) bool {
	if !(c.ssid == ssid && !c.started) {
		return false
	}
	err := waitUntilReconnection()
	if err != nil {
		c.logger.Println(err)
		return false
	}
	return true
}

func (c *Client) isEndOfTimeEntry(ssid string) bool {
	if !(c.ssid != ssid && c.started) {
		return false
	}
	err := waitUntilReconnection()
	if err != nil {
		c.logger.Println(err)
		return false
	}
	return true
}

func (c *Client) startTimeEntry() error {
	c.logger.Println("start new time entry")
	te, err := c.toggl.Start()
	if err != nil {
		return err
	}
	c.timeEntry = &te
	return nil
}

func (c *Client) stopTimeEntry(report *Report) error {
	c.logger.Println("stop time entry")
	_, err := c.toggl.Stop(*c.timeEntry, report)
	c.timeEntry = nil
	return err
}

func waitUntilReconnection() error {
	for {
		_, err := net.LookupHost("toggl.com")
		if err == nil {
			return nil
		}
		if _, ok := err.(*net.DNSError); !ok {
			return err
		}
	}
	return nil
}
