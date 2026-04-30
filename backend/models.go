package backend

import "time"

type CollectionUpdate struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Variables   map[string]string `json:"variables"`
}

type Collection struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Variables   map[string]string `json:"variables"`
	Folders     []Folder          `json:"folders"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

type Folder struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Method     string         `json:"method"`   // GET POST PUT PATCH DELETE HEAD OPTIONS *
	Path       string         `json:"path"`     // /api/users/:id or /api/*
	Responses  []MockResponse `json:"responses"`
	Strategy   string         `json:"strategy"` // fixed | cycle | random
	ActiveIdx  int            `json:"activeIdx"`
	DelayMs    int            `json:"delayMs"`
	Breakpoint bool           `json:"breakpoint"`
	ProxyURL   string         `json:"proxyUrl,omitempty"`
	WSEnabled  bool           `json:"wsEnabled"`
}

type MockResponse struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	StatusCode int          `json:"statusCode"`
	Body       string       `json:"body"`
	BodyType   string       `json:"bodyType"` // json | text | html | xml | none
	Headers    []KVPair     `json:"headers"`
	Cookies    []MockCookie `json:"cookies"`
}

type KVPair struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

type MockCookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Path     string `json:"path"`
	Domain   string `json:"domain"`
	MaxAge   int    `json:"maxAge"`
	HttpOnly bool   `json:"httpOnly"`
	Secure   bool   `json:"secure"`
	SameSite string `json:"sameSite"` // Strict | Lax | None
}

type ServerConfig struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Port                int      `json:"port"`
	CollectionIDs       []string `json:"collectionIds"`
	HTTPS               bool     `json:"https"`
	DisabledEndpointIDs []string `json:"disabledEndpointIds"`
}

type ServerInfo struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Port                int      `json:"port"`
	CollectionIDs       []string `json:"collectionIds"`
	HTTPS               bool     `json:"https"`
	Status              string   `json:"status"`
	ErrorMsg            string   `json:"errorMsg,omitempty"`
	RequestCount        int      `json:"requestCount"`
	DisabledEndpointIDs []string `json:"disabledEndpointIds"`
}

type RequestLogEntry struct {
	ID           string            `json:"id"`
	ServerID     string            `json:"serverId"`
	Timestamp    time.Time         `json:"timestamp"`
	Method       string            `json:"method"`
	Path         string            `json:"path"`
	Query        string            `json:"query"`
	ReqHeaders   map[string]string `json:"reqHeaders"`
	ReqBody      string            `json:"reqBody"`
	StatusCode   int               `json:"statusCode"`
	LatencyMs    int64             `json:"latencyMs"`
	Breakpointed bool              `json:"breakpointed"`
	BreakpointID string            `json:"breakpointId,omitempty"`
	IsProxy      bool              `json:"isProxy"`
	IsWS         bool              `json:"isWs"`
	Matched      bool              `json:"matched"`
}

type BreakpointHit struct {
	ID         string            `json:"id"`
	ServerID   string            `json:"serverId"`
	Timestamp  time.Time         `json:"timestamp"`
	Method     string            `json:"method"`
	Path       string            `json:"path"`
	Query      string            `json:"query"`
	ReqHeaders map[string]string `json:"reqHeaders"`
	ReqBody    string            `json:"reqBody"`
	Endpoint   Endpoint          `json:"endpoint"`
	Response   MockResponse      `json:"response"`
}
