package caching

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var RedisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func GetOrSetCache(cacheKey string, getDataFunc func() (map[string]interface{},error) ) (interface{},error){
	// GET FROM CACHE FIRST
	res,err := RedisClient.HGetAll(ctx,cacheKey).Result();
	if err == nil && len(res) > 0{
		fmt.Println("Get from cache")
		data := res["data"]
		return data,nil
	}
	// GET FROM DB IF CACHE EMPTY
	dataGetFromDb,err := getDataFunc();
	if err != nil{
		return nil,err
	}

	jsonData,err := json.Marshal(dataGetFromDb["data"])
	if err != nil {
		return nil,err
	}
	// Set to cache
	err = RedisClient.HSet(ctx,cacheKey, "data", jsonData).Err()
	if err != nil {
		return nil,err
	}
	err = RedisClient.Expire(ctx,cacheKey,time.Minute*15).Err()
	if err != nil {
		return nil,err
	}
	return dataGetFromDb["data"],nil
}

func InvalidCacheKey(cacheKey string) error{
	return RedisClient.Del(ctx,cacheKey).Err()
}

func SetCache(cacheKey string, data interface{}) error{
	jsonData,err := json.Marshal(data)
	if err != nil {
		return err
	}
	err =RedisClient.HSet(ctx,cacheKey,"data",jsonData).Err()
	return err
}




