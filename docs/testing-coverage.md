# Testing coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 80%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to
cover. Tests with cognitive complexity 0 would be tested by invocation.

## Packages

| #  | Status | Package                                           | Coverage | Cognitive | Lines |
| -- | ------ | ------------------------------------------------- | -------- | --------- | ----- |
| 0  | ✅      | titpetric/platform                     | 84.61%   | 34        | 311   |
| 1  | ✅      | titpetric/platform/cmd/platform        | 41.65%   | 1         | 16    |
| 2  | ✅      | titpetric/platform/internal            | 80.66%   | 22        | 192   |
| 3  | ✅      | titpetric/platform/internal/reflect    | 100.00%  | 7         | 31    |
| 4  | ✅      | titpetric/platform/internal/require    | 0.00%    | 0         | 0     |
| 5  | ✅      | titpetric/platform/module/autoload     | 100.00%  | 0         | 4     |
| 6  | ✅      | titpetric/platform/module/expvar       | 95.23%   | 1         | 19    |
| 7  | ✅      | titpetric/platform/module/theme        | 100.00%  | 0         | 11    |
| 8  | ❌      | titpetric/platform/module/user         | 60.41%   | 13        | 103   |
| 9  | ✅      | titpetric/platform/module/user/model   | 8.11%    | 1         | 82    |
| 10 | ❌      | titpetric/platform/module/user/service | 28.25%   | 27        | 285   |
| 11 | ❌      | titpetric/platform/module/user/storage | 20.00%   | 24        | 200   |
| 12 | ❌      | titpetric/platform/telemetry           | 70.90%   | 8         | 130   |

## Functions

