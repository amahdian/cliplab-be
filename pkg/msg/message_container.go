package msg

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/amahdian/cliplab-be/global/errs"
	"github.com/samber/lo"
)

type MessageContainer struct {
	fatalMessageGroups   map[string]*MessageGroup
	errorMessageGroups   map[string]*MessageGroup
	warningMessageGroups map[string]*MessageGroup
	infoMessageGroups    map[string]*MessageGroup
	hasError             bool
}

func NewMessageContainer() *MessageContainer {
	return &MessageContainer{
		fatalMessageGroups:   make(map[string]*MessageGroup),
		errorMessageGroups:   make(map[string]*MessageGroup),
		warningMessageGroups: make(map[string]*MessageGroup),
		infoMessageGroups:    make(map[string]*MessageGroup),
	}
}

func (mc *MessageContainer) Union(otherMessageContainers ...*MessageContainer) {
	for _, otherMessageContainer := range otherMessageContainers {
		for group, groupMessages := range otherMessageContainer.fatalMessageGroups {
			for _, message := range groupMessages.Messages {
				mc.AddFatal(group, message.Text)
			}
		}
		for group, groupMessages := range otherMessageContainer.errorMessageGroups {
			for _, message := range groupMessages.Messages {
				mc.AddError(group, message.Text)
			}
		}
		for group, groupMessages := range otherMessageContainer.warningMessageGroups {
			for _, message := range groupMessages.Messages {
				mc.AddWarning(group, message.Text)
			}
		}
		for group, groupMessages := range otherMessageContainer.infoMessageGroups {
			for _, message := range groupMessages.Messages {
				mc.AddInfo(group, message.Text)
			}
		}
	}
}

func (mc *MessageContainer) Count() int {
	count := 0
	count += countMessages(mc.fatalMessageGroups)
	count += countMessages(mc.errorMessageGroups)
	count += countMessages(mc.warningMessageGroups)
	count += countMessages(mc.infoMessageGroups)
	return count
}

func (mc *MessageContainer) ErrorCount() int {
	return countMessages(mc.errorMessageGroups)
}

func (mc *MessageContainer) WarningCount() int {
	return countMessages(mc.warningMessageGroups)
}

func (mc *MessageContainer) InfoCount() int {
	return countMessages(mc.infoMessageGroups)
}

func (mc *MessageContainer) HasError() bool {
	return mc.hasError
}

func (mc *MessageContainer) GetAll() []*MessageGroup {
	var messageGroupsArray []*MessageGroup
	//
	messageGroupsArray = append(messageGroupsArray, mc.getMessageGroupsByLevel(Fatal)...)
	messageGroupsArray = append(messageGroupsArray, mc.getMessageGroupsByLevel(Error)...)
	messageGroupsArray = append(messageGroupsArray, mc.getMessageGroupsByLevel(Warning)...)
	messageGroupsArray = append(messageGroupsArray, mc.getMessageGroupsByLevel(Info)...)
	//
	return messageGroupsArray
}

func (mc *MessageContainer) getMessagesByLevel(level MessageLevel) []*Message {
	var messageArray []*Message
	switch level {
	case Fatal:
		for _, messageGroups := range mc.fatalMessageGroups {
			messageArray = append(messageArray, messageGroups.Messages...)
		}
	case Error:
		for _, messageGroups := range mc.errorMessageGroups {
			messageArray = append(messageArray, messageGroups.Messages...)
		}
	case Warning:
		for _, messageGroups := range mc.warningMessageGroups {
			messageArray = append(messageArray, messageGroups.Messages...)
		}
	case Info:
		for _, messageGroups := range mc.infoMessageGroups {
			messageArray = append(messageArray, messageGroups.Messages...)
		}
	}
	//
	return messageArray
}

func (mc *MessageContainer) getMessageGroupsByLevel(level MessageLevel) []*MessageGroup {
	var groups []*MessageGroup
	switch level {
	case Fatal:
		groups = lo.Values(mc.fatalMessageGroups)
	case Error:
		groups = lo.Values(mc.errorMessageGroups)
	case Warning:
		groups = lo.Values(mc.warningMessageGroups)
	case Info:
		groups = lo.Values(mc.infoMessageGroups)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Order < groups[j].Order
	})
	return groups
}

func (mc *MessageContainer) GetErrors() []*Message {
	return mc.getMessagesByLevel(Error)
}

func (mc *MessageContainer) GetWarnings() []*Message {
	return mc.getMessagesByLevel(Warning)
}

func (mc *MessageContainer) GetInfos() []*Message {
	return mc.getMessagesByLevel(Info)
}

func (mc *MessageContainer) AddErr(err error) *MessageContainer {
	if err == nil {
		return mc
	}
	var customErr *errs.Error
	if errors.As(err, &customErr) {
		mc.AddError(customErr.Code.MessageGroup(), customErr.Error())
	} else {
		mc.AddError("Internal issue", err.Error())
	}
	return mc
}

func (mc *MessageContainer) AddFatal(group string, messageText string) {
	addMessage(mc.fatalMessageGroups, group, Fatal, messageText)
	mc.hasError = true
}

func (mc *MessageContainer) AddError(group string, messageText string) *MessageContainer {
	addMessage(mc.errorMessageGroups, group, Error, messageText)
	mc.hasError = true
	return mc
}

func (mc *MessageContainer) AddWarning(group string, messageText string) *MessageContainer {
	addMessage(mc.warningMessageGroups, group, Warning, messageText)
	return mc
}

func (mc *MessageContainer) AddInfo(group string, messageText string) *MessageContainer {
	addMessage(mc.infoMessageGroups, group, Info, messageText)
	return mc
}

func (mc *MessageContainer) AddErrorf(group string, format string, parameters ...any) *MessageContainer {
	messageText := fmt.Sprintf(format, parameters...)
	addMessage(mc.errorMessageGroups, group, Error, messageText)
	mc.hasError = true
	return mc
}

func (mc *MessageContainer) AddWarningf(group string, format string, parameters ...any) *MessageContainer {
	messageText := fmt.Sprintf(format, parameters...)
	addMessage(mc.warningMessageGroups, group, Warning, messageText)
	return mc
}

func (mc *MessageContainer) AddInfof(group string, format string, parameters ...any) *MessageContainer {
	messageText := fmt.Sprintf(format, parameters...)
	addMessage(mc.infoMessageGroups, group, Info, messageText)
	return mc
}

func MakePlainText(message string) string {
	newMessage := strings.ReplaceAll(message, string('"'), "`")
	newMessage = strings.ReplaceAll(message, "'", "`")
	newMessage = strings.ReplaceAll(newMessage, "\r\n", " ")
	newMessage = strings.ReplaceAll(newMessage, "\n", " ")
	newMessage = strings.ReplaceAll(newMessage, "\r", " ")
	newMessage = strings.ReplaceAll(newMessage, "\t", " ")
	return newMessage
}

func addMessage(messageGroups map[string]*MessageGroup, group string, level MessageLevel, text string) {
	if _, ok := messageGroups[group]; !ok {
		messageGroups[group] = &MessageGroup{
			Group:    group,
			Messages: make([]*Message, 0),
			Order:    len(messageGroups),
		}
	}
	messageGroup := messageGroups[group]
	messageGroup.Messages = append(messageGroup.Messages, &Message{
		Text:  text,
		Level: level,
	})
}

func countMessages(messageGroupMap map[string]*MessageGroup) int {
	count := 0
	for _, messageGroup := range messageGroupMap {
		count += len(messageGroup.Messages)
	}
	return count
}
