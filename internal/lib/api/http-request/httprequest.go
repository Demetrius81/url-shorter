package httprequest

type Request struct {
	URL   string `json:"url" validate: "required,url"`
	Alias string `json:"alias,omitempty"`
}
