package options

import "github.com/spf13/pflag"

// AliyunOSSOptions holds Alibaba Cloud OSS configuration.
type AliyunOSSOptions struct {
	Endpoint     string `yaml:"endpoint"`
	AccessKey    string `yaml:"access_key"`
	AccessSecret string `yaml:"access_secret"`
	Bucket       string `yaml:"bucket"`
}

func NewAliyunOSSOptions() *AliyunOSSOptions {
	return &AliyunOSSOptions{}
}

func (o *AliyunOSSOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.Endpoint, FlagNameFunc(fs, "endpoint"), o.Endpoint, "Aliyun OSS endpoint.")
	fs.StringVar(&o.AccessKey, FlagNameFunc(fs, "access-key"), o.AccessKey, "Aliyun OSS access key.")
	fs.StringVar(&o.AccessSecret, FlagNameFunc(fs, "access-secret"), o.AccessSecret, "Aliyun OSS access secret.")
	fs.StringVar(&o.Bucket, FlagNameFunc(fs, "bucket"), o.Bucket, "Aliyun OSS bucket name.")
}

func (o *AliyunOSSOptions) Kind() string {
	return "aliyunoss"
}
