
function setFileName(url) {
    console.log(url)
    const fileName = url.split(/(\\|\/)/g).pop();
    let description = document.querySelector('.description')
    description.innerText = fileName
    description.style.marginTop = "10px"

    let submitImg = document.getElementById('submitImg')
    if (description.innerText != "") {
        submitImg.style.display = "flex"
    }
}
window.onload = () => {
    let currentNotification = null;
    let socket = new WebSocket("ws://localhost:8080/ws");
    let postContainer = document.getElementById('containerPost');
    let notif = document.querySelector('.notification');
    let privateChatDiv = document.querySelector('.privateChatDiv');
    socket.addEventListener('message', function (event) {
        let data = JSON.parse(event.data);
        console.log(data, "here is the data")
        //Check if data is an array of posts
        if (data.type == "posts") {
            //delete all posts
            postContainer.innerHTML = "";
            data.data.forEach(post => {
                // Create a div for each post
                let div = document.createElement('div');
                div.classList.add('postsDiv');
                // Create a title element
                let title = document.createElement('h2');
                title.classList.add('post-title');
                title.innerText = post.title;
                // Create a category element
                let category = document.createElement('h3');
                category.classList.add('postCateg');
                let result = ""
                for (let i = 0; i < post.category.length; i++) {
                    if (post.category[i] == " ") {
                        result += " #"
                    } else {
                        result += post.category[i]
                    }
                }
                category.innerText = "#" + result;
                // Create a body element
                let body = document.createElement('p');
                body.classList.add('postBody');
                body.innerText = post.body;
                // Create a croix element
                // Create a date element
                let date = document.createElement('p');
                date.classList.add('date');
                date.innerText = post.postime;
                // Create a author element
                let username = document.createElement('p');
                username.classList.add('username');
                username.innerText = post.username || "Jesus";
                // create form for comments
                let formComments = document.createElement('form');
                formComments.classList.add('formComments');
                formComments.setAttribute('method', 'POST');
                formComments.onsubmit = function (event) {
                    event.preventDefault();
                    let formData = new FormData(formComments);
                    let comment = formData.get('comment');
                    socket.send(JSON.stringify({
                        type: "comment",
                        data: {
                            body: comment,
                            postid: post.id
                        }
                    }));
                    formComments.reset();
                }

                // create input for comments
                let inputComments = document.createElement('input');
                inputComments.classList.add('inputComments');
                inputComments.setAttribute('type', 'text');
                inputComments.setAttribute('name', 'comment');
                inputComments.setAttribute('placeholder', 'Comment');
                // create div for comments
                let commentsDiv = document.createElement('div');
                commentsDiv.classList.add('commentsDiv');
                // create button for comments
                let buttonComments = document.createElement('button');
                buttonComments.classList.add('buttonComments');
                buttonComments.setAttribute('type', 'submit');
                buttonComments.innerText = "Send";
                post.comments = post.comments || [];
                let numbersOfComments = document.createElement('p');
                numbersOfComments.classList.add('numbersOfComments');
                numbersOfComments.innerText = post.comments.length + " comment(s)";
                let fleche = document.createElement('div');
                fleche.classList.add('fleche');
                fleche.innerText = "↑";
                post.comments.forEach(comment => {
                    let commentDiv = document.createElement('div');
                    commentDiv.classList.add('commentDiv');
                    let commentBody = document.createElement('p');
                    commentBody.classList.add('commentBody');
                    commentBody.innerText = comment.body;
                    let commentDate = document.createElement('p');
                    commentDate.classList.add('commentDate');
                    //todays date with specific format
                    let today = new Date();
                    let dd = String(today.getDate()).padStart(2, '0');
                    let mm = String(today.getMonth() + 1).padStart(2, '0'); //January is 0!
                    let yyyy = today.getFullYear();
                    today = mm + '/' + dd + '/' + yyyy;
                    commentDate.innerText = "posted the: " + today;
                    let commentUsername = document.createElement('p');
                    commentUsername.classList.add('commentUsername');
                    commentUsername.innerText = comment.username;
                    commentDiv.appendChild(commentUsername);
                    commentDiv.appendChild(commentBody);
                    commentDiv.appendChild(commentDate);
                    commentsDiv.prepend(commentDiv);
                });
                fleche.addEventListener('click', function () {
                    if (inputComments.style.display == "none") {
                        numbersOfComments.style.display = "none";
                        fleche.innerText = "↑";
                        inputComments.style.removeProperty('display');
                        buttonComments.style.removeProperty('display');
                        commentsDiv.style.display = "block";
                    } else {
                        fleche.innerText = "↓";
                        numbersOfComments.style.display = "block";
                        inputComments.style.display = "none";
                        buttonComments.style.display = "none";
                        commentsDiv.style.display = "none";
                    }
                })
                div.appendChild(title);
                div.appendChild(category);
                div.appendChild(body);
                div.appendChild(date);
                div.appendChild(username);
                div.appendChild(formComments);
                div.appendChild(commentsDiv);
                formComments.appendChild(inputComments);
                formComments.appendChild(buttonComments);
                div.appendChild(commentsDiv);
                postContainer.prepend(div);
                div.appendChild(numbersOfComments);
                div.appendChild(fleche);
            });
        } else if (data.type == "messages") {
            messageContainer = document.querySelector('#messages');
            messageContainer.innerHTML = "";
            console.log(data.data, "here is the data")
            let messages = data.data.messages;
            messages.forEach((message, index) => {
                let div = document.createElement('div');
                div.classList.add('messageDiv');
                let notificationBody = document.querySelector('.notificationBody');
                //split message.username to ge the last word
                if (index === messages.length - 1 && data.data.initial != true && checkCookie("session")) {
                    let username = message.username.split(" ")[message.username.split(" ").length - 1]
                    console.log(username, chatUsername.innerText, "here are the usernames")
                    if (username != chatUsername.innerText) {
                        notificationBody.innerHTML = '<p class ="actualMessage">' + "<span style='color:rgba(30, 84, 114, 0.726);'>" + message.username + "</span>" + " : " + message.body + "</p>" + '<p class="notificationDate">' + message.created_at + "</p>";
                    }
                }
                let messageBody = document.createElement('p');
                messageBody.classList.add('messageBody');
                messageBody.innerText = message.body;
                let messageDate = document.createElement('p');
                messageDate.classList.add('messageDate');
                messageDate.innerText = message.created_at;
                let messageUsername = document.createElement('p');
                messageUsername.classList.add('messageUsername');
                messageUsername.innerText = message.username;
                div.appendChild(messageUsername);
                div.appendChild(messageBody);
                div.appendChild(messageDate);
                messageContainer.appendChild(div);
                //check if notif is empty
                if (notificationBody.innerText == "") {
                    notif.style.display = "none";
                } else {
                    notif.style.display = "flex";
                    notif.style.animation = "notificationAnimation 5s"
                    if (currentNotification) clearTimeout(currentNotification);
                    currentNotification = setTimeout(function () {
                        currentNotification = null;
                        notif.style.display = "none";
                        notificationBody.innerHTML = "";
                    }, 5000);
                }
            });
        } else if (data.type == "typing") {
            let isTypingBox = document.querySelector('.isTypingBox');
            let notifTyping = document.querySelector('.notifTyping');
            let username = data.data.username;
            notifTyping.innerText = username + " is typing...";
            isTypingBox.style.display = "flex";
            setTimeout(function () {
                isTypingBox.style.display = "none";
            }, 5000);
        } else if (data.type == "status") {
            let whoIsHere = document.querySelector('#whoIsHere');
            console.log(data.data, "here is the status")
            connected = data.data.connected;
            whoIsHere.innerHTML = "";
            let title = document.createElement('p');
            title.classList.add('whoIsHereTitle');
            title.innerText = "Connected Users";
            whoIsHere.appendChild(title);
            connected.forEach(user => {
                //get username and status 
                let username = user.username;
                let status = user.status;
                let div = document.createElement('div');
                div.classList.add('userDiv');
                let userStatus = document.createElement('p');
                userStatus.classList.add('userStatus');
                if (user.status == "Online") {
                    userStatus.innerHTML = "<span class='usernameConnected'>" + username + "</span>" + " " + "<span class='userStatusConnected'>" + status + "</span>";
                } else {
                    userStatus.innerHTML = "<span class='usernameConnected'>" + username + "</span>" + " " + "<span class='userStatusOffline'>" + status + "</span>";
                }

                div.appendChild(userStatus);
                whoIsHere.appendChild(div);
                div.addEventListener('click', function () {
                    console.log(user.username)
                    privateChatDiv.style.display = "flex";
                    let chatUsername = document.querySelector('.UnPm');
                    chatUsername.innerText = user.username;
                });
            })

        }
    });
    let crossPm = document.querySelector('.crossPm');
    crossPm.addEventListener('click', function () {
        privateChatDiv.style.display = "none";
    })



    let formBox = document.getElementsByClassName('formBox')[0];
    formBox.addEventListener('submit', function (event) {
        event.preventDefault();
        let formData = new FormData(formBox);
        let title = formData.get('title');
        let category = formData.get('category');
        let body = formData.get('body');
        console.log(formData, "here is the form data")
        socket.send(JSON.stringify({
            type: "post",
            data: {
                title: title,
                category: category,
                body: body
            }
        }));
        formBox.reset();
    })

    let register = document.getElementById('register');
    let authForm = document.getElementById('authForm');
    let registerForm = document.getElementById('registerForm');
    let signIn = document.getElementById('signIn');
    let forumPage = document.getElementById('forumPage')
    let AuthPage = document.getElementById('authPage')
    let menuBurger = document.getElementById('menuBurger')
    let burgerOpened = document.getElementById('burgerOpened')
    let logOut = document.querySelector('#logOut')
    let errorMsg = document.querySelector('.errorMsg')
    let profilePage = document.getElementById('profilPage')
    let profileBtn = document.getElementById('profil')
    let radioMale = document.querySelector('.radioMale')
    let radioFemale = document.querySelector('.radioFemale')
    let radioOther = document.querySelector('.radioOther')

    //check if cookie session in present
    let forumBurger = document.getElementById('forum')
    forumBurger.addEventListener('click', function () {
        chattingPage.style.display = "none";
    })
    //if radio male is checked female is unchecked
    radioMale.addEventListener('click', function () {
        radioFemale.checked = false;
        radioOther.checked = false;
        radioFemale.removeAttribute('required')
        radioOther.removeAttribute('required')
    })

    radioFemale.addEventListener('click', function () {
        radioMale.checked = false;
        radioOther.checked = false;
        radioMale.removeAttribute('required')
        radioOther.removeAttribute('required')
    })

    radioOther.addEventListener('click', function () {
        radioMale.checked = false;
        radioFemale.checked = false;
        radioMale.removeAttribute('required')
        radioFemale.removeAttribute('required')
    })

    register.addEventListener('click', function () {
        registerForm.style.animation = "none";
        authForm.style.animation = "dissapear 0.5s ease-out"
        setTimeout(function () {
            authForm.style.display = "none";
            registerForm.style.display = "flex";
            registerForm.style.animation = "appear 0.2s ease-in"
        }, 500);
    });

    signIn.addEventListener('click', function () {
        //set cookie to username
        authForm.style.animation = "none"
        registerForm.style.animation = "dissapear 0.5s ease-out"
        setTimeout(function () {
            registerForm.style.display = "none";
            authForm.style.display = "flex";
            authForm.style.animation = "appear 0.2s ease-in"
        }, 500);
    });

    function checkCookie(name) {
        var cookies = document.cookie.split(';');
        for (var i = 0; i < cookies.length; i++) {
            var cookie = cookies[i].trim();
            if (cookie.startsWith(name + '=')) {
                return true;
            }
        }
        return false;
    }

    // Call the function passing in the name of the cookie you want to check
    var isCookiePresent = checkCookie('session');
    if (isCookiePresent) {
        AuthPage.style.display = "none"
        forumPage.style.display = "flex"
    } else {
        AuthPage.style.display = "flex"
        forumPage.style.display = "none"
    }

    menuBurger.addEventListener('click', function () {
        console.log('clicked')
        if (burgerOpened.style.display == "none") {
            closeChat.style.display = "none"
            burgerOpened.style.display = 'flex'
            menuBurger.style.left = '200px'
        } else {
            closeChat.style.display = "flex"
            burgerOpened.style.display = 'none'
            menuBurger.style.left = '0px'
        }
    });

    logOut.addEventListener('click', function () {
        socket.send(JSON.stringify({
            type: "logout"
        }));
        document.cookie = "session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
        location.reload()
    });

    var cookies = document.cookie;

    // parse the string to get the value of the "error" cookie
    var errorCookie = cookies.split(';').find(cookie => cookie.trim().startsWith('error='));
    if (errorCookie) {
        var errorValue = errorCookie.split('"')[1];
        console.log(errorValue);
        errorMsg.innerHTML = errorValue
    } else {
        console.log('No "error" cookie found.');
    }

    addEventListener('load', function () {
        this.setTimeout(function () {
            document.cookie = "error=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
        }, 1000);
    })

    let image = document.querySelector('#profile-img')
    let image_a = document.querySelector('#profile-img-a')
    let username = document.querySelector('#profile-username')
    let nickname = document.querySelector('#profile-nickname')
    let description = document.querySelector('#profile-description')
    let profileCreationDate = document.querySelector('#profile-creation-date')
    let profileBirthday = document.querySelector('#profile-birthday')
    let profileAge = document.querySelector('#profile-age')
    let profileHobby = document.querySelector('#profile-hobby')
    let profileGender = document.querySelector('#profil-gender')

    profileBtn.addEventListener('click', function () {
        //set cookie
        if (profilePage.style.display == "none") {
            forumPage.style.display = "flex"
            burgerOpened.style.display = 'none'
            menuBurger.style.left = '0px'
            document.cookie = "chatt=false; path=/;";
            document.cookie = "profile=true; path=/;";
            profilePage.style.display = "flex"
            //fetch('/profile')
            fetch('/profile', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            })
                .then(response => response.json())
                .then(data => {
                    console.log(data)
                    image_a.href = data.Image
                    image.src = data.Image
                    if (data.Gender == "other") {
                        profileGender.innerText = "Other Gender, "
                    } else {
                        profileGender.innerText = data.Gender + ", "
                    }
                    username.innerText = data.Username
                    nickname.innerText = data.Nickname
                    description.innerText = data.Biography
                    profileCreationDate.innerText = "Profil created the " + data.CreatedAtStr
                    profileBirthday.innerText = data.Birthday
                    profileHobby.innerText = data.Hobby
                    if (profileBirthday.innerText == "Enter your birthday") {
                        profileAge.innerText = ""
                    } else {
                        profileAge.innerText = data.Age + " years old"
                    }

                })
        }
    });

    if (cookies.includes("profile=true")) {
        profilePage.style.display = "flex"
        //fetch('/profile')
        fetch('/profile', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
            .then(response => response.json())
            .then(data => {
                console.log(data)
                image_a.href = data.Image
                image.src = data.Image
                if (data.Gender == "other") {
                    profileGender.innerText = "Other Gender, "
                } else {
                    profileGender.innerText = data.Gender + ", "
                }
                username.innerText = data.Username
                nickname.innerText = data.Nickname
                description.innerText = data.Biography
                profileCreationDate.innerText = "Profil created the " + data.CreatedAtStr
                profileBirthday.innerText = data.Birthday
                profileHobby.innerText = data.Hobby
                if (profileBirthday.innerText == "Enter your birthday") {
                    profileAge.innerText = ""
                } else {
                    profileAge.innerText = data.Age + " years old"
                }
            })
    }

    let profileMenuBurger = document.getElementById('profilemenuBurger')
    let profileBurgerOpened = document.getElementById('profileburgerOpened')
    let forumBtn = document.getElementById('forum')

    profileMenuBurger.addEventListener('click', function () {
        console.log('clicked')
        if (profileBurgerOpened.style.display == "none") {
            profileBurgerOpened.style.display = 'flex'
        } else {
            profileBurgerOpened.style.display = 'none'
        }
    });

    forumBtn.addEventListener('click', function () {
        //set cookie to false
        document.cookie = "profile=false; path=/;";
        profilePage.style.display = "none"
        forumPage.style.display = "flex"
    })

    //if cookie profile=false 
    if (cookies.includes("profile=false")) {
        profilePage.style.display = "none"
    }

    profilelogOut.addEventListener('click', function () {
        document.cookie = "session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
        document.cookie = "profile=false; path=/;";
        location.reload()
    });

    let modifyNick = document.getElementById('modify-nick')
    let modifyBio = document.getElementById('modify-bio')
    let modifyBirthday = document.getElementById('modify-birthday')
    let inputNick = document.getElementById('input-nick')
    let inputBio = document.getElementById('input-bio')
    let inputBirthday = document.getElementById('input-birthday')
    let inputHobby = document.getElementById('input-hobby')
    let modifyHobby = document.getElementById('modify-hobby')

    modifyNick.addEventListener('click', function () {
        if (inputNick.style.display == "none") {
            nickname.style.display = "none"
            inputNick.style.display = "flex"
            inputNick.style.marginTop = "23px"
            modifyNick.style.marginTop = "23px"
            inputNick.style.marginBottom = "21px"
            modifyNick.style.marginBottom = "21px"
        } else {
            nickname.style.display = "flex"
            inputNick.style.display = "none"
            inputNick.style.marginTop = "0px"
            modifyNick.style.marginTop = "0px"
            inputNick.style.marginBottom = "0px"
            modifyNick.style.marginBottom = "0px"

        }
    })


    modifyBio.addEventListener('click', function () {
        if (inputBio.style.display == "none") {
            description.style.display = "none"
            inputBio.style.display = "flex"
            inputBio.style.marginTop = "13px"
        } else {
            description.style.display = "flex"
            inputBio.style.display = "none"
            inputBio.style.marginTop = "0px"
        }
    })

    modifyBirthday.addEventListener('click', function () {
        if (inputBirthday.style.display == "none") {
            profileBirthday.style.display = "none"
            inputBirthday.style.display = "flex"
            inputBirthday.style.marginBottom = "17px"
            modifyBirthday.style.marginBottom = "17px"
            modifyBirthday.style.marginTop = "15px"
            inputBirthday.style.marginTop = "15px"

        } else {
            profileBirthday.style.display = "flex"
            inputBirthday.style.display = "none"
            inputBirthday.style.marginBottom = "0px"
            modifyBirthday.style.marginBottom = "0px"
            modifyBirthday.style.marginTop = "0px"
            inputBirthday.style.marginTop = "0px"
        }
    })

    modifyHobby.addEventListener('click', function () {
        if (inputHobby.style.display == "none") {
            profileHobby.style.display = "none"
            inputHobby.style.display = "flex"
            inputHobby.style.marginBottom = "13px"
            modifyHobby.style.marginBottom = "13px"
            modifyHobby.style.marginTop = "15px"
            inputHobby.style.marginTop = "15px"

        } else {
            profileHobby.style.display = "flex"
            inputHobby.style.display = "none"
            inputHobby.style.marginBottom = "0px"
            modifyHobby.style.marginBottom = "0px"
            modifyHobby.style.marginTop = "0px"
            inputHobby.style.marginTop = "0px"
        }
    })



    let chatBox = document.querySelector('.chat')
    forumPage = document.getElementById('forum-triplePage')
    let chattingPage = document.getElementById('chattingPage')

    chatBox.addEventListener('click', function () {
        document.cookie = "chatt=true; path=/;";
        document.cookie = "profile=false; path=/;";
        console.log('clicked')
        if (forumPage.style.display == "none") {
            profilePage.style.display = "none"
            forumPage.style.display = "flex"
            chattingPage.style.display = "none"
        } else {
            profilePage.style.display = "none"
            forumPage.style.display = "none"
            chattingPage.style.display = "flex"
        }
    });

    notif.addEventListener('click', function () {
        if (chattingPage.style.display == "none") {
            document.cookie = "chatt=true; path=/;";
            document.cookie = "profile=false; path=/;";
            console.log('clicked')
            if (forumPage.style.display == "none") {
                profilePage.style.display = "none"
                forumPage.style.display = "flex"
                chattingPage.style.display = "none"
            } else {
                profilePage.style.display = "none"
                forumPage.style.display = "none"
                chattingPage.style.display = "flex"
            }
        }
    });

    let closeChat = document.getElementById('closeChat')
    cookies = document.cookie
    if (cookies.includes("chatt=true")) {
        forumPage.style.display = "none"
        profilePage.style.display = "none"
        chattingPage.style.display = "flex"
    } else {
        chattingPage.style.display = "none"
    }

    closeChat.addEventListener('click', function () {
        cookies = document.cookie = "chatt=false; path=/;";
        console.log('clicked')
        if (forumPage.style.display == "none") {
            forumPage.style.display = "flex"
            chattingPage.style.display = "none"
        } else {
            forumPage.style.display = "none"
            chattingPage.style.display = "flex"
        }
    });

    let sendMessage = document.getElementById('sendMessage')
    let chatUsername = document.getElementById('chatUsername')
    //get cookie name username
    cookies = document.cookie
    if (cookies.includes("username=")) {
        let username = cookies.split("username=")[1].split(";")[0]
        chatUsername.innerText = username
        console.log(username)
        console.log(chatUsername.innerText)
    }

    let writeMessage = document.getElementById('writeMessage')
    let chatHereBox = document.querySelector('#chatHere')

    chatHereBox.addEventListener('keydown', function (event) {
        socket.send(JSON.stringify({ type: "typing" }));
    });

    setTimeout(function () {
        if (checkCookie('session')) {
            console.log("connected")
            socket.send(JSON.stringify({
                type: "login",
            }));
        }
    }, 1000);

    sendMessage.addEventListener('click', function (event) {
        if (chatHereBox.value != "") {
            event.preventDefault()
            let formData = new FormData(writeMessage);
            let message = formData.get('chatHere');
            let username = chatUsername.innerText
            let time = new Date().toLocaleTimeString()
            console.log(message)
            console.log(username)
            console.log(time)
            socket.send(JSON.stringify({
                type: "message",
                data: {
                    body: message,
                    username: username,
                    time: time
                }
            }));
            message = document.getElementById('chatHere').value = ""
            writeMessage.reset()
        }
    });

    document.addEventListener("keypress", function (event) {
        if (event.key === "Enter" && chattingPage.style.display == "flex") {
            if (chatHereBox.value != "") {
                event.preventDefault()
                let formData = new FormData(writeMessage);
                let message = formData.get('chatHere');
                let username = chatUsername.innerText
                let time = new Date().toLocaleTimeString()
                console.log(message)
                console.log(username)
                console.log(time)
                socket.send(JSON.stringify({
                    type: "message",
                    data: {
                        body: message,
                        username: username,
                        time: time
                    }
                }));
                message = document.getElementById('chatHere').value = ""
                writeMessage.reset()
            }
        }
    });

};
