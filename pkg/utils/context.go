package utils

import (
	"context"

	"github.com/gunktp20/digital-hubx-be/pkg/constant"
)

func GetContextAuth(ctx context.Context) (token, user, email string) {
	var ok bool
	if token, ok = ctx.Value(constant.CtxToken).(string); !ok {
		token = ""
	}
	if user, ok = ctx.Value(constant.CtxName).(string); !ok {
		user = ""
	}
	if email, ok = ctx.Value(constant.CtxEmail).(string); !ok {
		email = ""
	}
	return
}
