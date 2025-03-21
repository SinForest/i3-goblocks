package module

import (
	"bufio"
	"fmt"
	"iter"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Module struct {
	logger       *log.Logger
	basePath     string
	tick         int
	clickHandler func(*Module, Click)
	stdIn        chan []byte
	DebugMode    bool
}

func New(logName, basePath string, tick int) *Module {
	var logger *log.Logger
	if logName != "" {
		logger = log.New(os.Stderr, fmt.Sprintf("[i3-goblocks][%s] ", logName), log.LstdFlags)
	}
	return &Module{
		logger:       logger,
		basePath:     basePath,
		tick:         tick,
		clickHandler: nil,
	}
}

func (m *Module) Log(v ...any) {
	m.logger.Println(v...)
}

func (m *Module) ScanSysFile(relPath string) iter.Seq[string] {
	f, err := os.Open(path.Join(m.basePath, relPath))
	if err != nil {
		m.logger.Fatalf("failed to open sys file %q: %v", relPath, err)
	}
	scanner := bufio.NewScanner(f)
	return func(yield func(string) bool) {
		defer f.Close()
		for {
			if !scanner.Scan() {
				return
			}
			if !yield(scanner.Text()) {
				return
			}
		}
	}
}

func (m *Module) ReadSysFile(relPath string) (string, error) {
	s, err := os.ReadFile(path.Join(m.basePath, relPath))
	if err != nil {
		return "", fmt.Errorf("failed to read sys file %q: %w", relPath, err)
	}
	return strings.TrimSpace(string(s)), nil
}

func (m *Module) MustReadSysFile(relPath string) string {
	s, err := m.ReadSysFile(relPath)
	if err != nil {
		m.logger.Fatalf("%v", err)
	}
	return s
}

func (m *Module) ReadFloat(relPath string) (float64, error) {
	s, err := m.ReadSysFile(relPath)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert %q to float: %w", s, err)
	}
	return f, nil
}

func (m *Module) MustReadFloat(relPath string) float64 {
	f, err := m.ReadFloat(relPath)
	if err != nil {
		m.logger.Fatalf("%v", err)
	}
	return f
}

func (m *Module) WriteSysFile(relPath string, inp string) {
	f, err := os.Create(path.Join(m.basePath, relPath))
	if err != nil {
		m.logger.Fatalf("failed to open sys file %q: %v", relPath, err)
	}
	defer f.Close()
	_, err = f.Write([]byte(inp))
	if err != nil {
		m.logger.Fatalf("failed to write to sys file %q: %v", relPath, err)
	}
}

func (m *Module) RegisterClickHandler(f func(*Module, Click)) {
	m.clickHandler = f
}

func (m *Module) initStdReader() {
	m.stdIn = make(chan []byte)
	go func() {
		scn := bufio.NewScanner(os.Stdin)
		for scn.Scan() {
			m.stdIn <- scn.Bytes()
		}
		fmt.Printf("DEBUG: Scanning Stdin stopped for some reason!")
	}()
}

type Click byte

const (
	BtnLeft Click = '1' + iota
	BtnMiddle
	BtnRight
	BtnWheelUp
	BtnWheelDown
)

func (m *Module) handleClick(inp []byte) {
	m.clickHandler(m, Click(inp[0]))
}

func (m *Module) Run(f func() error) {
	tick := func() {
		err := f()
		if err != nil {
			m.logger.Fatal(err.Error())
		}
	}
	tick()
	if m.tick <= 0 {
		//TODO: handle clicks here, too
		return
	}
	ticker := time.NewTicker(time.Duration(m.tick) * time.Second)
	if m.clickHandler != nil {
		m.initStdReader()
	}
	for {
		select {
		//TODO: handle refreshing signals?
		case <-ticker.C:
			tick()
		case inp := <-m.stdIn:
			m.handleClick(inp)
			tick()
		}
	}

}
