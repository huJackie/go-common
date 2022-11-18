package log

type Config struct {
	Level      string `json:"level" yaml:"level"`           // 日志级别
	Host       string `json:"host" yaml:"host"`             // os.Hostname()
	AppId      string `json:"appId" yaml:"appId"`           // 应用的ID号
	Path       string `json:"path" yaml:"path"`             // 输出路径
	Stdout     bool   `json:"stdout" yaml:"stdout"`         // 标准输出
	MaxSize    int    `json:"maxsize" yaml:"maxsize"`       // 单个文件日志大小
	MaxBackups int    `json:"maxBackups" yaml:"maxBackups"` // 最大文件数量
	MaxAge     int    `json:"maxage" yaml:"maxage"`         // 最大有效期
	Compress   bool   `json:"compress" yaml:"compress"`     // 是否压缩
}
