# whatsapp-api-client-golang-v2

Instructions for AI agents working with the repository.

---

## 0) Quick facts about this repository

This public Go library wraps GREEN-API's WhatsApp REST API. Its separate `voip` package implements call commands, reconnecting signaling and WebRTC audio; `opus` supplies an optional bundled codec. Consumers are application developers. This repository is not a deployed chatbot, webhook receiver or audio-device application.

- Module: `github.com/green-api/whatsapp-api-client-golang-v2`; root package name: `greenapi`.
- Toolchain: `go.mod` requires Go 1.23.0 and selects go1.23.2. Dependencies include fasthttp, coder/websocket, Pion WebRTC v4 and gopus.
- REST entry: construct `GreenAPI` in `types.go`, choose a category through `categories.go`, call its endpoint method. Responses are raw `APIResponse` envelopes.
- Calls entry: `voip.New(Config)` → `Connect(AudioFactory)` → `OpenAsync(ctx)`; construction does not dial or open signaling.
- Sources of truth: endpoint files plus `request.go`; `voip/calls.go`, `voip/signaling.go`, `voip/audio.go`, `voip/rtc_adapter.go`; `go.mod`; actual tests. README examples describe usage, not verified live behavior.
- Before transport, state, lifecycle or extension changes, read the matching flow in [ARCHITECTURE.md](ARCHITECTURE.md). Before any code contribution or review, read [CONTRIBUTING.md](CONTRIBUTING.md) for mandatory standards, examples and test gates. These linked files are not assumed to load automatically.

## 1) Agent role and operating modes

Act as a maintainer of a public integration library: explain capabilities accurately and preserve consumer contracts. Match work to the user's intent.

**Product/beginner explanation:** give a simple purpose, confirmed capabilities, prerequisites and one practical scenario. Include a short real public-API example in the first useful introduction when appropriate. Do not start an unrelated audit or imply application features ship with the SDK.

**Integration guidance:** establish whether the consumer uses root REST, VoIP REST, signaling or media. Explain credentials, response policy, lifecycle and actual public seams. Use environment-provided configuration; never invent usable account identifiers or claim local checks prove provider behavior.

**Engineering changes/review:** follow the numbered cycle below; read the relevant companions and source before editing. Distinguish implementation facts from mandatory new-code policy. Report findings at concrete paths and observable consequences.

**Application customization/White Label:** use the extension ladder in section 5. Keep application branding and business policy in consumer code. Confirm public constructor access and distribution terms before describing a fork or substitute as supported.

## 2) Product knowledge and honest presentation

Explain: “This Go SDK lets an application call GREEN-API to work with a WhatsApp account. The application supplies its business logic and credentials.” Confirmed capabilities include sending messages/media, managing account/contact/group operations, polling and deleting notifications, querying queues/journals/statuses, and call signaling with an application-supplied audio source/sink.

A practical scenario is a service that polls for notifications, processes each response using its own durable business rules and explicitly deletes the receipt after the chosen processing boundary. The SDK supplies individual HTTP operations; it does not provide durable processing, automatic acknowledgement, duplicate suppression or queue-delivery guarantees.

**Verified public API shape — consumer example.** This function makes a real request if called. The environment variables are application conventions, not SDK-provided configuration loading. Imports are `fmt`, `os` and the module aliased as `greenapi`.

```go
func sendGreeting(chatID string) error {
	api := greenapi.GreenAPI{
		APIURL:           os.Getenv("GREEN_API_URL"),
		MediaURL:         os.Getenv("GREEN_API_MEDIA_URL"),
		IDInstance:       os.Getenv("GREEN_API_ID"),
		APITokenInstance: os.Getenv("GREEN_API_TOKEN"),
	}
	if api.APIURL == "" || api.IDInstance == "" || api.APITokenInstance == "" {
		return fmt.Errorf("GREEN-API URL, instance ID and token are required")
	}
	response, err := api.Sending().SendMessage(chatID, "Hello")
	if err != nil {
		return fmt.Errorf("send greeting: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("send greeting: HTTP %d", response.StatusCode)
	}
	return nil
}
```

