# Testing coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 80%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to
cover. Tests with cognitive complexity 0 would be covered by invocation.

The storage package has integration tests behind a build tag. To run
integration tests, run `task integration`.

## Packages

| Status | Package                                           | Coverage | Cognitive | Lines |
| ------ | ------------------------------------------------- | -------- | --------- | ----- |
| ✅      | titpetric/platform                     | 84.61%   | 34        | 311   |
| ✅      | titpetric/platform/cmd/platform        | 41.65%   | 1         | 16    |
| ✅      | titpetric/platform/drivers             | 0.00%    | 0         | 0     |
| ❌      | titpetric/platform/internal            | 79.59%   | 24        | 203   |
| ✅      | titpetric/platform/internal/reflect    | 100.00%  | 7         | 31    |
| ✅      | titpetric/platform/internal/require    | 0.00%    | 0         | 0     |
| ✅      | titpetric/platform/module/autoload     | 100.00%  | 0         | 4     |
| ✅      | titpetric/platform/module/expvar       | 95.23%   | 1         | 19    |
| ✅      | titpetric/platform/module/theme        | 100.00%  | 0         | 11    |
| ❌      | titpetric/platform/module/user         | 54.76%   | 13        | 98    |
| ✅      | titpetric/platform/module/user/model   | 10.53%   | 3         | 88    |
| ❌      | titpetric/platform/module/user/service | 28.25%   | 27        | 285   |
| ❌      | titpetric/platform/module/user/storage | 44.01%   | 24        | 211   |
| ❌      | titpetric/platform/telemetry           | 78.59%   | 8         | 130   |

## Functions

