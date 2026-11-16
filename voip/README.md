# GREEN-API VoIP for Go

The `voip` package contains REST call commands (`calls.go`), reconnecting signaling (`signaling.go`), SDP and ICE mapping (`rtc.go`), audio contracts (`audio.go`), and a Pion media adapter (`rtc_adapter.go`). Methods named `Async` wait for their result; use a goroutine when concurrent execution is needed.

## Calling the API

```go
callsAPI := voip.New(voip.Config{
    Host: "https://api.green-api.com",
    IDInstance: os.Getenv("GREEN_API_ID"),
    APITokenInstance: os.Getenv("GREEN_API_TOKEN"),
})
state, err := callsAPI.GetStateAsync(ctx)
```

`Config.HTTPClient` and `Config.Logger` allow transport injection and logging configuration. The supplied HTTP client is copied and redirects are disabled. Without a custom client the REST transport has five-second connection/TLS/header and idle socket-read timeouts, no total timeout, environment proxy support, and no automatic command retries. All 2xx responses succeed; 204/empty responses have no result. Unexpected status and malformed JSON return errors independently of the main SDK's HTTP response policy. `RequestError` keeps the original cause available through `errors.Is` and `errors.As`.

Methods: `GetStateAsync`, `GetIceServersAsync`, `DialAsync`, `AcceptAsync`, `RejectAsync`, `HangUpAsync`, `WebsocketUrl`, `Connect`.

## Audio and connection lifecycle

```go
connection := callsAPI.Connect(audioFactory)
defer connection.CloseAsync()

connection.On("state", func(detail any) {
    state := detail.(voip.CallState)
    log.Println(state.State)
})
connection.OnAsync("incoming_call", func(detail any) {
    // Blocking/reentrant work belongs in OnAsync callbacks.
    if err := callsAPI.AcceptAsync(callContext); err != nil { return }
    if err := connection.StartAudioAsync(callContext); err != nil { /* handle */ }
})

openContext, cancel := context.WithTimeout(callContext, 30*time.Second)
defer cancel()
err := connection.OpenAsync(openContext)
```

For an outgoing call use `OpenAsync` → `DialAsync` → `StartAudioAsync`. For an incoming call use `OpenAsync` → incoming event → `AcceptAsync` → `StartAudioAsync`. Neither dial nor accept implicitly starts audio.

`CallAudio` contains `LocalTrack`, `OnRemoteTrack`, and `Close`. `AudioFactory` creates fresh resources for every bridge, including reconnects. The application owns these resources. `LocalAudioTrack.Recv` and the remote attachment callback must honor cancellation. `CallAudio.Close` must release the local source and sink. All attachment tasks are canceled and awaited before audio is disposed.

The adapter accepts `NewPCMTrack(source, codecFactory, channels)`. Frames are interleaved signed 16-bit PCM at 48 kHz, one or two channels, with PTS in samples per channel. The adapter encodes 20 ms Opus packets, paces outgoing audio, reorders incoming RTP and decodes it to PCM. `FrameAudioSink.AttachAsync` delivers frames sequentially and waits for each callback. Format conversion/resampling belongs in the source; arbitrary PyAV file formats are not implemented by this package. Codec instances are created per bridge and direction, never shared.

`NewRTCAdapter(servers)` creates the Pion implementation of `Bridge`. Custom `Bridge` implementations must call their `OnError` handler when the active peer fails, including ICE or media failures; `CallsConnection` uses that signal to close the failed bridge and emit an `error` event.

