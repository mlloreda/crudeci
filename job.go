package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Job struct {
	Name       string `json:"name"`
	Repo       string `json:"repo"`
	WorkingDir string `json:"workingDir"`
	Build      string `json:"build"` // TODO consolidate
	Test       string `json:"test"`  // TODO consolidate
	Successful bool
	Pipeline   Pipeline
	Steps      map[string]string
	Wd         string
	// contact string // TODO
}

// ParseJobs unmarshals the configuration file declaring job definitions, and
// populates the `Job` struct along with its embedded `Pipeline` struct
func ParseJobs(outDir, config string) ([]*Job, error) {
	file, err := os.Open(config)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var jobs []*Job
	err = json.NewDecoder(file).Decode(&jobs)
	if err != nil {
		return nil, fmt.Errorf("unable to decode jobs into struct: %v", err)
	}

	for _, job := range jobs {
		job.Wd = path.Join(outDir, strings.ToLower(job.Name), path.Base(job.Repo))
		pipeline, err := BuildPipeline(job)
		if err != nil {
			return nil, err
		}
		job.Pipeline = *pipeline
	}

	return jobs, nil
}

// InitJob initializes the job by cloning the source repository if necessary.
func InitJob(job *Job) error {
	if repoExists, _ := Exists(job.Wd); !repoExists {
		err := CloneRepo(job.Repo, job.Wd)
		if err != nil {
			return err
		}
	}

	return nil
}

// CloneRepo clones the repository passed, `repo`, in to the working directory,
// `wd`, in an idempotent manner.
func CloneRepo(repo, wd string) error {
	cmd := exec.Command("git", "clone", repo)
	cmd.Dir = path.Dir(wd)

	if dirExists, _ := Exists(cmd.Dir); !dirExists {
		if err := os.MkdirAll(cmd.Dir, 0777); err != nil {
			return err
		}
	}

	if repoDirExists, _ := Exists(wd); !repoDirExists {
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func (j Job) String() string {
	out := ""
	out += fmt.Sprintf("\n\tname: %s\n", j.Name)
	out += fmt.Sprintf("\trepo: %s\n", j.Repo)
	out += fmt.Sprintf("\twd: %s\n", j.Wd)
	out += fmt.Sprintf("\tbuild step: %v\n", j.Build)
	out += fmt.Sprintf("\ttest step: %v\n", j.Test)
	out += fmt.Sprintf("\tsuccess: %v\n", j.Successful)
	return out
}
