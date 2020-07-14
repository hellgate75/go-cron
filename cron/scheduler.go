package cron

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-cron/io"
	"github.com/hellgate75/go-cron/model"
	"path/filepath"
)

type scheduler struct{
	commands			[]model.CommandConfigRef
	syncRun				bool
	running				bool
	dir 				string
	file 				string
	enc 				io.Encoding

}

func (s *scheduler) IsRunning() bool {
	return s.running
}

func (s *scheduler) Start() error {
	panic("implement me")
}

func (s *scheduler) Stop() error {
	panic("implement me")
}

func (scheduler) RunOnce() error {
	panic("implement me")
}

func (s *scheduler) Running() []model.CommandValue {
	panic("implement me")
}

func (s *scheduler) Planned() []model.CommandConfig {
	panic("implement me")
}

func (s *scheduler) NextRunningTasks() []model.Execution {
	panic("implement me")
}

func (s *scheduler) AddAndSave(cmd model.CommandConfig, file string, enc io.Encoding) error {
	s.commands = append(s.commands, cmd)
	return s.save(file, enc)
}

func (s *scheduler) UpdateAndSave(cmd model.CommandConfig, index int, file string, enc io.Encoding) error {
	if index >= 0 && index < len(s.commands) {
		s.commands[index] = cmd
		return s.save(file, enc)
	}
	return errors.New(fmt.Sprintf("Index out of bound: %v, must be 0 <= x < %v ", index, len(s.commands)))
}

func (s *scheduler) DeleteAndSave(index int, file string, enc io.Encoding) error {
	if index >= 0 && index < len(s.commands) {
		var length = len(s.commands)
		if index == 0 {
			// truncate array head
			s.commands = s.commands[1:]
		} else if index == length - 1 {
			// truncate array tail
			var tailEndIndex = length - 1
			s.commands = s.commands[:tailEndIndex]
		} else {
			var chunk1 = s.commands[:index]
			var chunk2 = s.commands[index+1:]
			s.commands = chunk1
			s.commands = append(s.commands, chunk2...)
		}
		return s.save(file, enc)
	}
	return errors.New(fmt.Sprintf("Index out of bound: %v, must be 0 <= x < %v ", index, len(s.commands)))
}

func (s *scheduler) save() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	return io.SaveConfig(s.enc, s.file, model.SchedulerConfig{
		Sync: s.syncRun,
		Commands: s.commands,
	})
}

func (s *scheduler) LoadFrom(path string) {

}

func (s *scheduler) References() []model.CommandConfigRef {
	var out = make([]model.CommandConfigRef, 0)
	return out
}

func NewSchedulerFrom(file string, encoding io.Encoding, commands []model.CommandConfig,
	syncRun	bool) model.Scheduler {
	var refs = make([]model.CommandConfigRef, 0)
	var dir, _ = filepath.Split(file)
	if ! io.FileExists(dir) {
		_ = io.CreateFolder(dir, 0777)
	}
	var sc = &scheduler{refs, syncRun, false, dir, file, enc}
	return sc
}

