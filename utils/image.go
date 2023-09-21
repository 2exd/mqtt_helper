package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
)

var ImgChan = make(chan *image.RGBA, 10)

// DecodeImageFromBase64 从 Base64 字符串解码为 *image.RGBA
func DecodeImageFromBase64(base64String string) (*image.RGBA, error) {
	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	if rgbaImg, ok := img.(*image.RGBA); ok {
		return rgbaImg, nil
	}

	// 如果图像不是 *image.RGBA 类型，可以尝试进行类型转换
	return nil, fmt.Errorf("Decoded image is not of type *image.RGBA")
}
