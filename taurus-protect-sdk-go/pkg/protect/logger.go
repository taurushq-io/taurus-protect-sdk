package protect

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/service"
)

// Logging types are aliases to the service package, matching how errors.go aliases
// service.APIError: the interface has to live where the services can see it, and
// protect imports service, so the reverse would cycle.
type (
	// Logger receives SDK diagnostic events. ctx is first so a consumer whose logger
	// derives fields from it (a correlation id) keeps them.
	Logger = service.Logger
	// Field is one structured key/value pair. Metadata only — never payloads,
	// tokens or key material.
	Field = service.Field
	// NoopLogger discards everything and is the default.
	NoopLogger = service.NoopLogger
)
