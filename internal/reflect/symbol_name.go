package reflect

import (
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

var symbolCleanup = regexp.MustCompile(`[()*]`)

func SymbolName(in any) string {
	raw := readSymbolName(in)
	if idx := strings.LastIndex(raw, "/"); idx != -1 {
		raw = raw[idx+1:]
	}
	raw = strings.TrimSuffix(raw, "-fm")
	//	raw = strings.ToLower(raw)
	return symbolCleanup.ReplaceAllString(raw, "")
}

func readSymbolName(in any) string {
	if in == nil {
		return "<nil>"
	}

	// If it's a string, just return it
	if s, ok := in.(string); ok && s != "" {
		return s
	}

	v := reflect.ValueOf(in)
	t := v.Type()

	// If it's a function, get its runtime symbol name
	if t.Kind() == reflect.Func {
		if fn := runtime.FuncForPC(v.Pointer()); fn != nil {
			return fn.Name()
		}
	}

	// Otherwise, return the fully qualified type name
	if pkg := t.PkgPath(); pkg != "" {
		return pkg + "." + t.String()
	}

	return t.String()
}
