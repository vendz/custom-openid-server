// Function to make a signup API call
async function signup(userData, redirect, client_id) {
    const queryParams = new URLSearchParams({ return_to: redirect, client_id: client_id });
    const response = await fetch(`http://127.0.0.1:3000/api/v1/auth/createUser?${queryParams.toString()}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(userData)
    });

    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Signup failed');
    }

    const data = await response.json();
    return data;
}

// Function to handle user signup
async function handleSignup(userData, redirect, client_id) {
    try {
        const data = await signup(userData, redirect, client_id);
        setCookie(data, 7); // Save token in cookie for 7 days
        window.location.href = redirect;
    } catch (error) {
        console.error('Signup error:', error);
    }
}

function setCookie(data, days) {
    const expires = new Date(Date.now() + days * 24 * 60 * 60 * 1000).toUTCString();
    document.cookie = `token=${encodeURIComponent(data.access_token)}; expires=${expires}; path=/`;

    const user = {
        user: data.user,
        token: data.token
    };

    const localData = localStorage.getItem('userData');
    if(localData == null){
        let userData = [];
        userData.push(user);
        localStorage.setItem('userData', JSON.stringify(userData));
    } else {
        const parsedData = JSON.parse(localData);
        if(parsedData.find(user => user.user.email == email) == null){
            parsedData.push(user);
            localStorage.setItem('userData', JSON.stringify(parsedData));
        } else {
            parsedData.find(user => user.user.email == email).user = data.user;
            localStorage.setItem('userData', JSON.stringify(parsedData));
        }
    }
}