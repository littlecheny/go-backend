package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/littlecheny/go-backend/domain"
)

type redisService struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) domain.RedisService {
	return &redisService{
		client: client,
	}
}

func (r *redisService) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	
	// 将值序列化为JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	err = r.client.Set(ctx, key, jsonValue, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %v", key, err)
	}

	return nil
}

func (r *redisService) Get(key string) (string, error) {
	ctx := context.Background()
	
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key %s not found", key)
		}
		return "", fmt.Errorf("failed to get key %s: %v", key, err)
	}

	return val, nil
}

func (r *redisService) Del(key string) error {
	ctx := context.Background()
	
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %v", key, err)
	}

	return nil
}

// 辅助方法：获取并反序列化JSON
func (r *redisService) GetJSON(key string, dest interface{}) error {
	val, err := r.Get(key)
	if err != nil {
		return err
	}

	// 反序列化JSON
	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal value for key %s: %v", key, err)
	}

	return nil
}

func (r *redisService) Exists(key string) (bool, error) {
	ctx := context.Background()
	
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %v", key, err)
	}

	return exists > 0, nil
}

func (r *redisService) SetExpiration(key string, expiration time.Duration) error {
	ctx := context.Background()
	
	err := r.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration for key %s: %v", key, err)
	}

	return nil
}

func (r *redisService) GetTTL(key string) (time.Duration, error) {
	ctx := context.Background()
	
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL for key %s: %v", key, err)
	}

	return ttl, nil
}

func (r *redisService) SetHash(key, field string, value interface{}) error {
	ctx := context.Background()
	
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	err = r.client.HSet(ctx, key, field, jsonValue).Err()
	if err != nil {
		return fmt.Errorf("failed to set hash field %s:%s: %v", key, field, err)
	}

	return nil
}

func (r *redisService) GetHash(key, field string, dest interface{}) error {
	ctx := context.Background()
	
	val, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("hash field %s:%s not found", key, field)
		}
		return fmt.Errorf("failed to get hash field %s:%s: %v", key, field, err)
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal value for hash field %s:%s: %v", key, field, err)
	}

	return nil
}

func (r *redisService) GetAllHash(key string) (map[string]string, error) {
	ctx := context.Background()
	
	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get all hash fields for key %s: %v", key, err)
	}

	return result, nil
}

func (r *redisService) DeleteHash(key, field string) error {
	ctx := context.Background()
	
	err := r.client.HDel(ctx, key, field).Err()
	if err != nil {
		return fmt.Errorf("failed to delete hash field %s:%s: %v", key, field, err)
	}

	return nil
}

func (r *redisService) IncrementCounter(key string) (int64, error) {
	ctx := context.Background()
	
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment counter %s: %v", key, err)
	}

	return val, nil
}

func (r *redisService) DecrementCounter(key string) (int64, error) {
	ctx := context.Background()
	
	val, err := r.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement counter %s: %v", key, err)
	}

	return val, nil
}

func (r *redisService) AddToList(key string, values ...interface{}) error {
	ctx := context.Background()
	
	// 序列化所有值
	serializedValues := make([]interface{}, len(values))
	for i, value := range values {
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value at index %d: %v", i, err)
		}
		serializedValues[i] = jsonValue
	}

	err := r.client.RPush(ctx, key, serializedValues...).Err()
	if err != nil {
		return fmt.Errorf("failed to add to list %s: %v", key, err)
	}

	return nil
}

func (r *redisService) GetListRange(key string, start, stop int64) ([]string, error) {
	ctx := context.Background()
	
	result, err := r.client.LRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get list range for key %s: %v", key, err)
	}

	return result, nil
}

func (r *redisService) GetListLength(key string) (int64, error) {
	ctx := context.Background()
	
	length, err := r.client.LLen(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get list length for key %s: %v", key, err)
	}

	return length, nil
}

func (r *redisService) PopFromList(key string) (string, error) {
	ctx := context.Background()
	
	val, err := r.client.LPop(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("list %s is empty", key)
		}
		return "", fmt.Errorf("failed to pop from list %s: %v", key, err)
	}

	return val, nil
}

