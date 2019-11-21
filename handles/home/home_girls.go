package home

import (
	"encoding/json"
	"strings"
	"wegirl/gconst"
	"wegirl/pb"
	"wegirl/rconst"
	"wegirl/server"

	"github.com/crufter/goquery"
	"github.com/golang/protobuf/proto"
)

type girlsReq struct {
	CID  string `json:"cid"`
	Page int    `json:"page"`
}

type girlsRsp struct {
	Imgs []*rconst.HomeImg `json:"imgs"`
}

func girlsHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "home.girlsHandle")

	httpRsp := pb.HTTPResponse{
		Result: proto.Int32(int32(gconst.UnknownError)),
	}
	defer c.WriteJSONRsp(&httpRsp)

	// req
	req := &girlsReq{}
	if err := json.Unmarshal(c.Body, req); err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("请求信息解析失败")
		log.Errorf("code:%d msg:%s json Unmarshal err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	log.Info("tagsHandle enter:")

	// do something
	x, err := goquery.ParseUrl(queryServer)
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrHTTP))
		httpRsp.Msg = proto.String("gpquery失败")
		log.Errorf("code:%d msg:%s goquery parseurl err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	rspimgs := []*rconst.HomeImg{}
	item := x.Find(".panel-body ul.thumbnails li.span3 .thumbnail .img_single a")
	subitem := x.Find(".panel-body ul.thumbnails li.span3 .thumbnail .img_single a img")
	for i := range item {
		href := item.Eq(i).Attr("href")
		src := subitem.Eq(i).Attr("src")
		title := subitem.Eq(i).Attr("title")
		large := ""
		small := ""
		if src != "" {
			large = strings.Replace(src, "bmiddle", "large", -1)
			small = strings.Replace(src, "bmiddle", "small", -1)
		}

		tmp := &rconst.HomeImg{
			Title: title,
			Href:  href,
			Large: large,
			Thumb: src,
			Small: small,
		}

		rspimgs = append(rspimgs, tmp)
	}

	// rsp
	rsp := &girlsRsp{
		Imgs: rspimgs,
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

	// log.Info("girlsHandle rsp, rsp:", string(data))

	return
}
