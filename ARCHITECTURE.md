# Architecture — whatsapp-api-client-golang-v2

This document describes the current implementation. Mandatory standards for new and materially changed code are in [CONTRIBUTING.md](CONTRIBUTING.md); agents start with [AGENTS.md](AGENTS.md). The implementation does not yet meet every contribution standard.

---

## 1) Scope and entry points

This is a public Go SDK for GREEN-API's WhatsApp service, not a deployed application. Applications construct `greenapi.GreenAPI` and call category methods. Call signaling and media use the separate `voip` package; the optional `opus` package supplies a codec. Applications supply credentials, business processing, and audio sources/sinks. GREEN-API owns the account service, remote notification queue, call endpoints, signaling service, and ICE server configuration.

`go.mod` declares `github.com/green-api/whatsapp-api-client-golang-v2`, Go 1.23.0 and toolchain go1.23.2. REST uses `fasthttp`; signaling uses `coder/websocket`; media uses Pion WebRTC v4; the bundled codec uses `layeh.com/gopus`. Version and target changes can affect CGo and media behavior independently of REST.

Constructing `GreenAPI` merely supplies fields. `voip.New` prepares configuration and an HTTP client; it does not make a call. `Voip.Connect` creates the connection and registers handlers without opening signaling. `OpenAsync` starts the reconnect loop. Methods named `Async` wait for their result; the name does not imply fire-and-forget execution.

There is no root executable, webhook server, database, durable handler loop, microphone/speaker integration, resampling service, or workflow engine. Programs under `examples/` are consumers and can perform real account actions when executed.

## 2) Component ownership and coupling

| Source / symbols | Responsibility and caller | Why changes propagate |
| --- | --- | --- |
| `types.go`: `GreenAPI`, `APIResponse` | Consumer configuration; status/raw-body response envelope | Field or envelope changes affect every REST category and consumer decoding. |
| `categories.go` | Category constructors holding `GreenAPIInterface` | This is the public façade used by root examples and integrations. |
| `account.go`, `contacts.go`, `groups.go`, `statuses.go`, `journals.go`, `queues.go`, `readmark.go`, `service.go` | Endpoint arguments, options, JSON and delegation to `Request` | A changed JSON tag or omitted optional field changes wire behavior without changing a Go signature. |
| `sending.go` | Messages, polls, files, locations, contacts, interactive messages | Shares normal and media/multipart request routes; upload changes can affect several methods. |
| `receiving.go` | One poll, explicit receipt deletion, file download | Consumers decide when to delete; these methods do not run handlers. |
| `request.go` | Options, URL/header construction, multipart encoding, `fasthttp` transport | Host selection, response policy and buffer ownership are shared across categories. |
| `partner.go`, `GreenAPIPartner`, `PartnerRequest` | Partner methods, separate token and fixed partner host | Partner transport is separate code; a root request fix does not automatically fix it. |
| `voip/calls.go`: `Voip` | Context-aware call REST commands and connection construction | Its status/timeout policy differs from root REST. |
| `voip/calls.go`: `CallsConnection` | Signaling interpretation, bridge generation, startup completion, callbacks, teardown | State and bridge lifecycle are coupled; stale callbacks must not affect replacement resources. |
| `voip/signaling.go`: `ReconnectingSocket` | Socket, reconnect loop, heartbeat and event dispatch | Permanent/transient classification determines whether calls can resume. |
| `voip/rtc.go`: `Bridge`, ICE/SDP types | Peer contract and wire mapping | SDP/ICE normalization affects negotiation in both adapter and call coordinator. |
| `voip/rtc_adapter.go`: `RTCAdapter`, `PCMTrack` | Pion peer, Opus transport, outgoing/incoming PCM | Codec, frame size and cancellation changes affect media and shutdown. |
| `voip/audio.go`: `CallAudio`, `FrameAudioSink` | Application audio contracts and sequential frame delivery | Non-cooperative callbacks can prevent teardown from finishing. |
| `opus/opus.go`, `opus/opus_nocgo.go` | Bundled codec and CGo-free error fallback | Codec availability is a target/build property, not proof signaling is unavailable. |

