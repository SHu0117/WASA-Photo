<script>
export default {
	data: function() {
		return {
			errormsg: null,
			username: localStorage.getItem('username'),
			requesterID: localStorage.getItem('requesterID'),
			images: null,
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
            photoUsername: "",
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
            users: {
				user: [
					{
						id: 0,
						username: "",
					}
				],
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
				this.Stream = response.data ?? []
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
					this.$router.push({ path: '/users/' + this.searchUserUsername + '/profile' }) 
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
        async backToHomepage() {
			this.$router.push({ path: '/homepage' }) 		
		},
		async seeUserProfile(username) {
			try {
				let response1 = await this.$axios.get("users/" + username + "/profile", {
					headers: {
						Authorization: "Bearer " + localStorage.getItem("requesterID")
					}
				})
				this.profile = response1.data
				let response2 = await this.$axios.get("users/" + username + "/photos/", {
					headers: {
						Authorization: "Bearer " + localStorage.getItem("requesterID")
					}
				})
				this.Stream = response2.data ?? []
				for (let i = 0; i < this.Stream.length; i++) {
					this.Stream[i].file = 'data:image/*;base64,' + this.Stream[i].file
				}
				this.$router.push({ path: '/users/' + username + '/profile' })
			} catch (e) {
				if (e.response && e.response.status === 400) {
					this.errormsg = "Form error, please try again.";
					this.detailedmsg = null;
				} else if (e.response && e.response.status === 500) {
					this.errormsg = "No one found";
				} else if (e.response && e.response.status === 403) {
						this.errormsg = "Ops, cant's get the user profile, you've been banned by the user!";
				}else {
					this.errormsg = e.toString();
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
        async uncommentPhoto(username, Photo_id, Comment_id) {
			
			try {
				let response = await this.$axios.delete("/users/" + username + "/photos/" + Photo_id + "/comments/" + Comment_id, {
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
		async openFollowedLog() {
			try {
				let response = await this.$axios.get("/users/" + this.$route.params.username + "/following/", {
					headers: {
						Authorization: "Bearer " + localStorage.getItem("requesterID")
					}
				})
				this.users = response.data;
				var myModal = new bootstrap.Modal(document.getElementById('usersLogModal'));
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
		async openFollowerLog() {
			try {
				let response = await this.$axios.get("/users/" + this.$route.params.username + "/follower/", {
					headers: {
						Authorization: "Bearer " + localStorage.getItem("requesterID")
					}
				})
				this.users = response.data;
				var myModal = new bootstrap.Modal(document.getElementById('usersLogModal'));
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
        async openLikesLog(username, Photo_id) {
			try {
				let response = await this.$axios.get("/users/" + username + "/photos/" + Photo_id + "/likes", {
					headers: {
						Authorization: "Bearer " + localStorage.getItem("requesterID")
					}
				})
				this.users = response.data;
				var myModal = new bootstrap.Modal(document.getElementById('usersLogModal'));
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
				<h1 class="h2 my-auto"><strong>WASA-PHOTO - Search Result</strong></h1>
				<div class="d-flex align-items-center">
					<!-- Grouped Buttons and Input for Alignment -->
					<i class="fas fa-sync icon-spacing" @click="refresh()" aria-label="Refresh"></i>

					<i class="fas fa-sign-out-alt icon-spacing" @click="doLogout()" aria-label="Logout"></i>
					
					<i class="fas fa-home btn-home-icon icon-spacing" @click="backToHomepage()" aria-label="Go to homepage"></i>
				</div>
			</div>
		</div>

		<div>
			<div class="container-fluid">
				<div class="row">
					<!-- Profile Information Column -->
					<div class="col-md-4">
						<div class="container py-5">
							<div class="row align-items-center">
								<div class="card" style="border-radius: 15px; width:400px;height:200px;">
									<div class="card-body">
										<div class="d-flex text-black">										
											<div class="flex-grow-1 ms-3">
												<h5 class="mb-1" style="font-size:25px; font-weight:bold;">
													Visiting - <strong style="text-decoration:underline">{{ this.profile.username}} </strong>													
												</h5>
												<div class="d-flex justify-content-start rounded-3 p-2 mb-2" style="background-color: #efefef;">
													<div class="px-3">
														<p class="text-muted mb-1" style="font-weight:bold;">Post</p>
														<button class="number-button" @click="getPhotos()">{{ this.profile.photos }}</button>
													</div>
													<div class="px-3">
														<p class="text-muted mb-1" style="font-weight:bold;">Followers</p>
														<button class="number-button" @click="openFollowerLog()">{{ this.profile.followers }}</button>
													</div>
													<div class="px-3">
														<p class="text-muted mb-1" style="font-weight:bold;">Following</p>
														<button class="number-button" @click="openFollowedLog()">{{ this.profile.followed }}</button>
													</div>
												</div>
												<div class="d-flex pt-1">
													<button type="button" v-if="profile.isBanned==true" class="btn btn-custom btn-outline-primary me-1" @click="unbanUser()">Unban</button>
													<button type="button" v-if="profile.isBanned==false" class="btn btn-custom btn-outline-primary me-1" @click="banUser()">Ban</button>
													<button type="button" v-if="profile.isFollowed==false" class="btn btn-custom btn-primary" @click="followUser()">Follow</button>
													<button type="button" v-if="profile.isFollowed==true" class="btn btn-custom btn-primary" @click="unfollowUser()">Unfollow</button>
												</div>
											</div>
										</div>
									</div>
								</div>
							</div>
						</div>
					</div>

					<div class="modal fade" id="usersLogModal" tabindex="-1">
						<div class="modal-dialog">
							<div class="modal-content">
								<div class="modal-header">
								<h5 class="modal-title">Here's the list of users </h5>
								<button type="button" class="btn-close" data-bs-dismiss="modal">
									<span aria-hidden="true"></span>
								</button>
								</div>
								<div class="modal-body">
									<ul class="list-group custom-margin">
										<li class="list-group-item" v-for="user in this.users" :key="user.id">
											<strong><button type="button" class="btn btn-outline-primary ms-2 d-flex align-items-center" data-bs-dismiss="modal" @click="seeUserProfile(user.username)" :style="{ borderColor: 'white' }"><strong style="color: black">{{ user.username }}</strong></button></strong>
										</li>
									</ul>
									<p v-if="this.users==null" class="align-items-center justify-content-center"><strong style="font-size:20px">No users found</strong></p>
								</div>
								<div class="modal-footer">
									<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
								</div>
							</div>
						</div>
					</div>
		
					<!-- Photo Content Column -->
					<div class="col-md-8">
						<div class="row custom-margin">
							<h1 class="h2 my-auto" style="font-weight:bold;text-decoration:underline;"><strong>Photo uploaded by </strong><span><strong>{{ this.profile.username }}</strong></span></h1>
							<!-- Comments Log Modal -->
							<div class="modal fade" id="commentsLogModal" tabindex="-1">
								<div class="modal-dialog">
                                    <div class="modal-content">
                                        <div class="modal-header">
                                        <h5 class="modal-title">Comments Log</h5>
                                        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close">
                                            <span aria-hidden="true"></span>
                                        </button>
                                        </div>
                                        <div class="modal-body">
                                            <!-- Dynamic comments will be loaded here -->
                                            <ul class="list-group custom-margin">
                                                <li class="list-group-item" v-for="comment in Comments" :key="comment.id">
                                                    {{ comment.text }} --- by <strong style="margin-right:5px;">{{ comment.username }}</strong>
                                                    <i v-if="comment.username==this.username" class="fas fa-trash-alt" @click="uncommentPhoto(this.photoUsername, this.photoId, comment.id)" data-bs-dismiss="modal"></i>
                                                </li>
                                            </ul>
                                            <p v-if="Comments==null">No comments to display.</p>
                                        </div>
                                        <div class="modal-footer">
                                            <div class="btn-group mb-2">
                                                <input type="text" placeholder="Write a comment." id="commentText" v-model="commentText" aria-describedby="button-addon" style="width:300px;height:40px;">
                                                <button class="btn btn-primary" type="button" id="button-addon" @click="commentPhoto(this.photoUsername, this.photoId)" data-bs-dismiss="modal">Send</button>
                                            </div>
                                            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                                        </div>
								    </div>
							    </div>
                            </div>

                            <div class="modal fade" id="usersLogModal" tabindex="-1">
                                <div class="modal-dialog">
                                    <div class="modal-content">
                                        <div class="modal-header">
                                        <h5 class="modal-title">Here's the list of users </h5>
                                        <button type="button" class="btn-close" data-bs-dismiss="modal">
                                            <span aria-hidden="true"></span>
                                        </button>
                                        </div>
                                        <div class="modal-body">
                                            <ul class="list-group custom-margin">
                                                <li class="list-group-item" v-for="user in this.users" :key="user.id">
                                                    <strong><button type="button" class="btn btn-outline-primary ms-2 d-flex align-items-center" data-bs-dismiss="modal" @click="seeUserProfile(user.username)" :style="{ borderColor: 'white' }"><strong style="color: black">{{ user.username }}</strong></button></strong>
                                                </li>
                                            </ul>
                                            <p v-if="this.users==null" class="align-items-center justify-content-center"><strong style="font-size:20px">No users found</strong></p>
                                        </div>
                                        <div class="modal-footer">
                                            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            
                            
                            <div class="col-md-4 custom-margin" v-for="photo in this.Stream" :key="photo.id">
								<!-- Bootstrap card -->
								<div class="card">
								  <!-- Image at the top of the card -->
								  <img class="card-img-top" :src="photo.file" alt="Photo" style="width: 100%; height: 250px; object-fit: contain;background-color: black">
							  
								  <!-- Card body for text content -->
								  <div class="card-body">
									<h5 class="card-title d-flex align-items-center">
                                        Uploaded by :  <button type="button" class="btn btn-outline-primary ms-2 d-flex align-items-center" @click="seeUserProfile(photo.username)" :style="{ borderColor: 'white' }"><strong style="font-size: 20px;color: black">{{ photo.username }}</strong></button>
                                    </h5>
									<p class="card-text"><strong>Uploaded on : </strong> {{ new Date(photo.uploadtime).toLocaleString() }}</p>
									<p class="card-text d-flex align-items-center"><strong>Likes : </strong><button class="number-button" @click="openLikesLog(photo.username, photo.id)" :style="{ borderColor: 'white' }"><strong style="color: black">{{ photo.likesN }}</strong></button></p>
									<p class="card-text"><strong>Comments : </strong><strong style="font-size: 16px;">{{ photo.commentsN }}</strong></p>
								  </div>
							  
								  <!-- Card footer for buttons -->
								  <div class="card-footer">
									<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center">
										<button type="button" class="heart-button" :class="{ 'liked': photo.isliked }" @click="photo.isliked ? unlikePhoto(photo.username, photo.id) : likePhoto(photo.username, photo.id)">
                                            {{ photo.isliked ? '♥' : '♡' }}
                                        </button>
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

.btn {
    height: 40px;
    padding: 8px 12px;
    font-size: 1rem; 
    border-radius: 4px;
}

.btn-custom {
    width: 100px;  
}

.custom-margin {
	margin-bottom: 30px; 
  }

.heart-button {
background: transparent; 
border: none;
font-size: 40px;
cursor: pointer;
color: #ccc; 
}

/* Color change for liked state */
.heart-button[liked=true], .heart-button.liked {
color: red; 
}

.heart-button:hover {
	color: #ff6666; /* Lighter red on hover, adjust as needed */
}

.delete-button {
    background-color: transparent;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 5px;
}

.btn-home-icon {
	cursor: pointer; 
	font-size: 30px; 
	color: black; 
	transition: color 0.3s ease;
}

.btn-home-icon:hover {
color: #555; /* Darken the icon on hover */
}

.fas.fa-sign-out-alt {
	cursor: pointer;
	font-size: 24px; 
	color: black; 
	transition: color 0.3s;
}

.fas.fa-sign-out-alt:hover {
color: #555; /* Darken the icon on hover for feedback */
}

  /* Refresh icon button styling */
.fas.fa-sync {
	cursor: pointer;
	font-size: 24px;
	color: black; 
	transition: color 0.3s;
}

.fas.fa-sync:hover {
	color: #555; /* Darken the icon on hover for feedback */
}
.icon-spacing {
  margin: 0 10px; /* Adjust the margin as needed */
}

  
</style>

