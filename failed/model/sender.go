package model

type Sender interface {
	SendMsg(msgType int, msg []byte) error
}

type SendRequest struct {
	Sender   string
	Receiver string

	// 接收者类型，用布尔值表示，true 表示是用户，false表示是组
	ReceiverType bool
	Time         int64
}

func (info SendRequest) GetReceiver() (Sender, error) {
	if info.ReceiverType {
		return UserTable.GetMember(info.Receiver)
	}
	return GroupTable.GetMember(info.Receiver)
}
