package cron

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hellgate75/go-cron/io"
	"github.com/hellgate75/go-cron/model"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type scheduler struct {
	sync.RWMutex
	cache        	map[string]model.CommandConfig
	commands     	[]model.CommandConfigRef
	cacheCommands   []model.CommandConfigRef
	runningTasks 	[]model.Execution
	syncRun      	bool
	running      	bool
	dir          	string
	file         	string
	enc          	io.Encoding
	errors		 	chan error
	warnings		chan error
	execMutex		sync.Mutex
	uuid			string
}

func (s *scheduler) IsRunning() bool {
	return s.running
}

func (s *scheduler) Errors() chan error {
	return s.errors
}

func (s *scheduler) Warnings() chan error {
	return s.warnings
}

func (s *scheduler) Start() error {
	if s.running {
		return errors.New("scheduler is already running")
	}
	go func(scheduler *scheduler) {
		for scheduler.running {
			_ = s.loadExecutions()
			if checkNextSchedulerTasks(scheduler) {
				err := executeSchedulerTasks(scheduler)
				if err != nil {
					scheduler.running=false
					scheduler.errors <- errors.New(fmt.Sprintf("Fatal tasks execution error: %v", err))
				}
			} else {
				fmt.Println("No active Tasks for the scheduler")
				time.Sleep(5 * time.Second)
			}
		}
	}(s)
	return nil
}

func (s *scheduler) Wait() {
main:
	for s.running {
		select {
		case <-time.After(15 * time.Second):
			if !s.running {
				break main
			}
		}
	}
}

func (s *scheduler) Stop() error {
	if !s.running {
		return errors.New("scheduler not already running")
	}
	s.running = false
	return nil
}

func (s *scheduler) RunOnce() error {
	_ = s.loadExecutions()
	if checkNextSchedulerTasks(s) {
		err := executeSchedulerTasks(s)
		if err != nil {
			s.running=false
			err = errors.New(fmt.Sprintf("Fatal tasks execution error: %v", err))
		}
		return err
	} else {
		return errors.New(fmt.Sprint("No active Tasks for the scheduler"))
	}
}

func (s *scheduler) Running() []model.Execution {
	return s.runningTasks
}

func (s *scheduler) Planned() []model.CommandConfig {
	var out = make([]model.CommandConfig, 0)
	out = append(out, s.cacheValues()...)
	for _, cfg := range s.commands {
		if !s.cacheContains(cfg.UUID) {
			cf, err := s.loadItem(cfg.UUID)
			if err == nil {
				out = append(out, *cf)
			}
		} else {
			out = append(out, *s.cacheValue(cfg.UUID))
		}
	}
	return out
}


func (s *scheduler) IsExecutionStored(id string) bool {
	return filterFirstExecution(s.runningTasks, func(m model.Execution) bool { return m.UUID == id }) != nil
}


func (s *scheduler) ToExecutionWith(ref model.CommandConfigRef, cmd model.CommandConfig) *model.Execution {
	if exec := filterFirstExecution(s.runningTasks, func(m model.Execution) bool { return m.UUID == ref.UUID }); exec != nil {
		exec.UpdateNext()
		return exec
	} else {
		exec := model.Execution{
			UUID:    ref.UUID,
			Command: cmd,
			Map: make(map[string]interface{}),
			Times: 0,
			Scheduled: false,
		}
		exec.Reset()

		return &exec
	}
}

func (s *scheduler) ToExecution(ref model.CommandConfigRef) *model.Execution {
	if exec := filterFirstExecution(s.runningTasks, func(m model.Execution) bool { return m.UUID == ref.UUID }); exec == nil {
		exec.UpdateNext()
		return exec
	} else {
		var item *model.CommandConfig
		var err error
		var id = ref.UUID
		if ! s.cacheContains(id) {
			item, err = s.loadItem(id)
		} else {
			item = s.cacheValue(id)
		}
		if err == nil {
			exec := model.Execution{
				UUID:    ref.UUID,
				Command: *item,
				Map: make(map[string]interface{}),
				Times: 0,
				Scheduled: false,
			}
			exec.Reset()
			return &exec
		}
	}
	return nil
}
func (s *scheduler) ExecutedTasks() []model.Execution {
	var out = make([]model.Execution, 0)
	for _, ref := range s.commands {
		exec := s.ToExecution(ref)
		if exec != nil {
			exec.UpdateNext()
			if time.Since(exec.Next).Milliseconds() < 0 {
				out = append(out, *exec)
			}
		}
	}
	return out
}

