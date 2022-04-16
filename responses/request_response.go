package responses

type RequestResponseMessage struct {
	Message string `json:"message"`
}

type RequestResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}
