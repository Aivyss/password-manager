package pwmContext

import "sync"

var gCtxOnce sync.Once
var gCtx GlobalContext

type GlobalContext struct {
	BuildSecretKey string
}

func GetGlobalContext() GlobalContext {
	return gCtx
}

func NewGlobalContext(buildSecretKey string) *GlobalContext {
	gCtxOnce.Do(func() {
		gCtx = GlobalContext{
			BuildSecretKey: buildSecretKey,
		}
	})

	return &gCtx
}
