<?php session_start();  ?>
<!DOCTYPE html>
<html lang="en-US">
<title>斗地主 | madmic</title>
    <meta charset="UTF-8">
    
<head>

        </head>
        <?php
$id=$_SESSION["id"];
$roomid=$_SESSION["roomid"];
$arr_player=$_SESSION["player"];

$player=json_decode($arr_player);
$player2="";

for($i=0;$i<3;$i++){
    if($player[$i]==$id){
        $player1=$player[$i];
       
    }
    else{
        if($player2!=""){
            $player3=$player[$i];
        }else{
            $player2=$player[$i]; 
        }
        
    }
}

include_once(__DIR__."/api/dbconn.php");
$conn = new mysqli($servername, $username, $password, $dbname);
  // Check connection
  if ($conn->connect_error) {
    die("连接失败: " . $conn->connect_error);
                            } 
     $sql = "SELECT * FROM player where id = $player1";
    $result = $conn->query($sql);
    
    
    if ($result->num_rows > 0) {
    // 输出数据
    while($row = $result->fetch_assoc()) {
    
            $over=$row["over_card"];
            
            $arr_player1_card=$row["card"];
            $player1_name=$row["name"];
            $player1_touxiang=$row["touxiang"];
            $player1_weight=$row["weight"];
            $player1_identity=$row["identity"];

    }
}
$sql = "SELECT * FROM player where id = $player2";
    $result = $conn->query($sql);
    
    
    if ($result->num_rows > 0) {
    // 输出数据
    while($row = $result->fetch_assoc()) {
    
           
            $arr_player2_card=$row["card"];
 $player2_name=$row["name"];
            $player2_touxiang=$row["touxiang"];
            $player2_weight=$row["weight"];
            $player2_identity=$row["identity"];

    }
}
$sql = "SELECT * FROM player where id = $player3";
    $result = $conn->query($sql);
    
    
    if ($result->num_rows > 0) {
    // 输出数据
    while($row = $result->fetch_assoc()) {
    
 
            $arr_player3_card=$row["card"];
            $player3_name=$row["name"];
            $player3_touxiang=$row["touxiang"];
            $player3_weight=$row["weight"];
            $player3_identity=$row["identity"];
    }
}
$sql = "SELECT * FROM room where id = $roomid";
    $result = $conn->query($sql);
    
    
    if ($result->num_rows > 0) {
    // 输出数据
    while($row = $result->fetch_assoc()) {
    
 
            $round=$row["round"];
    }
}
?>
        <body  style="padding:0;margin:0;overflow:hidden;height: 100%;">  
        <div style="position:absolute; width:100%; height:100%; z-index:-5; left:0; top:0;">
<img  src="img/bg/Table_Dif12324.png" height="100%" width="100%" style="left:0; top:0;">
</div> 
        <div style="height:200px;bottom:120px; width:1200px; position: absolute; left:35%; " id="anniu" >
       
       
        </div>
      
        
<div >
<img src="img/ui/chatlog.png" width="105px" style="position: absolute;bottom:11%;left:3%;">
<img src="<?php echo $player1_touxiang; ?>" width="100px" style="position: absolute;bottom:11.3%;left:3.1%;">
<img src="img/ui/frame_short.png" width="100px"height="35px" style="position: absolute;bottom:7%;left:3.1%;">
<p style="z-index:1;position: absolute; bottom:6%;left:4.1%;font-size:18px; color:white;"><?php echo $player1_name; ?></p> 
</div>
<div><!--底牌 -->
<img src="img/ui/tab_choose.png"  style="position: absolute;z-index:-1;">

<?php

$overcard=json_decode($over);

$player1_card=json_decode($arr_player1_card);
$player2_card=json_decode($arr_player2_card);
$player3_card=json_decode($arr_player3_card);
echo "<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;";
for($i=0;$i<count($overcard);$i++){
    
   
    echo "<img  src='img/".$overcard[$i].".png' width='60px' height='80px' >";


}
?>


</div>

<div style="height:200px;bottom:50%; position: absolute; left:25%;width:700px;" id="txtHint"  ><!--出牌库 -->




