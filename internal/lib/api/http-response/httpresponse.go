package httpresponse

import resp "url-shorter/internal/lib/api/response"

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}
