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
	testImgs = [][]string{
		[]string{"https://b-ssl.duitang.com/uploads/item/201709/22/20170922162149_snyk3.jpeg", "https://img3.duitang.com/uploads/item/201602/18/20160218003123_aLMyv.jpeg", "https://img.zcool.cn/community/019e235c0a3f2da801209252c0497f.jpg@1280w_1l_2o_100sh.jpg"},
		[]string{"https://b-ssl.duitang.com/uploads/item/201809/26/20180926162125_vjbwi.jpg", "https://b-ssl.duitang.com/uploads/item/201607/25/20160725102949_2earM.jpeg", "https://img.zcool.cn/community/01045058a578d6a801219c77f4a434.jpg"},
		[]string{"https://hbimg.huabanimg.com/f6ee1b095215b8c2955fd07e56e2739b2671cd3579f0d-5feCNB_fw658", "https://img5.imgtn.bdimg.com/it/u=2717062052,3164034025&fm=26&gp=0.jpg", "https://cdn.duitang.com/uploads/item/201602/23/20160223124339_d2NkX.jpeg"},
		[]string{"https://b-ssl.duitang.com/uploads/item/201509/04/20150904014041_Lw8Cv.jpeg", "https://b-ssl.duitang.com/uploads/item/201602/10/20160210211239_JCnsw.jpeg", "https://cdn.duitang.com/uploads/item/201410/26/20141026191422_yEKyd.thumb.700_0.jpeg"},
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
