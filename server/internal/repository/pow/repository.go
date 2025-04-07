package pow

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/opravdin/word-of-wisdom/internal/logger"
	"github.com/redis/go-redis/v9"
)

// Default timeout values
const (
	// RedisOperationTimeout is the maximum duration for Redis operations
	RedisOperationTimeout = 5 * time.Second
)

const (
	// slidingWindowDuration is the duration of the sliding window for rate limiting
	slidingWindowDuration = time.Minute
	// unsolvedCountTTL is the TTL for unsolved challenge counts
	unsolvedCountTTL = time.Minute
	// hashLength is the number of characters to use from the hash for readability
	hashLength = 16
)

// hashIPAddress hashes an IP address to create a safe Redis key
// This prevents issues with special characters in IPv6 addresses
func hashIPAddress(ip string) string {
	hash := sha256.Sum256([]byte(ip))
	return hex.EncodeToString(hash[:])[:hashLength]
}

// DefaultRepository implements the Repository interface
type DefaultRepository struct {
	db             *redis.Client
	bucketCapacity int
	logger         logger.Logger
}

// NewRepository creates a new DefaultRepository
func NewRepository(db *redis.Client, bucketCapacity int, log logger.Logger) *DefaultRepository {
	return &DefaultRepository{
		db:             db,
		bucketCapacity: bucketCapacity,
		logger:         log,
	}
}

func (r *DefaultRepository) CreateTask(ctx context.Context, task Task, ttl time.Duration) error {
	r.logger.Debug("Creating task", "id", task.ID, "ttl", ttl)

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, RedisOperationTimeout)
	defer cancel()

	success, err := r.db.SetNX(timeoutCtx, fmt.Sprintf("task:%s", task.ID), task.Result, ttl).Result()
	if err != nil {
		r.logger.Error("Failed to create task", "id", task.ID, "error", err)
		return fmt.Errorf("failed to create task: %w", err)
	}
	if !success {
		r.logger.Debug("Task already exists", "id", task.ID)
		return ErrTaskAlreadyExists
	}
	r.logger.Debug("Task created successfully", "id", task.ID)
	return nil
}

func (r *DefaultRepository) GetTask(ctx context.Context, taskID string) (*Task, error) {
	r.logger.Debug("Getting task", "id", taskID)

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, RedisOperationTimeout)
	defer cancel()

	result, err := r.db.Get(timeoutCtx, fmt.Sprintf("task:%s", taskID)).Result()
	if err == redis.Nil {
		r.logger.Debug("Task not found", "id", taskID)
		return nil, ErrTaskNotFound
	}
	if err != nil {
		r.logger.Error("Redis operation failed", "id", taskID, "error", err)
		return nil, ErrRedisOperation
	}

	r.logger.Debug("Task found", "id", taskID)
	return &Task{
		ID:     taskID,
		Result: result,
	}, nil
}

// GetAndIncrementRequestCount returns and increments the number of requests made within the sliding window
// for a given IP address
func (r *DefaultRepository) GetAndIncrementRequestCount(ctx context.Context, ip string) (int64, error) {
	// Hash the IP address to avoid issues with special characters
	hashedIP := hashIPAddress(ip)
	key := fmt.Sprintf("difficulty:%s", hashedIP)
	now := time.Now().UnixMilli()
	windowStart := now - int64(slidingWindowDuration.Milliseconds())

	r.logger.Debug("Getting and incrementing request count", "ip", ip, "hashedIP", hashedIP)

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, RedisOperationTimeout)
	defer cancel()

	pipe := r.db.Pipeline()

	pipe.ZAdd(timeoutCtx, key, redis.Z{Score: float64(now), Member: now})

	pipe.ZRemRangeByScore(timeoutCtx, key, "0", fmt.Sprintf("%d", windowStart))

	pipe.ZCount(timeoutCtx, key, fmt.Sprintf("%d", windowStart), fmt.Sprintf("%d", now))

	pipe.Expire(timeoutCtx, key, slidingWindowDuration*2)

	results, err := pipe.Exec(timeoutCtx)
	if err != nil {
		r.logger.Error("Redis operation failed", "ip", ip, "error", err)
		return 0, ErrRedisOperation
	}

	// The 3rd command in the pipeline is count, that's what we return
	count := results[2].(*redis.IntCmd).Val()

	r.logger.Debug("Request count", "ip", ip, "count", count)
	return count, nil
}

