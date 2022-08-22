package main

import (
	"flag"
	"log"
	"os"
	"path"
)

var (
	config = flag.String("config", "config.json", "configuration file containing job declarations")
	outDir = flag.String("outdir", "out", "output directory for jobs")
)

func main() {
	log.SetOutput(os.Stdout)
	flag.Parse()
	if err := run(); err != nil {
		log.Printf("error: %v\n", err)
	}
}

// run contains the core logic of crudeci
func run() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	outDir := path.Join(cwd, *outDir)
	cfg := path.Join(cwd, *config)

	log.Printf("[INFO] Parsing jobs...")
	jobs, err := ParseJobs(outDir, cfg)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		err := InitJob(job)
		if err != nil {
			return err
		}
	}

	runSequentially(jobs)

	successTable := successTable()
	for _, job := range jobs {
		log.Printf("Job %s - %v\n", job.Name, successTable[job.Successful])
	}

	return nil
}

// runSequentially runs the jobs in a sequential manner.
func runSequentially(jobs []*Job) {
	for _, job := range jobs {
		log.Printf("[INFO] Running pipeline for %s", job.Name)
		for _, step := range job.Pipeline.Steps {
			// log.Printf("TASK IDX %d\n", idx)
			out, err := step.Run()
			if err != nil {
				log.Printf("> %s step\n❌ %s\n%s", step.(*Step).Name, out, err)
				job.Successful = false
				break
			}
			job.Successful = true
			log.Printf("> %s step\n✅ %s", step.(*Step).Name, out)
		}
	}
}

// runConcurrently runs the jobs in a concurrent manner. Runs an anonymous
// goroutine until a SIGINT or SIGTERM signal are sent, at which point the
// quitChannel channel unblocks, allowing the function to return.
//
// func runConcurrently(jobs []Job, intervalInSeconds int) {
// 	go func() {
// 		for {
// 			for _, job := range jobs {
// 				log.Printf("[INFO] Running pipeline for %s", job.Name)
// 				for _, step := range job.pipeline.Steps {
// 					out, err := step.Run()
// 					if err != nil {
// 						log.Printf("> %s step\n❌ %s\n%s", step.(*Step).Name, out, err)
// 						job.successful = false
// 						break
// 					}
// 					job.successful = true
// 					log.Printf("> %s step\n✅ %s", step.(*Step).Name, out)
// 				}
// 				log.Println("Job status:", job)
// 			}
// 			time.Sleep(time.Second * time.Duration(intervalInSeconds))
// 		}
// 	}()
// 	quitChannel := make(chan os.Signal, 1)
// 	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
// 	<-quitChannel
// 	log.Println("goodbye!")
// }

// Exists returns true if file/directory passed exist; false otherwise
func Exists(fn string) (bool, error) {
	_, err := os.Open(fn)
	if err != nil {
		return false, err
	}

	return true, nil
}

func successTable() map[bool]string {
	successTable := make(map[bool]string)
	successTable[true] = "PASS"
	successTable[false] = "FAIL"
	return successTable
}
