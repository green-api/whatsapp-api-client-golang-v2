# Contributing — whatsapp-api-client-golang-v2

These rules apply to human and AI contributions. [AGENTS.md](AGENTS.md) is the agent entry point; [ARCHITECTURE.md](ARCHITECTURE.md) describes current behavior and ownership. This document is mandatory policy for new and materially changed code, not a claim that all legacy code already complies.

---

## 1) Engineering goal and workflow

Deliver the smallest complete, safe and maintainable change, including appropriate tests and documentation. Preserve published and implicit contracts unless an intentional migration is approved. Fewer changed lines do not justify missing cleanup or testing.

1. **Classify the task:** explanation, endpoint change, notification semantics, upload, signaling/state, media/lifecycle or consumer integration. Identify the affected packages and consumers.
2. **Read the relevant architecture flow and source:** trace public entry through options, transport or state mutation to side effects and cleanup. Inspect existing tests before choosing a test seam.
3. **Write the contract:** inputs, outputs, errors, status policy, state/side effects, cancellation, resource owner and compatibility impact. For protocol changes list allowed, invalid, repeated and terminal events.
4. **Choose the smallest suitable extension:** reuse configuration and public seams first; keep application policy outside SDK transport. Explain any new interface by its real consumer and testing/coupling problem.
5. **Implement within touched scope:** meet the rules below immediately for new code and materially changed functions. Improve necessary legacy scope without unrelated sweeps or renaming public APIs for style.
6. **Verify by risk:** add local failure and boundary coverage; run focused checks, then broader checks when the change crosses package or concurrency boundaries. Never substitute live examples for tests.
7. **Review the diff:** preserve unrelated user changes, inspect credentials/log output, verify doc links and API examples, and check migration notes if behavior changes.
8. **Report evidence:** explain the problem, resulting behavior, preserved contracts, exact checks/results/skips, and remaining limits. A passing local test is not a provider guarantee.

## 2) Function size, naming and contracts

**New and materially changed functions must contain at most 50 nonblank, noncomment body lines; target about 30.** Count physical lines inside the function body, excluding blank and comment-only lines. Lines mixing code and comments count. Closures are functions too, and their code remains part of the enclosing function's body count. Do not combine statements, use giant closures, deep nesting or meaningless forwarding wrappers to evade the limit. No repository linter currently enforces this rule; review must.

Each function has one concrete reason to exist. A newcomer should understand its inputs, output/effect, errors, state changes and owner in about five minutes. Split validation, orchestration, encoding, I/O and cleanup where those are distinct responsibilities, without producing layers that add no meaning.

Use common English and domain terms, Go initialisms and package context. Prefer `ValidateReceiveTimeout` for error-returning validation and `HasPendingAudio` for a boolean predicate; vague `Handle`, `DoStuff` or misleading `Ensure` names hide contracts. Do not rename published methods just for taste. Signatures, types and documentation must expose valid inputs, nil/zero behavior, outputs, errors, side effects, concurrency and cancellation. Explain why a constraint exists instead of narrating syntax.

Legacy functions are not exempt once materially changed. Bring the touched function into compliance through a cohesive decomposition and regression tests; do not refactor unrelated functions. If a narrow exception is necessary, record the exact function/rule, reason, risk and follow-up/removal condition and obtain **explicit reviewer approval**. An agent cannot grant itself an exception.

### Pattern A — clear naming and decomposition

All examples in this document are **synthetic Go teaching examples**, not new SDK APIs or patches. The good examples are standalone packages and use only the standard library. Their policies illustrate future code design, not silently changed endpoint behavior.

**Bad:** this fragment conflates a predicate, validation, encoding and transport; it also swallows the outcome.

```go
func Ensure(v int) bool {
    if v != 0 && (v < 5 || v > 60) { return false }
    send("?receiveTimeout=" + strconv.Itoa(v)) // result ignored; zero changed
    return true
}
```

**Good:** a pure encoder with an explicit validator. In this synthetic policy zero omits the option. The current SDK does not enforce these bounds; introducing them there is an intentional compatibility change.

```go
package timeoutoption

import (
	"fmt"
	"net/url"
	"strconv"
)

// ValidateReceiveTimeout accepts omission (zero) or 5–60 seconds.
func ValidateReceiveTimeout(seconds int) error {
	if seconds != 0 && (seconds < 5 || seconds > 60) {
		return fmt.Errorf("receive timeout must be zero or 5–60 seconds: %d", seconds)
	}
	return nil
}

// ReceiveQuery returns a query without '?'; it performs no network I/O.
func ReceiveQuery(seconds int) (string, error) {
	if err := ValidateReceiveTimeout(seconds); err != nil {
		return "", err
	}
	if seconds == 0 {
		return "", nil
	}
	values := url.Values{}
	values.Set("receiveTimeout", strconv.Itoa(seconds))
	return values.Encode(), nil
}
```

