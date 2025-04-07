package getquote

import (
	"context"
	"testing"

	"github.com/opravdin/word-of-wisdom/internal/domain"
	loggermocks "github.com/opravdin/word-of-wisdom/internal/logger/mocks"
	"github.com/opravdin/word-of-wisdom/internal/usecase/getquote/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// Using the logger mock from internal/logger/mocks as per testing rules

func TestDefaultUsecase_Quote(t *testing.T) {
	// Define test cases
	tests := []struct {
		name       string
		input      context.Context
		expected   domain.Quote
		mockFunc   func(mockRepo *mocks.MockquotesRepository, mockLog *loggermocks.MockLogger)
		shouldFail bool
	}{
		{
			name:  "should_return_quote_from_repository",
			input: context.Background(),
			expected: domain.Quote{
				Text:   "The only true wisdom is in knowing you know nothing.",
				Author: "Socrates",
			},
			mockFunc: func(mockRepo *mocks.MockquotesRepository, mockLog *loggermocks.MockLogger) {
				// Set up logger expectations
				mockLog.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().With(gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().WithContext(gomock.Any()).Return(mockLog).AnyTimes()

				mockRepo.EXPECT().
					GetRandom(gomock.Any()).
					Return(domain.Quote{
						Text:   "The only true wisdom is in knowing you know nothing.",
						Author: "Socrates",
					})
			},
			shouldFail: false,
		},
		{
			name:     "should_return_empty_quote_if_repository_returns_empty",
			input:    context.Background(),
			expected: domain.Quote{},
			mockFunc: func(mockRepo *mocks.MockquotesRepository, mockLog *loggermocks.MockLogger) {
				// Set up logger expectations
				mockLog.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
				mockLog.EXPECT().With(gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().WithContext(gomock.Any()).Return(mockLog).AnyTimes()

				mockRepo.EXPECT().
					GetRandom(gomock.Any()).
					Return(domain.Quote{})
			},
			shouldFail: false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup controller and mock
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockquotesRepository(ctrl)
			mockLog := loggermocks.NewMockLogger(ctrl)
			tt.mockFunc(mockRepo, mockLog)

			// Create usecase with mock
			u := NewDefaultUsecase(mockRepo, mockLog)

			// Call the method
			got := u.Quote(tt.input)

			// Assert result
			if tt.shouldFail {
				// This function doesn't return an error, so we skip this check
				// In a real scenario, we would modify the function to return an error
				// and assert it here
			} else {
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
