package errdefs

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"runtime"
	"strings"
)

var (
	ErrInternal         = Status("internal_server_error", 500)
	ErrDatabase         = Status("database_error", 500)
	ErrPermissionDenied = Status("permission_denied", 403)
	ErrInvalidArgument  = Status("invalid_argument", 400)
	ErrAlreadyExists    = Status("already_exists", 409)
	ErrUnauthenticated  = Status("unauthenticated", 401)
	ErrResend           = Status("resend_error", 500)
	ErrUserNotFound     = Status("user_not_found", 404)
)

type Meta []any

type Error struct {
	ID     string         `json:"id"`
	Status int            `json:"status"`
	Title  string         `json:"title"`
	Detail string         `json:"detail"`
	Meta   map[string]any `json:"meta"`
	Frames stackTrace     `json:"-"`
}

type StatusFunc func(error, ...any) error

func Status(title string, status int) StatusFunc {
	return func(err error, vals ...any) error {
		e := &Error{
			ID:     errID(),
			Status: status,
			Title:  title,
			Detail: err.Error(),
			Meta:   make(map[string]any),
			Frames: newFrame(callers()),
		}

		for _, any := range vals {
			switch any := any.(type) {
			case Meta:
				e.Meta = appendMeta(e.Meta, any...)
			}
		}

		return e
	}
}

func (e *Error) Error() string {
	if e.Detail == "" {
		return e.Title
	}

	return e.Detail
}

func (e *Error) StackTrace() []string {
	if len(e.Frames) == 0 {
		return nil
	}
	var stack []string
	for _, frame := range e.Frames {
		stack = append(stack, frame.String())
	}
	return stack
}

func appendMeta(meta map[string]any, keyvals ...any) map[string]any {
	if meta == nil {
		meta = make(map[string]any)
	}
	var k string
	for n, v := range keyvals {
		if n%2 == 0 {
			k = fmt.Sprint(v)
		} else {
			meta[k] = v
		}
	}
	return meta
}

type frame struct {
	file           string
	lineNumber     int
	name           string
	programCounter uintptr
}

type stackTrace []*frame

func newFrame(pcs []uintptr) stackTrace {
	frames := []*frame{}

	for _, pc := range pcs {
		frame := &frame{programCounter: pc}
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			return frames
		}
		frame.name = trimPkgName(fn)

		frame.file, frame.lineNumber = fn.FileLine(pc - 1)
		frames = append(frames, frame)
	}

	return frames
}

func (f *frame) String() string {
	return fmt.Sprintf("%s:%d %s", f.file, f.lineNumber, f.name)
}

func trimPkgName(fn *runtime.Func) string {
	name := fn.Name()
	if ld := strings.LastIndex(name, "."); ld >= 0 {
		name = name[ld+1:]
	}

	return name
}

func callers() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])

	return pcs[0:n]
}

func errID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)

	return base64.StdEncoding.EncodeToString(b)
}

func IsUserNotFound(err error) bool {
	val, ok := err.(*Error)
	if !ok {
		return false
	}
	return val.Title == "user_not_found"
}
