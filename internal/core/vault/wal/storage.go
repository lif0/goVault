package wal

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"goVault/internal/configuration"
	"goVault/internal/pkg/directory"
	"goVault/internal/pkg/unit"
)

var (
	walSegmentNamef = "%i.wal"

	defaultMaxSegmentSize   = 10_000 // 10MB
	defaultWALRootDirectory = "/data/wal"
)

type walSegment struct {
	mu                sync.Mutex
	segmentDescriptor *os.File
	segmentName       string
	segmentSize       int

	rootDirectory  string
	maxSegmentSize int
}

func NewWALSegment(cfg configuration.WALConfig) (*walSegment, error) {
	maxSegmentSize := defaultMaxSegmentSize
	rootDir := defaultWALRootDirectory

	if cfg.MaxSegmentSize != "" {
		ms, err := unit.ParseDigitalStorage(cfg.MaxSegmentSize)
		if err != nil {
			return nil, err
		}

		maxSegmentSize = ms
	}

	if cfg.DataDirectory != "" {
		err := directory.TryCreateDirsByPath(cfg.DataDirectory)
		if err != nil {
			return nil, fmt.Errorf("failed create dir: %s", cfg.DataDirectory)
		}

		rootDir = cfg.DataDirectory
	}

	service := &walSegment{
		mu:                sync.Mutex{},
		segmentName:       "",
		segmentSize:       0,
		segmentDescriptor: nil,

		rootDirectory:  rootDir,
		maxSegmentSize: maxSegmentSize,
	}

	sDescriptor, sName, err := createSegmentDescriptor()
	if err != nil {
		return nil, fmt.Errorf("failed create segment file(%s): %v", cfg.DataDirectory, err)
	}

	service.segmentDescriptor, service.segmentName = sDescriptor, sName

	return service, nil
}

func (ws *walSegment) Write(data []byte) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	err := ws.writeUnsafe(data)
	if err != nil {
		return err
	}
	ws.segmentSize += len(data)

	return nil
}

func (ws *walSegment) writeUnsafe(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	if _, err := ws.segmentDescriptor.Write(data); err != nil {
		return errors.Join(ErrWalStorageErrorWriteDataToDisk, err)
	}

	if err := ws.segmentDescriptor.Sync(); err != nil {
		return errors.Join(ErrWalStorageErrorWriteDataToDisk, err)
	}

	return nil
}

func getSegmentName() string {
	return fmt.Sprintf(walSegmentNamef, time.Now().UTC().UnixNano())
}

func (ws *walSegment) recreateSegmentDescriptorUnsafe() error {
	err := ws.segmentDescriptor.Close()
	if err != nil {
		return err
	}

	sDescriptor, sName, err := createSegmentDescriptor()
	if err != nil {
		return err
	}

	ws.segmentDescriptor = sDescriptor
	ws.segmentName = sName
	ws.segmentSize = 0

	return nil
}

func createSegmentDescriptor() (segmentDescriptor *os.File, segmentName string, err error) {
	name := getSegmentName()
	nsd, err := os.Create(name)
	if err != nil {
		return nil, "", err
	}

	return nsd, name, nil
}
