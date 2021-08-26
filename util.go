package dola

import (
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"go.uber.org/multierr"
)

// Location attempts to write the name of the caller function's parent.
// This occurs when the pointer pc is set to 1 and when the compiler is queried for the function's name.
// The pointer's data type is set to the data type of the function that is currently being executed.
// The compiler is then queried to get the function's pointer. If it succeeds, the code then performs a location and completes the phrase
// If it cannot locate the function's pointer, it returns a question mark to indicate that it is unknown
func Location() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	xs := strings.SplitAfterN(fn.Name(), "/", 3) // nolint: gomnd

	return xs[len(xs)-1]
}

// Location2 implements the grandparent call interface
// and contains Call 'street calling' troubleshooting  and returns the name of the grandparent function
func Location2() string {
	pc, _, _, ok := runtime.Caller(2) // nolint:gomnd
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	xs := strings.SplitAfterN(fn.Name(), "/", 3) // nolint: gomnd

	return xs[len(xs)-1]
}

// +---------+
// | Checker |
// +---------+

// Checker is a simple tool to check if everything initialized is
// subsequently deinitialized.  Works from simple open/close calls to
// gourintes.
type resourceChecker struct {
	m         sync.Mutex
	resources map[string]int
}

// nolint:gochecknoglobals
var defaultResourceChecker = resourceChecker{
	m:         sync.Mutex{},
	resources: make(map[string]int),
}

func CheckerPush(xs ...string) {
	var name string

	switch len(xs) {
	case 0:
		name = Location2()
	case 1:
		name = xs[0]
	default:
		panic("invalid argument")
	}

	defaultResourceChecker.m.Lock()
	defaultResourceChecker.resources[name]++
	defaultResourceChecker.m.Unlock()
}

func CheckerPop(xs ...string) {
	var name string

	switch len(xs) {
	case 0:
		name = Location2()
	case 1:
		name = xs[0]
	default:
		panic("invalid argument")
	}

	defaultResourceChecker.m.Lock()
	defaultResourceChecker.resources[name]--
	defaultResourceChecker.m.Unlock()
}

// CheckerAssert should be defer-called in main().
func CheckerAssert() {
	What(log.Debug(), "checking resources...")
	time.Sleep(1 * time.Second)

	defaultResourceChecker.m.Lock()
	defer defaultResourceChecker.m.Unlock()

	for k, v := range defaultResourceChecker.resources {
		if v != 0 {
			What(log.Warn().Int("counter", v).Str("unit", k), "leaked resource")
		}
	}
}

// +----------------+
// | ErrorWaitGroup |
// +----------------+

type ErrorWaitGroup struct {
	err   error
	group sync.WaitGroup
	mutex sync.Mutex
}

func (m *ErrorWaitGroup) Add(delta int) {
	m.group.Add(delta)
}

func (m *ErrorWaitGroup) Done(right error) {
	m.mutex.Lock()
	m.err = multierr.Append(m.err, right)
	m.mutex.Unlock()
	m.group.Done()
}

func (m *ErrorWaitGroup) Wait() error {
	m.group.Wait()

	return m.err
}

// +----------+
// | Profiler |
// +----------+

type Profiler struct {
	Filename string
}

func NewProfiler(filename string) Profiler {
	p := Profiler{filename}

	f, err := os.Create(p.Filename)
	if err != nil {
		panic(err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}

	return p
}

func (p Profiler) Stop() {
	pprof.StopCPUProfile()
}

// +---------------+
// | Miscellaneous |
// +---------------+

// ConfigFile there is text in the `inp` string variable, if there is not, `inp` is set to an empty string.
// If `inp` is not empty, check to see if `inp` corresponds with multiple strings and whether `inp`'s corresponding string is a folder.
// The function is only necessary if `inp` is not the default configuration file.
// This utility function that is used to check for files with the name `inp` or `'inp`. These two strings are being output as such because of the file in the default configuration file, the `~/.dola/config.json` file. If `inp` does not correspond with the `inp` filename, then `inp` is set to `""` and is returned as such.
// If environment variable "DOLA_CONFIG" is found, it attempts to expand the string using ExpandUser(), also leading to an error if the string is empty.
// Also check if configuration file exists based on a user-operated environment variable, the executable running, and the user's home directory. If any of these conditions are met, the configuration file will be found. If none of the conditions are met, a configuration file will not be found
func ConfigFile(inp string) string {
	if inp != "" {
		path := ExpandUser(inp)
		if FileExists(path) {
			return path
		}
	}

	if env := os.Getenv("DOLA_CONFIG"); env != "" {
		path := ExpandUser(env)
		if FileExists(path) {
			return path
		}
	}

	if path := ExpandUser("~/.dola/config.json"); FileExists(path) {
		return path
	}

	return ""
}

// ExpandUser can take in a string and a string, and returns a string if the string paths are the same. This function also expands the tilde (~) to the current user's home directory.
// The userPath variable holds the result of all the string replacements, and then the result is returned.
// This function can be used to set files in one directory as the home for all changes that occur in directories the user has permissions to.
// The ~ is replaced with the current username's home directory and expanded in order to be a valid path, it is considered a "shortcut" for the current user's home directory.
// This file can then put to use with the other functions, to print the document's user directory and change the directory of the working directory.
func ExpandUser(path string) string {
	// Get user's home directory
	home := os.Getenv("HOME")
	// Expands ~
	var userPath = os.ExpandEnv(strings.Replace(path, "~", home, 1))
	// returns user's home path and home directory
	return userPath
}

// FileExists returns a boolean indicating whether the path exist. Path is the parameter for the filename to check
// os.Stat is a filesystem`s function which returns a code indicating whether the name is that of an existing file
// The return statement closing the function `FileExists` and informing that the File does not exist.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
