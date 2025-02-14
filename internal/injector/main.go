package injector

import (
	"go/parser"
	"go/token"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
	"github.com/gungun974/gonova/internal/logger"
)

func InjectMainDatabase(path string, projectName string) {
	f, err := decorator.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}

	addImport(f, projectName+"/internal/database", "")

	hasInject := false

	dstutil.Apply(f, nil, func(c *dstutil.Cursor) bool {
		n := c.Node()
		switch x := n.(type) {
		case *dst.BlockStmt:
			for i, stmt := range x.List {
				found := false

				dst.Inspect(stmt, func(n dst.Node) bool {
					callExpr, ok := n.(*dst.CallExpr)
					if !ok {
						return true
					}

					selectorExpr, ok := callExpr.Fun.(*dst.SelectorExpr)
					if !ok {
						return true
					}
					mod, ok := selectorExpr.X.(*dst.Ident)
					if !ok {
						return true
					}
					if mod.Name != "internal" {
						return true
					}
					if selectorExpr.Sel.Name != "NewContainer" {
						return true
					}

					found = true

					callExpr.Args = append(
						[]dst.Expr{dst.NewIdent("db")},
						callExpr.Args...,
					)

					return false
				})

				if !found {
					continue
				}

				assignStmt := &dst.AssignStmt{
					Lhs: []dst.Expr{&dst.Ident{Name: "db"}},
					Tok: token.DEFINE,
					Rhs: []dst.Expr{&dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   dst.NewIdent("database"),
							Sel: dst.NewIdent("Connect"),
						},
					}},
				}

				x.List = append(x.List[:i], append([]dst.Stmt{assignStmt}, x.List[i:]...)...)

				hasInject = true

				return true
			}
		}

		return true
	})

	if !hasInject {
		logger.InjectorLogger.Fatal("Failed to find `internal.NewContainer` in main.go")
	}

	err = writeDSTFileToPath(f, path)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}
}
