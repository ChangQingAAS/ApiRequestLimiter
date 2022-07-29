package limiter

import "fmt"

func getLimiterKey(user string) string {
	return fmt.Sprintf("%s_%s", "limiter", user)
}

func getLimiterLockKey(user string) string {
	return fmt.Sprintf("%s_%s", "limiterLock", user)
}
