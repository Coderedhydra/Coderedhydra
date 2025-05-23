// jules/kb/kb_test.go
package kb

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	// "github.com/stretchr/testify/assert" // Optional
)

func setupTestKB(t *testing.T) (*LocalKBService, string) {
	tempDir, err := os.MkdirTemp("", "jules-kb-test-data")
	if err != nil {
		t.Fatalf("Failed to create temp data dir: %v", err)
	}

	// Create some dummy JSON files
	article1 := KBArticle{ID: "TEST-001", Title: "Test SQLi Article", Tags: []string{"sqli", "test"}, Content: "SQL injection details..."}
	article2 := KBArticle{ID: "TEST-002", Title: "Test XSS Article", Tags: []string{"xss", "test"}, Content: "Cross-site scripting info..."}
	article3 := KBArticle{ID: "TEST-003", Title: "Generic Test Article", Tags: []string{"generic", "test"}, Content: "Some other details..."}

	for _, art := range []KBArticle{article1, article2, article3} {
		filePath := filepath.Join(tempDir, art.ID+".json")
		data, _ := json.Marshal(art)
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			t.Fatalf("Failed to write test article %s: %v", art.ID, err)
		}
	}
	
	// Add a non-JSON file to test robustness
	if err := os.WriteFile(filepath.Join(tempDir, "not-a-json.txt"), []byte("hello"), 0644); err != nil {
	    t.Fatalf("Failed to write non-json file: %v", err)
	}
    // Add a malformed JSON file
	if err := os.WriteFile(filepath.Join(tempDir, "malformed.json"), []byte("{not_json"), 0644); err != nil {
	    t.Fatalf("Failed to write malformed json file: %v", err)
	}


	kbService := NewLocalKBService()
	err = kbService.Initialize(tempDir) // Initialize with path to temp dir
	if err != nil {
		// The current Initialize logs errors but doesn't return them for file issues.
		// For a real test, we might want Initialize to return errors for critical load failures.
		t.Logf("KB Initialize logged errors (as expected for non-JSON/malformed files): %v", err)
	}
	return kbService, tempDir
}

func TestLocalKBService_InitializeAndLoad(t *testing.T) {
	kbService, tempDir := setupTestKB(t)
	defer os.RemoveAll(tempDir)

	if !kbService.Initialized {
		t.Fatal("KBService should be initialized")
	}
	// Check if mock data (hardcoded in Initialize) + file data is loaded.
	// Current Initialize loads mock data *then* file data, potentially overwriting.
	// For this test, we focus on file data being loaded.
	// The mock data in Initialize has 3 articles. The setupTestKB writes 3 more.
	// The current implementation of Initialize will have its hardcoded mock data overwritten if IDs clash,
	// or simply added to if IDs are unique. Let's assume they are unique for this count.
	expectedArticleCount := 3 + 3 // 3 from code, 3 from files
	if len(kbService.Articles) != expectedArticleCount {
		t.Errorf("Expected %d articles (3 mock + 3 file), got %d", expectedArticleCount, len(kbService.Articles))
	}


	art, err := kbService.GetArticleByID("TEST-001")
	if err != nil {
		t.Fatalf("GetArticleByID failed for TEST-001: %v", err)
	}
	if art.Title != "Test SQLi Article" {
		t.Errorf("Expected title 'Test SQLi Article', got '%s'", art.Title)
	}
	// assert.True(t, kbService.Initialized)
	// assert.GreaterOrEqual(t, len(kbService.Articles), 3) // 3 from files + mock
	// art, err := kbService.GetArticleByID("TEST-001")
	// assert.NoError(t, err)
	// assert.Equal(t, "Test SQLi Article", art.Title)
}

func TestLocalKBService_Search(t *testing.T) {
	kbService, tempDir := setupTestKB(t)
	defer os.RemoveAll(tempDir)

	// Test keyword search (should find TEST-001 from file)
	query1 := KBQuery{Keywords: []string{"SQL injection"}, Limit: 1}
	results1, err := kbService.Search(query1)
	if err != nil {
		t.Fatalf("Search failed for query1: %v", err)
	}
	if len(results1) != 1 {
		t.Fatalf("Expected 1 result for query1, got %d. Results: %+v", len(results1), results1)
	}
	if results1[0].ID != "TEST-001" {
		t.Errorf("Expected TEST-001 for query1, got %s", results1[0].ID)
	}
	// assert.NoError(t, err)
	// assert.Len(t, results1, 1)
	// assert.Equal(t, "TEST-001", results1[0].ID)

	// Test tag search (should find TEST-002 from file)
	query2 := KBQuery{Tags: []string{"xss"}, Limit: 1}
	results2, err := kbService.Search(query2)
	if err != nil {
		t.Fatalf("Search failed for query2: %v", err)
	}
	if len(results2) != 1 {
		t.Errorf("Expected 1 result for query2, got %d", len(results2))
	}
	if results2[0].ID != "TEST-002" {
		t.Errorf("Expected TEST-002 for query2, got %s", results2[0].ID)
	}

	// Test combined keyword and tag (should find TEST-003 from file)
	query3 := KBQuery{Keywords: []string{"details"}, Tags: []string{"generic"}, Limit: 1}
	results3, err := kbService.Search(query3)
	if err != nil {
		t.Fatalf("Search failed for query3: %v", err)
	}
	if len(results3) != 1 {
		t.Errorf("Expected 1 result for query3, got %d", len(results3))
	}
	if results3[0].ID != "TEST-003" {
		t.Errorf("Expected TEST-003 for query3, got %s", results3[0].ID)
	}
	
	// Test no results
	query4 := KBQuery{Keywords: []string{"nonexistentkeyword"}, Limit: 1}
	results4, err := kbService.Search(query4)
	if err != nil {
		t.Fatalf("Search failed for query4: %v", err)
	}
	if len(results4) != 0 {
		t.Errorf("Expected 0 results for query4, got %d", len(results4))
	}
}

func TestLocalKBService_GetArticleByID_NotFound(t *testing.T) {
	kbService, tempDir := setupTestKB(t)
	defer os.RemoveAll(tempDir)

	_, err := kbService.GetArticleByID("NONEXISTENT-ID")
	if err == nil {
		t.Error("Expected error for non-existent ID, got nil")
	}
	// assert.Error(t, err)
}
