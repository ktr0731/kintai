package main

import (
	"errors"
	"math/rand"

	"github.com/gedex/go-toggl/toggl"
)

type Toggler interface {
	Create(*toggl.TimeEntry) (*toggl.TimeEntry, error)
	Start(*toggl.TimeEntry) (*toggl.TimeEntry, error)
	Stop(int) (*toggl.TimeEntry, error)
}

type TogglClient struct {
	client *toggl.Client
}

func NewTogglClient(apiToken string) *TogglClient {
	return &TogglClient{
		client: toggl.NewClient(apiToken),
	}
}

func (c *TogglClient) Create(te *toggl.TimeEntry) (*toggl.TimeEntry, error) {
	return c.client.TimeEntries.Create(te)
}

func (c *TogglClient) Start(te *toggl.TimeEntry) (*toggl.TimeEntry, error) {
	return c.client.TimeEntries.Start(te)
}

func (c *TogglClient) Stop(id int) (*toggl.TimeEntry, error) {
	return c.client.TimeEntries.Stop(id)
}

type MockTogglClient struct {
	te *toggl.TimeEntry
}

func NewMockTogglClient(apiToken string) *MockTogglClient {
	return &MockTogglClient{}
}

func (c *MockTogglClient) Create(te *toggl.TimeEntry) (*toggl.TimeEntry, error) {
	id := rand.Intn(100)
	return &toggl.TimeEntry{ID: id}, nil
}

func (c *MockTogglClient) Start(te *toggl.TimeEntry) (*toggl.TimeEntry, error) {
	c.te = te
	return te, nil
}

func (c *MockTogglClient) Stop(id int) (*toggl.TimeEntry, error) {
	if c.te.ID != id {
		return nil, errors.New("no such id")
	}
	return c.te, nil
}