**Why:** callers can distinguish invalid input from transport failure and test encoding without an account. Decomposition exposes omission semantics instead of hiding them behind a boolean.

**Tests:** accept 0/5/60; reject -1/4/61; verify empty output on omission and error, and exact query spelling. Verify transport is not invoked on invalid input when this helper is integrated.

**Review:** reject ignored send errors, accidental default changes, vague predicate names and hidden I/O. A bounds check added to `OptionalReceiveTimeout` needs compatibility discussion because previously accepted inputs will fail.

## 3) Dependency direction and real extension boundaries

High-level orchestration must not know unnecessary provider, wire-format or storage details. At real integration/testing boundaries, define a small interface near the consumer, inject an adapter, and test with a fake. Keep fixed local details concrete. Do not add an interface to every struct, factories for imaginary options or speculative layers.

Root `GreenAPIInterface.Request` includes private `requestOptions`; it is not a freely implementable external transport interface. Use package-local tests or an application-owned capability interface instead of promising otherwise. For VoIP REST use `Config.HTTPClient` where it solves the problem; it does not configure signaling. `AudioFactory` and `CodecFactory` are public seams; `Voip.Connect` does not publicly accept a bridge factory. A new public seam must specify construction errors, ownership, nil/default behavior, concurrency, cancellation and compatibility before adding it.

### Pattern B — consumer contract, adapter and fake

**Bad:** consumer receipt orchestration constructs a provider URL and discards the result. Besides coupling policy to transport, this can acknowledge unsuccessful processing or expose tokens in logs.

```go
func Finish(receiptID int) {
    http.DefaultClient.Do(makeProviderDelete(receiptID)) // response/error lost
}
```

**Good:** a synthetic application's acknowledgement boundary. `ReceiptEndpoint` is an adapter input with a deliberately context-free call, similar to the root SDK limitation. It is not `greenapi.GreenAPIInterface` and does not claim to make root HTTP calls cancellable.

```go
package receiptpolicy

import (
	"errors"
	"fmt"
)

// ReceiptDeleter belongs to the consumer that decides when processing is done.
type ReceiptDeleter interface {
	DeleteReceipt(receiptID int) error
}

// ProcessReceipt calls delete only after process succeeds. It does not retry.
func ProcessReceipt(id int, process func() error, receipts ReceiptDeleter) error {
	if id <= 0 || process == nil || receipts == nil {
		return errors.New("positive receipt ID, processor and deleter required")
	}
	if err := process(); err != nil {
		return fmt.Errorf("process receipt: %w", err)
	}
	return receipts.DeleteReceipt(id)
}

// ReceiptEndpoint is a low-level synthetic endpoint returning HTTP status.
type ReceiptEndpoint func(int) (int, error)

type StatusAdapter struct{ Endpoint ReceiptEndpoint }

func (a StatusAdapter) DeleteReceipt(id int) error {
	if a.Endpoint == nil {
		return errors.New("receipt endpoint required")
	}
	status, err := a.Endpoint(id)
	if err != nil {
		return fmt.Errorf("delete receipt: %w", err)
	}
	if status < 200 || status >= 300 {
		return fmt.Errorf("delete receipt: HTTP %d", status)
	}
	return nil
}

// FakeReceipts is owned by a test; it is not safe for concurrent calls.
type FakeReceipts struct {
	IDs []int
	Err error
}

func (f *FakeReceipts) DeleteReceipt(id int) error {
	f.IDs = append(f.IDs, id)
	return f.Err
}
```

An application adapter around the real `ReceivingCategory.DeleteNotification` would extract `APIResponse.StatusCode` at this boundary. This example's positive-ID and processing policy belong to the application; they are not currently enforced by the SDK. Interface arguments require non-typed-nil implementations; ordinary `interface == nil` checks cannot detect every typed nil.

**Why:** processing rules depend only on deletion capability. The adapter owns status interpretation, while the fake records whether deletion happened. There is no fabricated `context.Context` guarantee over a context-free API. A context-capable implementation should expose and propagate context at its own real boundary.

**Tests:** success processes then deletes exactly the chosen ID; processing error prevents deletion; deletion error is returned and never automatically retried; zero ID/nil arguments fail before processing; adapter tests cover 199/200/299/300, endpoint failure and absent endpoint. Test duplicates under the application's chosen idempotency policy—this example does not supply durable exactly-once processing.

