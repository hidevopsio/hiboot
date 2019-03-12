module hidevops.io/hiboot

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Joker/jade v1.0.0
	github.com/Shopify/goreferrer v0.0.0-20180807163728-b9777dc9f9cc
	github.com/ajg/form v0.0.0-20160822230020-523a5da1a92f
	github.com/aymerick/raymond v2.0.2+incompatible
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.7.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/eknkc/amber v0.0.0-20171010120322-cdade1c07385
	github.com/fatih/camelcase v1.0.0
	github.com/fatih/structs v1.0.0
	github.com/flosch/pongo2 v0.0.0-20180809100617-24195e6d38b0
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gavv/monotime v0.0.0-20171021193802-6f8212e8d10d
	github.com/golang/mock v1.2.0
	github.com/golang/protobuf v1.2.0
	github.com/google/go-querystring v1.0.0
	github.com/gorilla/websocket v1.4.0
	github.com/hashicorp/go-version v1.0.0
	github.com/hashicorp/hcl v1.0.0
	github.com/imkira/go-interpol v1.1.0
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/iris-contrib/blackfriday v2.0.0+incompatible
	github.com/iris-contrib/formBinder v0.0.0-20171010160137-ad9fb86c356f
	github.com/iris-contrib/go.uuid v2.0.0+incompatible
	github.com/iris-contrib/httpexpect v0.0.0-20180314041918-ebe99fcebbce
	github.com/iris-contrib/i18n v0.0.0-20171121225848-987a633949d0
	github.com/iris-contrib/middleware v0.0.0-20171114084220-1060fbb0ce08
	github.com/json-iterator/go v1.1.5
	github.com/juju/errors v0.0.0-20180806074554-22422dad46e1
	github.com/kataras/golog v0.0.0-20180321173939-03be10146386
	github.com/kataras/iris v11.0.3+incompatible
	github.com/kataras/pio v0.0.0-20180511174041-a9733b5b6b83
	github.com/kataras/survey v2.0.0+incompatible
	github.com/klauspost/compress v1.4.0
	github.com/klauspost/cpuid v0.0.0-20180405133222-e7e905edc00e
	github.com/klauspost/crc32 v0.0.0-20161016154125-cb6bfca970f6
	github.com/magiconair/properties v1.8.0
	github.com/mattn/go-colorable v0.0.9
	github.com/mattn/go-isatty v0.0.4
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b
	github.com/microcosm-cc/bluemonday v1.0.1
	github.com/mitchellh/mapstructure v1.1.2
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 v1.0.1
	github.com/moul/http2curl v1.0.0
	github.com/pelletier/go-toml v1.2.0
	github.com/pmezard/go-difflib v1.0.0
	github.com/ryanuber/columnize v2.1.0+incompatible
	github.com/sergi/go-diff v1.0.0
	github.com/shurcooL/sanitized_anchor_name v0.0.0-20170918181015-86672fcb3f95
	github.com/sony/sonyflake v0.0.0-20160530021500-fa881fb1052b
	github.com/spf13/afero v1.1.2
	github.com/spf13/cast v1.3.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/jwalterweatherman v1.0.0
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.2.2
	github.com/valyala/bytebufferpool v1.0.0
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415
	github.com/xeipuuv/gojsonschema v0.0.0-20180816142147-da425ebb7609
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0
	github.com/yudai/gojsondiff v1.0.0
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82
	golang.org/x/crypto v0.0.0-20190211182817-74369b46fc67
	golang.org/x/net v0.0.0-20190213061140-3a22650c66bd
	golang.org/x/sys v0.0.0-20181029174526-d69651ed3497
	golang.org/x/text v0.3.1-0.20180807135948-17ff2d5776d2
	golang.org/x/vgo v0.0.0-20180912184537-9d567625acf4 // indirect
	google.golang.org/genproto v0.0.0-20190201180003-4b09977fb922
	google.golang.org/grpc v1.17.0
	gopkg.in/go-playground/validator.v8 v8.18.2
	gopkg.in/ini.v1 v1.38.3
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
