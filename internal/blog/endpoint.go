package blog

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/reyhardy/go-blog/db/scylladb"
	"github.com/reyhardy/go-blog/internal/ssevent"
	datastar "github.com/starfederation/datastar/sdk/go"
	"maragu.dev/gomponents"
)

type endpoint struct {
	svc Servicer
}

type API interface {
	GetHomePage(c echo.Context) error
	GetAddPostFormPage(c echo.Context) error
	GetAddPostFormFragment(c echo.Context) error
	GetEditForm(w http.ResponseWriter, r *http.Request)
	GetAllPostsFragment(c echo.Context) error
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

func (e *endpoint) GetHomePage(c echo.Context) error {
	return c.HTML(http.StatusOK, gomponents.NodeFunc(Home("post").Render).String())
}

func (e *endpoint) GetAddPostFormPage(c echo.Context) error {
	return c.HTML(http.StatusOK, gomponents.NodeFunc(Home("form").Render).String())
}

func (e *endpoint) GetAddPostFormFragment(c echo.Context) error {
	sse := ssevent.NewSSEvent(c.Response().Writer, c.Request())

	err := sse.MergeAllFragments(
		ssevent.Fragment{
			Node: AddForm(),
			Opts: ssevent.FragmentMergeOpts{
				datastar.WithMergeMode(datastar.FragmentMergeModeInner),
				datastar.WithSelectorID("main"),
			},
		},
		ssevent.Fragment{
			Node: NavbarBackButton(),
		},
	)
	if err != nil {
		return err
	}

	if err = sse.ReplaceURL(url.URL{Path: "/add-post"}); err != nil {
		return err
	}

	return nil
}

func (e *endpoint) GetEditForm(w http.ResponseWriter, r *http.Request) {
	for _, post := range posts {
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

			sse := ssevent.NewSSEvent(w, r)

			sse.MergeAllFragments(
				ssevent.Fragment{
					Node: EditForm(post),
					Opts: ssevent.FragmentMergeOpts{
						datastar.WithMergeMode(datastar.FragmentMergeModeOuter),
						datastar.WithSelectorID("main"),
					},
				},
				ssevent.Fragment{
					Node: NavbarBackButton(),
				},
			)

			sse.ReplaceURL(url.URL{Path: r.URL.Path})

			sse.MergeAllSignals(ssevent.Signals{Signal: inputSignal}, ssevent.Signals{Signal: viewSignal})
		}
	}
}

func (e *endpoint) GetAllPostsFragment(c echo.Context) error {
	res, err := e.svc.SelectAll(context.Background(), Keyspace)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("select all post error: %s", err))
	}

	sse := ssevent.NewSSEvent(c.Response().Writer, c.Request())

	err = sse.MergeAllFragments(
		ssevent.Fragment{
			Node: PostList(res),
			Opts: ssevent.FragmentMergeOpts{
				datastar.WithMergeMode(datastar.FragmentMergeModeInner),
				datastar.WithSelectorID("main"),
			},
		},
		ssevent.Fragment{
			Node: NavbarAddPostButton(),
		},
	)
	if err != nil {
		return err
	}

	if err = sse.ReplaceURL(url.URL{Path: "/"}); err != nil {
		return err
	}

	return nil
}

func (e *endpoint) GetPost(w http.ResponseWriter, r *http.Request) {
	for _, post := range posts {
		if post.ID == r.PathValue("id") {
			postParams := &PostParams{
				ID:      post.ID,
				Title:   post.Title,
				Author:  post.Author,
				Content: post.Content,
			}

			res, err := e.svc.Get(r.Context(), Keyspace, postParams)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "get post error: %s", err)
				return
			}

			fmt.Println("post", res)

			viewSignal := ViewSignal{
				View: "posts",
			}

			inputSignal := InputSignal{
				Input: Input{
					Title:   post.Title,
					Author:  post.Author,
					Content: post.Content,
				},
			}

			fmt.Println("query", r.URL.Query())

			sse := ssevent.NewSSEvent(w, r)

			sse.MergeAllSignals(ssevent.Signals{Signal: inputSignal}, ssevent.Signals{Signal: viewSignal})

			sse.MergeAllFragments(
				ssevent.Fragment{
					Node: PostCard(res),
					Opts: ssevent.FragmentMergeOpts{
						datastar.WithSelectorID(PostsId),
						datastar.WithMergeMode(datastar.FragmentMergeModeInner),
					},
				},
				ssevent.Fragment{
					Node: NavbarAddPostButton(),
				},
			)

			sse.ReplaceURL(url.URL{Path: homeURL})
		}
	}

	// postParams := &PostParams{
	// 	ID: r.PathValue("id"),
	// }

	// res, err := e.svc.Get(r.Context(), Keyspace, postParams)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintf(w, "get post error: %s", err)
	// 	return
	// }

	// viewSignal := ViewSignal{
	// 	View: "posts",
	// }

	// fmt.Println("query", r.URL.Query())

	// sse := ssevent.NewSSEvent(w, r)

	// sse.MergeAllFragments(
	// 	ssevent.Fragment{
	// 		Node: PostCard(res),
	// 		Opts: ssevent.FragmentMergeOpts{
	// 			datastar.WithSelectorID(PostsId),
	// 			datastar.WithMergeMode(datastar.FragmentMergeModeInner),
	// 		},
	// 	},
	// 	ssevent.Fragment{
	// 		Node: NavbarAddPostButton(),
	// 	},
	// )

	// sse.MergeAllSignals(ssevent.Signals{Signal: viewSignal})

	// sse.ReplaceURL(url.URL{Path: homeURL})
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

	sse := ssevent.NewSSEvent(w, r)

	sse.MergeAllFragments(
		ssevent.Fragment{
			Node: PostCard(res),
			Opts: ssevent.FragmentMergeOpts{
				datastar.WithMergePrepend(),
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

	sse := ssevent.NewSSEvent(w, r)

	sse.MergeAllFragments(
		ssevent.Fragment{
			Node: PostCard(res),
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

	err := e.svc.Delete(r.Context(), Keyspace, postParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "delete post error: %s", err)
		return
	}

	sse := ssevent.NewSSEvent(w, r)
	sse.RemoveAllFragments(fmt.Sprintf("#post-%s", postParams.ID))
}