## 3) Root REST request flow

**Entry and sequence:** `GreenAPI.Sending()` → `SendingCategory.SendMessage` → validation/options → `json.Marshal` → `GreenAPIInterface.Request` → `GreenAPI.request` → `fasthttp.Client.Do`. Other categories follow the same broad pattern; inspect the specific wrapper for validation and fields.

`Request` applies `requestOptions` and passes query, MIME/filename, form-data and media-host choices to the transport. The normal URL is `{APIURL}/waInstance{IDInstance}/{APIMethod}/{APITokenInstance}`. Media routing substitutes `MediaURL`. `GetParams` is appended verbatim; it is not always a query string. `PartnerRequest` instead uses the fixed partner URL with a partner token and a client name containing the configured email.

**Output and side effects:** `APIResponse` carries status code, status-message bytes, raw JSON body and a local timestamp. A remote operation may already have happened even if the response is lost. Non-2xx HTTP responses remain response data with a nil Go transport error. Validation, encoding, option and network failures return errors. There is no automatic retry or command queue.

**Ownership and cancellation:** each call creates a `fasthttp.Client`; active timeout fields are not configured. Public category calls have no context parameter. Requests and responses are acquired from pools and released on return. The implementation assigns `resp.Body()` and `resp.Header.StatusMessage()` directly into the envelope without copying. In particular, fasthttp's `Response.Body` contract says its bytes are valid only until release. This is an existing lifetime defect/risk, not a guarantee of independently owned response bytes. Changes to this code must account for buffer reuse; copying only in a consumer after return cannot establish safe ownership during concurrent pool reuse. This documentation does not repair the transport.

**Compatibility and verification:** URL suffix order, option encoding, status-vs-error behavior and response byte ownership affect all consumers. Local HTTP-server tests can capture paths, methods and JSON and return controlled status/body variants; package-local fakes can implement the otherwise restricted request interface. A transport fix should test repeated and concurrent requests retaining earlier responses, including multipart and partner paths. Do not infer successful WhatsApp delivery from an HTTP fixture.

## 4) Multipart and media flow

**Entry and sequence:** `SendingCategory.SendFileByUpload` prepares endpoint fields including a local `file` path and enables form-data; `GreenAPI.request` selects the URL and calls `MultipartRequest`. That helper decodes JSON into a map, writes string fields, opens the file, detects MIME type, builds a `file` part using the path's basename, copies bytes, closes the file and multipart writer, and returns a pooled request. The caller sends it and releases its request and response.

**Behavior:** the helper hardcodes HTTP POST regardless of its `method` parameter. It buffers the multipart body in memory. It uses unchecked string assertions for `file` and other fields, so unexpected types can panic. Missing file, invalid JSON, file/MIME/read/write/close errors can return errors. File close is both deferred and called explicitly on the successful encoding path; the explicit close error is checked. There are no upload retries or caller cancellation in this path.

**Boundary:** upload-by-URL, binary upload and multipart are not interchangeable. Read the exact `sending.go` method and its `WithMediaHost`, `WithSetMimetype` or `WithFormData` choices before changing shared code. MIME, filename and field spellings are endpoint contracts. A local multipart helper change can silently break a different sending operation.

**Verification:** no current root tests cover this route. Local server tests should parse multipart fields and bytes, verify selected host and filename/MIME, and exercise missing file, malformed map and file-read failures. Large-file behavior is bounded by available memory in the current implementation, not streaming. External uploads are unnecessary for encoding coverage.

## 5) Notification poll and receipt flow

```text
consumer → ReceiveNotification → one GET → APIResponse
consumer → inspect status / decode / process according to application policy
consumer → DeleteNotification(receiptID) → one DELETE → APIResponse
```

