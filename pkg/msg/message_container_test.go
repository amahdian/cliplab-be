package msg

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestMessageContainer(t *testing.T) {
	tests := []struct {
		name                       string
		msgContainer               *MessageContainer
		expectedMessageGroupOrders []string
		expectedMessageTextOrders  []string
		allMessagesCount           int
		errorCount                 int
	}{
		{
			name: "Test Scenario 1",
			msgContainer: NewMessageContainer().
				AddError("Rig Sheet is not found", "1.Rig Sheet is not found").
				AddError("Internal Error", "2.Internal Error").
				AddError("Internal Error", "3.Internal Error").
				AddError("Internal Error", "4.Internal Error").
				AddWarning("Zero issue value", "5.Zero issue value").
				AddError("X Internal Error", "6.X Internal Error"),
			expectedMessageGroupOrders: []string{"Rig Sheet is not found", "Internal Error", "X Internal Error", "Zero issue value"},
			expectedMessageTextOrders:  []string{"1.Rig Sheet is not found", "2.Internal Error", "3.Internal Error", "4.Internal Error", "6.X Internal Error", "5.Zero issue value"},
			allMessagesCount:           6,
			errorCount:                 5,
		},
		{
			name: "Test Scenario 2",
			msgContainer: NewMessageContainer().
				AddWarning("Zero issue value", "5.Zero issue value").
				AddError("X Internal Error", "6.X Internal Error").
				AddError("Internal Error", "3.Internal Error").
				AddError("Rig Sheet is not found", "1.Rig Sheet is not found").
				AddError("Internal Error", "2.Internal Error").
				AddError("Internal Error", "4.Internal Error"),
			expectedMessageGroupOrders: []string{"X Internal Error", "Internal Error", "Rig Sheet is not found", "Zero issue value"},
			expectedMessageTextOrders:  []string{"6.X Internal Error", "3.Internal Error", "2.Internal Error", "4.Internal Error", "1.Rig Sheet is not found", "5.Zero issue value"},
			allMessagesCount:           6,
			errorCount:                 5,
		},
		{
			name: "Test Scenario 3",
			msgContainer: NewMessageContainer().
				AddError("Rig Sheet is not found", "1.Rig Sheet is not found").
				AddError("Internal Error", "2.Internal Error").
				AddError("Internal Error", "3.Internal Error").
				AddError("X Internal Error", "6.X Internal Error").
				AddWarning("Zero issue value", "5.Zero issue value").
				AddError("Internal Error", "4.Internal Error"),
			expectedMessageGroupOrders: []string{"Rig Sheet is not found", "Internal Error", "X Internal Error", "Zero issue value"},
			expectedMessageTextOrders:  []string{"1.Rig Sheet is not found", "2.Internal Error", "3.Internal Error", "4.Internal Error", "6.X Internal Error", "5.Zero issue value"},
			allMessagesCount:           6,
			errorCount:                 5,
		},
	}
	//
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.msgContainer.GetAll()
			actualMessageGroupOrders := lo.Map(actual, func(item *MessageGroup, index int) string {
				return item.Group
			})
			require.Equal(t, tt.expectedMessageGroupOrders, actualMessageGroupOrders)
			//
			actualMessageTexts := make([]string, 0)
			for _, iterator := range actual {
				actualMessageTexts = append(actualMessageTexts, lo.Map(iterator.Messages, func(item *Message, index int) string {
					return item.Text
				})...)
			}
			require.Equal(t, tt.expectedMessageTextOrders, actualMessageTexts)
			//
			actualAllMessagesCount := tt.msgContainer.Count()
			require.Equal(t, tt.allMessagesCount, actualAllMessagesCount)
			//
			errorsCount := tt.msgContainer.ErrorCount()
			require.Equal(t, tt.errorCount, errorsCount)
		})
	}
}
