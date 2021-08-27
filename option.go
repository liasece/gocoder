package gocoder

// SetOption type
type SetOption struct {
	notCast *bool
}

// NewSetOpt func
func NewSetOpt() *SetOption {
	return &SetOption{}
}

// Cast func
func (o *SetOption) Cast(v bool) *SetOption {
	noCast := !v
	o.notCast = &noCast
	return o
}

// MergeSetOpt func
func MergeSetOpt(opts ...*SetOption) *SetOption {
	res := &SetOption{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.notCast != nil {
			res.notCast = opt.notCast
		}
	}
	return res
}

// ToCodeOption type
type ToCodeOption struct {
	pkgTool  PkgTool
	pkgName  *string
	noPretty *bool
}

// NewToCodeOpt func
func NewToCodeOpt() *ToCodeOption {
	return &ToCodeOption{}
}

// PkgTool func
func (o *ToCodeOption) PkgTool(v PkgTool) *ToCodeOption {
	o.pkgTool = v
	return o
}

// PkgTool func
func (o *ToCodeOption) PkgName(v string) *ToCodeOption {
	o.pkgName = &v
	return o
}

// PkgTool func
func (o *ToCodeOption) NoPretty(v bool) *ToCodeOption {
	o.noPretty = &v
	return o
}

// MergeToCodeOpt func
func MergeToCodeOpt(opts ...*ToCodeOption) *ToCodeOption {
	res := &ToCodeOption{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.pkgTool != nil {
			res.pkgTool = opt.pkgTool
		}
		if opt.pkgName != nil {
			res.pkgName = opt.pkgName
		}
		if opt.noPretty != nil {
			res.noPretty = opt.noPretty
		}
	}
	return res
}
