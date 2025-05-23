// jules/kb/kb_interface.go
package kb

// KBArticle represents a single piece of information in the knowledge base.
type KBArticle struct {
	ID          string   `json:"id"`          // Unique identifier (e.g., CVE-2023-XXXX, OWASP-A01, EDB-ID-YYYYY)
	Title       string   `json:"title"`       // Title of the article/entry
	Source      string   `json:"source"`      // e.g., "CVE Mitre", "OWASP", "Exploit-DB"
	Tags        []string `json:"tags"`        // Keywords for searching (e.g., "sqli", "xss", "apache", "nginx", "reflected")
	Content     string   `json:"content"`     // Textual content of the article (can be markdown, plain text)
	RelatedCVEs []string `json:"related_cves,omitempty"`
	ExploitRefs []string `json:"exploit_refs,omitempty"` // References to exploit code or details
	// Add more structured fields as needed, e.g., CVSS scores, affected software versions.
}

// KBQuery represents a query to search the knowledge base.
type KBQuery struct {
	Keywords    []string // Primary search terms (e.g., from error messages, software names)
	Tags        []string // Specific tags to filter by
	Limit       int      // Max number of results to return
	MatchType   string   // "any_keyword", "all_keywords", "exact_phrase" (placeholder)
}

// KnowledgeBase defines the interface for interacting with the vulnerability knowledge base.
type KnowledgeBase interface {
	// Initialize loads or sets up the knowledge base from a given path or configuration.
	Initialize(dataPath string) error

	// Search performs a query against the knowledge base.
	Search(query KBQuery) ([]KBArticle, error)

	// GetArticleByID retrieves a specific article by its ID.
	GetArticleByID(id string) (*KBArticle, error)

	// Name returns a descriptive name for the KB implementation.
	Name() string
}
