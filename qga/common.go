package qga

var l *Logger

// FileSystem struct
type FileSystem struct {
	Device  string
	Path    string
	Type    string
	Options []string
}

type ExecStatus struct {
	Exited   bool   `json:"exited"`
	ExitCode *int   `json:"exitcode,omitempty"`
	Signal   int    `json:"signal,omitempty"`
	OutData  string `json:"out-data,omitempty"`
	ErrData  string `json:"err-data,omitempty"`
	OutTrunc bool   `json:"out-truncated,omitempty"`
	ErrTrunc bool   `json:"err-truncated,omitempty"`
}
