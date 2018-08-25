package factory

type PostProcessor interface {
	BeforeInitialization()
	AfterInitialization()
}

type postProcessor struct{
}

var (
	postProcessors []PostProcessor
)

func init() {

}

func AddPostProcessor(p ...PostProcessor)  {
	postProcessors = append(postProcessors, p...)
}

func (p *postProcessor) BeforeInitialization()  {
	for _, processor := range postProcessors {
		processor.BeforeInitialization()
	}
}

func (p *postProcessor) AfterInitialization()  {
	for _, processor := range postProcessors {
		processor.AfterInitialization()
	}
}