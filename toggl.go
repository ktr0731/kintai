package main

import (
	"math/rand"

	"github.com/jason0x43/go-toggl"
)

type Toggler interface {
	Start() (toggl.TimeEntry, error)
	Stop(toggl.TimeEntry, *Report) (toggl.TimeEntry, error)
}

type TogglClient struct {
	client toggl.Session
}

func NewTogglClient(apiToken string) *TogglClient {
	return &TogglClient{
		client: toggl.OpenSession(apiToken),
	}
}

func (c *TogglClient) Start() (toggl.TimeEntry, error) {
	return c.client.StartTimeEntry("")
}

func (c *TogglClient) Stop(te toggl.TimeEntry, report *Report) (toggl.TimeEntry, error) {
	te.Description = report.Description
	_, err := c.client.UpdateTimeEntry(te)
	if err != nil {
		return toggl.TimeEntry{}, err
	}
	return c.client.StopTimeEntry(te)
}

type MockTogglClient struct {
	te *toggl.TimeEntry
}

func NewMockTogglClient(apiToken string) *MockTogglClient {
	return &MockTogglClient{}
}

func (c *MockTogglClient) Start() (toggl.TimeEntry, error) {
	id := rand.Intn(100)
	return toggl.TimeEntry{ID: id}, nil
}
