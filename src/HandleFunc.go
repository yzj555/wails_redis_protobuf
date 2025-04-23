package server

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

var FuncMap = make(map[string]func(params map[string]interface{}) interface{})

func InitFuncMap() {
	registerFunc("HandleAllMsgName", HandleAllMsgName)
	registerFunc("HandleAllKeys", HandleAllKeys)
	registerFunc("HandleDataByKey", HandleDataByKey)
}

func registerFunc(funcName string, fun func(params map[string]interface{}) interface{}) {
	FuncMap[funcName] = fun
}

func HandleAllMsgName(params map[string]interface{}) interface{} {
	return GetAllMsgName()
}

func HandleAllKeys(params map[string]interface{}) interface{} {
	likeKey := params["likeKey"].(string)
	if likeKey == "" {
		likeKey = "*"
	} else {
		likeKey = fmt.Sprintf("*%s*", likeKey)
	}
	result, err := rdb.Keys(ctx, likeKey).Result()
	if err != nil {
		return err.Error()
	}
	return result
}

func HandleDataByKey(params map[string]interface{}) interface{} {
	key := params["key"].(string)
	field := params["field"].(string)
	parseMode := params["parseMode"].(string)
	msgName := params["msg"].(string)
	result, err := queryKey(key, field, parseMode, msgName)
	if err != nil {
		return err.Error()
	}
	return result
}

func queryKey(key, field, parseMode, msgName string) (interface{}, error) {
	keyType, err := rdb.Type(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var result interface{}
	switch keyType {
	case "string":
		val, _ := rdb.Get(ctx, key).Result()
		result = parseString(val, msgName, parseMode)
	case "hash":
		if field != "" {
			val, _ := rdb.HGet(ctx, key, field).Result()
			result = parseString(val, msgName, parseMode)
		} else {
			fields, _ := rdb.HGetAll(ctx, key).Result()
			result = parseHash(fields, msgName, parseMode)
		}
	case "list":
		items, _ := rdb.LRange(ctx, key, 0, -1).Result()
		result = parseListOrSet(items, msgName, parseMode)
	case "set":
		items, _ := rdb.SMembers(ctx, key).Result()
		result = parseListOrSet(items, msgName, parseMode)
		// 其他类型：set、zset、stream 等
	case "zset":
		items, _ := rdb.ZRevRangeWithScores(ctx, key, 0, -1).Result()
		result = parseZset(items, msgName, parseMode)
	default:
		fmt.Println("未知类型或键不存在")
	}
	return result, nil
}

const (
	ParseModeProto   = "1"
	ParseModeSource  = "2"
	ParseModeMsgpack = "3"
)

func parseString(data string, msgName string, parseMode string) interface{} {
	if msgName == "" && parseMode != ParseModeMsgpack {
		return data
	}
	switch parseMode {
	case ParseModeProto:
		return UnmarshalByProto([]byte(data), msgName)
	case ParseModeMsgpack:
		return UnmarshalByMsgpack([]byte(data))
	default:
		return data
	}
}

func parseHash(fields map[string]string, msgName string, parseMode string) interface{} {
	tempMap := make(map[string]interface{})
	if msgName == "" && parseMode != ParseModeMsgpack {
		return fields
	}
	for k, v := range fields {
		tempMap[k] = parseString(v, msgName, parseMode)
		//tempMap[k] = server.UnmarshalByProto([]byte(v), msgName)
	}
	return tempMap
}

func parseListOrSet(items []string, msgName string, parseMode string) interface{} {
	tempMap := make([]interface{}, 0, len(items))
	if msgName == "" && parseMode != ParseModeMsgpack {
		return items
	}
	for _, v := range items {
		tempMap = append(tempMap, parseString(v, msgName, parseMode))
	}
	return tempMap
}

func parseZset(items []redis.Z, msgName string, parseMode string) interface{} {
	tempMap := make([]interface{}, 0, len(items))
	if msgName == "" && parseMode != ParseModeMsgpack {
		return items
	}
	for _, v := range items {
		temp := make(map[string]interface{})
		temp["score"] = v.Score
		temp["member"] = parseString(v.Member.(string), msgName, parseMode)
		tempMap = append(tempMap, temp)
	}
	return tempMap
}
