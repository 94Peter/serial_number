package model

import (
	"encoding/gob"
	"errors"
	"os"
	"sync"
	"time"
)

var (
	Err_PrefixNotFound = errors.New("prefix_not_found")
	Err_PrefixExist    = errors.New("prefix_exist")
)

type SerialMgr interface {
	CreateSerial(prefix string, startNumber int64) error
	UpdateSerial(prefix string, startNumber int64) error
	GetSerial(prefix string) (int64, error)
	ClearSerial(prefix string) error
	Persistance() error
}

type serial struct {
	serialMap   map[string]int64
	persistFile string
	persistTime time.Time
	modifyTime  time.Time
	sync.Mutex
}

func NewSerial(cfg *Config) SerialMgr {
	// check if file exists
	mgr := &serial{
		persistFile: cfg.PersistanceFile,
		serialMap:   make(map[string]int64),
	}
	if FileExists(cfg.PersistanceFile) {
		// if file exists, load from file
		f, err := os.Open(cfg.PersistanceFile)
		if err != nil {
			return nil
		}
		defer f.Close()
		dec := gob.NewDecoder(f)
		dec.Decode(&mgr.serialMap)
	}
	now := time.Now()
	mgr.persistTime = now
	mgr.modifyTime = now
	return mgr
}

func (s *serial) CreateSerial(prefix string, startNumber int64) error {
	if _, ok := s.serialMap[prefix]; ok {
		return Err_PrefixExist
	}
	s.Lock()
	s.serialMap[prefix] = startNumber
	s.modifyTime = time.Now()
	s.Unlock()
	return nil
}

func (s *serial) UpdateSerial(prefix string, startNumber int64) error {
	if _, ok := s.serialMap[prefix]; !ok {
		return Err_PrefixNotFound
	}
	s.Lock()
	s.serialMap[prefix] = startNumber
	s.modifyTime = time.Now()
	s.Unlock()
	return nil
}

func (s *serial) GetSerial(prefix string) (int64, error) {
	if _, ok := s.serialMap[prefix]; !ok {
		return 0, Err_PrefixNotFound
	}
	s.Lock()
	defer s.Unlock()
	result := s.serialMap[prefix] + 1
	s.serialMap[prefix] = result
	s.modifyTime = time.Now()
	return result, nil
}

func (s *serial) ClearSerial(prefix string) error {
	if _, ok := s.serialMap[prefix]; !ok {
		return Err_PrefixNotFound
	}
	s.Lock()
	delete(s.serialMap, prefix)
	s.modifyTime = time.Now()
	s.Unlock()
	return nil
}

func (s *serial) Persistance() error {
	if s.persistTime.After(s.modifyTime) {
		return nil
	}
	// save map to file
	f, err := os.Create(s.persistFile)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)

	s.Lock()
	defer s.Unlock()
	if err := enc.Encode(s.serialMap); err != nil {
		return err
	}
	s.persistTime = time.Now()
	return nil
}
