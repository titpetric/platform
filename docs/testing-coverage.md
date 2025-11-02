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
| 0  | ✅      | github.com/titpetric/platform                     | 82.01%   | 34        | 304   |
| 1  | ✅      | github.com/titpetric/platform/cmd/platform        | 41.65%   | 1         | 16    |
| 2  | ❌      | github.com/titpetric/platform/internal            | 74.31%   | 27        | 183   |
| 3  | ✅      | github.com/titpetric/platform/internal/reflect    | 100.00%  | 7         | 31    |
| 4  | ✅      | github.com/titpetric/platform/internal/require    | 0.00%    | 0         | 0     |
| 5  | ✅      | github.com/titpetric/platform/module/autoload     | 100.00%  | 0         | 4     |
| 6  | ✅      | github.com/titpetric/platform/module/expvar       | 93.33%   | 1         | 16    |
| 7  | ✅      | github.com/titpetric/platform/module/theme        | 100.00%  | 0         | 11    |
| 8  | ❌      | github.com/titpetric/platform/module/user         | 60.41%   | 13        | 103   |
| 9  | ✅      | github.com/titpetric/platform/module/user/model   | 8.11%    | 1         | 82    |
| 10 | ❌      | github.com/titpetric/platform/module/user/service | 28.25%   | 27        | 285   |
| 11 | ❌      | github.com/titpetric/platform/module/user/storage | 20.00%   | 24        | 200   |
| 12 | ❌      | github.com/titpetric/platform/telemetry           | 70.90%   | 8         | 130   |

## Functions

