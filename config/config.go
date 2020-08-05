package config

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

type Env struct {
	AWSClient *session.Session
}

func NewEnv(client *session.Session) *Env {
	return &Env{AWSClient: client}
}
