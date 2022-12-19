package ast

import "go/ast"

type Packages struct {
	List []*Package
}

type Package struct {
	*ast.Package
	Name  string
	Alias string
}

func (p *Packages) Add(pkgs ...*Package) {
	p.List = append(p.List, pkgs...)
}

func (p *Packages) MergeFrom(pkgs ...*Packages) {
	for _, pkg := range pkgs {
		for _, pkg := range pkg.List {
			pk := p.Get(pkg.Name)
			if pk != nil {
				pk.Name = pkg.Name
				pk.Alias = pkg.Alias
				{
					// merge ast.Package
					pk.Package.Name = pkg.Package.Name
					pk.Package.Scope = pkg.Package.Scope
					for k, v := range pkg.Package.Imports {
						pk.Package.Imports[k] = v
					}
					for k, v := range pkg.Package.Files {
						pk.Package.Files[k] = v
					}
				}
			} else {
				p.Add(pkg)
			}
		}
	}
}

func (p *Packages) Get(name string) *Package {
	for _, pkg := range p.List {
		if pkg.Name == name {
			return pkg
		}
	}
	return nil
}

// find package by alias
func (p *Packages) FindAlias(alias string) []*Package {
	var list []*Package
	for _, pkg := range p.List {
		if pkg.Alias == alias {
			list = append(list, pkg)
		}
	}
	return list
}
