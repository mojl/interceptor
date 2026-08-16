// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package netem

import (
	"fmt"
	"math"
	"math/rand/v2"
	"sync"
)

// LossModel decides whether the next datagram/packet should be lost.
type LossModel interface {
	// ShouldDrop returns true if the current packet should be dropped.
	ShouldDrop() bool
}

// lossless is the NoOp model.
type lossless struct{}

func (lossless) ShouldDrop() bool { return false }

// Random drops each packet independently with a probability rate.
type Random struct {
	rate float64

	rand *rand.Rand
	m    sync.Mutex
}

// NewRandom constructs a Random packet loss model. rate must be between 0 and 1.
func NewRandom(rate float64, seed uint64) (*Random, error) {
	if err := checkProbability("rate", rate); err != nil {
		return nil, err
	}

	// nolint:gosec for all rand operations since there is no security risk here.
	return &Random{rate: rate, rand: rand.New(rand.NewPCG(seed, ^seed))}, nil //nolint:gosec
}

// ShouldDrop drops packets randomly.
func (r *Random) ShouldDrop() bool {
	r.m.Lock()
	defer r.m.Unlock()

	return r.rand.Float64() < r.rate
}

// GilbertElliott uses a Gilbert-Elliott (burst loss) model.
type GilbertElliott struct {
	p        float64 // probability of starting bad (lossy) state.
	r        float64 // probability of exiting bad state.
	badRate  float64 // loss probability in bad state.
	goodRate float64 // loss probability in good state.
	bad      bool

	rand *rand.Rand
	m    sync.Mutex
}

// NewGilbertElliott constructs a new GilbertElliott packet loss model (burst
// mode). p, r, goodRate and badRate must all be between 0 and 1.
func NewGilbertElliott(p float64, r float64, goodRate float64, badRate float64, seed uint64) (*GilbertElliott, error) {
	params := map[string]float64{"p": p, "r": r, "goodRate": goodRate, "badRate": badRate}
	for name, prob := range params {
		if err := checkProbability(name, prob); err != nil {
			return nil, err
		}
	}

	return &GilbertElliott{
		p:        p,
		r:        r,
		goodRate: goodRate,
		badRate:  badRate,
		rand:     rand.New(rand.NewPCG(seed, ^seed)), //nolint:gosec
	}, nil
}

// ShouldDrop drops packets when the GE models tell us to do.
func (ge *GilbertElliott) ShouldDrop() bool {
	ge.m.Lock()
	defer ge.m.Unlock()

	ge.advance()

	rate := ge.goodRate
	if ge.bad {
		rate = ge.badRate
	}

	return ge.rand.Float64() < rate
}

// advance the state machine.
func (ge *GilbertElliott) advance() {
	if ge.bad {
		ge.bad = ge.rand.Float64() >= ge.r
	} else {
		ge.bad = ge.rand.Float64() < ge.p
	}
}

func checkProbability(name string, p float64) error {
	if p < 0 || p > 1 || math.IsNaN(p) {
		return fmt.Errorf("%w: %s = %v", ErrInvalidProbability, name, p)
	}

	return nil
}
