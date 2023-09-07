// Function to make a login API call
async function login(loginData, redirect, client_id) {
    const queryParams = new URLSearchParams({ return_to: redirect, client_id: client_id });
    const response = await fetch(`http://127.0.0.1:3000/api/v1/auth/login?${queryParams.toString()}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(loginData)
    });

    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Login failed');
    }

    const data = await response.json();
    return data; // Assuming the token is returned in the API response
}

// Function to handle user login
async function handleLogin(loginData, redirect, client_id) {
    try {
        const data = await login(loginData, redirect, client_id);
        setCookie(data, 7); // Save token in cookie for 7 days
        window.location.href = redirect;
    } catch (error) {
        console.error('Login error:', error.message);
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
        let isData = false;

        for(let i = 0; i < parsedData.length; i++){
            if(parsedData[i].user.email == user.user.email){
                parsedData[i] = user;
                localStorage.setItem('userData', JSON.stringify(parsedData));
                isData = true;
                break;
            }
        }

        if(!isData){
            parsedData.push(user);
            localStorage.setItem('userData', JSON.stringify(parsedData));
        }
    }
}
