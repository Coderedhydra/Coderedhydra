// jules/kb/local_kb_service.go
package kb

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LocalKBService is an implementation of KnowledgeBase that uses local files.
// For simplicity, this placeholder will assume KB articles are individual JSON files
// in a directory, or one large JSON file. More advanced would use a proper DB or index.
type LocalKBService struct {
	Articles    map[string]KBArticle // Articles loaded into memory, keyed by ID
	Initialized bool
	DataPath    string
}

func NewLocalKBService() *LocalKBService {
	return &LocalKBService{
		Articles: make(map[string]KBArticle),
	}
}

func (s *LocalKBService) Initialize(dataPath string) error {
	s.DataPath = dataPath
	fmt.Printf("KB_SERVICE: Initializing Local Knowledge Base from path: %s (Placeholder - loading mock data)\n", dataPath)

	// Placeholder: Load some mock data or attempt to load from dataPath
	// In a real scenario, this would walk dataPath, find JSON/Markdown files, parse them.
	mockArticles := []KBArticle{
		{
			ID:      "CVE-2021-44228",
			Title:   "Log4Shell - Apache Log4j Remote Code Execution",
			Source:  "CVE Mitre",
			Tags:    []string{"java", "log4j", "rce", "apache"},
			Content: "Apache Log4j2 <=2.14.1 JNDI features used in configuration, log messages, and parameters do not protect against attacker controlled LDAP and other JNDI related endpoints...",
		},
		{
			ID:      "OWASP-A03-Injection",
			Title:   "A03:2021 – Injection",
			Source:  "OWASP Top 10",
			Tags:    []string{"sqli", "xss", "nosqli", "command_injection", "owasp"},
			Content: "Injection flaws, such as SQL, NoSQL, Command Injection, etc., occur when untrusted data is sent to an interpreter as part of a command or query...",
		},
		{
			ID:      "EDB-ID-49502",
			Title:   "Apache Struts 2 - Remote Code Execution (S2-061)",
			Source:  "Exploit-DB",
			Tags:    []string{"java", "struts2", "rce", "apache", "ognl"},
			Content: "Apache Struts versions 2.0.0 to 2.5.25 are vulnerable to a Remote Code Execution vulnerability due to forced OGNL evaluation...",
		},
	}

	for _, article := range mockArticles {
		s.Articles[article.ID] = article
	}
	
	// Example: Try to load actual files if present (very basic)
	// This assumes each .json file in dataPath/cve/ or dataPath/owasp/ is a KBArticle
	// This is highly simplified. A real implementation would need more robust file handling.
	err := filepath.Walk(dataPath, func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			fmt.Printf("KB_SERVICE: Attempting to load KB article from %s\n", path)
			fileData, readErr := os.ReadFile(path)
			if readErr != nil {
				fmt.Printf("KB_SERVICE: Error reading file %s: %v\n", path, readErr)
				return nil // Continue walking
			}
			var article KBArticle
			jsonErr := json.Unmarshal(fileData, &article)
			if jsonErr != nil {
				fmt.Printf("KB_SERVICE: Error unmarshalling JSON from file %s: %v\n", path, jsonErr)
				return nil // Continue walking
			}
			if article.ID != "" { // Basic validation
				s.Articles[article.ID] = article
				fmt.Printf("KB_SERVICE: Loaded article '%s' from %s\n", article.ID, path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("KB_SERVICE: Error walking dataPath %s: %v\n", dataPath, err)
		// Not returning error for placeholder, just logging
	}


	s.Initialized = true
	fmt.Printf("KB_SERVICE: Loaded %d articles into the knowledge base.\n", len(s.Articles))
	return nil
}

func (s *LocalKBService) Search(query KBQuery) ([]KBArticle, error) {
	if !s.Initialized {
		return nil, fmt.Errorf("knowledge base not initialized")
	}
	var results []KBArticle
	// Simple search: iterates through all articles and checks for keyword/tag matches.
	// This is inefficient for large KBs; a real system would use an index (e.g., Bleve, Whoosh).
	for _, article := range s.Articles {
		match := false
		// Keyword match (any specified keyword found in title or content)
		if len(query.Keywords) > 0 {
			for _, keyword := range query.Keywords {
				if strings.Contains(strings.ToLower(article.Title), strings.ToLower(keyword)) ||
					strings.Contains(strings.ToLower(article.Content), strings.ToLower(keyword)) {
					match = true
					break
				}
			}
		} else {
			match = true // No keywords means match all (if tags also match or no tags specified)
		}

		if !match { continue } // Skip if keywords provided and none matched

		// Tag match (all specified tags must be present in article tags)
		if len(query.Tags) > 0 {
			tagMatchCount := 0
			for _, queryTag := range query.Tags {
				for _, articleTag := range article.Tags {
					if strings.EqualFold(queryTag, articleTag) {
						tagMatchCount++
						break
					}
				}
			}
			if tagMatchCount != len(query.Tags) { // All query tags must match
				match = false
			}
		}
		
		if match {
			results = append(results, article)
		}
	}

	// Apply limit
	if query.Limit > 0 && len(results) > query.Limit {
		results = results[:query.Limit]
	}
	fmt.Printf("KB_SERVICE: Search with keywords '%v' and tags '%v' found %d results.\n", query.Keywords, query.Tags, len(results))
	return results, nil
}

func (s *LocalKBService) GetArticleByID(id string) (*KBArticle, error) {
	if !s.Initialized {
		return nil, fmt.Errorf("knowledge base not initialized")
	}
	article, exists := s.Articles[id]
	if !exists {
		return nil, fmt.Errorf("article with ID '%s' not found", id)
	}
	return &article, nil
}

func (s *LocalKBService) Name() string {
	return "Local File-Based Knowledge Base (Placeholder)"
}
