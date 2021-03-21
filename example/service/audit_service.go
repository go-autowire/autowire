package service

import (
	"log"
)

// A EventSender represents interface
type EventSender interface {
	// Send function is simply dispatching event
	Send(event string)
}

// A AuditService represents a named struct
type AuditService struct {
	Type string
}

// Send is dispatching the event
func (AuditService) Send(event string) {
	log.Printf("auditClient Event %s sent", event)
}

// Close method release all occupied resources
func (a AuditService) Close() error {
	log.Println("Closing AuditService")
	return nil
}
