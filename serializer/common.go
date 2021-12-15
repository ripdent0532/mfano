package serializer

type Response struct {
	Code  int         `json:"code"`
	Msg   interface{} `json:"msg,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
	Page  Page        `json:"page,omitempty"`
}

type Page struct {
	Total int64 `json:"total"`
	Num   int   `json:"num"`
	Size  int   `json:"size"`
}
