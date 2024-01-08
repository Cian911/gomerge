package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	approve key.Binding
	close   key.Binding
	merge   key.Binding
}

func defaultKeyMappings() keyMap {
	return keyMap{
		approve: key.NewBinding(
			key.WithKeys("a"),
		),
		close: key.NewBinding(
			key.WithKeys("c"),
		),
		merge: key.NewBinding(
			key.WithKeys("m"),
		),
	}
}
