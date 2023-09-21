package utils

import "encoding/base64"

// ToBase64 将字节数组转换为Base64字符串
func ToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// FromBase64 将Base64字符串解码为字节数组
func FromBase64(base64Str string) ([]byte, error) {
	decodedData, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}
	return decodedData, nil
}
