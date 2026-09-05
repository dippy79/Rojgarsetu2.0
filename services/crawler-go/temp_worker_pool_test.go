package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rojgarsetu/crawler/internal"
	"github.com/rojgarsetu/crawler/internal/sources"
)

type fakeSource struct{ name string }

func (f fakeSource) FetchJobs() ([]sources.Job, error) {
	time.Sleep(350 * time.Millisecond)
	return []sources.Job{{Title: f.name, ApplyURL: "https://example.com/" + f.name}}, nil
}

func main() {
	os.Setenv("MAX_WORKERS", "2")
	jobSources := []internal.JobSource{
		fakeSource{name: "A"},
		fakeSource{name: "B"},
		fakeSource{name: "C"},
		fakeSource{name: "D"},
		fakeSource{name: "E"},
	}

	var active int
	var maxActive int
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 2)

	for _, src := range jobSources {
		wg.Add(1)
		go func(src internal.JobSource) {
			defer wg.Done()
			sem <- struct{}{}
			mu.Lock()
			active++
			if active > maxActive {
				maxActive = active
			}
			mu.Unlock()
			defer func() {
				mu.Lock()
				active--
				mu.Unlock()
				<-sem
			}()
			fmt.Printf("START %s active=%d\n", src.(fakeSource).name, active)
			time.Sleep(350 * time.Millisecond)
			fmt.Printf("END %s active=%d\n", src.(fakeSource).name, active)
		}(src)
	}

	wg.Wait()
	fmt.Printf("max_active=%d\n", maxActive)
	if maxActive != 2 {
		panic(fmt.Sprintf("expected max_active=2 got %d", maxActive))
	}
	fmt.Println("WORKER_POOL_CHECK=PASS")
}
