package backend

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const maxLogEntries = 500

type runningServer struct {
	config       ServerConfig
	httpServer   *http.Server
	status       string
	errorMsg     string
	requestCount int
	log          []RequestLogEntry
	logMu        sync.Mutex
	cancelFunc   context.CancelFunc
}

type ServerManager struct {
	mu            sync.RWMutex
	servers       map[string]*runningServer
	store         *CollectionStore
	cfgStore      *ServerConfigStore
	bpManager     *BreakpointManager
	appCtx        context.Context
	eventEmitter  func(string, ...interface{})
	cycleMu       sync.Mutex
	cycleCounters map[string]int // endpointID → next index
}

func NewServerManager(store *CollectionStore, cfgStore *ServerConfigStore, bpManager *BreakpointManager, appCtx context.Context, emitter func(string, ...interface{})) *ServerManager {
	sm := &ServerManager{
		servers:       make(map[string]*runningServer),
		store:         store,
		cfgStore:      cfgStore,
		bpManager:     bpManager,
		appCtx:        appCtx,
		eventEmitter:  emitter,
		cycleCounters: make(map[string]int),
	}
	// Load persisted server configs
	for _, cfg := range cfgStore.GetAll() {
		sm.servers[cfg.ID] = &runningServer{
			config: cfg,
			status: "stopped",
			log:    []RequestLogEntry{},
		}
	}
	return sm
}

func (sm *ServerManager) Create(cfg ServerConfig) ServerInfo {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if cfg.ID == "" {
		cfg.ID = uuid.New().String()
	}
	rs := &runningServer{config: cfg, status: "stopped", log: []RequestLogEntry{}}
	sm.servers[cfg.ID] = rs
	sm.cfgStore.Set(cfg)
	return sm.toInfo(rs)
}

func (sm *ServerManager) Update(cfg ServerConfig) (ServerInfo, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	rs, ok := sm.servers[cfg.ID]
	if !ok {
		return ServerInfo{}, fmt.Errorf("server not found")
	}
	if rs.status == "running" {
		return ServerInfo{}, fmt.Errorf("stop the server before editing")
	}
	rs.config = cfg
	sm.cfgStore.Set(cfg)
	return sm.toInfo(rs), nil
}

func (sm *ServerManager) Start(id string) error {
	sm.mu.Lock()
	rs, ok := sm.servers[id]
	sm.mu.Unlock()
	if !ok {
		return fmt.Errorf("server %s not found", id)
	}
	if rs.status == "running" {
		return nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sm.handleRequest(w, r, id)
	})

	addr := fmt.Sprintf(":%d", rs.config.Port)
	httpSrv := &http.Server{Addr: addr, Handler: mux}

	if rs.config.HTTPS {
		cert, err := generateSelfSignedCert()
		if err != nil {
			sm.mu.Lock()
			rs.status = "error"
			rs.errorMsg = err.Error()
			sm.mu.Unlock()
			sm.eventEmitter("server:status", sm.GetAll())
			return err
		}
		httpSrv.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	}

	_, cancel := context.WithCancel(sm.appCtx)

	sm.mu.Lock()
	rs.httpServer = httpSrv
	rs.cancelFunc = cancel
	rs.status = "running"
	rs.errorMsg = ""
	sm.mu.Unlock()

	go func() {
		var err error
		if rs.config.HTTPS {
			err = httpSrv.ListenAndServeTLS("", "")
		} else {
			err = httpSrv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			sm.mu.Lock()
			rs.status = "error"
			rs.errorMsg = err.Error()
			sm.mu.Unlock()
			sm.eventEmitter("server:status", sm.GetAll())
		}
	}()

	sm.eventEmitter("server:status", sm.GetAll())
	return nil
}

func (sm *ServerManager) Stop(id string) error {
	sm.mu.RLock()
	rs, ok := sm.servers[id]
	sm.mu.RUnlock()
	if !ok {
		return fmt.Errorf("server %s not found", id)
	}

	if rs.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = rs.httpServer.Shutdown(ctx)
	}
	if rs.cancelFunc != nil {
		rs.cancelFunc()
	}

	sm.mu.Lock()
	rs.status = "stopped"
	rs.httpServer = nil
	sm.mu.Unlock()

	sm.eventEmitter("server:status", sm.GetAll())
	return nil
}

func (sm *ServerManager) Delete(id string) {
	_ = sm.Stop(id)
	sm.mu.Lock()
	delete(sm.servers, id)
	sm.mu.Unlock()
	sm.cfgStore.Delete(id)
	sm.eventEmitter("server:status", sm.GetAll())
}

func (sm *ServerManager) GetAll() []ServerInfo {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	result := make([]ServerInfo, 0, len(sm.servers))
	for _, rs := range sm.servers {
		result = append(result, sm.toInfo(rs))
	}
	return result
}