Prerequisites: Go 1.23+, an authorized GREEN-API instance and network access for actual account operations. Media methods may need `MediaURL`. The snippet does not validate all configuration or prove message delivery. Root calls have no context argument/configured timeout and return non-2xx statuses in `APIResponse`, not as Go errors. Avoid logging underlying request errors without considering URL/token disclosure.

For calls, outgoing order is Open → Dial → StartAudio; incoming order is Open → incoming event → Accept → StartAudio. Dial/accept do not start audio. Bundled Opus availability depends on CGo/target requirements; a custom codec can be supplied. Microphone/speaker I/O and format conversion are application work. See [voip/README.md](voip/README.md) for the full consumer lifecycle.

## 3) Agent engineering protocol

1. **Classify and bound:** identify affected package, user-visible outcome and change category. Inspect working-tree status and preserve unrelated changes.
2. **Load appropriate guidance:** read CONTRIBUTING before coding/review; read the relevant ARCHITECTURE flow for transport, receipt, state, audio or extension work. Consult localized/VoIP README usage when changing public behavior.
3. **Trace the actual path:** public entry → encoding/validation → transport or state owner → side effect → failure/cleanup. Read the specific wrapper and tests; do not extrapolate from a similarly named endpoint.
4. **Identify contracts:** record signatures, payload/omission, URL order, status/error policy, callback order, retry/timeouts, resource owner and cancellation. For state work enumerate invalid, repeated and terminal events and generation identity.
5. **Choose extension level:** reuse existing behavior/configuration; prefer consumer composition/public seams. Add a narrow interface only for demonstrated coupling/testing needs. Check constructor access before promising substitution.
6. **Apply corporate rules:** name contracts clearly, separate responsibilities, meet function size limits, centralize constrained transitions and define cleanup/error reporting. Request an explicit reviewer exception only if a concrete compliant solution cannot satisfy the change; never self-approve one.
7. **Implement the smallest complete diff:** include required local tests/docs; improve materially touched legacy code without unrelated sweeps. Do not change runtime code during a documentation-only task.
8. **Verify according to risk:** run focused local checks and broader race/CGo-free checks for cross-package or lifecycle changes. Inspect failures/skips and fix relevant issues. Do not execute credentialed examples as local verification.
9. **Report truthfully:** explain changed behavior, preserved contracts, exact checks/results and unresolved limits. Separate checks run from recommended commands and local evidence from live-provider claims.

## 4) Architecture and code placement

```text
root consumer → category file → Request/request.go → fasthttp → GREEN-API
call consumer → Voip REST + CallsConnection → ReconnectingSocket → callsRtc
                              ↘ RTCAdapter/Pion ↔ application CallAudio + CodecFactory
```

Endpoint payloads/options belong in their category (`sending.go`, `receiving.go`, `account.go`, etc.); shared request mechanics belong in `request.go`. Partner requests have a separate path. `types.go` owns configuration/envelope shapes. Do not put durable workflow policy into endpoint wrappers.

`CallsConnection` in `voip/calls.go` owns negotiation and bridge lifetime, serializes protocol mutations with `op`, guards shared data with `mu`, and uses generations to reject stale work. `ReconnectingSocket` owns reconnect/read/heartbeat. `RTCAdapter` owns the peer and media transport; application resources returned through `AudioFactory` are released by the supplied `CallAudio.Close`. A change to one owner's lifetime affects the others during cancellation and resume.

Root transport currently creates a client per call and releases pooled responses while returning their body/status-message slices without copying. Body lifetime is therefore unsafe to assume after release. Document this limitation honestly; do not replicate it in new code or claim a copy after return fixes concurrent reuse. A transport repair needs response-retention tests across normal/multipart/partner paths.

## 5) Integration and White Label

Customization ladder: **configuration → consumer code → existing public extension → narrow reusable option/interface → core change → fork**. Use the earliest level that completely meets the real need. There is no White Label configuration subsystem or branded end-user product in this repository.

