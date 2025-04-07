package pow

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/opravdin/word-of-wisdom/internal/configuration/env"
	loggermocks "github.com/opravdin/word-of-wisdom/internal/logger/mocks"
	"github.com/opravdin/word-of-wisdom/internal/pow/mocks"
	"github.com/opravdin/word-of-wisdom/internal/protocol"
	randommocks "github.com/opravdin/word-of-wisdom/internal/random/mocks"
	powrepo "github.com/opravdin/word-of-wisdom/internal/repository/pow"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// Test constants
const (
	// Common test values
	testClientIP    = "127.0.0.1"
	testSeed        = "test-seed"
	testNonce       = "test-nonce"
	testValidUUID   = "00000000-0000-0000-0000-000000000000"
	testInvalidUUID = "invalid-uuid"

	// Error messages
	errDB       = "db error"
	errNotFound = "not found"
	errRandom   = "random error"

	// Challenge config values
	testScryptN         = 16384
	testScryptR         = 8
	testScryptP         = 1
	testKeyLen          = 32
	testMaxUnsolved     = 5
	testChallengeTTL    = 5 * time.Minute
	testRequestsPerDiff = 10
	testMaxDiffLevel    = 5

	// Test request counts
	testRequestCount    = int64(5)
	testUnsolvedCount   = int64(1)
	testTooManyUnsolved = int64(6)

	// Test difficulty level
	testDifficultyLevel = 1
)

