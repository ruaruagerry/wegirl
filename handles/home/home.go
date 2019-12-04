package home

import (
	"math/rand"
	"wegirl/rconst"
)

var (
	droprand *rand.Rand
)

const (
	queryServer = "http://www.dbmeinv.com"
)

var (
	testTags = []string{"男孩", "女孩", "可爱", "唯美"}
	testImgs = []string{
		"http://b-ssl.duitang.com/uploads/item/201709/22/20170922162149_snyk3.jpeg",
		"http://b-ssl.duitang.com/uploads/item/201809/26/20180926162125_vjbwi.jpg",
		"http://hbimg.huabanimg.com/f6ee1b095215b8c2955fd07e56e2739b2671cd3579f0d-5feCNB_fw658",
		"http://b-ssl.duitang.com/uploads/item/201509/04/20150904014041_Lw8Cv.jpeg",
	}
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
