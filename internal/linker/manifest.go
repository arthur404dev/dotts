package linker

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Manifest struct {
	path    string
	entries map[string]LinkEntry
	mu      sync.RWMutex
}

func LoadManifest(dataDir string) (*Manifest, error) {
	path := filepath.Join(dataDir, "manifest.json")
	m := &Manifest{
		path:    path,
		entries: make(map[string]LinkEntry),
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return m, nil
	}
	if err != nil {
		return nil, err
	}

	var entries []LinkEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	for _, entry := range entries {
		m.entries[entry.Target] = entry
	}

	return m, nil
}

func (m *Manifest) Add(entry LinkEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[entry.Target] = entry
}

func (m *Manifest) Remove(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, target)
}

func (m *Manifest) HasEntry(target string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.entries[target]
	return ok
}

func (m *Manifest) Get(target string) (LinkEntry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	entry, ok := m.entries[target]
	return entry, ok
}

func (m *Manifest) Entries() []LinkEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entries := make([]LinkEntry, 0, len(m.entries))
	for _, entry := range m.entries {
		entries = append(entries, entry)
	}
	return entries
}

func (m *Manifest) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.entries)
}

func (m *Manifest) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entries := make([]LinkEntry, 0, len(m.entries))
	for _, entry := range m.entries {
		entries = append(entries, entry)
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(m.path), 0755); err != nil {
		return err
	}

	return os.WriteFile(m.path, data, 0644)
}

func (m *Manifest) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = make(map[string]LinkEntry)
}
