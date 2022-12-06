package reply

type (
	Value struct {
		Key   string
		Flags int16
		Bytes int
		Cas   int64
	}

	Error struct{}

	ClientError struct {
		Error string
	}

	ServerError struct {
		Error string
	}

	End struct {
	}

	Stored struct {
	}

	Deleted struct {
	}

	NotFound struct {
	}
)

type Reply interface {
	isReply()
}

var (
	_ Reply = (*Value)(nil)
	_ Reply = (*Error)(nil)
	_ Reply = (*ClientError)(nil)
	_ Reply = (*ServerError)(nil)
	_ Reply = (*End)(nil)
	_ Reply = (*Stored)(nil)
	_ Reply = (*Deleted)(nil)
	_ Reply = (*NotFound)(nil)
)

func (*Value) isReply()       {}
func (*Error) isReply()       {}
func (*ClientError) isReply() {}
func (*ServerError) isReply() {}
func (*End) isReply()         {}
func (*Stored) isReply()      {}
func (*Deleted) isReply()     {}
func (*NotFound) isReply()    {}
