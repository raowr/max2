/**
 * /^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,11}$/ //数字和字母组合
 * 验证用户名是否合法
 * @param val
 * @returns {boolean}
 */
export const check_userName = val => {
    let telreg = /^[0-9A-Za-z]{6,11}$/;
    if (!telreg.test(val)) {
        return false;
    } else {
        return true;
    }
}
/**
 * 验证手机号格式
 * @param phone
 * @returns {boolean}
 */
export const check_Mobile = phone => {
    let telreg = /^[1][3,4,5,6,7,8,9][0-9]{9}$/;
    phone = phone.toString().replace(/\s*/g, "")
    if (!telreg.test(phone)) {
        return false;
    } else {
        return true;
    }
}
/**
 * 验证短信验证码格式
 * @param smsCode
 * @returns {boolean}
 */
export const check_smsCode = smsCode =>{
    let telreg = /^[0-9]{6}$/;
    if (!telreg.test(smsCode)) {
        return false;
    } else {
        return true;
    }
}
/**
 * 验证中文真实姓名
 * @param val
 * @returns {boolean}
 */
export const check_realName = val => {
    let reg = /^[\u4e00-\u9fa5]{2,}$/;
    if (reg.test(val)) {
        return true;
    } else {
        return false;
    }
}

/**
 * 验证邮箱格式
 * @param mail
 * @returns {boolean}
 */
export const check_Mail = mail => {
    let telreg = /^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$/;
    if (!telreg.test(mail)) {
        return false;
    } else {
        return true;
    }
}
/**
 * 判断是否为空
 * @param Content
 * @returns {boolean}
 */
export const checkAllNull = Content => {
    let reg = /^[\s]+$/;
    if (reg.test(Content)) {
        return false;
    } else {
        return true;
    }
}

/**
 * 判断是否为有空格
 * @param Content
 * @returns {boolean}
 */
export const checkNull = Content =>{
    let reg = /[\s]+/;
    if (reg.test(Content)) {
        return false;
    } else {
        return true;
    }
}

/**
 * 判断是否是纯数字
 * @param val
 * @returns {boolean}
 */
export const check_Number = num => {
    let reg = /^[0-9]+$/;
    if (!reg.test(num)) {
        return false;
    } else {
        return true;
    }
}

/**
 * 判断是否是数字或者带小数的数字
 * @param val
 * @returns {boolean}
 */
 export const check_FloatNumber = (num) => {
    let reg = /^[1-9]\d{0,}\.{0,1}\d{0,}$|^[0]\.\d{0,}[1-9]$/
    if (!reg.test(num)) {
        return false;
    } else {
        return true;
    }
};

/**
 * 验证密码格式 6-12 数字加字母组合
 * @param password
 * @returns {boolean}
 */
export const check_pwd = password => {
    let pwd = /^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,12}$/;
    if (!pwd.test(password)) {
        return false;
    } else {
        return true;
    }
}

/**
 * 验证输入格式 是否为正整数
 * @param password
 * @returns {boolean}
 */
export const check_positiveInt = num => {
    let regu = /^[1-9]\d*$/;
    if (!regu.test(num)) {
        return false;
    } else {
        return true;
    }
}

/**
 * 浏览器类型
 * @returns {string}
 */
