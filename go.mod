module github.com/sonatype-nexus-community/ahab

go 1.14

require (
	github.com/common-nighthawk/go-figure v0.0.0-20200604155835-c37800f1341b
	github.com/jedib0t/go-pretty/v6 v6.0.3
	github.com/logrusorgru/aurora v0.0.0-20200102142835-e9ef32dff381
	github.com/shopspring/decimal v1.2.0
	github.com/sirupsen/logrus v1.6.0
	github.com/sonatype-nexus-community/go-sona-types v0.0.2
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	golang.org/x/sys v0.0.0-20200602225109-6fdc65e7d980 // indirect
)

replace github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2

replace golang.org/x/crypto => golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9

replace golang.org/x/text => golang.org/x/text v0.3.3
