const apiHost = "http://localhost:8080"
$(function () {
    includes = $('[data-include]');
    jQuery.each(includes, function(){
        file = $(this).data('include') + '.html';
        $(this).load(file);
    });
})

function bindFormValid(formId) {
    let elements = document.getElementById(formId).querySelectorAll("input[required]")
    for (let i = 0; i < elements.length; i++) {
        elements[i].addEventListener('invalid', function (e) {
           this.style.backgroundColor = '#ffdddd'
        })

        elements[i].addEventListener('focus', function (e) {
            this.style.removeProperty('background-color')
        })
    }
}

function checkForm(formId) {
    let pass = true
    inputs = document.getElementById(formId).querySelectorAll("input")
    inputs.forEach(function (input) {
        if (!input.checkValidity()) {
            pass = false
        }
    })
    return pass
}

function sessionSet(key, value) {
    if (typeof value === 'object') {
        value = JSON.stringify(value)
    }
    localStorage.setItem(key, value)
}

function sessionGet(key) {
    const value = localStorage.getItem(key) || ''
    return JSON.parse(value)
}

function sessionRemove(key) {
    localStorage.removeItem(key)
}

function sessionClearAll() {
    localStorage.clear()
}

function queryParam(queryMap) {
    let queryParamArray = new Array()
    queryMap.forEach(function (value, key) {
        if (value !== undefined && value != '') {
            queryParamArray.push(key + "=" + value)
        }
    })
    let queryParam = ''
    if (queryParamArray.length > 0) {
        queryParam = "?" + queryParamArray.join("&")
    }
    return queryParam
}

function formatDate(date) {
    var arr = date.split("T");
    var d = arr[0];
    var darr = d.split('-');
    var t = arr[1];
    var tarr = t.split('.');
    var marr = tarr[0].split(':');
    var dd = parseInt(darr[0]) + "/" + parseInt(darr[1]) + "/" + parseInt(darr[2]) + " " + parseInt(marr[0]) + ":" + parseInt(marr[1]) + ":" + parseInt(marr[2]);
    return formatDateTime(dd);
}

function formatDateTime(date) {
    let time = new Date(Date.parse(date));
    let Y = time.getFullYear() + '-';
    let M = addZero(time.getMonth() + 1) + '-';
    let D = addZero(time.getDate()) + ' ';
    let h = addZero(time.getHours()) + ':';
    let m = addZero(time.getMinutes()) + ':';
    let s = addZero(time.getSeconds());
    return Y + M + D + h + m + s;
}
function addZero(num) {
    return num < 10 ? '0' + num : num;
}