**Entry and sequence:** `ReceivingCategory.ReceiveNotification` applies options, includes `?receiveTimeout=N` only when nonzero, and sends a GET request without a body. `OptionalReceiveTimeout` documents 5–60 seconds/default 5 but does not validate bounds. `DeleteNotification` appends `/{receiptID}` through `WithGetParams` to the assembled endpoint: `/waInstance{ID}/deleteNotification/{token}/{receiptID}`. The receipt is not inserted before the token.

**Output and ownership:** both operations return the root raw response envelope. No polling goroutine, dispatch, receipt storage or automatic acknowledgement exists. The application owns decoding, processing, duplicate handling and deletion timing. `examples/receiveNotification/main.go` demonstrates a single poll, not a durable consumer.

**Failure semantics:** root transport and status policy apply. A processing failure does not trigger any SDK action; an application that deletes before completing its chosen processing boundary risks losing recoverability. Conversely, a delete failure can leave processing uncertain. Queue retention, redelivery, ordering and exactly-once behavior are provider properties not established by this code. Adding automatic deletion would change the product contract, not just simplify a helper.

**Verification:** the root regression test checks that the polling GET has no body, omits the timeout when absent or zero, and includes an explicit timeout in the query. Changes to deletion or response handling should cover the DELETE suffix and error statuses. Consumer workflow tests should cover processing failure before deletion, deletion failure and duplicate input under the consumer's chosen policy.

## 6) VoIP REST and signaling flow

### 6.1) Call commands

`voip.New(Config)` copies a supplied HTTP client and overrides redirect handling to reject following redirects. Its default transport uses environment proxies, disables keep-alives and sets five-second connection/TLS/header and idle read timeouts, with no total request timeout. The caller's context is carried into `net/http`. `Config.HTTPClient` affects these REST calls, not the WebSocket factory.

`GetStateAsync`, `GetIceServersAsync`, `DialAsync`, `AcceptAsync`, `RejectAsync` and `HangUpAsync` pass through `Voip.request`. It closes response bodies. All 2xx statuses succeed; 204/empty bodies skip decoding. Other statuses and malformed JSON fail. `RequestError` retains an underlying cause; caller cancellation is returned when detected. Commands are not automatically retried: a lost response does not prove dial/accept/hangup was never executed.

These semantics are intentionally described separately from root REST. A shared response-policy refactor would change callers' expectations. `voip/calls_test.go` covers local REST status/body behavior; it does not validate live account authorization or provider idempotency.

### 6.2) Socket open, reconnect and close

**Entry:** `Voip.Connect` → `WebsocketUrl` → `NewReconnectingSocket` → `CallsConnection.OpenAsync(ctx)`. The URL converts HTTP(S) to WS(S), uses `callsRtc`, and clears query/fragment. Invalid host errors are retained for opening. The first open starts an owned background connect loop.

**Sequence and output:** dial → install socket → `connect` event → sequential read/dispatch and heartbeat → disconnect classification → delay/retry. Valid JSON messages are copied before dispatch; invalid JSON is ignored. There is no outbound queue: `Send` fails without a connected socket. Listener snapshots preserve synchronous registration order; panics are recovered. Async application callbacks have different completion ordering.

**Failure policy:** reconnect delay grows 500 ms → 1 → 2 → 4 → 8 → 10 seconds and stays capped; successful connection resets it. There is no jitter. HTTP handshake 401/403 and WebSocket close codes 4000–4999 permanently stop retries. `HandshakeRefusedError` omits the token-bearing URL from its displayed error while retaining response/cause information. Generic connection failures retry. Heartbeat ticks every 20 seconds and allows 20 seconds for ping completion; handshake timeout is 10 seconds and read limit is 1 MiB.

**Lifetime:** opening-context cancellation stops the caller's wait, not the background loop. `CloseAsync` cancels and waits for that loop if started and closes the socket; repeated close is safe. A failed open therefore still requires owner cleanup if the connection is abandoned. A change that binds the background loop to one opening context would change reconnect lifetime.