| Status | Package                                           | Function                         | Coverage | Cognitive |
| ------ | ------------------------------------------------- | -------------------------------- | -------- | --------- |
| ✅      | titpetric/platform                     | Error                            | 100.00%  | 1         |
| ✅      | titpetric/platform                     | JSON                             | 100.00%  | 1         |
| ✅      | titpetric/platform                     | New                              | 83.30%   | 1         |
| ✅      | titpetric/platform                     | NewOptions                       | 100.00%  | 0         |
| ✅      | titpetric/platform                     | NewTestOptions                   | 100.00%  | 0         |
| ✅      | titpetric/platform                     | Platform.Context                 | 0.00%    | 0         |
| ✅      | titpetric/platform                     | Platform.Register                | 100.00%  | 0         |
| ❌      | titpetric/platform                     | Platform.Start                   | 73.30%   | 7         |
| ✅      | titpetric/platform                     | Platform.Stats                   | 100.00%  | 0         |
| ✅      | titpetric/platform                     | Platform.Stop                    | 100.00%  | 0         |
| ✅      | titpetric/platform                     | Platform.URL                     | 0.00%    | 0         |
| ✅      | titpetric/platform                     | Platform.Use                     | 100.00%  | 0         |
| ✅      | titpetric/platform                     | Platform.Wait                    | 100.00%  | 0         |
| ✅      | titpetric/platform                     | Platform.setup                   | 66.70%   | 2         |
| ✅      | titpetric/platform                     | Platform.setupListener           | 71.40%   | 2         |
| ✅      | titpetric/platform                     | Register                         | 100.00%  | 0         |
| ✅      | titpetric/platform                     | Registry.Clone                   | 100.00%  | 0         |
| ✅      | titpetric/platform                     | Registry.Close                   | 86.70%   | 8         |
| ✅      | titpetric/platform                     | Registry.Register                | 100.00%  | 0         |
| ✅      | titpetric/platform                     | Registry.Start                   | 81.80%   | 7         |
| ✅      | titpetric/platform                     | Registry.Stats                   | 100.00%  | 0         |
| ✅      | titpetric/platform                     | Registry.Use                     | 100.00%  | 0         |
| ✅      | titpetric/platform                     | Start                            | 75.00%   | 1         |
| ✅      | titpetric/platform                     | UnimplementedModule.Mount        | 100.00%  | 0         |
| ✅      | titpetric/platform                     | UnimplementedModule.Name         | 0.00%    | 0         |
| ✅      | titpetric/platform                     | UnimplementedModule.Start        | 100.00%  | 0         |
| ✅      | titpetric/platform                     | UnimplementedModule.Stop         | 100.00%  | 0         |
| ✅      | titpetric/platform                     | Use                              | 100.00%  | 0         |
| ✅      | titpetric/platform                     | init                             | 100.00%  | 0         |
| ✅      | titpetric/platform                     | setupConnections                 | 100.00%  | 4         |
| ✅      | titpetric/platform/cmd/platform        | main                             | 0.00%    | 0         |
| ✅      | titpetric/platform/cmd/platform        | start                            | 83.30%   | 1         |
| ❌      | titpetric/platform/internal            | ContextValue[T].Get              | 0.00%    | 1         |
| ✅      | titpetric/platform/internal            | ContextValue[T].Set              | 0.00%    | 0         |
| ✅      | titpetric/platform/internal            | CountRoutes                      | 100.00%  | 2         |
| ✅      | titpetric/platform/internal            | DatabaseOption.Apply             | 75.00%   | 1         |
| ✅      | titpetric/platform/internal            | DatabaseProvider.Connect         | 66.70%   | 2         |
| ✅      | titpetric/platform/internal            | DatabaseProvider.Open            | 100.00%  | 0         |
| ✅      | titpetric/platform/internal            | DatabaseProvider.Register        | 100.00%  | 0         |
| ✅      | titpetric/platform/internal            | DatabaseProvider.cached          | 100.00%  | 5         |
| ✅      | titpetric/platform/internal            | DatabaseProvider.parseCredential | 100.00%  | 1         |
| ✅      | titpetric/platform/internal            | DatabaseProvider.with            | 83.30%   | 7         |
| ✅      | titpetric/platform/internal            | NewContextValue                  | 100.00%  | 0         |
| ✅      | titpetric/platform/internal            | NewDatabaseProvider              | 100.00%  | 0         |
| ✅      | titpetric/platform/internal            | NewServer                        | 100.00%  | 0         |
| ✅      | titpetric/platform/internal            | NewTemplate                      | 100.00%  | 0         |
| ✅      | titpetric/platform/internal            | PrintRoutes                      | 100.00%  | 0         |
| ❌      | titpetric/platform/internal            | Template.Render                  | 0.00%    | 1         |
| ✅      | titpetric/platform/internal            | Transaction                      | 66.70%   | 3         |
| ✅      | titpetric/platform/internal            | ULID                             | 100.00%  | 0         |
| ✅      | titpetric/platform/internal            | addOptionToDSN                   | 100.00%  | 1         |
| ✅      | titpetric/platform/internal            | cleanDSN                         | 100.00%  | 0         |
| ✅      | titpetric/platform/internal/reflect    | SymbolName                       | 100.00%  | 1         |
| ✅      | titpetric/platform/internal/reflect    | readSymbolName                   | 100.00%  | 6         |
| ✅      | titpetric/platform/module/autoload     | init                             | 100.00%  | 0         |
| ✅      | titpetric/platform/module/expvar       | Handler.Mount                    | 100.00%  | 0         |
| ✅      | titpetric/platform/module/expvar       | Handler.Start                    | 85.70%   | 1         |
| ✅      | titpetric/platform/module/expvar       | NewHandler                       | 100.00%  | 0         |
| ✅      | titpetric/platform/module/theme        | NewOptions                       | 100.00%  | 0         |
| ❌      | titpetric/platform/module/user         | GetSessionUser                   | 0.00%    | 8         |
| ✅      | titpetric/platform/module/user         | Handler.Mount                    | 100.00%  | 0         |
| ✅      | titpetric/platform/module/user         | Handler.Name                     | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user         | Handler.Start                    | 83.30%   | 2         |
| ✅      | titpetric/platform/module/user         | Handler.Stop                     | 100.00%  | 0         |
| ❌      | titpetric/platform/module/user         | IsLoggedIn                       | 0.00%    | 3         |
| ✅      | titpetric/platform/module/user         | NewHandler                       | 100.00%  | 0         |
| ✅      | titpetric/platform/module/user/model   | NewUser                          | 100.00%  | 0         |
| ✅      | titpetric/platform/module/user/model   | NewUserGroup                     | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | User.GetCreatedAt                | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | User.GetDeletedAt                | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | User.GetFirstName                | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | User.GetID                       | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | User.GetLastName                 | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | User.GetUpdatedAt                | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | User.IsActive                    | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | User.SetCreatedAt                | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | User.SetDeletedAt                | 100.00%  | 0         |
| ✅      | titpetric/platform/module/user/model   | User.SetUpdatedAt                | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | User.String                      | 100.00%  | 1         |
| ✅      | titpetric/platform/module/user/model   | UserAuth.GetCreatedAt            | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserAuth.GetEmail                | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserAuth.GetPassword             | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserAuth.GetUpdatedAt            | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserAuth.GetUserID               | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserAuth.SetCreatedAt            | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserAuth.SetUpdatedAt            | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserAuth.Valid                   | 100.00%  | 2         |
| ✅      | titpetric/platform/module/user/model   | UserGroup.GetCreatedAt           | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserGroup.GetID                  | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserGroup.GetTitle               | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserGroup.GetUpdatedAt           | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserGroup.SetCreatedAt           | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserGroup.SetUpdatedAt           | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserGroup.String                 | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserGroupMember.GetJoinedAt      | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserGroupMember.GetUserGroupID   | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserGroupMember.GetUserID        | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserGroupMember.SetJoinedAt      | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserSession.GetCreatedAt         | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserSession.GetExpiresAt         | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserSession.GetID                | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserSession.GetUserID            | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserSession.SetCreatedAt         | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/model   | UserSession.SetExpiresAt         | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/service | NewService                       | 75.00%   | 1         |
| ✅      | titpetric/platform/module/user/service | Service.Close                    | 100.00%  | 0         |
| ❌      | titpetric/platform/module/user/service | Service.Error                    | 0.00%    | 2         |
| ✅      | titpetric/platform/module/user/service | Service.GetError                 | 0.00%    | 0         |
| ❌      | titpetric/platform/module/user/service | Service.Login                    | 0.00%    | 5         |
| ❌      | titpetric/platform/module/user/service | Service.LoginView                | 0.00%    | 8         |
| ❌      | titpetric/platform/module/user/service | Service.Logout                   | 0.00%    | 2         |
| ✅      | titpetric/platform/module/user/service | Service.LogoutView               | 0.00%    | 0         |
| ✅      | titpetric/platform/module/user/service | Service.Mount                    | 100.00%  | 0         |
| ❌      | titpetric/platform/module/user/service | Service.Register                 | 0.00%    | 4         |
| ✅      | titpetric/platform/module/user/service | Service.RegisterView             | 0.00%    | 0         |
| ❌      | titpetric/platform/module/user/service | Service.View                     | 0.00%    | 3         |
| ✅      | titpetric/platform/module/user/service | Service.initTemplates            | 92.30%   | 2         |
| ✅      | titpetric/platform/module/user/storage | DB                               | 100.00%  | 0         |
| ✅      | titpetric/platform/module/user/storage | NewSessionStorage                | 100.00%  | 0         |
| ✅      | titpetric/platform/module/user/storage | NewUserStorage                   | 100.00%  | 0         |
| ❌      | titpetric/platform/module/user/storage | SessionStorage.Create            | 0.00%    | 1         |
| ✅      | titpetric/platform/module/user/storage | SessionStorage.Delete            | 85.70%   | 1         |
| ✅      | titpetric/platform/module/user/storage | SessionStorage.Get               | 63.60%   | 4         |
| ❌      | titpetric/platform/module/user/storage | UserStorage.Authenticate         | 34.80%   | 8         |
| ❌      | titpetric/platform/module/user/storage | UserStorage.Create               | 0.00%    | 7         |
| ❌      | titpetric/platform/module/user/storage | UserStorage.Get                  | 0.00%    | 1         |
| ❌      | titpetric/platform/module/user/storage | UserStorage.GetGroups            | 0.00%    | 1         |
| ❌      | titpetric/platform/module/user/storage | UserStorage.Update               | 0.00%    | 1         |
| ✅      | titpetric/platform/telemetry           | CaptureError                     | 100.00%  | 0         |
| ✅      | titpetric/platform/telemetry           | Fatal                            | 0.00%    | 0         |
| ✅      | titpetric/platform/telemetry           | Middleware                       | 100.00%  | 0         |
| ✅      | titpetric/platform/telemetry           | Monitor.Enabled                  | 100.00%  | 0         |
| ✅      | titpetric/platform/telemetry           | Monitor.SetEnabled               | 100.00%  | 0         |
| ✅      | titpetric/platform/telemetry           | Monitor.Touch                    | 100.00%  | 2         |
| ✅      | titpetric/platform/telemetry           | NewMonitor                       | 100.00%  | 0         |
| ✅      | titpetric/platform/telemetry           | Open                             | 75.00%   | 1         |
| ✅      | titpetric/platform/telemetry           | Start                            | 100.00%  | 0         |
| ✅      | titpetric/platform/telemetry           | StartAuto                        | 100.00%  | 0         |
| ✅      | titpetric/platform/telemetry           | StartRequest                     | 0.00%    | 0         |
| ✅      | titpetric/platform/telemetry           | init                             | 66.70%   | 2         |
| ✅      | titpetric/platform/telemetry           | initOpenTelemetry                | 80.00%   | 3         |

