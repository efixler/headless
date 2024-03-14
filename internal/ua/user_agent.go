package ua

const (
	Firefox88 = "Mozilla/5.0 (X11; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0"
	Safari537 = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3830.0 Safari/537.36"
)

type Arg struct {
	userAgent string
}

func (u Arg) String() string {
	return u.userAgent
}

func (u *Arg) UnmarshalText(text []byte) error {
	agent := string(text)
	switch agent {
	case ":firefox:":
		u.userAgent = Firefox88
	case ":safari:":
		u.userAgent = Safari537
	default:
		u.userAgent = agent
	}
	return nil
}
