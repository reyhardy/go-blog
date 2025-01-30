package blog

import (
	"context"
	"fmt"
	"net/http"

	"github.com/reyhardy/go-blog/db/scylladb"
	datastar "github.com/starfederation/datastar/sdk/go"
)

type endpoint struct {
	svc servicer
}

type API interface {
	GetHome(w http.ResponseWriter, r *http.Request)
	GetPost(w http.ResponseWriter, r *http.Request)
	AddPost(w http.ResponseWriter, r *http.Request)
	DeletePost(w http.ResponseWriter, r *http.Request)
}

func NewAPI(session scylladb.Client) API {
	return &endpoint{
		svc: newService(session),
	}
}

func (e *endpoint) GetHome(w http.ResponseWriter, r *http.Request) {
	Home().Render(w)
}

func (e *endpoint) GetPost(w http.ResponseWriter, r *http.Request) {
	res, err := e.svc.SelectAll(context.Background(), Keyspace)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "select all post error: %s", err)
		return
	}

	sse := datastar.NewSSE(w, r)

	err = sse.MergeFragments(RenderPostList(res).String())
	if err != nil {
		sse.ConsoleError(err)
	}
}

func (e *endpoint) AddPost(w http.ResponseWriter, r *http.Request) {
	postParams := &PostParams{
		Title:   r.FormValue("title"),
		Author:  r.FormValue("author"),
		Content: r.FormValue("content"),
	}

	res, err := e.svc.Add(r.Context(), Keyspace, postParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "add post error: %s", err)
		return
	}

	sse := datastar.NewSSE(w, r)
	err = sse.MergeFragments(RenderPostCard(res).String(), datastar.WithMergeAppend(), datastar.WithSelectorID("postlist"))
	if err != nil {
		sse.ConsoleError(err)
	}
}

func (e *endpoint) DeletePost(w http.ResponseWriter, r *http.Request) {
	postParams := &PostParams{
		ID: r.PathValue("id"),
	}

	err := e.svc.Delete(r.Context(), Keyspace, postParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "delete post error: %s", err)
		return
	}

	sse := datastar.NewSSE(w, r)
	err = sse.RemoveFragments(fmt.Sprintf("#post-%v", postParams.ID))
	if err != nil {
		sse.ConsoleError(err)
	}
}
