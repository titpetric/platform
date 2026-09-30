# Testing coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 80%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to cover. Tests with cognitive complexity 0 would be covered by invocation.

## Packages

| Status | Package              | Coverage | Cognitive | Lines |
|--------|----------------------|----------|-----------|-------|
| ✅     | .                    | 93.71%   | 142       | 994   |
| ✅     | cmd                  | 46.67%   | 2         | 20    |
| ✅     | cmd/platform         | 0.00%    | 0         | 3     |
| ✅     | internal             | 88.51%   | 46        | 246   |
| ✅     | internal/assert      | 98.57%   | 67        | 270   |
| ✅     | internal/httpcontext | 100.00%  | 1         | 18    |
| ✅     | internal/pidfile     | 95.83%   | 11        | 44    |

## Functions

| Status | Package              | Function                         | Coverage | Cognitive |
|--------|----------------------|----------------------------------|----------|-----------|
| ✅     | .                    | Error                            | 100.00%  | 1         |
| ✅     |                      | FromContext                      | 100.00%  | 0         |
| ✅     |                      | JSON                             | 100.00%  | 1         |
| ✅     |                      | Manager.Context                  | 100.00%  | 0         |
| ✅     |                      | Manager.Platform                 | 100.00%  | 1         |
| ✅     |                      | Manager.Reload                   | 100.00%  | 4         |
| ✅     |                      | Manager.Start                    | 86.84%   | 17        |
| ✅     |                      | Manager.Stop                     | 96.43%   | 4         |
| ✅     |                      | Manager.URL                      | 100.00%  | 0         |
| ✅     |                      | Manager.Wait                     | 100.00%  | 0         |
| ✅     |                      | Manager.logger                   | 66.67%   | 1         |
| ✅     |                      | Manager.retire                   | 100.00%  | 1         |
| ✅     |                      | Manager.startGeneration          | 93.94%   | 4         |
| ✅     |                      | Manager.url                      | 100.00%  | 1         |
| ✅     |                      | Manager.watch                    | 100.00%  | 1         |
| ✅     |                      | New                              | 88.24%   | 6         |
| ✅     |                      | NewManager                       | 87.50%   | 2         |
| ✅     |                      | NewOptions                       | 100.00%  | 0         |
| ✅     |                      | NewTelemetryModule               | 85.71%   | 2         |
| ✅     |                      | NewTestOptions                   | 100.00%  | 0         |
| ✅     |                      | NewUnimplementedModule           | 100.00%  | 0         |
| ✅     |                      | Options.env                      | 100.00%  | 1         |
| ✅     |                      | Options.envBool                  | 100.00%  | 1         |
| ✅     |                      | Options.envCSV                   | 100.00%  | 1         |
| ✅     |                      | Param                            | 100.00%  | 1         |
| ✅     |                      | Platform.Context                 | 100.00%  | 0         |
| ✅     |                      | Platform.Find                    | 100.00%  | 0         |
| ✅     |                      | Platform.Register                | 100.00%  | 0         |
| ✅     |                      | Platform.Start                   | 96.67%   | 7         |
| ✅     |                      | Platform.Stats                   | 100.00%  | 0         |
| ✅     |                      | Platform.Stop                    | 95.83%   | 7         |
| ✅     |                      | Platform.URL                     | 100.00%  | 0         |
| ✅     |                      | Platform.Use                     | 100.00%  | 0         |
| ✅     |                      | Platform.Wait                    | 100.00%  | 0         |
| ✅     |                      | Platform.logger                  | 100.00%  | 1         |
| ✅     |                      | Platform.observe                 | 100.00%  | 1         |
| ✅     |                      | Platform.setup                   | 100.00%  | 4         |
| ✅     |                      | Platform.setupListener           | 100.00%  | 3         |
| ✅     |                      | Platform.setupRequestContext     | 100.00%  | 0         |
| ✅     |                      | QueryParam                       | 100.00%  | 0         |
| ✅     |                      | RegisterFunc                     | 100.00%  | 0         |
| ✅     |                      | Registry.Cleanup                 | 100.00%  | 0         |
| ✅     |                      | Registry.Clone                   | 100.00%  | 1         |
| ✅     |                      | Registry.Close                   | 100.00%  | 0         |
| ✅     |                      | Registry.Find                    | 100.00%  | 0         |
| ✅     |                      | Registry.Register                | 100.00%  | 0         |
| ✅     |                      | Registry.RegisterFunc            | 100.00%  | 0         |
| ✅     |                      | Registry.Start                   | 97.73%   | 2         |
| ✅     |                      | Registry.Stats                   | 100.00%  | 0         |
| ✅     |                      | Registry.Use                     | 100.00%  | 0         |
| ✅     |                      | Registry.cleanup                 | 100.00%  | 0         |
| ✅     |                      | Registry.close                   | 100.00%  | 2         |
| ✅     |                      | Registry.find                    | 80.95%   | 10        |
| ✅     |                      | Registry.materialize             | 100.00%  | 3         |
| ✅     |                      | Registry.mount                   | 90.91%   | 4         |
| ✅     |                      | Registry.register                | 100.00%  | 0         |
| ✅     |                      | Registry.start                   | 100.00%  | 3         |
| ✅     |                      | Registry.startModule             | 100.00%  | 0         |
| ✅     |                      | Registry.stats                   | 100.00%  | 0         |
| ✅     |                      | Registry.stopModule              | 80.00%   | 3         |
| ✅     |                      | SetupConnections                 | 100.00%  | 0         |
| ✅     |                      | Start                            | 100.00%  | 1         |
| ✅     |                      | TelemetryModule.Middleware       | 100.00%  | 0         |
| ✅     |                      | TelemetryModule.Mount            | 100.00%  | 0         |
| ✅     |                      | TelemetryModule.Tracer           | 100.00%  | 0         |
| ✅     |                      | TestMiddleware                   | 100.00%  | 0         |
| ✅     |                      | Transaction                      | 75.00%   | 3         |
| ✅     |                      | URLParam                         | 100.00%  | 0         |
| ✅     |                      | UnimplementedModule.Mount        | 100.00%  | 1         |
| ✅     |                      | UnimplementedModule.Name         | 66.67%   | 1         |
| ✅     |                      | UnimplementedModule.Start        | 100.00%  | 1         |
| ✅     |                      | UnimplementedModule.Stop         | 66.67%   | 1         |
| ✅     |                      | Use                              | 100.00%  | 0         |
| ✅     |                      | filter                           | 80.95%   | 10        |
| ✅     |                      | generationListener.Accept        | 90.00%   | 4         |
| ✅     |                      | generationListener.Addr          | 100.00%  | 0         |
| ✅     |                      | generationListener.Close         | 100.00%  | 2         |
| ✅     |                      | init                             | 100.00%  | 0         |
| ✅     |                      | listenerURL                      | 85.71%   | 2         |
| ✅     |                      | loggerFromContext                | 100.00%  | 1         |
| ✅     |                      | newSharedListener                | 100.00%  | 0         |
| ✅     |                      | registration.instance            | 100.00%  | 1         |
| ✅     |                      | routePattern                     | 66.67%   | 2         |
| ✅     |                      | setupConnections                 | 100.00%  | 8         |
| ✅     |                      | sharedListener.Close             | 100.00%  | 1         |
| ✅     |                      | sharedListener.handoff           | 100.00%  | 1         |
| ✅     |                      | sharedListener.next              | 100.00%  | 1         |
| ✅     | cmd                  | Main                             | 46.67%   | 2         |
| ✅     | internal             | CountRoutes                      | 100.00%  | 2         |
| ✅     |                      | DatabaseOption.Apply             | 75.00%   | 1         |
| ✅     |                      | DatabaseProvider.Connect         | 85.71%   | 2         |
| ✅     |                      | DatabaseProvider.Open            | 100.00%  | 0         |
| ✅     |                      | DatabaseProvider.Register        | 100.00%  | 0         |
| ✅     |                      | DatabaseProvider.cached          | 94.12%   | 5         |
| ✅     |                      | DatabaseProvider.evict           | 100.00%  | 3         |
| ✅     |                      | DatabaseProvider.parseCredential | 80.00%   | 4         |
| ✅     |                      | DatabaseProvider.with            | 91.67%   | 7         |
| ✅     |                      | NewDatabaseProvider              | 100.00%  | 0         |
| ✅     |                      | PrintRoutes                      | 100.00%  | 0         |
| ✅     |                      | addOption                        | 100.00%  | 1         |
| ✅     |                      | cleanDSN                         | 100.00%  | 3         |
| ✅     |                      | databaseOption                   | 100.00%  | 3         |
| ✅     |                      | handlerName                      | 90.91%   | 7         |
| ✅     |                      | hasOption                        | 66.67%   | 1         |
| ✅     |                      | isSQLiteMemoryDSN                | 100.00%  | 2         |
| ✅     |                      | joinDSN                          | 66.67%   | 1         |
| ✅     |                      | joinOption                       | 100.00%  | 1         |
| ✅     |                      | splitDSN                         | 100.00%  | 1         |
| ✅     | internal/assert      | Assert                           | 100.00%  | 1         |
| ✅     |                      | CheckEquals                      | 100.00%  | 1         |
| ✅     |                      | Contains                         | 100.00%  | 1         |
| ✅     |                      | Empty                            | 100.00%  | 1         |
| ✅     |                      | Equal                            | 100.00%  | 1         |
| ✅     |                      | EqualValues                      | 100.00%  | 1         |
| ✅     |                      | Error                            | 100.00%  | 1         |
| ✅     |                      | ErrorContains                    | 100.00%  | 2         |
| ✅     |                      | ErrorIs                          | 100.00%  | 1         |
| ✅     |                      | Errorf                           | 100.00%  | 0         |
| ✅     |                      | Fail                             | 100.00%  | 1         |
| ✅     |                      | False                            | 100.00%  | 1         |
| ✅     |                      | Greater                          | 100.00%  | 1         |
| ✅     |                      | IsIncreasing                     | 100.00%  | 5         |
| ✅     |                      | Len                              | 100.00%  | 2         |
| ✅     |                      | Nil                              | 100.00%  | 1         |
| ✅     |                      | NoError                          | 100.00%  | 1         |
| ✅     |                      | NotContains                      | 100.00%  | 1         |
| ✅     |                      | NotEmpty                         | 100.00%  | 1         |
| ✅     |                      | NotEqual                         | 100.00%  | 1         |
| ✅     |                      | NotNil                           | 100.00%  | 1         |
| ✅     |                      | ObjectsAreEqualValues            | 100.00%  | 0         |
| ✅     |                      | Run                              | 100.00%  | 10        |
| ✅     |                      | TestSuite.SetT                   | 100.00%  | 0         |
| ✅     |                      | TestSuite.T                      | 100.00%  | 0         |
| ✅     |                      | True                             | 100.00%  | 1         |
| ✅     |                      | compare                          | 80.00%   | 2         |
| ✅     |                      | containsElement                  | 100.00%  | 9         |
| ✅     |                      | hasLength                        | 100.00%  | 1         |
| ✅     |                      | isEmpty                          | 100.00%  | 2         |
| ✅     |                      | isNil                            | 100.00%  | 2         |
| ✅     |                      | orderedCompare                   | 92.86%   | 14        |
| ✅     | internal/httpcontext | NewValue                         | 100.00%  | 0         |
| ✅     |                      | Value[T].Get                     | 100.00%  | 0         |
| ✅     |                      | Value[T].GetContext              | 100.00%  | 1         |
| ✅     |                      | Value[T].Set                     | 100.00%  | 0         |
| ✅     |                      | Value[T].SetContext              | 100.00%  | 0         |
| ✅     | internal/pidfile     | New                              | 100.00%  | 0         |
| ✅     |                      | Pidfile.Remove                   | 85.71%   | 6         |
| ✅     |                      | Pidfile.Write                    | 100.00%  | 2         |
| ✅     |                      | Read                             | 100.00%  | 3         |