**Review:** reject provider URL construction in business policy, premature deletion, interfaces without consumers, fake cancellation guarantees and retries after uncertain side effects without an explicit policy.

## 4) Explicit state transition discipline

When correctness depends on allowed states/events, new work must use a centralized explicit FSM or equivalent transition model with one mutation owner. Define invalid, repeated and terminal semantics and transition tests. Separate transition planning from I/O when feasible. Do not apply a state machine to every trivial flag.

Current call behavior is distributed across handlers, pointers, locks and generation counters in `voip/calls.go`. Changes must preserve its event ordering and resume semantics, or explicitly migrate them. The synthetic machine below must not replace the provider's call states or introduce a different close policy by accident.

### Pattern C — transition planning and one owner

**Bad:** separate goroutines assign `state = "active"` after an answer and `state = "closed"` on shutdown. A late answer can resurrect closed work; no component enforces the allowed transitions.

```go
func onAnswer() { state = "active" }
func onClose()  { state = "closed" }
```

**Good:** a complete synthetic startup machine. An owner applies a pure transition function under its mutex; I/O stays outside this example.

```go
package startupstate

import (
	"fmt"
	"sync"
)

type State uint8

const (
	New State = iota
	Starting
	Active
	Closed
)

type Event uint8

const (
	Start Event = iota
	Answer
	Close
)

// Next rejects invalid input without changing state. Close is idempotent.
func Next(state State, event Event) (State, error) {
	if state > Closed || event > Close {
		return state, fmt.Errorf("unknown transition: %d/%d", state, event)
	}
	if event == Close {
		return Closed, nil
	}
	if state == New && event == Start {
		return Starting, nil
	}
	if state == Starting && event == Answer {
		return Active, nil
	}
	if state == Active && event == Answer {
		return Active, nil // repeated answer explicitly has no effect
	}
	return state, fmt.Errorf("invalid transition: %d/%d", state, event)
}

// Session's zero value is New. Only Apply mutates state.
type Session struct {
	mu    sync.Mutex
	state State
}

func (s *Session) Apply(event Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	next, err := Next(s.state, event)
	if err != nil {
		return err
	}
	s.state = next
	return nil
}
```

**Why:** late answer after close is rejected; duplicate answer has a documented result. Rule tests do not require sockets. A mutex prevents races but alone would not define valid transitions. A real peer implementation also needs generation identity: a valid event from an old peer must not mutate the current peer.

**Tests:** cover every state/event pair, invalid enum values, repeated Start rejection, repeated Answer in Active, repeated Close, Answer before Start and after Close; failures leave state unchanged. Exercise concurrent Apply with `-race`. In repository tests also cover stale factory/peer callbacks, transient loss during resume, permanent refusal, early candidates and cancellation competing with answer.

**Review:** reject scattered assignments, undocumented duplicate behavior, transition I/O while holding unnecessary locks, and generation removal without a replacement invariant. Reviewers must see the changed transition matrix or equivalent explicit model and its regression tests.

## 5) Resource ownership and cancellation

Every connection, goroutine, timer, queue, file and stream must have a named owner and end-of-life path. The owner handles partial acquisition, cancellation, cleanup and reporting. Upper levels call a simple `Close`/`Stop` method rather than coordinating adapter internals. Document idempotency, concurrent invocation and cleanup-error behavior precisely; never promise bounded close when it waits for uncooperative callbacks.

For this repository, `CallsConnection` cancels and waits attachments before closing peer and application audio. `CallAudio.Close` owns source/sink disposal; factories must return fresh resources and clean partial acquisitions on error. `Recv` and attachments must honor cancellation. No unowned background tasks or hidden infinite queues are allowed in new code.

### Pattern D — partial acquisition, cancellable work and joined cleanup

**Bad:** this fragment loses the first resource if the second acquisition fails and ignores every close error.

```go
func run(ctx context.Context) error {
    source, err := openSource(ctx)
    if err != nil { return err }
    sink, err := openSink(ctx)
    if err != nil { return err } // source leaked
    defer source.Close()
    defer sink.Close()
    return transfer(ctx, source, sink)
}
```

**Good:** a synchronous synthetic runner acquires source then sink; its owner closes sink then source. Callbacks are required to respect context and to release partial allocations on their own error returns.