</div>
<div style="height:200px;bottom:50%; position: absolute; left:25%" ><!--赢了 -->
<h2 id="win"></h2>



</div>
<div style="height:200px;bottom:10px; position: absolute; left:10%" ><!--player1 牌 -->
<?php





if(count($player1_card)>17){
    $player_card_size='90px';
}else{
    $player_card_size='100px';
}

for($i=0;$i<count($player1_card);$i++){
    if($i<10){
        $inum="0".strval($i);
       
    }else{
        $inum=$i;
    }
   
    echo "<img onClick='chupai()'  src='img/".$player1_card[$i].".png' width='".$player_card_size."' height='120px'  id='".$inum."'>";


}
?>


</div>


<div style="height:200px;top:20%; position: absolute; left:3%;width:110px;"  ><!--player3 左牌-->
<?php
echo "<img src='img/ui/chatlog.png' width='105px' style='position: absolute;left:3%;z-index:-1;'>";
echo "<img src='".$player3_touxiang."' width='100px' style='position: absolute;z-index:1;left:5%;top:1%'>";
echo "<img src='img/ui/frame_short.png' width='100px'height='35px' style='position: absolute;top:55%;left:3.1%;'>";
echo "<p style='z-index:1;position: absolute; bottom:22%;left:28.1%;font-size:18px; color:white;'>".$player3_name."</p>";

echo "<img src='img/54.png' width='100px'style='position: absolute;z-index:1;left:5%;top:80%'>>";
?>
<h2  id="player3" style='z-index:1;position: absolute; top:135%;left:28.1%;font-size:18px; color:white;'>剩<?php echo count($player3_card);  ?>张</h2>
</div>

<div style="height:200px;top:20%; position: absolute; right:3%; float:right;width:110px;" ><!--player2 右牌-->
<?php
echo "<img src='img/ui/chatlog.png' width='105px' style='position: absolute;left:3%;z-index:-1;'>";
echo "<img src='".$player2_touxiang."' width='100px'style='position: absolute;z-index:1;left:5%;top:1%'>";
echo "<img src='img/ui/frame_short.png' width='100px'height='35px' style='position: absolute;top:55%;left:3.1%;'>";
echo "<p style='z-index:1;position: absolute; bottom:22%;left:28.1%;font-size:18px; color:white;'>".$player2_name."</p>";
echo "<img src='img/54.png' width='100px'style='position: absolute;z-index:1;left:5%;top:80%'>>";
?>
<h2  id="player2" style='z-index:1;position: absolute; top:135%;left:28.1%;font-size:18px; color:white;'>剩<?php echo count($player2_card);  ?>张</h2>
</div>
<?php 

if($id==$player[0] && $round=="0" && $player1_identity==""){
    echo "<img src='img/btn_jiaodizhu_up.png' style='position: absolute;top:67%;left:40%' onclick='jiaodizhu()'>
    <img src='img/btn_bujiao2.png' style='position: absolute;top:67%;left: 52%'  onclick='bujiao()'>";
}else if($id==$player[1] && $round=="1"  && ($player2_weight=="1" || $player3_weight=="1") && $player1_identity==""){
    echo "<img src='img/btn_qiangdizhu_up.png' style='position: absolute;top:67%;left:40%' onclick='jiaodizhu()'>
    <img src='img/btn_bujiao2.png' style='position: absolute;top:67%;left: 52%'  onclick='bujiao()'>";
} else if($id==$player[1] && $round=="1"  && ($player2_weight=="0" || $player3_weight=="0") && $player1_identity==""){
    echo "<img src='img/btn_jiaodizhu_up.png' style='position: absolute;top:67%;left:40%' onclick='jiaodizhu()'>
    <img src='img/btn_bujiao2.png' style='position: absolute;top:67%;left: 52%'  onclick='bujiao()'>";
} else if($id==$player[2] && $round=="2"  && ($player2_weight=="1" || $player3_weight=="1") && $player1_identity==""){
    echo "<img src='img/btn_qiangdizhu_up.png' style='position: absolute;top:67%;left:40%' onclick='jiaodizhu()'>
    <img src='img/btn_bujiao2.png' style='position: absolute;top:67%;left: 52%'  onclick='bujiao()'>";
} else if($id==$player[2] && $round=="2"  && ($player2_weight=="0" || $player3_weight=="0") && $player1_identity==""){
    echo "<img src='img/btn_jiaodizhu_up.png' style='position: absolute;top:67%;left:40%' onclick='jiaodizhu()'>
    <img src='img/btn_bujiao2.png' style='position: absolute;top:67%;left: 52%'  onclick='bujiao()'>";
}

