package blog

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/reyhardy/go-blog/db/scylladb"
	"github.com/reyhardy/go-blog/internal/ssevent"
	datastar "github.com/starfederation/datastar/sdk/go"
)

type endpoint struct {
	svc Servicer
}

type API interface {
	GetHome(w http.ResponseWriter, r *http.Request)
	GetAddForm(w http.ResponseWriter, r *http.Request)
	GetEditForm(w http.ResponseWriter, r *http.Request)
	GetPost(w http.ResponseWriter, r *http.Request)
	AddPost(w http.ResponseWriter, r *http.Request)
	EditPost(w http.ResponseWriter, r *http.Request)
	DeletePost(w http.ResponseWriter, r *http.Request)
}

func NewAPI(session scylladb.Client) API {
	return &endpoint{
		svc: NewService(session),
	}
}

const homeURL string = "/home"

func (e *endpoint) GetHome(w http.ResponseWriter, r *http.Request) {
	Home().Render(w)
}

func (e *endpoint) GetAddForm(w http.ResponseWriter, r *http.Request) {
	inputSignal := InputSignal{
		Input: Input{
			Title:   "",
			Author:  "",
			Content: "",
		},
	}

	viewSignal := ViewSignal{
		View: "form",
	}

	fmt.Println("add form URL: \n", r.URL.Query())

	sse := ssevent.NewSSEvent(w, r)

	sse.MergeAllFragments(
		ssevent.Fragment{
			Node: AddForm(),
			Opts: ssevent.FragmentMergeOpts{
				datastar.WithMergeMode(datastar.FragmentMergeModeInner),
				datastar.WithSelectorID(FormId),
			},
		},
		ssevent.Fragment{
			Node: NavbarBackButton(),
		},
	)

	sse.ReplaceURL(url.URL{Path: r.URL.Path})

	sse.MergeAllSignals(ssevent.Signals{Signal: inputSignal}, ssevent.Signals{Signal: viewSignal})
}

func (e *endpoint) GetEditForm(w http.ResponseWriter, r *http.Request) {
	for _, post := range postList {
		if post.ID == r.PathValue("id") {

			inputSignal := InputSignal{
				Input: Input{
					Title:   post.Title,
					Author:  post.Author,
					Content: post.Content,
				},
			}

			viewSignal := ViewSignal{
				View: "form",
			}

			fmt.Println("edit form URL: \n", r.URL.Query())

			sse := ssevent.NewSSEvent(w, r)

			sse.MergeAllFragments(
				ssevent.Fragment{
					Node: EditForm(post),
					Opts: ssevent.FragmentMergeOpts{
						datastar.WithMergeMode(datastar.FragmentMergeModeInner),
						datastar.WithSelectorID(FormId),
					},
				},
				ssevent.Fragment{
					Node: NavbarBackButton(),
					Opts: ssevent.FragmentMergeOpts{
						datastar.WithSelectorID(NavbarButtonId),
					},
				},
			)

			sse.ReplaceURL(url.URL{Path: r.URL.Path})

			sse.MergeAllSignals(ssevent.Signals{Signal: inputSignal}, ssevent.Signals{Signal: viewSignal})
		}
	}
}

func (e *endpoint) GetPost(w http.ResponseWriter, r *http.Request) {
	res, err := e.svc.SelectAll(context.Background(), Keyspace)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "select all post error: %s", err)
		return
	}

	fmt.Println("get post url: \n", r.URL)

	viewSignal := ViewSignal{
		View: "posts",
	}

	sse := ssevent.NewSSEvent(w, r)

	sse.MergeAllFragments(
		ssevent.Fragment{
			Node: PostList(res),
			Opts: ssevent.FragmentMergeOpts{
				datastar.WithMergeMode(datastar.FragmentMergeModeInner),
				datastar.WithSelectorID(PostsId),
			},
		},
		ssevent.Fragment{
			Node: NavbarAddPostButton(),
		},
	)

	sse.ReplaceURL(url.URL{Path: homeURL})

	sse.MergeAllSignals(ssevent.Signals{Signal: viewSignal})
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

	viewSignal := ViewSignal{
		View: "posts",
	}

	fmt.Println("add post url: \n", r.URL)

	sse := ssevent.NewSSEvent(w, r)

	sse.MergeAllFragments(
		ssevent.Fragment{
			Node: PostList(res),
			Opts: ssevent.FragmentMergeOpts{
				datastar.WithMergeMode(datastar.FragmentMergeModeInner),
				datastar.WithSelectorID(PostsId),
			},
		},
		ssevent.Fragment{
			Node: NavbarAddPostButton(),
		},
	)

	sse.MergeAllSignals(ssevent.Signals{Signal: viewSignal})

	sse.ReplaceURL(url.URL{Path: homeURL})
}

func (e *endpoint) EditPost(w http.ResponseWriter, r *http.Request) {
	postParams := &PostParams{
		ID:      r.PathValue("id"),
		Title:   r.FormValue("title"),
		Author:  r.FormValue("author"),
		Content: r.FormValue("content"),
	}

	inputSignal := InputSignal{
		Input: Input{
			Title:   postParams.Title,
			Author:  postParams.Author,
			Content: postParams.Content,
		},
	}

	viewSignal := ViewSignal{
		View: "posts",
	}

	res, err := e.svc.Update(r.Context(), Keyspace, postParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "update post error: %s", err)
		return
	}

	fmt.Println("edit post url: \n", r.URL)

	sse := ssevent.NewSSEvent(w, r)

	sse.MergeAllFragments(
		ssevent.Fragment{
			Node: PostList(res),
			Opts: ssevent.FragmentMergeOpts{
				datastar.WithMergeMode(datastar.FragmentMergeModeInner),
				datastar.WithSelectorID(PostsId),
			},
		},
		ssevent.Fragment{
			Node: NavbarAddPostButton(),
		},
	)

	sse.MergeAllSignals(ssevent.Signals{Signal: inputSignal}, ssevent.Signals{Signal: viewSignal})

	sse.ReplaceURL(url.URL{Path: homeURL})
}

func (e *endpoint) DeletePost(w http.ResponseWriter, r *http.Request) {
	postParams := &PostParams{
		ID: r.PathValue("id"),
	}

	res, err := e.svc.Delete(r.Context(), Keyspace, postParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "delete post error: %s", err)
		return
	}

	fmt.Println("delete post url: \n", r.URL)

	sse := ssevent.NewSSEvent(w, r)

	sse.MergeAllFragments(
		ssevent.Fragment{
			Node: PostList(res),
			Opts: ssevent.FragmentMergeOpts{
				datastar.WithMergeMode(datastar.FragmentMergeModeInner),
				datastar.WithSelectorID(PostsId),
			},
		},
	)
}
