package shared

import (
	"github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
	"zenithar.org/go/common/registrar"
)

var (
	// Config contains flags passed to the API
	Config *Flags
	// Log is the API's logrus instance
	Log *logrus.Logger
	// Raven is sentry client
	Raven *raven.Client
	// Registrar is the bean registry
	Registrar registrar.Registrar
)