**Verification:** `voip/signaling_test.go` uses local HTTP/WebSocket and socket doubles for refusal, reconnect, dispatch and closure. Preserve heartbeat shutdown, close-before-open, cancellation while waiting, and distinction between permanent refusal and transient loss. Successful local reconnection does not prove TURN or provider connectivity.

## 7) Call state, negotiation and audio lifetime

### 7.1) State owners and event consequences

`CallsConnection` uses `op` to serialize protocol mutations and `mu` to protect shared data. State strings, peer/pending pointers, resume flags and generation counters together form its current transition model; there is no single formal FSM table. `onMessage`, `onState`, `onSocketDisconnect`, startup/answer handlers and teardown participate. Contribution policy requires explicit transition discipline for new changes, without pretending the current representation is already centralized.

| Input/condition | Current consequence |
| --- | --- |
| Any decoded state | Store it, then emit `state`. |
| Transition into `inc-call`, with nonnil info | Emit `incoming_call`; a repeated `inc-call` does not emit it again. |
| `inc-call`, `out-call` or `on-call` → non-call state | Reject pending startup, stop/tear down bridge, emit `end_call`. Unknown strings count as non-call. |
| Transient disconnect during a recognized call with a peer | Dispose bridge, preserve resume intent; a later active state may create fresh audio without dial/accept. |
| Disconnect without resumable work / permanent refusal | Clear applicable state/resume intent; may emit `end_call` with `connection-lost`; permanent refusal prevents resume. |
| Closed connection receives signaling | Ignore it. Startup on closed/active/pending connection is rejected. |
| Repeated answer after remote description set | Ignore duplicate. Answer without peer has no effect. |
| Generic server error while startup pending | Fail startup and tear down. After startup, emit error without automatically destroying the bridge. |

`CallInfo.Name` is a pointer so absent/null differs from an empty string. Event payload types are listed in [voip/README.md](voip/README.md). `On` callbacks run synchronously and cannot block or invoke lifecycle methods: protocol locks may still be held. `OnAsync` starts goroutines; completion can be out of order and application work needs its own lifetime control. Both return unsubscribe functions. Callback panic recovery is not a replacement for application error handling.

### 7.2) Bridge negotiation

```text
Open → Dial (outgoing) / incoming event → Accept (incoming)
     → StartAudio → fresh CallAudio → ICE servers → peer → AddTrack
     → offer + local gathering → send offer and candidates
     → remote answer → buffered remote candidates → startup complete
```

Dial and accept do not start audio. `StartAudioAsync` reserves pending startup and a media context, invokes `AudioFactory`, validates all `CallAudio` fields, fetches ICE servers, creates the bridge, binds handlers and adds the track. Offer creation and local gathering precede sending. Gathered candidates are also sent individually afterward. Early remote candidates queue until the answer; application order is preserved, including end-of-candidates.

Startup completes when the answer and buffered candidates have been applied, not when audio is audible. There is no built-in answer deadline. Context cancellation during startup rejects it and tears down resources. After successful negotiation, the bridge lifetime is detached from that startup context and ends on explicit stop/close or relevant call/peer events.

Factory, ICE, peer, offer, send, answer and candidate failures can abort startup. A factory must clean its own partial acquisitions when returning an error; once valid audio is accepted, the connection invokes its close callback at teardown. Incomplete or superseded audio is closed if possible. Generation checks prevent old factory results, track callbacks and peer errors from corrupting a replacement. An active peer failure tears down that generation and emits an error. This makes generation checks part of correctness, not an optional optimization.

### 7.3) Teardown and media contracts

`teardownBridge` increments the generation and detaches current pointers, cancels media, cancels/waits remote attachments, optionally sends a bounded ten-second `stop`, closes the peer, then calls application `CallAudio.Close`. It joins errors on its return path; some internal callers deliberately discard or log teardown errors. The architecture does not promise all cleanup failures reach the consumer.

