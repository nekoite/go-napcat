package utils

import "github.com/tidwall/gjson"

func IsRawMessageApiResp(resp []byte) bool {
	return gjson.GetBytes(resp, "echo").Exists()
}
