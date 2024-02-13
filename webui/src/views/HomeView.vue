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
			commentText: "",
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
				banned: 0,
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
			newUsername: "",
		
		}
	},
	methods: {
		async refresh() {
			this.myOwnProfile()
			this.getStream()
		},
		async selectFile() {
			this.images = this.$refs.file.files[0] // Get the first file
		},
		async Upload() {
			if (this.images === null) {
				this.errormsg = "Please select a file to upload."
			} else {
				// Prepare the FormData object to send the file
				let formData = new FormData();
    			formData.append('image', this.images);
				try {
					let response = await this.$axios.post("/users/" + this.username + "/photos/", formData, {
						headers: {
							Authorization: "Bearer " + localStorage.getItem("requesterID"),
							'Content-Type': 'multipart/form-data'
						}
					})
					this.profile = response.data
					this.errormsg = "Photo uploaded successfully."
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
		async setNewUsername() {
            if (this.newUsername == "") {
                this.errormsg = "Please insert a valid username."
            } else {
                try {
                    let response = await this.$axios.put("/users/" + this.username, this.newUsername , {
                        headers: {
                            Authorization: "Bearer " + localStorage.getItem("requesterID"),
							'Content-Type': 'application/json'
                        }
                    })
                    this.clear = response.data
                    localStorage.setItem("username", this.newUsername);
                    this.profile.username = this.newUsername;
                    this.username = this.newUsername;
                    this.refresh()
                } catch (e) {
                    if (e.response && e.response.status === 400) {
                        this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
                        this.detailedmsg = null;
                    } else if (e.response && e.response.status === 500) {
                        this.errormsg = "An internal error occurred. We will be notified. Please try again later.";
                        this.detailedmsg = e.toString();
                    } else {
                        this.errormsg = e.toString();
                        this.detailedmsg = null;
                    }
                }
            }

        },
		async getStream() {
			try {
				let response = await this.$axios.get("/users/" + this.username + "/homepage", {
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
					this.errormsg = e.toString() + e.response.status.toString();
				}
			}
		},
		async getPhotos() {
			try {
				let response = await this.$axios.get("/users/" + this.username + "/photos/", {
					headers: {
						Authorization: "Bearer " + localStorage.getItem("requesterID")
					}
				})
				this.Stream = response.data
				for (let i = 0; i < this.Stream.length; i++) {
					this.Stream[i].file = 'data:image/*;base64,' + this.Stream[i].file
				}
				this.$router.push({ path: '/photos' }) 
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
					this.$router.push({ path: '/users/' + this.usernameToSearch + '/profile' }) // Da cambiare
				} catch (e) {
					if (e.response && e.response.status === 400) {
						this.errormsg = "Form error, please try again.";
					} else if (e.response && e.response.status === 500) {
						this.errormsg = "No one found";
					}
					else if (e.response && e.response.status === 403) {
						this.errormsg = "Ops, cant's get the user profile, you've been banned by the user!";
					} else {
						this.errormsg = e.toString();
					}
				}
			}
		},
		async myOwnProfile() {
			try {
				let response = await this.$axios.get("users/" + this.username + "/profile", {
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
		async seeUserProfile(username) {
			try {
				let response = await this.$axios.get("users/" + username + "/profile", {
					headers: {
						Authorization: "Bearer " + localStorage.getItem("requesterID")
					}
				})
				this.profile = response.data
				this.$router.push({ path: '/users/' + username + '/profile' }) // Da cambiare
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
		async openSetUsernameLog() {
			try {
				var myModal = new bootstrap.Modal(document.getElementById('setUsernameLogModal'));
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
				let response = await this.$axios.get("/users/" + this.username + "/following/", {
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
		async openBannedLog() {
			try {
				let response = await this.$axios.get("/users/" + this.username + "/banned/", {
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
		async doLogout() {
			localStorage.removeItem("requesterID")
			localStorage.removeItem("username")
			this.$router.push({ path: '/' })
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
				<h1 class="h2 my-auto"><strong>Home page - Welcome back to WASA_PHOTO {{ this.username }}</strong></h1>
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
						<div>
							<!-- File Chooser and Upload Button -->
							<div class="btn-group mb-2">
								<input type="file" accept="image/*" class="btn" @change="selectFile" ref="file" style="background-color: #f8f9fa; border-color: #ced4da;">
								<button type="button" class="btn btn-sm" @click="Upload" style="background-color: #28a745; color: white;">
									Upload
								</button>
							</div>
						</div>
						<div class="profile-infobox">
							<h2 style="font-size:40px">My Profile</h2>
							<div class="info-row">
								<p style="font-size:30px"><strong>Username:</strong><button type="button" class="btn btn-outline-primary ms-2 d-flex align-items-center" @click="openSetUsernameLog()" :style="{ borderColor: 'white' }"><strong style="font-size: 30px;color: black">{{ this.profile.username }}</strong></button></p>
							</div>
							<div class="info-row">
								<p style="font-size:30px"><strong>Followers:</strong><span><strong>{{ this.profile.followers }}</strong></span></p>
							</div>
							<div class="info-row">
								<p style="font-size:30px"><strong>Following:</strong><button type="button" class="btn btn-outline-primary ms-2 d-flex align-items-center" @click="openFollowedLog()" :style="{ borderColor: 'white' }"><strong style="font-size: 30px;color: black">{{ this.profile.followed }}</strong></button></p>
							</div>
							<div class="info-row">
								<p style="font-size:30px"><strong>Banned:</strong><button type="button" class="btn btn-outline-primary ms-2 d-flex align-items-center" @click="openBannedLog()" :style="{ borderColor: 'white' }"><strong style="font-size: 30px;color: black">{{ this.profile.banned }}</strong></button></p>
							</div>
							<div class="info-row">
								<p style="font-size:30px"><strong>Photos:</strong><button type="button" class="btn btn-outline-primary ms-2 d-flex align-items-center" @click="getPhotos()" :style="{ borderColor: 'white' }"><strong style="font-size: 30px;color: black">{{ this.profile.photos }}</strong></button></p>
							</div>
						</div>
						<div class="modal fade" id="setUsernameLogModal" tabindex="-1">
							<div class="modal-dialog">
								<div class="modal-content">
									<div class="modal-header">
									<h5 class="modal-title">Set new username</h5>
									<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close">
										<span aria-hidden="true"></span>
									</button>
									</div>
									<div class="modal-body">
										<div class="btn-group mb-2">
											<input type="text" placeholder="Set a new username." id="newUsername" v-model="newUsername" aria-describedby="button-addon" style="width:300px;height:40px;">
											<button class="btn btn-primary" type="button" id="button-addon" @click="setNewUsername()" data-bs-dismiss="modal">Confirm</button>
										</div>									
									</div>
									<div class="modal-footer">
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
					</div>
		
					<!-- Photo Content Column -->
					<div class="col-md-8">
						<div class="row">
							<h1 class="h2 my-auto"><strong style="font-size:40px">Here's the stream of photos </strong><span><strong>{{ this.profile.username }}</strong></span></h1>
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
											{{ comment.text }} --- by <strong style="margin-right:5px">{{ comment.username }}</strong>
											<button v-if="comment.username==this.username" class="btn-primary" type="button" id="button-addon" @click="uncommentPhoto(this.photoUsername, this.photoId, comment.id)" data-bs-dismiss="modal" style="background-color: #dc3545; color: white;">Delete</button>
											</li>
										</ul>										
										<p v-if="Comments==null"><strong>This photo has no comment yet!!</strong></p>										
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
						<div class="col-md-4 custom-margin" v-for="photo in Stream" :key="photo.id">
							<!-- Bootstrap card -->
							<div class="card">
								<!-- Image at the top of the card -->
								<img class="card-img-top" :src="photo.file" alt="Photo" style="width: 100%; height: 250px; object-fit: contain;background-color: black;">
							
								<!-- Card body for text content -->
								<div class="card-body">
								<h5 class="card-title d-flex align-items-center">
									Uploaded by :  <button type="button" class="btn btn-outline-primary ms-2 d-flex align-items-center" @click="seeUserProfile(photo.username)" :style="{ borderColor: 'white' }"><strong style="font-size: 20px;color: black">{{ photo.username }}</strong></button>
								</h5>
								<p class="card-text"><strong>Uploaded on : </strong> {{ new Date(photo.uploadtime).toLocaleString() }}</p>
								<p class="card-text d-flex align-items-center"><strong>Likes : </strong><button type="button" class="btn btn-outline-primary ms-2 d-flex align-items-center" @click="openLikesLog(photo.username, photo.id)" :style="{ borderColor: 'white' }"><strong style="color: black">{{ photo.likesN }}</strong></button></p>
								<p class="card-text"><strong>Comments  : </strong><strong style="font-size: 17px;">{{ photo.commentsN }}</strong></p>
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
/* Apply consistent height and styling to all buttons */
.btn {
    height: 40px; /* Ensure consistent height */
    padding: 8px 12px; /* Adjust padding as needed */
    font-size: 1rem; /* Ensure consistent font size */
    border-radius: 4px; /* Optional: Ensure consistent border radius */
}

.btn-custom {
    width: 100px;  
}



.profile-infobox {
    padding: 8px;
	width: 400px;
    border-radius: 4px;
    border: 1px solid #ddd; 
    background-color: #f8f9fa; /* Light gray background */
}

.profile-infobox .info-row p {
    display: flex;
    justify-content: space-between;
    margin-bottom: 0.5rem;
}

.profile-infobox .info-row p strong {
    margin-right: 0.5rem;
}

.profile-infobox .info-row p span {
    margin-left: auto;
}

.custom-margin {
	margin-bottom: 30px; /* or any other value */
}

.heart-button {
	background: transparent; /* Ensures no background color */
	border: none;
	font-size: 40px; /* Adjust size as needed */
	cursor: pointer;
	color: #ccc; /* Default color for unliked state */
  }
  
  /* Color change for liked state */
  .heart-button[liked=true], .heart-button.liked {
	color: red; /* Keeps the heart red when liked */
  }
  
  /* Optional: Different color on hover for visual feedback */
  .heart-button:hover {
	color: #ff6666; /* Lighter red on hover, adjust as needed */
  }
  
  /* Ensure the heart remains red when liked, even on hover */
  .heart-button.liked:hover {
	color: red; /* Keeps the heart red when liked, even on hover */
  }
  
</style>

