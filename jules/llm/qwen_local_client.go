// jules/llm/qwen_local_client.go
package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath" // For robust path construction
	"strings"       // For processing KB content

	"github.com/jules-org/jules/kb" // Import the Knowledge Base package
)

type QwenLocalClient struct {
	ModelPath        string // Path to the model files (used by Python script)
	PythonScriptPath string // Path to the qwen_model_handler.py script
	PythonExecutable string // Path to python interpreter (e.g., "python3" or "/usr/bin/python3")
	KnowledgeBase    kb.KnowledgeBase // Instance of the knowledge base service
	Initialized      bool
}

// NewQwenLocalClient creates a new client.
// It now accepts a KnowledgeBase instance.
func NewQwenLocalClient(pythonExecutable, scriptDir, modelPath string, knowledgeService kb.KnowledgeBase) *QwenLocalClient {
	return &QwenLocalClient{
		ModelPath:        modelPath,
		PythonScriptPath: filepath.Join(scriptDir, "qwen_model_handler.py"),
		PythonExecutable: pythonExecutable,
		KnowledgeBase:    knowledgeService, // Store the KB service
	}
}

func (q *QwenLocalClient) Initialize(configPath string) error {
	if _, err := exec.LookPath(q.PythonExecutable); err != nil {
		return fmt.Errorf("python executable '%s' not found in PATH: %w", q.PythonExecutable, err)
	}
	fmt.Printf("LLM_CLIENT: Initializing Qwen LLM via Python script: %s (Model: %s)\n", q.PythonScriptPath, q.ModelPath)
	if q.KnowledgeBase != nil {
		fmt.Printf("LLM_CLIENT: Knowledge Base '%s' is configured.\n", q.KnowledgeBase.Name())
	} else {
		fmt.Printf("LLM_CLIENT: No Knowledge Base is configured.\n")
	}
	q.Initialized = true
	return nil
}

// callPythonScript helper remains largely the same
func (q *QwenLocalClient) callPythonScript(request LLMRequest) ([]byte, error) {
	if !q.Initialized {
		fmt.Println("LLM_CLIENT: QwenLocalClient not explicitly initialized, attempting lazy initialization...")
		err := q.Initialize("") 
		if err != nil {
			return nil, fmt.Errorf("lazy initialization failed: %w", err)
		}
	}
	
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal LLMRequest to JSON: %w", err)
	}

	cmd := exec.Command(q.PythonExecutable, q.PythonScriptPath, "--model_path", q.ModelPath, "--request_json", string(requestJSON))
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// fmt.Printf("LLM_CLIENT: Executing Python LLM handler: %s %s --model_path %s --request_json '%s'\n", q.PythonExecutable, q.PythonScriptPath, q.ModelPath, string(requestJSON))

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("python script execution failed: %w. Stderr: %s", err, stderr.String())
	}
	return stdout.Bytes(), nil
}


func (q *QwenLocalClient) AnalyzeInputContext(req LLMRequest) (*LLMAnalysisResponse, error) {
	req.TaskType = "context_analysis"
	outputBytes, err := q.callPythonScript(req)
	if err != nil {
		return nil, fmt.Errorf("python script call for context analysis failed: %w", err)
	}
	var response LLMAnalysisResponse
	err = json.Unmarshal(outputBytes, &response)
	if err != nil {
		var pyError map[string]string
		if json.Unmarshal(outputBytes, &pyError) == nil && pyError["error"] != "" {
			return nil, fmt.Errorf("python script returned error for context_analysis: %s. Raw: %s", pyError["error"], string(outputBytes))
		}
		return nil, fmt.Errorf("failed to unmarshal JSON response from python script for context_analysis: %w. Raw output: %s", err, string(outputBytes))
	}
	fmt.Printf("LLM_CLIENT: Received LLM Analysis for '%s': Potential Vulns: %v\n", req.InputContext.Name, response.PotentialVulns)
	return &response, nil
}

