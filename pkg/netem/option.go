// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package netem

// InterceptorOption is a configuration option for netem interceptors.
type InterceptorOption func(*Interceptor) error

// RTPUpModel uses a provided model for uplink packet loss.
func RTPUpModel(m LossModel) InterceptorOption {
	return func(i *Interceptor) error {
		if m == nil {
			return ErrNilModel
		}
		i.rtpUpModel = m

		return nil
	}
}

// RTPDownModel uses a provided model for downlink packet loss.
func RTPDownModel(m LossModel) InterceptorOption {
	return func(i *Interceptor) error {
		if m == nil {
			return ErrNilModel
		}
		i.rtpDownModel = m

		return nil
	}
}

// RTCPUpModel uses a provided model for uplink packet loss.
func RTCPUpModel(m LossModel) InterceptorOption {
	return func(i *Interceptor) error {
		if m == nil {
			return ErrNilModel
		}
		i.rtcpUpModel = m

		return nil
	}
}

// RTCPDownModel uses a provided model for downlink packet loss.
func RTCPDownModel(m LossModel) InterceptorOption {
	return func(i *Interceptor) error {
		if m == nil {
			return ErrNilModel
		}
		i.rtcpDownModel = m

		return nil
	}
}
