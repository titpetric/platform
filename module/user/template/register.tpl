<h1>Register</h1>
<form method="POST" action="/register">
    <label>
        First Name:
        <input type="text" name="first_name" required>
    </label><br><br>
    <label>
        Last Name:
        <input type="text" name="last_name" required>
    </label><br><br>
    <label>
        Email:
        <input type="email" name="email" required>
    </label><br><br>
    <label>
        Password:
        <input type="password" name="password" required>
    </label><br><br>
    <button type="submit">Register</button>
</form>
<p>Already have an account? <a href="/login">Login here</a></p>
