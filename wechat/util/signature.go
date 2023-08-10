package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"sort"
	"strings"
)

// 微信签名算法方式
const (
	SignTypeMD5        = `MD5`
	SignTypeHMACSHA256 = `HMAC-SHA256`
)

// ECB provides confidentiality by assigning a fixed ciphertext block to each plaintext block.
// See NIST SP 800-38A, pp 08-09
// reference: https://codereview.appspot.com/7860047/patch/23001/24001
type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

// NewECBDecryptor returns a BlockMode which decrypts in electronic code book mode, using the given Block.
// ECBDecryptor -
type ECBDecryptor ecb

func NewECBDecryptor(b cipher.Block) cipher.BlockMode {
	return (*ECBDecryptor)(newECB(b))
}

// BlockSize implement BlockMode.BlockSize
func (x *ECBDecryptor) BlockSize() int {
	return x.blockSize
}

// CryptBlocks implement BlockMode.CryptBlocks
func (x *ECBDecryptor) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// Signature sha1签名
func Signature(params ...string) string {
	sort.Strings(params)
	h := sha1.New()
	for _, s := range params {
		_, _ = io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// CalculateSign 计算签名
func CalculateSign(content, signType, key string) (string, error) {
	var h hash.Hash
	if signType == SignTypeHMACSHA256 {
		h = hmac.New(sha256.New, []byte(key))
	} else {
		h = md5.New()
	}

	if _, err := h.Write([]byte(content)); err != nil {
		return ``, err
	}
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil))), nil
}

// ParamSign 计算所传参数的签名
func ParamSign(p map[string]string, key string) (string, error) {
	bizKey := "&key=" + key
	str := OrderParam(p, bizKey)

	var signType string
	switch p["sign_type"] {
	case SignTypeMD5, SignTypeHMACSHA256:
		signType = p["sign_type"]
	case ``:
		signType = SignTypeMD5
	default:
		return ``, errors.New(`invalid sign_type`)
	}

	return CalculateSign(str, signType, key)
}

// OrderParam order params
func OrderParam(p map[string]string, bizKey string) (returnStr string) {
	keys := make([]string, 0, len(p))
	for k := range p {
		if k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buf bytes.Buffer
	for _, k := range keys {
		if p[k] == "" {
			continue
		}
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(p[k])
	}
	buf.WriteString(bizKey)
	returnStr = buf.String()
	return
}

// AesECBDecrypt will decrypt data with PKCS5Padding
func AesECBDecrypt(ciphertext []byte, aesKey []byte) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	// ECB mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	NewECBDecryptor(block).CryptBlocks(ciphertext, ciphertext)
	return PKCS5UnPadding(ciphertext), nil
}

// PKCS5UnPadding -
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
