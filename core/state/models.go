package state

type Project struct {
	UUID             string                   `json:"uuid" toml:"uuid" form:"uuid" query:"uuid"`
	Name             string                   `json:"name" toml:"name" form:"name" query:"name"`
	Path             string                   `json:"path" toml:"path" form:"path" query:"path"`
	ExternalLink     string                   `json:"external_link" toml:"external_link" form:"external_link" query:"external_link"`
	Models           map[string]*Model        `json:"-" toml:"models" form:"models" query:"models"`
	Images           map[string]*ProjectImage `json:"-" toml:"images" form:"images" query:"images"`
	Slices           map[string]*Slice        `json:"-" toml:"slices" form:"slices" query:"slices"`
	Files            map[string]*ProjectFile  `json:"-" toml:"files" form:"files" query:"files"`
	Tags             []string                 `json:"tags" toml:"tags" form:"tags" query:"tags"`
	DefaultImagePath string                   `json:"default_image_path" toml:"default_image_path" form:"default_image_path" query:"default_image_path"`
	Initialized      bool                     `json:"initialized" toml:"initialized" form:"initialized" query:"initialized"`
}

type Model struct {
	SHA1      string   `json:"sha1" toml:"sha1" form:"sha1" query:"sha1"`
	Name      string   `json:"name" toml:"name" form:"name" query:"name"`
	Path      string   `json:"path" toml:"path" form:"path" query:"path"`
	FileName  string   `json:"file_name" toml:"file_name" form:"file_name" query:"file_name"`
	Tags      []string `json:"tags" toml:"tags" form:"tags" query:"tags"`
	Extension string   `json:"extension" toml:"extension" form:"extension" query:"extension"`
	MimeType  string   `json:"mime_type" toml:"mime_type" form:"mime_type" query:"mime_type"`
}

type ProjectImage struct {
	SHA1      string `json:"sha1" toml:"sha1" form:"sha1" query:"sha1"`
	Name      string `json:"name" toml:"name" form:"name" query:"name"`
	Path      string `json:"path" toml:"path" form:"path" query:"path"`
	Extension string `json:"extension" toml:"extension" form:"extension" query:"extension"`
	ModelSHA1 string `json:"model_sha1" toml:"model_sha1" form:"model_sha1" query:"model_sha1"`
	MimeType  string `json:"mime_type" toml:"mime_type" form:"mime_type" query:"mime_type"`
}

type ProjectFile struct {
	SHA1      string `json:"sha1" toml:"sha1" form:"sha1" query:"sha1"`
	Name      string `json:"name" toml:"name" form:"name" query:"name"`
	Path      string `json:"path" toml:"path" form:"path" query:"path"`
	FileName  string `json:"file_name" toml:"file_name" form:"file_name" query:"file_name"`
	Extension string `json:"extension" toml:"extension" form:"extension" query:"extension"`
	MimeType  string `json:"mime_type" toml:"mime_type" form:"mime_type" query:"mime_type"`
}

type Slice struct {
	SHA1       string        `json:"sha1" toml:"sha1" form:"sha1" query:"sha1"`
	Name       string        `json:"name" toml:"name" form:"name" query:"name"`
	Path       string        `json:"path" toml:"path" form:"path" query:"path"`
	Extension  string        `json:"extension" toml:"extension" form:"extension" query:"extension"`
	ModelsSHA1 []string      `json:"models_sha1" toml:"models_sha1" form:"models_sha1" query:"models_sha1"`
	MimeType   string        `json:"mime_type" toml:"mime_type" form:"mime_type" query:"mime_type"`
	Image      *ProjectImage `json:"image" toml:"image" form:"image" query:"image"`
	Slicer     string        `json:"slicer" toml:"slicer" form:"slicer" query:"slicer"`
	Filament   *Filament     `json:"filament" toml:"filament" form:"filament" query:"filament"`
	Cost       float64       `json:"cost" toml:"cost" form:"cost" query:"cost"`
	LayerCount int           `json:"layer_count" toml:"layer_count" form:"layer_count" query:"layer_count"`
	Duration   string        `json:"duration" toml:"duration" form:"duration" query:"duration"`
}

type Filament struct {
	Length float64 `json:"length" toml:"length" form:"length" query:"length"`
	Mass   float64 `json:"mass" toml:"mass" form:"mass" query:"mass"`
	Weight float64 `json:"weight" toml:"weight" form:"weight" query:"weight"`
}
