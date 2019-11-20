package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

/* 微信登录相关 */
const (
	WXLoginURL = "https://api.weixin.qq.com/sns/jscode2session"
	// WXAppID     = "wxe9e3a15f9de9a419"
	// WXAppSecret = "32dfe47807b5961ff26d57c63a2bb936"
	WXAppID     = "wxe6dbc93d12e0240a"
	WXAppSecret = "7545cb3d76732c85f3d6a23875a51fc5"
)

// LoadAccessTokenReply 微信拉取access token回复
type LoadAccessTokenReply struct {
	SessionKey string `json:"session_key"`
	OpenID     string `json:"openid"`
	UnionID    string `json:"unionid"`

	ErrorCode int    `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}

// WeiXinUserPlusInfo 微信用户信息
type WeiXinUserPlusInfo struct {
	OpenID    string `json:"openId"`
	NickName  string `json:"nickName"`
	Gender    int32  `json:"gender"`
	AvatarURL string `json:"avatarUrl"`
	UnionID   string `json:"unionId"`
}

// WeixinGetUserInfo 微信获取用户信息
func WeixinGetUserInfo(wechatCode string) (*LoadAccessTokenReply, error) {
	log := logrus.WithField("func", "WeixinGetUserInfo")

	urlGetAccessToken := fmt.Sprintf("%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", WXLoginURL, WXAppID, WXAppSecret, wechatCode)

	reply := &LoadAccessTokenReply{}
	err := loadDataUseHTTPGet(urlGetAccessToken, reply)
	if err != nil {
		log.Error("WeixinGetUserInfo, err :", err)
		return nil, err
	}

	return reply, nil
}

// GetWeiXinPlusUserInfo 获取微信更多信息
func GetWeiXinPlusUserInfo(sessionkey string, encrypteddata string, iv string) (*WeiXinUserPlusInfo, error) {
	log := logrus.WithField("func", "GetWeiXinPlusUserInfo")

	log.Infof("sessionkey:%s encrypteddata:%s iv:%s", sessionkey, encrypteddata, iv)

	skey, err := base64.StdEncoding.DecodeString(sessionkey)
	if err != nil {
		return nil, err
	}

	sdata, err := base64.StdEncoding.DecodeString(encrypteddata)
	if err != nil {
		return nil, err
	}

	siv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}

	databyte := PswDecrypt(string(sdata), string(skey), string(siv))
	log.Info("GetWeiXinPlusUserInfo data:", string(databyte))

	userplusinfo := &WeiXinUserPlusInfo{}
	err = json.Unmarshal(databyte, &userplusinfo)
	if err != nil {
		return nil, err
	}

	return userplusinfo, nil
}

// PswDecrypt TODO
func PswDecrypt(src string, skey string, siv string) []byte {
	key := []byte(skey)
	iv := []byte(siv)
	data := []byte(src)

	var err error

	origData, err := Aes128Decrypt(data, key, iv)
	if err != nil {
		panic(err)
	}
	return origData
}

// Aes128Decrypt TODO
func Aes128Decrypt(crypted, key []byte, IV []byte) ([]byte, error) {
	if key == nil || len(key) != 16 {
		return nil, nil
	}
	if IV != nil && len(IV) != 16 {
		return nil, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, IV[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

// PKCS5UnPadding TODO
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func loadDataUseHTTPGet(url string, jsonStruct interface{}) error {
	log := logrus.WithField("func", "loadDataUseHTTPGet")

	resp, err := http.Get(url)
	if err != nil {
		log.Error("loadDataUseHTTPGet, err :", err)
		return err
	}
	log.Info("url:", url)
	// 确保body关闭
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(jsonStruct)
	log.Info("jsonStruct:", jsonStruct)
	return err
}
