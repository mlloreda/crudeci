package main

import (
	"path"
	"strings"
)

type Pipeline struct {
	Steps []Runner
}

// BuildPipeline constructs and returns a pipeline object from the embedded
// `step`s within the `job` passed in as an argument
func BuildPipeline(job *Job) (*Pipeline, error) {

	// ugly \HACK to iterate map in order :/
	var keys []string
	job.Steps = make(map[string]string)
	if len(job.Build) != 0 {
		key := "build"
		job.Steps[key] = job.Build
		keys = append(keys, key)
	}
	if len(job.Test) != 0 {
		key := "test"
		job.Steps[key] = job.Test
		keys = append(keys, key)
	}

	steps := []Runner{}
	for _, k := range keys {
		cmd := job.Steps[k]
		cmdSlice := strings.Split(cmd, " ")
		step := &Step{
			Name: k,
			Bin:  cmdSlice[0],
			Args: cmdSlice[1:],
			wd:   path.Join(job.Wd, job.WorkingDir),
		}

		steps = append(steps, step)
	}

	return &Pipeline{Steps: steps}, nil
}