`opus.New` provides the bundled libopus codec through `layeh.com/gopus`, pinned in `go.mod`. On amd64/386 it needs CGo and a C compiler; other supported gopus targets use system libopus (see that dependency's build requirements). CGo-free builds of the core SDK and signaling work normally. In a CGo-free build `opus.New` returns an explicit error; provide another `CodecFactory` to use PCM audio without it. The pinned gopus embeds libopus 1.1.2: codec updates and platform coverage can be handled in this isolated adapter.

### Call behavior

- States: `idle`, `inc-call`, `out-call`, `on-call`; `CallInfo.Name` is `*string` to distinguish absent/null from an empty name.
- WebSocket backoff: 500 ms, 1 s, 2 s, 4 s, 8 s, then 10 s indefinitely. Every successful connection resets it. No jitter or outgoing-message queue.
- WebSocket close codes 4000–4999 and HTTP handshake status 401/403 permanently stop reconnection. A refused handshake returns `HandshakeRefusedError` with the HTTP response and original dial error, emits a permanent `disconnect`, and sets `Refused()`. Other handshake failures remain retryable. Explicit close stops attempts.
- WebSocket opening timeout only stops the caller's wait, not the background reconnect loop. `CloseAsync` ends that loop.
- Heartbeat uses a 20-second ping interval and 20-second pong wait; handshake timeout is 10 seconds and message size limit is 1 MiB.
- Internal signaling frames are dispatched sequentially. A state transition completes bridge teardown before the next frame is handled.
- Local ICE gathering finishes before sending the offer; gathered candidates are also sent as separate frames afterward. Early remote candidates wait for the answer. End-of-candidates is supported.
- `StartAudioAsync` succeeds when the remote answer and buffered ICE candidates have been applied, not when audio starts arriving. There is no built-in answer timeout. Canceling its context during startup aborts negotiation; after success the bridge lives until stop, close, or a call event.
- Duplicate bridge startup is rejected. Generation checks discard stale factory results and track callbacks.
- A transient disconnect disposes the old bridge. A new server call state can recreate the audio and peer without dialing/accepting again, including after another transient disconnect during resume. Permanent refusal prevents resume.
- Repeated answers are ignored after the first answer is applied. A failed ICE or DTLS connection reports `error` and disposes the current bridge; late errors from an older bridge do not affect a replacement.
- Errors during pending negotiation fail startup and clean up; a generic server error after the answer emits an error event without automatically destroying the bridge.
- `StopAudioAsync` sends `stop` when a peer exists and disposes audio. `CloseAsync` disposes the bridge and socket without sending `stop`. Neither sends `callsHangUp`: use `HangUpAsync` explicitly. A stop write is bounded by ten seconds so cleanup cannot hang indefinitely.

The adapter negotiates Opus audio only. It does not flush an incomplete 20 ms PCM frame on EOF; callers that need the final samples sent must pad the source frame. Local Go integration tests verify bidirectional non-silent PCM audio over WebRTC.

### Events

| Event | Payload |
| --- | --- |
| `connect` | `nil` |
| `disconnect` | `DisconnectDetail` |
| `state` | `CallState` |
| `incoming_call` | `CallInfo` |
| `remote_track` | `RemoteAudioTrack` |
| `end_call` | `EndCallDetail` (`Reason`, optional `Cause`) |
| `error` | `ErrorDetail` |

`On` callbacks run synchronously: keep them short and do not call connection lifecycle methods from them. Use `OnAsync` for blocking or reentrant work. Async callbacks can finish out of order and must manage their own long-running work. `On`/`OnAsync` return an unsubscribe function; invoke it or pass it to `Off`. Callback panics are logged without terminating signaling. `NewFrameAudioSink` accepts an optional error handler; otherwise it logs consumer errors.

## Runnable examples

Set `GREEN_API_ID`, `GREEN_API_TOKEN`, optionally `GREEN_API_HOST`. The host must serve the call endpoints and `callsRtc` WebSocket. Provide a 48 kHz signed 16-bit PCM WAV (mono or stereo).

```sh
go run ./examples/voip/headless_call 77001234567 ./audio.wav
go run ./examples/voip/incoming_call ./audio.wav
```

The examples send audio from the file and print received frame counts, with no audio hardware. They reopen the file for every bridge. Ctrl+C closes local call resources; applications call `HangUpAsync` explicitly when they intend to end the server call. EOF stops outgoing audio; it does not hang up the call.

## Verification

```sh
go test -race ./...
CGO_ENABLED=0 go test ./...
```

Tests cover REST/status/body policy, candidate ordering, backoff/permanent refusal, ordered dispatch, pending negotiation, close/stop/end semantics, cancellation, fresh audio on resume, real peer audio, and the complete REST → WebSocket → WebRTC path against a local server.

Local tests need no account credentials. Actual account calls and TURN deployment require an authorized GREEN-API instance and network access; local tests do not establish those deployment properties.
