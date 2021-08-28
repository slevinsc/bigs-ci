package core

type Release interface {
	PullCode()(string,error)
	Build()error
	Deploy()error
}

type EventData struct {
	ServiceName  string
	Branch       string
	CommitID     string
	GitSshUrl    string
	GitHttpUrl   string
	MainFileDir string
}