func TestService_CreateChallenge(t *testing.T) {
	tests := []struct {
		name       string
		clientIP   string
		expected   *protocol.Challenge
		mockFunc   func(mockRepo *mocks.MockRepository, mockUtils *mocks.MockPoWUtilsInterface, mockLog *loggermocks.MockLogger, mockRandom *randommocks.MockProvider)
		shouldFail bool
	}{
		{
			name:     "should_create_challenge_successfully",
			clientIP: testClientIP,
			expected: &protocol.Challenge{
				ID:              "test-id", // ID will be generated with UUID
				Seed:            testSeed,
				DifficultyLevel: testDifficultyLevel,
				ScryptN:         testScryptN,
				ScryptR:         testScryptR,
				ScryptP:         testScryptP,
				KeyLen:          testKeyLen,
			},
			mockFunc: func(mockRepo *mocks.MockRepository, mockUtils *mocks.MockPoWUtilsInterface, mockLog *loggermocks.MockLogger, mockRandom *randommocks.MockProvider) {
				// Set up logger expectations
				mockLog.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().With(gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().WithContext(gomock.Any()).Return(mockLog).AnyTimes()

				mockRepo.EXPECT().
					IncrementUnsolvedCount(gomock.Any(), testClientIP).
					Return(testUnsolvedCount, nil)

				mockRepo.EXPECT().
					GetAndIncrementRequestCount(gomock.Any(), testClientIP).
					Return(testRequestCount, nil)

				mockUtils.EXPECT().
					CalculateDifficultyLevel(testRequestCount).
					Return(testDifficultyLevel)

				mockUtils.EXPECT().
					GenerateRandomSeed().
					Return(testSeed, nil)

				mockRepo.EXPECT().
					CreateTask(gomock.Any(), gomock.Any(), testChallengeTTL).
					DoAndReturn(func(_ context.Context, task powrepo.Task, ttl time.Duration) error {
						assert.Equal(t, testSeed, task.Seed)
						assert.Equal(t, testDifficultyLevel, task.DifficultyLevel)
						assert.Equal(t, testChallengeTTL, ttl)
						return nil
					})
			},
			shouldFail: false,
		},
		{
			name:     "should_fail_when_too_many_unsolved_challenges",
			clientIP: testClientIP,
			expected: nil,
			mockFunc: func(mockRepo *mocks.MockRepository, mockUtils *mocks.MockPoWUtilsInterface, mockLog *loggermocks.MockLogger, mockRandom *randommocks.MockProvider) {
				// Set up logger expectations
				mockLog.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().With(gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().WithContext(gomock.Any()).Return(mockLog).AnyTimes()

				mockRepo.EXPECT().
					IncrementUnsolvedCount(gomock.Any(), testClientIP).
					Return(testTooManyUnsolved, nil)
			},
			shouldFail: true,
		},
		{
			name:     "should_fail_when_increment_unsolved_count_fails",
			clientIP: testClientIP,
			expected: nil,
			mockFunc: func(mockRepo *mocks.MockRepository, mockUtils *mocks.MockPoWUtilsInterface, mockLog *loggermocks.MockLogger, mockRandom *randommocks.MockProvider) {
				// Set up logger expectations
				mockLog.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().With(gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().WithContext(gomock.Any()).Return(mockLog).AnyTimes()

				mockRepo.EXPECT().
					IncrementUnsolvedCount(gomock.Any(), testClientIP).
					Return(int64(0), errors.New(errDB))
			},
			shouldFail: true,
		},
		{
			name:     "should_fail_when_get_request_count_fails",
			clientIP: testClientIP,
			expected: nil,
			mockFunc: func(mockRepo *mocks.MockRepository, mockUtils *mocks.MockPoWUtilsInterface, mockLog *loggermocks.MockLogger, mockRandom *randommocks.MockProvider) {
				// Set up logger expectations
				mockLog.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().With(gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().WithContext(gomock.Any()).Return(mockLog).AnyTimes()

				mockRepo.EXPECT().
					IncrementUnsolvedCount(gomock.Any(), testClientIP).
					Return(testUnsolvedCount, nil)

				mockRepo.EXPECT().
					GetAndIncrementRequestCount(gomock.Any(), testClientIP).
					Return(int64(0), errors.New(errDB))
			},
			shouldFail: true,
		},
		{
			name:     "should_fail_when_generate_seed_fails",
			clientIP: testClientIP,
			expected: nil,
			mockFunc: func(mockRepo *mocks.MockRepository, mockUtils *mocks.MockPoWUtilsInterface, mockLog *loggermocks.MockLogger, mockRandom *randommocks.MockProvider) {
				// Set up logger expectations
				mockLog.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().With(gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().WithContext(gomock.Any()).Return(mockLog).AnyTimes()

				mockRepo.EXPECT().
					IncrementUnsolvedCount(gomock.Any(), testClientIP).
					Return(testUnsolvedCount, nil)

				mockRepo.EXPECT().
					GetAndIncrementRequestCount(gomock.Any(), testClientIP).
					Return(testRequestCount, nil)

				mockUtils.EXPECT().
					CalculateDifficultyLevel(testRequestCount).
					Return(testDifficultyLevel)

				mockUtils.EXPECT().
					GenerateRandomSeed().
					Return("", errors.New(errRandom))
			},
			shouldFail: true,
		},
		{
			name:     "should_fail_when_create_task_fails",
			clientIP: testClientIP,
			expected: nil,
			mockFunc: func(mockRepo *mocks.MockRepository, mockUtils *mocks.MockPoWUtilsInterface, mockLog *loggermocks.MockLogger, mockRandom *randommocks.MockProvider) {
				// Set up logger expectations
				mockLog.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().With(gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().WithContext(gomock.Any()).Return(mockLog).AnyTimes()

				mockRepo.EXPECT().
					IncrementUnsolvedCount(gomock.Any(), testClientIP).
					Return(testUnsolvedCount, nil)

				mockRepo.EXPECT().
					GetAndIncrementRequestCount(gomock.Any(), testClientIP).
					Return(testRequestCount, nil)

				mockUtils.EXPECT().
					CalculateDifficultyLevel(testRequestCount).
					Return(testDifficultyLevel)

				mockUtils.EXPECT().
					GenerateRandomSeed().
					Return(testSeed, nil)

				mockRepo.EXPECT().
					CreateTask(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New(errDB))
			},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mocks.NewMockRepository(ctrl)
			mockUtils := mocks.NewMockPoWUtilsInterface(ctrl)
			log := loggermocks.NewMockLogger(ctrl)
			mockRandom := randommocks.NewMockProvider(ctrl)

			// Create config
			config := &env.PowConfig{
				ScryptN:                       testScryptN,
				ScryptR:                       testScryptR,
				ScryptP:                       testScryptP,
				KeyLen:                        testKeyLen,
				MaxUnsolvedChallenges:         testMaxUnsolved,
				ChallengeTTL:                  testChallengeTTL,
				RequestsPerDifficultyIncrease: testRequestsPerDiff,
				MaxDifficultyLevel:            testMaxDiffLevel,
			}

			// Set up mock expectations
			tt.mockFunc(mockRepo, mockUtils, log, mockRandom)

			// Create service with mocks
			service := &Service{
				repo:                  mockRepo,
				config:                config,
				utils:                 mockUtils,
				maxUnsolvedChallenges: config.MaxUnsolvedChallenges,
				logger:                log,
				random:                mockRandom,
			}

			// Call the method
			challenge, err := service.CreateChallenge(context.Background(), tt.clientIP)

			// Assert results
			if tt.shouldFail {
				assert.Error(t, err)
				assert.Nil(t, challenge)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, challenge)
				// We can't check the ID directly as it's generated with uuid.New()
				// but we can check other fields
				assert.Equal(t, tt.expected.Seed, challenge.Seed)
				assert.Equal(t, tt.expected.DifficultyLevel, challenge.DifficultyLevel)
				assert.Equal(t, tt.expected.ScryptN, challenge.ScryptN)
				assert.Equal(t, tt.expected.ScryptR, challenge.ScryptR)
				assert.Equal(t, tt.expected.ScryptP, challenge.ScryptP)
				assert.Equal(t, tt.expected.KeyLen, challenge.KeyLen)
			}
		})
	}
}

func TestService_ValidateChallenge(t *testing.T) {
	tests := []struct {
		name         string
		clientIP     string
		challengeID  string
		nonce        string
		mockFunc     func(mockRepo *mocks.MockRepository, mockUtils *mocks.MockPoWUtilsInterface, mockLog *loggermocks.MockLogger, mockRandom *randommocks.MockProvider)
		shouldFail   bool
		errorMessage string
	}{
		{
			name:        "should_validate_challenge_successfully",
			clientIP:    testClientIP,
			challengeID: testValidUUID,
			nonce:       testNonce,
			mockFunc: func(mockRepo *mocks.MockRepository, mockUtils *mocks.MockPoWUtilsInterface, mockLog *loggermocks.MockLogger, mockRandom *randommocks.MockProvider) {
				// Set up logger expectations
				mockLog.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().With(gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().WithContext(gomock.Any()).Return(mockLog).AnyTimes()

				mockRepo.EXPECT().
					GetTask(gomock.Any(), testValidUUID).
					Return(&powrepo.Task{
						ID:              testValidUUID,
						Seed:            testSeed,
						DifficultyLevel: testDifficultyLevel,
					}, nil)

				mockUtils.EXPECT().
					VerifySolution(testValidUUID, testSeed, testNonce, testDifficultyLevel).
					Return(true)

				// Expect DeleteTask to be called to prevent replay attacks
				mockRepo.EXPECT().
					DeleteTask(gomock.Any(), testValidUUID).
					Return(nil)

				// Set up random.Float32() to return a value > 0.05 to test regular decrement path
				mockRandom.EXPECT().
					Float32().
					Return(float32(0.1)) // 10% > 5%, so not lucky
				mockRepo.EXPECT().
					DecrementUnsolvedCount(gomock.Any(), testClientIP).
					Return(nil)
			},
			shouldFail: false,
		},
		{
			name:        "should_fail_with_invalid_uuid",
			clientIP:    testClientIP,
			challengeID: testInvalidUUID,
			nonce:       testNonce,
			mockFunc: func(mockRepo *mocks.MockRepository, mockUtils *mocks.MockPoWUtilsInterface, mockLog *loggermocks.MockLogger, mockRandom *randommocks.MockProvider) {
				// Set up logger expectations
				mockLog.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().With(gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().WithContext(gomock.Any()).Return(mockLog).AnyTimes()

				// No expectations - validation fails before repo call
			},
			shouldFail:   true,
			errorMessage: "invalid challenge ID format",
		},
		{
			name:        "should_fail_when_get_task_fails",
			clientIP:    testClientIP,
			challengeID: testValidUUID,
			nonce:       testNonce,
			mockFunc: func(mockRepo *mocks.MockRepository, mockUtils *mocks.MockPoWUtilsInterface, mockLog *loggermocks.MockLogger, mockRandom *randommocks.MockProvider) {
				// Set up logger expectations
				mockLog.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().With(gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().WithContext(gomock.Any()).Return(mockLog).AnyTimes()

				mockRepo.EXPECT().
					GetTask(gomock.Any(), testValidUUID).
					Return(nil, errors.New(errNotFound))
			},
			shouldFail:   true,
			errorMessage: "invalid challenge",
		},
		{
			name:        "should_fail_with_invalid_solution",
			clientIP:    testClientIP,
			challengeID: testValidUUID,
			nonce:       testNonce,
			mockFunc: func(mockRepo *mocks.MockRepository, mockUtils *mocks.MockPoWUtilsInterface, mockLog *loggermocks.MockLogger, mockRandom *randommocks.MockProvider) {
				// Set up logger expectations
				mockLog.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().With(gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().WithContext(gomock.Any()).Return(mockLog).AnyTimes()

				mockRepo.EXPECT().
					GetTask(gomock.Any(), testValidUUID).
					Return(&powrepo.Task{
						ID:              testValidUUID,
						Seed:            testSeed,
						DifficultyLevel: testDifficultyLevel,
					}, nil)

				mockUtils.EXPECT().
					VerifySolution(testValidUUID, testSeed, testNonce, testDifficultyLevel).
					Return(false)
			},
			shouldFail:   true,
			errorMessage: "invalid challenge solution",
		},
		{
			name:        "should_fail_when_decrement_unsolved_count_fails",
			clientIP:    testClientIP,
			challengeID: testValidUUID,
			nonce:       testNonce,
			mockFunc: func(mockRepo *mocks.MockRepository, mockUtils *mocks.MockPoWUtilsInterface, mockLog *loggermocks.MockLogger, mockRandom *randommocks.MockProvider) {
				// Set up logger expectations
				mockLog.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().With(gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().WithContext(gomock.Any()).Return(mockLog).AnyTimes()

				mockRepo.EXPECT().
					GetTask(gomock.Any(), testValidUUID).
					Return(&powrepo.Task{
						ID:              testValidUUID,
						Seed:            testSeed,
						DifficultyLevel: testDifficultyLevel,
					}, nil)

				mockUtils.EXPECT().
					VerifySolution(testValidUUID, testSeed, testNonce, testDifficultyLevel).
					Return(true)

				// Expect DeleteTask to be called
				mockRepo.EXPECT().
					DeleteTask(gomock.Any(), testValidUUID).
					Return(nil)

				// Set up random.Float32() to return a value > 0.05 to test regular decrement path
				mockRandom.EXPECT().
					Float32().
					Return(float32(0.1)) // 10% > 5%, so not lucky

				mockRepo.EXPECT().
					DecrementUnsolvedCount(gomock.Any(), testClientIP).
					Return(errors.New(errDB))
			},
			shouldFail:   true,
			errorMessage: "failed to update challenge count",
		},
		{
			name:        "should_handle_lucky_client_with_double_decrement",
			clientIP:    testClientIP,
			challengeID: testValidUUID,
			nonce:       testNonce,
			mockFunc: func(mockRepo *mocks.MockRepository, mockUtils *mocks.MockPoWUtilsInterface, mockLog *loggermocks.MockLogger, mockRandom *randommocks.MockProvider) {
				// Set up logger expectations
				mockLog.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().With(gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().WithContext(gomock.Any()).Return(mockLog).AnyTimes()

				mockRepo.EXPECT().
					GetTask(gomock.Any(), testValidUUID).
					Return(&powrepo.Task{
						ID:              testValidUUID,
						Seed:            testSeed,
						DifficultyLevel: testDifficultyLevel,
					}, nil)

				mockUtils.EXPECT().
					VerifySolution(testValidUUID, testSeed, testNonce, testDifficultyLevel).
					Return(true)

				// Expect DeleteTask to be called
				mockRepo.EXPECT().
					DeleteTask(gomock.Any(), testValidUUID).
					Return(nil)

				// Set up random.Float32() to return a value < 0.05 to test lucky path
				mockRandom.EXPECT().
					Float32().
					Return(float32(0.01)) // 1% < 5%, so lucky

				// Expect DecrementUnsolvedCountBy to be called with 2 for lucky clients
				mockRepo.EXPECT().
					DecrementUnsolvedCountBy(gomock.Any(), testClientIP, 2).
					Return(nil)
			},
			shouldFail: false,
		},
		{
			name:        "should_handle_delete_task_failure",
			clientIP:    testClientIP,
			challengeID: testValidUUID,
			nonce:       testNonce,
			mockFunc: func(mockRepo *mocks.MockRepository, mockUtils *mocks.MockPoWUtilsInterface, mockLog *loggermocks.MockLogger, mockRandom *randommocks.MockProvider) {
				// Set up logger expectations
				mockLog.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().With(gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().WithContext(gomock.Any()).Return(mockLog).AnyTimes()

				mockRepo.EXPECT().
					GetTask(gomock.Any(), testValidUUID).
					Return(&powrepo.Task{
						ID:              testValidUUID,
						Seed:            testSeed,
						DifficultyLevel: testDifficultyLevel,
					}, nil)

				mockUtils.EXPECT().
					VerifySolution(testValidUUID, testSeed, testNonce, testDifficultyLevel).
					Return(true)

				// DeleteTask fails but validation should continue
				mockRepo.EXPECT().
					DeleteTask(gomock.Any(), testValidUUID).
					Return(errors.New(errDB))

				// Set up random.Float32() to return a value > 0.05 to test regular decrement path
				mockRandom.EXPECT().
					Float32().
					Return(float32(0.1)) // 10% > 5%, so not lucky

				// Regular decrement should still happen
				mockRepo.EXPECT().
					DecrementUnsolvedCount(gomock.Any(), testClientIP).
					Return(nil)
			},
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepository(ctrl)
			mockUtils := mocks.NewMockPoWUtilsInterface(ctrl)
			log := loggermocks.NewMockLogger(ctrl)
			mockRandom := randommocks.NewMockProvider(ctrl)

			// Create config
			config := &env.PowConfig{
				ScryptN:                       testScryptN,
				ScryptR:                       testScryptR,
				ScryptP:                       testScryptP,
				KeyLen:                        testKeyLen,
				MaxUnsolvedChallenges:         testMaxUnsolved,
				ChallengeTTL:                  testChallengeTTL,
				RequestsPerDifficultyIncrease: testRequestsPerDiff,
				MaxDifficultyLevel:            testMaxDiffLevel,
			}

			// Set up mock expectations
			tt.mockFunc(mockRepo, mockUtils, log, mockRandom)

			// Create service with mocks
			service := &Service{
				repo:                  mockRepo,
				config:                config,
				utils:                 mockUtils,
				maxUnsolvedChallenges: config.MaxUnsolvedChallenges,
				logger:                log,
				random:                mockRandom,
			}

			// Call the method
			err := service.ValidateChallenge(context.Background(), tt.clientIP, tt.challengeID, tt.nonce)

			// Assert results
			if tt.shouldFail {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
