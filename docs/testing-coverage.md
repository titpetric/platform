# Testing coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 80%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to
cover. Tests with cognitive complexity 0 would be covered by invocation.

## Packages

| Status | Package         | Coverage | Cognitive | Lines |
|--------|-----------------|----------|-----------|-------|
| ✅     | .               | 89.48%   | 122       | 900   |
| ✅     | cmd             | 46.70%   | 2         | 20    |
| ✅     | cmd/platform    | 0.00%    | 0         | 3     |
| ✅     | internal        | 81.12%   | 31        | 176   |
| ✅     | pkg/assert      | 0.00%    | 0         | 0     |
| ✅     | pkg/drivers     | 0.00%    | 0         | 0     |
| ✅     | pkg/httpcontext | 100.00%  | 1         | 18    |
| ✅     | pkg/require     | 0.00%    | 0         | 0     |
| ✅     | pkg/ulid        | 100.00%  | 0         | 13    |

## Functions

| Status | Package         | Function                         | Coverage | Cognitive |
|--------|-----------------|----------------------------------|----------|-----------|
| ✅     | .               | Error                            | 100.00%  | 1         |
| ✅     |                 | FromContext                      | 100.00%  | 0         |
| ✅     |                 | JSON                             | 100.00%  | 1         |
| ✅     |                 | Manager.Platform                 | 100.00%  | 1         |
| ✅     |                 | Manager.Reload                   | 100.00%  | 1         |
| ✅     |                 | Manager.Start                    | 83.90%   | 11        |
| ✅     |                 | Manager.Stop                     | 100.00%  | 2         |
| ✅     |                 | Manager.URL                      | 100.00%  | 0         |
| ✅     |                 | Manager.Wait                     | 100.00%  | 0         |
| ✅     |                 | Manager.logger                   | 66.70%   | 1         |
| ✅     |                 | Manager.retire                   | 100.00%  | 1         |
| ✅     |                 | Manager.startGeneration          | 92.30%   | 4         |
| ✅     |                 | Manager.watch                    | 100.00%  | 1         |
| ✅     |                 | New                              | 88.20%   | 6         |
| ✅     |                 | NewManager                       | 87.50%   | 2         |
| ✅     |                 | NewOptions                       | 100.00%  | 0         |
| ✅     |                 | NewTelemetryModule               | 85.70%   | 2         |
| ✅     |                 | NewTestOptions                   | 100.00%  | 0         |
| ✅     |                 | NewUnimplementedModule           | 100.00%  | 0         |
| ✅     |                 | Options.env                      | 100.00%  | 1         |
| ✅     |                 | Options.envBool                  | 100.00%  | 1         |
| ✅     |                 | Options.envCSV                   | 100.00%  | 1         |
| ✅     |                 | Param                            | 100.00%  | 1         |
| ✅     |                 | Platform.Context                 | 100.00%  | 0         |
| ✅     |                 | Platform.Find                    | 100.00%  | 0         |
| ✅     |                 | Platform.Register                | 100.00%  | 0         |
| ✅     |                 | Platform.Start                   | 96.40%   | 6         |
| ✅     |                 | Platform.Stats                   | 100.00%  | 0         |
| ✅     |                 | Platform.Stop                    | 100.00%  | 4         |
| ✅     |                 | Platform.URL                     | 100.00%  | 0         |
| ✅     |                 | Platform.Use                     | 100.00%  | 0         |
| ✅     |                 | Platform.Wait                    | 100.00%  | 0         |
| ✅     |                 | Platform.logger                  | 100.00%  | 1         |
| ✅     |                 | Platform.observe                 | 100.00%  | 1         |
| ✅     |                 | Platform.setup                   | 100.00%  | 4         |
| ✅     |                 | Platform.setupListener           | 100.00%  | 3         |
| ✅     |                 | Platform.setupRequestContext     | 100.00%  | 0         |
| ✅     |                 | QueryParam                       | 100.00%  | 0         |
| ✅     |                 | RegisterFunc                     | 100.00%  | 0         |
| ✅     |                 | Registry.Cleanup                 | 100.00%  | 0         |
| ✅     |                 | Registry.Clone                   | 100.00%  | 1         |
| ✅     |                 | Registry.Close                   | 100.00%  | 0         |
| ✅     |                 | Registry.Find                    | 86.20%   | 10        |
| ✅     |                 | Registry.Register                | 100.00%  | 0         |
| ✅     |                 | Registry.RegisterFunc            | 100.00%  | 0         |
| ✅     |                 | Registry.Start                   | 95.70%   | 3         |
| ✅     |                 | Registry.Stats                   | 100.00%  | 0         |
| ✅     |                 | Registry.Use                     | 100.00%  | 0         |
| ✅     |                 | Registry.close                   | 100.00%  | 3         |
| ✅     |                 | Registry.filter                  | 81.00%   | 10        |
| ✅     |                 | Registry.materialize             | 100.00%  | 3         |
| ✅     |                 | Registry.mount                   | 90.90%   | 4         |
| ✅     |                 | Registry.register                | 100.00%  | 0         |
| ✅     |                 | Registry.start                   | 100.00%  | 3         |
| ✅     |                 | Registry.startModule             | 100.00%  | 0         |
| ✅     |                 | Registry.stopModule              | 80.00%   | 3         |
| ✅     |                 | SetupConnections                 | 100.00%  | 0         |
| ✅     |                 | Start                            | 75.00%   | 1         |
| ✅     |                 | TelemetryModule.Middleware       | 100.00%  | 0         |
| ✅     |                 | TelemetryModule.Mount            | 100.00%  | 0         |
| ✅     |                 | TelemetryModule.Tracer           | 100.00%  | 0         |
| ✅     |                 | TestMiddleware                   | 100.00%  | 0         |
| ✅     |                 | Transaction                      | 75.00%   | 3         |
| ✅     |                 | URLParam                         | 100.00%  | 0         |
| ✅     |                 | UnimplementedModule.Mount        | 100.00%  | 1         |
| ✅     |                 | UnimplementedModule.Name         | 66.70%   | 1         |
| ✅     |                 | UnimplementedModule.Start        | 100.00%  | 1         |
| ✅     |                 | UnimplementedModule.Stop         | 66.70%   | 1         |
| ✅     |                 | Use                              | 100.00%  | 0         |
| ✅     |                 | generationListener.Accept        | 90.00%   | 4         |
| ✅     |                 | generationListener.Addr          | 100.00%  | 0         |
| ✅     |                 | generationListener.Close         | 100.00%  | 2         |
| ✅     |                 | listenerURL                      | 100.00%  | 0         |
| ✅     |                 | loggerFromContext                | 100.00%  | 1         |
| ✅     |                 | newSharedListener                | 100.00%  | 0         |
| ✅     |                 | registration.instance            | 100.00%  | 1         |
| ✅     |                 | routePattern                     | 66.70%   | 2         |
| ✅     |                 | setupConnections                 | 100.00%  | 4         |
| ✅     |                 | sharedListener.Close             | 100.00%  | 1         |
| ✅     |                 | sharedListener.handoff           | 100.00%  | 1         |
| ✅     |                 | sharedListener.next              | 100.00%  | 1         |
| ✅     | cmd             | Main                             | 46.70%   | 2         |
| ✅     | internal        | CountRoutes                      | 100.00%  | 2         |
| ✅     |                 | DatabaseOption.Apply             | 75.00%   | 1         |
| ✅     |                 | DatabaseProvider.Connect         | 66.70%   | 2         |
| ✅     |                 | DatabaseProvider.Open            | 100.00%  | 0         |
| ✅     |                 | DatabaseProvider.Register        | 100.00%  | 0         |
| ✅     |                 | DatabaseProvider.cached          | 92.90%   | 5         |
| ✅     |                 | DatabaseProvider.parseCredential | 80.00%   | 4         |
| ✅     |                 | DatabaseProvider.with            | 83.30%   | 7         |
| ✅     |                 | NewDatabaseProvider              | 100.00%  | 0         |
| ✅     |                 | PrintRoutes                      | 100.00%  | 0         |
| ✅     |                 | addOptionToDSN                   | 100.00%  | 1         |
| ✅     |                 | cleanDSN                         | 100.00%  | 3         |
| ✅     |                 | databaseOption                   | 100.00%  | 2         |
| ✅     |                 | isSQLiteMemoryDSN                | 100.00%  | 2         |
| ✅     | pkg/httpcontext | NewValue                         | 100.00%  | 0         |
| ✅     |                 | Value[T].Get                     | 100.00%  | 0         |
| ✅     |                 | Value[T].GetContext              | 100.00%  | 1         |
| ✅     |                 | Value[T].Set                     | 100.00%  | 0         |
| ✅     |                 | Value[T].SetContext              | 100.00%  | 0         |
| ✅     | pkg/ulid        | Parse                            | 100.00%  | 0         |
| ✅     |                 | String                           | 100.00%  | 0         |
| ✅     |                 | ULID                             | 100.00%  | 0         |
| ✅     |                 | Valid                            | 100.00%  | 0         |
