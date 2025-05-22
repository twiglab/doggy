package doggy

import "runtime"

const version = "0.1.10"

var (
	GitCommit string //Git提交号
	BuildTime string //编译时间
)

type Ver struct {
	Version   string
	GitCommit string
	BuildTime string
	GoVersion string
	OsArch    string
}

func Version() *Ver {
	return &Ver{
		Version:   version,
		GitCommit: GitCommit,
		BuildTime: BuildTime,
		GoVersion: runtime.Version(),
		OsArch:    runtime.GOOS + "/" + runtime.GOARCH,
	}
}
