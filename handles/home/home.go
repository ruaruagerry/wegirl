package home

import "wegirl/rconst"

const (
	queryServer = "http://www.dbmeinv.com"
)

// tagid 排序重载
type tagid []*rconst.HomeTag

func (a tagid) Len() int {
	return len(a)
}

func (a tagid) Less(i, j int) bool {
	if a[i].CID < a[j].CID {
		return true
	}
	return false
}

func (a tagid) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
