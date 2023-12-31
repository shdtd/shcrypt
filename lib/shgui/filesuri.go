package shgui

type FilesURI struct {
	SrcURI, KeyURI, OutURI string
}

func (f *FilesURI) SetSrcURI(uri string) {
	f.SrcURI = uri
}

func (f *FilesURI) SetKeyURI(uri string) {
	f.KeyURI = uri
}

func (f *FilesURI) SetOutURI(uri string) {
	f.OutURI = uri
}
