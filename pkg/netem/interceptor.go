// SPDX-FileCopyrightText: 2026 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package netem

import (
	"maps"

	"github.com/pion/interceptor"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
)

// Interceptor emulates a lossy network link, provide a new model or reuse one
// for each direction in accordance with your use case.
type Interceptor struct {
	interceptor.NoOp
	rtpUpModel    LossModel
	rtpDownModel  LossModel
	rtcpUpModel   LossModel
	rtcpDownModel LossModel
}

// NewInterceptor returns a new InterceptorFactory.
func NewInterceptor(opts ...InterceptorOption) (*InterceptorFactory, error) {
	return &InterceptorFactory{opts: opts}, nil
}

// BindLocalStream implements interceptor.Interceptor.
func (i *Interceptor) BindLocalStream(
	_ *interceptor.StreamInfo,
	writer interceptor.RTPWriter,
) interceptor.RTPWriter {
	return interceptor.RTPWriterFunc(func(
		header *rtp.Header,
		payload []byte,
		attributes interceptor.Attributes,
	) (int, error) {
		// A dropped packet must look like a successful send to us.
		if i.rtpUpModel.ShouldDrop() {
			return header.MarshalSize() + len(payload), nil
		}

		return writer.Write(header, payload, attributes)
	})
}

// BindRemoteStream implements interceptor.Interceptor.
func (i *Interceptor) BindRemoteStream(
	_ *interceptor.StreamInfo,
	reader interceptor.RTPReader,
) interceptor.RTPReader {
	return interceptor.RTPReaderFunc(func(
		buf []byte,
		attributes interceptor.Attributes,
	) (int, interceptor.Attributes, error) {
		for {
			// Copy per attempt, a dropped read caches into the map it is given.
			n, attr, err := reader.Read(buf, maps.Clone(attributes))
			if err != nil {
				return 0, nil, err
			}
			// wait for a successful packet and return it
			if !i.rtpDownModel.ShouldDrop() {
				return n, attr, nil
			}
		}
	})
}

// BindRTCPWriter implements interceptor.Interceptor.
func (i *Interceptor) BindRTCPWriter(writer interceptor.RTCPWriter) interceptor.RTCPWriter {
	return interceptor.RTCPWriterFunc(func(
		pkts []rtcp.Packet,
		attributes interceptor.Attributes,
	) (int, error) {
		if i.rtcpUpModel.ShouldDrop() {
			n := 0
			for _, pkt := range pkts {
				n += pkt.MarshalSize()
			}

			return n, nil // simulate a successful write
		}

		return writer.Write(pkts, attributes)
	})
}

// BindRTCPReader implements interceptor.Interceptor.
func (i *Interceptor) BindRTCPReader(reader interceptor.RTCPReader) interceptor.RTCPReader {
	return interceptor.RTCPReaderFunc(func(
		buf []byte,
		attributes interceptor.Attributes,
	) (int, interceptor.Attributes, error) {
		for {
			// Copy per attempt, a dropped read caches into the map it is given.
			n, attr, err := reader.Read(buf, maps.Clone(attributes))
			if err != nil {
				return 0, nil, err
			}
			if !i.rtcpDownModel.ShouldDrop() {
				return n, attr, nil
			}
		}
	})
}
