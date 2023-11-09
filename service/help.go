package service

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/lufeijun/go-tool-aliyun/common"
	"github.com/tjfoc/gmsm/sm3"
)

func HexEncode(raw []byte) *string {
	return common.String(hex.EncodeToString(raw))
}

func Hash(raw []byte, signatureAlgorithm *string) []byte {
	signType := common.StringValue(signatureAlgorithm)
	if signType == "ACS3-HMAC-SHA256" || signType == "ACS3-RSA-SHA256" {
		h := sha256.New()
		h.Write(raw)
		return h.Sum(nil)
	} else if signType == "ACS3-HMAC-SM3" {
		h := sm3.New()
		h.Write(raw)
		return h.Sum(nil)
	}
	return nil
}