`StopAudioAsync` rejects pending startup, clears resume intent and sends `stop` if a peer exists. `CloseAsync` marks the connection closed, tears down without `stop`, then closes signaling. Repeated close returns nil rather than replaying the first cleanup error. Neither operation calls `HangUpAsync`; local disposal is distinct from ending the server call.

Applications create fresh local source/sink resources on every factory invocation, including resume. They implement cancellation-aware `Recv` and remote attachment callbacks and release their sources/sinks in `CallAudio.Close`. Waiting for attachments before disposal prevents use-after-close; a callback that ignores cancellation can hang that wait. `FrameAudioSink.AttachAsync` launches frame delivery with sequential receive/callback processing for each attachment; its close cancels and waits for active work.

`RTCAdapter` negotiates Opus only. `NewPCMTrack` supplies a codec factory and one or two channels. Frames are interleaved signed 16-bit PCM at 48 kHz; PTS is samples per channel. Outgoing packets use 960 samples/channel (20 ms), paced by the adapter; incoming RTP is reordered and decoded. Codec instances belong to a bridge and direction, not shared sessions. Incomplete final frames are not flushed on EOF; sources needing those samples sent must pad. EOF does not hang up a call. Resampling and audio-device I/O belong to the application.

**Verification and change risk:** `voip/calls_test.go` covers startup, answer/candidate ordering, cancellation, generation protection, disconnect/resume and cleanup. `voip/audio_test.go`, `voip/rtc_test.go`, `voip/rtc_adapter_test.go` and `opus/opus_test.go` cover audio/mapping/codec behavior. `test/voip/calls_integration_test.go` exercises local REST → WebSocket → Pion and bidirectional PCM; it skips when the bundled codec is unavailable. Changing teardown order, callback dispatch or frame sizing can break an otherwise unchanged consumer during reconnect. Race and failure-path tests matter as much as a successful offer/answer.

## 8) Integration and extension limits

Configuration and consumer code are the first customization points. Keep branding and durable business behavior in the consuming application. There is no White Label product/configuration subsystem here.

- Root categories expose `GreenAPIInterface`, but its `Request` signature contains unexported `requestOptions`. External packages cannot ordinarily implement this interface by declaring their own matching method. Do not promise unrestricted root transport injection.
- `voip.Config.HTTPClient` and logger are public seams. The copied HTTP client affects VoIP REST only; redirect behavior is replaced. It does not configure root `fasthttp` or WebSocket dialing.
- `Bridge` is exported, and `NewRTCAdapter` is public, but `Voip.Connect` selects a private bridge factory. There is no public bridge-factory option through this constructor. A custom bridge must report failures via `OnError`, but substituting it in the standard connection path requires a deliberate API change.
- `AudioFactory`, `CallAudio`, track contracts and `CodecFactory` are genuine public extension points with cancellation and fresh-resource constraints.
- `opus.New` has a CGo implementation and a CGo-free implementation returning an error. C compiler/system libopus requirements depend on the pinned gopus target. A custom codec can enable media without the bundled codec; signaling does not require it.

The project uses the [MIT License](LICENSE), which permits modification and redistribution subject to its terms. File/dependency coverage and distribution rights need separate review before a distribution decision.

## 9) Verification reality and external unknowns

[CONTRIBUTING.md](CONTRIBUTING.md) defines task-dependent checks. `voip/README.md` documents `go test -race ./...` and `CGO_ENABLED=0 go test ./...`. No local CI workflow, lint/build configuration or release compatibility policy is present in the inspected checkout. Existing tests are executable evidence only when run and their results reported; their presence does not imply they pass.

Local tests do not verify live GREEN-API endpoint compatibility, account/queue delivery guarantees, TURN deployment, network quality or behavior across every Go/CGo target. Root REST has no current test suite. The response-buffer lifetime issue and multipart assertions above are observed implementation limitations, not standards for new code. This document introduces no runtime changes.
