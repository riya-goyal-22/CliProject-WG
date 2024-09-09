package interfaces

type PromptInterface interface {
	Run() (string, error)
}

//type PromptWrapper struct {
//	*promptui.Prompt
//}
//
//func (pw *PromptWrapper) Run() (string, error) {
//	return pw.Prompt.Run()
//}
