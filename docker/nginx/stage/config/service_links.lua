ngx.say([[
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Service Links - Login</title>
    <style>
        /* General styles */
        body {
            font-family: 'Arial', sans-serif;
            margin: 0;
            padding: 0;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
            background: linear-gradient(135deg, #6a11cb, #2575fc);
            color: #fff;
        }

        .container {
            background: rgba(255, 255, 255, 0.1);
            padding: 2rem;
            border-radius: 15px;
            backdrop-filter: blur(10px);
            box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
            width: 100%;
            max-width: 400px;
            animation: fadeIn 1.2s ease-in-out;
        }

        h1 {
            font-size: 2.5rem;
            margin-bottom: 20px;
            text-align: center;
        }

        .form-group {
            margin-bottom: 20px;
        }

        label {
            display: block;
            margin-bottom: 8px;
            font-weight: bold;
        }

        input {
            width: 100%;
            padding: 12px;
            border: none;
            border-radius: 8px;
            background: rgba(255, 255, 255, 0.2);
            color: #fff;
            font-size: 16px;
            box-sizing: border-box;
        }

        input::placeholder {
            color: rgba(255, 255, 255, 0.7);
        }

        .button {
            display: inline-block;
            width: 100%;
            padding: 15px;
            font-size: 1rem;
            font-weight: bold;
            text-align: center;
            text-decoration: none;
            color: #fff;
            background: linear-gradient(135deg, #ff7eb3, #ff758c);
            border: none;
            border-radius: 30px;
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
            transition: all 0.4s ease;
            cursor: pointer;
        }

        .button:hover {
            background: linear-gradient(135deg, #ff758c, #ff7eb3);
            transform: scale(1.02);
            box-shadow: 0 8px 20px rgba(0, 0, 0, 0.3);
        }

        .error-message {
            color: #ff4444;
            text-align: center;
            margin-top: 10px;
            display: none;
        }

        #services-container {
            display: none;
        }

        .button-container {
            display: flex;
            flex-wrap: wrap;
            justify-content: center;
            gap: 15px;
            animation: slideIn 1s ease-in-out;
        }

        /* Animations */
        @keyframes fadeIn {
            from {
                opacity: 0;
                transform: translateY(-20px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        @keyframes slideIn {
            from {
                opacity: 0;
                transform: translateX(-30px);
            }
            to {
                opacity: 1;
                transform: translateX(0);
            }
        }
    </style>
</head>
<body>
    <div id="login-container" class="container">
        <h1>Login</h1>
        <form id="login-form">
            <div class="form-group">
                <label for="username">Username</label>
                <input type="text" id="username" placeholder="Enter username" required>
            </div>
            <div class="form-group">
                <label for="password">Password</label>
                <input type="password" id="password" placeholder="Enter password" required>
            </div>
            <button type="submit" class="button">Login</button>
            <div id="error-message" class="error-message">Invalid username or password</div>
        </form>
    </div>

    <div id="services-container" class="container">
        <h1>Available Services</h1>
        <div class="button-container">
            <a href="/Ai/Free/swagger" class="button">AI Service</a>
            <a href="/Blog/Free/swagger" class="button">Blog Service</a>
            <a href="/Drive/Free/swagger" class="button">Drive Service</a>
            <a href="/Order/Free/swagger" class="button">Order Service</a>
            <a href="/Payment/Free/swagger" class="button">Payment Service</a>
            <a href="/Product/Free/swagger" class="button">Product Service</a>
            <a href="/Search/Free/swagger" class="button">Search Service</a>
            <a href="/Site/Free/swagger" class="button">Site Service</a>
            <a href="/Support/Free/swagger" class="button">Support Service</a>
            <a href="/User/Free/swagger" class="button">User Service</a>
        </div>
    </div>

    <script>
        document.getElementById('login-form').addEventListener('submit', function(e) {
            e.preventDefault();
            
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            const errorMessage = document.getElementById('error-message');
            const loginContainer = document.getElementById('login-container');
            const servicesContainer = document.getElementById('services-container');

            if (username === 'lionheart' && password === 'admin') {
                loginContainer.style.display = 'none';
                servicesContainer.style.display = 'block';
            } else {
                errorMessage.style.display = 'block';
            }
        });
    </script>
</body>
</html>

]])
