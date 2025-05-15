function handleLogin(event: Event): void {
    event.preventDefault();
    
    const formData = {
        username: (document.getElementById('username') as HTMLInputElement).value,
        password: (document.getElementById('password') as HTMLInputElement).value
    };

    fetch('/api/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
        credentials: 'include' // Important for cookies
    })
    .then(response => {
        if (response.ok) {
            window.location.href = '/view'; // Redirect on success
        } else {
            return response.json().then(errorData => {
                throw new Error(errorData.message || 'Login failed');
            });
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert(error.message || 'An error occurred during login');
    });
}

document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }
});