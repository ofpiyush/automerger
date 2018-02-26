package automerger

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (p *PushEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("X-GitHub-Event") != "push" {
		httpErrors(w, http.StatusAccepted, ErrUnknownEvent)
		return
	}
	signature := r.Header.Get("X-Hub-Signature")
	if signature == "" {
		httpErrors(w, http.StatusBadRequest, ErrBadSignature)
		return
	}
	var v PushEvent
	body, err := ioutil.ReadAll(r.Body)
	if httpErrors(w, http.StatusAccepted, err) {
		return
	}

	if httpErrors(w, http.StatusBadRequest, verifySignature(p.Config.Secret, signature, body)) {
		return
	}

	if httpErrors(w, http.StatusAccepted, json.Unmarshal(body, &v)) {
		return
	}
	if v.Ref != p.Config.Ref {
		// We don't want to deal with these cases
		httpErrors(w, http.StatusAccepted, ErrUnknownBranch)
		return
	}
	go p.MergeBranches(URL(p.Config.ApiURL, "repos", v.Repository.FullName), v.Pusher.Name, v.Installation.ID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "Scheduled merging",
		"repo":   v.Repository.FullName,
	})
}

func (p *PushEventHandler) MergeBranches(repoURL, assignee string, installationID int) {
	token, errs := p.GetToken(installationID)
	if len(handleErrs(errs...)) > 0 {
		return
	}
	var branches []Namer
	if len(handleErrs(GithubRequest("GET", URL(repoURL, "branches"), token, http.StatusOK, nil, &branches)...)) > 0 {
		return
	}
	for _, branch := range branches {
		// Todo: Use workers
		go func(branchName string) {
			handleErrs(p.MergeBranch(repoURL, token, branchName, assignee)...)
		}(branch.Name)
	}
}

func (p *PushEventHandler) MergeBranch(repoURL, token, branch, assignee string) []error {
	// Ignore the default branch
	if branch == p.Config.Branch {
		return nil
	}
	var prURL = URL(repoURL, "pulls")
	var req = &PullRequest{
		Head:  p.Config.Branch,
		Base:  branch,
		Title: fmt.Sprintf("Automerge %s into %s", p.Config.Branch, branch),
	}
	var prRes PullRequestResponse
	var errs []error
	errs = GithubRequest("POST", prURL, token, http.StatusCreated, req, &prRes)
	if errs != nil {
		return errs
	}
	var mrReq = MergeRequest{
		MergeMethod: p.Config.MergeMethod,
	}
	errs = GithubRequest("PUT", URL(prURL, prRes.Number, "merge"), token, http.StatusOK, mrReq, nil)
	if errs == nil {
		return nil
	}
	var assignees Assignees
	if len(p.Config.Assignees) > 0 {
		assignees = p.Config.Assignees
	} else {
		assignees = Assignees{assignee}
	}

	data, err := json.Marshal(&PullRequestAssignees{Assignees: assignees})
	if err != nil {
		return append(errs, err)
	}

	errs2 := GithubRequest("POST", URL(repoURL, "issues", prRes.Number, "assignees"), token, http.StatusCreated, data, nil)
	if errs2 != nil {
		return append(errs, errs2...)
	}

	return errs
}

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}

func verifySignature(secret []byte, signature string, body []byte) error {

	const signaturePrefix = "sha1="
	const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return ErrBadSignature
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))

	if !hmac.Equal(signBody(secret, body), actual) {
		return ErrBadSignature
	}
	return nil
}
