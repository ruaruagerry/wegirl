package demo

import (
	"wegirl/gconst"
	"wegirl/pb"
	"wegirl/server"

	"github.com/golang/protobuf/proto"
)

func helloHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "demo.helloHandle")

	httpRsp := pb.HTTPResponse{
		Result: proto.Int32(int32(gconst.UnknownError)),
	}
	defer c.WriteJSONRsp(&httpRsp)

	// req
	// req := &pb.HelloReq{}
	// if err := json.Unmarshal(c.Body, req); err != nil {
	// 	httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
	// 	httpRsp.Msg = proto.String("请求信息解析失败")
	// 	log.Errorf("code:%d msg:%s json Unmarshal err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
	// 	return
	// }

	// log.Info("helloHandle enter, req:", string(c.Body))

	// conn := c.RedisConn
	// playerid := c.UserID

	// // redis multi get
	// conn.Send("MULTI")
	// redisMDArray, err := redis.Values(conn.Do("EXEC"))
	// if err != nil {
	// 	httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
	// 	httpRsp.Msg = proto.String("统一获取缓存操作失败")
	// 	log.Errorf("code:%d msg:%s redisMDArray Values err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
	// 	return
	// }

	// // do something

	// // redis multi set
	// conn.Send("MULTI")
	// _, err = conn.Do("EXEC")
	// if err != nil {
	// 	httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
	// 	httpRsp.Msg = proto.String("统一存储缓存操作失败")
	// 	log.Errorf("code:%d msg:%s exec err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
	// 	return
	// }

	// // rsp
	// rsp := &pb.HelloRsp{}
	// data, err := json.Marshal(rsp)
	// if err != nil {
	// 	httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
	// 	httpRsp.Msg = proto.String("返回信息marshal解析失败")
	// 	log.Errorf("code:%d msg:%s json marshal err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
	// 	return
	// }
	httpRsp.Result = proto.Int32(int32(gconst.Success))
	// httpRsp.Data = data

	// log.Info("helloHandle rsp, rsp:", string(data))

	return
}
