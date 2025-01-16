package schema

type Response[data any] struct {
	StatusCode uint16 `json:"status"`
	Payload    data   `json:"payload"`
}
