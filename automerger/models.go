package automerger

import "encoding/json"

type Namer struct {
	Name string `json:"name"`
}

type PullRequestResponse struct {
	Number int `json:"number"`
}

type MergeRequest struct {
	MergeMethod string `json:"merge_method"`
}

type PullRequest struct {
	Title string `json:"title"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}

type Repository struct {
	FullName string `json:"full_name"`
}

type Installation struct {
	ID int `json:"id"`
}

type PushEvent struct {
	Ref          string       `json:"ref"`
	Repository   Repository   `json:"repository"`
	Installation Installation `json:"installation"`
	Pusher       Namer        `json:"pusher"`
}

type PushEventHandler struct {
	Config   *Config
	GetToken func(int) (string, []error)
}

type PullRequestAssignees struct {
	Assignees Assignees `json:"assignees"`
}

type ErrorResp struct {
	Errors []error
}

func (e *ErrorResp) MarshalJSON() ([]byte, error) {
	type ErrJson struct {
		Errors []string `json:"errors"`
	}
	var k = &ErrJson{Errors: make([]string, len(e.Errors), len(e.Errors))}
	for i, err := range e.Errors {
		k.Errors[i] = err.Error()
	}

	return json.Marshal(k)
}
