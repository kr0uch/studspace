package logger

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kr0uch/studspace/pkg/utils"
)

type dailyFile struct {
	dir  string
	mu   sync.Mutex
	day  string
	file *os.File
}

func newDailyFile(dir string) *dailyFile {
	return &dailyFile{dir: dir}
}

func (w *dailyFile) Write(p []byte) (int, error) {
	day := time.Now().In(utils.Moscow).Format("2006-01-02")

	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.rotate(day); err != nil {
		return 0, err
	}
	return w.file.Write(p)
}

func (w *dailyFile) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file == nil {
		return nil
	}
	return w.file.Sync()
}

func (w *dailyFile) rotate(day string) error {
	if w.file != nil && w.day == day {
		return nil
	}
	if w.file != nil {
		_ = w.file.Close()
		w.file = nil
	}

	file, err := os.OpenFile(filepath.Join(w.dir, day+".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	w.file = file
	w.day = day
	return nil
}
