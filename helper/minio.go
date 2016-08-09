package helper

import "github.com/scompo/data-management/domain"

func NewProject(*domain.Project) (error, domain.Project) {
	return nil, domain.Project{}
}

func ProjectExists(name string) (error, bool) {
	return nil, false
}

func All() []string {
	return make([]string, 0)
}
