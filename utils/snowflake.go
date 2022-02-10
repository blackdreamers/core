package utils

import "github.com/sony/sonyflake"

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func NextID() (uint64, error) {
	return sf.NextID()
}
