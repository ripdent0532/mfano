page = new Page({
    loadDataFn: loadUsers,
    onPageChangeFn: function (num) {
        const username = document.getElementById('searchUserName').value
        loadUsers(num, username)
    }
})

$(function () {
    page.init()
    page.loadData()
    bindFormValid('addUserForm')
})

function loadUsers(pageNum, username) {
    $.ajax({
        url: apiHost + '/users',
        method: 'GET',
        dataType: 'JSON',
        data: {'pageNum': pageNum, 'userName': username},
        async: true,
        xhrFields: {
            withCredentials: true
        },
        crossDomain: true,
        success: function (resp) {
            $('#userList').empty()
            $('#usersTemplate').tmpl(resp.data).appendTo('#userList')
            page.options.total = resp.total
            page.options.currentPage = resp.num
            page.options.pageSize = resp.size
            page.renderPageNum()
        }
    })
}

function search() {
    const username = document.getElementById('searchUserName').value
    loadUsers(0, username)
}

function addUser() {
    let username;
    let nickname;
    let email;
    if (checkForm('addUserForm')) {
        username = document.getElementById('username').value
        nickname = document.getElementById('nickname').value
        email = document.getElementById('email').value
        $.ajax({
            url: apiHost + '/user',
            method: 'POST',
            dataType: 'JSON',
            data: {'username': username, 'nickname': nickname, 'email': email},
            async: true,
            xhrFields: {
                withCredentials: true
            },
            crossDomain: true,
            beforeSend: function () {
                document.getElementById('addAnimation').classList.toggle('visually-hidden')
                document.getElementById('addBtnTxt').classList.toggle('visually-hidden')
            },
            success: function (data) {
                showToast('success', {"message": "操作成功"})
                setTimeout(function () {
                    hiddenToast()
                    clearInputForm()
                    loadUsers()
                }, 1500);
            },
            error: function (resp) {
                if (resp.status === 401) {
                    alert("别闹！")
                }
            },
            complete: function () {
                document.getElementById('addAnimation').classList.toggle('visually-hidden')
                document.getElementById('addBtnTxt').classList.toggle('visually-hidden')
            }
        })
    }
}

function clearInputForm() {
    $('#username').val("")
    $('#nickname').val("")
    $('#email').val("")
}