func (r *redisService) AddToSet(key string, members ...interface{}) error {
	ctx := context.Background()
	
	// 序列化所有成员
	serializedMembers := make([]interface{}, len(members))
	for i, member := range members {
		jsonValue, err := json.Marshal(member)
		if err != nil {
			return fmt.Errorf("failed to marshal member at index %d: %v", i, err)
		}
		serializedMembers[i] = jsonValue
	}

	err := r.client.SAdd(ctx, key, serializedMembers...).Err()
	if err != nil {
		return fmt.Errorf("failed to add to set %s: %v", key, err)
	}

	return nil
}

func (r *redisService) GetSetMembers(key string) ([]string, error) {
	ctx := context.Background()
	
	members, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get set members for key %s: %v", key, err)
	}

	return members, nil
}

func (r *redisService) IsSetMember(key string, member interface{}) (bool, error) {
	ctx := context.Background()
	
	jsonValue, err := json.Marshal(member)
	if err != nil {
		return false, fmt.Errorf("failed to marshal member: %v", err)
	}

	isMember, err := r.client.SIsMember(ctx, key, jsonValue).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check set membership for key %s: %v", key, err)
	}

	return isMember, nil
}

func (r *redisService) RemoveFromSet(key string, members ...interface{}) error {
	ctx := context.Background()
	
	// 序列化所有成员
	serializedMembers := make([]interface{}, len(members))
	for i, member := range members {
		jsonValue, err := json.Marshal(member)
		if err != nil {
			return fmt.Errorf("failed to marshal member at index %d: %v", i, err)
		}
		serializedMembers[i] = jsonValue
	}

	err := r.client.SRem(ctx, key, serializedMembers...).Err()
	if err != nil {
		return fmt.Errorf("failed to remove from set %s: %v", key, err)
	}

	return nil
}

func (r *redisService) GetKeys(pattern string) ([]string, error) {
	ctx := context.Background()
	
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys with pattern %s: %v", pattern, err)
	}

	return keys, nil
}

func (r *redisService) FlushDB() error {
	ctx := context.Background()
	
	err := r.client.FlushDB(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to flush database: %v", err)
	}

	return nil
}

func (r *redisService) Ping() error {
	ctx := context.Background()
	
	_, err := r.client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to ping Redis: %v", err)
	}

	return nil
}

// 实现接口要求的方法

// SetBalance 缓存余额信息
func (r *redisService) SetBalance(address string, balance string, expiration time.Duration) error {
	key := fmt.Sprintf("balance:%s", address)
	return r.Set(key, balance, expiration)
}

// GetBalance 获取缓存的余额信息
func (r *redisService) GetBalance(address string) (string, error) {
	key := fmt.Sprintf("balance:%s", address)
	return r.Get(key)
}

// SetGasPrice 缓存Gas价格
func (r *redisService) SetGasPrice(network string, gasPrice string, expiration time.Duration) error {
	key := fmt.Sprintf("gas_price:%s", network)
	return r.Set(key, gasPrice, expiration)
}

// GetGasPrice 获取缓存的Gas价格
func (r *redisService) GetGasPrice(network string) (string, error) {
	key := fmt.Sprintf("gas_price:%s", network)
	return r.Get(key)
}

// SetTransaction 缓存交易信息
func (r *redisService) SetTransaction(hash string, tx *domain.TransactionResponse, expiration time.Duration) error {
	key := fmt.Sprintf("transaction:%s", hash)
	return r.Set(key, tx, expiration)
}

// GetTransaction 获取缓存的交易信息
func (r *redisService) GetTransaction(hash string) (*domain.TransactionResponse, error) {
	key := fmt.Sprintf("transaction:%s", hash)
	var tx domain.TransactionResponse
	err := r.GetJSON(key, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// SetBlock 缓存区块信息
func (r *redisService) SetBlock(number uint64, block *domain.BlockInfo, expiration time.Duration) error {
	key := fmt.Sprintf("block:%d", number)
	return r.Set(key, block, expiration)
}

// GetBlock 获取缓存的区块信息
func (r *redisService) GetBlock(number uint64) (*domain.BlockInfo, error) {
	key := fmt.Sprintf("block:%d", number)
	var block domain.BlockInfo
	err := r.GetJSON(key, &block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}