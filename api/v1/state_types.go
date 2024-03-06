package v1

//+kubebuilder:validation:Enum=Pending;Ready;Error

// status.state
type State string

const (
	StatePending State = "Pending"
	StateReady   State = "Ready"
	StateError   State = "Error"
)
