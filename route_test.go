package wstf

func NewFakeRequest() (*Request) {
	var req Request
	req.Params = map[string]string{}
	return &req
}

func NewFakeResponse() (*Response) {
	var res Response
	return &res
}