package response

type Response struct {
    Metadata Metadata    `json:"metadata"`
    Data     interface{} `json:"data"`
}

func New(metadata Metadata, data interface{}) *Response {
    return &Response{Metadata: metadata, Data: data}
}

type Metadata struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func Success(metadata Metadata, data interface{}) *Response {
    return New(metadata, data)
}

func Error(metadata Metadata, err error) *Response {
    // handle error, maybe put err into metadata
    return New(metadata, nil)
}
