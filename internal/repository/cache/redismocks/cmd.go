package redismocks

//go:generate mockgen -package=redismocks -destination=./cmd.mock.go github.com/redis/go-redis/v9 Cmdable
