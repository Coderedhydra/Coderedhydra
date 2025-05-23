// jules/core/core_test.go
package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
	// "github.com/stretchr/testify/assert" // A popular assertion library, optional
	"github.com/jules-org/jules/kb"  // For mock KB
	"github.com/jules-org/jules/llm" // For mock LLM
)

// Basic test to ensure Go testing is set up
func TestGoTestSetup(t *testing.T) {
	if false {
		t.Errorf("Go test setup is not working (this should not fail)")
	}
}

func TestJulesHTTPClient_SendCustomRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.Header.Get("Content-Type") == "application/json; charset=utf-8" {
			var data map[string]string
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				http.Error(w, "bad json", http.StatusBadRequest)
				return
			}
			if data["key"] == "payload" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"success from json"}`))
				return
			}
		}
		if r.Method == "GET" && r.URL.Query().Get("param") == "payload_in_query" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success from query"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Hello, test server!`))
	}))
	defer server.Close()

	client, err := NewJulesHTTPClient(5) // 5s timeout
	if err != nil {
		t.Fatalf("Failed to create JulesHTTPClient: %v", err)
	}

	// Test GET request with query parameters
	detailsGet := CustomRequestDetails{
		URL:             server.URL,
		Method:          "GET",
		QueryParameters: url.Values{"param": []string{"payload_in_query"}},
	}
	respGet, bodyGet, errGet := client.SendCustomRequest(detailsGet)
	if errGet != nil {
		t.Fatalf("SendCustomRequest (GET) failed: %v", errGet)
	}
	if respGet.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK for GET; got %v", respGet.Status)
	}
	if string(bodyGet) != `{"status":"success from query"}` {
		t.Errorf("Unexpected body for GET: %s", string(bodyGet))
	}
	// assert.NoError(t, errGet)
	// assert.Equal(t, http.StatusOK, respGet.StatusCode)
	// assert.Contains(t, string(bodyGet), "success from query")

	// Test POST request with JSON body
	detailsPost := CustomRequestDetails{
		URL:      server.URL,
		Method:   "POST",
		BodyType: BodyTypeJSON,
		JSONBody: map[string]string{"key": "payload"},
		Headers:  map[string]string{"X-Custom-Header": "JulesTest"},
	}
	respPost, bodyPost, errPost := client.SendCustomRequest(detailsPost)
	if errPost != nil {
		t.Fatalf("SendCustomRequest (POST) failed: %v", errPost)
	}
	if respPost.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK for POST; got %v", respPost.Status)
	}
	if respPost.Header.Get("X-Custom-Header") != "" { // Server doesn't reflect this header back
		// This test would be for client correctly *sending* it.
		// To test if server received it, server would need to echo it.
	}
	if string(bodyPost) != `{"status":"success from json"}` {
		t.Errorf("Unexpected body for POST: %s", string(bodyPost))
	}
	// assert.NoError(t, errPost)
	// assert.Equal(t, http.StatusOK, respPost.StatusCode)
	// assert.Contains(t, string(bodyPost), "success from json")
}

// --- Mock LLM for Plugin System Test ---
type MockTestLLM struct {
	KnowledgeBase kb.KnowledgeBase // Embed KB if RefineWithKnowledge needs it directly
}

func (m *MockTestLLM) Initialize(configPath string) error { return nil }
func (m *MockTestLLM) AnalyzeInputContext(req llm.LLMRequest) (*llm.LLMAnalysisResponse, error) {
	return &llm.LLMAnalysisResponse{InputClassification: map[string]string{"type": "mock"}, PotentialVulns: []string{"mock-vuln"}}, nil
}
func (m *MockTestLLM) GeneratePayloads(req llm.LLMRequest) (*llm.LLMPayloadResponse, error) {
	payload := "<llm_generated_payload_for_" + req.InputContext.Name + ">"
	if len(req.KnowledgeBaseHits) > 0 {
		payload = "<llm_refined_payload_for_" + req.InputContext.Name + "_using_kb: " + req.KnowledgeBaseHits[0] + ">"
	}
	return &llm.LLMPayloadResponse{GeneratedPayloads: []string{payload}, ConfidenceScores: []float64{0.9}}, nil
}
func (m *MockTestLLM) RefineWithKnowledge(keywords []string, originalContext llm.LLMRequest) (*llm.LLMPayloadResponse, error) {
	// Simplified: directly call GeneratePayloads, assuming KB results are now in originalContext.KnowledgeBaseHits
	// A real RefineWithKnowledge in QwenLocalClient would first query its KB instance.
	// Here, we simulate that the KB results are already populated.
	if m.KnowledgeBase != nil && len(keywords) > 0 {
		// Simulate KB search and adding to originalContext
		originalContext.KnowledgeBaseHits = append(originalContext.KnowledgeBaseHits, "mock_kb_hit_from_keywords: "+keywords[0])
	}
	return m.GeneratePayloads(originalContext)
}
func (m *MockTestLLM) Name() string { return "MockTestLLM" }

