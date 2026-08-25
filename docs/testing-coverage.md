# Testing coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 80%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to
cover. Tests with cognitive complexity 0 would be covered by invocation.

## Packages

| Status | Package                                       | Coverage | Cognitive | Lines |
|--------|-----------------------------------------------|----------|-----------|-------|
| ✅     | github.com/titpetric/platform                 | 89.25%   | 122       | 1041  |
| ✅     | github.com/titpetric/platform/cmd             | 46.70%   | 2         | 30    |
| ✅     | github.com/titpetric/platform/cmd/platform    | 0.00%    | 0         | 3     |
| ✅     | github.com/titpetric/platform/internal        | 81.12%   | 31        | 195   |
| ✅     | github.com/titpetric/platform/pkg/assert      | 0.00%    | 0         | 0     |
| ✅     | github.com/titpetric/platform/pkg/drivers     | 0.00%    | 0         | 0     |
| ✅     | github.com/titpetric/platform/pkg/httpcontext | 100.00%  | 1         | 23    |
| ✅     | github.com/titpetric/platform/pkg/require     | 0.00%    | 0         | 0     |
| ✅     | github.com/titpetric/platform/pkg/ulid        | 100.00%  | 0         | 18    |

## Functions

| Status | Package                                       | Function                         | Coverage | Cognitive |
|--------|-----------------------------------------------|----------------------------------|----------|-----------|
| ✅     | github.com/titpetric/platform                 | Error                            | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | FromContext                      | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | JSON                             | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | Manager.Platform                 | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | Manager.Reload                   | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | Manager.Start                    | 83.90%   | 11        |
| ✅     | github.com/titpetric/platform                 | Manager.Stop                     | 100.00%  | 2         |
| ✅     | github.com/titpetric/platform                 | Manager.URL                      | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Manager.Wait                     | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Manager.logger                   | 66.70%   | 1         |
| ✅     | github.com/titpetric/platform                 | Manager.retire                   | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | Manager.startGeneration          | 92.30%   | 4         |
| ✅     | github.com/titpetric/platform                 | Manager.watch                    | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | New                              | 88.20%   | 6         |
| ✅     | github.com/titpetric/platform                 | NewManager                       | 87.50%   | 2         |
| ✅     | github.com/titpetric/platform                 | NewOptions                       | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | NewTelemetryModule               | 90.00%   | 2         |
| ✅     | github.com/titpetric/platform                 | NewTestOptions                   | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | NewUnimplementedModule           | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Options.env                      | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | Options.envBool                  | 75.00%   | 1         |
| ✅     | github.com/titpetric/platform                 | Options.envCSV                   | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | Param                            | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | Platform.Context                 | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Platform.Find                    | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Platform.Register                | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Platform.Start                   | 96.40%   | 6         |
| ✅     | github.com/titpetric/platform                 | Platform.Stats                   | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Platform.Stop                    | 100.00%  | 4         |
| ✅     | github.com/titpetric/platform                 | Platform.URL                     | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Platform.Use                     | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Platform.Wait                    | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Platform.logger                  | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | Platform.observe                 | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | Platform.setup                   | 100.00%  | 4         |
| ✅     | github.com/titpetric/platform                 | Platform.setupListener           | 100.00%  | 3         |
| ✅     | github.com/titpetric/platform                 | Platform.setupRequestContext     | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | QueryParam                       | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | RegisterFunc                     | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Registry.Cleanup                 | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Registry.Clone                   | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | Registry.Close                   | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Registry.Find                    | 86.20%   | 10        |
| ✅     | github.com/titpetric/platform                 | Registry.Register                | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Registry.RegisterFunc            | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Registry.Start                   | 95.70%   | 3         |
| ✅     | github.com/titpetric/platform                 | Registry.Stats                   | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Registry.Use                     | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Registry.close                   | 100.00%  | 3         |
| ✅     | github.com/titpetric/platform                 | Registry.filter                  | 81.00%   | 10        |
| ✅     | github.com/titpetric/platform                 | Registry.materialize             | 100.00%  | 3         |
| ✅     | github.com/titpetric/platform                 | Registry.mount                   | 90.90%   | 4         |
| ✅     | github.com/titpetric/platform                 | Registry.register                | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Registry.start                   | 100.00%  | 3         |
| ✅     | github.com/titpetric/platform                 | Registry.startModule             | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Registry.stopModule              | 80.00%   | 3         |
| ✅     | github.com/titpetric/platform                 | SetupConnections                 | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Start                            | 75.00%   | 1         |
| ✅     | github.com/titpetric/platform                 | TelemetryModule.Middleware       | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | TelemetryModule.Mount            | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | TelemetryModule.Tracer           | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | TestMiddleware                   | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | Transaction                      | 75.00%   | 3         |
| ✅     | github.com/titpetric/platform                 | URLParam                         | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | UnimplementedModule.Mount        | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | UnimplementedModule.Name         | 66.70%   | 1         |
| ✅     | github.com/titpetric/platform                 | UnimplementedModule.Start        | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | UnimplementedModule.Stop         | 66.70%   | 1         |
| ✅     | github.com/titpetric/platform                 | Use                              | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | generationListener.Accept        | 90.00%   | 4         |
| ✅     | github.com/titpetric/platform                 | generationListener.Addr          | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | generationListener.Close         | 100.00%  | 2         |
| ✅     | github.com/titpetric/platform                 | listenerURL                      | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | loggerFromContext                | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | newSharedListener                | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform                 | registration.instance            | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | routePattern                     | 66.70%   | 2         |
| ✅     | github.com/titpetric/platform                 | setupConnections                 | 100.00%  | 4         |
| ✅     | github.com/titpetric/platform                 | sharedListener.Close             | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | sharedListener.handoff           | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform                 | sharedListener.next              | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform/cmd             | Main                             | 46.70%   | 2         |
| ✅     | github.com/titpetric/platform/internal        | CountRoutes                      | 100.00%  | 2         |
| ✅     | github.com/titpetric/platform/internal        | DatabaseOption.Apply             | 75.00%   | 1         |
| ✅     | github.com/titpetric/platform/internal        | DatabaseProvider.Connect         | 66.70%   | 2         |
| ✅     | github.com/titpetric/platform/internal        | DatabaseProvider.Open            | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform/internal        | DatabaseProvider.Register        | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform/internal        | DatabaseProvider.cached          | 92.90%   | 5         |
| ✅     | github.com/titpetric/platform/internal        | DatabaseProvider.parseCredential | 80.00%   | 4         |
| ✅     | github.com/titpetric/platform/internal        | DatabaseProvider.with            | 83.30%   | 7         |
| ✅     | github.com/titpetric/platform/internal        | NewDatabaseProvider              | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform/internal        | PrintRoutes                      | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform/internal        | addOptionToDSN                   | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform/internal        | cleanDSN                         | 100.00%  | 3         |
| ✅     | github.com/titpetric/platform/internal        | databaseOption                   | 100.00%  | 2         |
| ✅     | github.com/titpetric/platform/internal        | isSQLiteMemoryDSN                | 100.00%  | 2         |
| ✅     | github.com/titpetric/platform/pkg/httpcontext | NewValue                         | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform/pkg/httpcontext | Value[T].Get                     | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform/pkg/httpcontext | Value[T].GetContext              | 100.00%  | 1         |
| ✅     | github.com/titpetric/platform/pkg/httpcontext | Value[T].Set                     | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform/pkg/httpcontext | Value[T].SetContext              | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform/pkg/ulid        | Parse                            | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform/pkg/ulid        | String                           | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform/pkg/ulid        | ULID                             | 100.00%  | 0         |
| ✅     | github.com/titpetric/platform/pkg/ulid        | Valid                            | 100.00%  | 0         |
