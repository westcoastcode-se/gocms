package security

import "context"

// Fet a user from the supplied context
func GetUserFromContext(ctx context.Context) *User {
	if ret, ok := ctx.Value(SessionKey).(*User); ok {
		return ret
	}
	return NotLoggedInUser
}

func PutUserInContext(user *User, ctx context.Context) context.Context {
	return context.WithValue(ctx, SessionKey, user)
}
