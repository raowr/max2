<?php
session_start();

$arr_player=$_GET["player"];//当局游戏的三个人的玩家id
$id=$_GET["id"];//当前登录的玩家id
$prepare=$_GET["prepare"];//当局游戏的三个人的当前状态
$roomid=$_GET["roomid"];
include_once(__DIR__."/api/dbconn.php");
$conn = new mysqli($servername, $username, $password, $dbname);
//检查用户的牌是否打完
$url = 'http://1.13.182.150/ddz/api/findplayer/?id='.$id;
 $json = file_get_contents($url);
 $re = json_decode($json,true);

 $player1_card= $re["card"];

$url = 'http://1.13.182.150/ddz/api/findplayer/next/?id='.$id."&player=".$arr_player;

 $json = file_get_contents($url);
 $re = json_decode($json,true);
  $player2_card= $re["card"];
 $url = 'http://1.13.182.150/ddz/api/findplayer/up/?id='.$id."&player=".$arr_player;
 $json = file_get_contents($url);
 $re = json_decode($json,true);
 $player3_card= $re["card"];

if($player1_card=="[]" || $player1_card=="" || $player2_card=="" || $player3_card=="" || $player2_card=="[]" || $player3_card=="[]"){
     // Check connection
  if ($conn->connect_error) {
    die("连接失败: " . $conn->connect_error);
                            } 
     $sql = "UPDATE room SET begin='off',prepare='[1,0,0]' WHERE id='$roomid'";
    if ($conn->query($sql) === TRUE) {
       echo "-2";
        exit();
      
    } else {
      
        echo "Error: " . $sql . "<br>" . $conn->error;
    }
    
    
   
}


  // Check connection
  if ($conn->connect_error) {
    die("连接失败: " . $conn->connect_error);
                            } 
     $sql = "SELECT * FROM room where id = $roomid";
    $result = $conn->query($sql);
     if ($res=mysqli_query($conn,$sql))
    {
        // 返回记录数
        $sums=mysqli_num_rows($res);
    }
    
    if ($result->num_rows > 0) {
    // 输出数据
    while($row = $result->fetch_assoc()) {
            $data_prepare=$row["prepare"];
    }
}else{
    echo json_encode(array("code"=>"404"));
    exit();
}

if($data_prepare==$prepare){

    echo "1";
}else{

    echo "-1";
    $_SESSION["prepare"]= $data_prepare;
}

?>