<script>
export default {
    components: {},
    data: function () {
        return {
            errormsg: null,
            username: "",
            user: {
                id: 0,
                username: "",
            },
        }
    },
    methods: {
        async doLogin() {
            if (this.username == "") {
                this.errormsg = "Username cannot be empty.";
            } else {
                try {
                    let response = await this.$axios.post("/session", { username: this.username })
                    this.user = response.data
                    localStorage.setItem("requesterID", this.user.id);
                    localStorage.setItem("username", this.user.username);
                    this.$router.push({ path: '/session' })
                } catch (e) {
                    if (e.response && e.response.status === 400) {
                        this.errormsg = "Form error, please try again";
                        this.detailedmsg = null;
                    } else if (e.response && e.response.status === 500) {
                        this.errormsg = "An internal error occurred. Please try again later.";
                        this.detailedmsg = e.toString();
                    } else {
                        this.errormsg = e.toString();
                        this.detailedmsg = null;
                    }
                }
            }

        }
    },
    mounted() {

    }

}
</script>


<template>
    <div class="login-container">
        <h1 class="title">WASA-PHOTO</h1>
        <div class="input-group">
            <label for="username" class="form-label">Please enter your username to login:</label>
            <input type="text" id="username" v-model="username" placeholder="Enter your username">
            <button @click="doLogin">Login</button>
        </div>
        <p v-if="errormsg" class="error-msg">{{ errormsg }}</p>
    </div>
</template>

<style>
.login-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100vh;
    font-family: 'Arial', sans-serif;
}

.title {
    margin-bottom: 20px;
    color: #333;
    font-size: 6rem; /* Larger font size */
    font-family: 'Helvetica Neue', sans-serif; /* Different font */
    text-align: center; /* Ensure title is centered at the top */
}

.input-group {
    display: flex;
    flex-direction: column; /* Align label and input vertically */
    align-items: center;
    margin-bottom: 15px;
}

.input-group label {
    margin-bottom: 10px; /* Adjust spacing */
    font-size: 1.2rem; /* Larger label text */
    color: #555;
    text-align: center; /* Center the label text */
}

.input-group input[type="text"] {
    width: 300px;
    padding: 8px;
    font-size: 1rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    margin-bottom: 10px; /* Adjust spacing */
}

.input-group button {
    width: 100px;
    padding: 8px 12px;
    font-size: 1rem;
    background-color: #28a745;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
}

.input-group button:hover {
    background-color: #218838;
}

.error-msg {
    color: #d9534f;
    font-size: 0.9rem;
}
</style>