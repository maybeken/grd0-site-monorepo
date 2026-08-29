package schema

type Equipment struct {
	Camera string `json:"camera,omitempty"`
	Lens   string `json:"lens,omitempty"`
}

type Exif struct {
	Datetime  string    `json:"datetime,omitempty"`
	Shutter   string    `json:"shutter,omitempty"`
	Fstop     string    `json:"fstop,omitempty"`
	Iso       *uint16    `json:"iso,omitempty"`
	Focal     *float32   `json:"focal,omitempty"`
	Equipment Equipment `json:"equipment,omitempty"`
}

type Asset struct {
	Filename   string `json:"filename"`
	Exif       Exif   `json:"exif,omitempty"`
	Collection string `json:"collection"`
}

type AssetFileList map[string][]Asset
