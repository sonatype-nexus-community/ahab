module github.com/sonatype-nexus-community/ahab

go 1.17

require (
	github.com/common-nighthawk/go-figure v0.0.0-20200609044655-c4b36f998cf2
	github.com/jedib0t/go-pretty/v6 v6.0.5
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/mitchellh/go-homedir v1.1.0
	github.com/shopspring/decimal v1.2.0
	github.com/sirupsen/logrus v1.7.0
	github.com/sonatype-nexus-community/go-sona-types v0.1.3
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/package-url/packageurl-go v0.1.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/recoilme/pudge v1.0.3 // indirect
	github.com/spf13/afero v1.1.2 // indirect
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/ini.v1 v1.51.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

// fix vulnerability: CVE-2020-15114 in etcd v3.3.13+incompatible
replace github.com/coreos/etcd => github.com/coreos/etcd v3.3.24+incompatible

// fix vulnerability: CVE-2021-3121 in github.com/gogo/protobuf v1.2.1
replace github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2

// fix vulnerability: CVE-2022-21698 in github.com/prometheus/client_golang v0.9.3
replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v1.11.1

// fix vulnerability: CVE-2021-38561 in golang.org/x/text v0.3.3
// fix vulnerability: CVE-2022-32149 in golang.org/x/text v0.3.7
replace golang.org/x/text => golang.org/x/text v0.3.8
