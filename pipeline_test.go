package main_test

import (
	"strings"
	"testing"

	crudeci "github.com/mlloreda/crudeci"
)

func TestBuildPipeline(t *testing.T) {
	cases := []struct {
		name       string
		build      string
		test       string
		expLen     int
		expSuccess bool
		// steps map[string]string
	}{
		{
			name:       "success case",
			build:      "pwd",
			test:       "ls -lah",
			expLen:     2,
			expSuccess: true,
			// steps: map[string]string{"build": "pwd", "test": "ls -lah"},
		},
		{
			name:       "failure case (# of steps mismatch)",
			build:      "pwd",
			expLen:     2,
			expSuccess: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			j := testCreateJob(t, tc.build, tc.test)
			pipeline, err := crudeci.BuildPipeline(j)
			if err != nil {
				t.Fatalf("%v", err)
			}
			if tc.expLen != len(pipeline.Steps) && tc.expSuccess == true {
				t.Fatalf("expected len %d, got len %d", tc.expLen, len(pipeline.Steps))
			}
			for _, step := range pipeline.Steps {
				s := step.(*crudeci.Step)
				cmd := strings.TrimSpace(s.Bin + " " + strings.Join(s.Args, " "))

				if s.Name == "build" {
					if cmd != tc.build {
						t.Fatalf("expected %s, got %s", tc.build, cmd)
					}
				} else if s.Name == "test" {
					if cmd != tc.test {
						t.Fatalf("expected %s, got %s", tc.test, cmd)
					}
				}
			}
		})
	}
}

func testCreateJob(t *testing.T, build, test string) *crudeci.Job {
	t.Helper()

	return &crudeci.Job{
		Name:       "testJob",
		Repo:       "https://github.com/mlloreda/dotfiles",
		WorkingDir: ".",
		Build:      build,
		Test:       test,
	}
}