- `GreenAPI` permits configured API/media hosts and instance fields. Its `GreenAPIInterface.Request` has unexported `requestOptions`; an ordinary external implementation cannot declare that exact signature. Package-local fakes or a consumer-owned capability interface are different solutions.
- `voip.Config.HTTPClient` is copied and redirects are disabled; it affects only VoIP REST. It does not replace root fasthttp or signaling transport.
- `AudioFactory`, `CallAudio`, track interfaces and `CodecFactory` are public seams. Factories must create fresh resources for each bridge/resume.
- `Bridge` is exported but standard `Voip.Connect` hardcodes a private bridge factory. Do not claim custom bridge injection is already supported there.
- Keep branding, storage, receipt policy and audio hardware in the application. Verify constructor requirements and owner cleanup for every substitution.

The project uses the [MIT License](LICENSE), which permits modification and redistribution subject to its terms. Establish applicable file/dependency terms separately before a distribution decision.

## 6) Public and implicit contracts

Root REST constructs `/waInstance{ID}/{method}/{token}`, appending query or suffix afterward. Category options/JSON tags/omission and API-versus-media host are compatibility surfaces. Non-2xx responses remain envelopes; local/transport failures are errors. VoIP REST instead errors on non-2xx/malformed JSON, accepts all 2xx and skips decoding for 204/empty bodies. Neither automatically retries commands.

Polling and deletion are separate. `DeleteNotification` uses `/deleteNotification/{token}/{receiptID}` after the instance prefix. No handler, automatic acknowledgement or provider retention/ordering guarantee is implemented locally.

Synchronous `On` callbacks cannot block or call lifecycle methods; `OnAsync` allows that work in goroutines but completion can be out of order. Opening-context expiry does not stop the reconnect loop; close abandoned connections explicitly. Startup-context cancellation aborts negotiation; after success active media has its own lifetime. Stop/close do not hang up the remote call.

Audio must be 48 kHz signed 16-bit interleaved mono/stereo; Opus packets use 960 samples/channel per 20 ms. Cancellation-aware sources/attachments and fresh per-bridge/per-direction codecs are required. Teardown cancels and waits attachments before peer/audio disposal. Incomplete final frames are not flushed on EOF; EOF is not hangup.

## 7) Decision rule — REST endpoints and response contracts

**When:** adding/changing an endpoint, option, shared URL, status behavior or response handling.

**Read:** `categories.go`, the exact category method, `types.go:APIResponse`, `request.go:Request`, `GreenAPI.request`; include `PartnerRequest` if shared behavior is affected.

**Determine:** HTTP method/path, JSON tags and omitted fields, validation, query/suffix construction, media/form-data/MIME choices, error/status policy and byte ownership. Identify all callers of a shared helper.

**Preserve:** category façade, module/API shape, raw envelope and non-2xx-in-response semantics unless intentionally migrated. Do not add a context parameter or retry as an undocumented behavior change.

**Verify:** package-local fake or local HTTP server assertions for method/path/body/options, invalid input, transport failure and non-2xx. Root has no current tests: new transport contracts require tests. For response ownership fixes retain multiple responses during repeated/concurrent calls and cover all affected paths; never use a live account to prove encoding.

## 8) Decision rule — notification polling and acknowledgement

**When:** changing `ReceiveNotification`, `DeleteNotification`, queue examples or adding consumer processing.

**Read:** `receiving.go`, `queues.go`, `examples/receiveNotification/main.go`; ARCHITECTURE section 5 for exact receipt suffix and ownership.

**Determine:** who decodes/status-checks, processes durably, handles duplicates and decides deletion; whether failures leave the receipt unresolved. Zero timeout is omitted; 5–60 is documented but not enforced. New validation changes accepted inputs.

**Preserve:** poll/delete separation, receipt identifier and URL ordering, one-call SDK behavior. Do not imply local queues, handlers or exactly-once delivery exist.

**Verify:** omitted/explicit timeout, GET and DELETE construction, status failures. If adding application workflow, test process success before delete, process failure preventing delete, delete failure and repeat processing under its explicit policy. Do not invent provider retention guarantees.

## 9) Decision rule — upload and media routing

**When:** changing `SendFileByUpload`, `UploadFile`, MIME/filename options or `MultipartRequest`.

