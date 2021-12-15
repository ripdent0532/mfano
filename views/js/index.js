page = new Page(
    {
        'loadDataFn': loadProject,
        'onPageChangeFn': onPageChange
    })
window.onload = function () {
    page.init()
    page.loadData()
    bindFormValid('uploadForm')
    loadGroup()
}

function onPageChange(num) {
    loadProject(num, getSearchGroupItemId())
}

function loadProject(pagenum, groupId) {
    $.ajax({
        url: apiHost + '/v1/projects',
        method: 'GET',
        dataType: 'json',
        async: true,
        data: {page_num: pagenum, group_id: groupId},
        xhrFields: {
            withCredentials: true
        },
        crossDomain: true,
        success: function (resp) {
            document.getElementById('projectList').innerHTML = ''
            $('#projectsTemplate').tmpl(resp.data).appendTo('#projectList')
            page.options.total = resp['page']['total']
            page.options.currentPage = resp['page']['num']
            page.options.pageSize = resp['page']['size']
            page.renderPageNum()
        },
        error: function (resp) {
            if (resp.status === 401) {
                window.location.href = "/login.html"
            }
        }
    })
}
function submit() {
    let projectGroup;
    if (checkForm('uploadForm')) {
        const formData = new FormData();
        const files = $("#files").prop("files");
        projectGroup = $('#projectGroup').val()
        for (const file of files) {
            formData.append(file.name, file)
        }
        $.ajax({
            url: apiHost + '/v1/project?group=' + projectGroup,
            type: 'POST',
            async: true,
            data: formData,
            cache: false,
            processData: false,
            contentType: false,
            xhrFields: {
                withCredentials: true
            },
            crossDomain: true,
            beforeSend: function () {
                $('#uploadAnimation').toggleClass("visually-hidden")
                $('#uploadButtonTxt').toggleClass("visually-hidden")
            },
            success: function () {
                showToast('success',{"message": "操作成功"})
                setTimeout(function () {
                    hiddenToast()
                    clearInputFileForm()
                    loadProject(0, getSearchGroupItemId())
                    loadGroup()
                }, 1500);
            },
            error: function (resp) {
                if (resp.status === 401) {
                    alert("别闹！")
                }
            },
            complete: function () {
                $('#uploadAnimation').toggleClass("visually-hidden")
                $('#uploadButtonTxt').toggleClass("visually-hidden")
            }
        })
    }
}

function clearInputFileForm() {
    $('#files').val("")
    $('#projectGroup').val("")
}

function loadGroup() {
    $.ajax({
        url: apiHost + '/v1/groups',
        method: 'GET',
        async: true,
        xhrFields: {
            withCredentials: true
        },
        crossDomain: true,
        success: function (resp) {
            if (resp.code === 0) {
                $('#groups').empty()
                $('#searchGroups').empty()
                $('#projectGroupTemplate').tmpl(resp.data).appendTo('#groups')
                $('#projectGroupTemplate').tmpl(resp.data).appendTo('#searchGroups')
            }
        }
    })
}

function search() {
    loadProject(0, getSearchGroupItemId())
}

function getSearchGroupItemId() {
    let selectGroup = $('#searchProjectGroup').val()
    let id
    for (const group of $('#searchGroups').children()) {
        if (group.value === selectGroup) {
            id = group.dataset['id']
            break
        }
    }
    return id
}