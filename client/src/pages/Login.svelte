<script>
    import Input from '../components/Input.svelte';
    import user_logo from '../vectors/man-user.svg';
    import { push } from 'svelte-spa-router';
    import { server_name } from '../server_data';

    const username_input_id = "login-username-field";
    const password_input_id = "password-username-field";
    
    const authenticateUser = (username, password) => {
        fetch(`${server_name}/login?username=${username}&password=${password}`)
            .then(promise => {
                if(promise.status == 200) {
                    // user is valid
                    promise.json().then(response => {
                        push(`/home/${response.response}`);
                    });
                } else if(promise.status == 400) {
                    // user doesnt exists
                    alert("user doesnt exists");
                } else if(promise.status == 203) {
                    // wrong password
                    alert("wrong password");
                }
            })
    }

    const inputRecived = e => {
        if(e.key.toLowerCase() === 'enter') {
            isLoginReady();
        }
    };

    const isLoginReady = () => {
        const username_element = document.getElementById(username_input_id);
        const password_element = document.getElementById(password_input_id);
        if(username_element.value !== "" && password_element.value !== "") {
            return authenticateUser(username_element.value, password_element.value);
        }
    };


</script>

<style>
    #login-page-content{
        display: flex;
        width: 40%;
        height: 100vh;
        flex-direction: column;
        margin: 0 auto;
        justify-content: center;
        align-items: center;
    }

    #login-logo-container {
        display: flex;
        background-color: var(--theme-color);
        width: 20vh;
        height: 20vh;
        justify-content: center;
        align-items: center;
        border-radius: 50%;
        box-shadow: inset 0 0 15px 5px #000000ae;
    }

    .login-field {
        margin-top: 5vh;
    }

    :global(#login-logo-container svg){
        fill: var(--dark-color);
        /* clip-path: circle(50%); */
        width: 80%;
    }
</style>

<div id="login-page-content">
    <div class="login-field">
        <!-- Login header -->
        <div id="login-logo-container">
            {@html user_logo}
        </div>
    </div>
    <div class="login-field">
        <Input 
            input_id={username_input_id} 
            input_placeholder="username"
            onKeypressed={inputRecived}
        />
    </div>
    <div class="login-field">
        <Input 
            input_id={password_input_id} 
            input_placeholder="password"
            input_type="password"
            onKeypressed={inputRecived}
        />
    </div>
</div>