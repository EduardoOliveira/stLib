package state

type Project struct {
	UUID             string                   `json:"uuid" toml:"uuid" form:"uuid" query:"uuid"`
	Name             string                   `json:"name" toml:"name" form:"name" query:"name"`
	Path             string                   `json:"path" toml:"path" form:"path" query:"path"`
	ExternalLink     string                   `json:"external_link" toml:"external_link" form:"external_link" query:"external_link"`
	Models           map[string]*Model        `json:"-" toml:"models" form:"models" query:"models"`
	Images           map[string]*ProjectImage `json:"-" toml:"images" form:"images" query:"images"`
	Slices           map[string]*Slice        `json:"-" toml:"slices" form:"slices" query:"slices"`
	Tags             []string                 `json:"tags" toml:"tags" form:"tags" query:"tags"`
	DefaultModel     *Model                   `json:"default_model" toml:"default_model" form:"default_model" query:"default_model"`
	DefaultImagePath string                   `json:"default_image_path" toml:"default_image_path" form:"default_image_path" query:"default_image_path"`
	Initialized      bool                     `json:"initialized" toml:"initialized" form:"initialized" query:"initialized"`
}

type Model struct {
	SHA1     string `json:"sha1" toml:"sha1" form:"sha1" query:"sha1"`
	Name     string `json:"name" toml:"name" form:"name" query:"name"`
	Path     string `json:"path" toml:"path" form:"path" query:"path"`
	FileName string `json:"file_name" toml:"file_name" form:"file_name" query:"file_name"`
	Tags     []*Tag `json:"tags" toml:"tags" form:"tags" query:"tags"`
}

type Tag struct {
	Name  string `json:"name" toml:"name" form:"name" query:"name"`
	Value string `json:"value" toml:"value" form:"value" query:"value"`
}

type ProjectImage struct {
	SHA1      string `json:"sha1" toml:"sha1" form:"sha1" query:"sha1"`
	Name      string `json:"name" toml:"name" form:"name" query:"name"`
	Path      string `json:"path" toml:"path" form:"path" query:"path"`
	Extension string `json:"extension" toml:"extension" form:"extension" query:"extension"`
	ModelSHA1 string `json:"model_sha1" toml:"model_sha1" form:"model_sha1" query:"model_sha1"`
}

type Slice struct {
	SHA1       string   `json:"sha1" toml:"sha1" form:"sha1" query:"sha1"`
	Name       string   `json:"name" toml:"name" form:"name" query:"name"`
	Path       string   `json:"path" toml:"path" form:"path" query:"path"`
	Extension  string   `json:"extension" toml:"extension" form:"extension" query:"extension"`
	ModelsSHA1 []string `json:"models_sha1" toml:"models_sha1" form:"models_sha1" query:"models_sha1"`
}
