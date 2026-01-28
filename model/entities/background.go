package entities

type BackgroundRecord struct {
	ID        string
	Approvers []string
	Refusers  []string
	Sender    string
}