**Read:** those methods in `sending.go`, request option functions and multipart/normal branches in `request.go`.

**Determine:** actual host/method, local path requirements, string-valued JSON map, field names and filename/MIME behavior. Existing helper hardcodes POST, buffers the full body and uses unchecked assertions; these are limitations, not new-code standards.

**Preserve:** correct field/host routing and file closure. Keep binary, URL and multipart uploads distinct. For deliberate input-error improvements, test what used to panic and document the new failure contract.

**Verify:** temporary files/local servers for part bytes/fields/basename/MIME, media host selection, missing files, malformed JSON/types and read errors. No external upload is needed. Check response lifetime in the multipart return path if transport is touched.

## 10) Decision rule — signaling, state and callbacks

**When:** changing reconnect, `CallsConnection` events, negotiation, stop/end/disconnect or resume.

**Read:** `voip/calls.go:Connect`, `onMessage`, `onState`, `onSocketDisconnect`, `StartAudioAsync`, `teardownBridge`; `voip/signaling.go`; `voip/calls_test.go`, `voip/signaling_test.go`; ARCHITECTURE sections 6–7.

**Determine:** mutation/lock owner, current peer generation, pending completion and resume intent; synchronous versus asynchronous callbacks. Define invalid/repeated/terminal events and whether a context controls waiting, startup or resource lifetime.

**Preserve:** backoff 500 ms doubling to 10 s/reset on connect; permanent HTTP 401/403 and WS 4000–4999 refusal; early ICE ordered after answer; ignored duplicate answers; generation guards. State is stored/emitted before incoming/end handling; recognized call → non-call tears down before end event. Transient loss can resume with fresh audio; permanent refusal cannot. Generic server errors after startup differ from peer failure. Stop sends bounded `stop`; close does not; neither hangs up.

**Verify:** local socket/call tests for permanent/transient loss, open timeout plus close, ordered dispatch, late answer/peer errors, cancellation competing with negotiation, repeated close, call-end and repeated disconnect during resume. Use race testing for changed concurrency. Never equate reconnect with permission to retry call commands.

## 11) Decision rule — audio/provider substitution and lifecycle

**When:** changing audio source/sink contracts, codec/peer integration, factory construction or cleanup.

**Read:** `voip/audio.go`, `voip/rtc.go:Bridge`, `voip/rtc_adapter.go`, `opus/opus.go`, `opus/opus_nocgo.go`, `prepareAudio`, `preparePeer`, `teardownBridge` in `voip/calls.go`, and [voip/README.md](voip/README.md).

**Determine:** whether the proposed seam is actually publicly injectable; who owns each resource on success and partial failure; cancellation response, close errors, frame/channel/EOF rules and target CGo constraints. `Bridge.OnError` must report active peer failures.

**Preserve:** fresh application resources per bridge/resume; per-direction codecs; cancellation and attachment wait before peer then audio close; stale-result disposal. A factory returning an error must clean its partial acquisitions. Do not transfer ownership merely by renaming a helper.

**Verify:** factory/ICE/peer/offer/answer failures, cancel during startup/attachment, stale result after stop, codec failure, mono/stereo/invalid format and final partial frame. Use audio/adapter/codec tests plus local Pion integration for cross-layer changes; report codec skips. No microphone or live calls are needed for local protocol tests.

## 12) Corporate engineering standards and review

**Mandatory for new/materially changed functions:** at most **50 nonblank, noncomment body lines**, target about 30; one concrete responsibility. Comment-only/blank lines do not count; code with trailing comments does. Closures count too and remain inside their enclosing body's count. No giant closures, compressed statements, deep nesting or meaningless wrappers to evade the rule. No current linter enforces it; reviewers must. Touched legacy functions comply within a cohesive scope; only explicit reviewers may approve a narrowly documented exception.

Use plain English/domain names and explicit input/output/error/state/cancellation contracts. Define small consumer-owned interfaces only at real integration/testing boundaries; inject adapters and test fakes there. Keep fixed details concrete. Constrained transition systems require one mutation owner and an explicit FSM/equivalent model with invalid/repeated/terminal tests, not scattered flags. Do not force FSMs on trivial flags or migrate call semantics inadvertently.

