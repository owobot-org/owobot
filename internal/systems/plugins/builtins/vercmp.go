package builtins

import "go.elara.ws/vercmp"

type vercmpAPI struct{}

func (vercmpAPI) Newer(v1, v2 string) bool {
	return vercmp.Compare(v1, v2) == 1
}

func (vercmpAPI) Older(v1, v2 string) bool {
	return vercmp.Compare(v1, v2) == -1
}

func (vercmpAPI) Equal(v1, v2 string) bool {
	return vercmp.Compare(v1, v2) == 0
}

func (vercmpAPI) Compare(v1, v2 string) int {
	return vercmp.Compare(v1, v2)
}