export const check_browser = () => {
    var ua = navigator.userAgent.toLowerCase(); //取得浏览器的userAgent字符串
    if (ua.indexOf("opera") > 0) {
        return 'opera';
    } else if (ua.indexOf("baidu") > 0) {
        return 'baidu';
    } else if (ua.indexOf("qq") > 0) {
        //qq logic
        return 'qq';
    } else if (ua.indexOf("liebao") > 0) {
        //liebao logic
        return 'liebao';
    } else if (ua.indexOf("uc") > 0) {
        //uc logic
        return 'uc';
    } else if (ua.indexOf("chrome") > 0) {
        //chrome logic
        return 'chrome';
    } else if (ua.indexOf("safari") > 0) {
        //safari logic
        return 'safari';
    } else {
        //others
        return 'other';
    }
}
//获取url 参数
export const urlparams = url => {
    url = url || document.location.href;
    var data = {};
    var pos = url.indexOf('?');
    if (pos != -1) {
        var query = url.substring(pos);
        if (query) {
            var ll = query.substring(1).split('&');
            var i = 0;
            for (; i < ll.length; i++) {
                var arr = ll[i].split('=');
                if (arr.length == 2)
                    data[arr[0]] = decodeURI(arr[1]);
            }
        }
    }
    return data;
}
//获取时间
export const getForMatDate = day => {
    var today = new Date();
    var targetday_milliseconds = today.getTime() + 1000 * 60 * 60 * 24 * day;
    today.setTime(targetday_milliseconds);
    var tYear = today.getFullYear();
    var tMonth = today.getMonth();
    var tDate = today.getDate();
    tMonth = doHandleMonth(tMonth + 1);
    tDate = doHandleMonth(tDate);
    return tYear + "-" + tMonth + "-" + tDate;
}
//获取一个月前时间，格式YYYY-MM-DD
export const getMonthFormatDate = () => {
    return getForMatDate(-29)
}
//获取一周前时间，格式YYYY-MM-DD
export const getsunFormatDate = () => {
    return getForMatDate(-6)
}
//获取三天前时间，格式YYYY-MM-DD
export const getshrFormatDate = () => {
    return getForMatDate(-2)
}
//获取昨天时间，格式YYYY-MM-DD
export const getYesToDate = () => {
    return getForMatDate(-1)
}
//获取当前时间，格式YYYY-MM-DD hh:mm:ss
export const getNowFormatDateAndTime = () => {
    var date = new Date();
    var seperator1 = "-";
    var seperator2 = ":";
    var month = date.getMonth() + 1;
    var strDate = date.getDate();
    if (month >= 1 && month <= 9) {
        month = "0" + month;
    }
    if (strDate >= 0 && strDate <= 9) {
        strDate = "0" + strDate;
    }
    var currentdate = date.getFullYear() + seperator1 + month + seperator1 + strDate
        + " " + (date.getHours()<10?'0'+date.getHours():date.getHours()) + seperator2 + (date.getMinutes()<10?'0'+date.getMinutes():date.getMinutes())
        + seperator2 + (date.getSeconds()<10?'0'+date.getSeconds():date.getSeconds());
    return currentdate;
}
export const doHandleMonth = month => {
    var m = month;
    if (month.toString().length == 1) {
        m = "0" + month;
    }
    return m;
}
export const trimStr =str => {
    return str.replace(/(^\s*)|(\s*$)/g, "");
}


/**
 * @param {Number} times 回调函数需要执行的次数
 * @param {Function} callback 回调函数
 */
export const doCustomTimes = (times, callback) => {
    let i = -1
    while (++i < times) {
        callback(i)
    }
}


/**
 * @param {*} obj1 对象
 * @param {*} obj2 对象
 * @description 判断两个对象是否相等，这两个对象的值只能是数字或字符串
 */
export const objEqual = (obj1, obj2) => {
    const keysArr1 = Object.keys(obj1)
    const keysArr2 = Object.keys(obj2)
    if (keysArr1.length !== keysArr2.length) return false
    else if (keysArr1.length === 0 && keysArr2.length === 0) return true
    /* eslint-disable-next-line */
    else return !keysArr1.some(key => obj1[key] != obj2[key])
}

/**
 * @description 根据name/params/query判断两个路由对象是否相等
 * @param {*} route1 路由对象
 * @param {*} route2 路由对象
 */
export const routeEqual = (route1, route2) => {
    const params1 = route1.params || {}
    const params2 = route2.params || {}
    const query1 = route1.query || {}
    const query2 = route2.query || {}
    return (route1.name === route2.name) && objEqual(params1, params2) && objEqual(query1, query2)
}

