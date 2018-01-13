package acl

func (a *ACL) AllowNamespaceRead() bool {
	switch {
	case a.management:
		return true
	case a.node == PolicyWrite:
		return true
	case a.node == PolicyRead:
		return true
	default:
		return false
	}
}
