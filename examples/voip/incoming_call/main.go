// Answer one incoming call with a WAV source.
package main

import (
	"context"
	"fmt"
	"github.com/green-api/whatsapp-api-client-golang-v2/examples/voip/internal/example"
	"github.com/green-api/whatsapp-api-client-golang-v2/voip"
	"os"
	"os/signal"
	"time"
)

func run() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("usage: incoming_call <audio.wav>")
	}

	config, e := example.Config()

	if e != nil {
		return e
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()

	api := voip.New(config)
	calls := api.Connect(example.AudioFactory(os.Args[1]))

	defer calls.CloseAsync()

	incoming := make(chan voip.CallInfo, 1)

	calls.On("incoming_call", func(d any) {
		select {
		case incoming <- d.(voip.CallInfo):

		default:
		}
	})

	calls.On("end_call", func(any) { cancel() })
	calls.On("error", func(d any) { fmt.Fprintln(os.Stderr, d) })
	openCtx, stop := context.WithTimeout(ctx, 30*time.Second)

	defer stop()

	if e = calls.OpenAsync(openCtx); e != nil {
		return e
	}

	select {
	case <-incoming:

	case <-ctx.Done():
		return nil
	}

	if state := calls.State(); state == nil || state.State != "inc-call" {
		return nil
	}

	if e = api.AcceptAsync(ctx); e != nil {
		return e
	}

	if e = calls.StartAudioAsync(ctx); e != nil && ctx.Err() == nil {
		return e
	}

	<-ctx.Done()

	return nil
}

func main() {
	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