if($id==$player[0] && $round=="0" && $player1_identity!=""){
    echo "<img src='img/btn_chupai.png' style='position: absolute;top:67%;left:40%' onclick='jiaodizhu()'>
    <img src='img/btn_bujiao2.png' style='position: absolute;top:67%;left: 52%'  onclick='bujiao()'>";

}


?>

<script>
  
    
    function jiaodizhu(){
        if (window.XMLHttpRequest)
{
    // IE7+, Firefox, Chrome, Opera, Safari 浏览器执行的代码
    xmlhttp=new XMLHttpRequest();
}
else
{    
    //IE6, IE5 浏览器执行的代码
    xmlhttp=new ActiveXObject("Microsoft.XMLHTTP");
}
xmlhttp.onreadystatechange=function()
{
    if (xmlhttp.readyState==4 && xmlhttp.status==200)
    {
        location.reload();
        
    }
}
xmlhttp.open("GET","api/up/index.php?id=<?php echo $player1; ?>&id2=<?php echo $player2; ?>&id3=<?php echo $player3; ?>",true);
xmlhttp.send();




        }
        function chupai(){
            
            
            var el = window.event.target;
            var chupainum=el.id;
         
                if(document.getElementById(chupainum).offsetHeight==120){
                    document.getElementById(chupainum).style.height="150px";     
                }else{
                    document.getElementById(chupainum).style.height="120px";     
                }

                if(cp.search(chupainum)==-1){
                    cp=cp+chupainum+",";
                    
                }else{
                    cp=cp.replace(chupainum+",","");
                }
                
        }
        function chupai2(){
           // alert(cp);
            cp=cp.slice(0,cp.length-1);
            window.location.href="chupai.php?cp=" + cp; 
            
        }
        //出牌库刷新显示
        setInterval(function(){
        if (window.XMLHttpRequest)
    {
        // IE7+, Firefox, Chrome, Opera, Safari 浏览器执行的代码
        xmlhttp=new XMLHttpRequest();
    }
    else
    {    
        //IE6, IE5 浏览器执行的代码
        xmlhttp=new ActiveXObject("Microsoft.XMLHTTP");
    }
    xmlhttp.onreadystatechange=function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)
        {
            document.getElementById("txtHint").innerHTML=xmlhttp.responseText;
        }
    }
    xmlhttp.open("GET","chupaiku.php",true);
    xmlhttp.send();


    },100);
    
    //玩家3牌数显示
    setInterval(function(){
        if (window.XMLHttpRequest)
    {
        // IE7+, Firefox, Chrome, Opera, Safari 浏览器执行的代码
        xmlhttp2=new XMLHttpRequest();
    }
    else
    {    
        //IE6, IE5 浏览器执行的代码
        xmlhttp2=new ActiveXObject("Microsoft.XMLHTTP");
    }
    xmlhttp2.onreadystatechange=function()
    {
        if (xmlhttp2.readyState==4 && xmlhttp2.status==200)
        {
            document.getElementById("player3").innerHTML=xmlhttp2.responseText;
        }
    }
    xmlhttp2.open("GET","paishu3.php",true);
    xmlhttp2.send();


    },100);
      //玩家2牌数显示
      setInterval(function(){
        if (window.XMLHttpRequest)
    {
        // IE7+, Firefox, Chrome, Opera, Safari 浏览器执行的代码
        xmlhttp3=new XMLHttpRequest();
    }
    else
    {    
        //IE6, IE5 浏览器执行的代码
        xmlhttp3=new ActiveXObject("Microsoft.XMLHTTP");
    }
    xmlhttp3.onreadystatechange=function()
    {
        if (xmlhttp3.readyState==4 && xmlhttp3.status==200)
        {
            document.getElementById("player2").innerHTML=xmlhttp3.responseText;
        }
    }
    xmlhttp3.open("GET","paishu2.php",true);
    xmlhttp3.send();


    },100);
    //判断赢了没有
      setInterval(function(){
        if (window.XMLHttpRequest)
    {
        // IE7+, Firefox, Chrome, Opera, Safari 浏览器执行的代码
        xmlhttp4=new XMLHttpRequest();
    }
    else
    {    
        //IE6, IE5 浏览器执行的代码
        xmlhttp4=new ActiveXObject("Microsoft.XMLHTTP");
    }
    xmlhttp4.onreadystatechange=function()
    {
        if (xmlhttp4.readyState==4 && xmlhttp4.status==200)
        {
            document.getElementById("win").innerHTML=xmlhttp4.responseText;
        }
    }
    xmlhttp4.open("GET","win.php",true);
    xmlhttp4.send();


    },100);
