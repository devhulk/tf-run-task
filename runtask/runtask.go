package runtask

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/devhulk/test-task/jwt"
)

// TFCInitRequest - The request TFC sends to our runtask
type TFCInitRequest struct {
	PayloadVersion        int    `json:"payload_version"`
	AccessToken           string `json:"access_token"`
	TaskResultID          string `json:"task_result_id"`
	TaskResultCallbackURL string `json:"task_result_callback_url"`
	RunAppURL             string `json:"run_app_url"`
	RunID                 string `json:"run_id"`
	RunMessage            string `json:"run_message"`
	RunCreatedAt          string `json:"run_created_ad"`
	RunCreatedBy          string `json:"run_created_by"`
	WorkspaceID           string `json:"workspace_id"`
	WorkspaceName         string `json:"workspace_name"`
	WorkspaceAppURL       string `json:"workspace_app_url"`
	OrganizationName      string `json:"organization_name"`
	PlanJSONApiURL        string `json:"plan_json_api_url"`
	VCSRepoURL            string `json:"vcs_repo_url"`
	VCSBranch             string `json:"vcs_branch"`
	VCSPullRequestURL     string `json:"vcs_pull_request_url"`
	VCSCommitURL          string `json:"vcs_commit_url"`
}

// TFCTaskResponse - Response that we send to TFC after we finish 'task processing'
type TFCTaskResponse struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Status  string `json:"status"`
			Message string `json:"message"`
			URL     string `json:"url"`
		} `json:"attributes"`
	} `json:"data"`
}

type TFCVariableReq struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Key         string `json:"key"`
			Value       string `json:"value"`
			Description string `json:"description"`
			Category    string `json:"category"`
			Hcl         bool   `json:"hcl"`
			Sensitive   bool   `json:"sensitive"`
		} `json:"attributes"`
		Relationships struct {
			Workspace struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"workspace"`
		} `json:"relationships"`
	} `json:"data"`
}

// TaskHandler - handles initial connection from TFC
func TaskHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("SERVER HIT!")
	w.Header().Set("Content-Type", "application/json")
	tfcreq := TFCInitRequest{}
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatalf("Error body parsing.(404) Err: %s", err)
	}
	json.Unmarshal(bodyBytes, &tfcreq)
	fmt.Printf("CallbackURL: %s", tfcreq.TaskResultCallbackURL)
	fmt.Printf("AccessToken: %s", tfcreq.AccessToken)
	fmt.Printf("Authorization: Bearer %s", tfcreq.AccessToken)

	w.WriteHeader(http.StatusOK)
	handleCallback(&tfcreq)

	return

}

// HandleCallback - evaluate task and execute callback to tfc
func handleCallback(t *TFCInitRequest) {

	var response *TFCTaskResponse

	fmt.Println("Formulating Callback Response...")

	fmt.Println("Deciding if you pass or fail my amazing test...")

	// update workspace variable  with JWT token
	_, err := setTFCVariable(t)
	if err != nil {
		log.Fatalf("Error setting JWT token. (404) Err: %s", err)
		response = passOrFail("fail")
		return
	} else {
		response = passOrFail("pass")
	}

	taskResult, err := json.Marshal(&response)
	if err != nil {
		log.Fatalf("Error happened in task result json marshal. (404) Err: %s", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPatch, t.TaskResultCallbackURL, bytes.NewBuffer(taskResult))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	if err != nil {
		log.Fatalf("Error happened in callback. (404) Err: %s", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))

	fmt.Println("Executing callback...")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error happened in callback client call. (404) Err: %s", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading callback response body. Err: %s ", err)
	}

	log.Println(string(body))
	return
}

func passOrFail(decision string) *TFCTaskResponse {
	// TODO: Pass or Fail based on OPA policy

	pass := TFCTaskResponse{
		Data: struct {
			Type       string `json:"type"`
			Attributes struct {
				Status  string `json:"status"`
				Message string `json:"message"`
				URL     string `json:"url"`
			} `json:"attributes"`
		}{
			Type: "task-results",
			Attributes: struct {
				Status  string `json:"status"`
				Message string `json:"message"`
				URL     string `json:"url"`
			}{
				Status:  "passed",
				Message: "YOUUUU SHALLLL PASSSSSSS",
				URL:     "https://devhulk.ddns.net",
			},
		},
	}

	fail := TFCTaskResponse{
		Data: struct {
			Type       string `json:"type"`
			Attributes struct {
				Status  string `json:"status"`
				Message string `json:"message"`
				URL     string `json:"url"`
			} `json:"attributes"`
		}{
			Type: "task-results",
			Attributes: struct {
				Status  string `json:"status"`
				Message string `json:"message"`
				URL     string `json:"url"`
			}{
				Status:  "failed",
				Message: "You were the CHOSEN ONEEE! IT WAS SAID YOU WOULD DESTROY THE SITH NOT JOIN THEM!",
				URL:     "https://devhulk.ddns.net",
			},
		},
	}

	if decision == "pass" {
		return &pass
	}
	return &fail
}

func setTFCVariable(t *TFCInitRequest) (string, error) {

	jwtToken, err := jwt.GetJWT()
	if err != nil {
		log.Fatalf("Error creating JWT: %v", err)
	}

	tfcVar := &TFCVariableReq{
		Data: struct {
			Type       string "json:\"type\""
			Attributes struct {
				Key         string "json:\"key\""
				Value       string "json:\"value\""
				Description string "json:\"description\""
				Category    string "json:\"category\""
				Hcl         bool   "json:\"hcl\""
				Sensitive   bool   "json:\"sensitive\""
			} "json:\"attributes\""
			Relationships struct {
				Workspace struct {
					Data struct {
						ID   string "json:\"id\""
						Type string "json:\"type\""
					} "json:\"data\""
				} "json:\"workspace\""
			} "json:\"relationships\""
		}{
			Type: "vars",
			Attributes: struct {
				Key         string "json:\"key\""
				Value       string "json:\"value\""
				Description string "json:\"description\""
				Category    string "json:\"category\""
				Hcl         bool   "json:\"hcl\""
				Sensitive   bool   "json:\"sensitive\""
			}{
				Key:         "jwt_token",
				Value:       jwtToken,
				Description: "A cool JWT being set by an even cooler run-task",
				Category:    "terraform",
				Hcl:         false,
				Sensitive:   false,
			},
			Relationships: struct {
				Workspace struct {
					Data struct {
						ID   string "json:\"id\""
						Type string "json:\"type\""
					} "json:\"data\""
				} "json:\"workspace\""
			}{
				Workspace: struct {
					Data struct {
						ID   string "json:\"id\""
						Type string "json:\"type\""
					} "json:\"data\""
				}{
					Data: struct {
						ID   string "json:\"id\""
						Type string "json:\"type\""
					}{
						ID:   t.WorkspaceID,
						Type: "workspaces",
					},
				},
			},
		},
	}

	tfcVarReq, err := json.Marshal(&tfcVar)
	if err != nil {
		log.Fatalf("Error happened in TFC Variable json marshal. (404) Err: %s", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://app.terraform.io/api/v2/workspaces/%v/vars", t.WorkspaceID), bytes.NewBuffer(tfcVarReq))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	if err != nil {
		log.Fatalf("Error happened in callback. (404) Err: %s", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))

	fmt.Println("Executing callback...")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error happened in callback client call. (404) Err: %s", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading callback response body. Err: %s ", err)
	}

	return fmt.Sprint(body), err
}
