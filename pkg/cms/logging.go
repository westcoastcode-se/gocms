package cms

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/westcoastcode-se/gocms/pkg/security"
	"net/http"
)

const ContextKeyRequestID = "request_id"
const UserKey = "user"

func GetRequestID(ctx context.Context) string {
	if ret, ok := ctx.Value(ContextKeyRequestID).(string); ok {
		return ret
	}
	return uuid.New().String()
}

// Create a logger that contains useful context information
func Log(ctx context.Context) logrus.FieldLogger {
	id := GetRequestID(ctx)
	user := security.GetUserFromContext(ctx).Name
	logger := logrus.StandardLogger()
	return logger.WithField(ContextKeyRequestID, id).WithField(UserKey, user)
}

func LogFromRequest(r *http.Request) logrus.FieldLogger {
	return Log(r.Context())
}
