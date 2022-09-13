package main

type MessageInfo struct {
	fName     string
	lName     string
	username  string
	chatId    int64
	userId    int64
	text      string
	isGroup   bool
	messageId int
}

func NewMessageInfo(fname, lname, username string, chatId, userId int64, messageId int, isGroup bool) *MessageInfo {
	return &MessageInfo{
		fName:     fname,
		lName:     lname,
		username:  username,
		chatId:    chatId,
		userId:    userId,
		messageId: messageId,
		isGroup:   isGroup,
	}
}
