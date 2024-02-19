package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Approve key.Binding
	Close   key.Binding
	Merge   key.Binding
	Remove  key.Binding
	Open    key.Binding
	Help    key.Binding
	Quit    key.Binding
}

func defaultKeyMappings() keyMap {
	return keyMap{
		Approve: key.NewBinding(
			key.WithKeys("a"),
      key.WithHelp("a", "approve pr"),
		),
		Close: key.NewBinding(
			key.WithKeys("c"),
      key.WithHelp("c", "close pr"),
		),
		Merge: key.NewBinding(
			key.WithKeys("m"),
      key.WithHelp("m", "merge pr"),
		),
		Remove: key.NewBinding(
			key.WithKeys("r"),
      key.WithHelp("r", "remove selected"),
		),
    Open: key.NewBinding(
      key.WithKeys("o"),
      key.WithHelp("o", "open/close selected view"),
    ),
    Help: key.NewBinding(
      key.WithKeys("?"),
      key.WithHelp("?", "toggle help"),
    ),
    Quit: key.NewBinding(
      key.WithKeys("q", "esc", "ctrl+c"),
      key.WithHelp("q", "quit"),
    ),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Approve, k.Merge, k.Close, k.Open}, // first column
		{k.Help, k.Quit},                // second column
	}
}
