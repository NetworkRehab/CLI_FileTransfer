package main

import (
    "fmt"
    "github.com/charmbracelet/bubbles/textinput"
    tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    protocols        []string
    cursor           int
    state            string
    sourceInput      textinput.Model
    destinationInput textinput.Model
    err              error
    transferComplete bool
}

func initialModel() model {
    m := model{
        protocols:        []string{"azureblob", "cifs", "sftp", "s3", "local"},
        state:            "selecting_protocol",
        transferComplete: false,
    }
    m.sourceInput = textinput.New()
    m.sourceInput.Placeholder = "Enter source path"
    m.sourceInput.Focus()
    m.sourceInput.CharLimit = 150
    m.sourceInput.Width = 50

    m.destinationInput = textinput.New()
    m.destinationInput.Placeholder = "Enter destination path"
    m.destinationInput.CharLimit = 150
    m.destinationInput.Width = 50
    return m
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch m.state {
    case "selecting_protocol":
        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
            case "up":
                if m.cursor > 0 {
                    m.cursor--
                }
            case "down":
                if m.cursor < len(m.protocols)-1 {
                    m.cursor++
                }
            case "enter":
                m.sourceInput.Focus()
                m.state = "entering_source"
            case "esc":
                return m, tea.Quit
            case "ctrl+c", "q":
                return m, tea.Quit
            }
        }
    case "entering_source":
        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
            case "esc":
                m.state = "selecting_protocol"
                m.sourceInput.Reset()
                return m, nil
            case "enter":
                if m.sourceInput.Value() != "" {
                    m.sourceInput.Blur()
                    m.destinationInput.Focus()
                    m.state = "entering_destination"
                }
            }
        }
        var cmd tea.Cmd
        m.sourceInput, cmd = m.sourceInput.Update(msg)
        return m, cmd
    case "entering_destination":
        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
            case "esc":
                m.state = "entering_source"
                m.destinationInput.Reset()
                m.sourceInput.Focus()
                return m, nil
            case "enter":
                if m.destinationInput.Value() != "" {
                    m.destinationInput.Blur()
                    m.state = "transferring"
                    go func() {
                        err := transferFile(m.protocols[m.cursor], m.sourceInput.Value(), m.destinationInput.Value())
                        if err != nil {
                            m.err = err
                        }
                        m.transferComplete = true
                    }()
                }
            }
        }
        var cmd tea.Cmd
        m.destinationInput, cmd = m.destinationInput.Update(msg)
        return m, cmd
    case "transferring":
        if m.transferComplete {
            if m.err != nil {
                return m, tea.Quit
            }
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m model) View() string {
    var s string

    s += "File Transfer Utility\n\n"

    switch m.state {
    case "selecting_protocol":
        s += "Select a protocol (↑/↓ to move, enter to select):\n\n"
        for i, choice := range m.protocols {
            cursor := " "
            if m.cursor == i {
                cursor = ">"
            }
            s += fmt.Sprintf("%s %s\n", cursor, choice)
        }
        s += "\n\nPress 'q' or Ctrl+C to quit"
        return s
    case "entering_source":
        s += "Enter source path (esc to go back):\n"
        s += m.sourceInput.View()
        s += "\n\nPress 'q' or Ctrl+C to quit"
        return s
    case "entering_destination":
        s += "Enter destination path (esc to go back):\n"
        s += m.destinationInput.View()
        s += "\n\nPress 'q' or Ctrl+C to quit"
        return s
    case "transferring":
        if m.transferComplete {
            if m.err != nil {
                return fmt.Sprintf("Error: %v\nPress any key to exit.", m.err)
            }
            return "Transfer completed successfully!\nPress any key to exit."
        }
        return "Transferring file...\nPlease wait..."
    }
    return ""
}
