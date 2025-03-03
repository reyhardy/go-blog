package blog

type Input struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

type InputSignal struct {
	Input Input `json:"input"`
}

type ViewSignal struct {
	View string `json:"view"`
}
