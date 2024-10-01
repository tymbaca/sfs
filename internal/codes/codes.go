package codes

type Code = uint64

const (
	Ok         Code = 10
	NotFound   Code = 20
	InvalidReq Code = 21
	Internal   Code = 30
)
