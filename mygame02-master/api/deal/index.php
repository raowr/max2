<?php
session_start();
//当存在get值id
if(isset($_GET["id1"]) && isset($_GET["id2"]) && isset($_GET["id3"])  && isset($_GET["roomid"])){
    if(!empty($_GET["id1"]) && !empty($_GET["id2"]) && !empty($_GET["id3"]) && !empty($_GET["roomid"])){
            $id1=$_GET["id1"];
            $id2=$_GET["id2"];
            $id3=$_GET["id3"];
            $roomid=$_GET["roomid"];
    }else{
         header('HTTP/1.1 404 Not Found');

    echo json_encode(array("status"=>"404"));
    exit();
    }

}else{
  header('HTTP/1.1 404 Not Found');

    echo json_encode(array("status"=>"404"));
    exit();
}
$_SESSION["player"]="[".$id1.",".$id2.",".$id3."]";
//斗地主牌库
 $arr = array(
    0=>'大王',1=>'小王',
    2=>'方片2',3=>'梅花2',4=>'红心2',5=>'黑桃2',
    6=>'方片A',7=>'梅花A',8=>'红心A',9=>'黑桃A',
    10=>'方片K',11=>'梅花K',12=>'红心K',13=>'黑桃K',
    14=>'方片Q',15=>'梅花Q',16=>'红心Q',17=>'黑桃Q',
    18=>'方片J',19=>'梅花J',20=>'红心J',21=>'黑桃J',
    22=>'方片10',23=>'梅花10',24=>'红心10',25=>'黑桃10',
    26=>'方片9',27=>'梅花9',28=>'红心9',29=>'黑桃9',
    30=>'方片8',31=>'梅花8',32=>'红心8',33=>'黑桃8',
    34=>'方片7',35=>'梅花7',36=>'红心7',37=>'黑桃7',
    38=>'方片6',39=>'梅花6',40=>'红心6',41=>'黑桃6',
    42=>'方片5',43=>'梅花5',44=>'红心5',45=>'黑桃5',
    46=>'方片4',47=>'梅花4',48=>'红心4',49=>'黑桃4',
    50=>'方片3',51=>'梅花3',52=>'红心3',53=>'黑桃3',
  );
  $num=Array();
  for($i=0;$i<54;$i++){
    $num[$i]=$i;
  }
  //print_r($num);
  
  //打乱牌库
  shuffle($num);
  //创建手牌数组
  $player1_card=array();
  $player2_card=array();
  $player3_card=array();
  //创建底牌数组
  $over_card=array();
  //玩家1的手牌
  for($i=0;$i<17;$i++){
    $player1_card[$i]=$num[$i];

  }
  //玩家2的手牌
  for($i=17;$i<34;$i++){
    $player2_card[$i]=$num[$i];

  }
  //玩家3的手牌
  for($i=34;$i<51;$i++){
    $player3_card[$i]=$num[$i];

  }
   //底牌
   for($i=51;$i<54;$i++){
    $over_card[$i]=$num[$i];

  }

  $player2_card=array_values($player2_card);
  $player3_card=array_values($player3_card);
  $over_card=array_values($over_card);
  sort($player1_card);
  sort($player2_card);
     sort($player3_card);
  $json1=json_encode($player1_card);
  $json2=json_encode($player2_card);
  $json3=json_encode($player3_card);
  $json4=json_encode($over_card);
  include_once(__DIR__."/../dbconn.php");

  $conn = new mysqli($servername, $username, $password, $dbname);
// Check connection
if ($conn->connect_error) {
   
    die("连接失败: " . $conn->connect_error);
    }

    $sql = "UPDATE player SET card='$json1' WHERE id='$id1'";

    if ($conn->query($sql) === TRUE) {
      $play1=true;
    } else {
        $play1=false;
        
        echo "Error: " . $sql . "<br>" . $conn->error;
    }
    $sql = "UPDATE player SET card='$json2' WHERE id='$id2'";
    if ($conn->query($sql) === TRUE) {
      $play2=true;
    } else {
        $play2=false;
        echo "Error: " . $sql . "<br>" . $conn->error;
    }
    $sql = "UPDATE player SET card='$json3' WHERE id='$id3'";
    if ($conn->query($sql) === TRUE) {
      $play3=true;
    } else {
        $play3=false;
        echo "Error: " . $sql . "<br>" . $conn->error;
    }
    $sql = "UPDATE player SET over_card='$json4',identity='',weight='0' WHERE id='$id1' or id='$id2' or id='$id3'";
    if ($conn->query($sql) === TRUE) {
      $over=true;
    } else {
        $over=false;
        echo "Error: " . $sql . "<br>" . $conn->error;
    }
$sql = "UPDATE room SET prepare='[1,0,0]',discard='',discard_player='',discard_type='',discard_weight='',begin='ok' WHERE id='$roomid'";
    if ($conn->query($sql) === TRUE) {
      $over=true;
    } else {
        $over=false;
        echo "Error: " . $sql . "<br>" . $conn->error;
    }
echo json_encode(array("status"=>"200")); 
?>