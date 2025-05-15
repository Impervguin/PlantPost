function handleRegister(event: Event): void {
    event.preventDefault();
    
    const formData = {
        email: (document.getElementById('email') as HTMLInputElement).value,
        username: (document.getElementById('username') as HTMLInputElement).value,
        password: (document.getElementById('password') as HTMLInputElement).value,
        passwordConfirm: (document.getElementById('password-confirm') as HTMLInputElement).value
    };

    if (formData.password !== formData.passwordConfirm) {
        alert('Passwords do not match');
        return;
    }

    if (!formData.email.includes('@')) {
        alert('Invalid email');
        return;
    }

    const requestData = {
        email: formData.email,
        username: formData.username,
        password: formData.password
    };

    fetch('/api/auth/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
    })
    .then(response => {
        if (response.ok) {
            window.location.href = '/view/login'; // Redirect on success
        } else {
            return response.json().then(errorData => {
                throw new Error(errorData.message || 'Register failed');
            });
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert(error.message || 'An error occurred during login');
    });
}

document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', handleRegister);
    }
});