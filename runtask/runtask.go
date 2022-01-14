package runtask

import (
	"log"
  "encoding/json"
	"net/http"
  "io"
)

// TFCInitRequest - The request TFC sends to our runtask
type TFCInitRequest struct {
  PayloadVersion int `json:"payload_version"`
  AccessToken string `json:"access_token"`
  TaskResultID string `json:"task_result_id"`
  TaskResultCallbackURL string `json:"task_result_callback_url"`
  RunAppURL string `json:"run_app_url"`
  RunID string `json:"run_id"`
  RunMessage string `json:"run_message"`
  RunCreatedAt string `json:"run_created_ad"`
  RunCreatedBy string `json:"run_created_by"`
  WorkspaceID string `json:"workspace_id"`
  WorkspaceName string `json:"workspace_name"`
  WorkspaceAppURL string `json:"workspace_app_url"`
  OrganizationName string `json:"organization_name"`
  PlanJSONApiURL string `json:"plan_json_api_url"`
  VCSRepoURL string `json:"vcs_repo_url"`
  VCSBranch string `json:"vcs_branch"`
  VCSPullRequestURL string `json:"vcs_pull_request_url"`
  VCSCommitURL string `json:"vcs_commit_url"`
}

// TFCTaskResponse - Response that we send to TFC after we finish 'task processing'
type TFCTaskResponse struct {
  Data struct {
    Type string `json:"data"`
    Attributes struct {
      Status string `json:"status"`
      Message string `json:"message"`
      URL string `json:"url"`
    } `json:"attributes"`
  } `json:"data"`
}

// RootHandler - Root placeholder
func RootHandler(w http.ResponseWriter, req *http.Request) {

  w.WriteHeader(http.StatusOK)
  w.Header().Set("Content-Type", "application/json")

  resp := make(map[string]string)
  resp["message"] = "Status OK: (200)"

  jsonResp, err := json.Marshal(resp)
  if err != nil {
    log.Fatalf("Error happened in JSON marshal. Err: %s", err)
  }
  w.Write(jsonResp)

  return

}

// InitHandler - handles initial connection from TFC
func InitHandler(w http.ResponseWriter, req *http.Request) {

  w.Header().Set("Content-Type", "application/json")
  bodyBytes, err := io.ReadAll(req.Body)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    log.Fatalf("Error body parsing.(404) Err: %s", err)
  }
  resp := make(map[string]string)
  resp["body"] = string(bodyBytes)

  jsonResp, err := json.Marshal(resp)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    log.Fatalf("Error happened in JSON marshal. (404) Err: %s", err)
  }

  w.WriteHeader(http.StatusOK)
  w.Write(jsonResp)


  return

}