//刷新消息
setInterval(function(){
        if (window.XMLHttpRequest)
    {
        // IE7+, Firefox, Chrome, Opera, Safari 浏览器执行的代码
        xmlhttp5=new XMLHttpRequest();
    }
    else
    {    
        //IE6, IE5 浏览器执行的代码
        xmlhttp5=new ActiveXObject("Microsoft.XMLHTTP");
    }
    xmlhttp5.onreadystatechange=function()
    {
        if (xmlhttp5.readyState==4 && xmlhttp5.status==200)
        {
            document.getElementById("wobuyao").innerHTML=xmlhttp5.responseText;
        }
    }
    xmlhttp5.open("GET","wobuyao.php",true);
    xmlhttp5.send();


    },100);
    //刷新按钮
setInterval(function(){
        if (window.XMLHttpRequest)
    {
        // IE7+, Firefox, Chrome, Opera, Safari 浏览器执行的代码
        xmlhttp6=new XMLHttpRequest();
    }
    else
    {    
        //IE6, IE5 浏览器执行的代码
        xmlhttp6=new ActiveXObject("Microsoft.XMLHTTP");
    }
    xmlhttp6.onreadystatechange=function()
    {
        if (xmlhttp6.readyState==4 && xmlhttp6.status==200)
        {
            document.getElementById("anniu").innerHTML=xmlhttp6.responseText;
        }
    }
    xmlhttp6.open("GET","anniu.php",true);
    xmlhttp6.send();


    },1000);
    function buyao(){
           // alert(cp);
           document.getElementById("wobuyao").innerHTML="玩家1选择了不出！";
           setTimeout(function(){
            document.getElementById("wobuyao").innerHTML="";
	},3000);
            window.location.href="buchu.php?player=1"; 

        }
        function win(){
            alert('玩家1赢了！');
            window.location.href='gameover.php?winer=1';
        }
        function win2(){
            alert('玩家2赢了！');
            window.location.href='gameover.php?winer=2';
        }
        function win3(){
            alert('玩家3赢了！');
            window.location.href='gameover.php?winer=3';
        }
       
        
        setInterval(function(){
           
            if(document.getElementById("win").innerHTML=="玩家1赢了"){
                document.getElementById("win").innerHTML='-1';
                alert('玩家1赢了！');
                window.location.href='gameover.php?winer=1';
        }
        if(document.getElementById("win").innerHTML=="玩家2赢了"){
                document.getElementById("win").innerHTML='-1';
                alert('玩家2赢了！');
                window.location.href='gameover.php?winer=2';
        }
        if(document.getElementById("win").innerHTML=="玩家3赢了"){
                document.getElementById("win").innerHTML='-1';
                alert('玩家3赢了！');
                window.location.href='gameover.php?winer=3';
        }
        },1000);
        setInterval(function(){
           
           if(document.getElementById("huihenum").innerHTML=="消息：轮到玩家1出牌" || document.getElementById("huihenum").innerHTML=="消息：等待地主出牌" ){
             
               document.getElementById("anniu").style.display='block';
               

       }else{
        document.getElementById("anniu").style.display='none';
       
       }
       
       },100);
</script>



        </body>

</html>