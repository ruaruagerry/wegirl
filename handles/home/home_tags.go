package home

import (
	"encoding/json"
	"sort"
	"wegirl/gconst"
	"wegirl/pb"
	"wegirl/rconst"
	"wegirl/server"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
)

type tagsRsp struct {
	Tags []*rconst.HomeTags `json:"tags"`
}

func tagsHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "home.tagsHandle")

	httpRsp := pb.HTTPResponse{
		Result: proto.Int32(int32(gconst.UnknownError)),
	}
	defer c.WriteJSONRsp(&httpRsp)

	log.Info("tagsHandle enter:")

	conn := c.RedisConn

	// redis multi get
	conn.Send("MULTI")
	conn.Send("HGETALL", rconst.HashHomeTagsConfig)
	redisMDArray, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
		httpRsp.Msg = proto.String("统一获取缓存操作失败")
		log.Errorf("code:%d msg:%s redisMDArray Values err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	hometagsmap, _ := redis.StringMap(redisMDArray[0], nil)

	// do something
	rsptags := []*rconst.HomeTags{}
	for _, v := range hometagsmap {
		tmp := &rconst.HomeTags{}
		err := json.Unmarshal([]byte(v), tmp)
		if err != nil {
			httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
			httpRsp.Msg = proto.String("导航栏标签unmarshal解析失败")
			log.Errorf("code:%d msg:%s tag unmarshal err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
			return
		}

		rsptags = append(rsptags, tmp)
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

	log.Info("tagsHandle rsp, rsp:", string(data))

	return
}
