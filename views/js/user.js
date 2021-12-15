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
    let queryMap = new Map()
    queryMap.set('page_num', pageNum)
    queryMap.set('username', username)
    $.ajax({
        url: apiHost + '/v1/users' + queryParam(queryMap),
        method: 'GET',
        dataType: 'json',
        async: true,
        xhrFields: {
            withCredentials: true
        },
        crossDomain: true,
        success: function (resp) {
            $('#userList').empty()
            if (resp["data"] != null) {
                $('#usersTemplate').tmpl(resp.data).appendTo('#userList')
            }
            page.options.total = resp["page"].total
            page.options.currentPage = resp["page"].num
            page.options.pageSize = resp["page"].size
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
            url: apiHost + '/v1/user',
            method: 'POST',
            dataType: 'JSON',
            data: JSON.stringify({'username': username, 'nickname': nickname, 'email': email}),
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
                if (data.code === 0) {
                    showToast('success', {"message": data["msg"]})
                }
                if (data.code === 40001) {
                    showToast('warning', {"message": data["msg"]})
                }

                setTimeout(function () {
                    hiddenToast()
                    clearInputForm()
                    loadUsers()
                }, 1500);
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

