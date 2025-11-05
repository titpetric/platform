# Testing coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 80%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to
cover. Tests with cognitive complexity 0 would be covered by invocation.

## Packages

| Status | Package                                       | Coverage | Cognitive | Lines |
| ------ | --------------------------------------------- | -------- | --------- | ----- |
| ✅      | titpetric/platform                 | 83.47%   | 51        | 439   |
| ✅      | titpetric/platform/cmd             | 87.50%   | 2         | 23    |
| ✅      | titpetric/platform/cmd/platform    | 0.00%    | 0         | 3     |
| ✅      | titpetric/platform/internal        | 93.16%   | 19        | 135   |
| ✅      | titpetric/platform/pkg/assert      | 0.00%    | 0         | 0     |
| ✅      | titpetric/platform/pkg/drivers     | 0.00%    | 0         | 0     |
| ✅      | titpetric/platform/pkg/httpcontext | 100.00%  | 1         | 22    |
| ✅      | titpetric/platform/pkg/reflect     | 100.00%  | 7         | 31    |
| ✅      | titpetric/platform/pkg/require     | 0.00%    | 0         | 0     |
| ❌      | titpetric/platform/pkg/telemetry   | 57.44%   | 8         | 130   |
| ✅      | titpetric/platform/pkg/ulid        | 100.00%  | 0         | 20    |

## Functions

