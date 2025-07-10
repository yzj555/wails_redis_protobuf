package server

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"sort"
)

var FuncMap = make(map[string]func(params map[string]interface{}) interface{})

func InitFuncMap() {
	registerFunc("HandleConfig", HandleConfig)
	registerFunc("HandleChangeServer", HandleChangeServer)
	registerFunc("HandleAllMsgName", HandleAllMsgName)
	registerFunc("HandleAllKeys", HandleAllKeys)
	registerFunc("HandleDataByKey", HandleDataByKey)
	registerFunc("HandleModifyData", HandleModifyData)
	registerFunc("HandleParseData", HandleParseData)
}

func registerFunc(funcName string, fun func(params map[string]interface{}) interface{}) {
	FuncMap[funcName] = fun
}

func HandleConfig(params map[string]interface{}) interface{} {
	return Config
}

func HandleChangeServer(params map[string]interface{}) interface{} {
	name := params["name"].(string)
	result := make(map[string]interface{})
	if name == "" {
		result["data"] = nil
		result["error"] = "name is empty"
		return result
	}
	if name == Config.CurrentRedis {
		result["data"] = nil
		result["error"] = "current server is " + name
		return result

	}
	ChangeRedisServer(name)
	result["data"] = Config.CurrentRedis
	result["error"] = nil
	return result
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
	result, err := scanAllKeys(ctx, likeKey, 1000)
	sort.Strings(result)
	if err != nil {
		return err.Error()
	}
	fmt.Println("likeKey:::::", likeKey)
	return result
}

func HandleDataByKey(params map[string]interface{}) interface{} {
	key := params["key"].(string)
	//field := params["field"].(string)
	//parseMode := params["parseMode"].(string)
	//msgName := params["msg"].(string)
	//result, err := queryKey(key, field, parseMode, msgName)
	result, err := queryDataByKey(key)
	if err != nil {
		return err.Error()
	}
	return result
}

func HandleModifyData(params map[string]interface{}) interface{} {
	result := make(map[string]interface{})

	key := params["key"].(string)
	parseMode := params["parseMode"].(string)
	msgName := params["msg"].(string)
	data := params["data"].(string)

	var dbData string

	switch parseMode {
	case ParseModeProto:
		fmt.Println("jsonData:::::  ", data)
		dMsg, errMsg := MarshalByProto([]byte(data), msgName)
		if dMsg == nil {
			result["data"] = nil
			result["error"] = errMsg
			return result
		}
		//把消息体序列化成二进制数据
		pbStr, err := dMsg.Marshal()
		if err != nil {
			fmt.Printf("Marshal err, %v\n", err)
			result["data"] = nil
			result["error"] = err.Error()
			return result
		}
		dbData = string(pbStr)
	case ParseModeMsgpack:
		result["data"] = nil
		result["error"] = "msgpack方式不允许修改！"
		return result
		//fields, err := MarshalByMsgPack(data, msgName)
		//if err != nil {
		//	result["data"] = nil
		//	result["error"] = err.Error()
		//	return result
		//}
		//dbData = fields
	case ParseModeSource: //源数据类型当做string处理，不作处理
		dbData = data
	default:
		result["data"] = nil
		result["error"] = "parseMode is empty"
		return result
	}
	if dbData == "" {
		result["data"] = nil
		result["error"] = "data is empty"
		return result
	}
	keyType, err := rdb.Type(ctx, key).Result()
	if err != nil {
		result["data"] = nil
		result["error"] = err.Error()
		return result
	}
	switch keyType {
	case "string":
		err = rdb.Set(ctx, key, dbData, 0).Err()
		if err != nil {
			result["data"] = nil
			result["error"] = err.Error()
			return result
		}
	case "hash":
		field := params["field"].(string)
		err = rdb.HSet(ctx, key, field, dbData).Err()
		if err != nil {
			result["data"] = nil
			result["error"] = err.Error()
			return result
		}
	case "list":
		field := params["field"].(int64)
		err = rdb.LSet(ctx, key, field, dbData).Err()
		if err != nil {
			result["data"] = nil
			result["error"] = err.Error()
			return result
		}
	case "set":
		oldVal := params["oldItem"].(string)
		flag, err := rdb.SIsMember(ctx, key, oldVal).Result()
		if err != nil {
			result["data"] = nil
			result["error"] = err.Error()
			return result
		}
		if !flag {
			result["data"] = nil
			result["error"] = "item is not exist"
			return result
		}
		_, err = rdb.SRem(ctx, key, oldVal).Result()
		if err != nil {
			result["data"] = nil
			result["error"] = err.Error()
			return result
		}
		err = rdb.SAdd(ctx, key, dbData).Err()
		if err != nil {
			result["data"] = nil
			result["error"] = err.Error()
			return result
		}
	case "zset":
		member := params["field"].(string)
		score := params["score"].(float64)
		_, err := rdb.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Result()
		if err != nil {
			result["data"] = nil
			result["error"] = err.Error()
			return result
		}
	default:
		fmt.Println("未知类型或键不存在")
		result["data"] = nil
		result["error"] = "unknown type or key not exist"
		return result
	}
	result["data"] = data
	result["error"] = ""
	return result
}

