package backend

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// ImportFromOpenAPI parses an OpenAPI 3.x JSON spec and creates a collection.
func ImportFromOpenAPI(store *CollectionStore, data string) (*Collection, error) {
	var spec map[string]interface{}
	if err := json.Unmarshal([]byte(data), &spec); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}

	title := "Imported Collection"
	description := ""
	if info, ok := spec["info"].(map[string]interface{}); ok {
		if t, ok := info["title"].(string); ok {
			title = t
		}
		if d, ok := info["description"].(string); ok {
			description = d
		}
	}

	c := store.Create(title, description)

	paths, ok := spec["paths"].(map[string]interface{})
	if !ok {
		return &c, nil
	}

	folderMap := make(map[string]string) // tag → folderID

	for path, pathItem := range paths {
		pathItemMap, ok := pathItem.(map[string]interface{})
		if !ok {
			continue
		}
		for rawMethod, operation := range pathItemMap {
			method := strings.ToUpper(rawMethod)
			if method == "PARAMETERS" || method == "SERVERS" || method == "SUMMARY" || method == "DESCRIPTION" {
				continue
			}
			opMap, ok := operation.(map[string]interface{})
			if !ok {
				continue
			}

			tag := "Default"
			if tags, ok := opMap["tags"].([]interface{}); ok && len(tags) > 0 {
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
				ID:        uuid.New().String(),
				Method:    method,
				Path:      path,
				Strategy:  "fixed",
				Responses: []MockResponse{},
			}

			if opID, ok := opMap["operationId"].(string); ok {
				ep.Name = opID
			} else if summary, ok := opMap["summary"].(string); ok {
				ep.Name = summary
			} else {
				ep.Name = method + " " + path
			}

			if responses, ok := opMap["responses"].(map[string]interface{}); ok {
				for statusStr, respObj := range responses {
					statusCode := 200
					fmt.Sscanf(statusStr, "%d", &statusCode)

					mockResp := MockResponse{
						ID:         uuid.New().String(),
						Name:       statusStr,
						StatusCode: statusCode,
						BodyType:   "json",
						Body:       "{}",
						Headers:    []KVPair{},
						Cookies:    []MockCookie{},
					}

					if respMap, ok := respObj.(map[string]interface{}); ok {
						if content, ok := respMap["content"].(map[string]interface{}); ok {
							for _, mediaType := range content {
								if mt, ok := mediaType.(map[string]interface{}); ok {
									if example, ok := mt["example"]; ok {
										if b, err := json.MarshalIndent(example, "", "  "); err == nil {
											mockResp.Body = string(b)
										}
									}
								}
								break
							}
						}
					}
					ep.Responses = append(ep.Responses, mockResp)
				}
			}

			if len(ep.Responses) == 0 {
				ep.Responses = []MockResponse{{
					ID:         uuid.New().String(),
					Name:       "200 OK",
					StatusCode: 200,
					Body:       "{}",
					BodyType:   "json",
					Headers:    []KVPair{},
					Cookies:    []MockCookie{},
				}}
			}

			updated, err := addEndpointToFolder(store, c.ID, folderID, ep)
			if err == nil && updated != nil {
				c = *updated
			}
		}
	}

	return &c, nil
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
