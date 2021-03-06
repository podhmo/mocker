package mocker

import (
	"bytes"
	"fmt"
	"go/build"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/golang/mock/mockgen/model"
	format "golang.org/x/tools/imports"
)

type mocker struct {
	src, dst, prefix, suffix, pkg, selfpkg string
	intfs                                  []string
}

func New(src, dst, pkg, prefix, suffix, selfpkg string, intfs []string) (*mocker, error) {
	return &mocker{
		src:     src,
		dst:     dst,
		pkg:     pkg,
		prefix:  prefix,
		suffix:  suffix,
		selfpkg: selfpkg,
		intfs:   intfs,
	}, nil
}

func (m *mocker) Mock() error {
	pkg, err := ParseFile(m.src)
	if err != nil {
		log.Fatalf("loading src failed: %v", err)
	}

	outputPackagePath := m.selfpkg
	dstDir, _ := filepath.Abs(filepath.Dir(m.dst))
	for _, prefix := range build.Default.SrcDirs() {
		if strings.HasPrefix(dstDir, prefix) {
			if rel, err := filepath.Rel(prefix, dstDir); err == nil {
				outputPackagePath = rel
				break
			}
		}
	}

	dst := os.Stdout
	if len(m.dst) > 0 {
		if err := os.MkdirAll(filepath.Dir(m.dst), os.ModePerm); err != nil {
			log.Fatalf("Unable to create directory: %v", err)
		}
		f, err := os.Create(m.dst)
		if err != nil {
			log.Fatalf("Failed opening destination file: %v", err)
		}
		defer f.Close()
		dst = f
	}

	if m.pkg == "" {
		m.pkg = pkg.Name
	}

	g := new(generator)
	g.filename = m.src
	g.prefix = m.prefix
	g.suffix = m.suffix
	if err := g.Generate(pkg, m.pkg, m.intfs, outputPackagePath); err != nil {
		return err
	}

	if _, err := dst.Write(g.Output()); err != nil {
		log.Fatalf("Failed writing to destination: %v", err)
	}

	return nil
}

type generator struct {
	buf            bytes.Buffer
	pkgMap         map[string]string // import path to pkg name
	mockNames      map[string]string
	filename       string
	indent         string
	prefix, suffix string
}

func (g *generator) p(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, g.indent+format+"\n", args...)
}

func (g *generator) in() {
	g.indent += "\t"
}

func (g *generator) out() {
	if len(g.indent) > 0 {
		g.indent = g.indent[0 : len(g.indent)-1]
	}
}

// Output returns the generator's output, formatted in the standard Go style.
func (g *generator) Output() []byte {
	options := &format.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
	}
	src, err := format.Process(g.filename, g.buf.Bytes(), options)
	if err != nil {
		log.Fatalf("Failed to format generated source code: %s\n%s", err, g.buf.String())
	}
	return src
}

