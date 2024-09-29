package enums

type Status int

const (
	Queue Status = iota + 1
	Process
	Accepted
	WrongAnswer
	TimeLimitExceeded
	CompilationError
	RuntimeErrorSIGSEGV
	RuntimeErrorSIGXFSZ
	RuntimeErrorSIGFPE
	RuntimeErrorSIGABRT
	RuntimeErrorNZEC
	RuntimeErrorOther
	InternalError
	ExecFormatError
)

type StatusInfo struct {
	ID   int
	Name string
}

var statusMap = map[Status]StatusInfo{
	Queue:               {1, "In Queue"},
	Process:             {2, "Processing"},
	Accepted:            {3, "Accepted"},
	WrongAnswer:         {4, "Wrong Answer"},
	TimeLimitExceeded:   {5, "Time Limit Exceeded"},
	CompilationError:    {6, "Compilation Error"},
	RuntimeErrorSIGSEGV: {7, "Runtime Error (SIGSEGV)"},
	RuntimeErrorSIGXFSZ: {8, "Runtime Error (SIGXFSZ)"},
	RuntimeErrorSIGFPE:  {9, "Runtime Error (SIGFPE)"},
	RuntimeErrorSIGABRT: {10, "Runtime Error (SIGABRT)"},
	RuntimeErrorNZEC:    {11, "Runtime Error (NZEC)"},
	RuntimeErrorOther:   {12, "Runtime Error (Other)"},
	InternalError:       {13, "Internal Error"},
	ExecFormatError:     {14, "Exec Format Error"},
}

func FindRuntimeErrorByStatusCode(statusCode int) Status {
	switch statusCode {
	case 11:
		return RuntimeErrorSIGSEGV
	case 25:
		return RuntimeErrorSIGXFSZ
	case 8:
		return RuntimeErrorSIGFPE
	case 6:
		return RuntimeErrorSIGABRT
	default:
		return RuntimeErrorOther
	}
}

func GetStatusInfo(status Status) (StatusInfo, bool) {
	info, found := statusMap[status]
	return info, found
}