```go
package ownedstream

import (
	"context"
	"errors"
	"io"
	"sync"
)

type Open func(context.Context) (io.Closer, error)

// Pair must be created by Acquire and not copied. Close waits once and caches
// its error. Close does not itself cancel or wait for external work.
type Pair struct {
	source, sink io.Closer
	once         sync.Once
	closeErr     error
}

func (p *Pair) Close() error {
	p.once.Do(func() {
		p.closeErr = errors.Join(p.sink.Close(), p.source.Close())
	})
	return p.closeErr
}

// Acquire requires nonnil openers returning a nonnil resource on success.
func Acquire(ctx context.Context, source, sink Open) (*Pair, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	first, err := source(ctx)
	if err != nil {
		return nil, err
	}
	second, err := sink(ctx)
	if err != nil {
		return nil, errors.Join(err, first.Close())
	}
	return &Pair{source: first, sink: second}, nil
}

// Run owns the pair. Work must be nonnil, return on cancellation, finish child work
// before returning, and must not retain or close either borrowed resource.
func Run(ctx context.Context, source, sink Open,
	work func(context.Context, io.Closer, io.Closer) error) (err error) {
	pair, err := Acquire(ctx, source, sink)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, pair.Close()) }()
	if err := ctx.Err(); err != nil {
		return err
	}
	return work(ctx, pair.source, pair.sink)
}
```

**Why:** every successfully acquired resource has an owner even on partial failure, and operation/cleanup errors remain available through `errors.Is`. `sync.Once` makes concurrent/repeated Close return the stored result; it does not make closing during active work safe. Only the owning runner closes after work finishes. This example has no background loop and no close deadline: blocking open/work/close implementations can still block it. Actual `CallsConnection.CloseAsync` returns nil on repeated calls instead of caching errors; do not transfer this synthetic policy to it implicitly.

**Tests:** normal work and reverse acquisition cleanup order; source acquisition error; sink error plus source-close error; canceled context before acquisition and during work; work error joined with both close errors; exactly-once close including concurrent close and repeated calls returning the stored error. Test that borrowed resources remain open during work and no child work survives its return. Fakes must make cancellation deterministic rather than rely on sleeps.

**Review:** reject missing ownership, premature close, dropped cleanup errors, invented idempotency, unbounded waits claimed to be bounded, or callbacks that cannot stop. Ensure integration tests exercise cancellation during factory, ICE fetch, offer/answer and remote attachment, not just final Close.

## 6) Errors, retries, concurrency and technical debt

- Preserve underlying errors with `%w` or `errors.Join` when adding useful context. Inspect status independently of Go error in root REST. Do not normalize VoIP's non-2xx errors into root semantics or vice versa without a migration.
- Do not swallow errors in new code. If an error is intentionally secondary during teardown, report it through the agreed callback/logger/return path and explain the decision. Do not log full token-bearing URLs, credentials, SDP/ICE secrets or private message bodies; even retained underlying errors may contain sensitive URLs.
- Declare which operations can be retried and what makes repetition safe. Socket reconnect is not permission to retry dial, accept, send-message or receipt deletion. Preserve backoff, permanent refusal and ordering unless intentionally changed and tested.
- For concurrency, identify lock ownership, callback reentrancy, cancellation and a bounded admission policy. `On` must stay short and must not call lifecycle methods; `OnAsync` permits blocking work but does not ensure bounded concurrency or ordered completion. New asynchronous work needs an owner and stop path.
- Use context only where cancellation reaches the operation. `OpenAsync` timeout leaves reconnect running; startup cancellation is different from active media lifetime. Never hide these distinctions behind a convenient helper.
- Prefer current dependencies and the standard library. Reject duplicate logic, hidden globals, vague names, unnecessary interfaces, provider logic in consumer policy and unbounded concurrency. A workaround must have a documented scope, reason and removal condition; "temporary" alone is insufficient.
- Respect response buffer ownership: root transport currently returns pooled bytes after release. Do not copy this pattern. A repair needs local lifetime/reuse tests across affected return paths; do not claim a consumer-side copy after return proves safety.

## 7) Compatibility and application customization

Compatibility includes Go signatures, import paths, JSON tags/omission, host selection, receipt URL order, raw response shape, status/error distinctions, timeout/retry policy, callback ordering/payloads, close semantics, audio formats and generation behavior. Tests must cover the observable contract rather than mirror implementation steps. Intentional breaks require an explicit migration plan, consumer impact, updated examples/docs and reviewer approval; do not silently change semantics under a refactor title.

For branded or White Label applications, use this ladder: configuration → consumer composition → existing public extension → narrow reusable option/interface justified by a real consumer → core behavior change → fork. Verify constructor access and ownership before promising substitution. Keep branding, persistence, receipt-processing policy and audio-device integration outside generic SDK internals. There is no built-in White Label subsystem.