func (sm *ServerManager) GetLog(serverID string) []RequestLogEntry {
	sm.mu.RLock()
	rs, ok := sm.servers[serverID]
	sm.mu.RUnlock()
	if !ok {
		return nil
	}
	rs.logMu.Lock()
	defer rs.logMu.Unlock()
	cp := make([]RequestLogEntry, len(rs.log))
	copy(cp, rs.log)
	return cp
}

func (sm *ServerManager) ClearLog(serverID string) {
	sm.mu.RLock()
	rs, ok := sm.servers[serverID]
	sm.mu.RUnlock()
	if !ok {
		return
	}
	rs.logMu.Lock()
	rs.log = nil
	rs.requestCount = 0
	rs.logMu.Unlock()
	sm.eventEmitter("server:status", sm.GetAll())
}

func (sm *ServerManager) toInfo(rs *runningServer) ServerInfo {
	rs.logMu.Lock()
	cnt := rs.requestCount
	rs.logMu.Unlock()
	return ServerInfo{
		ID:           rs.config.ID,
		Name:         rs.config.Name,
		Port:         rs.config.Port,
		CollectionID: rs.config.CollectionID,
		HTTPS:        rs.config.HTTPS,
		Status:       rs.status,
		ErrorMsg:     rs.errorMsg,
		RequestCount: cnt,
	}
}

var wsUpgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func (sm *ServerManager) handleRequest(w http.ResponseWriter, r *http.Request, serverID string) {
	start := time.Now()

	sm.mu.RLock()
	rs, ok := sm.servers[serverID]
	sm.mu.RUnlock()
	if !ok {
		http.Error(w, "server gone", 500)
		return
	}

	reqHeaders := make(map[string]string)
	for k, v := range r.Header {
		reqHeaders[k] = strings.Join(v, ", ")
	}

	var reqBody string
	if r.Body != nil {
		b, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		reqBody = string(b)
		r.Body = io.NopCloser(strings.NewReader(reqBody))
	}

	entry := RequestLogEntry{
		ID:         uuid.New().String(),
		ServerID:   serverID,
		Timestamp:  time.Now(),
		Method:     r.Method,
		Path:       r.URL.Path,
		Query:      r.URL.RawQuery,
		ReqHeaders: reqHeaders,
		ReqBody:    reqBody,
	}

	_, ep := sm.findEndpoint(rs.config.CollectionID, r.Method, r.URL.Path)
	entry.Matched = ep != nil

	if ep == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		fmt.Fprint(w, `{"error":"no matching endpoint","path":"`+r.URL.Path+`"}`)
		entry.StatusCode = 404
		entry.LatencyMs = time.Since(start).Milliseconds()
		sm.appendLog(rs, entry)
		sm.eventEmitter("request:logged", entry)
		return
	}

	if ep.DelayMs > 0 {
		time.Sleep(time.Duration(ep.DelayMs) * time.Millisecond)
	}

	if ep.WSEnabled && websocket.IsWebSocketUpgrade(r) {
		sm.handleWebSocket(w, r, rs, ep, &entry, start)
		return
	}

	if ep.ProxyURL != "" {
		sm.handleProxy(w, r, ep.ProxyURL, rs, &entry, start)
		return
	}

	resp := sm.selectResponse(ep)
	if resp == nil {
		w.WriteHeader(200)
		entry.StatusCode = 200
		entry.LatencyMs = time.Since(start).Milliseconds()
		sm.appendLog(rs, entry)
		sm.eventEmitter("request:logged", entry)
		return
	}

	if ep.Breakpoint {
		hit := BreakpointHit{
			ID:         uuid.New().String(),
			ServerID:   serverID,
			Timestamp:  time.Now(),
			Method:     r.Method,
			Path:       r.URL.Path,
			Query:      r.URL.RawQuery,
			ReqHeaders: reqHeaders,
			ReqBody:    reqBody,
			Endpoint:   *ep,
			Response:   *resp,
		}
		entry.Breakpointed = true
		entry.BreakpointID = hit.ID

		released, modifiedResp := sm.bpManager.Hold(hit, func(h BreakpointHit) {
			sm.eventEmitter("breakpoint:hit", h)
		})

		if !released {
			w.WriteHeader(200)
			entry.StatusCode = 200
			entry.LatencyMs = time.Since(start).Milliseconds()
			sm.appendLog(rs, entry)
			sm.eventEmitter("request:logged", entry)
			return
		}
		resp = modifiedResp
	}

	sm.writeResponse(w, resp)
	entry.StatusCode = resp.StatusCode
	entry.LatencyMs = time.Since(start).Milliseconds()
	sm.appendLog(rs, entry)
	sm.eventEmitter("request:logged", entry)
}

