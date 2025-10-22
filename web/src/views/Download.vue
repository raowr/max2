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
                    <a
                        v-bind:href="downloadUrl"
                        class="ys-download-pc__btn"
                        ref="downloadbtn"
                        id="ys-download-pc__btn"
                    >
                        <img
                            src="@/assets/img/download/Android.png"
                            alt=""
                            srcset=""
                        />
                        <img
                            src="@/assets/img/download/Android_light.png"
                            alt=""
                        />
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
                <a
                    href="javascript:;"
                    class="download-m__btn"
                    @click="downloadApp"
                >
                    <img
                        src="@/assets/img/download/download_mb.png"
                        alt=""
                    />
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
import QRCode from "qrcodejs2";
import { androidOrIOS } from "@/utils/tools";
export default {
    name: "download",
    data() {
        return {
          apkData:null,
          androidOrIOS: androidOrIOS ,
          downloadUrl: process.env.VUE_APP_APK_URL,
          playSwitch: false
        };
    },
    created() {
      this.initGamePlaySwitch()
    },
    mounted() {
        //加载apk 版本文件   
        var ysbtn = this.$refs.downloadbtn;
        ysbtn.style.visibility='hidden'

        window.apkVersionHandler = function(filename){
            console.info(filename)
            ysbtn.href = ysbtn.href + filename;ysbtn.style.visibility="visible"
        }

        const script = document.createElement('script');
        script.type = 'text/javascript'
        script.src = process.env.VUE_APP_APK_URL + 'output-metadata.js?'+new Date().getTime();
        // ysbtn.appendChild(script);
        document.body.appendChild(script)
        //版本文件加载错误，使用默认
        script.onerror = function(){
            ysbtn.style.visibility='visible'
            ysbtn.href = ysbtn.href + "2max.apk"
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
            var qrCode = new QRCode(this.$refs.qrCodeUrl, {
                text: process.env.VUE_APP_DOWNLOAD_URL, // 待生成为二维码的内容
                width: 102,
                height: 102,
                correctLevel: QRCode.CorrectLevel.H,
            });
        },
        closeCode() {
            this.$refs.qrCodeUrl.innerHTML = "";
        },
        downloadApp() {
            var ysbtn = this.$refs.downloadbtn;
            if( ysbtn.style.visibility=='hidden'){
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
<style lang="less" scoped>
.download {
    position: absolute;
    width: 100%;
    height: 100%;
    background: url("@/assets/img/cg/cg.png") no-repeat;
    z-index: 10;
    .download-pc {
        position: absolute;
        height: 125px;
        bottom: 100px;
        left: 50%;
        transform: translateX(-50%);
        .download-pc__box--download {
            display: flex;
            justify-content: center;
            margin: 0 auto;
            // width: 540px;
            height: 100%;
            align-items: center;
            .download-pc__qr {
                position: relative;
                width: 102px;
                height: 102px;
                margin-right: 25px;
                // transition: all 0.5s ease;
                // &:hover {
                //     transform: scale(1.1);
                // }
            }
            .download-pc__group {
                display: flex;
                flex-direction: column;
                justify-content: center;
                align-items: center;
                height: 112px;
                padding: 3px 0;
                margin-right: 15px;
                .ys-download-pc__btn {
                    width: 160px;
                    height: 52px;
                    display: block;
                    cursor: pointer;
                    transition: all 0.5s ease;
                    position: relative;
                    &:hover {
                        img {
                            &:nth-child(1) {
                                opacity: 0;
                            }
                            &:nth-child(2) {
                                opacity: 1;
                            }
                        }
                    }
                    img {
                        width: 160px;
                        height: 52px;
                        position: absolute;
                        transition: all 0.3s ease;
                        &:nth-child(1) {
                            opacity: 1;
                        }
                        &:nth-child(2) {
                            opacity: 0;
                        }
                    }
                    span {
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
                }
            }
            .download-pc-btn-pc {
                width: 111px;
                height: 111px;
                display: block;
                cursor: pointer;
                img {
                    width: 111px;
                    height: 111px;
                    position: absolute;
                    transition: all 0.3s ease;

                    &:nth-child(1) {
                        opacity: 1;
                    }
                    &:nth-child(2) {
                        opacity: 0;
                    }
                }
                &:hover {
                    img {
                        &:nth-child(1) {
                            opacity: 0;
                        }
                        &:nth-child(2) {
                            opacity: 1;
                        }
                    }
                }
            }
        }
    }
    .download-m__dl-wrapper {
        display: none;
    }
    video {
        object-fit: fill;
        position: relative;
        z-index: -1;
    }
}
@media screen and (max-width: 575px) {
    video {
        display: none;
    }
    .download {
        background: url("@/assets/img/download/mb-bg.png") no-repeat;
        background-size: cover;
        background-position: center;
        .download-pc {
            display: none;
        }
        .download-m__dl-wrapper {
            display: block;
            position: absolute;
            bottom: 8%;
            width: 100%;
            .download-m-btn {
                margin-bottom: 30px;
                .download-m__btn {
                    width: 180px;
                    height: 100%;
                    cursor: pointer;
                    position: relative;
                    img {
                        width: 100%;
                        height: 100%;
                    }
                    span {
                        position: absolute;
                        text-align: center;
                        left: 50%;
                        top: 50%;
                        transform: translate(-50%, -50%);
                        font-size: 20px;
                        color: #fff;
                        line-height: 60px;
                    }
                }
            }
            .download-m__txts {
                p {
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
                    background: url("@/assets/img/download/textBg.png")
                        no-repeat;
                    background-size: cover;
                    text-align: center;
                    &.ys-download-m__txt2 {
                        span {
                            margin-right: 10px;
                        }
                    }
                }
            }
        }
    }
}
</style>
