// Package errors defines application error codes for the apiserver module.
//
// Error code ranges:
//   - 10000-10999: Auth module
//   - 11000-11999: OAuth module
//   - 20000-20999: Book module
//   - 30000-30999: Image module
//   - 40000-40999: AccountSecret module
//   - 50000-50999: Email module
package errors

import (
	pkgerr "github.com/mgcis-cn/ibookfs/pkg/errors"
)

// Re-export pkg/errors functions for convenience within apiserver.
var (
	FromError   = pkgerr.FromError
	IsCode      = pkgerr.IsCode
	WithMessage = pkgerr.WithMessage
	WithCause   = pkgerr.WithCause
	Wrap        = pkgerr.Wrap
)
