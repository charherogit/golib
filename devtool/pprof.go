package devtool

import (
	"context"
	"expvar"
	"fmt"
	"github.com/arl/statsviz"
	"github.com/google/gops/agent"
	"go-micro.dev/v4/debug/profile"
	"net/http"
	wp "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
)

type farmProf struct {
	sync.Mutex
	running bool
	server  *http.Server
}

func init() {
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)
}

func NewPProf(addr string) profile.Profile {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", wp.Index)
	mux.HandleFunc("/debug/pprof/cmdline", wp.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", wp.Profile)
	mux.HandleFunc("/debug/pprof/symbol", wp.Symbol)
	mux.HandleFunc("/debug/pprof/trace", wp.Trace)
	mux.Handle("/debug/vars", expvar.Handler())
	_ = statsviz.Register(mux)

	return &farmProf{
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (p *farmProf) Start() error {
	p.Lock()
	defer p.Unlock()

	if p.running {
		return nil
	}

	ch := make(chan error)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go func() {
		if err := p.server.ListenAndServe(); err != nil {
			ch <- err
			close(ch)
			p.Lock()
			p.running = false
			p.Unlock()
		}
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		p.running = true
		return nil
	}
}

func (p *farmProf) Stop() error {
	p.Lock()
	defer p.Unlock()

	if !p.running {
		return nil
	}

	p.running = false

	return p.server.Shutdown(context.Background())
}

func (p *farmProf) String() string {
	return "pkg.pprof"
}

func (p *farmProf) Port() string {
	if p == nil || p.server == nil {
		return ""
	}
	return strings.ReplaceAll(p.server.Addr, ":", "")
}

func InitGoPs() error {
	return agent.Listen(agent.Options{
		// Addr: "0.0.0.0:6060", // 外部可以用过 gops cmd 192.168.196.18:6060 来查看信息
		// ConfigDir:       "", // 默认 ~/.config/gops
	})
}

func CPUProfile(output string) func() {
	if f, err := os.Create(output); err != nil {
		return func() {
			_, _ = fmt.Fprintf(os.Stderr, "CPUProfile: %v\n", err)
		}
	} else if err = pprof.StartCPUProfile(f); err != nil {
		return func() {
			_, _ = fmt.Fprintf(os.Stderr, "CPUProfile: %v\n", err)
		}
	} else {
		return func() {
			pprof.StopCPUProfile()
			if err = f.Close(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "CPUProfile: %v\n", err)
			}
		}
	}
}

func MemProfile(output string) func() {
	f, err := os.Create(output)
	if err != nil {
		return func() {
			_, _ = fmt.Fprintf(os.Stderr, "MemProfile: %v\n", err)
		}
	}
	runtime.GC()
	if err := pprof.WriteHeapProfile(f); err != nil {
		return func() {
			_, _ = fmt.Fprintf(os.Stderr, "MemProfile: %v\n", err)
		}
	}
	return func() {
		if err := f.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "MemProfile: %v\n", err)
		}
	}
}
