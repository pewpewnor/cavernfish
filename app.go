package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"cavernfish/backend"

	"github.com/google/uuid"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx        context.Context
	store      *backend.CollectionStore
	srvManager *backend.ServerManager
	bpManager  *backend.BreakpointManager
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	home, _ := os.UserHomeDir()
	dataDir := filepath.Join(home, ".cavernfish")

	a.bpManager = backend.NewBreakpointManager()

	store, err := backend.NewCollectionStore(filepath.Join(dataDir, "collections"))
	if err != nil {
		wailsRuntime.LogError(ctx, fmt.Sprintf("store init failed: %v", err))
	}
	a.store = store

	cfgStore, err := backend.NewServerConfigStore(filepath.Join(dataDir, "servers.json"))
	if err != nil {
		wailsRuntime.LogError(ctx, fmt.Sprintf("server config store init failed: %v", err))
	}

	emitter := func(event string, data ...any) {
		wailsRuntime.EventsEmit(ctx, event, data...)
	}
	a.srvManager = backend.NewServerManager(store, cfgStore, a.bpManager, ctx, emitter)
}

func (a *App) GetCollections() []backend.Collection {
	return a.store.GetAll()
}

func (a *App) CreateCollection(name, description string) backend.Collection {
	return a.store.Create(name, description)
}

func (a *App) UpdateCollection(id string, upd backend.CollectionUpdate) (*backend.Collection, error) {
	c, ok := a.store.Update(id, func(c *backend.Collection) {
		if upd.Name != "" {
			c.Name = upd.Name
		}
		c.Description = upd.Description
		if upd.Variables != nil {
			c.Variables = upd.Variables
		}
		c.UpdatedAt = time.Now()
	})
	if !ok {
		return nil, fmt.Errorf("collection not found")
	}
	return c, nil
}

func (a *App) DeleteCollection(id string) {
	a.store.Delete(id)
}

func (a *App) ExportCollection(id string) (string, error) {
	return a.store.Export(id)
}

func (a *App) ImportCollection(data string) (*backend.Collection, error) {
	return a.store.Import(data)
}

func (a *App) ImportFromOpenAPI(data string) (*backend.Collection, error) {
	return backend.ImportFromOpenAPI(a.store, data)
}

func (a *App) CreateFolder(collectionID, name string) (*backend.Collection, error) {
	c, ok := a.store.Update(collectionID, func(c *backend.Collection) {
		c.Folders = append(c.Folders, backend.Folder{
			ID:        uuid.New().String(),
			Name:      name,
			Endpoints: []backend.Endpoint{},
		})
		c.UpdatedAt = time.Now()
	})
	if !ok {
		return nil, fmt.Errorf("collection not found")
	}
	return c, nil
}

func (a *App) RenameFolder(collectionID, folderID, name string) (*backend.Collection, error) {
	c, ok := a.store.Update(collectionID, func(c *backend.Collection) {
		for i := range c.Folders {
			if c.Folders[i].ID == folderID {
				c.Folders[i].Name = name
				break
			}
		}
		c.UpdatedAt = time.Now()
	})
	if !ok {
		return nil, fmt.Errorf("collection not found")
	}
	return c, nil
}

func (a *App) DeleteFolder(collectionID, folderID string) (*backend.Collection, error) {
	c, ok := a.store.Update(collectionID, func(c *backend.Collection) {
		for i := range c.Folders {
			if c.Folders[i].ID == folderID {
				c.Folders = append(c.Folders[:i], c.Folders[i+1:]...)
				break
			}
		}
		c.UpdatedAt = time.Now()
	})
	if !ok {
		return nil, fmt.Errorf("collection not found")
	}
	return c, nil
}

