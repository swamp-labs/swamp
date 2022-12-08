package httpreq

import (
	"crypto/tls"
	"net"
	"net/http/httptrace"
	"sync/atomic"
	"time"
)

// Sample represents calculated duration of an http request trace
type Sample struct {
	EndTime time.Time

	// Total connect time (Connecting + TLSHandshaking)
	ConnDuration time.Duration

	// Total request duration, excluding DNS lookup and connect time.
	ReqDuration    time.Duration
	DNSDuration    time.Duration // DNS lookup
	WaitingConn    time.Duration // Waiting to acquire a connection.
	Connecting     time.Duration // Connecting to remote host.
	TLSHandshaking time.Duration // Executing TLS handshake.
	Sending        time.Duration // Writing request.
	WaitingResp    time.Duration // Waiting for first byte.
	Receiving      time.Duration // Receiving response.

	// Detailed connection information.
	ConnReused     bool
	ConnRemoteAddr net.Addr
}

// A Trace represents detailed information about a http request. Info are returned using
// httptrace package
type Trace struct {
	getConn              int64
	connectStart         int64
	connectDone          int64
	tlsHandshakeStart    int64
	tlsHandshakeDone     int64
	dnsStart             int64
	dnsDone              int64
	gotConn              int64
	wroteRequest         int64
	gotFirstResponseByte int64
	connReused           bool
	connRemoteAddr       net.Addr
}

// Trace associates http hooks with functions created in this package
func (t *Trace) Trace() *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn:              t.GetConn,
		ConnectStart:         t.ConnectStart,
		ConnectDone:          t.ConnectDone,
		TLSHandshakeStart:    t.TLSHandshakeStart,
		TLSHandshakeDone:     t.TLSHandshakeDone,
		GotConn:              t.GotConn,
		WroteRequest:         t.WroteRequest,
		GotFirstResponseByte: t.GotFirstResponseByte,
		DNSStart:             t.DNSStart,
		DNSDone:              t.DNSDone,
	}
}

func now() int64 {
	return time.Now().UnixNano()
}

// GetConn is called before a connection is created or
// retrieved from an idle pool. The hostPort is the
// "host:port" of the target or proxy. GetConn is called even
// if there's already an idle cached connection available.
func (t *Trace) GetConn(hostPort string) {
	t.getConn = now()
}

// ConnectStart is called when a new connection's Dial begins.
// If net.Dialer.DualStack (IPv6 "Happy Eyeballs") support is
// enabled, this may be called multiple times.
func (t *Trace) ConnectStart(network, addr string) {
	// CompareAndSwapInt64 ensures only the first connectStart is used
	// if the func is called multiple times
	atomic.CompareAndSwapInt64(&t.connectStart, 0, now())

}

// GotConn is called after a successful connection is
// obtained. There is no hook for failure to obtain a
// connection; instead, use the error from
// Transport.RoundTrip.
func (t *Trace) GotConn(info httptrace.GotConnInfo) {
	t.gotConn = now()
	t.connReused = info.Reused
	t.connRemoteAddr = info.Conn.RemoteAddr()
	_, isTLS := info.Conn.(*tls.Conn)

	// If connection is reused we need to attribute a value
	// for connection attribute & tlsHandshake
	if info.Reused {
		atomic.SwapInt64(&t.connectStart, now())
		atomic.SwapInt64(&t.connectDone, now())
		if isTLS {
			atomic.SwapInt64(&t.tlsHandshakeStart, now())
			atomic.SwapInt64(&t.tlsHandshakeDone, now())
		}
	}

}

// ConnectDone is called when a new connection's Dial
// completes. The provided err indicates whether the
// connection completed successfully.
// If net.Dialer.DualStack ("Happy Eyeballs") support is
// enabled, this may be called multiple times.
func (t *Trace) ConnectDone(network, addr string, err error) {
	if err == nil {
		atomic.CompareAndSwapInt64(&t.connectDone, 0, now())
	}
}

