// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package event

import (
	"sort"

	"github.com/gin-gonic/gin"
)

var eventsListeners map[string][]EventListener

// Interface to implment to create a type of event
type EventInterface interface {
	GetParentEventType() string
}

// Type of Event to handle actions of resource
type ResourceActionEvent struct {
	Resource interface{}
	Action   string
}

// Type of Event to handle requests events
type RequestEvent struct {
	Context *gin.Context
}

// Event listener to add on list of eventsListeners
type EventListener struct {
	Type     string
	Handler  func(c *gin.Context, e EventInterface) error
	Priority int
}

// Get parent type of the event
func (rae *ResourceActionEvent) GetParentEventType() string {
	return EVENT_RESOURCE_ACTION
}

// Get parent type of the event
func (rae *RequestEvent) GetParentEventType() string {
	return ""
}

// Reset event listeners
func ResetEventListeners() {
	eventsListeners = map[string][]EventListener{}
}

// Register a new listener
func RegisterEventListener(el EventListener) {
	if _, ok := eventsListeners[el.Type]; !ok {
		eventsListeners[el.Type] = []EventListener{}
	}
	eventsListeners[el.Type] = append(eventsListeners[el.Type], el)
}

// Dispatch an event and browse the list of registered event listeners to apply the handler
func DispatchEvent(c *gin.Context, eventType string, event EventInterface) error {
	// dispatch parent event
	if event.GetParentEventType() != "" {
		err := dispatchEvent(c, event.GetParentEventType(), event)
		if err != nil {
			return err
		}
	}
	return dispatchEvent(c, eventType, event)
}

func dispatchEvent(c *gin.Context, eventType string, event EventInterface) error {
	if _, ok := eventsListeners[eventType]; !ok {
		return nil
	}
	sort.SliceStable(eventsListeners[eventType], func(i, j int) bool {
		return eventsListeners[eventType][i].Priority > eventsListeners[eventType][j].Priority
	})
	for _, el := range eventsListeners[eventType] {
		err := el.Handler(c, event)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	ResetEventListeners()
}
