package ssevent

import (
	"net/http"
	"net/url"

	datastar "github.com/starfederation/datastar/sdk/go"
	"maragu.dev/gomponents"
)

type serverSentEvent struct {
	sse *datastar.ServerSentEventGenerator
}

type SSE interface {
	MergeAllFragments(fragments ...Fragment) error
	MergeAllSignals(signals ...Signals) error
	RemoveAllFragments(selector string) error
	ReplaceURL(u url.URL) error
	Redirect(u url.URL) error
}

func NewSSEvent(w http.ResponseWriter, r *http.Request) SSE {
	sse := datastar.NewSSE(w, r)
	return &serverSentEvent{sse}
}

type FragmentMergeOpts []datastar.MergeFragmentOption
type SignalMergeOpts []datastar.MergeSignalsOption

// type fragmentRemoveOpts []datastar.RemoveFragmentsOption

type Fragment struct {
	Node gomponents.Node
	Opts FragmentMergeOpts
}

type Signals struct {
	Signal any
	Opts   SignalMergeOpts
}

func (f *Fragment) NodeString() string {
	return gomponents.NodeFunc(f.Node.Render).String()
}

func (sse *serverSentEvent) MergeAllFragments(fragments ...Fragment) error {
	for _, f := range fragments {
		err := sse.sse.MergeFragments(f.NodeString(), f.Opts...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sse *serverSentEvent) MergeAllSignals(signals ...Signals) error {
	for _, s := range signals {
		err := sse.sse.MarshalAndMergeSignals(s.Signal, s.Opts...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sse *serverSentEvent) RemoveAllFragments(selector string) error {
	return sse.sse.RemoveFragments(selector)
}

func (sse *serverSentEvent) ReplaceURL(u url.URL) error {
	return sse.sse.ReplaceURL(u)
}

func (sse *serverSentEvent) Redirect(u url.URL) error {
	return sse.sse.Redirect(u.Path)
}
