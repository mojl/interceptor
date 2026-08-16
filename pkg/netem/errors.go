// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package netem

import "errors"

var (
	// ErrNilModel is returned when a nil LossModel is passed to an option.
	ErrNilModel = errors.New("loss model cannot be nil")
	// ErrInvalidProbability is returned when a probability is outside 0 and 1.
	ErrInvalidProbability = errors.New("probability must be in between 0 and 1")
)
