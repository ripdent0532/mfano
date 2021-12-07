function Page(options) {
    let defaultOptions = {
        elementId: 'pagination',
        pageNumId: 'page-num',
        currentPage: 1,
        totalPage: 1,
        loadDataFn: 'loadData',
        onPageChangeFn: 'onPageChange'
    }
    this.options = mixin(defaultOptions, options)
}

function mixin(origin, src) {
    let dest = origin
    Object.keys(src).forEach(function (key, index) {
        dest[key] = src[key]
    })
    return dest
}

Page.prototype = (function () {
    function initPage() {
        let pageElement = document.getElementById(this.options.elementId)
        pageElement.style.marginTop = "1rem"
        pageElement.style.marginRight = "4rem"

        let ulElement = document.createElement('ul')
        ulElement.classList.add("pagination", "pagination-sm", "justify-content-end")
        pageElement.appendChild(ulElement)

        let prevPageHtml = '<li class="page-item page-action prev-page disabled"><a class="page-link" href="#">上一页</a></li>'
        let nextPageHtml = '<li class="page-item page-action next-page"><a class="page-link" href="#">下一页</a></li>'
        let contextPageHtml = '<li class="page-item disabled"><a class="page-link" id="page-num">1/7</a></li>'
        ulElement.innerHTML = prevPageHtml + contextPageHtml + nextPageHtml

        bindClickEven.apply(this)
    }

    function bindClickEven() {
        let self = this
        let pageElement = document.getElementById(this.options.elementId)
        let pageBoxes = pageElement.querySelectorAll('.page-action')
        for (let pageBox of pageBoxes) {
            pageBox.addEventListener('click', function () {
                // 判断是否可以点击
                if (this.classList.contains('disabled')) {
                    return
                }
                let currentPage = self.options.currentPage
                let classList = this.classList
                if (classList.contains('next-page')) {
                    self.options.currentPage = currentPage + 1
                } else {
                    self.options.currentPage = currentPage - 1
                }
                renderPageNum.apply(self)
                self.options.onPageChangeFn(self.options.currentPage)
            })
        }
    }

    function loadData() {
        this.options.loadDataFn.apply()
    }

    function renderPageNum() {
        calculatePage.apply(this)
        let elementId = this.options.elementId
        let currentPage = this.options.currentPage
        let totalPage = this.options.totalPage
        let nextPageElement = document.getElementById(elementId).querySelector('.next-page')
        let prevPageElement = document.getElementById(elementId).querySelector('.prev-page')
        if (currentPage == 0) {
            prevPageElement.classList.add('disabled')
        } else {
            prevPageElement.classList.remove('disabled')
        }
        if (currentPage + 1 == totalPage) {
            nextPageElement.classList.add('disabled')
        } else {
            nextPageElement.classList.remove('disabled')
        }


        document.getElementById(this.options.pageNumId).innerHTML = (currentPage + 1) + '/' + totalPage
    }

    function calculatePage() {
        let total = this.options.total
        let size = this.options.pageSize
        this.options.totalPage = Math.ceil(total / size) == 0 ? 1 : Math.ceil(total / size)
    }

    return {
        constructor:Page,
        init: function () {
            this._(initPage)
        },
        loadData: function () {
            this._(loadData)
        },
        renderPageNum: function () {
            this._(renderPageNum)
        },
        _: function (fun, arguments) {
            const that = this
            return fun.apply(that, arguments)
        }
    }
})();