// --- Mock KB for Plugin System Test ---
type MockTestKB struct{}

func (mkb *MockTestKB) Initialize(dataPath string) error { return nil }
func (mkb *MockTestKB) Search(query kb.KBQuery) ([]kb.KBArticle, error) {
	if len(query.Keywords) > 0 && query.Keywords[0] == "error_trigger_keyword" {
		return []kb.KBArticle{{ID: "KB-ERR-01", Title: "Error Related Info", Content: "Info about error"}}, nil
	}
	return []kb.KBArticle{}, nil
}
func (mkb *MockTestKB) GetArticleByID(id string) (*kb.KBArticle, error) { return &kb.KBArticle{ID: id}, nil }
func (mkb *MockTestKB) Name() string                                     { return "MockTestKB" }

// Mock ScanPlugin for testing PluginManager
type MockPluginForManagerTest struct {
	ShouldError bool
	UsesLLM     bool
	UsesKB      bool
}

func (p *MockPluginForManagerTest) Name() string        { return "MockPluginForManager" }
func (p *MockPluginForManagerTest) Description() string { return "A mock plugin for manager tests." }
func (p *MockPluginForManagerTest) RunScan(targetURL string, inputCtx *llm.FormInputContext, httpClient *JulesHTTPClient, llmClient llm.JulesLLM) ([]Issue, error) {
	if p.ShouldError {
		return nil, fmt.Errorf("mock plugin intentional error")
	}
	var issues []Issue
	payload := "static_payload"
	llmSource := ""
	var kbArticles []string

	if p.UsesLLM && llmClient != nil && inputCtx != nil {
		var llmResp *llm.LLMPayloadResponse
		var err error
		if p.UsesKB {
			// Simulate keywords that would trigger KB search in a real scenario
			llmResp, err = llmClient.RefineWithKnowledge([]string{"error_trigger_keyword"}, llm.LLMRequest{InputContext: *inputCtx, ScanType: "test_scan"})
			llmSource = "llm_refined_from_kb"
			// This kbArticle ID should match what MockTestKB would return for "error_trigger_keyword"
			// or what MockTestLLM's RefineWithKnowledge simulates adding.
			// Based on MockTestLLM's RefineWithKnowledge, it adds "mock_kb_hit_from_keywords: error_trigger_keyword"
			// to KnowledgeBaseHits. So, if the plugin is to populate ContributingKBArticles based on
			// what was *used* by the LLM from its own KB search (done inside RefineWithKnowledge),
			// then the plugin would need to somehow know which KB articles were actually used.
			// For this test, we'll assume the plugin *knows* a KB article influenced the payload
			// because it asked for KB refinement. A more complex setup might involve the LLM
			// response itself indicating which KB articles were influential.
			kbArticles = append(kbArticles, "mock_kb_hit_from_keywords: error_trigger_keyword")
		} else {
			llmResp, err = llmClient.GeneratePayloads(llm.LLMRequest{InputContext: *inputCtx, ScanType: "test_scan"})
			llmSource = "llm_generated"
		}
		if err == nil && llmResp != nil && len(llmResp.GeneratedPayloads) > 0 {
			payload = llmResp.GeneratedPayloads[0]
		}
	}

	issues = append(issues, Issue{
		Type:                   "Mock Vuln",
		Severity:               "Info",
		URL:                    targetURL,
		InputName:              inputCtx.Name,
		Description:            "Mock issue from " + p.Name(),
		PayloadUsed:            payload,
		LLMPayloadSource:       llmSource,
		ContributingKBArticles: kbArticles,
	})
	return issues, nil
}

