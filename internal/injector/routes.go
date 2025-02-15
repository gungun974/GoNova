package injector

import (
	"go/parser"
	"go/token"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/gungun974/gonova/internal/logger"
)

func InjectRoutesNewRouter(path string, routerName string, routeURL string) {
	f, err := decorator.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}

	found := false

	var routerVariableName *string

inject_loop:
	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*dst.FuncDecl)
		if !ok {
			continue
		}

		if funcDecl.Name.Name != "MainRouter" {
			continue
		}

		containerVariable := "c"

		for _, field := range funcDecl.Type.Params.List {
			selectorExpr, ok := field.Type.(*dst.SelectorExpr)
			if !ok {
				continue
			}

			mod, ok := selectorExpr.X.(*dst.Ident)
			if !ok {
				continue
			}

			if mod.Name != "internal" || selectorExpr.Sel.Name != "Container" {
				continue
			}

			if len(field.Names) < 1 {
				continue
			}

			containerVariable = field.Names[0].Name

			break
		}

		for j, stmt := range funcDecl.Body.List {
			switch x := stmt.(type) {
			case *dst.AssignStmt:
				if routerVariableName != nil {
					continue
				}
				var newRouterIndex *int

				for i, rightStmt := range x.Rhs {
					switch x := rightStmt.(type) {
					case *dst.CallExpr:
						selectorExpr, ok := x.Fun.(*dst.SelectorExpr)
						if !ok {
							continue
						}

						mod, ok := selectorExpr.X.(*dst.Ident)
						if !ok {
							continue
						}

						if mod.Name == "chi" && selectorExpr.Sel.Name == "NewRouter" {
							newRouterIndex = &i
							break
						}
					}
				}

				if newRouterIndex == nil {
					continue
				}

				if *newRouterIndex >= len(x.Lhs) {
					continue
				}

				expr := x.Lhs[*newRouterIndex]
				ident, ok := expr.(*dst.Ident)
				if !ok {
					continue
				}
				routerVariableName = &ident.Name
			case *dst.ExprStmt:
				if routerVariableName == nil {
					continue
				}
				callExpr, ok := x.X.(*dst.CallExpr)
				if !ok {
					continue
				}

				selectorExpr, ok := callExpr.Fun.(*dst.SelectorExpr)
				if !ok {
					continue
				}

				variable, ok := selectorExpr.X.(*dst.Ident)
				if !ok {
					continue
				}

				if variable.Name != *routerVariableName || selectorExpr.Sel.Name != "Mount" {
					continue
				}

				if len(callExpr.Args) < 1 {
					continue
				}

				argBasicLit, ok := callExpr.Args[0].(*dst.BasicLit)
				if !ok {
					continue
				}

				if argBasicLit.Kind != token.STRING || argBasicLit.Value != "\"/\"" {
					continue
				}

				assignStmt := &dst.ExprStmt{
					X: &dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   dst.NewIdent(*routerVariableName),
							Sel: dst.NewIdent("Mount"),
						},
						Args: []dst.Expr{
							&dst.BasicLit{
								Kind:  token.STRING,
								Value: "\"" + routeURL + "\"",
							},
							&dst.CallExpr{
								Fun: dst.NewIdent(routerName),
								Args: []dst.Expr{
									dst.NewIdent(containerVariable),
								},
							},
						},
					},
				}

				funcDecl.Body.List = append(funcDecl.Body.List[:j], append([]dst.Stmt{assignStmt}, funcDecl.Body.List[j:]...)...)

				found = true
				break inject_loop
			case *dst.ReturnStmt:
				if routerVariableName == nil {
					continue
				}

				for _, result := range x.Results {
					ident, ok := result.(*dst.Ident)
					if !ok {
						continue
					}

					if ident.Name != *routerVariableName {
						continue
					}

					assignStmt := &dst.ExprStmt{
						X: &dst.CallExpr{
							Fun: &dst.SelectorExpr{
								X:   dst.NewIdent(*routerVariableName),
								Sel: dst.NewIdent("Mount"),
							},
							Args: []dst.Expr{
								&dst.BasicLit{
									Kind:  token.STRING,
									Value: "\"" + routeURL + "\"",
								},
								&dst.CallExpr{
									Fun: dst.NewIdent(routerName),
									Args: []dst.Expr{
										dst.NewIdent(containerVariable),
									},
								},
							},
						},
					}

					funcDecl.Body.List = append(funcDecl.Body.List[:j], append([]dst.Stmt{assignStmt}, funcDecl.Body.List[j:]...)...)

					found = true
					break inject_loop
				}
			}
		}
	}

	if !found {
		logger.InjectorLogger.Fatal("Failed to find `chi.Mux` in /internal/routes/routes.go")
	}

	err = writeDSTFileToPath(f, path)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}
}