| Status | Package                                       | Function                         | Coverage | Cognitive |
| ------ | --------------------------------------------- | -------------------------------- | -------- | --------- |
| ✅      | titpetric/platform                 | Error                            | 100.00%  | 1         |
| ✅      | titpetric/platform                 | FromContext                      | 0.00%    | 0         |
| ✅      | titpetric/platform                 | FromRequest                      | 0.00%    | 0         |
| ✅      | titpetric/platform                 | JSON                             | 100.00%  | 1         |
| ✅      | titpetric/platform                 | New                              | 83.30%   | 1         |
| ✅      | titpetric/platform                 | NewOptions                       | 100.00%  | 0         |
| ✅      | titpetric/platform                 | NewTestOptions                   | 100.00%  | 0         |
| ✅      | titpetric/platform                 | Options.env                      | 75.00%   | 1         |
| ✅      | titpetric/platform                 | Platform.Context                 | 0.00%    | 0         |
| ✅      | titpetric/platform                 | Platform.Find                    | 100.00%  | 0         |
| ✅      | titpetric/platform                 | Platform.Register                | 100.00%  | 0         |
| ❌      | titpetric/platform                 | Platform.Start                   | 73.30%   | 7         |
| ✅      | titpetric/platform                 | Platform.Stats                   | 100.00%  | 0         |
| ✅      | titpetric/platform                 | Platform.Stop                    | 100.00%  | 0         |
| ✅      | titpetric/platform                 | Platform.URL                     | 100.00%  | 0         |
| ✅      | titpetric/platform                 | Platform.Use                     | 100.00%  | 0         |
| ✅      | titpetric/platform                 | Platform.Wait                    | 100.00%  | 0         |
| ✅      | titpetric/platform                 | Platform.bindContext             | 100.00%  | 0         |
| ✅      | titpetric/platform                 | Platform.setup                   | 71.40%   | 2         |
| ✅      | titpetric/platform                 | Platform.setupListener           | 71.40%   | 2         |
| ✅      | titpetric/platform                 | Platform.setupRequestContext     | 100.00%  | 0         |
| ✅      | titpetric/platform                 | Register                         | 100.00%  | 0         |
| ✅      | titpetric/platform                 | Registry.Clone                   | 100.00%  | 0         |
| ✅      | titpetric/platform                 | Registry.Close                   | 86.70%   | 8         |
| ✅      | titpetric/platform                 | Registry.Find                    | 81.20%   | 8         |
| ✅      | titpetric/platform                 | Registry.Register                | 100.00%  | 0         |
| ✅      | titpetric/platform                 | Registry.Start                   | 80.00%   | 1         |
| ✅      | titpetric/platform                 | Registry.Stats                   | 100.00%  | 0         |
| ✅      | titpetric/platform                 | Registry.Use                     | 100.00%  | 0         |
| ✅      | titpetric/platform                 | Registry.mount                   | 83.30%   | 4         |
| ✅      | titpetric/platform                 | Registry.start                   | 75.00%   | 3         |
| ✅      | titpetric/platform                 | Start                            | 75.00%   | 1         |
| ✅      | titpetric/platform                 | TestMiddleware                   | 100.00%  | 0         |
| ✅      | titpetric/platform                 | Transaction                      | 66.70%   | 3         |
| ✅      | titpetric/platform                 | UnimplementedModule.Mount        | 100.00%  | 1         |
| ✅      | titpetric/platform                 | UnimplementedModule.Name         | 66.70%   | 1         |
| ✅      | titpetric/platform                 | UnimplementedModule.Start        | 66.70%   | 1         |
| ✅      | titpetric/platform                 | UnimplementedModule.Stop         | 66.70%   | 1         |
| ✅      | titpetric/platform                 | Use                              | 100.00%  | 0         |
| ✅      | titpetric/platform                 | init                             | 100.00%  | 0         |
| ✅      | titpetric/platform                 | setupConnections                 | 100.00%  | 4         |
| ✅      | titpetric/platform/cmd             | Main                             | 87.50%   | 2         |
| ✅      | titpetric/platform/cmd/platform    | main                             | 0.00%    | 0         |
| ✅      | titpetric/platform/internal        | CountRoutes                      | 100.00%  | 2         |
| ✅      | titpetric/platform/internal        | DatabaseOption.Apply             | 75.00%   | 1         |
| ✅      | titpetric/platform/internal        | DatabaseProvider.Connect         | 66.70%   | 2         |
| ✅      | titpetric/platform/internal        | DatabaseProvider.Open            | 100.00%  | 0         |
| ✅      | titpetric/platform/internal        | DatabaseProvider.Register        | 100.00%  | 0         |
| ✅      | titpetric/platform/internal        | DatabaseProvider.cached          | 92.90%   | 5         |
| ✅      | titpetric/platform/internal        | DatabaseProvider.parseCredential | 100.00%  | 1         |
| ✅      | titpetric/platform/internal        | DatabaseProvider.with            | 83.30%   | 7         |
| ✅      | titpetric/platform/internal        | NewDatabaseProvider              | 100.00%  | 0         |
| ✅      | titpetric/platform/internal        | PrintRoutes                      | 100.00%  | 0         |
| ✅      | titpetric/platform/internal        | addOptionToDSN                   | 100.00%  | 1         |
| ✅      | titpetric/platform/internal        | cleanDSN                         | 100.00%  | 0         |
| ✅      | titpetric/platform/pkg/httpcontext | NewValue                         | 100.00%  | 0         |
| ✅      | titpetric/platform/pkg/httpcontext | Value[T].Get                     | 100.00%  | 0         |
| ✅      | titpetric/platform/pkg/httpcontext | Value[T].GetContext              | 100.00%  | 1         |
| ✅      | titpetric/platform/pkg/httpcontext | Value[T].Set                     | 100.00%  | 0         |
| ✅      | titpetric/platform/pkg/httpcontext | Value[T].SetContext              | 100.00%  | 0         |
| ✅      | titpetric/platform/pkg/reflect     | SymbolName                       | 100.00%  | 1         |
| ✅      | titpetric/platform/pkg/reflect     | readSymbolName                   | 100.00%  | 6         |
| ✅      | titpetric/platform/pkg/telemetry   | CaptureError                     | 100.00%  | 0         |
| ✅      | titpetric/platform/pkg/telemetry   | Fatal                            | 0.00%    | 0         |
| ✅      | titpetric/platform/pkg/telemetry   | Middleware                       | 0.00%    | 0         |
| ✅      | titpetric/platform/pkg/telemetry   | Monitor.Enabled                  | 100.00%  | 0         |
| ✅      | titpetric/platform/pkg/telemetry   | Monitor.SetEnabled               | 100.00%  | 0         |
| ✅      | titpetric/platform/pkg/telemetry   | Monitor.Touch                    | 100.00%  | 2         |
| ✅      | titpetric/platform/pkg/telemetry   | NewMonitor                       | 100.00%  | 0         |
| ❌      | titpetric/platform/pkg/telemetry   | Open                             | 0.00%    | 1         |
| ✅      | titpetric/platform/pkg/telemetry   | Start                            | 100.00%  | 0         |
| ✅      | titpetric/platform/pkg/telemetry   | StartAuto                        | 0.00%    | 0         |
| ✅      | titpetric/platform/pkg/telemetry   | StartRequest                     | 0.00%    | 0         |
| ✅      | titpetric/platform/pkg/telemetry   | init                             | 66.70%   | 2         |
| ✅      | titpetric/platform/pkg/telemetry   | initOpenTelemetry                | 80.00%   | 3         |
| ✅      | titpetric/platform/pkg/ulid        | Parse                            | 100.00%  | 0         |
| ✅      | titpetric/platform/pkg/ulid        | String                           | 100.00%  | 0         |
| ✅      | titpetric/platform/pkg/ulid        | ULID                             | 100.00%  | 0         |
| ✅      | titpetric/platform/pkg/ulid        | Valid                            | 100.00%  | 0         |

