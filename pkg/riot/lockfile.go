package riot

import (
	"github.com/cesoun/vsk/pkg/riot/errors"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// Lockfile defined by the Riot lockfile content: name:pid:port:password:protocol
type Lockfile struct {
	Name     string
	PID      int
	Port     int
	Password string
	Protocol string
}

// NewLockfile attempts to read the Riot lockfile and parse the contents
func NewLockfile() (*Lockfile, error) {
	lockfilePath, err := GetLockfilePath()
	if err != nil {
		return nil, err
	}

	if !DoesLockfileExist(lockfilePath) {
		return nil, errors.ErrNoLockfileFound
	}

	fbytes, err := os.ReadFile(lockfilePath)
	if err != nil {
		return nil, errors.ErrLockfileRead
	}

	lockfileParts := strings.Split(string(fbytes), ":")
	if len(lockfileParts) != 5 {
		return nil, errors.ErrLockfileBadLength
	}

	pid, err := strconv.Atoi(lockfileParts[1])
	if err != nil {
		return nil, errors.ErrLockfilePidAtoi
	}

	port, err := strconv.Atoi(lockfileParts[2])
	if err != nil {
		return nil, errors.ErrLockfilePortAtoi
	}

	return &Lockfile{
		Name:     lockfileParts[0],
		PID:      pid,
		Port:     port,
		Password: lockfileParts[3],
		Protocol: lockfileParts[4],
	}, nil
}

// GetLockfilePath returns the lockfile path as a string
func GetLockfilePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", errors.ErrFailedToFindUserCache
	}

	lockfilePath := filepath.Join(cacheDir, "Riot Games", "Riot Client", "Config", "lockfile")

	return lockfilePath, nil
}

// GetConfigPath returns the config path as a string
func GetConfigPath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", errors.ErrFailedToFindUserCache
	}

	lockfilePath := path.Join(cacheDir, "Riot Games", "Riot Client", "Config")

	return lockfilePath, nil
}

// DoesLockfileExist checks if the lockfile currently exists
func DoesLockfileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}

	return true
}
