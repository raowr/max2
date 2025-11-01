<template>
    <div class="download">
        <div v-show="this.playSwitch" class="download-pc">
            <div class="download-pc__box--download">
                <div class="download-pc__qr" ref="qrCodeUrl"></div>
                <div class="download-pc__group">
                    <!-- <a href="" class="ys-download-pc__btn">
                        <img
                            src="@/assets/img/download/apple.png"
                            alt=""
                            srcset=""
                        />
                        <img
                            src="@/assets/img/download/apple_light.png"
                            alt=""
                            srcset=""
                        />
                    </a> -->
                    <a v-bind:href="downloadUrl" class="ys-download-pc__btn" ref="downloadbtn" id="ys-download-pc__btn"
                        v-show="showDownloadButton">
                        <img src="@/assets/img/download/Android.png" alt="" srcset="" />
                        <img src="@/assets/img/download/Android_light.png" alt="" />
                        <span> Android </span>
                    </a>
                </div>
                <!-- <div class="download-pc-btn-pc">
                    <img src="@/assets/img/download/PC.png" alt="" />
                    <img src="@/assets/img/download/pc_light.png" alt="" />
                </div> -->
            </div>
        </div>
        <div v-show="this.playSwitch" class="download-m__dl-wrapper">
            <div class="download-m-btn flex-r-c-c">
                <a v-bind:href="downloadUrl" class="download-m__btn" ref="downloadbtn">
                    <img src="@/assets/img/download/download_mb.png" alt="" />
                    <span>Download</span>
                </a>
            </div>
            <div class="download-m__txts">
                <!--                <p><span>开发者名称：广州多米科技有限公司</span></p>-->
                <!--                <p class="ys-download-m__txt2">-->
                <!--                    <span>当前版本：1.0.0</span><span>更新时间：2023.6.1</span>-->
                <!--                </p>-->
            </div>
        </div>
        <!--        <video-->
        <!--            src="../../../assets/video/download.mp4"-->
        <!--            autoplay-->
        <!--            loop-->
        <!--            muted-->
        <!--            width="100%"-->
        <!--            height="100%"-->
        <!--        ></video>-->
        <!-- <script type="text/javascript" src="https://res.igplayplus.com/client/uat/igapp/apk/output-metadata.js"></script> -->
        <!-- <script>
            var localHandler = function(data){
                var ysbtn = this.$refs.downloadbtn;
                ysbtn.href = ysbtn.href + this.apkData.elements[0].outputFile
            }
        </script> -->
    </div>
</template>

<!-- <script type="text/javascript">
var localHandler = function(data){
    var ysbtn = this.$refs.downloadbtn;
    ysbtn.href = ysbtn.href + this.apkData.elements[0].outputFile
}
</script> -->


