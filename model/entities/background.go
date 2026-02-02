package entities

type BackgroundRecord struct {
	ID        string
	Approvers []string
	Refusers  []string
	Role      string
	Sender    string
}
