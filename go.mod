module hidevops.io/hiboot

require (
	cloud.google.com/go v0.36.0 // indirect
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/deckarep/golang-set v1.7.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatih/camelcase v1.0.0
	github.com/fatih/structtag v1.0.0
	github.com/gojektech/valkyrie v0.0.0-20190210220504-8f62c1e7ba45
	github.com/golang/mock v1.2.0
	github.com/golang/protobuf v1.2.0
	github.com/gorilla/mux v1.7.0
	github.com/gorilla/websocket v1.4.0
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/iris-contrib/formBinder v0.0.0-20190104093907-fbd5963f41e1 // indirect
	github.com/iris-contrib/go.uuid v2.0.0+incompatible
	github.com/iris-contrib/httpexpect v0.0.0-20180314041918-ebe99fcebbce
	github.com/iris-contrib/middleware v0.0.0-20171114084220-1060fbb0ce08
	github.com/kataras/golog v0.0.0-20180321173939-03be10146386
	github.com/kataras/iris v11.0.3+incompatible
	github.com/kataras/pio v0.0.0-20180511174041-a9733b5b6b83
	github.com/mitchellh/mapstructure v1.1.2
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/opentracing/opentracing-go v1.0.2
	github.com/pkg/errors v0.8.0
	github.com/sony/sonyflake v0.0.0-20160530021500-fa881fb1052b
	github.com/spf13/afero v1.1.2
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/stretchr/testify v1.2.2
	github.com/uber-go/atomic v1.3.2 // indirect
	github.com/uber/jaeger-client-go v2.15.0+incompatible
	github.com/uber/jaeger-lib v1.5.0+incompatible // indirect
	github.com/valyala/bytebufferpool v1.0.0
	go.uber.org/atomic v1.3.2 // indirect
	golang.org/x/net v0.0.0-20190213061140-3a22650c66bd
	google.golang.org/grpc v1.17.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2
	gopkg.in/yaml.v2 v2.2.1
	hidevops.io/viper v1.3.2
)

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.36.0
	golang.org/x/build => github.com/golang/build v0.0.0-20190215225244-0261b66eb045
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20181030022821-bc7917b19d8f
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190212162250-21964bba6549
	golang.org/x/lint => github.com/golang/lint v0.0.0-20181217174547-8f45f776aaf1
	golang.org/x/net => github.com/golang/net v0.0.0-20181029044818-c44066c5c816
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20181017192945-9dcd33a902f4
	golang.org/x/perf => github.com/golang/perf v0.0.0-20190124201629-844a5f5b46f4
	golang.org/x/sync => github.com/golang/sync v0.0.0-20181221193216-37e7f081c4d4
	golang.org/x/sys => github.com/golang/sys v0.0.0-20181029174526-d69651ed3497
	golang.org/x/text => github.com/golang/text v0.3.0
	golang.org/x/time => github.com/golang/time v0.0.0-20180412165947-fbb02b2291d2
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190214204934-8dcb7bc8c7fe
	golang.org/x/vgo => github.com/golang/vgo v0.0.0-20180912184537-9d567625acf4
	google.golang.org/api => github.com/googleapis/googleapis v0.0.0-20190215163516-1a4f0f12777d
	google.golang.org/appengine => github.com/golang/appengine v1.4.0
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20190215211957-bd968387e4aa
	google.golang.org/grpc => github.com/grpc/grpc-go v1.14.0
)
