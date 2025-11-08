package analyzer

type AnalyzedDependency interface {
	GetDependencies() []AnalyzedDependency

	GetName() string

	GetNewFunction() string

	GetPkgPath() string

	GetImportName() string

	GetType() string
}

type AnalyzedDatabaseDependency struct{}

func (a *AnalyzedDatabaseDependency) GetDependencies() []AnalyzedDependency {
	return []AnalyzedDependency{}
}

func (a *AnalyzedDatabaseDependency) GetName() string {
	return ""
}

func (a *AnalyzedDatabaseDependency) GetNewFunction() string {
	return ""
}

func (a *AnalyzedDatabaseDependency) GetPkgPath() string {
	return ""
}

func (a *AnalyzedDatabaseDependency) GetImportName() string {
	return ""
}

func (a *AnalyzedDatabaseDependency) GetType() string {
	return ""
}
