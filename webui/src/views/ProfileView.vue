<script>
export default {
	data: function() {
		return {
			errormsg: null,
			successmsg: null,
			username: localStorage.getItem('username'),
			requesterID: localStorage.getItem('requesterID'),
			loading: false,
			some_data: null,
			images: null,
			image: null,
			clear: null,
			Comments: {
				comment: [
					{
						id: 0,
						user_id: 0,
						username: "",
						photo_id: 0,
						photo_Owner: 0,
						text: "",
					}
				],
			},
			comment: "",
			Stream: {
				photo: [
					{
						id: 0,
						useriId: 0,
						username: "",
						likesN: 0,
						commentsN: 0,
						uploadtime: "",
						file: "",
						isliked: null,
					}
				],
			},
			usernameToSearch: "",
            photoUsernmae: "",
            photoId: 0,
			like: {
				ID: 0,
				User_id: 0,
				Photo_id: 0,
				Photo_user: 0,
			},
			profile: {
				requester_id: 0,
				user_id: 0,
				username: "",
				followers: 0,
				followed: 0,
				photos: 0,
				IsFollowed: null,
				IsBanned: null,
			},
		
		}
	},
	methods: {
		async refresh() {
			this.userProfile()
			this.getPhotos()
		},
        async userProfile() {
			try {
				let response = await this.$axios.get("users/" + this.$route.params.username + "/profile", {
					headers: {
						Authorization: "Bearer " + localStorage.getItem("requesterID")
					}
				})
				this.profile = response.data
			} catch (e) {
				if (e.response && e.response.status === 400) {
					this.errormsg = "Form error, please try again.";
					this.detailedmsg = null;
				} else if (e.response && e.response.status === 500) {
					this.errormsg = "No one found";
				} else {
					this.errormsg = e.toString();
				}
			}
		},
		async getPhotos() {
			try {
				let response = await this.$axios.get("/users/" + this.$route.params.username + "/photos/", {
					headers: {
						Authorization: "Bearer " + localStorage.getItem("requesterID")
					}
				})
				this.Stream = response.data
				for (let i = 0; i < this.Stream.length; i++) {
					this.Stream[i].file = 'data:image/*;base64,' + this.Stream[i].file
				}
			} catch (e) {
				if (e.response && e.response.status === 400) {
					this.errormsg = "Form error, please try again";
				} else if (e.response && e.response.status === 500) {
					this.errormsg = "An internal error occurred. Please try again later.";
				} else {
					this.errormsg = e.toString();
				}
			}
		},
		async SearchUser() {
			if (this.usernameToSearch === this.username) {
				this.errormsg = "You can't search yourself."
			} else if (this.searchUserUsername === "") {
				this.errormsg = "Please insert a valid username."
			} else {
				try {
					let response = await this.$axios.get("users/" + this.usernameToSearch + "/profile", {
						headers: {
							Authorization: "Bearer " + localStorage.getItem("requesterID")
						}
					})
					this.profile = response.data
					this.$router.push({ path: '/users/' + this.searchUserUsername + '/profile' }) // Da cambiare
				} catch (e) {
					if (e.response && e.response.status === 400) {
						this.errormsg = "Form error, please try again.";
						this.detailedmsg = null;
					} else if (e.response && e.response.status === 500) {
						this.errormsg = "No one found";
					} else {
						this.errormsg = e.toString();
					}
				}
			}
		},
		async commentPhoto(username, Photo_id) {
			if (this.commentText === "") {
				this.errormsg = "Emtpy comment is not valid."
			} else {
				try {
					let response = await this.$axios.post("/users/" + username + "/photos/" + Photo_id + "/comments/", { text: this.commentText }, {
						headers: {
							Authorization: "Bearer " + localStorage.getItem("requesterID")
						}
					})
					this.clear = response.data
					this.refresh()
				} catch (e) {
					if (e.response && e.response.status === 400) {
						this.errormsg = "Form error, please try again.";
						this.detailedmsg = null;
					} else if (e.response && e.response.status === 500) {
						this.errormsg = "An internal error occurred. Please try again later.";
					} else {
						this.errormsg = e.toString();
					}
				}
			}
		},
		async openCommentsLog(username, Photo_id) {
			try {
				let response = await this.$axios.get("/users/" + username + "/photos/" + Photo_id + "/comments/", {
					headers: {
						Authorization: "Bearer " + localStorage.getItem("requesterID")
					}
				})
				this.Comments = response.data;
				this.photoUsername = username;
				this.photoId = Photo_id;
				var myModal = new bootstrap.Modal(document.getElementById('commentsLogModal'));
    			myModal.show();
			} catch (e) {
				if (e.response && e.response.status === 400) {
					this.errormsg = "Form error, please try again.";
				} else if (e.response && e.response.status === 500) {
					this.errormsg = "An internal error occurred. Please try again later.";
				} else {
					this.errormsg = e.toString();
				}
			}
		},
		async likePhoto(username, Photo_id) {
			try {
				let response = await this.$axios.put("/users/" + username + "/photos/" + Photo_id + "/likes/" + this.username, {}, {
					headers: {
						Authorization: "Bearer " + localStorage.getItem("requesterID")
					}
				})
				this.clear = response.data
				this.refresh()
			} catch (e) {
				if (e.response && e.response.status === 400) {
					this.errormsg = "Form error, please try again.";
					this.detailedmsg = null;
				} else if (e.response && e.response.status === 500) {
					this.errormsg = "An internal error occurred. Please try again later.";
				} else {
					this.errormsg = e.toString();
				}
			}
		},
		async unlikePhoto(username, Photo_id) {

			try {
				let response = await this.$axios.delete("/users/" + username + "/photos/" + Photo_id + "/likes/" + this.username, {
					headers: {
						Authorization: "Bearer " + localStorage.getItem("requesterID")
					}
				})
				this.clear = response.data
				this.refresh()
			} catch (e) {
				if (e.response && e.response.status === 400) {
					this.errormsg = "Form error, please try again.";
					this.detailedmsg = null;
				} else if (e.response && e.response.status === 500) {
					this.errormsg = "An internal error occurred. Please try again later.";
				} else {
					this.errormsg = e.toString();
				}
			}
		},
        async followUser() {
			try {
				let response = await this.$axios.put("/users/" + this.username + "/following/" + this.$route.params.username, {}, {
					headers: {
						Authorization: "Bearer " + localStorage.getItem("requesterID")
					}
				})
				this.clear = response.data
				this.refresh()
			} catch (e) {
				if (e.response && e.response.status === 400) {
					this.errormsg = "Form error, please try again.";
					this.detailedmsg = null;
				} else if (e.response && e.response.status === 500) {
					this.errormsg = "An internal error occurred. Please try again later.";
				} else {
					this.errormsg = e.toString();
				}
			}
		},
        async unfollowUser() {
			try {
				let response = await this.$axios.delete("/users/" + this.username + "/following/" + this.$route.params.username, {
					headers: {
						Authorization: "Bearer " + localStorage.getItem("requesterID")
					}
				})
				this.clear = response.data
				this.refresh()
			} catch (e) {
				if (e.response && e.response.status === 400) {
					this.errormsg = "Form error, please try again.";
					this.detailedmsg = null;
				} else if (e.response && e.response.status === 500) {
					this.errormsg = "An internal error occurred. Please try again later.";
				} else {
					this.errormsg = e.toString();
				}
			}
		},
        async banUser() {
			try {
				let response = await this.$axios.put("/users/" + this.username + "/banned/" + this.$route.params.username, {}, {
					headers: {
						Authorization: "Bearer " + localStorage.getItem("requesterID")
					}
				})
				this.clear = response.data
				this.refresh()
			} catch (e) {
				if (e.response && e.response.status === 400) {
					this.errormsg = "Form error, please try again.";
					this.detailedmsg = null;
				} else if (e.response && e.response.status === 500) {
					this.errormsg = "An internal error occurred. Please try again later.";
				} else {
					this.errormsg = e.toString();
				}
			}
		},
        async unbanUser() {
			try {
				let response = await this.$axios.delete("/users/" + this.username + "/banned/" + this.$route.params.username, {
					headers: {
						Authorization: "Bearer " + localStorage.getItem("requesterID")
					}
				})
				this.clear = response.data
				this.refresh()
			} catch (e) {
				if (e.response && e.response.status === 400) {
					this.errormsg = "Form error, please try again.";
					this.detailedmsg = null;
				} else if (e.response && e.response.status === 500) {
					this.errormsg = "An internal error occurred. Please try again later.";
				} else {
					this.errormsg = e.toString();
				}
			}
		},
		async doLogout() {
			localStorage.removeItem("requesterID")
			localStorage.removeItem("username")
			this.$router.push({ path: '/' })
		},
		async ViewProfile() {
			this.$router.push({ path: '/users/' + this.username + '/profile' })
		},
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		<div>
			<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
				<h1 class="h2 my-auto">Home page - Welcome back to WASA_PHOTO {{ this.username }}</h1>
				<div class="d-flex align-items-center">
					<!-- Grouped Buttons and Input for Alignment -->
					<button type="button" class="btn btn-primary mx-1" @click="refresh" style="background-color: #28a745; color: white;">
						Refresh
					</button>
					<button type="button" class="btn btn-primary mx-1" @click="doLogout" style="background-color: #dc3545; color: white;">
						Logout
					</button>
					<input type="text" id="usernameToSearch" v-model="usernameToSearch" class="form-control mx-1" placeholder="Search a user in WASA-PHOTO." style="height:40px;">
					<button class="btn btn-primary mx-1" type="button" @click="SearchUser">Search</button>
				</div>
			</div>
		</div>
		
		<div>
			<div class="container-fluid">
				<div class="row">
					<!-- Profile Information Column -->
					<div class="col-md-4">
						<!-- Existing Content Column -->
						<div class="profile-info mt-3">
							<h2>{{ this.profile.username }}'s Profile</h2>
							<div class="info-row">
								<p><strong>Username:</strong><span><strong>{{ this.profile.username }}</strong></span></p>
							</div>
							<div class="info-row">
								<p><strong>Followers:</strong><span><strong>{{ this.profile.followers }}</strong></span></p>
							</div>
							<div class="info-row">
								<p><strong>Following:</strong><span><strong>{{ this.profile.followed }}</strong></span></p>
							</div>
							<div class="info-row">
								<p><strong>Photos:</strong><span><strong>{{ this.profile.photos }}</strong></span></p>
							</div>
                            <div class="info-row">
                                <button type="button" v-if="profile.isBanned==true" class="btn btn-primary btn-custom" @click="unbanUser()" style="background-color: #28a745; color: white;">UnBAN</button>
							    <button type="button" v-if="profile.isBanned==false" class="btn btn-primary btn-custom" @click="banUser()" style="background-color: #dc3545; color: white;">BAN</button>									
                            </div>
                            <div class="info-row">
                                <button type="button" v-if="profile.isFollowed==true" class="btn btn-primary btn-custom" @click="unfollowUser(this.username)" style="background-color: #dc3545; color: white;">UnFollow</button>
							    <button type="button" v-if="profile.isFollowed==false" class="btn btn-primary btn-custom" @click="followUser()" style="background-color: #28a745; color: white;">Follow</button>
                            </div>
                        </div>
					</div>
		
					<!-- Photo Content Column -->
					<div class="col-md-8">
						<div class="row">
							<h1 class="h2 my-auto"><strong>Photo uploaded by </strong><span><strong>{{ this.profile.username }}</strong></span></h1>
							<!-- Comments Log Modal -->
							<div class="modal fade" id="commentsLogModal" tabindex="-1">
								<div class="modal-dialog">
                                    <div class="modal-content">
                                        <div class="modal-header">
                                        <h5 class="modal-title">Comments Log</h5>
                                        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close">
                                            <span aria-hidden="true">&times;</span>
                                        </button>
                                        </div>
                                        <div class="modal-body">
                                            <!-- Dynamic comments will be loaded here -->
                                            <ul class="list-group custom-margin">
                                                <li class="list-group-item" v-for="comment in Comments" :key="comment.id">
                                                {{ comment.text }} --- by <strong>{{ comment.username }}</strong>
                                                </li>
                                            </ul>
                                            <p v-if="Comments==null">No comments to display.</p>
                                            <div class="btn-group mb-2">
                                                <input type="text" placeholder="Write a comment." id="commentText" v-model="commentText" aria-describedby="button-addon" style="width:300px;height:40px;">
                                                <button class="btn btn-primary" type="button" id="button-addon" @click="commentPhoto(this.photoUsername, this.photoId)" data-bs-dismiss="modal">Send</button>
                                            </div>
                                        </div>
                                        <div class="modal-footer">
                                            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                                        </div>
								    </div>
							    </div>
                            </div>
                            
                            
                            <div class="col-md-4 custom-margin" v-for="photo in Stream" :key="photo.id">
								<!-- Bootstrap card -->
								<div class="card">
								  <!-- Image at the top of the card -->
								  <img class="card-img-top" :src="photo.file" alt="Photo" style="width: 100%; height: auto;">
							  
								  <!-- Card body for text content -->
								  <div class="card-body">
									<h5 class="card-title">Uploaded by :  {{ photo.username }}</h5>
									<p class="card-text"><strong>Uploaded on : </strong> {{ new Date(photo.uploadtime).toLocaleString() }}</p>
									<p class="card-text"><strong>Likes : </strong>{{ photo.likesN }}</p>
									<p class="card-text"><strong>Comments : </strong>{{ photo.commentsN }}</p>
								  </div>
							  
								  <!-- Card footer for buttons -->
								  <div class="card-footer">
									<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center">
										<button type="button" v-if="photo.isliked==false" class="btn btn-primary btn-custom" @click="likePhoto(photo.username, photo.id)">Like</button>
										<button type="button" v-if="photo.isliked==true" class="btn btn-primary btn-custom" @click="unlikePhoto(photo.username, photo.id)" style="background-color: #dc3545; color: white;">Unike</button>
										<button type="button" class="btn btn-secondary btn-custom" @click="openCommentsLog(photo.username, photo.id)">Comments</button>										
									</div>
								  </div>
								</div>
								<!-- End of Bootstrap card -->
							  </div>
							  
						</div>
					</div>
				</div>
			</div>	
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		</div>
	</div>
</template>




<style>
/* Apply consistent height and styling to all buttons */
.btn {
    height: 40px; /* Ensure consistent height */
    padding: 8px 12px; /* Adjust padding as needed */
    font-size: 1rem; /* Ensure consistent font size */
    border-radius: 4px; /* Optional: Ensure consistent border radius */
}

.btn-custom {
    width: 100px;  /* or any other value */
}

/* Specific styles for input to ensure it matches button height */
.input-group input[type="text"] {
    height: 40px; /* Match button height */
    width: 400px; /* Specific width for the input */
    padding: 8px;
    font-size: 1rem;
    border: 1px solid #ddd;
    border-radius: 4px;
}

.profile-info {
    /* Optional: Add some styling to your profile info for better presentation */
    padding: 8px;
	width: 200px;
    border-radius: 4px;
    border: 1px solid #ddd; /* Just an example, adjust as needed */
    background-color: #f8f9fa; /* Light gray background */
}

.profile-info .info-row p {
    display: flex;
    justify-content: space-between;
    margin-bottom: 0.5rem;
}

.profile-info .info-row p strong {
    margin-right: 0.5rem;
}

.profile-info .info-row p span {
    margin-left: auto;
}

.custom-margin {
	margin-bottom: 30px; /* or any other value */
  }
</style>


  