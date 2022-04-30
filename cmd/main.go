package main

import (
	"context"
	"fmt"
	"github.com/thunder33345/certsurfer"
	"os"
	"os/signal"
	"strings"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig)
		<-sig
		fmt.Printf("signal: ok got cancel\n")
		cancel()
	}()

	fmt.Printf("main: running\n")
	start := time.Now()
	err := do(ctx)
	end := time.Now()
	fmt.Printf("open: done")
	if err != nil {
		fmt.Printf("open: got error: %v\n", err)
	}
	fmt.Printf("open: lasted for: %v\n", end.Sub(start))

	fmt.Printf("main: waiting..\n")
	time.Sleep(time.Second)
	fmt.Printf("main: quit..\n")
}

func do(ctx context.Context) error {
	ch := make(chan certsurfer.MixedData, 10)
	go func() {
		start := time.Now()
		var last certsurfer.PingStatus
	out:
		for {
			select {
			case d, ok := <-ch:
				if !ok {
					break out
				}
				if h, ok := d.AsHeartbeat(); ok {
					fmt.Printf("heartbeat latency: %g\n", h.Timestamp)
				} else if u, ok := d.AsCertificate(); ok {
					fmt.Printf("new cert(domain): %s\n", strings.Join(u.Data.LeafCert.AllDomains, ", "))
				} else if s, ok := d.AsDomain(); ok {
					fmt.Printf("new domains: %s\n", strings.Join(s.Data, ", "))
				} else if typ, data, ok := d.AsUnknown(); ok {
					fmt.Printf("unknown (%s): %s\n", typ, string(data))
				}

				if p, ok := d.Ping(); ok {
					if p != last {
						fmt.Printf("Ping: Latency: %v, Last updated: %v, Since start: %v\n", p.Latency, time.Since(p.Time), time.Since(start))
						last = p
					}
				}
			}
		}
		fmt.Printf("routine: Channel Closed\n")
	}()
	defer close(ch)
	err := certsurfer.Open(ch, certsurfer.WithContext(ctx), certsurfer.WithDomainStream())
	return err
}