func TestPluginManager_RunScans_WithLLMAndKB(t *testing.T) {
	pm := NewPluginManager()
	pm.RegisterPlugin(&MockPluginForManagerTest{UsesLLM: true, UsesKB: true})

	mockLLM := &MockTestLLM{KnowledgeBase: &MockTestKB{}} // Give mock LLM a mock KB
	mockHTTPClient, _ := NewJulesHTTPClient(5)

	contexts := []llm.FormInputContext{
		{Name: "param1", Type: "text", HtmlContext: "input_value"},
	}

	issues := pm.RunScans("http://example.com", contexts, mockHTTPClient, mockLLM)

	if len(issues) != 1 {
		t.Fatalf("Expected 1 issue, got %d", len(issues))
	}
	if issues[0].InputName != "param1" {
		t.Errorf("Expected issue for param1, got %s", issues[0].InputName)
	}
	if issues[0].LLMPayloadSource != "llm_refined_from_kb" {
		t.Errorf("Expected LLM source 'llm_refined_from_kb', got '%s'", issues[0].LLMPayloadSource)
	}
	if !strings.Contains(issues[0].PayloadUsed, "mock_kb_hit_from_keywords: error_trigger_keyword") {
		t.Errorf("Expected payload to contain KB refinement, got '%s'", issues[0].PayloadUsed)
	}
	if len(issues[0].ContributingKBArticles) == 0 || !strings.Contains(issues[0].ContributingKBArticles[0], "mock_kb_hit_from_keywords: error_trigger_keyword") {
		t.Errorf("Expected ContributingKBArticles to reflect KB usage, got %v", issues[0].ContributingKBArticles)
	}
	// assert.Len(t, issues, 1)
	// assert.Equal(t, "param1", issues[0].InputName)
	// assert.Equal(t, "llm_refined_from_kb", issues[0].LLMPayloadSource)
	// assert.Contains(t, issues[0].PayloadUsed, "mock_kb_hit_from_keywords: error_trigger_keyword")
}

func TestReportGenerator_JSON(t *testing.T) { // Renamed from TestGenerateJSONReport
	tempFile, err := os.CreateTemp("", "test_report.*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up

	issues := []Issue{
		{
			Type: "TestVuln", Severity: "High", URL: "http://example.com/vuln", InputName: "q",
			Description:            "A test vulnerability",
			PayloadUsed:            "<script>alert(1)</script>",
			LLMPayloadSource:       "llm_generated",
			LLMConfidence:          0.9,
			ContributingKBArticles: []string{"KB-001"}, // This makes it KB-refined too
		},
		{Type: "AnotherVuln", Severity: "Medium", URL: "http://example.com/page", InputName: "p"},
	}

	err = GenerateJSONReport("scan001", "http://example.com", time.Now(), issues, tempFile.Name())
	if err != nil {
		t.Fatalf("GenerateJSONReport failed: %v", err)
	}

	reportData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read report file: %v", err)
	}
	var reportOutput Report
	if err := json.Unmarshal(reportData, &reportOutput); err != nil {
		t.Fatalf("Failed to unmarshal report JSON: %v", err)
	}

	if reportOutput.TotalIssuesFound != 2 {
		t.Errorf("Expected 2 total issues, got %d", reportOutput.TotalIssuesFound)
	}
	if reportOutput.TotalLLMAssistedFindings != 1 {
		t.Errorf("Expected 1 LLM assisted finding, got %d", reportOutput.TotalLLMAssistedFindings)
	}
	if reportOutput.TotalKBRefinedFindings != 1 { // Because ContributingKBArticles is present
		t.Errorf("Expected 1 KB refined finding, got %d", reportOutput.TotalKBRefinedFindings)
	}
	// assert.NoError(t, err)
	// assert.Equal(t, 2, reportOutput.TotalIssuesFound)
	// assert.Equal(t, 1, reportOutput.TotalLLMAssistedFindings)
}

// TODO: Add more specific unit tests for each module.
// Consider table-driven tests for multiple scenarios.