<script>
import QRCode from 'qrcode';
import { androidOrIOS } from "@/utils/tools";
export default {
    name: "download",
    data() {
        return {
            apkData: null,
            androidOrIOS: androidOrIOS,
            downloadUrl: import.meta.env.VITE_APP_APK_URL,
            playSwitch: false,
            showDownloadButton: false // 添加这个属性控制按钮可见性
        };
    },
    created() {
        this.initGamePlaySwitch()
    },
    mounted() {
        // 保存this引用，以便在回调函数中使用
        const vm = this;

        window.apkVersionHandler = function (filename) {
            console.info(filename)
            // 关键修复：更新Vue的响应式数据，而不是直接操作DOM
            vm.downloadUrl = vm.downloadUrl + filename;
            vm.showDownloadButton = true; // 添加一个响应式属性来控制可见性
            console.log(vm.downloadUrl)
        }

        const script = document.createElement('script');
        script.type = 'text/javascript'
        script.src = import.meta.env.VITE_APP_APK_URL + 'output-metadata.js?' + new Date().getTime();
        document.body.appendChild(script)

        //版本文件加载错误，使用默认
        script.onerror = function () {
            vm.downloadUrl = vm.downloadUrl + "2max.apk";
            vm.showDownloadButton = true;
        }
        this.createQrCode();
    },
    methods: {
        async initGamePlaySwitch() {
            await this.$api.getGamePlaySwitch({}).then((res) => {
                if (res.code == 0) {
                    console.log("initGamePlaySwitch data", res.data);
                    this.playSwitch = res.data.switch != "关闭" ? true : false;
                } else {
                    this.$Message.error(res.msg);
                }
            });
        },
        createQrCode() {
            console.log("--------------")
            console.log(import.meta.env.VITE_APP_DOWNLOAD_URL)

            // 添加数据验证
            const url = import.meta.env.VITE_APP_DOWNLOAD_URL;
            if (!url || typeof url !== 'string') {
                console.error('无效的URL数据:', url);
                this.$refs.qrCodeUrl.innerHTML = '<div style="width:100%;height:100%;display:flex;align-items:center;justify-content:center;">二维码生成失败</div>';
                return;
            }

            try {
                // 清除之前可能存在的内容
                this.$refs.qrCodeUrl.innerHTML = '';

                // 创建一个canvas元素
                const canvas = document.createElement('canvas');
                canvas.width = 102; // 与容器宽度匹配
                canvas.height = 102; // 与容器高度匹配
                this.$refs.qrCodeUrl.appendChild(canvas);

                // 使用配置选项尝试生成二维码
                QRCode.toCanvas(canvas, url, {
                    errorCorrectionLevel: 'M',
                    margin: 1,
                    width: 102,
                    color: {
                        dark: '#000000', // 二维码颜色
                        light: '#ffffff'  // 背景色
                    }
                }, (error) => {
                    if (error) {
                        console.error('二维码生成错误:', error);
                        this.$refs.qrCodeUrl.innerHTML = '<div style="width:100%;height:100%;display:flex;align-items:center;justify-content:center;">二维码生成失败</div>';
                    } else {
                        console.log('二维码生成成功');
                    }
                });
            } catch (e) {
                console.error('二维码生成异常:', e);
                this.$refs.qrCodeUrl.innerHTML = '<div style="width:100%;height:100%;display:flex;align-items:center;justify-content:center;">二维码生成失败</div>';
            }
        },
        async initGamePlaySwitch() {
            try {
                // 检查 $api 是否存在
                if (!this.$api || typeof this.$api.getGamePlaySwitch !== 'function') {
                    console.warn('$api 未正确配置，使用默认值');
                    // 设置默认值，确保页面能正常显示
                    this.playSwitch = true;
                    return;
                }

                await this.$api.getGamePlaySwitch({}).then((res) => {
                    if (res.code == 0) {
                        console.log("initGamePlaySwitch data", res.data);
                        this.playSwitch = res.data.switch != "关闭" ? true : false;
                    } else {
                        // 如果API调用失败，使用默认值
                        console.warn('获取游戏开关状态失败，使用默认值', res.msg);
                        this.playSwitch = true;
                    }
                });
            } catch (error) {
                console.error('获取游戏开关状态时发生错误:', error);
                // 发生异常时也使用默认值
                this.playSwitch = true;
            }
        },
        closeCode() {
            this.$refs.qrCodeUrl.innerHTML = "";
        },
        downloadApp() {
            var ysbtn = this.$refs.downloadbtn;
            if (ysbtn.style.visibility == 'hidden') {
                return
            }
            if (this.androidOrIOS() == 1) {
                // window.location.href = this.ppro;
                window.location.href = ysbtn.href;
            } else if (this.androidOrIOS() == 2) {
                window.location.href = "";
            }
        },
    },
};
</script>
<style scoped>
.download {
    position: absolute;
    width: 100%;
    height: 100%;
    background: url("@/assets/img/cg/cg.png") no-repeat;
    z-index: 10;
}

.download .download-pc {
    position: absolute;
    height: 125px;
    bottom: 100px;
    left: 50%;
    transform: translateX(-50%);
}

.download .download-pc .download-pc__box--download {
    display: flex;
    justify-content: center;
    margin: 0 auto;
    height: 100%;
    align-items: center;
}

