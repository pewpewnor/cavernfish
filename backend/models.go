package backend

import (
	"encoding/json"
	"time"
)

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
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Method     string       `json:"method"`
	Path       string       `json:"path"`
	StatusCode int          `json:"statusCode"`
	Body       string       `json:"body"`
	BodyType   string       `json:"bodyType"`
	Headers    []KVPair     `json:"headers"`
	Cookies    []MockCookie `json:"cookies"`
	DelayMs    int          `json:"delayMs"`
	Breakpoint bool         `json:"breakpoint"`
	ProxyURL   string       `json:"proxyUrl,omitempty"`
	WSEnabled  bool         `json:"wsEnabled"`
}

type endpointAlias Endpoint

type endpointWire struct {
	endpointAlias
	Responses []MockResponse `json:"responses,omitempty"`
	ActiveIdx int            `json:"activeIdx,omitempty"`
	Strategy  string         `json:"strategy,omitempty"`
}

func (e *Endpoint) UnmarshalJSON(data []byte) error {
	var w endpointWire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	*e = Endpoint(w.endpointAlias)
	if e.StatusCode == 0 && len(w.Responses) > 0 {
		idx := w.ActiveIdx
		if idx < 0 || idx >= len(w.Responses) {
			idx = 0
		}
		r := w.Responses[idx]
		e.StatusCode = r.StatusCode
		e.Body = r.Body
		e.BodyType = r.BodyType
		e.Headers = r.Headers
		e.Cookies = r.Cookies
	}
	return nil
}

func (e *Endpoint) ToResponse() MockResponse {
	return MockResponse{
		StatusCode: e.StatusCode,
		Body:       e.Body,
		BodyType:   e.BodyType,
		Headers:    append([]KVPair(nil), e.Headers...),
		Cookies:    append([]MockCookie(nil), e.Cookies...),
	}
}

type MockResponse struct {
	StatusCode int          `json:"statusCode"`
	Body       string       `json:"body"`
	BodyType   string       `json:"bodyType"`
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
	SameSite string `json:"sameSite"`
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
