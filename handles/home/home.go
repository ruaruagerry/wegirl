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
		[]string{"https://pic4.zhimg.com/50/v2-6f3d900004a053b38736d70a95e1ff4d_hd.jpg", "https://pic1.zhimg.com/80/v2-e127049dbe074821174cfd16954ac00d_hd.jpg", "https://img.zcool.cn/community/019e235c0a3f2da801209252c0497f.jpg@1280w_1l_2o_100sh.jpg"},
		[]string{"https://pic1.zhimg.com/80/v2-79dde9e02680f1594b159d3b00058d6f_hd.jpg", "https://pic4.zhimg.com/80/v2-516691a7fc4f536f5587146b36376314_hd.jpg", "https://img.zcool.cn/community/01045058a578d6a801219c77f4a434.jpg"},
		[]string{"https://pic1.zhimg.com/80/v2-5dee6f4a28ba3252e03c2d2268eae93f_hd.jpg", "https://pic2.zhimg.com/80/v2-176bb323efbf92d0e139d43cc38590ff_hd.jpg", "https://pic1.zhimg.com/80/v2-47534989141c9418e6a7992ebcbed3e9_hd.jpg"},
		[]string{"https://pic1.zhimg.com/80/v2-b4eebad17f6637f4c231ab34dae45473_hd.jpg", "https://pic4.zhimg.com/80/v2-3d78919445f889a0bc8903ed40abb100_hd.jpg", "https://pic4.zhimg.com/80/v2-63ed0eaf5efd35ae2aa400ea4f2c2b17_hd.jpg"},
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