.download .download-pc .download-pc__box--download .download-pc__qr {
    position: relative;
    width: 102px;
    height: 102px;
    margin-right: 25px;
}

.download .download-pc .download-pc__box--download .download-pc__group {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 112px;
    padding: 3px 0;
    margin-right: 15px;
}

.download .download-pc .download-pc__box--download .download-pc__group .ys-download-pc__btn {
    width: 160px;
    height: 52px;
    display: block;
    cursor: pointer;
    transition: all 0.5s ease;
    position: relative;
}

.download .download-pc .download-pc__box--download .download-pc__group .ys-download-pc__btn:hover img:nth-child(1) {
    opacity: 0;
}

.download .download-pc .download-pc__box--download .download-pc__group .ys-download-pc__btn:hover img:nth-child(2) {
    opacity: 1;
}

.download .download-pc .download-pc__box--download .download-pc__group .ys-download-pc__btn img {
    width: 160px;
    height: 52px;
    position: absolute;
    transition: all 0.3s ease;
}

.download .download-pc .download-pc__box--download .download-pc__group .ys-download-pc__btn img:nth-child(1) {
    opacity: 1;
}

.download .download-pc .download-pc__box--download .download-pc__group .ys-download-pc__btn img:nth-child(2) {
    opacity: 0;
}

.download .download-pc .download-pc__box--download .download-pc__group .ys-download-pc__btn span {
    width: 160px;
    height: 52px;
    top: 0;
    left: 0;
    line-height: 52px;
    text-align: center;
    font-size: 20px;
    color: #fff;
    position: absolute;
    text-indent: 20%;
}

.download .download-pc .download-pc__box--download .download-pc-btn-pc {
    width: 111px;
    height: 111px;
    display: block;
    cursor: pointer;
}

.download .download-pc .download-pc__box--download .download-pc-btn-pc img {
    width: 111px;
    height: 111px;
    position: absolute;
    transition: all 0.3s ease;
}

.download .download-pc .download-pc__box--download .download-pc-btn-pc img:nth-child(1) {
    opacity: 1;
}

.download .download-pc .download-pc__box--download .download-pc-btn-pc img:nth-child(2) {
    opacity: 0;
}

.download .download-pc .download-pc__box--download .download-pc-btn-pc:hover img:nth-child(1) {
    opacity: 0;
}

.download .download-pc .download-pc__box--download .download-pc-btn-pc:hover img:nth-child(2) {
    opacity: 1;
}

.download .download-m__dl-wrapper {
    display: none;
}

.download video {
    object-fit: fill;
    position: relative;
    z-index: -1;
}

@media screen and (max-width: 575px) {
    video {
        display: none;
    }

    .download {
        background: url("@/assets/img/download/mb-bg.png") no-repeat;
        background-size: cover;
        background-position: center;
    }

    .download .download-pc {
        display: none;
    }

    .download .download-m__dl-wrapper {
        display: block;
        position: absolute;
        bottom: 8%;
        width: 100%;
    }

    .download .download-m__dl-wrapper .download-m-btn {
        margin-bottom: 30px;
    }

    .download .download-m__dl-wrapper .download-m-btn .download-m__btn {
        width: 180px;
        height: 100%;
        cursor: pointer;
        position: relative;
    }

    .download .download-m__dl-wrapper .download-m-btn .download-m__btn img {
        width: 100%;
        height: 100%;
    }

    .download .download-m__dl-wrapper .download-m-btn .download-m__btn span {
        position: absolute;
        text-align: center;
        left: 50%;
        top: -300%;
        transform: translate(-50%, -50%);
        font-size: 40px;
        color: #fff;
        line-height: 60px;
    }

    .download .download-m__dl-wrapper .download-m__txts p {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 85%;
        height: 26px;
        margin: 0 auto;
        font-size: 16px;
        color: #fff;
        text-align: center;
        margin-bottom: 8px;
        background: url("@/assets/img/download/textBg.png") no-repeat;
        background-size: cover;
        text-align: center;
    }

    .download .download-m__dl-wrapper .download-m__txts p.ys-download-m__txt2 span {
        margin-right: 10px;
    }
}
</style>
