package jsonstream

type (
	Small struct {
		Name string `json:"name"`
	}
	Large struct {
		Name        string   `json:"name,omitempty"`
		FullName    string   `json:"full_name,omitempty"`
		Occupations []string `json:"occupations,omitempty"`
	}
)

var configs = []struct {
	name   string
	config *Config
}{
	{"default", DefaultConfig},
	{"std", StdConfig},
	{"go-json", GoJSONConfig},
}

var smallSamples = []Small{
	{
		Name: "Barbie",
	},
	{
		Name: "Barbie",
	},
	{
		Name: "Ken",
	},
}

var benchSample = Small{
	Name: "Barbie",
}
