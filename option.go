package certsurfer

import (
	"context"
	"net/http"
	"time"
)

const calidogFullStream = "wss://certstream.calidog.io/full-stream"
const calidogDomainStream = "wss://certstream.calidog.io/domains-only"

type config struct {
	endpoint     string
	ctx          context.Context
	client       *http.Client
	header       http.Header
	pingInterval time.Duration
	pingTimeout  time.Duration
	readTimeout  time.Duration
}

type Option func(opt config) config

func defaultConfig() config {
	c := *http.DefaultClient
	c.Timeout = time.Second * 5
	return config{
		endpoint:     calidogFullStream,
		ctx:          context.Background(),
		client:       &c,
		header:       nil,
		pingInterval: time.Second * 60,
		pingTimeout:  time.Second * 30,
		readTimeout:  time.Second * 30,
	}
}

func applyOpts(opt config, conf []Option) config {
	for _, c := range conf {
		opt = c(opt)
	}
	return opt
}

//WithFullStream uses the full stream endpoint, which includes certificate data
func WithFullStream() Option {
	return func(opt config) config {
		opt.endpoint = calidogFullStream
		return opt
	}
}

//WithDomainStream uses the domain stream endpoint that only includes partial data
func WithDomainStream() Option {
	return func(opt config) config {
		opt.endpoint = calidogDomainStream
		return opt
	}
}

//WithEndpoint uses a custom certstream url as endpoint
//this will work with both domains and full cert endpoints
func WithEndpoint(url string) Option {
	return func(opt config) config {
		opt.endpoint = url
		return opt
	}
}

//WithContext uses a context to cancel the stream
//setting nil will fall back to background context
func WithContext(ctx context.Context) Option {
	return func(opt config) config {
		if ctx == nil {
			opt.ctx = context.Background()
		} else {
			opt.ctx = ctx
		}
		return opt
	}
}

//WithClient uses a custom http client
func WithClient(client *http.Client) Option {
	return func(opt config) config {
		opt.client = client
		return opt
	}
}

//WithHeader adds custom headers to the request
func WithHeader(header http.Header) Option {
	return func(opt config) config {
		opt.header = header
		return opt
	}
}

//WithPingInterval sets the interval between pings
func WithPingInterval(interval time.Duration) Option {
	return func(opt config) config {
		opt.pingInterval = interval
		return opt
	}
}

//WithPingTimeout sets the timeout for reading pings
//if the timeout has been exceeded, the stream will be closed
func WithPingTimeout(timeout time.Duration) Option {
	return func(opt config) config {
		opt.pingTimeout = timeout
		return opt
	}
}

//WithReadTimeout sets the timeout for reading data
func WithReadTimeout(timeout time.Duration) Option {
	return func(opt config) config {
		opt.readTimeout = timeout
		return opt
	}
}
