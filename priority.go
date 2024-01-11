package wrr

type Priority struct {
	Name   string
	Force  bool
	Weight int
}

var (
	LowPriority    = &Priority{"low", false, 20}
	MediumPriority = &Priority{"medium", false, 30}
	HighPriority   = &Priority{"high", false, 50}
	ForcePriority  = &Priority{"force", true, 0}
)
