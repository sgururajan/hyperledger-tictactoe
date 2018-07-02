package providers

type predicate func() bool
type applier func()

type setter struct {
	isSet bool
}

func (s *setter) set(current interface{}, check predicate, apply applier) {
	if current == nil && (check == nil || check()) {
		apply()
		s.isSet = true
	}
}
