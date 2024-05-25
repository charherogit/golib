package devtool

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4/debug/profile"
	"golib/config"
	log "golib/logutil"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func init() {
	AddCommand(defaultCommands()())
}

var cmds = make(map[string]any)

func AddCommand(m map[string]any) {
	for k, v := range m {
		set(cmds, k, v)
	}
}

func set(dst map[string]any, name string, v any) {
	vk, ok := v.(map[string]any)
	if !ok {
		dst[name] = v
		return
	}
	vv, ok := dst[name]
	if !ok {
		vv = make(map[string]any)
		dst[name] = vv
	}
	vvk := vv.(map[string]any)
	for k, v := range vk {
		set(vvk, k, v)
	}
}

type item struct {
	name string
	val  any
}

func help(w io.Writer) {
	for _, v := range toList(cmds) {
		inspectCmd(w, 0, v.name, v.val)
	}
}

func inspectCmd(w io.Writer, d int, name string, v any) {
	if vv, ok := v.(map[string]any); ok {
		_, _ = fmt.Fprintf(w, "%s%s\n", strings.Repeat(" ", d), name)
		for _, vvv := range toList(vv) {
			inspectCmd(w, d+1, vvv.name, vvv.val)
		}
	} else {
		i, _ := getSign(v)
		_, _ = fmt.Fprintf(w, "%s%s %v\n", strings.Repeat(" ", d), name, i)
	}
}

func toList(m map[string]any) (l []item) {
	for k, v := range m {
		l = append(l, item{
			name: k, val: v,
		})
	}
	sort.Slice(l, func(i, j int) bool {
		return l[i].name < l[j].name
	})
	return
}

func defaultCommands() func() map[string]any {
	var pprof profile.Profile
	return func() map[string]any {
		return map[string]any{
			"build":     config.GetBuildInfo,
			"pid":       os.Getpid,
			"forceExit": os.Exit,
			"log": map[string]any{
				"setLv": log.SetLevel,
				"lvs": func(w io.Writer) {
					writeln(w, logrus.AllLevels)
				},
				"query": func(w io.Writer) {
					writeln(w, log.GetLevel())
				},
			},
			"prof": map[string]any{
				"start": func(addr string) error {
					if pprof == nil {
						p := NewPProf(addr)
						if err := p.Start(); err != nil {
							return err
						}
						pprof = p
					} else {
						return fmt.Errorf("pprof is already running port: %s", pprof.(*farmProf).Port())
					}
					return nil
				},
				"stop": func() error {
					if pprof != nil {
						if err := pprof.Stop(); err != nil {
							return err
						}
						pprof = nil
					}
					return nil
				},
			},
			"env": func(w io.Writer) {
				for _, v := range os.Environ() {
					writeln(w, v)
				}
			},
		}
	}
}

func doCommand(w io.Writer, args []string) bool {
	cc := cmds
	for {
		if len(args) == 0 {
			return true
		}
		v, ok := cc[args[0]]
		if !ok {
			return writeln(w, "command not found", args)
		}
		args = args[1:]
		if vv, ok := v.(map[string]any); ok {
			cc = vv
		} else {
			return callFn(w, v, args)
		}
	}
}

func getSign(f any) ([]string, []string) {
	in := make([]string, 0)
	out := make([]string, 0)

	ft := reflect.TypeOf(f)
	for i := 0; i < ft.NumIn(); i++ {
		if ft.In(i).Implements(reflect.TypeOf((*io.Writer)(nil)).Elem()) {
			continue
		}
		in = append(in, ft.In(i).String())
	}
	for i := 0; i < ft.NumOut(); i++ {
		out = append(out, ft.Out(i).String())
	}
	return in, out
}

// TODO: 变长参数
func parseVal(w io.Writer, ft reflect.Type, args []string) ([]reflect.Value, error) {
	callArgs := make([]reflect.Value, 0, ft.NumIn())
	for i := 0; i < ft.NumIn(); i++ {
		if i == 0 && ft.In(i).Implements(reflect.TypeOf((*io.Writer)(nil)).Elem()) {
			callArgs = append(callArgs, reflect.ValueOf(w))
			args = append([]string{""}, args...)
			continue
		}
		v := reflect.New(ft.In(i)).Elem()
		if len(args) > i {
			if ok, err := SetBuiltInVal(v, args[i]); err != nil {
				return nil, err
			} else if !ok {
				return nil, fmt.Errorf("can't setFn %s to %s", args[i], ft.In(i).Kind())
			}
		}
		callArgs = append(callArgs, v)
	}
	return callArgs, nil
}

func SetBuiltInVal(v reflect.Value, s string) (bool, error) {
	if !v.CanSet() {
		return false, nil
	}
	switch v.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		val, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return false, err
		}
		v.SetUint(val)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return false, err
		}
		v.SetInt(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return false, err
		}
		v.SetFloat(val)
	case reflect.String:
		v.SetString(s)
	case reflect.Bool:
		val, err := strconv.ParseBool(s)
		if err != nil {
			return false, err
		}
		v.SetBool(val)
	default:
		return false, nil
	}
	return true, nil
}

// 返回值只能是error或者无
func callFn(w io.Writer, fn any, args []string) bool {
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		return true
	}
	ft := v.Type()

	callArgs, err := parseVal(w, ft, args)
	if err != nil {
		write(w, err)
		return true
	}
	res := v.Call(callArgs)
	ret := make([]any, 0)
	for _, v := range res {
		if v.IsValid() {
			ret = append(ret, v.Interface())
		}
	}
	if len(ret) != 0 {
		_ = writeln(w, ret...)
	}
	return true
}
