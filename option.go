package gocoder

// SetOption type
type SetOption struct {
	notCast *bool
}

// NewSetOpt func
func NewSetOpt() *SetOption {
	return &SetOption{
		notCast: nil,
	}
}

// Cast func
func (o *SetOption) Cast(v bool) *SetOption {
	noCast := !v
	o.notCast = &noCast
	return o
}

// MergeSetOpt func
func MergeSetOpt(opts ...*SetOption) *SetOption {
	res := &SetOption{
		notCast: nil,
	}
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
	pkgPath  *string
	noPretty *bool
}

// NewToCodeOpt func
func NewToCodeOpt() *ToCodeOption {
	var res ToCodeOption
	return &res
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

func (o *ToCodeOption) PkgPath(v string) *ToCodeOption {
	o.pkgPath = &v
	return o
}

func (o *ToCodeOption) GetPkgPath() *string {
	return o.pkgPath
}

// PkgTool func
func (o *ToCodeOption) NoPretty(v bool) *ToCodeOption {
	o.noPretty = &v
	return o
}

// MergeToCodeOpt func
func MergeToCodeOpt(opts ...*ToCodeOption) *ToCodeOption {
	var res ToCodeOption
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
		if opt.pkgPath != nil {
			res.pkgPath = opt.pkgPath
		}
		if opt.noPretty != nil {
			res.noPretty = opt.noPretty
		}
	}
	return &res
}