func (q *QwenLocalClient) GeneratePayloads(req LLMRequest) (*LLMPayloadResponse, error) {
	req.TaskType = "payload_generation"
	outputBytes, err := q.callPythonScript(req)
	if err != nil {
		return nil, fmt.Errorf("python script call for payload generation failed: %w", err)
	}
	var response LLMPayloadResponse
	err = json.Unmarshal(outputBytes, &response)
	if err != nil {
		var pyError map[string]string
		if json.Unmarshal(outputBytes, &pyError) == nil && pyError["error"] != "" {
			return nil, fmt.Errorf("python script returned error for payload_generation: %s. Raw: %s", pyError["error"], string(outputBytes))
		}
		return nil, fmt.Errorf("failed to unmarshal JSON response from python script for payload_generation: %w. Raw output: %s", err, string(outputBytes))
	}
	fmt.Printf("LLM_CLIENT: Received %d LLM Payloads for '%s' (ScanType: %s)\n", len(response.GeneratedPayloads), req.InputContext.Name, req.ScanType)
	return &response, nil
}

// RefineWithKnowledge now uses the configured KnowledgeBase to augment the LLMRequest.
// keywordsForKB are terms extracted from scanner context (e.g., error messages, server tech).
// originalContext is the LLMRequest that led to wanting KB refinement.
func (q *QwenLocalClient) RefineWithKnowledge(keywordsForKB []string, originalContext LLMRequest) (*LLMPayloadResponse, error) {
	if q.KnowledgeBase == nil {
		fmt.Println("LLM_CLIENT: KnowledgeBase not configured, cannot refine with knowledge. Proceeding with original context.")
		// Fallback to standard payload generation without KB insights
		return q.GeneratePayloads(originalContext)
	}
	if len(keywordsForKB) == 0 {
		fmt.Println("LLM_CLIENT: No keywords provided for KB search. Proceeding with original context for refinement.")
		return q.GeneratePayloads(originalContext)
	}

	fmt.Printf("LLM_CLIENT: Refining LLM request for input '%s' using KB with keywords: %v\n", originalContext.InputContext.Name, keywordsForKB)
	
	kbQuery := kb.KBQuery{
		Keywords:  keywordsForKB,
		Limit:     3, // Limit number of KB articles to avoid overly long prompts
		MatchType: "any_keyword", // Placeholder, KB service implements actual search logic
	}
	articles, err := q.KnowledgeBase.Search(kbQuery)
	if err != nil {
		fmt.Printf("LLM_CLIENT: Error searching knowledge base: %v. Proceeding without KB insights.\n", err)
		// Fallback to standard payload generation
		return q.GeneratePayloads(originalContext)
	}

	if len(articles) == 0 {
		fmt.Printf("LLM_CLIENT: No relevant articles found in KB for keywords: %v. Proceeding without KB insights.\n", keywordsForKB)
		return q.GeneratePayloads(originalContext)
	}

	// Prepare KB hits for the LLM prompt
	if originalContext.KnowledgeBaseHits == nil {
		originalContext.KnowledgeBaseHits = []string{}
	}
	for _, article := range articles {
		// Combine title and a snippet of content. Be mindful of prompt length.
		hit := fmt.Sprintf("KB ID: %s, Title: %s, Source: %s, Tags: [%s], Snippet: %s...",
			article.ID, article.Title, article.Source, strings.Join(article.Tags, ", "), article.Content)
		maxSnippetLength := 150 // Max length for snippet from content
		if len(article.Content) > maxSnippetLength {
			hit = fmt.Sprintf("KB ID: %s, Title: %s, Source: %s, Tags: [%s], Snippet: %s...",
				article.ID, article.Title, article.Source, strings.Join(article.Tags, ", "), article.Content[:maxSnippetLength])
		} else {
			hit = fmt.Sprintf("KB ID: %s, Title: %s, Source: %s, Tags: [%s], Snippet: %s",
				article.ID, article.Title, article.Source, strings.Join(article.Tags, ", "), article.Content)
		}
		originalContext.KnowledgeBaseHits = append(originalContext.KnowledgeBaseHits, hit)
	}
	
	fmt.Printf("LLM_CLIENT: Augmented LLM request with %d hits from KB.\n", len(articles))

	// The task type is still payload generation, but now with added KB context
	originalContext.TaskType = "payload_generation" 
	return q.GeneratePayloads(originalContext)
}

func (q *QwenLocalClient) Name() string {
	kbName := "None"
	if q.KnowledgeBase != nil {
		kbName = q.KnowledgeBase.Name()
	}
	return fmt.Sprintf("Qwen3-8B-AWQ (Local via %s %s, KB: %s)", q.PythonExecutable, filepath.Base(q.PythonScriptPath), kbName)
}
