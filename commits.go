//
// Copyright 2017, Sander van Harmelen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package gitlab

import (
	"fmt"
	"net/url"
	"time"
)

// CommitsService handles communication with the commit related methods
// of the GitLab API.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html
type CommitsService struct {
	client *Client
}

// Commit represents a GitLab commit.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html
type Commit struct {
	ID             string           `bson:"id" json:"id"`
	ShortID        string           `bson:"short_id" json:"short_id"`
	Title          string           `bson:"title" json:"title"`
	AuthorName     string           `bson:"author_name" json:"author_name"`
	AuthorEmail    string           `bson:"author_email" json:"author_email"`
	AuthoredDate   *time.Time       `bson:"authored_date" json:"authored_date"`
	CommitterName  string           `bson:"committer_name" json:"committer_name"`
	CommitterEmail string           `bson:"committer_email" json:"committer_email"`
	CommittedDate  *time.Time       `bson:"committed_date" json:"committed_date"`
	CreatedAt      *time.Time       `bson:"created_at" json:"created_at"`
	Message        string           `bson:"message" json:"message"`
	ParentIDs      []string         `bson:"parent_ids" json:"parent_ids"`
	Stats          *CommitStats     `bson:"stats" json:"stats"`
	Status         *BuildStateValue `bson:"status" json:"status"`
}

// CommitStats represents the number of added and deleted files in a commit.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html
type CommitStats struct {
	Additions int `bson:"additions" json:"additions"`
	Deletions int `bson:"deletions" json:"deletions"`
	Total     int `bson:"total" json:"total"`
}

func (c Commit) String() string {
	return Stringify(c)
}

// ListCommitsOptions represents the available ListCommits() options.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html#list-repository-commits
type ListCommitsOptions struct {
	ListOptions
	RefName   *string    `url:"ref_name,omitempty" bson:"ref_name,omitempty" json:"ref_name,omitempty"`
	Since     *time.Time `url:"since,omitempty" bson:"since,omitempty" json:"since,omitempty"`
	Until     *time.Time `url:"until,omitempty" bson:"until,omitempty" json:"until,omitempty"`
	Path      *string    `url:"path,omitempty" bson:"path,omitempty" json:"path,omitempty"`
	All       *bool      `url:"all,omitempty" bson:"all,omitempty" json:"all,omitempty"`
	WithStats *bool      `url:"with_stats,omitempty" bson:"with_stats,omitempty" json:"with_stats,omitempty"`
}

