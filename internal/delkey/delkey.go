package delkey

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"os"
	"strings"
	"time"
)

func logger(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
}

func writeLog(logger *os.File, format string, a ...interface{}) {
	t := time.Now().Format("2006/01/02 15:04:05.000")
	_, _ = logger.WriteString(fmt.Sprintf(t+"|"+format+"\n", a...))
}

func redisClient(addr, auth string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: addr, Password: auth, DB: db})
}

func Processor(host, port, auth, match string, count int64, ttl int8, logDir string) {
	if host == "" {
		log.Fatal("请输入Redis服务地址")
	}
	if port == "" {
		log.Fatal("请输入Redis服务端口")
	}

	rds := redisClient(host+":"+port, auth, 0)
	defer func() {
		if err := rds.Close(); err != nil {
			log.Printf("关闭Redis实例出错：%v", err)
		}
	}()

	if rds == nil {
		log.Fatal("创建Redis实例失败")
	}
	if err := rds.Ping().Err(); err != nil {
		log.Fatalf("Ping Redis服务出错：%v", err)
	}
	match = strings.TrimSpace(match)
	if match == "" {
		log.Fatal("请输入要删除keys的指定前缀")
	}
	if count <= 0 {
		log.Fatal("请输入每次扫描的返回数量")
	}
	if ttl != -1 && ttl != 1 {
		log.Fatal("请输入要清除缓存的类型(-1永久缓存|1所有)")
	}

	var cursor, succeed uint64
	start := time.Now()

	filename := logDir + "/key(" + match + ")_" + start.Format("2006-01-02T15-04-05") + ".log"
	logger, err := logger(filename)
	if err != nil {
		log.Fatalf("open log file \"%s\" error: %v", filename, err)
	}
	defer logger.Close()

	for {
		var keys []string
		var err error

		if keys, cursor, err = rds.Scan(cursor, match+`*`, count).Result(); err != nil {
			log.Fatalf("redis Scan err: %v", err)
		}

		if len(keys) > 0 {
			for _, key := range keys {
				if key == "" {
					continue
				}
				var keyType string
				if keyType, err = rds.Type(key).Result(); err != nil {
					log.Printf("key: %s, redis.Type error: %v\n", key, err)
					continue
				}
				var keyTtl time.Duration
				if keyTtl, err = rds.TTL(key).Result(); err != nil {
					log.Printf("key: %s, redis.TTL error: %v\n", key, err)
					continue
				}
				if keyTtl.Seconds() == -2 {
					continue
				}
				if ttl == -1 && keyTtl.Seconds() != -1 {
					continue
				}
				var value interface{}
				switch keyType {
				case "string":
					value, err = rds.Get(key).Result()
				case "set":
					value, err = rds.SMembers(key).Result()
				case "zset":
					value, err = rds.ZRevRangeWithScores(key, 0, -1).Result()
				case "list":
					value, err = rds.LRange(key, 0, -1).Result()
				case "hash":
					value, err = rds.HGetAll(key).Result()
				default:
					log.Printf("key: %s, type: %s undefined\n", key, keyType)
					continue
				}
				if err != nil {
					log.Printf("key: %s, redis.Value error: %v\n", key, err)
					continue
				}

				if _, err = rds.Del(key).Result(); err != nil {
					log.Printf("key: %s, redis.Del error: %v\n", key, err)
					continue
				} else {
					succeed++
					bts, _ := json.Marshal(value)
					writeLog(logger, "key=%s|type=%s|value=%s", key, keyType, string(bts))
				}
			}
		}

		if cursor == 0 {
			break
		}
	}
	log.Printf("deleted %d, takes %v\n", succeed, time.Now().Sub(start))
}
