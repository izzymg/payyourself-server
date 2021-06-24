function debugLine(...text) {
    const area = document.querySelector(".debug-area")
    text.forEach(v => {
        area.value += v
        area.value += "\n"
    })
    console.log(...text)
}

function clearDebug() {
    const area = document.querySelector(".debug-area")
    area.value = ""
}

function setNameElement(name) {
    const nameEle = document.querySelector(".logged-in .name")
    nameEle.textContent = `Welcome ${name}, you're signed in!`
}

function setPhotoElement(photoUrl) {
    const nameEle = document.querySelector(".logged-in .photo")
    nameEle.setAttribute("src", photoUrl)
}

function onSignInSuccess(googleUser) {
    debugLine("onSignInSuccess fired")

    const activeUser = googleUser.getBasicProfile()

    setNameElement(activeUser.getName())
    setPhotoElement(activeUser.getImageUrl())

}

function onSignInFail(e) {
    debugLine("onSignInFail fired " + e)
}


function onSignOut() {
    debugLine("onSignOut fired")
}


function onSignOutClicked() {
    debugLine("onSignOutClicked fired")
    try {
        const auth = gapi.auth2.getAuthInstance()
        auth.signOut()
        .then(done => debugLine("signOut done: " + done))
        .catch(e => debugLine("signOut fail:" + e))
    } catch(e) {
        debugLine("signOut failed: " + e)
    }
}

function getUser() {
    const user = gapi.auth2.getAuthInstance().currentUser.get()
    if(!user) {
        throw "No user"
    }
    return user
}

function getToken() {
    const token = getUser().getAuthResponse().id_token
    if(!token) {
        throw "No token"
    }
    return token
}

async function getUserSave() {
    debugLine("getting user save")
    
    const token = getToken().toString()
    debugLine("got token")

    try {
        debugLine("fetch()")
        const res = await fetch("http://localhost:5000/v1", {
            method: "GET",
            headers: {
                "g-token": token,
            },
        })
        debugLine("fetch done")
        if(res.status != 200) {
            console.debug(res)
            debugLine(`fetch returned bad status ${res.status}: ${res.statusText} `)
        }
    } catch(e) {
        debugLine("fetch network failed", e)
    }

}
function removeUserSave() {
    debugLine("removing user save")
}
function saveUserSave() {
    debugLine("saving user save")
}


function appInit() {
    clearDebug()
    debugLine("appInit fired")
    gapi.load("auth2", function() {
        debugLine("auth2 loaded")
        gapi.auth2.init({
            client_id: "43786040065-agi8tcvp56den4bjehkq7cpgovjgjkdk.apps.googleusercontent.com"
        })
        debugLine("auth2 initialized")
        
        gapi.signin2.render("my-signin2", {
            scope: "profile",
            width: 240,
            height: 50,
            onsuccess: onSignInSuccess,
            onfailure: onSignInFail,
        })
        gapi.auth2.getAuthInstance().isSignedIn.listen(
            (isSignIn) => {
                if(!isSignIn) {
                    onSignOut()
                }
        })

        const signOutBtn = document.querySelector(".sign-out-btn")
        signOutBtn.addEventListener("click", onSignOutClicked)

        document.querySelector(".get-user-save-btn").addEventListener("click", getUserSave)
        document.querySelector(".remove-user-save-btn").addEventListener("click", removeUserSave)
        document.querySelector(".save-user-save-btn").addEventListener("click", saveUserSave)
    })
}