func (a *App) CreateEndpoint(collectionID, folderID string, ep backend.Endpoint) (*backend.Collection, error) {
	if ep.ID == "" {
		ep.ID = uuid.New().String()
	}
	if ep.Method == "" {
		ep.Method = "GET"
	}
	if ep.StatusCode == 0 {
		ep.StatusCode = 200
	}
	if ep.BodyType == "" {
		ep.BodyType = "json"
	}
	if ep.Body == "" {
		ep.Body = `{"message": "OK"}`
	}
	if ep.Headers == nil {
		ep.Headers = []backend.KVPair{}
	}
	if ep.Cookies == nil {
		ep.Cookies = []backend.MockCookie{}
	}
	c, ok := a.store.Update(collectionID, func(c *backend.Collection) {
		for i := range c.Folders {
			if c.Folders[i].ID == folderID {
				c.Folders[i].Endpoints = append(c.Folders[i].Endpoints, ep)
				break
			}
		}
		c.UpdatedAt = time.Now()
	})
	if !ok {
		return nil, fmt.Errorf("collection not found")
	}
	return c, nil
}

func (a *App) UpdateEndpoint(collectionID, folderID, endpointID string, ep backend.Endpoint) (*backend.Collection, error) {
	c, ok := a.store.Update(collectionID, func(c *backend.Collection) {
		for i := range c.Folders {
			if c.Folders[i].ID == folderID {
				for j := range c.Folders[i].Endpoints {
					if c.Folders[i].Endpoints[j].ID == endpointID {
						ep.ID = endpointID
						c.Folders[i].Endpoints[j] = ep
						break
					}
				}
				break
			}
		}
		c.UpdatedAt = time.Now()
	})
	if !ok {
		return nil, fmt.Errorf("collection not found")
	}
	return c, nil
}

func (a *App) DeleteEndpoint(collectionID, folderID, endpointID string) (*backend.Collection, error) {
	c, ok := a.store.Update(collectionID, func(c *backend.Collection) {
		for i := range c.Folders {
			if c.Folders[i].ID == folderID {
				for j := range c.Folders[i].Endpoints {
					if c.Folders[i].Endpoints[j].ID == endpointID {
						c.Folders[i].Endpoints = append(c.Folders[i].Endpoints[:j], c.Folders[i].Endpoints[j+1:]...)
						break
					}
				}
				break
			}
		}
		c.UpdatedAt = time.Now()
	})
	if !ok {
		return nil, fmt.Errorf("collection not found")
	}
	return c, nil
}

func (a *App) GetServers() []backend.ServerInfo {
	return a.srvManager.GetAll()
}

func (a *App) CreateServer(cfg backend.ServerConfig) backend.ServerInfo {
	if cfg.ID == "" {
		cfg.ID = uuid.New().String()
	}
	return a.srvManager.Create(cfg)
}

func (a *App) UpdateServer(cfg backend.ServerConfig) (backend.ServerInfo, error) {
	return a.srvManager.Update(cfg)
}

func (a *App) IsServerRunning(id string) bool {
	return a.srvManager.IsRunning(id)
}

func (a *App) StartServer(id string) error {
	return a.srvManager.Start(id)
}

func (a *App) StopServer(id string) error {
	return a.srvManager.Stop(id)
}

func (a *App) RestartServer(id string) error {
	if err := a.srvManager.Stop(id); err != nil {
		return err
	}
	return a.srvManager.Start(id)
}

func (a *App) DeleteServer(id string) {
	a.srvManager.Delete(id)
}

func (a *App) GetRequestLog(serverID string) []backend.RequestLogEntry {
	return a.srvManager.GetLog(serverID)
}

func (a *App) ClearRequestLog(serverID string) {
	a.srvManager.ClearLog(serverID)
}

func (a *App) GetPendingBreakpoints() []backend.BreakpointHit {
	return a.bpManager.GetPending()
}

func (a *App) ReleaseBreakpoint(id string, resp backend.MockResponse) bool {
	released := a.bpManager.Release(id, resp)
	if released {
		wailsRuntime.EventsEmit(a.ctx, "breakpoint:released", id)
	}
	return released
}

func (a *App) DiscardBreakpoint(id string) bool {
	discarded := a.bpManager.Discard(id)
	if discarded {
		wailsRuntime.EventsEmit(a.ctx, "breakpoint:released", id)
	}
	return discarded
}
