<div class="row justify-content-center my-5">
  <div class="col-md-6 col-lg-4">
    <div class="card shadow-sm rounded-4">
      <div class="card-body">
        <h1 class="mb-4">Register</h1>

        <form method="POST" action="/register">
          <table class="table border-0 align-middle">
            <tbody>
              <tr>
                <td class="border-0 w-25">
                  <label for="first_name" class="form-label fw-semibold mb-0">First Name:</label>
                </td>
                <td class="border-0">
                  <input
                    type="text"
                    name="first_name"
                    id="first_name"
                    class="form-control"
                    value="{{.Form.first_name}}"
                    required
                  />
                </td>
              </tr>
              <tr>
                <td class="border-0">
                  <label for="last_name" class="form-label fw-semibold mb-0">Last Name:</label>
                </td>
                <td class="border-0">
                  <input
                    type="text"
                    name="last_name"
                    id="last_name"
                    class="form-control"
                    value="{{.Form.last_name}}"
                    required
                  />
                </td>
              </tr>
              <tr>
                <td class="border-0">
                  <label for="email" class="form-label fw-semibold mb-0">Email:</label>
                </td>
                <td class="border-0">
                  <input
                    type="email"
                    name="email"
                    id="email"
                    class="form-control"
                    required
                    value="{{.Form.email}}"
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
            <button type="submit" class="btn btn-success btn-lg rounded-3">
              Register
            </button>
          </div>
        </form>

        <p class="text-center mt-4 mb-0">
          Already have an account?
          <a href="/login" class="text-decoration-none">Login here</a>
        </p>
      </div>
    </div>
  </div>
</div>
