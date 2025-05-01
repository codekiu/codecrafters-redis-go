package storage

type Information struct {
	Role               string
	Master_Replid      string
	Master_Repl_Offset string
}

func NewInformation(role string, masterReplid string, masterReplOffset string) *Information {
	return &Information{
		Role:               role,
		Master_Replid:      masterReplid,
		Master_Repl_Offset: masterReplOffset,
	}
}
