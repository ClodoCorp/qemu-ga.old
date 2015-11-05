package qga

type ExecStatus struct {
	Exited   bool   `json:"exited"`
	ExitCode *int   `json:"exitcode,omitempty"`
	Signal   int    `json:"signal,omitempty"`
	OutData  string `json:"out-data,omitempty"`
	ErrData  string `json:"err-data,omitempty"`
	OutTrunc bool   `json:"out-truncated,omitempty"`
	ErrTrunc bool   `json:"err-truncated,omitempty"`
}

var execStatuses map[int]*ExecStatus

func init() {
	execStatuses = make(map[int]*ExecStatus)
}
