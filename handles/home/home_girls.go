package home

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"wegirl/gconst"
	"wegirl/pb"
	"wegirl/rconst"
	"wegirl/server"
	"wegirl/servercfg"

	"github.com/crufter/goquery"
	"github.com/garyburd/redigo/redis"
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

	log.Info("girlsHandle enter:", string(c.Body))

	conn := c.RedisConn

	rspimgs := []*rconst.HomeImg{}
	if servercfg.ForTestOnly {
		if req.Page == 1 {
			cidint, _ := strconv.Atoi(req.CID)

			for _, v := range testImgs[cidint] {
				tmp := &rconst.HomeImg{
					Large: v,
					Thumb: v,
					Small: v,
				}

				rspimgs = append(rspimgs, tmp)
			}
		}
	} else {
		// redis multi get
		conn.Send("MULTI")
		conn.Send("SRANDMEMBER", rconst.SetHomeGoodImages)
		redisMDArray, err := redis.Values(conn.Do("EXEC"))
		if err != nil {
			httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
			httpRsp.Msg = proto.String("统一获取缓存操作失败")
			log.Errorf("code:%d msg:%s redisMDArray Values err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
			return
		}

		goodimage, _ := redis.String(redisMDArray[0], nil)

		// do something
		queryurl := fmt.Sprintf("%s?cid=%s&page=%d", queryServer, req.CID, req.Page)
		x, err := goquery.ParseUrl(queryurl)
		if err != nil {
			httpRsp.Result = proto.Int32(int32(gconst.ErrHTTP))
			httpRsp.Msg = proto.String("gpquery失败")
			log.Errorf("code:%d msg:%s goquery parseurl err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
			return
		}

		item := x.Find(".panel-body ul.thumbnails li.span3 .thumbnail .img_single a")
		subitem := x.Find(".panel-body ul.thumbnails li.span3 .thumbnail .img_single a img")
		randidx := droprand.Intn(len(item) / 2)
		for i := range item {
			filter := false
			href := item.Eq(i).Attr("href")
			src := subitem.Eq(i).Attr("src")
			if i == randidx && goodimage != "" {
				src = goodimage
				filter = true
			}
			title := subitem.Eq(i).Attr("title")
			large := ""
			small := ""
			if src != "" {
				large = strings.Replace(src, "bmiddle", "large", -1)
				small = strings.Replace(src, "bmiddle", "small", -1)
			}

			tmp := &rconst.HomeImg{
				Title:  url.QueryEscape(title),
				Href:   href,
				Large:  large,
				Thumb:  src,
				Small:  small,
				Filter: filter,
			}

			rspimgs = append(rspimgs, tmp)
		}
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
