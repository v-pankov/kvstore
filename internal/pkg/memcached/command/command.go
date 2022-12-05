package command

type (
	Set struct {
		Key     string
		Flags   int16
		ExpTime int
		Bytes   int
		NoReply bool
	}

	Get struct {
		Keys []string
	}

	Gat struct {
		ExpTime int
		Keys    []string
	}

	Delete struct {
		Key     string
		NoReply bool
	}
)

type Command interface {
	isCommand()
}

var (
	_ Command = (*Set)(nil)
	_ Command = (*Get)(nil)
	_ Command = (*Gat)(nil)
	_ Command = (*Delete)(nil)
)

func (*Set) isCommand()    {}
func (*Get) isCommand()    {}
func (*Gat) isCommand()    {}
func (*Delete) isCommand() {}
