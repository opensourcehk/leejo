package data

type Resp struct {
	Status  string      `json:"status"`
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

func (r *Resp) GetCode() int {
	return r.Code
}
