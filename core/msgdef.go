package core

type MsgRequest struct {
	id         int64
	sender     string
	topic      string
	account    string
	recipients []string
	canReply   bool
	data       interface{}
}

type MsgResponse struct {
	id        int64
	sender    string
	recipient string
	account   string
	data      interface{}
}

