package module

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

type Module struct {
	Logger   *log.Logger
	BasePath string
}

func New(logName, basePath string) *Module {
	return &Module{
		Logger:   log.New(os.Stderr, fmt.Sprintf("[i3-goblocks][%s] ", logName), log.LstdFlags),
		BasePath: basePath,
	}
}

func (m *Module) ReadSysFile(relPath string) string {
	s, err := os.ReadFile(path.Join(m.BasePath, relPath))
	if err != nil {
		m.Logger.Fatalf("failed to read sys file %q: %v", relPath, err)
	}
	return strings.TrimSpace(string(s))
}

func (m *Module) ReadFloat(relPath string) float64 {
	s := m.ReadSysFile(relPath)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		m.Logger.Fatalf("failed to convert %q to float: %v", s, err)
	}
	return f
}