The project uses the [MIT License](LICENSE), which permits modification and redistribution subject to its terms. Confirm applicable source and dependency terms separately for a distribution decision. This guide neither changes licensing nor certifies its scope.

## 8) Task-dependent verification

Go 1.23+ is required by the module, which also selects toolchain go1.23.2. Use an appropriate C compiler for bundled codec/race builds. System libopus needs depend on the pinned gopus target. CGo-free `opus.New` returns an error; it is not successful media coverage.

These commands are local checks, **not a CI-backed guarantee**. No local CI or lint configuration is present in this checkout. Dependency downloads, compiler caches and local listening sockets may require environment permissions; do not substitute real account examples if sandboxing blocks them.

```sh
go version
go test ./voip
go test ./opus
go test ./test/voip
go test -race ./...
CGO_ENABLED=0 go test ./...
git diff --check
```

Choose focused packages/tests for a narrow change; for broad call/lifecycle/media changes run the full race and CGo-free variants documented in [voip/README.md](voip/README.md). Use `go test .` after adding root tests. Format changed Go files with `gofmt`; avoid formatting unrelated legacy files. Tests under `./...` compile example packages but do not run their `main` functions. Review any newly introduced test initialization for external side effects.

| Change | Required evidence |
| --- | --- |
| REST endpoint/options/status | Local server or package-local fake: exact method/path/query/JSON, omitted fields, validation/transport/non-2xx outcomes. Root currently has no tests; new transport contracts require coverage. |
| Notifications | Exact timeout omission/value and token/receipt suffix; consumer tests for processing failure, deletion failure and repeats if processing logic is introduced. No inferred retention guarantees. |
| Upload/media host | Temporary file plus local server: fields, filename, MIME, bytes, selected host, malformed types/missing file/read errors; cleanup on failure. |
| Signaling/state/callback | `voip/signaling_test.go`, `voip/calls_test.go`: transient/permanent cases, ordered dispatch, callback reentrancy rules, duplicate/late events, cancel/open/close and resume races. |
| Audio/codec/lifecycle | `voip/audio_test.go`, `voip/rtc_test.go`, `voip/rtc_adapter_test.go`, `opus/opus_test.go`: format/channel/EOF boundaries, cancellation, factory failure, attachment wait and close ownership. |
| Cross-layer call behavior | `test/voip/calls_integration_test.go`: local HTTP/WebSocket/Pion negotiation and bidirectional audio; report skips when codec unavailable. |
| Documentation only | Validate relative links/fences, source/API assertions, instruction size where relevant, sample syntax/compilation and `git diff --check`. No need to run credentialed programs or the full media suite for prose edits. |

Use deterministic fakes/local servers and bounded test deadlines. Test success, failure, boundaries, invalid/repeated transitions, cancellation and cleanup as applicable. Race tests do not prove all schedules, and mock/local tests do not prove provider guarantees. Run live account actions only when explicitly intended and authorized with reviewed credentials/targets; never commit those credentials.

## 9) Mandatory review decision

A change **fails review** for any mandatory corporate-rule violation, unapproved compatibility break, avoidable technical debt, unclear ownership, sensitive data exposure, unverifiable correctness or missing risk-appropriate tests. Reviewers must specifically reject:

- A new/materially changed function over 50 counted body lines, artificial splitting to game the rule, vague contracts or multiple unrelated responsibilities.
- Interface/factory layers without a real consumer, or assumptions that private root request options/public `Bridge` imply unrestricted public injection.
- Root HTTP errors ignored because `err == nil`, automatic notification acknowledgement added silently, or token/receipt path reordered.
- Upload assertions copied without input checks in new code; shared pooled response buffers exposed as owned output.
- Scattered state mutation, lost generation guards, reordered early ICE, retries after permanent refusal or command retries conflated with reconnect.
- Blocking lifecycle calls from `On`, unowned callback goroutines, uncancellable audio waits or disposal before attachment completion.
- Synthetic examples described as shipping APIs, codec skips described as audio passes, or local tests described as live provider validation.

Only an explicit reviewer may approve a narrowly described policy exception. Record rule, scope, reason, risk, validation and removal/follow-up conditions in the review. Compatibility breaks additionally need the migration plan. Existing violations do not authorize new ones, and unrelated cleanup must not obscure the requested change.

A reviewable contribution includes the concrete problem and resulting behavior, consumer impact, changed contracts (or those preserved), exact verification results and material limits. Authors and agents must not claim checks that were not executed.
