package proto

import "github.com/cafebazaar/booker-resources/common"

func ReplyPropertiesTemplate() *ReplyProperties {
	return &ReplyProperties{
		ServerVersion: common.Version,
	}
}
