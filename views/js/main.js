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