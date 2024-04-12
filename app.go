package main

import (
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type intervalTick time.Time

type Level int8

const (
	easy   Level = 9
	medium Level = 16
	hard   Level = 25
)

type Status int8

const (
	stopped Status = iota
	started
	completed
)

type App struct {
	board        Board
	status       Status
	level        Level
	wrongs       int
	elapsed      int
	wrongTimer   timer.Model
	showAllTimer timer.Model
	completedAt  time.Time
	startedAt    time.Time
}

func NewApp() App {
	return App{
		board:   Board{},
		status:  stopped,
		level:   medium,
		wrongs:  0,
		elapsed: 0,
	}
}

func (m App) Init() tea.Cmd {
	return nil
}

func (m App) restart() (App, tea.Cmd) {
	m.board = NewBoard(m.level)
	m.completedAt = time.Time{}
	m.elapsed = 0
	m.wrongs = 0
	m.startedAt = time.Now()
	for i := range m.board.cells {
		m.board.cells[i].matched = true
	}

	m.showAllTimer = timer.NewWithInterval(2*time.Second, 100*time.Millisecond)
	if m.status != started {
		m.status = started
		return m, tea.Batch(m.showAllTimer.Init(), m.tick())
	} else {
		return m, m.showAllTimer.Init()
	}
}

func (m App) completed() (App, tea.Cmd) {
	m.status = completed
	return m, m.wrongTimer.Stop()
}

func (m App) tick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return intervalTick(t)
	})
}

func (m App) menuView() string {
	s := "    **** Match Game ****\n\n"
	s += "    press s to start the game\n"
	s += "    press h to go left\n"
	s += "    press j to go up\n"
	s += "    press k to go down\n"
	s += "    press l to go right\n"
	s += "    press q to exit\n"
	return s
}