func (s *scheduler) NextRunningTasks() []model.Execution {
	var out = make([]model.Execution, 0)
	for _, ref := range s.commands {
		exec := s.ToExecution(ref)
		if exec != nil {
			exec.UpdateNext()
			if time.Since(exec.Next).Milliseconds() >= 0 {
				out = append(out, *exec)
			}
		}
	}
	return out
}

func (s *scheduler) AddToCache(cmd model.CommandConfig) error {
	var err error
	var id = uuid.New().String()
	var ref = model.CommandConfigRef{
		UUID:     id,
		Command:  cmd.Command,
		Created:  time.Now(),
		Updated:  time.Now(),
		FirstRun: time.Now(),
		LastRun:  time.Now(),
	}
	s.cache[id] = cmd
	if err != nil {
		return err
	}
	s.cacheCommands = append(s.cacheCommands, ref)
	return err
}

func (s *scheduler) AddAndPersist(cmd model.CommandConfig) error {
	var err error
	var id = uuid.New().String()
	var ref = model.CommandConfigRef{
		UUID:     id,
		Command:  cmd.Command,
		Created:  time.Now(),
		Updated:  time.Now(),
		FirstRun: time.Now(),
		LastRun:  time.Now(),
	}
	err = s.saveItem(id, cmd)
	if err != nil {
		return err
	}
	s.commands = append(s.commands, ref)
	err = s.save()
	return err
}

func (s *scheduler) UpdateToCache(cmd model.CommandConfig, index int) error {
	var err error
	if index >= 0 && index < len(s.cacheCommands) {
		s.cacheCommands[index].Command = cmd.Command
		s.cacheCommands[index].Updated = time.Now()
		s.cache[s.cacheCommands[index].UUID] = cmd
	} else {
		return errors.New(fmt.Sprintf("Index out of bound: %v, must be 0 <= x < %v ", index, len(s.commands)))
	}
	return err
}

func (s *scheduler) UpdateAndPersist(cmd model.CommandConfig, index int) error {
	var err error
	if index >= 0 && index < len(s.commands) {
		s.commands[index].Command = cmd.Command
		s.commands[index].Updated = time.Now()
		err = s.saveItem(s.commands[index].UUID, cmd)
		if err != nil {
			return err
		}
		err = s.save()
	} else {
		return errors.New(fmt.Sprintf("Index out of bound: %v, must be 0 <= x < %v ", index, len(s.commands)))
	}
	return err
}

func (s *scheduler) cacheValue(id string) *model.CommandConfig {
	if v, ok := s.cache[id]; ok {
		return &v
	}
	return nil
}


func (s *scheduler) cacheValues() []model.CommandConfig {
	var out = make([]model.CommandConfig, 0)
	for _, v := range s.cache{
		out = append(out, v)
	}
	return out
}


func (s *scheduler) cacheContains(id string) bool {
	if _, ok := s.cache[id]; ok {
		return true
	}
	return false
}

func (s *scheduler) DeleteFromCache(index int) error {
	var err error
	if index >= 0 && index < len(s.cacheCommands) {
		var length = len(s.cacheCommands)
		var tailEndIndex = length - 1
		var id = s.cacheCommands[index].UUID
		if s.cacheContains(id) {
			delete(s.cache, id)
			if index == 0 {
				// truncate array head
				s.cacheCommands = s.cacheCommands[1:]
			} else if index == tailEndIndex {
				// truncate array tail
				s.cacheCommands = s.cacheCommands[:tailEndIndex]
			} else {
				var chunk1 = s.cacheCommands[:index]
				var chunk2 = s.cacheCommands[index+1:]
				s.cacheCommands = chunk1
				s.cacheCommands = append(s.cacheCommands, chunk2...)
			}
			err = s.removeExecution(id)
		} else {
			err = errors.New(fmt.Sprintf("At index: %v, the UUID: %s has no value", index, id))
		}
	} else {
		return errors.New(fmt.Sprintf("Index out of bound: %v, must be 0 <= x < %v ", index, len(s.cacheCommands)))
	}
	return err
}

func (s *scheduler) DeleteAndPersist(index int) error {
	if index >= 0 && index < len(s.commands) {
		var err error
		var id = s.commands[index].UUID
		err = s.deleteItem(id)
		if err != nil {
			return err
		}
		var length = len(s.commands)
		if index == 0 {
			// truncate array head
			s.commands = s.commands[1:]
		} else if index == length-1 {
			// truncate array tail
			var tailEndIndex = length - 1
			s.commands = s.commands[:tailEndIndex]
		} else {
			var chunk1 = s.commands[:index]
			var chunk2 = s.commands[index+1:]
			s.commands = chunk1
			s.commands = append(s.commands, chunk2...)
		}
		err = s.save()
		if err != nil {
			return err
		}
		err = s.removeExecution(id)
		if err != nil {
			return err
		}
		err = s.saveExecutions()
		return err
	}
	return errors.New(fmt.Sprintf("Index out of bound: %v, must be 0 <= x < %v ", index, len(s.commands)))
}