// ListCommits gets a list of repository commits in a project.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html#list-commits
func (s *CommitsService) ListCommits(pid interface{}, opt *ListCommitsOptions, options ...OptionFunc) ([]*Commit, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/commits", url.QueryEscape(project))

	req, err := s.client.NewRequest("GET", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var c []*Commit
	resp, err := s.client.Do(req, &c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, err
}

// FileAction represents the available actions that can be performed on a file.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html#create-a-commit-with-multiple-files-and-actions
type FileAction string

// The available file actions.
const (
	FileCreate FileAction = "create"
	FileDelete FileAction = "delete"
	FileMove   FileAction = "move"
	FileUpdate FileAction = "update"
)

// CommitAction represents a single file action within a commit.
type CommitAction struct {
	Action       FileAction `url:"action" bson:"action" json:"action"`
	FilePath     string     `url:"file_path" bson:"file_path" json:"file_path"`
	PreviousPath string     `url:"previous_path,omitempty" bson:"previous_path,omitempty" json:"previous_path,omitempty"`
	Content      string     `url:"content,omitempty" bson:"content,omitempty" json:"content,omitempty"`
	Encoding     string     `url:"encoding,omitempty" bson:"encoding,omitempty" json:"encoding,omitempty"`
}

// CommitRef represents the reference of branches/tags in a commit.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/commits.html#get-references-a-commit-is-pushed-to
type CommitRef struct {
	Type string `bson:"type" json:"type"`
	Name string `bson:"name" json:"name"`
}

// GetCommitRefsOptions represents the available GetCommitRefs() options.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/commits.html#get-references-a-commit-is-pushed-to
type GetCommitRefsOptions struct {
	ListOptions
	Type *string `url:"type,omitempty" bson:"type,omitempty" json:"type,omitempty"`
}

// GetCommitRefs gets all references (from branches or tags) a commit is pushed to
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/commits.html#get-references-a-commit-is-pushed-to
func (s *CommitsService) GetCommitRefs(pid interface{}, sha string, opt *GetCommitRefsOptions, options ...OptionFunc) ([]CommitRef, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/commits/%s/refs", url.QueryEscape(project), sha)

	req, err := s.client.NewRequest("GET", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var cs []CommitRef
	resp, err := s.client.Do(req, &cs)
	if err != nil {
		return nil, resp, err
	}

	return cs, resp, err
}

// GetCommit gets a specific commit identified by the commit hash or name of a
// branch or tag.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html#get-a-single-commit
func (s *CommitsService) GetCommit(pid interface{}, sha string, options ...OptionFunc) (*Commit, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/commits/%s", url.QueryEscape(project), sha)

	req, err := s.client.NewRequest("GET", u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	c := new(Commit)
	resp, err := s.client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, err
}

// CreateCommitOptions represents the available options for a new commit.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html#create-a-commit-with-multiple-files-and-actions
type CreateCommitOptions struct {
	Branch        *string         `url:"branch" bson:"branch" json:"branch"`
	CommitMessage *string         `url:"commit_message" bson:"commit_message" json:"commit_message"`
	StartBranch   *string         `url:"start_branch,omitempty" bson:"start_branch,omitempty" json:"start_branch,omitempty"`
	Actions       []*CommitAction `url:"actions" bson:"actions" json:"actions"`
	AuthorEmail   *string         `url:"author_email,omitempty" bson:"author_email,omitempty" json:"author_email,omitempty"`
	AuthorName    *string         `url:"author_name,omitempty" bson:"author_name,omitempty" json:"author_name,omitempty"`
}

// CreateCommit creates a commit with multiple files and actions.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html#create-a-commit-with-multiple-files-and-actions
func (s *CommitsService) CreateCommit(pid interface{}, opt *CreateCommitOptions, options ...OptionFunc) (*Commit, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/commits", url.QueryEscape(project))

	req, err := s.client.NewRequest("POST", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var c *Commit
	resp, err := s.client.Do(req, &c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, err
}

// Diff represents a GitLab diff.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html
type Diff struct {
	Diff        string `bson:"diff" json:"diff"`
	NewPath     string `bson:"new_path" json:"new_path"`
	OldPath     string `bson:"old_path" json:"old_path"`
	AMode       string `bson:"a_mode" json:"a_mode"`
	BMode       string `bson:"b_mode" json:"b_mode"`
	NewFile     bool   `bson:"new_file" json:"new_file"`
	RenamedFile bool   `bson:"renamed_file" json:"renamed_file"`
	DeletedFile bool   `bson:"deleted_file" json:"deleted_file"`
}

func (d Diff) String() string {
	return Stringify(d)
}

// GetCommitDiffOptions represents the available GetCommitDiff() options.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/commits.html#get-the-diff-of-a-commit
type GetCommitDiffOptions ListOptions

// GetCommitDiff gets the diff of a commit in a project..
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/commits.html#get-the-diff-of-a-commit
func (s *CommitsService) GetCommitDiff(pid interface{}, sha string, opt *GetCommitDiffOptions, options ...OptionFunc) ([]*Diff, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/commits/%s/diff", url.QueryEscape(project), sha)

	req, err := s.client.NewRequest("GET", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var d []*Diff
	resp, err := s.client.Do(req, &d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, err
}

// CommitComment represents a GitLab commit comment.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html
type CommitComment struct {
	Note     string `bson:"note" json:"note"`
	Path     string `bson:"path" json:"path"`
	Line     int    `bson:"line" json:"line"`
	LineType string `bson:"line_type" json:"line_type"`
	Author   Author `bson:"author" json:"author"`
}

// Author represents a GitLab commit author
type Author struct {
	ID        int        `bson:"id" json:"id"`
	Username  string     `bson:"username" json:"username"`
	Email     string     `bson:"email" json:"email"`
	Name      string     `bson:"name" json:"name"`
	State     string     `bson:"state" json:"state"`
	Blocked   bool       `bson:"blocked" json:"blocked"`
	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
}

func (c CommitComment) String() string {
	return Stringify(c)
}

// GetCommitCommentsOptions represents the available GetCommitComments() options.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/commits.html#get-the-comments-of-a-commit
type GetCommitCommentsOptions ListOptions

// GetCommitComments gets the comments of a commit in a project.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/commits.html#get-the-comments-of-a-commit
func (s *CommitsService) GetCommitComments(pid interface{}, sha string, opt *GetCommitCommentsOptions, options ...OptionFunc) ([]*CommitComment, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/commits/%s/comments", url.QueryEscape(project), sha)

	req, err := s.client.NewRequest("GET", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var c []*CommitComment
	resp, err := s.client.Do(req, &c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, err
}

// PostCommitCommentOptions represents the available PostCommitComment()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/commits.html#post-comment-to-commit
type PostCommitCommentOptions struct {
	Note     *string `url:"note,omitempty" bson:"note,omitempty" json:"note,omitempty"`
	Path     *string `url:"path" bson:"path" json:"path"`
	Line     *int    `url:"line" bson:"line" json:"line"`
	LineType *string `url:"line_type" bson:"line_type" json:"line_type"`
}

// PostCommitComment adds a comment to a commit. Optionally you can post
// comments on a specific line of a commit. Therefor both path, line_new and
// line_old are required.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/commits.html#post-comment-to-commit
func (s *CommitsService) PostCommitComment(pid interface{}, sha string, opt *PostCommitCommentOptions, options ...OptionFunc) (*CommitComment, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/commits/%s/comments", url.QueryEscape(project), sha)

	req, err := s.client.NewRequest("POST", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	c := new(CommitComment)
	resp, err := s.client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, err
}

// GetCommitStatusesOptions represents the available GetCommitStatuses() options.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html#get-the-status-of-a-commit
type GetCommitStatusesOptions struct {
	ListOptions
	Ref   *string `url:"ref,omitempty" bson:"ref,omitempty" json:"ref,omitempty"`
	Stage *string `url:"stage,omitempty" bson:"stage,omitempty" json:"stage,omitempty"`
	Name  *string `url:"name,omitempty" bson:"name,omitempty" json:"name,omitempty"`
	All   *bool   `url:"all,omitempty" bson:"all,omitempty" json:"all,omitempty"`
}

// CommitStatus represents a GitLab commit status.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html#get-the-status-of-a-commit
type CommitStatus struct {
	ID          int        `bson:"id" json:"id"`
	SHA         string     `bson:"sha" json:"sha"`
	Ref         string     `bson:"ref" json:"ref"`
	Status      string     `bson:"status" json:"status"`
	Name        string     `bson:"name" json:"name"`
	TargetURL   string     `bson:"target_url" json:"target_url"`
	Description string     `bson:"description" json:"description"`
	CreatedAt   *time.Time `bson:"created_at" json:"created_at"`
	StartedAt   *time.Time `bson:"started_at" json:"started_at"`
	FinishedAt  *time.Time `bson:"finished_at" json:"finished_at"`
	Author      Author     `bson:"author" json:"author"`
}

// GetCommitStatuses gets the statuses of a commit in a project.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html#get-the-status-of-a-commit
func (s *CommitsService) GetCommitStatuses(pid interface{}, sha string, opt *GetCommitStatusesOptions, options ...OptionFunc) ([]*CommitStatus, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/commits/%s/statuses", url.QueryEscape(project), sha)

	req, err := s.client.NewRequest("GET", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var cs []*CommitStatus
	resp, err := s.client.Do(req, &cs)
	if err != nil {
		return nil, resp, err
	}

	return cs, resp, err
}

// SetCommitStatusOptions represents the available SetCommitStatus() options.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html#post-the-status-to-commit
type SetCommitStatusOptions struct {
	State       BuildStateValue `url:"state" bson:"state" json:"state"`
	Ref         *string         `url:"ref,omitempty" bson:"ref,omitempty" json:"ref,omitempty"`
	Name        *string         `url:"name,omitempty" bson:"name,omitempty" json:"name,omitempty"`
	Context     *string         `url:"context,omitempty" bson:"context,omitempty" json:"context,omitempty"`
	TargetURL   *string         `url:"target_url,omitempty" bson:"target_url,omitempty" json:"target_url,omitempty"`
	Description *string         `url:"description,omitempty" bson:"description,omitempty" json:"description,omitempty"`
}

// SetCommitStatus sets the status of a commit in a project.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html#post-the-status-to-commit
func (s *CommitsService) SetCommitStatus(pid interface{}, sha string, opt *SetCommitStatusOptions, options ...OptionFunc) (*CommitStatus, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/statuses/%s", url.QueryEscape(project), sha)

	req, err := s.client.NewRequest("POST", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var cs *CommitStatus
	resp, err := s.client.Do(req, &cs)
	if err != nil {
		return nil, resp, err
	}

	return cs, resp, err
}

// GetMergeRequestsByCommit gets merge request associated with a commit.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/commits.html#list-merge-requests-associated-with-a-commit
func (s *CommitsService) GetMergeRequestsByCommit(pid interface{}, sha string, options ...OptionFunc) ([]*MergeRequest, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/commits/%s/merge_requests",
		url.QueryEscape(project), url.QueryEscape(sha))

	req, err := s.client.NewRequest("GET", u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var mrs []*MergeRequest
	resp, err := s.client.Do(req, &mrs)
	if err != nil {
		return nil, resp, err
	}

	return mrs, resp, err
}

// CherryPickCommitOptions represents the available options for cherry-picking a commit.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html#cherry-pick-a-commit
type CherryPickCommitOptions struct {
	TargetBranch *string `url:"branch" bson:"branch,omitempty" json:"branch,omitempty"`
}

// CherryPickCommit sherry picks a commit to a given branch.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/commits.html#cherry-pick-a-commit
func (s *CommitsService) CherryPickCommit(pid interface{}, sha string, opt *CherryPickCommitOptions, options ...OptionFunc) (*Commit, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/commits/%s/cherry_pick",
		url.QueryEscape(project), url.QueryEscape(sha))

	req, err := s.client.NewRequest("POST", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var c *Commit
	resp, err := s.client.Do(req, &c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, err
}
