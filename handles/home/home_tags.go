package home

import (
	"encoding/json"
	"sort"
	"strings"
	"wegirl/gconst"
	"wegirl/pb"
	"wegirl/rconst"
	"wegirl/server"

	"github.com/crufter/goquery"
	"github.com/golang/protobuf/proto"
)

type tagsRsp struct {
	Tags []*rconst.HomeTag `json:"tags"`
}

func tagsHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "home.tagsHandle")

	httpRsp := pb.HTTPResponse{
		Result: proto.Int32(int32(gconst.UnknownError)),
	}
	defer c.WriteJSONRsp(&httpRsp)

	log.Info("tagsHandle enter:")

	// do something
	x, err := goquery.ParseUrl(queryServer)
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrHTTP))
		httpRsp.Msg = proto.String("gpquery失败")
		log.Errorf("code:%d msg:%s goquery parseurl err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	rsptags := []*rconst.HomeTag{}
	xnodes := x.Find(".panel-heading ul.nav li a")
	for i := range xnodes {
		title := xnodes.Eq(i).Html()
		href := xnodes.Eq(i).Attr("href")
		hrefs := strings.Split(href, "=")
		if len(hrefs) == 2 {
			tmp := &rconst.HomeTag{
				CID:   hrefs[1],
				Title: title,
			}

			rsptags = append(rsptags, tmp)
		}
	}

	sort.Stable(tagid(rsptags))

	// rsp
	rsp := &tagsRsp{
		Tags: rsptags,
	}
	data, err := json.Marshal(rsp)
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("返回信息marshal解析失败")
		log.Errorf("code:%d msg:%s json marshal err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}
	httpRsp.Result = proto.Int32(int32(gconst.Success))
	httpRsp.Data = data

	// log.Info("tagsHandle rsp, rsp:", string(data))

	return
}
