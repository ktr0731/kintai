package main

import (
	"context"
	"log"
	"time"

	"github.com/gedex/go-toggl/toggl"
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
		toggl:    NewMockTogglClient(apiToken),
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
				c.logger.Println("start new time entry")
			case c.isEndOfTimeEntry(ssid):
				ctx, cancel := context.WithCancel(context.Background())

				closeCh := make(chan struct{})
				srv, err := NewServer(c.logger, closeCh)
				if err != nil {
					return err
				}
				srv.Start(ctx)

				open.Run("http://127.0.0.1:8080/report")

				// stop server when report submitted
				c.logger.Println("waiting for submitting report...")
				<-closeCh
				cancel()

				c.started = false
				if err := c.stopTimeEntry(); err != nil {
					return err
				}
				c.logger.Println("stop time entry")
			}
		}
	}
	return nil
}

func (c *Client) isStartNewTimeEntry(ssid string) bool {
	return c.ssid == ssid && !c.started
}

// TODO: 時間差で終える
func (c *Client) isEndOfTimeEntry(ssid string) bool {
	return c.ssid != ssid && c.started
}

func (c *Client) startTimeEntry() error {
	te, err := c.toggl.Create(&toggl.TimeEntry{})
	if err != nil {
		return err
	}
	te, err = c.toggl.Start(te)
	if err != nil {
		return err
	}
	c.timeEntry = te
	return nil
}

func (c *Client) stopTimeEntry() error {
	_, err := c.toggl.Stop(c.timeEntry.ID)
	c.timeEntry = nil
	return err
}
