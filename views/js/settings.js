window.onload = function () {
    let user = sessionGet("user")
    document.getElementById('name').value = user.Name
    document.getElementById("nickname").value = user.NickName
    document.getElementById("email").value = user.Email
    document.getElementById('gravatar').setAttribute('src', user.AvatarURL)

    let passwordElement = document.getElementById('password')
    let newPasswordElement = document.getElementById('newPassword')
    let confirmPasswordElement = document.getElementById('confirmPassword')
    let submitBtn = document.getElementById('updateBtnTxt').parentElement

    document.getElementById('submit').addEventListener('click', function () {
        let name = document.getElementById('name').value
        let nickName = document.getElementById('nickname').value
        let email = document.getElementById('email').value
        let postJson = {'username': name, 'nickname': nickName, 'email': email}
        if (passwordElement.value != '' && newPasswordElement.value != '' && confirmPasswordElement.value != '') {
            postJson = {
                'name': name,
                'nickname': nickName,
                'email': email,
                'originPassword': passwordElement.value,
                'newPassword': newPasswordElement.value,
                'confirmPassword': confirmPasswordElement.value
            }
        }
        $.ajax({
            url: apiHost + '/user/update',
            method: 'POST',
            dataType: 'JSON',
            data: postJson,
            async: true,
            xhrFields: {
                withCredentials: true
            },
            crossDomain: true,
            beforeSend: function () {
                document.getElementById('updateAnimation').classList.toggle('visually-hidden')
                document.getElementById('updateBtnTxt').classList.toggle('visually-hidden')
            },
            success: function (data) {
                showToast('success', {message: "更新用户信息成功"})
                setTimeout(function () {
                    hiddenToast()
                }, 1500);
            },
            error: function (resp) {
                if (resp.status === 401) {
                    alert("别闹！")
                }
                if (resp.status === 406) {
                    showToast('warning', resp.responseJSON)
                }
            },
            complete: function () {
                document.getElementById('updateAnimation').classList.toggle('visually-hidden')
                document.getElementById('updateBtnTxt').classList.toggle('visually-hidden')
            }
        })
    })

    function changePasswordValid() {
        // 检查原始密码框是否有内容输入
        if (passwordElement.value == '') {
            showValidTooltip(passwordElement)
            submitBtn.setAttribute('disabled', "")
            return;
        }

        let newPassword = newPasswordElement.value
        let confirmPassword = confirmPasswordElement.value

        if (newPassword == '' || confirmPassword == '') {
            submitBtn.setAttribute('disabled', "")
            return
        } else if(newPassword != confirmPassword) {
            showValidTooltip(this)
            submitBtn.setAttribute('disabled', "")
        } else {
            hiddenValidTooltip(this)
            submitBtn.removeAttribute('disabled')

            // 确认密码框提示移除
            let confirmPasswordValidTooltipEle = confirmPasswordElement.parentElement.querySelector(".valid-tooltip")
            let newPasswordValidTooltipEle = newPasswordElement.parentElement.querySelector('.valid-tooltip')
            confirmPasswordValidTooltipEle.style.removeProperty('display')
            confirmPasswordValidTooltipEle.style.removeProperty('position')
            newPasswordValidTooltipEle.style.removeProperty('display')
            newPasswordValidTooltipEle.style.removeProperty('position')
        }
    }

    function showValidTooltip(that) {
        let validTooltip = that.parentElement.querySelector('.valid-tooltip')
        validTooltip.style.display = 'block'
        validTooltip.style.position = 'relative'
    }

    function hiddenValidTooltip(that) {
        let validTooltip = that.parentElement.querySelector('.valid-tooltip')
        validTooltip.style.removeProperty('display')
        validTooltip.style.removeProperty('position')
    }

    passwordElement.addEventListener('input', function () {
        let newPassword = newPasswordElement.value
        let confirmPassword = confirmPasswordElement.value
        if (this.value != '') {
            hiddenValidTooltip(this)
           if (newPassword == '' || confirmPassword == '') {
               submitBtn.setAttribute('disabled', "")
           }
        } else {
            showValidTooltip(this)
        }
    })

    newPasswordElement.addEventListener('input', changePasswordValid)
    confirmPasswordElement.addEventListener('input', changePasswordValid)
}