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
                    let response = await this.$axios.post("/session", {username: this.username})
                    this.user = response.data
                    localStorage.setItem("requesterID", this.user.id);
                    localStorage.setItem("username", this.user.username);
                    this.$router.replace('/homepage')
                } catch (e) {
                    if (e.response && e.response.status === 400) {
                        this.errormsg = "Form error, please try again";
                    } else if (e.response && e.response.status === 500) {
                        this.errormsg = e.toString();
                    } else {
                        this.errormsg = e.toString();
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
    <div class="logInBox">
        <div class="info">
          <h4>WELCOME!</h4>
        </div>
        <div class="inputArea">
          <div class="wordArea">
            <input class="userName" id="username" v-model="username" placeholder="Enter your username" style="font-size:25px;">
          </div>
          <div class="buttonArea">
            <input type="button" class="logButton" value="Log In" @click="doLogin">
          </div>
        </div>
        <p v-if="errormsg" class="error-msg">{{ errormsg }}</p>
    </div>

</template>

<style>
.logInBox{
    width:500px;
    height:200px;
    background:#abc;
    position:fixed;
    top:50%;
    left:50%;
    margin-top:-150px;
    margin-left:-250px;
    opacity:0.8;
}

.info{
width:100%;
height:30%;
background:#123;
text-align:center;
}

h4{
color:white;
margin:auto;
position:relative;
top:5px;
font-size:45px;
}

.inputArea{
width:60%;
height:70%;
margin:auto;
}

.wordArea{
width:80%;
margin:auto;
}

.userName{
width:100%;
height:40px;
margin-top:15px;
}

.buttonArea{
position:relative;
top:10%;
width:40%;
margin:auto;
}

.logButton{
width:100%;
height:40px;
border-radius:10px;
outline:none;
border:1px solid #abc;
}

.click:hover{
background:#bbb;
border:#bbb;
}

.click:active{
margin-left:1px;
margin-top:1px;
}

.error-msg {
    color: #d9534f;
    font-size: 0.9rem;
}
</style>