package backend

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// CollectionStore persists collections as individual JSON files under dataDir.
type CollectionStore struct {
	mu          sync.RWMutex
	dataDir     string
	collections map[string]*Collection
}

func NewCollectionStore(dataDir string) (*CollectionStore, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	cs := &CollectionStore{
		dataDir:     dataDir,
		collections: make(map[string]*Collection),
	}
	return cs, cs.load()
}

func (cs *CollectionStore) load() error {
	files, err := filepath.Glob(filepath.Join(cs.dataDir, "*.json"))
	if err != nil {
		return err
	}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		var c Collection
		if err := json.Unmarshal(data, &c); err != nil {
			continue
		}
		cs.collections[c.ID] = &c
	}
	return nil
}

func (cs *CollectionStore) save(c *Collection) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(cs.dataDir, c.ID+".json"), data, 0644)
}

func (cs *CollectionStore) GetAll() []Collection {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	result := make([]Collection, 0, len(cs.collections))
	for _, c := range cs.collections {
		result = append(result, *c)
	}
	return result
}

func (cs *CollectionStore) Get(id string) (*Collection, bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	c, ok := cs.collections[id]
	if !ok {
		return nil, false
	}
	cp := *c
	return &cp, true
}

func (cs *CollectionStore) Create(name, description string) Collection {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	c := Collection{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Variables:   make(map[string]string),
		Folders:     []Folder{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	cs.collections[c.ID] = &c
	_ = cs.save(&c)
	return c
}

func (cs *CollectionStore) Update(id string, updateFn func(*Collection)) (*Collection, bool) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	c, ok := cs.collections[id]
	if !ok {
		return nil, false
	}
	updateFn(c)
	_ = cs.save(c)
	cp := *c
	return &cp, true
}

func (cs *CollectionStore) Delete(id string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	delete(cs.collections, id)
	_ = os.Remove(filepath.Join(cs.dataDir, id+".json"))
}

func (cs *CollectionStore) Export(id string) (string, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	c, ok := cs.collections[id]
	if !ok {
		return "", fmt.Errorf("collection %s not found", id)
	}
	data, err := json.MarshalIndent(c, "", "  ")
	return string(data), err
}

func (cs *CollectionStore) Import(data string) (*Collection, error) {
	var c Collection
	if err := json.Unmarshal([]byte(data), &c); err != nil {
		return nil, err
	}
	c.ID = uuid.New().String()
	c.UpdatedAt = time.Now()
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.collections[c.ID] = &c
	_ = cs.save(&c)
	return &c, nil
}

// ServerConfigStore persists server configs in a single JSON file.
type ServerConfigStore struct {
	mu      sync.RWMutex
	path    string
	configs map[string]ServerConfig
}

func NewServerConfigStore(path string) (*ServerConfigStore, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	s := &ServerConfigStore{path: path, configs: make(map[string]ServerConfig)}
	return s, s.load()
}

func (s *ServerConfigStore) load() error {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	var configs []ServerConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return err
	}
	for _, c := range configs {
		s.configs[c.ID] = c
	}
	return nil
}

func (s *ServerConfigStore) save() error {
	configs := make([]ServerConfig, 0, len(s.configs))
	for _, c := range s.configs {
		configs = append(configs, c)
	}
	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *ServerConfigStore) GetAll() []ServerConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]ServerConfig, 0, len(s.configs))
	for _, c := range s.configs {
		result = append(result, c)
	}
	return result
}

func (s *ServerConfigStore) Set(cfg ServerConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.configs[cfg.ID] = cfg
	_ = s.save()
}

func (s *ServerConfigStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.configs, id)
	_ = s.save()
}
