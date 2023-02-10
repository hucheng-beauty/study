package request

type Common struct {
	Action string `json:"Action"`
}

type PageInfo struct {
	Offset int `json:"Offset"`
	Limit  int `json:"Limit"`
}

/*
	TODO add request list
*/
