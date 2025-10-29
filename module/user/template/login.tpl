<div class="row justify-content-center my-5">
  <div class="col-md-6 col-lg-4">
    <div class="card shadow-sm rounded-4">
      <div class="card-body">
        <h1 class="mb-4">Login</h1>

        <form method="POST" action="/login">
          <table class="table border-0 align-middle">
            <tbody>
              <tr>
                <td class="border-0 w-25">
                  <label for="email" class="form-label fw-semibold mb-0">Email:</label>
                </td>
                <td class="border-0">
                  <input
                    type="email"
                    name="email"
                    id="email"
                    class="form-control"
                    value="{{.Form.email}}"
                    required
                  />
                </td>
              </tr>
              <tr>
                <td class="border-0">
                  <label for="password" class="form-label fw-semibold mb-0">Password:</label>
                </td>
                <td class="border-0">
                  <input
                    type="password"
                    name="password"
                    id="password"
                    class="form-control"
                    required
                  />
                </td>
              </tr>
            </tbody>
          </table>

          {{ if .ErrorMessage }}
          <div class="alert alert-warning text-center mt-3">
            {{.ErrorMessage}}
          </div>
          {{ end }}

          <div class="d-grid mt-4">
            <button type="submit" class="btn btn-primary btn-lg rounded-3">
              Login
            </button>
          </div>
        </form>

        <p class="text-center mt-4 mb-0">
          Don't have an account?
          <a href="/register" class="text-decoration-none">Register here</a>
        </p>
      </div>
    </div>
  </div>
</div>
