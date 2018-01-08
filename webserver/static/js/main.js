

$(document).ready(function () {
    var nicknameInput = $("#nickname-input");
    if (nicknameInput !== null) {
        nicknameInput.val(getCookie("nickname"));
    }
    $('#login-form').on('submit', function (event) {
        event.preventDefault();
        login();
    });
    $('#user-creation-form').on('submit', function (event) {
        event.preventDefault();
        createUser();
    });
    $('#change-password-form').on('submit', function (event) {
        event.preventDefault();
        changePassword();
    });
    $('.user-creation-input').on('keyup', function (event) {
        $("#user-creation-error").hide();
    });
    $('.login-input').on('keyup', function (event) {
        hideCredentialsError();
    });
    $('#edit-post-form').on('submit', function (event) {
        event.preventDefault();
        editPost();
    });
});

function login() {
    $.ajax({
        url: "/login",
        type: "POST",
        data: $('#login-form').serialize(),
        success: function (result) {
            if (result === "success") {
                $("#credentials-valid").show();
                setTimeout(function () {
                    location.reload();
                }, 1000);
            } else {
                $("#credentials-invalid").show();
            }
        },
        error: function (err) {
            $("#credentials-error").show();
            // alert(err);
        }
    });
}

function hideCredentialsError() {
    $("#credentials-invalid").hide();
    $("#credentials-error").hide();
}

function addTag() {
    var input = document.getElementById("keyword-input");
    if (input.value.length === 0) {
        document.getElementById("keyword-input").style.border = "1px solid red";
        return;
    }
    var elementCount = document.getElementById("tag-container").childElementCount;
    var container = document.createElement("div");
    container.classList.add("input-group");
    container.classList.add("entry-tag-container");
    container.id = "tag-container-" + elementCount;
    var tag = document.createElement("input");
    tag.value = input.value;
    tag.classList.add("form-control");
    tag.classList.add("entry-tag");
    tag.name = "tag";
    tag.readOnly = true;
    var span = document.createElement("span");
    var button = document.createElement("button");
    span.classList.add("input-group-btn");
    button.classList.add("btn");
    button.classList.add("btn-secondary");
    button.innerHTML = "x";
    button.onclick = function () {
        removeTag("tag-container-" + elementCount)
    };
    span.appendChild(button);
    container.appendChild(tag);
    container.appendChild(span);
    input.value = "";
    document.getElementById("tag-container").appendChild(container);
}

function resetTagInput() {
    document.getElementById("keyword-input").style.border = "1px solid rgba(0,0,0,.15)";
}

function removeTag(elementId) {
    document.getElementById(elementId).remove();
}

function saveNickname() {
    var nickname = document.getElementById("nickname-input").value;
    document.cookie = "nickname=" + nickname;
}

function getCookie(name) {
    match = document.cookie.match(new RegExp(name + '=([^;]+)'));
    if (match) return match[1];
}

function verifyComment(postId, commentId) {
    $.ajax({
        url: "?verify",
        type: "POST",
        data: {"postId": postId, "commentId": commentId},
        success: function (result) {
            if (result === "true") {
                $("#not-verified-status-" + commentId).replaceWith("<span class='verification-status verification-verified'>Verified</span>");
            }
        },
        error: function (err) {
            alert(err);
        }
    });
}

function deletePost(postId) {
    $.ajax({
        url: "?delete",
        type: "POST",
        data: {"postId": postId},
        success: function (result) {
            if (result === "true"){
                window.location = "/"
            } else {
                alert("Something went wrong.")
            }
        },
        error: function (err) {
            alert(err);
        }
    });
}


function createUser() {
    $.ajax({
        url: "?newUser",
        type: "POST",
        data: $('#user-creation-form').serialize(),
        success: function (result) {
            var msg = result.split("#");
            if (msg[1].length === 0){
                alert("User '" + msg[0] + "' created!");
                $("#user-creation-form")[0].reset();
            } else {
                var err = $("#user-creation-error");
                err.text(msg[1]);
                err.css('display','block');
            }
        },
        error: function (err) {
            alert(err);
        }
    });
}

function changePassword() {
    $.ajax({
        url: "?password",
        type: "POST",
        data: $('#change-password-form').serialize(),
        success: function (result) {
            if (result.length === 0){
                alert("Password successfully changed!");
                window.location = "/";
            } else {
                var err = $("#change-password-error");
                err.text(result);
                err.css('display','block');
            }
        },
        error: function (err) {
            alert(err);
        }
    });
}

function requestMorePosts(index){
    $.ajax({
        url: "?more=" + index,
        type: "GET",
        success: function (result) {
            $("#more-content-placeholder").replaceWith(result)
        },
        error: function (err) {
            alert(err);
        }
    });
}

