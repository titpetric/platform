<div class="row justify-content-center my-5">
  <div class="col-md-6 col-lg-4">
    <div class="card shadow-sm rounded-4">
      <div class="card-body">
        <h1 class="mb-4">Welcome, {{.User}}</h1>

        <form method="POST" action="/logout">
          <div class="d-grid mt-3">
            <button type="submit" class="btn btn-danger btn-lg rounded-3">
              Logout
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</div>
