package dict

import (
	"strings"
	"time"
)

// PushEventPayload contains the information for GitLab's push event
type PushEventPayload struct {
	ObjectKind        string     `json:"object_kind"`
	Before            string     `json:"before"`
	After             string     `json:"after"`
	Ref               string     `json:"ref"`
	CheckoutSHA       string     `json:"checkout_sha"`
	UserID            int64      `json:"user_id"`
	UserName          string     `json:"user_name"`
	UserUsername      string     `json:"user_username"`
	UserEmail         string     `json:"user_email"`
	UserAvatar        string     `json:"user_avatar"`
	ProjectID         int64      `json:"project_id"`
	Project           Project    `json:"Project"`
	Repository        Repository `json:"repository"`
	Commits           []Commit   `json:"commits"`
	TotalCommitsCount int64      `json:"total_commits_count"`
}


// Commit contains all of the GitLab commit information
type Commit struct {
	ID        string     `json:"id"`
	Message   string     `json:"message"`
	Timestamp customTime `json:"timestamp"`
	URL       string     `json:"url"`
	Author    Author     `json:"author"`
	//Added     []string   `json:"added"`
	//Modified  []string   `json:"modified"`
	//Removed   []string   `json:"removed"`
}


// Project contains all of the GitLab project information
type Project struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	WebURL            string `json:"web_url"`
	AvatarURL         string `json:"avatar_url"`
	GitSSSHURL        string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	Namespace         string `json:"namespace"`
	VisibilityLevel   int64  `json:"visibility_level"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
	Homepage          string `json:"homepage"`
	URL               string `json:"url"`
	SSHURL            string `json:"ssh_url"`
	HTTPURL           string `json:"http_url"`
}

// Repository contains all of the GitLab repository information
type Repository struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}


// Author contains all of the GitLab author information
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}


type customTime struct {
	time.Time
}

func (t *customTime) UnmarshalJSON(b []byte) (err error) {
	layout := []string{
		"2006-01-02 15:04:05 MST",
		"2006-01-02 15:04:05 Z07:00",
		"2006-01-02 15:04:05 Z0700",
		time.RFC3339,
	}
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return
	}
	for _, l := range layout {
		t.Time, err = time.Parse(l, s)
		if err == nil {
			break
		}
	}
	return
}