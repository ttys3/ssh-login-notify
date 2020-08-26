package pam

import "os"

/*
	pam_exec is a PAM module that can be used to run an external command.
	The child's environment is set to the current PAM environment list, as returned by pam_getenvlist(3) In addition,
	the following PAM items are exported as environment variables:
	PAM_RHOST, PAM_RUSER, PAM_SERVICE, PAM_TTY, PAM_USER and PAM_TYPE,
	which contains one of the module types: account, auth, password, open_session and close_session.
	Commands called by pam_exec need to be aware of that the user can have control over the environment.

	http://www.linux-pam.org/Linux-PAM-html/sag-pam_exec.html
*/
type PAMEnv struct {
	PAM_RHOST string
	PAM_RUSER string
	PAM_SERVICE string
	PAM_TTY string
	PAM_USER string
	PAM_TYPE string
}

func NewPAMEnv() *PAMEnv {
	return &PAMEnv{}
}

func (p *PAMEnv) Init() *PAMEnv {
	p.PAM_RHOST = os.Getenv("PAM_RHOST")
	p.PAM_RUSER = os.Getenv("PAM_RUSER")
	p.PAM_SERVICE = os.Getenv("PAM_SERVICE")
	p.PAM_TTY = os.Getenv("PAM_TTY")
	p.PAM_USER = os.Getenv("PAM_USER")
	p.PAM_TYPE = os.Getenv("PAM_TYPE")
	return p
}