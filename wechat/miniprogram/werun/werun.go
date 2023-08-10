package werun

import (
	"encoding/json"

	"github.com/lufeijun/go-tool-wechat/wechat/miniprogram/encryptor"
	"github.com/lufeijun/go-tool-wechat/wechat/miniprogram/wcontext"
)

// WeRun 微信运动
type WeRun struct {
	*wcontext.Context
}

// Data 微信运动数据
type Data struct {
	StepInfoList []struct {
		Timestamp int `json:"timestamp"`
		Step      int `json:"step"`
	} `json:"stepInfoList"`
}

// NewWeRun 实例化
func NewWeRun(ctx *wcontext.Context) *WeRun {
	return &WeRun{Context: ctx}
}

// GetWeRunData 解密数据
func (werun *WeRun) GetWeRunData(sessionKey, encryptedData, iv string) (*Data, error) {
	cipherText, err := encryptor.GetCipherText(sessionKey, encryptedData, iv)
	if err != nil {
		return nil, err
	}
	var weRunData Data
	err = json.Unmarshal(cipherText, &weRunData)
	if err != nil {
		return nil, err
	}
	return &weRunData, nil
}