// DNSStart is called when a DNS lookup begins.
func (t *Trace) DNSStart(dns httptrace.DNSStartInfo) {
	t.dnsStart = now()
}

// DNSDone is called when a DNS lookup ends.
func (t *Trace) DNSDone(httptrace.DNSDoneInfo) {
	t.dnsDone = now()
}

// TLSHandshakeStart is called when the TLS handshake is started. When
// connecting to an HTTPS site via an HTTP proxy, the handshake happens
// after the CONNECT request is processed by the proxy.
func (t *Trace) TLSHandshakeStart() {
	atomic.CompareAndSwapInt64(&t.tlsHandshakeStart, 0, now())
}

// TLSHandshakeDone is called after the TLS handshake with either the
// successful handshake's connection state, or a non-nil error on handshake
// failure.
func (t *Trace) TLSHandshakeDone(state tls.ConnectionState, err error) {
	if err == nil {
		atomic.CompareAndSwapInt64(&t.tlsHandshakeDone, 0, now())
	}
}

// WroteRequest is called with the result of writing the
// request and any body. It may be called multiple times
// in the case of retried requests.
func (t *Trace) WroteRequest(info httptrace.WroteRequestInfo) {
	if info.Err == nil {
		t.wroteRequest = now()
	}
}

// GotFirstResponseByte is called when the first byte of the response
// headers is available.
func (t *Trace) GotFirstResponseByte() {
	t.gotFirstResponseByte = now()
}

// Done function calculates duration variables to generate a metric sample
func (t *Trace) Done() *Sample {

	s := Sample{}

	s.ConnReused = t.connReused
	s.ConnRemoteAddr = t.connRemoteAddr

	if t.gotConn != 0 && t.getConn != 0 && t.gotConn > t.getConn {
		s.WaitingConn = time.Duration(t.gotConn - t.getConn)
	}
	if t.connectDone != 0 && t.connectStart != 0 {
		s.Connecting = time.Duration(t.connectDone - t.connectStart)
	}
	if t.tlsHandshakeDone != 0 && t.tlsHandshakeStart != 0 {
		s.TLSHandshaking = time.Duration(t.tlsHandshakeDone - t.tlsHandshakeStart)
	}
	if t.dnsStart != 0 && t.dnsDone != 0 {
		s.DNSDuration = time.Duration(t.dnsDone - t.dnsStart)
	}
	if t.wroteRequest != 0 {
		if t.tlsHandshakeDone != 0 {
			// If the request was sent over TLS, we need to use
			// TLS Handshake Done time to calculate sending duration
			s.Sending = time.Duration(t.wroteRequest - t.tlsHandshakeDone)

		} else if t.connectDone != 0 {
			// Otherwise, use the end of the normal connection
			s.Sending = time.Duration(t.wroteRequest - t.connectDone)

		} else {
			// Finally, this handles the strange HTTP/2 case where the GotConn() hook
			// gets called first, but with Reused=false
			s.Sending = time.Duration(t.wroteRequest - t.gotConn)
		}
		if t.gotFirstResponseByte != 0 {
			// We started receiving at least some response back

			if t.gotFirstResponseByte > t.wroteRequest {
				// For some requests, especially HTTP/2, the server starts responding before the
				// client has finished sending the full request
				s.WaitingResp = time.Duration(t.gotFirstResponseByte - t.wroteRequest)
			}
		} else {
			// The server never responded to our request
			s.WaitingResp = time.Now().Sub(time.Unix(0, t.wroteRequest))
		}
		if t.gotFirstResponseByte != 0 {
			// Time elapsed since we receive the first byte
			s.Receiving = time.Now().Sub(time.Unix(0, t.gotFirstResponseByte))
		}
	}
	s.EndTime = time.Now()
	s.ConnDuration = s.Connecting + s.TLSHandshaking
	s.ReqDuration = s.Sending + s.WaitingResp + s.Receiving
	return &s
}
