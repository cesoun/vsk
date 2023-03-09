package errors

// LockfileError defines an error related to Lockfile processing
type LockfileError int

const (
	ErrFailedToFindUserCache LockfileError = iota
	ErrNoLockfileFound
	ErrLockfileRead
	ErrLockfileBadLength
	ErrLockfilePidAtoi
	ErrLockfilePortAtoi
)

func (l LockfileError) Error() string {
	switch l {
	case ErrFailedToFindUserCache:
		return "failed to resolve the system's cache directory"
	case ErrNoLockfileFound:
		return "no lockfile was detected, is valorant running"
	case ErrLockfileRead:
		return "failed to read the lockfile bytes"
	case ErrLockfileBadLength:
		return "lockfile length was not equal to 5"
	case ErrLockfilePidAtoi:
		return "failed to parse pid to int (atoi)"
	case ErrLockfilePortAtoi:
		return "failed to parse port to int (atoi)"
	default:
		return "unknown lockfile error"
	}
}