func (sm *ServerManager) writeResponse(w http.ResponseWriter, resp *MockResponse) {
	for _, h := range resp.Headers {
		if h.Enabled && h.Key != "" {
			w.Header().Set(h.Key, h.Value)
		}
	}
	for _, c := range resp.Cookies {
		http.SetCookie(w, &http.Cookie{
			Name:     c.Name,
			Value:    c.Value,
			Path:     c.Path,
			Domain:   c.Domain,
			MaxAge:   c.MaxAge,
			HttpOnly: c.HttpOnly,
			Secure:   c.Secure,
		})
	}
	if w.Header().Get("Content-Type") == "" {
		switch resp.BodyType {
		case "json":
			w.Header().Set("Content-Type", "application/json")
		case "html":
			w.Header().Set("Content-Type", "text/html")
		case "xml":
			w.Header().Set("Content-Type", "application/xml")
		default:
			w.Header().Set("Content-Type", "text/plain")
		}
	}
	w.WriteHeader(resp.StatusCode)
	fmt.Fprint(w, resp.Body)
}

func (sm *ServerManager) handleProxy(w http.ResponseWriter, r *http.Request, proxyURL string, rs *runningServer, entry *RequestLogEntry, start time.Time) {
	entry.IsProxy = true
	target, err := url.Parse(proxyURL)
	if err != nil {
		http.Error(w, "invalid proxy URL", 500)
		entry.StatusCode = 500
		entry.LatencyMs = time.Since(start).Milliseconds()
		sm.appendLog(rs, *entry)
		sm.eventEmitter("request:logged", *entry)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	r.URL.Host = target.Host
	r.URL.Scheme = target.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = target.Host
	rw := &statusCapture{ResponseWriter: w, status: 200}
	proxy.ServeHTTP(rw, r)
	entry.StatusCode = rw.status
	entry.LatencyMs = time.Since(start).Milliseconds()
	sm.appendLog(rs, *entry)
	sm.eventEmitter("request:logged", *entry)
}

type statusCapture struct {
	http.ResponseWriter
	status int
}

func (sc *statusCapture) WriteHeader(status int) {
	sc.status = status
	sc.ResponseWriter.WriteHeader(status)
}

func (sm *ServerManager) handleWebSocket(w http.ResponseWriter, r *http.Request, rs *runningServer, ep *Endpoint, entry *RequestLogEntry, start time.Time) {
	entry.IsWS = true
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		entry.StatusCode = 400
		entry.LatencyMs = time.Since(start).Milliseconds()
		sm.appendLog(rs, *entry)
		sm.eventEmitter("request:logged", *entry)
		return
	}
	defer conn.Close()

	resp := sm.selectResponse(ep)
	if resp != nil && resp.Body != "" {
		_ = conn.WriteMessage(websocket.TextMessage, []byte(resp.Body))
	}

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	entry.StatusCode = 101
	entry.LatencyMs = time.Since(start).Milliseconds()
	sm.appendLog(rs, *entry)
	sm.eventEmitter("request:logged", *entry)
}

func (sm *ServerManager) findEndpoint(collectionID, method, path string) (*Collection, *Endpoint) {
	coll, ok := sm.store.Get(collectionID)
	if !ok {
		return nil, nil
	}
	for i := range coll.Folders {
		for j := range coll.Folders[i].Endpoints {
			ep := &coll.Folders[i].Endpoints[j]
			if strings.EqualFold(ep.Method, method) || ep.Method == "*" {
				if matched, _ := matchPath(ep.Path, path); matched {
					return coll, ep
				}
			}
		}
	}
	return coll, nil
}

func (sm *ServerManager) selectResponse(ep *Endpoint) *MockResponse {
	if len(ep.Responses) == 0 {
		return nil
	}
	switch ep.Strategy {
	case "cycle":
		sm.cycleMu.Lock()
		idx := sm.cycleCounters[ep.ID]
		sm.cycleCounters[ep.ID] = (idx + 1) % len(ep.Responses)
		sm.cycleMu.Unlock()
		return &ep.Responses[idx]
	case "random":
		return &ep.Responses[rand.Intn(len(ep.Responses))]
	default: // fixed
		idx := ep.ActiveIdx
		if idx >= len(ep.Responses) {
			idx = 0
		}
		return &ep.Responses[idx]
	}
}

func (sm *ServerManager) appendLog(rs *runningServer, entry RequestLogEntry) {
	rs.logMu.Lock()
	rs.requestCount++
	rs.log = append(rs.log, entry)
	if len(rs.log) > maxLogEntries {
		rs.log = rs.log[len(rs.log)-maxLogEntries:]
	}
	rs.logMu.Unlock()
}

func matchPath(pattern, path string) (bool, map[string]string) {
	if pattern == path {
		return true, nil
	}
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	// Trailing wildcard
	if len(patternParts) > 0 && patternParts[len(patternParts)-1] == "*" {
		if len(pathParts) < len(patternParts)-1 {
			return false, nil
		}
		patternParts = patternParts[:len(patternParts)-1]
		if len(pathParts) > len(patternParts) {
			pathParts = pathParts[:len(patternParts)]
		}
	}

	if len(patternParts) != len(pathParts) {
		return false, nil
	}

	params := make(map[string]string)
	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			params[part[1:]] = pathParts[i]
		} else if part != pathParts[i] {
			return false, nil
		}
	}
	return true, params
}