// save the configuration to the file
func (s *scheduler) save() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		s.Unlock()
	}()
	s.Lock()
	err = io.SaveConfig(s.enc, s.file, model.SchedulerConfig{
		Sync:     s.syncRun,
		Commands: s.commands,
	})
	return err
}

func (s *scheduler) loadConfig() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		s.Unlock()
	}()
	s.Lock()
	var config = model.SchedulerConfig{}
	err = io.ReadConfig(s.enc, s.file, &config)
	if err == nil {
		s.syncRun = config.Sync
		s.commands = config.Commands
	}
	return err
}

// Load running tasks cache form his file
func (s *scheduler) loadExecutions() error {
	var err error
	var config =make([]model.Execution, 0)
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		s.execMutex.Unlock()
	}()
	s.execMutex.Lock()
	var file = fmt.Sprintf("%s%c%s.gob", s.dir, os.PathSeparator, "executions")
	err = io.ReadNative(file, &config)
	if err == nil {
		s.runningTasks = config
	}
	return err
}

// Save running tasks cache to his file
func (s *scheduler) saveExecutions() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		s.execMutex.Unlock()
	}()
	s.execMutex.Lock()
	var file = fmt.Sprintf("%s%c%s.gob", s.dir, os.PathSeparator, "executions")
	var config = make([]model.Execution, 0)
	for _, rt := range s.runningTasks {
		if ! s.cacheContains(rt.UUID) {
			config = append(config, rt)
		}
	}
	err = io.SaveNative(file, config)
	return err
}

// Save running tasks cache to his file
func (s *scheduler) trimExpiredExecutions() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		s.execMutex.Unlock()
	}()
	s.execMutex.Lock()
	var config = make([]model.Execution, 0)
	for _, rt := range s.runningTasks {
		if ! rt.Expired() {
			config = append(config, rt)
		}
	}
	s.runningTasks = config
	return err
}

// Save running tasks cache to his file
func (s *scheduler) removeExecution(id string) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		s.execMutex.Unlock()
	}()
	s.execMutex.Lock()
	var config = make([]model.Execution, 0)
	for _, rt := range s.runningTasks {
		if rt.UUID != id {
			config = append(config, rt)
		}
	}
	s.runningTasks = config
	return err
}

// Load a single command config form his file
func (s *scheduler) loadItem(id string) (*model.CommandConfig, error) {
	if _, ok := itemsLock[id]; !ok {
		itemsLock[s.uuid][id] = &sync.Mutex{}
	}
	var err error
	var config *model.CommandConfig
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		itemsLock[s.uuid][id].Unlock()
	}()
	itemsLock[s.uuid][id].Lock()
	var file = fmt.Sprintf("%s%c%s.gob", s.dir, os.PathSeparator, id)
	err = io.ReadNative(file, &config)
	return config, err
}

// save a single command config to his file
func (s *scheduler) saveItem(id string, config model.CommandConfig) error {
	var err error
	if _, ok := itemsLock[s.uuid][id]; !ok {
		itemsLock[s.uuid][id] = &sync.Mutex{}
	}
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		itemsLock[s.uuid][id].Unlock()
	}()
	itemsLock[s.uuid][id].Lock()
	var file = fmt.Sprintf("%s%c%s.gob", s.dir, os.PathSeparator, id)
	err = io.SaveNative(file, &config)
	return err
}

// save a single command config to his file
func (s *scheduler) deleteItem(id string) error {
	var err error
	if _, ok := itemsLock[s.uuid][id]; !ok {
		itemsLock[s.uuid][id] = &sync.Mutex{}
	}
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		itemsLock[s.uuid][id].Unlock()
	}()
	itemsLock[s.uuid][id].Lock()
	var file = fmt.Sprintf("%s%c%s.%s", s.dir, os.PathSeparator, id, s.enc.String())
	err = io.DeleteFile(file)
	return err
}

func (s *scheduler) Load() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		s.Unlock()
	}()
	s.Lock()
	var config = model.SchedulerConfig{}
	err = io.ReadConfig(s.enc, s.file, &config)
	if err == nil {
		s.syncRun = config.Sync
		s.commands = config.Commands
		err = s.loadExecutions()
	}
	return err
}

