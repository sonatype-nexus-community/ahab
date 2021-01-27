module github.com/sonatype-nexus-community/ahab

go 1.14

require (
	github.com/common-nighthawk/go-figure v0.0.0-20200609044655-c4b36f998cf2
	github.com/jedib0t/go-pretty/v6 v6.0.5
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/mitchellh/go-homedir v1.1.0
	github.com/shopspring/decimal v1.2.0
	github.com/sirupsen/logrus v1.7.0
	github.com/sonatype-nexus-community/go-sona-types v0.0.10
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f // indirect
)

replace github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2

replace golang.org/x/crypto => golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9

replace golang.org/x/text => golang.org/x/text v0.3.3

// fix vulnerability: CVE-2020-15114 in etcd v3.3.13+incompatible
replace github.com/coreos/etcd => github.com/coreos/etcd v3.3.24+incompatible
