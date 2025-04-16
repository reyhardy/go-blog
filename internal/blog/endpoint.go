package blog

import (
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
	AddPostEcho(c echo.Context) error
	GetPost(w http.ResponseWriter, r *http.Request)
	EditPost(w http.ResponseWriter, r *http.Request)
	// DeletePost(w http.ResponseWriter, r *http.Request)
	DeletePostEcho(c echo.Context) error
}

func NewAPI(session scylladb.Client) API {
	return &endpoint{
		svc: NewService(session),
	}
}

const homeURL string = "/"

func (e *endpoint) GetHomePage(c echo.Context) error {
	return c.HTML(http.StatusOK, gomponents.NodeFunc(Home(PostsId).Render).String())
}

func (e *endpoint) GetAddPostFormPage(c echo.Context) error {
	return c.HTML(http.StatusOK, gomponents.NodeFunc(Home(FormId).Render).String())
}

func (e *endpoint) GetAddPostFormFragment(c echo.Context) error {
	viewSignal := ViewSignal{
		View: FormId,
	}

	sse := ssevent.NewSSEvent(c.Response().Writer, c.Request())

	err := sse.MergeAllFragments(
		ssevent.Fragment{
			Node: FormAddPost(),
			Opts: ssevent.FragmentMergeOpts{
				datastar.WithMergeInner(),
				datastar.WithSelectorID(FormId),
			},
		},
		ssevent.Fragment{
			Node: NavbarBackButton(),
		},
	)
	if err != nil {
		return err
	}

	sse.MergeAllSignals(ssevent.Signals{Signal: viewSignal})

	if err = sse.ReplaceURL(url.URL{Path: "/form/add-post"}); err != nil {
		return err
	}

	return c.String(http.StatusOK, "add post form fragment loaded")
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
						datastar.WithMergeInner(),
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
	res, err := e.svc.SelectAll(c.Request().Context(), Keyspace)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("select all post error: %s", err))
	}

	viewSignal := ViewSignal{
		View: PostsId,
	}

	sse := ssevent.NewSSEvent(c.Response().Writer, c.Request())

	err = sse.MergeAllFragments(
		ssevent.Fragment{
			Node: PostList(res),
			Opts: ssevent.FragmentMergeOpts{
				datastar.WithMergeInner(),
				datastar.WithSelectorID(PostsId),
			},
		},
		ssevent.Fragment{
			Node: NavbarAddPostButton(),
		},
	)
	if err != nil {
		return err
	}

	sse.MergeAllSignals(ssevent.Signals{Signal: viewSignal})

	if err = sse.ReplaceURL(url.URL{Path: homeURL}); err != nil {
		return err
	}

	return c.String(http.StatusOK, "get all post fragment loaded")
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

func (e *endpoint) AddPostEcho(c echo.Context) error {
	postParams := &PostParams{
		Title:   c.FormValue("title"),
		Author:  c.FormValue("author"),
		Content: c.FormValue("content"),
	}

	_, err := e.svc.Add(c.Request().Context(), Keyspace, postParams)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("add post error: %s", err))
	}

	sse := ssevent.NewSSEvent(c.Response(), c.Request())

	// err = sse.MergeAllFragments(
	// 	ssevent.Fragment{
	// 		Node: PostCard(res),
	// 		Opts: ssevent.FragmentMergeOpts{
	// 			datastar.WithSelectorID(PostsId),
	// 			datastar.WithMergePrepend(),
	// 		},
	// 	},
	// 	ssevent.Fragment{
	// 		Node: NavbarAddPostButton(),
	// 	},
	// )
	// if err != nil {
	// 	return err
	// }

	// err = sse.ReplaceURL(url.URL{Path: homeURL})
	// if err != nil {
	// 	return err
	// }

	sse.Redirect(url.URL{Path: homeURL})

	return c.String(http.StatusOK, fmt.Sprintf("successfully added %s to posts", postParams.Title))
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

func (e *endpoint) DeletePostEcho(c echo.Context) error {
	postParams := &PostParams{
		ID: c.Param("id"),
	}

	err := e.svc.Delete(c.Request().Context(), Keyspace, postParams)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("delete post error: %s", err))
	}

	sse := ssevent.NewSSEvent(c.Response(), c.Request())
	err = sse.RemoveAllFragments(fmt.Sprintf("#post-%s", postParams.ID))
	if err != nil {
		return c.String(http.StatusBadRequest, "error remove fragment")
	}

	fmt.Printf("post id: %s \n", postParams.ID)

	return c.NoContent(http.StatusOK)
}
