package config

var BaseInfo = &BaseConfig{}

type BaseConfig struct {
    IDLInfo      IDLConfig    `yaml:"idl_info"`
    PlatformInfo SourceConfig `yaml:"platform_info"` // default: tcc
    ServiceInfo  SourceConfig `yaml:"service_info"`  // default: tcc
}

type IDLConfig struct {
    // Source enum(bam、file)
    Source string `yaml:"source"`

    // RootDir is the root directory of idl files only for file source.
    RootDir string `yaml:"root_dir"`

    // IndexFile is the index of file only for file source
    // which key is the psm+version, and value is the index of idl.
    IndexFile string `yaml:"index_file"`
}

type SourceConfig struct {
    Type           string `yaml:"type"`            // enum(tcc、file)
    Root           string `yaml:"root"`            // file: namespace(psm); tcc: root directory
    UpdateInterval int    `yaml:"update_interval"` // unit: second, default: 10
}
