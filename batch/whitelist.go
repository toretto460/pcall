package batch

// Whitelist is a map of whitelisted items
type Whitelist map[string]bool

// WhitelistFromList creates a Whitelist from a plain list
func WhitelistFromList(list []string) Whitelist {
	w := make(Whitelist, len(list))

	for _, item := range list {
		w[item] = true
	}

	return w
}

func (w Whitelist) isWhitelisted(item string) bool {
	if _, ok := w[item]; ok {
		return true
	}

	return false
}