func queryDataByKey(key string) (interface{}, error) {
	keyType, err := rdb.Type(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var data interface{}
	switch keyType {
	case "string":
		data, _ = rdb.Get(ctx, key).Result()
	case "hash":
		temp, _ := rdb.HGetAll(ctx, key).Result()
		resArr := make([]map[string]interface{}, 0)
		for k, v := range temp {
			resArr = append(resArr, map[string]interface{}{
				"field": k,
				"value": v,
			})
		}
		data = resArr
	case "list":
		temp, _ := rdb.LRange(ctx, key, 0, -1).Result()
		resArr := make([]map[string]interface{}, 0)
		for k, v := range temp {
			resArr = append(resArr, map[string]interface{}{
				"field": k,
				"value": v,
			})
		}
		data = resArr
	case "set":
		temp, _ := rdb.SMembers(ctx, key).Result()
		resArr := make([]map[string]interface{}, 0)
		for k, v := range temp {
			resArr = append(resArr, map[string]interface{}{
				"field": k,
				"value": v,
			})
		}
		data = resArr
	case "zset":
		temp, _ := rdb.ZRevRangeWithScores(ctx, key, 0, -1).Result()
		resArr := make([]map[string]interface{}, 0)
		for _, v := range temp {
			resArr = append(resArr, map[string]interface{}{
				"field": v.Member,
				"value": v.Score,
			})
		}
		data = resArr
	default:
		fmt.Println("未知类型或键不存在")
		return nil, nil
	}
	result := make(map[string]interface{})
	result["dataType"] = keyType
	result["data"] = data
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

func HandleParseData(params map[string]interface{}) interface{} {
	result := make(map[string]interface{})

	key := params["key"].(string)
	keyType, err := rdb.Type(ctx, key).Result()
	if err != nil {
		result["data"] = nil
		result["error"] = err.Error()
		return result
	}

	var val string
	switch keyType {
	case "string":
		val, err = rdb.Get(ctx, key).Result()
		if err != nil {
			result["data"] = nil
			result["error"] = err.Error()
			return result
		}
	case "hash":
		field := params["field"].(string)
		if field == "" {
			result["data"] = nil
			result["error"] = "field is empty"
			return result
		}
		val, err = rdb.HGet(ctx, key, field).Result()
		if err != nil {
			result["data"] = nil
			result["error"] = err.Error()
			return result
		}
	case "list":
		index := params["field"].(int64)
		val, err = rdb.LIndex(ctx, key, index).Result()
		if err != nil {
			result["data"] = nil
			result["error"] = err.Error()
			return result
		}
	case "set":
		oldVal := params["oldItem"].(string)
		flag, err := rdb.SIsMember(ctx, key, oldVal).Result()
		if err != nil {
			result["data"] = nil
			result["error"] = err.Error()
			return result
		}
		if !flag {
			result["data"] = nil
			result["error"] = "item is not exist"
			return result
		}
		val = oldVal
	case "zset":
		oldVal := params["field"].(string)
		score, err := rdb.ZScore(ctx, key, oldVal).Result()
		if err != nil {
			result["data"] = nil
			result["error"] = err.Error()
			return result
		}
		result["data"] = score
		result["error"] = ""
		return result

	default:
		result["data"] = nil
		result["error"] = "未知类型或键不存在"
		return result
	}
	msgName := params["msg"].(string)
	parseMode := params["parseMode"].(string)
	data := parseString(val, msgName, parseMode)
	result["data"] = data
	result["error"] = ""
	return result
}
