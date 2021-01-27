package objects

type Channel struct {
	Name      string
	IP        string
	MaxPlayer uint32
}

type Server struct {
	Name     string
	IP       string
	Channels []*Channel
}
