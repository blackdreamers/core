package utils

import (
	"context"
	"strings"

	"go-micro.dev/v4/metadata"

	"github.com/blackdreamers/core/consts"
)

func GetClientName(ctx context.Context) (string, bool) {
	md, ok := metadata.FromContext(ctx)
	if ok {
		return md.Get(consts.SrvNameKey)
	}
	return "", false
}

func GetClientSrvName(ctx context.Context) (string, bool) {
	val, ok := GetClientName(ctx)
	if ok {
		idx := strings.Index(val, consts.Delimiter)
		if idx != len(val) {
			return val[idx+1:], true
		}
	}
	return "", false
}
