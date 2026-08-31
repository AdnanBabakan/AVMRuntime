package sections

type Section struct {
	Name    string
	Content string
}

type Tree map[string]string

type SectionPolicy struct {
	Required  bool
	Singleton bool
}

var Policy = map[string]SectionPolicy{
	"target": {
		Required:  true,
		Singleton: true,
	},
	"namespace": {
		Required:  true,
		Singleton: true,
	},
	"entry": {
		Required:  false,
		Singleton: true,
	},
	"capabilities": {
		Required:  false,
		Singleton: false,
	},
	"batteries": {
		Required:  false,
		Singleton: false,
	},
	"deps": {
		Required:  false,
		Singleton: false,
	},
	"include": {
		Required:  false,
		Singleton: false,
	},
	"exports": {
		Required:  false,
		Singleton: false,
	},
	"constants": {
		Required:  false,
		Singleton: false,
	},
	"reg_aliases": {
		Required:  false,
		Singleton: false,
	},
	"type_aliases": {
		Required:  false,
		Singleton: false,
	},
	"macros": {
		Required:  false,
		Singleton: false,
	},
	"code": {
		Required:  true,
		Singleton: false,
	},
}

type FullTree struct {
	Target *TargetSection
}
