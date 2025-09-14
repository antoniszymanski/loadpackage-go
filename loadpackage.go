// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package loadpackage

import (
	"errors"
	"strings"

	"golang.org/x/tools/go/packages"
)

func Load(pattern string, cfg *packages.Config) (*packages.Package, error) {
	// https://pkg.go.dev/cmd/go#hdr-Package_lists_and_patterns
	// https://pkg.go.dev/golang.org/x/tools/go/packages#pkg-overview
	switch pattern {
	case "main", "pattern=main",
		"all", "pattern=all",
		"std", "pattern=std",
		"cmd", "pattern=cmd",
		"tool", "pattern=tool":
		return nil, errors.New("pattern cannot be a reserved name")
	}
	if strings.Contains(pattern, "...") {
		return nil, errors.New("pattern cannot contain wildcards")
	}
	pkgs, err := packages.Load(cfg, pattern)
	if err != nil {
		return nil, err
	}
	if len(pkgs) != 1 {
		return nil, errors.New("expected exactly one package")
	}
	if err = Validate(pkgs[0]); err != nil {
		return nil, err
	}
	return pkgs[0], nil
}

func Validate(pkg *packages.Package) error {
	if pkg == nil {
		return errors.New("package is nil")
	}
	if len(pkg.Errors) == 0 && (pkg.Module == nil || pkg.Module.Error == nil) {
		return nil
	}
	var e Error
	e.Errors = pkg.Errors
	if pkg.Module != nil {
		e.ModuleError = pkg.Module.Error
	}
	return &e
}

type Error struct {
	Errors      []packages.Error
	ModuleError *packages.ModuleError
}

func (e *Error) Error() string {
	var sb strings.Builder
	for i, err := range e.Errors {
		sb.WriteString(err.Error())
		if e.ModuleError != nil || i < len(e.Errors)-1 {
			sb.WriteByte('\n')
		}
	}
	if e.ModuleError != nil {
		sb.WriteString(e.ModuleError.Err)
	}
	return sb.String()
}

func (e *Error) Unwrap() []error {
	errs := make([]error, 0, len(e.Errors))
	for _, err := range e.Errors {
		errs = append(errs, err)
	}
	return errs
}