func (s *scheduler) References() []model.CommandConfigRef {
	return s.commands
}

func (s *scheduler) Destroy(saveState bool)  {
	delete(itemsLock, s.uuid)
	var err error
	if s.IsRunning() {
		err = s.Stop()
		if err != nil {
			s.errors <- err
		}
	}
	if saveState {
		err = s.save()
		if err != nil {
			s.errors <- err
		}
		err = s.trimExpiredExecutions()
		if err != nil {
			s.errors <- err
		}
		err = s.saveExecutions()
		if err != nil {
			s.errors <- err
		}
	}
	close(s.errors)
	close(s.warnings)
	s.runningTasks = make([]model.Execution, 0)
	s.cacheCommands = make([]model.CommandConfigRef, 0)
	s.commands = make([]model.CommandConfigRef, 0)
	s.cache  = make(map[string]model.CommandConfig)
}

// Load an existing scheduler, add the given scheduler config items and save the config file.
func LoadSchedulerWith(file string, encoding io.Encoding, commands []model.CommandConfig,
	syncRun bool) (model.Scheduler, []error) {
	var errorsList = make([]error, 0)
	var dir, _ = filepath.Split(file)
	if !io.FileExists(dir) {
		_ = io.CreateFolder(dir, 0777)
	}
	var sc = &scheduler{sync.RWMutex{}, make(map[string]model.CommandConfig), make([]model.CommandConfigRef, 0), make([]model.CommandConfigRef, 0), make([]model.Execution, 0), syncRun, false, dir, file, encoding, make(chan error), make(chan error), sync.Mutex{}, uuid.New().String()}
	itemsLock[sc.uuid] = make(map[string]*sync.Mutex)
	if file != "" {
		var err = sc.Load()
		if err != nil {
			errorsList = append(errorsList, err)
		}
	}
	for _, c := range commands {
		var err = sc.AddAndPersist(c)
		if err != nil {
			errorsList = append(errorsList, err)
		}
	}
	return sc, errorsList
}

// Load an existing scheduler and return the component.
func LoadSchedulerFrom(file string, encoding io.Encoding, syncRun bool) (model.Scheduler, error) {
	var err error
	var dir, _ = filepath.Split(file)
	if !io.FileExists(dir) {
		_ = io.CreateFolder(dir, 0777)
	}
	var sc = &scheduler{sync.RWMutex{}, make(map[string]model.CommandConfig), make([]model.CommandConfigRef, 0), make([]model.CommandConfigRef, 0), make([]model.Execution, 0), syncRun, false, dir, file, encoding, make(chan error), make(chan error), sync.Mutex{}, uuid.New().String()}
	itemsLock[sc.uuid] = make(map[string]*sync.Mutex)
	if file != "" {
		err = sc.Load()
	}
	return sc, err
}

// Create a new empty scheduler and save the config file.
func NewEmptyScheduler(file string, encoding io.Encoding, syncRun bool) (model.Scheduler, error) {
	var err error
	var dir, _ = filepath.Split(file)
	if !io.FileExists(dir) {
		_ = io.CreateFolder(dir, 0777)
	}
	var sc = &scheduler{sync.RWMutex{}, make(map[string]model.CommandConfig), make([]model.CommandConfigRef, 0), make([]model.CommandConfigRef, 0), make([]model.Execution, 0), syncRun, false, dir, file, encoding, make(chan error), make(chan error), sync.Mutex{}, uuid.New().String()}
	itemsLock[sc.uuid] = make(map[string]*sync.Mutex)
	if file != "" {
		err = sc.save()
	}
	return sc, err
}

// Create a new scheduler from given command config and save the config file.
func NewSchedulerWith(file string, encoding io.Encoding, commands []model.CommandConfig,
	syncRun bool) (model.Scheduler, []error) {
	var errorsList = make([]error, 0)
	var dir, _ = filepath.Split(file)
	if !io.FileExists(dir) {
		_ = io.CreateFolder(dir, 0777)
	}
	var sc = &scheduler{sync.RWMutex{}, make(map[string]model.CommandConfig), make([]model.CommandConfigRef, 0), make([]model.CommandConfigRef, 0), make([]model.Execution, 0), syncRun, false, dir, file, encoding, make(chan error), make(chan error), sync.Mutex{}, uuid.New().String()}
	itemsLock[sc.uuid] = make(map[string]*sync.Mutex)
	for _, c := range commands {
		var err = sc.AddAndPersist(c)
		if err != nil {
			errorsList = append(errorsList, err)
		}
	}
	return sc, errorsList
}
