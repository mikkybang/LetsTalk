var vue = new Vue({
	el: "#login",
	data: {
		alert: {
			message: "hello world"
		},
		login: {
			login: "",
			password: ""
		},
		shake: false,
		good: "",
		fake: {
			login: "vincent",
			password: "admin"
		}
	},
	computed: {
		isShake: function(){
			console.log(this.shake);
			if(this.shake == true){
				return 'shake'
			}
			return 'none'
		}
	},
	methods: {
		onSubmit: function(event) {
			event.preventDefault();
			this.shake = false
			setTimeout(function(){
				if (
				this.fake.login == this.login.login &&
				this.fake.password == this.login.password
			) {
				this.alert.message = "Hello Huston !";
			} else {
				this.shake = true;
				this.alert.message = "Huston, we got a problem !";
			}
			},3000)
			console.log(this.shake)
			
		}
	}
});
