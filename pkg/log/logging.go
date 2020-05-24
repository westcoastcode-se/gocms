package log

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
)

const RequestIdKey = "request_id"
const UserKey = "user"

func SetUserName(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func SetRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIdKey, id)
}

func GetUserName(ctx context.Context) string {
	if ret, ok := ctx.Value(UserKey).(string); ok {
		return ret
	}
	return ""
}

func GetRequestID(ctx context.Context) string {
	if ret, ok := ctx.Value(RequestIdKey).(string); ok {
		return ret
	}
	return ""
}

// Create a logger that contains useful context information
func FromContext(ctx context.Context) logrus.FieldLogger {
	id := GetRequestID(ctx)
	user := GetUserName(ctx)
	logger := logrus.StandardLogger()
	return logger.WithField(RequestIdKey, id).WithField(UserKey, user)
}

func LogFromRequest(r *http.Request) logrus.FieldLogger {
	return FromContext(r.Context())
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Infof(format, args)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Warnf(format, args)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Errorf(format, args)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		FromContext(ctx).Debugf(format, args)
	}
}