Name every resource owner and stop path. Handle partial acquisition and cleanup errors; document idempotency and blocking limitations. Protect payloads, retries/timeouts, ordering, observability and compatibility. Reject duplicate logic, hidden globals, swallowed errors, speculative layers, unbounded work and workarounds without scope/reason/removal conditions.

**Synthetic resource pattern — not an SDK API.** Caller supplies nonnil functions; successful `open` returns a nonnil resource. `use` borrows it, honors cancellation, returns only after child work ends and never closes/retains it. A failing opener cleans its own partial acquisitions. Imports: `context`, `errors`, `io`.

```go
func withResource(ctx context.Context,
	open func(context.Context) (io.Closer, error),
	use func(context.Context, io.Closer) error) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}
	resource, err := open(ctx)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, resource.Close()) }()
	if err := ctx.Err(); err != nil {
		return err
	}
	return use(ctx, resource)
}
```

This owner closes exactly once and preserves work plus close errors. It cannot forcibly cancel an uncooperative implementation or bound Close. Test acquire failure, success, cancellation and joined work/close errors. Read CONTRIBUTING for complete bad/good patterns covering naming/decomposition, interfaces/adapters/fakes, transitions and partial resource acquisition.

A proposed change **fails review** for a mandatory rule violation, unapproved compatibility break, avoidable debt, unclear ownership, sensitive data exposure, unverifiable correctness or missing risk-appropriate tests. An exception requires explicit reviewer approval of its scope, reason, risk, validation and removal/follow-up condition; agents cannot waive rules.

## 13) Verification and completion

Safe local checks depend on the change. `voip/README.md` documents these broad variants:

```sh
go test -race ./...
CGO_ENABLED=0 go test ./...
```

Focused checks are `go test ./voip`, `go test ./opus`, `go test ./test/voip`; use `go test .` for newly added root tests. Format changed Go files with `gofmt` and run `git diff --check`. Broad tests compile examples but do not execute their `main` functions. Local tests use loopback servers and may need socket/cache permissions. Dependency downloads and the bundled codec may require network/build tooling; CGo-free media integration can skip when `opus.New` is unavailable.

For documentation-only work validate links, source/API claims, Markdown and sample compilation; the full media suite is not required. No local CI/lint workflow is present, so never label these commands CI enforcement. Never claim existing tests passed without running them. Report skipped checks and material gaps, including root test absence and live GREEN-API/TURN properties not established by local testing.

---

## 14) Do Not

- Do not execute credentialed examples or send/dial/delete/hang up on a real account as routine verification.
- Do not log full token-bearing URLs, private payloads or ICE credentials, or copy inline credential-like examples into new docs.
- Do not treat nil root REST error as HTTP success, or conflate root and VoIP response policies.
- Do not add automatic receipt deletion, command retries or provider guarantees silently.
- Do not block/reenter lifecycle from `On`, remove generation guards or dispose audio before attachments finish.
- Do not equate open timeout with loop shutdown, stop/close with hangup, or EOF with call completion.
- Do not promise public bridge injection or safe pooled-response lifetime that current code does not provide.
- Do not present mandatory new-code policy as existing implementation, self-approve exceptions, or alter unrelated work.

## 15) Pre-change checklist

- [ ] Task and affected consumers identified; working-tree changes preserved.
- [ ] CONTRIBUTING read for code work; matching ARCHITECTURE flow and source/tests read.
- [ ] Actual endpoint/options/status/receipt or state/audio contracts recorded.
- [ ] Extension access verified; minimal complete approach selected without speculative layering.
- [ ] Function size, naming, dependency direction and transition policy satisfied in touched scope.
- [ ] Owners, partial failure, cancellation, generation and cleanup-error behavior explicit.
- [ ] Local success/failure/boundary/race tests chosen for risk; live actions excluded unless intended.
- [ ] Compatibility/migration and README/doc impacts addressed; no secrets or unsupported licensing claims.
- [ ] Completion will report actual checks, skips, preserved behavior and external unknowns.
