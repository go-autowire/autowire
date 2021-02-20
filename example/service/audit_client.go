package service

import (
	"github.com/go-autowire/autowire"
	"log"
)

func init() {
	autowire.Autowire(&AuditClient{})
}

type AuditEventSender interface {
	Send(event string)
}

type AuditClient struct {
	Type string
}

func (AuditClient) Send(event string) {
	log.Printf("auditClient Event %s sent", event)
}
