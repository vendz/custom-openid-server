async function sso(userData, redirect, client_id) {
    const queryParams = new URLSearchParams({ return_to: redirect, client_id: client_id });
    const response = await fetch(`http://127.0.0.1:3000/api/v1/sso?${queryParams.toString()}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${userData.token}`
        },
    });

    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Login failed');
    }

    const data = await response.json();
    return data; // Assuming the token is returned in the API response
}

async function handleSSO(data, redirect, client_id) {
    try {
        const userData = await sso(data, redirect, client_id);
        setCookie(userData, 7); // Save token in cookie for 7 days
        window.location.href = redirect;
        console.log('Login success:', userData);
    } catch (error) {
        console.error('Login error:', error.message);
    }
}

function setCookie(data, days) {
    const expires = new Date(Date.now() + days * 24 * 60 * 60 * 1000).toUTCString();
    document.cookie = `token=${encodeURIComponent(data.access_token)}; expires=${expires}; path=/`;
}