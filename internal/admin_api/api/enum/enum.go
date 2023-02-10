package enum

var EmptyEnumRequest = map[string]string{
	"Empty": "",
}

var EmptyEnumResponse = map[string]string{
	"": "Empty",
}

// RtCodeEnumRequest 1-成功,0-失败
var RtCodeEnumRequest = map[string]int{
	"Success": 1,
	"Failed":  0,
}

var RtCodeEnumResponse = map[int]string{
	1: "Success",
	0: "Failed",
}
