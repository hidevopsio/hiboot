package factory

type PostProcessor interface {
	BeforeInitialization()
	AfterInitialization()
}

var (
	postProcessors []PostProcessor
)

func init() {

}

func AddPostProcessor(p ...PostProcessor)  {
	postProcessors = append(postProcessors, p...)
}