func (r *DefaultRepository) IncrementUnsolvedCount(ctx context.Context, ip string) (int64, error) {
	// Hash the IP address to avoid issues with special characters
	hashedIP := hashIPAddress(ip)
	key := fmt.Sprintf("unsolved:%s", hashedIP)

	r.logger.Debug("Incrementing unsolved count", "ip", ip, "hashedIP", hashedIP)

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, RedisOperationTimeout)
	defer cancel()

	pipe := r.db.Pipeline()

	// Increment
	incrCmd := pipe.Incr(timeoutCtx, key)
	// Set/Reset TTL
	pipe.Expire(timeoutCtx, key, unsolvedCountTTL)

	if _, err := pipe.Exec(timeoutCtx); err != nil && err != redis.Nil {
		r.logger.Error("Redis operation failed", "ip", ip, "error", err)
		return 0, ErrRedisOperation
	}

	// Return the new count after increment
	count, err := incrCmd.Result()
	if err != nil && err != redis.Nil {
		r.logger.Error("Redis operation failed", "ip", ip, "error", err)
		return 0, ErrRedisOperation
	}

	r.logger.Debug("Unsolved count incremented", "ip", ip, "count", count)
	return count, nil
}

func (r *DefaultRepository) DecrementUnsolvedCount(ctx context.Context, ip string) error {
	return r.DecrementUnsolvedCountBy(ctx, ip, 1)
}

func (r *DefaultRepository) DecrementUnsolvedCountBy(ctx context.Context, ip string, count int) error {
	// Hash the IP address to avoid issues with special characters
	hashedIP := hashIPAddress(ip)
	key := fmt.Sprintf("unsolved:%s", hashedIP)

	r.logger.Debug("Decrementing unsolved count", "ip", ip, "hashedIP", hashedIP, "count", count)

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, RedisOperationTimeout)
	defer cancel()

	pipe := r.db.Pipeline()

	// Decrement by count but not below 0
	pipe.Eval(timeoutCtx, `
		local count = redis.call('get', KEYS[1])
		if count and tonumber(count) > 0 then
			local newCount = math.max(0, tonumber(count) - tonumber(ARGV[1]))
			return redis.call('set', KEYS[1], newCount)
		end
		return 0
	`, []string{key}, count)
	// Reset TTL
	pipe.Expire(timeoutCtx, key, unsolvedCountTTL)

	if _, err := pipe.Exec(timeoutCtx); err != nil && err != redis.Nil {
		r.logger.Error("Redis operation failed", "ip", ip, "error", err)
		return ErrRedisOperation
	}

	r.logger.Debug("Unsolved count decremented", "ip", ip, "by", count)
	return nil
}

func (r *DefaultRepository) GetUnsolvedCount(ctx context.Context, ip string) (int64, error) {
	// Hash the IP address to avoid issues with special characters
	hashedIP := hashIPAddress(ip)
	key := fmt.Sprintf("unsolved:%s", hashedIP)

	r.logger.Debug("Getting unsolved count", "ip", ip, "hashedIP", hashedIP)

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, RedisOperationTimeout)
	defer cancel()

	count, err := r.db.Get(timeoutCtx, key).Int64()
	if err == redis.Nil {
		r.logger.Debug("No unsolved count found", "ip", ip)
		return 0, nil
	}
	if err != nil {
		r.logger.Error("Redis operation failed", "ip", ip, "error", err)
		return 0, ErrRedisOperation
	}

	r.logger.Debug("Got unsolved count", "ip", ip, "count", count)
	return count, nil
}

func (r *DefaultRepository) DeleteTask(ctx context.Context, taskID string) error {
	key := fmt.Sprintf("task:%s", taskID)

	r.logger.Debug("Deleting task", "id", taskID)

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, RedisOperationTimeout)
	defer cancel()

	result, err := r.db.Del(timeoutCtx, key).Result()
	if err != nil {
		r.logger.Error("Failed to delete task", "id", taskID, "error", err)
		return ErrRedisOperation
	}

	if result == 0 {
		r.logger.Debug("Task not found for deletion", "id", taskID)
		return ErrTaskNotFound
	}

	r.logger.Debug("Task deleted successfully", "id", taskID)
	return nil
}