| #   | Status | Package                                           | Function                         | Coverage | Cognitive |
| --- | ------ | ------------------------------------------------- | -------------------------------- | -------- | --------- |
| 0   | ✅      | titpetric/platform                     | Error                            | 100.00%  | 1         |
| 1   | ✅      | titpetric/platform                     | JSON                             | 100.00%  | 1         |
| 2   | ✅      | titpetric/platform                     | New                              | 83.30%   | 1         |
| 3   | ✅      | titpetric/platform                     | NewOptions                       | 100.00%  | 0         |
| 4   | ✅      | titpetric/platform                     | NewTestOptions                   | 100.00%  | 0         |
| 5   | ✅      | titpetric/platform                     | Platform.Context                 | 0.00%    | 0         |
| 6   | ✅      | titpetric/platform                     | Platform.Register                | 100.00%  | 0         |
| 7   | ❌      | titpetric/platform                     | Platform.Start                   | 73.30%   | 7         |
| 8   | ✅      | titpetric/platform                     | Platform.Stats                   | 100.00%  | 0         |
| 9   | ✅      | titpetric/platform                     | Platform.Stop                    | 100.00%  | 0         |
| 10  | ✅      | titpetric/platform                     | Platform.URL                     | 0.00%    | 0         |
| 11  | ✅      | titpetric/platform                     | Platform.Use                     | 100.00%  | 0         |
| 12  | ✅      | titpetric/platform                     | Platform.Wait                    | 100.00%  | 0         |
| 13  | ✅      | titpetric/platform                     | Platform.setup                   | 66.70%   | 2         |
| 14  | ✅      | titpetric/platform                     | Platform.setupListener           | 71.40%   | 2         |
| 15  | ✅      | titpetric/platform                     | Register                         | 100.00%  | 0         |
| 16  | ✅      | titpetric/platform                     | Registry.Clone                   | 100.00%  | 0         |
| 17  | ✅      | titpetric/platform                     | Registry.Close                   | 86.70%   | 8         |
| 18  | ✅      | titpetric/platform                     | Registry.Register                | 100.00%  | 0         |
| 19  | ✅      | titpetric/platform                     | Registry.Start                   | 81.80%   | 7         |
| 20  | ✅      | titpetric/platform                     | Registry.Stats                   | 100.00%  | 0         |
| 21  | ✅      | titpetric/platform                     | Registry.Use                     | 100.00%  | 0         |
| 22  | ✅      | titpetric/platform                     | Start                            | 75.00%   | 1         |
| 23  | ✅      | titpetric/platform                     | UnimplementedModule.Mount        | 100.00%  | 0         |
| 24  | ✅      | titpetric/platform                     | UnimplementedModule.Name         | 0.00%    | 0         |
| 25  | ✅      | titpetric/platform                     | UnimplementedModule.Start        | 100.00%  | 0         |
| 26  | ✅      | titpetric/platform                     | UnimplementedModule.Stop         | 100.00%  | 0         |
| 27  | ✅      | titpetric/platform                     | Use                              | 100.00%  | 0         |
| 28  | ✅      | titpetric/platform                     | init                             | 100.00%  | 0         |
| 29  | ✅      | titpetric/platform                     | setupConnections                 | 100.00%  | 4         |
| 30  | ✅      | titpetric/platform/cmd/platform        | main                             | 0.00%    | 0         |
| 31  | ✅      | titpetric/platform/cmd/platform        | start                            | 83.30%   | 1         |
| 32  | ❌      | titpetric/platform/internal            | ContextValue[T].Get              | 0.00%    | 1         |
| 33  | ✅      | titpetric/platform/internal            | ContextValue[T].Set              | 0.00%    | 0         |
| 34  | ✅      | titpetric/platform/internal            | CountRoutes                      | 100.00%  | 2         |
| 35  | ✅      | titpetric/platform/internal            | DatabaseOption.Apply             | 75.00%   | 1         |
| 36  | ✅      | titpetric/platform/internal            | DatabaseProvider.Connect         | 75.00%   | 1         |
| 37  | ✅      | titpetric/platform/internal            | DatabaseProvider.Open            | 100.00%  | 0         |
| 38  | ✅      | titpetric/platform/internal            | DatabaseProvider.Register        | 100.00%  | 0         |
| 39  | ✅      | titpetric/platform/internal            | DatabaseProvider.cached          | 100.00%  | 5         |
| 40  | ✅      | titpetric/platform/internal            | DatabaseProvider.parseCredential | 100.00%  | 1         |
| 41  | ✅      | titpetric/platform/internal            | DatabaseProvider.with            | 83.30%   | 7         |
| 42  | ✅      | titpetric/platform/internal            | NewContextValue                  | 100.00%  | 0         |
| 43  | ✅      | titpetric/platform/internal            | NewDatabaseProvider              | 100.00%  | 0         |
| 44  | ✅      | titpetric/platform/internal            | NewServer                        | 100.00%  | 0         |
| 45  | ✅      | titpetric/platform/internal            | NewTemplate                      | 100.00%  | 0         |
| 46  | ✅      | titpetric/platform/internal            | PrintRoutes                      | 100.00%  | 0         |
| 47  | ❌      | titpetric/platform/internal            | Template.Render                  | 0.00%    | 1         |
| 48  | ✅      | titpetric/platform/internal            | Transaction                      | 80.00%   | 2         |
| 49  | ✅      | titpetric/platform/internal            | ULID                             | 100.00%  | 0         |
| 50  | ✅      | titpetric/platform/internal            | addOptionToDSN                   | 100.00%  | 1         |
| 51  | ✅      | titpetric/platform/internal            | cleanDSN                         | 100.00%  | 0         |
| 52  | ✅      | titpetric/platform/internal/reflect    | SymbolName                       | 100.00%  | 1         |
| 53  | ✅      | titpetric/platform/internal/reflect    | readSymbolName                   | 100.00%  | 6         |
| 54  | ✅      | titpetric/platform/module/autoload     | init                             | 100.00%  | 0         |
| 55  | ✅      | titpetric/platform/module/expvar       | Handler.Mount                    | 100.00%  | 0         |
| 56  | ✅      | titpetric/platform/module/expvar       | Handler.Start                    | 85.70%   | 1         |
| 57  | ✅      | titpetric/platform/module/expvar       | NewHandler                       | 100.00%  | 0         |
| 58  | ✅      | titpetric/platform/module/theme        | NewOptions                       | 100.00%  | 0         |
| 59  | ✅      | titpetric/platform/module/user         | DB                               | 100.00%  | 0         |
| 60  | ❌      | titpetric/platform/module/user         | GetSessionUser                   | 0.00%    | 8         |
| 61  | ✅      | titpetric/platform/module/user         | Handler.Mount                    | 100.00%  | 0         |
| 62  | ✅      | titpetric/platform/module/user         | Handler.Name                     | 0.00%    | 0         |
| 63  | ✅      | titpetric/platform/module/user         | Handler.Start                    | 83.30%   | 2         |
| 64  | ✅      | titpetric/platform/module/user         | Handler.Stop                     | 100.00%  | 0         |
| 65  | ❌      | titpetric/platform/module/user         | IsLoggedIn                       | 0.00%    | 3         |
| 66  | ✅      | titpetric/platform/module/user         | NewHandler                       | 100.00%  | 0         |
| 67  | ✅      | titpetric/platform/module/user/model   | NewUser                          | 100.00%  | 0         |
| 68  | ✅      | titpetric/platform/module/user/model   | NewUserGroup                     | 0.00%    | 0         |
| 69  | ✅      | titpetric/platform/module/user/model   | User.GetCreatedAt                | 0.00%    | 0         |
| 70  | ✅      | titpetric/platform/module/user/model   | User.GetDeletedAt                | 0.00%    | 0         |
| 71  | ✅      | titpetric/platform/module/user/model   | User.GetFirstName                | 0.00%    | 0         |
| 72  | ✅      | titpetric/platform/module/user/model   | User.GetID                       | 0.00%    | 0         |
| 73  | ✅      | titpetric/platform/module/user/model   | User.GetLastName                 | 0.00%    | 0         |
| 74  | ✅      | titpetric/platform/module/user/model   | User.GetUpdatedAt                | 0.00%    | 0         |
| 75  | ✅      | titpetric/platform/module/user/model   | User.IsActive                    | 0.00%    | 0         |
| 76  | ✅      | titpetric/platform/module/user/model   | User.SetCreatedAt                | 0.00%    | 0         |
| 77  | ✅      | titpetric/platform/module/user/model   | User.SetDeletedAt                | 100.00%  | 0         |
| 78  | ✅      | titpetric/platform/module/user/model   | User.SetUpdatedAt                | 0.00%    | 0         |
| 79  | ✅      | titpetric/platform/module/user/model   | User.String                      | 100.00%  | 1         |
| 80  | ✅      | titpetric/platform/module/user/model   | UserAuth.GetCreatedAt            | 0.00%    | 0         |
| 81  | ✅      | titpetric/platform/module/user/model   | UserAuth.GetEmail                | 0.00%    | 0         |
| 82  | ✅      | titpetric/platform/module/user/model   | UserAuth.GetPassword             | 0.00%    | 0         |
| 83  | ✅      | titpetric/platform/module/user/model   | UserAuth.GetUpdatedAt            | 0.00%    | 0         |
| 84  | ✅      | titpetric/platform/module/user/model   | UserAuth.GetUserID               | 0.00%    | 0         |
| 85  | ✅      | titpetric/platform/module/user/model   | UserAuth.SetCreatedAt            | 0.00%    | 0         |
| 86  | ✅      | titpetric/platform/module/user/model   | UserAuth.SetUpdatedAt            | 0.00%    | 0         |
| 87  | ✅      | titpetric/platform/module/user/model   | UserGroup.GetCreatedAt           | 0.00%    | 0         |
| 88  | ✅      | titpetric/platform/module/user/model   | UserGroup.GetID                  | 0.00%    | 0         |
| 89  | ✅      | titpetric/platform/module/user/model   | UserGroup.GetTitle               | 0.00%    | 0         |
| 90  | ✅      | titpetric/platform/module/user/model   | UserGroup.GetUpdatedAt           | 0.00%    | 0         |
| 91  | ✅      | titpetric/platform/module/user/model   | UserGroup.SetCreatedAt           | 0.00%    | 0         |
| 92  | ✅      | titpetric/platform/module/user/model   | UserGroup.SetUpdatedAt           | 0.00%    | 0         |
| 93  | ✅      | titpetric/platform/module/user/model   | UserGroup.String                 | 0.00%    | 0         |
| 94  | ✅      | titpetric/platform/module/user/model   | UserGroupMember.GetJoinedAt      | 0.00%    | 0         |
| 95  | ✅      | titpetric/platform/module/user/model   | UserGroupMember.GetUserGroupID   | 0.00%    | 0         |
| 96  | ✅      | titpetric/platform/module/user/model   | UserGroupMember.GetUserID        | 0.00%    | 0         |
| 97  | ✅      | titpetric/platform/module/user/model   | UserGroupMember.SetJoinedAt      | 0.00%    | 0         |
| 98  | ✅      | titpetric/platform/module/user/model   | UserSession.GetCreatedAt         | 0.00%    | 0         |
| 99  | ✅      | titpetric/platform/module/user/model   | UserSession.GetExpiresAt         | 0.00%    | 0         |
| 100 | ✅      | titpetric/platform/module/user/model   | UserSession.GetID                | 0.00%    | 0         |
| 101 | ✅      | titpetric/platform/module/user/model   | UserSession.GetUserID            | 0.00%    | 0         |
| 102 | ✅      | titpetric/platform/module/user/model   | UserSession.SetCreatedAt         | 0.00%    | 0         |
| 103 | ✅      | titpetric/platform/module/user/model   | UserSession.SetExpiresAt         | 0.00%    | 0         |
| 104 | ✅      | titpetric/platform/module/user/service | NewService                       | 75.00%   | 1         |
| 105 | ✅      | titpetric/platform/module/user/service | Service.Close                    | 100.00%  | 0         |
| 106 | ❌      | titpetric/platform/module/user/service | Service.Error                    | 0.00%    | 2         |
| 107 | ✅      | titpetric/platform/module/user/service | Service.GetError                 | 0.00%    | 0         |
| 108 | ❌      | titpetric/platform/module/user/service | Service.Login                    | 0.00%    | 5         |
| 109 | ❌      | titpetric/platform/module/user/service | Service.LoginView                | 0.00%    | 8         |
| 110 | ❌      | titpetric/platform/module/user/service | Service.Logout                   | 0.00%    | 2         |
| 111 | ✅      | titpetric/platform/module/user/service | Service.LogoutView               | 0.00%    | 0         |
| 112 | ✅      | titpetric/platform/module/user/service | Service.Mount                    | 100.00%  | 0         |
| 113 | ❌      | titpetric/platform/module/user/service | Service.Register                 | 0.00%    | 4         |
| 114 | ✅      | titpetric/platform/module/user/service | Service.RegisterView             | 0.00%    | 0         |
| 115 | ❌      | titpetric/platform/module/user/service | Service.View                     | 0.00%    | 3         |
| 116 | ✅      | titpetric/platform/module/user/service | Service.initTemplates            | 92.30%   | 2         |
| 117 | ✅      | titpetric/platform/module/user/storage | NewSessionStorage                | 100.00%  | 0         |
| 118 | ✅      | titpetric/platform/module/user/storage | NewUserStorage                   | 100.00%  | 0         |
| 119 | ❌      | titpetric/platform/module/user/storage | SessionStorage.Create            | 0.00%    | 1         |
| 120 | ✅      | titpetric/platform/module/user/storage | SessionStorage.Delete            | 0.00%    | 0         |
| 121 | ❌      | titpetric/platform/module/user/storage | SessionStorage.Get               | 0.00%    | 4         |
| 122 | ❌      | titpetric/platform/module/user/storage | UserStorage.Authenticate         | 0.00%    | 8         |
| 123 | ❌      | titpetric/platform/module/user/storage | UserStorage.Create               | 0.00%    | 8         |
| 124 | ❌      | titpetric/platform/module/user/storage | UserStorage.Get                  | 0.00%    | 1         |
| 125 | ❌      | titpetric/platform/module/user/storage | UserStorage.GetGroups            | 0.00%    | 1         |
| 126 | ❌      | titpetric/platform/module/user/storage | UserStorage.Update               | 0.00%    | 1         |
| 127 | ✅      | titpetric/platform/telemetry           | CaptureError                     | 100.00%  | 0         |
| 128 | ✅      | titpetric/platform/telemetry           | Fatal                            | 0.00%    | 0         |
| 129 | ✅      | titpetric/platform/telemetry           | Middleware                       | 100.00%  | 0         |
| 130 | ✅      | titpetric/platform/telemetry           | Monitor.Enabled                  | 100.00%  | 0         |
| 131 | ✅      | titpetric/platform/telemetry           | Monitor.SetEnabled               | 100.00%  | 0         |
| 132 | ✅      | titpetric/platform/telemetry           | Monitor.Touch                    | 100.00%  | 2         |
| 133 | ✅      | titpetric/platform/telemetry           | NewMonitor                       | 100.00%  | 0         |
| 134 | ✅      | titpetric/platform/telemetry           | Open                             | 75.00%   | 1         |
| 135 | ✅      | titpetric/platform/telemetry           | Start                            | 100.00%  | 0         |
| 136 | ✅      | titpetric/platform/telemetry           | StartAuto                        | 0.00%    | 0         |
| 137 | ✅      | titpetric/platform/telemetry           | StartRequest                     | 0.00%    | 0         |
| 138 | ✅      | titpetric/platform/telemetry           | init                             | 66.70%   | 2         |
| 139 | ✅      | titpetric/platform/telemetry           | initOpenTelemetry                | 80.00%   | 3         |