/**
 * 判断打开的标签列表里是否已存在这个新添加的路由对象
 */
export const routeHasExist = (cacheList, routeItem) => {
    let len = cacheList.length
    let res = false
    doCustomTimes(len, (index) => {
        if (routeEqual(cacheList[index], routeItem)) res = true
    })
    return res
}

/**
 * @param {*} list 现有标签导航列表
 * @param {*} newRoute 新添加的路由原信息对象
 * @description 如果该newRoute已经存在则不再添加
 */
export const getNewCacheList = (list, newRoute) => {
    const { name, path, meta } = newRoute
    let newList = [...list]
    if (routeHasExist(newList,newRoute) || (newRoute.meta && !newRoute.meta.cache) ) return newList
    else newList.push({ name, path, meta })
    return newList
}

export const localSave = (key, value) => {
    localStorage.setItem(key, value)
}

export const localRead = (key) => {
    return localStorage.getItem(key) || ''
}

// scrollTop animation
export const scrollTop = (el, from = 0, to, duration = 500, endCallback) => {
    if (!window.requestAnimationFrame) {
        window.requestAnimationFrame = (
            window.webkitRequestAnimationFrame ||
            window.mozRequestAnimationFrame ||
            window.msRequestAnimationFrame ||
            function (callback) {
                return window.setTimeout(callback, 1000 / 60)
            }
        )
    }
    const difference = Math.abs(from - to)
    const step = Math.ceil(difference / duration * 50)

    const scroll = (start, end, step) => {
        if (start === end) {
            endCallback && endCallback()
            return
        }

        let d = (start + step > end) ? end : start + step
        if (start > end) {
            d = (start - step < end) ? end : start - step
        }

        if (el === window) {
            window.scrollTo(d, d)
        } else {
            el.scrollTop = d
        }
        window.requestAnimationFrame(() => scroll(d, end, step))
    }
    scroll(from, to, step)
}

/**
 * @description 本地存储和获取缓存视图
 */
export const setCacheListInLocalstorage = list => {
    localStorage._cacheList = JSON.stringify(list)
}
/**
 * @returns {Array} 其中的每个元素只包含路由原信息中的name, path, meta三项
 */
export const getCacheListFromLocalstorage = () => {
    const list = localStorage._cacheList
    return list ? JSON.parse(list) : []
}
/**
 * 检测是否移动端
 * @returns {boolean}
 */
export const isH5 = () => {
    if((navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i))){
        return true;
    }else{
        return false;
    }
}
/**
 * 是否是安卓或ios
 * @returns {number}
 */
export const androidOrIOS = () => {
    var u = navigator.userAgent;
    var osType;
    if (u.indexOf('Android') > -1 || u.indexOf('Linux') > -1) {
        // 安卓手机
        osType = 1
    } else if (u.indexOf('iPhone') > -1 ||  u.indexOf('iPad') > -1) {
        osType = 2
        // 苹果手机
    }else {
        osType = 0
    }
    return osType;
}
/**
 * 获取url的参数
 * @param url
 */
export const appUrlParams = (url) => {
    url = url || window.location.href;
    var data = {};
    var pos = url.indexOf('?');
    if(pos != -1){
        var query = url.substring(pos);
        if(query){
            var ll = query.substring(1).split('&');
            var i = 0;
            for(; i < ll.length; i++){
                var arr = ll[i].split('=');
                if(arr.length == 2)
                    data[arr[0]] = decodeURI(arr[1]);
            }
        }
    }
    return data;
}
/**
 * 返回重组后的url
 * @param url 分享Url
 * @param r 上级ID
 * @param c 渠道
 * @returns {string}
 */
export const reformUrl = (url='',r='',c='') => {
    let customUrl = '';
    let pos = url.indexOf('?');
    customUrl = pos != -1 ?  url.substring(0,pos) : url;
    return `${customUrl}?r=${r}&c=${c}`
}

