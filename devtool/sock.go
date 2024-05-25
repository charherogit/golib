package devtool

import (
	"encoding/csv"
	"golib/config"
	log "golib/logutil"

	"fmt"
	"io"
	"net"
	"path/filepath"
	"syscall"
)

func ListenSock() error {
	pth := "." + config.GetSessionName() + ".sock"
	baseDir := config.GetEnv("SOCK_PATH", "")
	if baseDir != "" {
		pth = filepath.Join(baseDir, config.GetSessionName()+".sock")
	}
	_ = syscall.Unlink(pth)
	l, err := net.Listen("unix", pth)
	if err != nil {
		return err
	}
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}
			go handle(c)
		}
	}()
	return nil
}

func handle(c net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			_, _ = fmt.Fprintln(c, r)
		}
		_ = c.Close()
	}()

	r := csv.NewReader(c)
	r.Comma = ' '
	r.Comment = '#'
	r.TrimLeadingSpace = true
	r.LazyQuotes = true
	r.FieldsPerRecord = -1

	for {
		if !write(c, "> ") {
			return
		}
		args, err := r.Read()
		if err != nil {
			if _, err = c.Write([]byte(err.Error())); err != nil {
				log.Error(err)
			}
			return
		}
		if len(args) == 0 {
			continue
		}
		switch args[0] {
		case "help":
			help(c)
		case "exit", "q":
			return
		default:
			if !doCommand(c, args) {
				return
			}
		}
	}
}

func writeln(w io.Writer, a ...any) bool {
	_, err := fmt.Fprintln(w, a...)
	if err != nil {
		log.Error(err)
	}
	return err == nil
}

func write(w io.Writer, a ...any) bool {
	_, err := fmt.Fprint(w, a...)
	if err != nil {
		log.Error(err)
	}
	return err == nil
}
