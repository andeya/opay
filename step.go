package opay

type (
	// handling order's action
	Step int
)

// Six order processing behavior states
const (
	FAIL      Step = UNSET - 2 //Processing failed
	CANCEL    Step = UNSET - 1 //Cancel order
	UNSET     Step = 0         //Not set
	PEND      Step = UNSET + 1 //Wait for processing
	DO        Step = UNSET + 2 //Is being processed
	SUCCEED   Step = UNSET + 3 //Processing success
	SYNC_DEAL Step = UNSET + 4 //Processing success synchronously
)

var (
	steps = map[Step]bool{
		FAIL:      true,
		CANCEL:    true,
		UNSET:     true,
		PEND:      true,
		DO:        true,
		SUCCEED:   true,
		SYNC_DEAL: true,
	}
)
