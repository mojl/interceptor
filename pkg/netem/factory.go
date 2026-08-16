// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package netem

import (
	"github.com/pion/interceptor"
)

// InterceptorFactory is an interceptor.Factory for netem Interceptors.
type InterceptorFactory struct {
	opts []InterceptorOption
}

// NewInterceptor constructs a new netem Interceptor, lossless by default.
func (f *InterceptorFactory) NewInterceptor(_ string) (interceptor.Interceptor, error) {
	netemInterceptor := &Interceptor{
		NoOp:          interceptor.NoOp{},
		rtpUpModel:    lossless{},
		rtpDownModel:  lossless{},
		rtcpUpModel:   lossless{},
		rtcpDownModel: lossless{},
	}

	for _, opt := range f.opts {
		if err := opt(netemInterceptor); err != nil {
			return nil, err
		}
	}

	return netemInterceptor, nil
}
