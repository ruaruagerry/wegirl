package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"wegirl/gconst"
	"wegirl/pb"

	"github.com/golang/protobuf/proto"

	"github.com/garyburd/redigo/redis"
	"github.com/go-xorm/xorm"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// StupidContext simple http request context
type StupidContext struct {
	// RedisConn redis 连接对象
	RedisConn redis.Conn
	RedisPool *redis.Pool
	// DbConn mysql连接对象
	DbConn    *xorm.Engine
	LogDbConn *xorm.Engine
	// MqChannel mq 通道
	MqChannel *Channel
	// Log 日志输出
	Log *log.Entry

	// UserID 当前请求的用户ID
	UserID string

	// Query cached query string
	// 余下的处理代码应该复用该query
	Query url.Values

	// Params httprouter params
	Params httprouter.Params

	// body bytes array, if exist
	Body []byte

	r *http.Request

	W http.ResponseWriter
}

// StupidHandle stupid handle
type StupidHandle func(*StupidContext)

// newReqContext 新建一个context
func newReqContext(r *http.Request, requiredUserID bool) (*StupidContext, int) {
	// parse userID from token, if exist
	userID := ""

	query := r.URL.Query()
	tk := r.Header.Get("Session")
	if requiredUserID {
		var errCode int
		// try to parse token to get userID
		userID, errCode = parseTK(tk)
		if errCode != errTokenSuccess {
			return nil, errCode
		}
	}

	// construct context
	ctx := &StupidContext{}
	ctx.UserID = userID
	ctx.Query = query
	// TODO: with or without IP address?
	ctx.Log = log.WithField("userID", userID)

	return ctx, errTokenSuccess
}

// WriteRsp send protobuf messsage to peer
func (ctx *StupidContext) WriteRsp(msg *pb.HTTPResponse) {
	if msg.GetResult() != int32(gconst.Success) {
		if msg.GetMsg() == "" {
			result := gconst.Error(msg.GetResult())
			msg.Msg = proto.String(result.String())
		}
	}

	bytes, err := proto.Marshal(msg)
	if err != nil {
		ctx.Log.Panic("WriteRsp panic:", err)
	}

	ctx.W.Write(bytes)
}

// WriteJSONRsp send protobuf messsage to peer
func (ctx *StupidContext) WriteJSONRsp(msg *pb.HTTPResponse) {
	if msg.GetResult() != int32(gconst.Success) {
		if msg.GetMsg() == "" {
			result := gconst.Error(msg.GetResult())
			msg.Msg = proto.String(result.String())
		}
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		ctx.Log.Panic("WriteRsp panic:", err)
	}

	ctx.W.Write(bytes)
}

// replyClientWithTokenError reply to client for decoding token failed
func replyClientWithTokenError(w http.ResponseWriter, errCode int) {
	result := int32(gconst.ErrTokenFormat)
	switch errCode {
	case errTokenEmpty:
		result = int32(gconst.ErrTokenEmpty)
	case errTokenDecrypt:
		result = int32(gconst.ErrTokenDecrypt)
	case errTokenFormat:
		result = int32(gconst.ErrTokenFormat)
	case errTokenExpired:
		result = int32(gconst.ErrTokenExpired)
	}

	httpRsp := pb.HTTPResponse{}
	httpRsp.Result = &result

	// 退出函数时发送回复给客户端
	bytes, err := json.Marshal(&httpRsp)
	if err != nil {
		log.Panic("WriteRsp panic:", err)
	}

	w.Write(bytes)
}

// wrapGetHandleInternal 包装 get handle
func wrapGetHandleInternal(handle StupidHandle, requiredUserID bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		conn := pool.Get()
		defer conn.Close()

		ctx, ecode := newReqContext(r, requiredUserID)
		if ecode != errTokenSuccess {
			replyClientWithTokenError(w, ecode)
			return
		}

		ctx.RedisConn = conn
		ctx.RedisPool = pool

		ctx.DbConn = dbConnect
		ctx.LogDbConn = logDbConnect
		ctx.MqChannel = mqChannel
		ctx.Params = params

		ctx.r = r
		ctx.W = w
		handle(ctx)
	}
}

// wrapPostHandleInternal 包装 post handle
func wrapPostHandleInternal(handle StupidHandle, requiredUserID bool) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		conn := pool.Get()
		defer conn.Close()

		ctx, ecode := newReqContext(r, requiredUserID)
		if ecode != errTokenSuccess {
			replyClientWithTokenError(w, ecode)
			return
		}

		ctx.RedisConn = conn
		ctx.RedisPool = pool
		ctx.DbConn = dbConnect
		ctx.LogDbConn = logDbConnect
		ctx.MqChannel = mqChannel
		ctx.Params = params

		// read all body bytes
		// Read body
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		ctx.Body = b
		ctx.r = r
		ctx.W = w

		handle(ctx)
	}
}

// RegisterGetHandle 注册http get handle
func RegisterGetHandle(subPath string, handle StupidHandle) {
	log.Info("RegisterGetHandle:", subPath)
	if subPath[0] != '/' {
		log.Panic("subPath must begin with '/', :", subPath)
	}

	path := rootPath + subPath
	h, _, _ := rootRouter.Lookup("GET", path)
	if h != nil {
		log.Panic("subPath with 'GET' has been register, subPath:", subPath)
	}

	rootRouter.GET(path, wrapGetHandleInternal(handle, true))
}

// RegisterPostHandle 注册http post handle
func RegisterPostHandle(subPath string, handle StupidHandle) {
	log.Info("RegisterPostHandle:", subPath)
	if subPath[0] != '/' {
		log.Panic("RegisterPostHandle subPath must begin with '/', :", subPath)
	}

	path := rootPath + subPath
	h, _, _ := rootRouter.Lookup("POST", path)
	if h != nil {
		log.Panic("RegisterPostHandle subPath with 'POST' has been register, subPath:", subPath)
	}

	rootRouter.POST(path, wrapPostHandleInternal(handle, true))
}

// RegisterPostHandleNoUserID 注册http post handle
func RegisterPostHandleNoUserID(subPath string, handle StupidHandle) {
	log.Info("RegisterPostHandleNoUserID:", subPath)
	if subPath[0] != '/' {
		log.Panic("RegisterPostHandleNoUserID subPath must begin with '/', :", subPath)
	}

	path := rootPath + subPath
	h, _, _ := rootRouter.Lookup("POST", path)
	if h != nil {
		log.Panic("RegisterPostHandleNoUserID subPath with 'POST' has been register, subPath:", subPath)
	}

	rootRouter.POST(path, wrapPostHandleInternal(handle, false))
}

// RegisterGetHandleNoUserID 注册http post handle
func RegisterGetHandleNoUserID(subPath string, handle StupidHandle) {
	log.Info("RegisterGetHandleNoUserID:", subPath)
	if subPath[0] != '/' {
		log.Panic("RegisterGetHandleNoUserID subPath must begin with '/', :", subPath)
	}

	path := rootPath + subPath
	h, _, _ := rootRouter.Lookup("GET", path)
	if h != nil {
		log.Panic("subPath with 'GET' has been register, subPath:", subPath)
	}

	rootRouter.GET(path, wrapGetHandleInternal(handle, false))
}
