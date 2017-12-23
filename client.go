package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gedex/go-toggl/toggl"
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
	ch := time.Tick(c.duration)
	for range ch {
		ssid, err := GetSSID()
		if err != nil {
			msg := fmt.Sprintf("failed to get SSID: %s", err)
			c.logger.Println(msg)
			return errors.New(msg)
		}
		switch {
		case c.isStartNewTimeEntry(ssid):
			c.started = true
			if err := c.startTimeEntry(); err != nil {
				return err
			}
			c.logger.Println("start new time entry")
		case c.isEndOfTimeEntry(ssid):
			c.started = false
			if err := c.stopTimeEntry(); err != nil {
				return err
			}
			c.logger.Println("stop time entry")
		}
	}
	c.logger.Println("pouring finished")
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
