package main

import (
	"fmt"
	"os"
	"time"

	timer "github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case intervalTick:
		if m.status == started {
			m.elapsed += 1
			return m, m.tick()
		}

	case timer.TickMsg:
		var cmd tea.Cmd

		if msg.ID == m.showAllTimer.ID() {
			m.showAllTimer, cmd = m.showAllTimer.Update(msg)
			if m.showAllTimer.Timedout() {
				for i := range m.board.cells {
					m.board.cells[i].matched = false
				}
			}
		}

		if msg.ID == m.wrongTimer.ID() {
			m.wrongTimer, cmd = m.wrongTimer.Update(msg)
			if m.wrongTimer.Timedout() {
				m.board = m.board.match()
			}
		}

		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "s":
			return m.restart()

		case "enter", " ":
			if m.status == started {
				if !m.board.isCellAlreadyMatched() {
					if len(m.board.selection) == 0 {
						m.board = m.board.selectCell()
					} else {
						if m.wrongTimer.Timedout() {
							m.board = m.board.selectCell()
							if m.board.isSelectionMatched() {
								m.board = m.board.match()
								if m.board.isCompleted() {
									return m.completed()
								}
							} else {
								m.wrongs += 1
								m.wrongTimer = timer.NewWithInterval(time.Second, 100*time.Millisecond)
								return m, m.wrongTimer.Init()
							}
						}
					}
				}
			}

		case "h":
			if m.status == started {
				if m.board.cursor-1 >= 0 {
					m.board.cursor--
				}
			}

		case "l":
			if m.status == started {
				if m.board.cursor+1 < len(m.board.cells) {
					m.board.cursor++
				}
			}

		case "k":
			if m.status == started {
				if m.board.cursor-m.board.cols >= 0 {
					m.board.cursor -= m.board.cols
				}
			}

		case "j":
			if m.status == started {
				if m.board.cursor+m.board.cols < len(m.board.cells) {
					m.board.cursor += m.board.cols
				}
			}
		}
	}
	return m, nil
}

func (m App) View() string {
	var s string
	if m.status == started {
		s = m.board.view()
		s += fmt.Sprintf("     Wrong choices: %d\n", m.wrongs)
		s += fmt.Sprintf("     Time elapsed: %ds\n", m.elapsed)
	} else {
		if m.status == completed {
			s += fmt.Sprintf("     Wrong choices: %d\n", m.wrongs)
			s += fmt.Sprintf("     Time elapsed: %ds\n", m.elapsed)
			s += "-----------------\n\n"
		}
		s += m.menuView()
	}
	return s
}

func main() {
	p := tea.NewProgram(NewApp())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
