
const setInfo = function (UserInfo) {
    window.localStorage.setItem('UserInfo', JSON.stringify(UserInfo))
}


const getInfo = function () {
    return JSON.parse(window.localStorage.getItem('UserInfo'))
}

export { setInfo, getInfo }