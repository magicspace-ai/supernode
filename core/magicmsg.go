package core

type MsgRequest struct {
	id         string
	sender     string
	topic      string
	account    string
	recipients []string
	doReply    bool
	data       interface{}
}

type MsgResponse struct {
	id        string
	sender    string
	recipient string
	account   string
	data      interface{}
}