// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// SPDX-License-Identifier: MIT

package dataminer

import (
	"errors"
)

const (
	// Symbl proxy valid ports
	DefaultStartPort int = 30000
	DefaultEndPort   int = 34999

	// server side events ports for client notifications will add 10000 to the configured or
	// DefaultStartPort to start numbering from there
	DefaultNotificationPortOffset int    = 10000
	DefaultNotificationPath       string = "notifications"
)

var (
	// ErrInvalidInput required input was not found
	ErrInvalidInput = errors.New("required input was not found")
)