func (g *generator) Generate(pkg *model.Package, pkgName string, intfs []string, outputPkgPath string) error {
	g.p("// Code generated by mocker. DO NOT EDIT.")
	g.p("// github.com/travisjeffery/mocker")
	if g.filename != "" {
		g.p("// Source: %v", g.filename)
	}
	g.p("")

	im := pkg.Imports()

	sortedPaths := make([]string, len(im), len(im))
	sortedPaths = append(sortedPaths, "sync")
	i := 0
	for path := range im {
		sortedPaths[i] = path
		i++
	}
	sort.Strings(sortedPaths)

	g.pkgMap = make(map[string]string, len(im))
	names := make(map[string]bool, len(im))
	for _, path := range sortedPaths {
		base := sanitize(path)

		pkgName := base
		i := 0
		for names[pkgName] || token.Lookup(pkgName).IsKeyword() {
			pkgName = base + strconv.Itoa(i)
			i++
		}

		g.pkgMap[path] = pkgName
		names[pkgName] = true
	}

	g.p("package %v", pkgName)
	g.p("")
	g.p("import (")
	g.in()

	for path, pkg := range g.pkgMap {
		if path == outputPkgPath {
			continue
		}
		g.p("%v %q", pkg, path)
	}
	g.out()
	g.p(")")

	for _, intf := range pkg.Interfaces {
		if !contains(intfs, intf.Name) {
			continue
		}
		if err := g.GenerateInterface(intf, outputPkgPath); err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) GenerateInterface(intf *model.Interface, outputPkgPath string) error {
	mockType := g.mockName(intf.Name)

	g.p("")
	g.p("// %v is a mock of %v interface", mockType, intf.Name)
	g.p("type %v struct {", mockType)
	g.in()

	for _, m := range intf.Methods {
		g.p("lock%v sync.Mutex", m.Name)

		argNames := g.getArgNames(m)
		argTypes := g.getArgTypes(m, outputPkgPath)
		argString := makeArgString(argNames, argTypes)

		rets := make([]string, len(m.Out))
		for i, p := range m.Out {
			rets[i] = p.Type.String(g.pkgMap, outputPkgPath)
		}
		retString := strings.Join(rets, ", ")
		if len(rets) > 1 {
			retString = "(" + retString + ")"
		}
		if retString != "" {
			retString = " " + retString
		}

		g.p("%vFunc func(%v) %v", m.Name, argString, retString)
		g.p("")
	}

	g.p("calls struct {")
	g.in()
	for _, m := range intf.Methods {
		g.p("%v []struct {", m.Name)
		g.in()

		argNames := g.getArgNames(m)
		argTypes := g.getArgTypes(m, outputPkgPath)

		for i, name := range argNames {
			s := fmt.Sprintf("%v %v", strings.Title(name), argTypes[i])
			s = strings.Replace(s, "...", "[]", -1)
			g.p(s)
		}

		g.out()
		g.p("}")
	}
	g.out()
	g.p("}")

	g.out()
	g.p("}")
	g.p("")

	return g.GenerateMethods(mockType, intf, outputPkgPath)
}

func (g *generator) GenerateMethods(mockType string, intf *model.Interface, outputPkgPath string) error {
	for _, m := range intf.Methods {
		g.p("")
		if err := g.GenerateMethod(mockType, m, outputPkgPath); err != nil {
			return err
		}
		g.p("")
	}

	g.p("// Reset resets the calls made to the mocked methods.")
	g.p("func (m *%v) Reset() {", mockType)
	g.in()
	for _, m := range intf.Methods {
		g.p("m.lock%v.Lock()", m.Name)
		g.p("m.calls.%v = nil", m.Name)
		g.p("m.lock%v.Unlock()", m.Name)
	}
	g.out()
	g.p("}")

	return nil
}

func (g *generator) GenerateMethod(mockType string, m *model.Method, outputPkgPath string) error {
	argNames := g.getArgNames(m)
	argTypes := g.getArgTypes(m, outputPkgPath)
	argString := makeArgString(argNames, argTypes)

	rets := make([]string, len(m.Out))
	for i, p := range m.Out {
		rets[i] = p.Type.String(g.pkgMap, outputPkgPath)
	}
	retString := strings.Join(rets, ", ")
	if len(rets) > 1 {
		retString = "(" + retString + ")"
	}
	if retString != "" {
		retString = " " + retString
	}

	ia := newIdentifierAllocator(argNames)
	idRecv := ia.allocateIdentifier("m")

	g.p("// %v mocks base method by wrapping the associated func.", m.Name)
	g.p("func (%v *%v) %v(%v)%v {", idRecv, mockType, m.Name, argString, retString)
	g.in()
	g.p("%s.lock%s.Lock()", idRecv, m.Name)
	g.p("defer %s.lock%s.Unlock()", idRecv, m.Name)

	g.p("")
	g.p("if %v.%vFunc == nil {", idRecv, m.Name)
	g.in()
	g.p("panic(\"mocker: %v.%vFunc is nil but %v.%v was called.\")", mockType, m.Name, mockType, m.Name)
	g.out()
	g.p("}")
	g.p("")

	g.p("call := struct {")
	g.in()
	for i, name := range argNames {
		s := fmt.Sprintf("%v %v", strings.Title(name), argTypes[i])
		s = strings.Replace(s, "...", "[]", -1)
		g.p(s)
	}
	g.out()
	g.p("}{")
	g.in()
	for _, name := range argNames {
		g.p("%v: %v,", strings.Title(name), name)
	}
	g.out()
	g.p("}")
	g.p("")

	g.p("%v.calls.%v = append(%v.calls.%v, call)", idRecv, m.Name, idRecv, m.Name)
	g.p("")

	var callArgs string
	if len(argNames) > 0 {
		callArgs = strings.Join(argNames, ", ")
	}

	if m.Variadic != nil {
		callArgs += "..."
	}

	if len(m.Out) == 0 {
		g.p(`%v.%vFunc(%v)`, idRecv, m.Name, callArgs)
	} else {
		g.p(`return %v.%vFunc(%v)`, idRecv, m.Name, callArgs)
	}

	g.out()
	g.p("}")
	g.p("")

	g.p("// %vCalled returns true if %v was called at least once.", m.Name, m.Name)
	g.p("func (%v *%v) %vCalled() bool {", idRecv, mockType, m.Name)
	g.in()
	g.p("%s.lock%s.Lock()", idRecv, m.Name)
	g.p("defer %s.lock%s.Unlock()", idRecv, m.Name)
	g.p("")
	g.p("return len(%v.calls.%v) > 0", idRecv, m.Name)
	g.out()
	g.p("}")

	g.p("// %vCalls returns the calls made to %v.", m.Name, m.Name)
	g.p("func (%v *%v) %vCalls() []struct {", idRecv, mockType, m.Name)

	g.in()
	for i, name := range argNames {
		s := fmt.Sprintf("%v %v", strings.Title(name), argTypes[i])
		s = strings.Replace(s, "...", "[]", -1)
		g.p(s)
	}
	g.out()
	g.p("} {")

	g.in()
	g.p("%s.lock%s.Lock()", idRecv, m.Name)
	g.p("defer %s.lock%s.Unlock()", idRecv, m.Name)
	g.p("")
	g.p("return %v.calls.%v", idRecv, m.Name)
	g.out()
	g.p("}")

	return nil
}

func makeArgString(argNames, argTypes []string) string {
	args := make([]string, len(argNames))
	for i, name := range argNames {
		// specify the type only once for consecutive args of the same type
		if i+1 < len(argTypes) && argTypes[i] == argTypes[i+1] {
			args[i] = name
		} else {
			args[i] = name + " " + argTypes[i]
		}
	}
	return strings.Join(args, ", ")
}

func (g *generator) getArgNames(m *model.Method) []string {
	argNames := make([]string, len(m.In))
	for i, p := range m.In {
		name := p.Name
		if name == "" {
			name = fmt.Sprintf("arg%d", i)
		}
		argNames[i] = name
	}
	if m.Variadic != nil {
		name := m.Variadic.Name
		if name == "" {
			name = fmt.Sprintf("arg%d", len(m.In))
		}
		argNames = append(argNames, name)
	}
	return argNames
}

func (g *generator) getArgTypes(m *model.Method, pkgOverride string) []string {
	argTypes := make([]string, len(m.In))
	for i, p := range m.In {
		argTypes[i] = p.Type.String(g.pkgMap, pkgOverride)
	}
	if m.Variadic != nil {
		argTypes = append(argTypes, "..."+m.Variadic.Type.String(g.pkgMap, pkgOverride))
	}
	return argTypes
}

// The name of the mock type to use for the given interface identifier.
func (g *generator) mockName(typeName string) string {
	if mockName, ok := g.mockNames[typeName]; ok {
		return mockName
	}
	return g.prefix + typeName + g.suffix
}

type identifierAllocator map[string]struct{}

func newIdentifierAllocator(taken []string) identifierAllocator {
	a := make(identifierAllocator, len(taken))
	for _, s := range taken {
		a[s] = struct{}{}
	}
	return a
}

func (o identifierAllocator) allocateIdentifier(want string) string {
	id := want
	for i := 2; ; i++ {
		if _, ok := o[id]; !ok {
			o[id] = struct{}{}
			return id
		}
		id = want + "_" + strconv.Itoa(i)
	}
}

// sanitize cleans up a string to make a suitable package name.
func sanitize(s string) string {
	t := ""
	for _, r := range s {
		if t == "" {
			if unicode.IsLetter(r) || r == '_' {
				t += string(r)
				continue
			}
		} else {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				t += string(r)
				continue
			}
		}
		t += "_"
	}
	if t == "_" {
		t = "x"
	}
	return t
}

func contains(sl []string, s string) bool {
	for _, e := range sl {
		if e == s {
			return true
		}
	}
	return false
}
