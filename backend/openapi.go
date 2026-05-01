package backend

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func ImportFromOpenAPI(store *CollectionStore, data string) (*Collection, error) {
	var spec map[string]any
	if err := json.Unmarshal([]byte(data), &spec); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}

	title := "Imported Collection"
	description := ""
	if info, ok := spec["info"].(map[string]any); ok {
		if t, ok := info["title"].(string); ok {
			title = t
		}
		if d, ok := info["description"].(string); ok {
			description = d
		}
	}

	c := store.Create(title, description)

	paths, ok := spec["paths"].(map[string]any)
	if !ok {
		return &c, nil
	}

	folderMap := make(map[string]string)

	for path, pathItem := range paths {
		pathItemMap, ok := pathItem.(map[string]any)
		if !ok {
			continue
		}
		for rawMethod, operation := range pathItemMap {
			method := strings.ToUpper(rawMethod)
			if method == "PARAMETERS" || method == "SERVERS" || method == "SUMMARY" || method == "DESCRIPTION" {
				continue
			}
			opMap, ok := operation.(map[string]any)
			if !ok {
				continue
			}

			tag := "Default"
			if tags, ok := opMap["tags"].([]any); ok && len(tags) > 0 {
				if t, ok := tags[0].(string); ok {
					tag = t
				}
			}

			if _, exists := folderMap[tag]; !exists {
				updated, err := createFolderInCollection(store, c.ID, tag)
				if err == nil && updated != nil {
					c = *updated
					for _, f := range c.Folders {
						if f.Name == tag {
							folderMap[tag] = f.ID
							break
						}
					}
				}
			}
			folderID := folderMap[tag]

			ep := Endpoint{
				ID:         uuid.New().String(),
				Method:     method,
				Path:       path,
				StatusCode: 200,
				BodyType:   "json",
				Body:       "{}",
				Headers:    []KVPair{},
				Cookies:    []MockCookie{},
			}

			if opID, ok := opMap["operationId"].(string); ok {
				ep.Name = opID
			} else if summary, ok := opMap["summary"].(string); ok {
				ep.Name = summary
			} else {
				ep.Name = method + " " + path
			}

			if responses, ok := opMap["responses"].(map[string]any); ok {
				bestStatus, bestBody := pickResponse(responses)
				ep.StatusCode = bestStatus
				if bestBody != "" {
					ep.Body = bestBody
				}
			}

			updated, err := addEndpointToFolder(store, c.ID, folderID, ep)
			if err == nil && updated != nil {
				c = *updated
			}
		}
	}

	return &c, nil
}

func pickResponse(responses map[string]any) (int, string) {
	preferred := []string{"200", "201", "204"}
	for _, p := range preferred {
		if r, ok := responses[p]; ok {
			return parseStatus(p), extractExample(r)
		}
	}
	for k, r := range responses {
		if strings.HasPrefix(k, "2") {
			return parseStatus(k), extractExample(r)
		}
	}
	for k, r := range responses {
		return parseStatus(k), extractExample(r)
	}
	return 200, ""
}

func parseStatus(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 200
	}
	return n
}

func extractExample(respObj any) string {
	respMap, ok := respObj.(map[string]any)
	if !ok {
		return ""
	}
	content, ok := respMap["content"].(map[string]any)
	if !ok {
		return ""
	}
	for _, mediaType := range content {
		mt, ok := mediaType.(map[string]any)
		if !ok {
			continue
		}
		if example, ok := mt["example"]; ok {
			if b, err := json.MarshalIndent(example, "", "  "); err == nil {
				return string(b)
			}
		}
	}
	return ""
}

func createFolderInCollection(store *CollectionStore, collectionID, name string) (*Collection, error) {
	c, ok := store.Update(collectionID, func(c *Collection) {
		c.Folders = append(c.Folders, Folder{
			ID:        uuid.New().String(),
			Name:      name,
			Endpoints: []Endpoint{},
		})
	})
	if !ok {
		return nil, fmt.Errorf("collection not found")
	}
	return c, nil
}

func addEndpointToFolder(store *CollectionStore, collectionID, folderID string, ep Endpoint) (*Collection, error) {
	c, ok := store.Update(collectionID, func(c *Collection) {
		for i := range c.Folders {
			if c.Folders[i].ID == folderID {
				c.Folders[i].Endpoints = append(c.Folders[i].Endpoints, ep)
				break
			}
		}
	})
	if !ok {
		return nil, fmt.Errorf("collection not found")
	}
	return c, nil
}
