# Testing Rules

This document outlines testing standards for the project. All AI code generation agents must comply with these rules when generating test code.

## Table-Driven Tests
- Use slices of structs with test cases
- Include name, input, expected output, and mockFunc fields
- Run each case with t.Run(tt.name, func(t *testing.T) {...})

## Mocking with gomock
- Use Uber's gomock package
- Define interfaces in deps.go files
- Generate mocks with: `//go:generate mockgen -source=${GOFILE} -package=mocks -destination=./mocks/deps.go`
- Create mocks on per-test basis with new controller for each test
- Set up expectations in mockFunc
- Exception: The logger interface should be defined in internal/logger and mocks for it should be used from internal/mocks/logger

## Test Case Naming
- Use snake_case format (e.g., should_return_error_when_invalid)
- Be descriptive of scenario and expected outcome

## Code Quality
- Don't leave unnecessary comments in tests
- Avoid using magic constants in test cases; define them in const () section at the top of the file. This includes test IDs, values of structures and so on.

## Sample Test
```go
import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "go.uber.org/mock/gomock"
)

func TestService_ProcessData(t *testing.T) {
    // Define constants for test values
    const (
        validInput = "valid data"
        successStatus = "success"
        processedData = "processed"
    )
    
    tests := []struct {
        name       string
        input      string
        expected   Result
        mockFunc   func(mockRepo *mocks.MockRepository)
        shouldFail bool
    }{
        {
            name:  "should_process_valid_data_successfully",
            input: validInput,
            expected: Result{Status: successStatus, Data: processedData},
            mockFunc: func(mockRepo *mocks.MockRepository) {
                mockRepo.EXPECT().Save(gomock.Any(), validInput).Return(nil)
            },
            shouldFail: false,
        },
        {
            name:     "should_return_error_for_invalid_data",
            input:    "",
            expected: Result{},
            mockFunc: func(mockRepo *mocks.MockRepository) {
                // No expectations needed - validation fails before repo call
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
            tt.mockFunc(mockRepo)
            service := NewService(mockRepo)

            // Execute
            result, err := service.ProcessData(context.Background(), tt.input)

            // Assert
            if tt.shouldFail {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected, result)
            }
        })
    }
}