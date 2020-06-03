package models

type HostEntry struct {
	Domain  string
	A       string
	AOld    string
	AAAA    string
	AAAAOld string
	Ttl     int
}

type Zone struct {
	Name        string
	Hostentries []HostEntry
}
