// secret service implementation according to:
// http://standards.freedesktop.org/secret-service
package service

// create and initialize a new session
func NewSession(parent *Service) *Session {
	session := &Session{}
	session.Parent = parent
	// Sessions don't need to get persistent in db so no need for 'Update'
	return session
}

func NewCliSession(parent *Service) *CliSession {
	CliSession := &CliSession{}
	CliSession.Parent = parent
	return CliSession
}

// CreateMethodFromPath returns a.b.c.Foo when session path
// is /a/b/c/xyz and passed method is 'Foo'
func (s *Session) CreateMethodFromPath(method string) string {
	_, child := Path2Name(string(s.ObjectPath), method)
	return child
}
