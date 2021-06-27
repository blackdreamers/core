package utils

import (
	"context"
	"strings"

	"github.com/blackdreamers/core/constant"
	"github.com/blackdreamers/go-micro/v3/metadata"
)

func GetClientSrvName(ctx context.Context) (string, bool) {
	md, ok := metadata.FromContext(ctx)
	if ok {
		return md.Get("Micro-From-Service")
	}
	return "", false
}

func GetClientName(ctx context.Context) (string, bool) {
	val, ok := GetClientSrvName(ctx)
	if ok {
		idx := strings.Index(val, constant.Delimiter)
		if idx != len(val) {
			return val[idx+1:], true
		}
	}
	return "", false
}