| #   | Status | Package                                           | Function                       | Coverage | Cognitive |
| --- | ------ | ------------------------------------------------- | ------------------------------ | -------- | --------- |
| 0   | ✅      | github.com/titpetric/platform                     | Error                          | 100.00%  | 1         |
| 1   | ✅      | github.com/titpetric/platform                     | JSON                           | 100.00%  | 1         |
| 2   | ✅      | github.com/titpetric/platform                     | New                            | 83.30%   | 1         |
| 3   | ✅      | github.com/titpetric/platform                     | NewOptions                     | 100.00%  | 0         |
| 4   | ✅      | github.com/titpetric/platform                     | NewTestOptions                 | 100.00%  | 0         |
| 5   | ✅      | github.com/titpetric/platform                     | Platform.Context               | 0.00%    | 0         |
| 6   | ✅      | github.com/titpetric/platform                     | Platform.Register              | 100.00%  | 0         |
| 7   | ❌      | github.com/titpetric/platform                     | Platform.Start                 | 69.60%   | 11        |
| 8   | ✅      | github.com/titpetric/platform                     | Platform.Stats                 | 100.00%  | 0         |
| 9   | ✅      | github.com/titpetric/platform                     | Platform.Stop                  | 100.00%  | 0         |
| 10  | ✅      | github.com/titpetric/platform                     | Platform.URL                   | 0.00%    | 0         |
| 11  | ✅      | github.com/titpetric/platform                     | Platform.Use                   | 0.00%    | 0         |
| 12  | ✅      | github.com/titpetric/platform                     | Platform.Wait                  | 100.00%  | 0         |
| 13  | ✅      | github.com/titpetric/platform                     | Register                       | 100.00%  | 0         |
| 14  | ✅      | github.com/titpetric/platform                     | Registry.Clone                 | 100.00%  | 0         |
| 15  | ✅      | github.com/titpetric/platform                     | Registry.Close                 | 86.70%   | 8         |
| 16  | ✅      | github.com/titpetric/platform                     | Registry.Register              | 100.00%  | 0         |
| 17  | ✅      | github.com/titpetric/platform                     | Registry.Start                 | 81.80%   | 7         |
| 18  | ✅      | github.com/titpetric/platform                     | Registry.Stats                 | 100.00%  | 0         |
| 19  | ✅      | github.com/titpetric/platform                     | Registry.Use                   | 100.00%  | 0         |
| 20  | ✅      | github.com/titpetric/platform                     | Start                          | 75.00%   | 1         |
| 21  | ✅      | github.com/titpetric/platform                     | UnimplementedModule.Mount      | 100.00%  | 0         |
| 22  | ✅      | github.com/titpetric/platform                     | UnimplementedModule.Name       | 0.00%    | 0         |
| 23  | ✅      | github.com/titpetric/platform                     | UnimplementedModule.Start      | 100.00%  | 0         |
| 24  | ✅      | github.com/titpetric/platform                     | UnimplementedModule.Stop       | 100.00%  | 0         |
| 25  | ✅      | github.com/titpetric/platform                     | Use                            | 100.00%  | 0         |
| 26  | ✅      | github.com/titpetric/platform                     | init                           | 100.00%  | 0         |
| 27  | ✅      | github.com/titpetric/platform                     | setupConnections               | 100.00%  | 4         |
| 28  | ✅      | github.com/titpetric/platform/cmd/platform        | main                           | 0.00%    | 0         |
| 29  | ✅      | github.com/titpetric/platform/cmd/platform        | start                          | 83.30%   | 1         |
| 30  | ❌      | github.com/titpetric/platform/internal            | ContextValue[T].Get            | 0.00%    | 1         |
| 31  | ✅      | github.com/titpetric/platform/internal            | ContextValue[T].Set            | 0.00%    | 0         |
| 32  | ✅      | github.com/titpetric/platform/internal            | CountRoutes                    | 100.00%  | 2         |
| 33  | ✅      | github.com/titpetric/platform/internal            | DatabaseProvider.Connect       | 75.00%   | 1         |
| 34  | ✅      | github.com/titpetric/platform/internal            | DatabaseProvider.Open          | 100.00%  | 0         |
| 35  | ✅      | github.com/titpetric/platform/internal            | DatabaseProvider.Register      | 100.00%  | 0         |
| 36  | ✅      | github.com/titpetric/platform/internal            | DatabaseProvider.cached        | 100.00%  | 5         |
| 37  | ✅      | github.com/titpetric/platform/internal            | DatabaseProvider.with          | 88.20%   | 13        |
| 38  | ✅      | github.com/titpetric/platform/internal            | NewContextValue                | 100.00%  | 0         |
| 39  | ✅      | github.com/titpetric/platform/internal            | NewDatabaseProvider            | 100.00%  | 0         |
| 40  | ✅      | github.com/titpetric/platform/internal            | NewTemplate                    | 100.00%  | 0         |
| 41  | ✅      | github.com/titpetric/platform/internal            | PrintRoutes                    | 100.00%  | 0         |
| 42  | ❌      | github.com/titpetric/platform/internal            | Template.Render                | 0.00%    | 1         |
| 43  | ❌      | github.com/titpetric/platform/internal            | Transaction                    | 0.00%    | 3         |
| 44  | ✅      | github.com/titpetric/platform/internal            | ULID                           | 100.00%  | 0         |
| 45  | ✅      | github.com/titpetric/platform/internal            | addOptionToDSN                 | 100.00%  | 1         |
| 46  | ✅      | github.com/titpetric/platform/internal            | cleanDSN                       | 100.00%  | 0         |
| 47  | ✅      | github.com/titpetric/platform/internal/reflect    | SymbolName                     | 100.00%  | 1         |
| 48  | ✅      | github.com/titpetric/platform/internal/reflect    | readSymbolName                 | 100.00%  | 6         |
| 49  | ✅      | github.com/titpetric/platform/module/autoload     | init                           | 100.00%  | 0         |
| 50  | ✅      | github.com/titpetric/platform/module/expvar       | Handler.Mount                  | 100.00%  | 0         |
| 51  | ✅      | github.com/titpetric/platform/module/expvar       | Handler.Start                  | 80.00%   | 1         |
| 52  | ✅      | github.com/titpetric/platform/module/expvar       | NewHandler                     | 100.00%  | 0         |
| 53  | ✅      | github.com/titpetric/platform/module/theme        | NewOptions                     | 100.00%  | 0         |
| 54  | ✅      | github.com/titpetric/platform/module/user         | DB                             | 100.00%  | 0         |
| 55  | ❌      | github.com/titpetric/platform/module/user         | GetSessionUser                 | 0.00%    | 8         |
| 56  | ✅      | github.com/titpetric/platform/module/user         | Handler.Mount                  | 100.00%  | 0         |
| 57  | ✅      | github.com/titpetric/platform/module/user         | Handler.Name                   | 0.00%    | 0         |
| 58  | ✅      | github.com/titpetric/platform/module/user         | Handler.Start                  | 83.30%   | 2         |
| 59  | ✅      | github.com/titpetric/platform/module/user         | Handler.Stop                   | 100.00%  | 0         |
| 60  | ❌      | github.com/titpetric/platform/module/user         | IsLoggedIn                     | 0.00%    | 3         |
| 61  | ✅      | github.com/titpetric/platform/module/user         | NewHandler                     | 100.00%  | 0         |
| 62  | ✅      | github.com/titpetric/platform/module/user/model   | NewUser                        | 100.00%  | 0         |
| 63  | ✅      | github.com/titpetric/platform/module/user/model   | NewUserGroup                   | 0.00%    | 0         |
| 64  | ✅      | github.com/titpetric/platform/module/user/model   | User.GetCreatedAt              | 0.00%    | 0         |
| 65  | ✅      | github.com/titpetric/platform/module/user/model   | User.GetDeletedAt              | 0.00%    | 0         |
| 66  | ✅      | github.com/titpetric/platform/module/user/model   | User.GetFirstName              | 0.00%    | 0         |
| 67  | ✅      | github.com/titpetric/platform/module/user/model   | User.GetID                     | 0.00%    | 0         |
| 68  | ✅      | github.com/titpetric/platform/module/user/model   | User.GetLastName               | 0.00%    | 0         |
| 69  | ✅      | github.com/titpetric/platform/module/user/model   | User.GetUpdatedAt              | 0.00%    | 0         |
| 70  | ✅      | github.com/titpetric/platform/module/user/model   | User.IsActive                  | 0.00%    | 0         |
| 71  | ✅      | github.com/titpetric/platform/module/user/model   | User.SetCreatedAt              | 0.00%    | 0         |
| 72  | ✅      | github.com/titpetric/platform/module/user/model   | User.SetDeletedAt              | 100.00%  | 0         |
| 73  | ✅      | github.com/titpetric/platform/module/user/model   | User.SetUpdatedAt              | 0.00%    | 0         |
| 74  | ✅      | github.com/titpetric/platform/module/user/model   | User.String                    | 100.00%  | 1         |
| 75  | ✅      | github.com/titpetric/platform/module/user/model   | UserAuth.GetCreatedAt          | 0.00%    | 0         |
| 76  | ✅      | github.com/titpetric/platform/module/user/model   | UserAuth.GetEmail              | 0.00%    | 0         |
| 77  | ✅      | github.com/titpetric/platform/module/user/model   | UserAuth.GetPassword           | 0.00%    | 0         |
| 78  | ✅      | github.com/titpetric/platform/module/user/model   | UserAuth.GetUpdatedAt          | 0.00%    | 0         |
| 79  | ✅      | github.com/titpetric/platform/module/user/model   | UserAuth.GetUserID             | 0.00%    | 0         |
| 80  | ✅      | github.com/titpetric/platform/module/user/model   | UserAuth.SetCreatedAt          | 0.00%    | 0         |
| 81  | ✅      | github.com/titpetric/platform/module/user/model   | UserAuth.SetUpdatedAt          | 0.00%    | 0         |
| 82  | ✅      | github.com/titpetric/platform/module/user/model   | UserGroup.GetCreatedAt         | 0.00%    | 0         |
| 83  | ✅      | github.com/titpetric/platform/module/user/model   | UserGroup.GetID                | 0.00%    | 0         |
| 84  | ✅      | github.com/titpetric/platform/module/user/model   | UserGroup.GetTitle             | 0.00%    | 0         |
| 85  | ✅      | github.com/titpetric/platform/module/user/model   | UserGroup.GetUpdatedAt         | 0.00%    | 0         |
| 86  | ✅      | github.com/titpetric/platform/module/user/model   | UserGroup.SetCreatedAt         | 0.00%    | 0         |
| 87  | ✅      | github.com/titpetric/platform/module/user/model   | UserGroup.SetUpdatedAt         | 0.00%    | 0         |
| 88  | ✅      | github.com/titpetric/platform/module/user/model   | UserGroup.String               | 0.00%    | 0         |
| 89  | ✅      | github.com/titpetric/platform/module/user/model   | UserGroupMember.GetJoinedAt    | 0.00%    | 0         |
| 90  | ✅      | github.com/titpetric/platform/module/user/model   | UserGroupMember.GetUserGroupID | 0.00%    | 0         |
| 91  | ✅      | github.com/titpetric/platform/module/user/model   | UserGroupMember.GetUserID      | 0.00%    | 0         |
| 92  | ✅      | github.com/titpetric/platform/module/user/model   | UserGroupMember.SetJoinedAt    | 0.00%    | 0         |
| 93  | ✅      | github.com/titpetric/platform/module/user/model   | UserSession.GetCreatedAt       | 0.00%    | 0         |
| 94  | ✅      | github.com/titpetric/platform/module/user/model   | UserSession.GetExpiresAt       | 0.00%    | 0         |
| 95  | ✅      | github.com/titpetric/platform/module/user/model   | UserSession.GetID              | 0.00%    | 0         |
| 96  | ✅      | github.com/titpetric/platform/module/user/model   | UserSession.GetUserID          | 0.00%    | 0         |
| 97  | ✅      | github.com/titpetric/platform/module/user/model   | UserSession.SetCreatedAt       | 0.00%    | 0         |
| 98  | ✅      | github.com/titpetric/platform/module/user/model   | UserSession.SetExpiresAt       | 0.00%    | 0         |
| 99  | ✅      | github.com/titpetric/platform/module/user/service | NewService                     | 75.00%   | 1         |
| 100 | ✅      | github.com/titpetric/platform/module/user/service | Service.Close                  | 100.00%  | 0         |
| 101 | ❌      | github.com/titpetric/platform/module/user/service | Service.Error                  | 0.00%    | 2         |
| 102 | ✅      | github.com/titpetric/platform/module/user/service | Service.GetError               | 0.00%    | 0         |
| 103 | ❌      | github.com/titpetric/platform/module/user/service | Service.Login                  | 0.00%    | 5         |
| 104 | ❌      | github.com/titpetric/platform/module/user/service | Service.LoginView              | 0.00%    | 8         |
| 105 | ❌      | github.com/titpetric/platform/module/user/service | Service.Logout                 | 0.00%    | 2         |
| 106 | ✅      | github.com/titpetric/platform/module/user/service | Service.LogoutView             | 0.00%    | 0         |
| 107 | ✅      | github.com/titpetric/platform/module/user/service | Service.Mount                  | 100.00%  | 0         |
| 108 | ❌      | github.com/titpetric/platform/module/user/service | Service.Register               | 0.00%    | 4         |
| 109 | ✅      | github.com/titpetric/platform/module/user/service | Service.RegisterView           | 0.00%    | 0         |
| 110 | ❌      | github.com/titpetric/platform/module/user/service | Service.View                   | 0.00%    | 3         |
| 111 | ✅      | github.com/titpetric/platform/module/user/service | Service.initTemplates          | 92.30%   | 2         |
| 112 | ✅      | github.com/titpetric/platform/module/user/storage | NewSessionStorage              | 100.00%  | 0         |
| 113 | ✅      | github.com/titpetric/platform/module/user/storage | NewUserStorage                 | 100.00%  | 0         |
| 114 | ❌      | github.com/titpetric/platform/module/user/storage | SessionStorage.Create          | 0.00%    | 1         |
| 115 | ✅      | github.com/titpetric/platform/module/user/storage | SessionStorage.Delete          | 0.00%    | 0         |
| 116 | ❌      | github.com/titpetric/platform/module/user/storage | SessionStorage.Get             | 0.00%    | 4         |
| 117 | ❌      | github.com/titpetric/platform/module/user/storage | UserStorage.Authenticate       | 0.00%    | 8         |
| 118 | ❌      | github.com/titpetric/platform/module/user/storage | UserStorage.Create             | 0.00%    | 8         |
| 119 | ❌      | github.com/titpetric/platform/module/user/storage | UserStorage.Get                | 0.00%    | 1         |
| 120 | ❌      | github.com/titpetric/platform/module/user/storage | UserStorage.GetGroups          | 0.00%    | 1         |
| 121 | ❌      | github.com/titpetric/platform/module/user/storage | UserStorage.Update             | 0.00%    | 1         |
| 122 | ✅      | github.com/titpetric/platform/telemetry           | CaptureError                   | 100.00%  | 0         |
| 123 | ✅      | github.com/titpetric/platform/telemetry           | Fatal                          | 0.00%    | 0         |
| 124 | ✅      | github.com/titpetric/platform/telemetry           | Middleware                     | 100.00%  | 0         |
| 125 | ✅      | github.com/titpetric/platform/telemetry           | Monitor.Enabled                | 100.00%  | 0         |
| 126 | ✅      | github.com/titpetric/platform/telemetry           | Monitor.SetEnabled             | 100.00%  | 0         |
| 127 | ✅      | github.com/titpetric/platform/telemetry           | Monitor.Touch                  | 100.00%  | 2         |
| 128 | ✅      | github.com/titpetric/platform/telemetry           | NewMonitor                     | 100.00%  | 0         |
| 129 | ✅      | github.com/titpetric/platform/telemetry           | Open                           | 75.00%   | 1         |
| 130 | ✅      | github.com/titpetric/platform/telemetry           | Start                          | 100.00%  | 0         |
| 131 | ✅      | github.com/titpetric/platform/telemetry           | StartAuto                      | 0.00%    | 0         |
| 132 | ✅      | github.com/titpetric/platform/telemetry           | StartRequest                   | 0.00%    | 0         |
| 133 | ✅      | github.com/titpetric/platform/telemetry           | init                           | 66.70%   | 2         |
| 134 | ✅      | github.com/titpetric/platform/telemetry           | initOpenTelemetry              | 80.00%   | 3         |

