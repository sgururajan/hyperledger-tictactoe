package pathUtil

import (
	"go/build"
	"path/filepath"
	"strings"
	"os"
	"bytes"
)

func goPath() string {
	gpDefault:= build.Default.GOPATH
	gpSplited:= filepath.SplitList(gpDefault)
	return gpSplited[0]
}

func lookupEnvVar(v string) (string, bool) {
	if strings.EqualFold(v, "GOPATH") {
		return goPath(), true
	}

	return os.LookupEnv(v)
}

func Substitute(path string) string{
	const(
		sepPrefix="${"
		sepSuffix="}"
	)

	splits:= strings.Split(path, sepPrefix)
	var buffer bytes.Buffer
	buffer.WriteString(splits[0])

	for _, s:= range splits[1:] {
		subst, rest:= substituteEnvVar(s, sepPrefix, sepSuffix)
		buffer.WriteString(subst)
		buffer.WriteString(rest)
	}

	return buffer.String()
}

func substituteEnvVar(s string, noMatch string, sep string) (string, string){
	endPos:= strings.Index(s, sep)
	if endPos== -1 {
		return  noMatch, s
	}

	v, ok:= lookupEnvVar(s[:endPos])
	if !ok {
		return  noMatch, s
	}

	return  v, s[endPos+1:]
}

func EnsureDirectory(path string, dirName string) {
	if _,err:= os.Stat(filepath.Join(path, dirName)); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(path, dirName), os.ModePerm)
	}
}

func GetExecutablePath() string {
	exePath, _:= os.Executable()
	exePath = filepath.Dir(exePath)
	return exePath
}