package wstf

const LocalsKeyRegisteredPushers = "_wstf.inner.registered-pushers"

func (m *Connection) SetRegisteredPushersToLocals(pushers map[string]*Pusher) {
	m.Locals[LocalsKeyRegisteredPushers] = pushers
}

// Try getting the registered pushers, which may be nil and error may be found.
func (m *Connection) TryGettingRegisteredPushers() (map[string]*Pusher, bool, bool) {
	target := m.Locals[LocalsKeyRegisteredPushers]
	isInitialized := target != nil
	if !isInitialized {
		return nil, isInitialized, true
	}
	pushers, ok := target.(map[string]*Pusher)
	return pushers, isInitialized, ok
}

// Get the map of the registered pushers, which will be initialized if it is nil; return nil if error found.
func (m *Connection) GetRegisteredPushers(res *Response) (map[string]*Pusher, bool) {
	pushers, isInitialized, ok := m.TryGettingRegisteredPushers()
	if !ok {
		// Failed to convert the target resource to the expected type.
		return nil, false
	}
	if !isInitialized {
		pushers = make(map[string]*Pusher)
		m.SetRegisteredPushersToLocals(pushers)
	}
	return pushers, true
}

func (m *Connection) GetRegisteredPusher(res *Response) *Pusher {
	pushers, ok := m.GetRegisteredPushers(res)
	if !ok {
		return nil
	}

	key := res.GetPusherId()
	if pushers[key] == nil {
		pushers[key] = res.NewPusher()
	}
	return pushers[key]
}

func (m *Connection) HasRegisteredPushers() *bool {
	pushers, isInitialized, ok := m.TryGettingRegisteredPushers()
	if !ok {
		// Initialized but Error
		return nil
	}
	if !isInitialized {
		return &isInitialized
	}
	has := len(pushers) > 0
	return &has
}

func (m *Connection) RemoveRegisteredPusher(res *Response) *Pusher {
	pushers, _, _ := m.TryGettingRegisteredPushers()
	if pushers != nil {
		return nil
	}

	key := res.GetPusherId()
	t := pushers[key]
	if pushers[key] != nil {
		delete(pushers, key)
	}

	if len(pushers) == 0 {
		// Reset the local map when it is empty.
		m.SetRegisteredPushersToLocals(nil)
	}

	return t
}
