package certsurfer

import (
	"context"
	"errors"
	"github.com/goccy/go-json"
	"nhooyr.io/websocket"
	"time"
)

//Open a connection and feed the data into the given channel
//this function will block until the connection gets closed, or an error occurs, or cancelled by context
//Option can be used to configure the connection
//The channel should always have enough buffer to receive the data, otherwise the data will be dropped
func Open(stream chan<- MixedData, options ...Option) error {
	opt := applyOpts(defaultConfig(), options)

	fn := func(data MixedData) {
		select {
		case stream <- data:
		default:
		}
	}

	err := open(fn, opt)
	if err != nil {
		return err
	}
	return nil
}

//OpenWithCallback open a connection and feed the data into the given callback
//this function will block until the connection gets closed, or an error occurs, or cancelled by context
//Option can be used to configure the connection
//The callback must not block as the callback gets called in the main loop,
//if the callback blocks, this function cannot exit until the callback returns
//the callback should spawn a goroutine if the callback needs to block
func OpenWithCallback(f func(MixedData), options ...Option) error {
	opt := applyOpts(defaultConfig(), options)
	err := open(f, opt)
	if err != nil {
		return err
	}
	return nil
}

func open(f func(MixedData), opt config) error {
	con, _, err := websocket.Dial(context.Background(), opt.endpoint, &websocket.DialOptions{
		HTTPClient: opt.client,
		HTTPHeader: opt.header,
	})

	defer func() {
		if err == nil || errors.Is(err, opt.ctx.Err()) {
			_ = con.Close(websocket.StatusNormalClosure, "")
		} else if _, ok := err.(*json.UnmarshalTypeError); ok {
			_ = con.Close(websocket.StatusInternalError, "invalid json error")
		} else {
			_ = con.Close(websocket.StatusInternalError, "")
		}
	}()

	if err != nil {
		return dialError{err: err}
	}

	tick := time.NewTicker(opt.pingInterval)
	defer tick.Stop()
	var pingStat PingStatus
	pingC := make(chan PingStatus, 1)
	for {
		select {
		case <-opt.ctx.Done():
			return nil
		case <-tick.C:
			go func() { //uncertain if this is even necessary or helps
				ctx, cancel := context.WithTimeout(opt.ctx, opt.pingTimeout)
				defer cancel()
				start := time.Now()
				_ = con.Ping(ctx) //errors will get pass to next con.read call
				pingC <- PingStatus{
					Latency: time.Since(start),
					Time:    time.Now(),
				}
			}()
		case pingStat = <-pingC:
		default:
			var bs []byte
			callCtx(func(ctx context.Context) {
				_, bs, err = con.Read(ctx)
			}, opt.readTimeout)
			if err != nil && !errors.Is(err, opt.ctx.Err()) {
				return readError{err: err}
			}
			var v MixedData
			err = json.Unmarshal(bs, &v)
			if err != nil {
				return jsonError{err: err}
			}
			v.pingStat = pingStat
			f(v)
		}
	}
}

func callCtx(f func(ctx context.Context), duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	f(ctx)
}
