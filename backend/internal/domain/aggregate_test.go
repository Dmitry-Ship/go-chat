package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultGetEvents(t *testing.T) {
	aggregate := aggregate{}

	assert.Equal(t, len(aggregate.GetEvents()), 0)
}

func TestAddEvents(t *testing.T) {
	aggregate := aggregate{}

	event := domainEvent{
		name: "event",
	}

	aggregate.AddEvent(event)

	assert.Equal(t, aggregate.GetEvents()[len(aggregate.GetEvents())-1], event)
}
