package response

type BytesResponse struct {
	Body       []byte
	StatusCode int
}

func NewBytesResponse(b []byte, sc int) *BytesResponse {
	return &BytesResponse{
		Body:       b,
		StatusCode: sc,
	}
}
