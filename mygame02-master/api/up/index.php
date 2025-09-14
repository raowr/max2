<?php
//当存在get值id
if(isset($_GET["id"]) && isset($_GET["id2"]) && isset($_GET["id3"]) && isset($_GET["roomid"]) && isset($_GET["prepare"]) ){
    if(!empty($_GET["id"]) && !empty($_GET["id2"]) && !empty($_GET["id3"]) && !empty($_GET["roomid"]) && !empty($_GET["prepare"])){
            $id=$_GET["id"];
            $id2=$_GET["id2"];
            $id3=$_GET["id3"];
            $roomid=$_GET["roomid"];
            $prepare=$_GET["prepare"];
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

//跳过回合
 $url = 'http://1.13.182.150/ddz/api/state/?roomid='.$roomid."&id=".$id;
 $json = file_get_contents($url);
 $re = json_decode($json,true);
 $ex=$re["ex"];
 if($ex==0){
     $next=1;
 }else if($ex==1){
     $next=2;
     
 }else{
     $next=0;
 }
  $re2=ltrim(rtrim($prepare,"]"),"[");
 $arr_prepare= explode(',',$re2);
//print_r($arr_prepare);
 

include_once(__DIR__."/../dbconn.php");
$conn = new mysqli($servername, $username, $password, $dbname);
  // Check connection
  if ($conn->connect_error) {
    die("连接失败: " . $conn->connect_error);
                            } 
     $sql = "SELECT * FROM player";
    $result = $conn->query($sql);
     if ($res=mysqli_query($conn,$sql))
    {
        // 返回记录数
        $sums=mysqli_num_rows($res);
    }
    
    if ($result->num_rows > 0) {
    // 输出数据
    while($row = $result->fetch_assoc()) {
        if($row["id"]===$id){
            $weight=$row["weight"]; 
          }

    }
}

if($weight=="0"){ //第一次叫地主
    $weight+=1;
    if ($conn->connect_error) {
        die("连接失败: " . $conn->connect_error);
        }
    
        $sql = "UPDATE player SET weight='$weight' WHERE id='$id'";
        if ($conn->query($sql) === TRUE) {
          $play1=true;
        } else {
            $play1=false;
            echo "Error: " . $sql . "<br>" . $conn->error;
        }
   $arr_prepare[$ex]=0;     
 $arr_prepare[$next]=2;
 $json=json_encode($arr_prepare);
$sql = "UPDATE room SET prepare='$json' WHERE id='$roomid'";
        if ($conn->query($sql) === TRUE) {
          $play1=true;
        } else {
            $play1=false;
            echo "Error: " . $sql . "<br>" . $conn->error;
        }

 echo json_encode(array("status"=>"200","weight"=>"1"));

}else if($weight=="1"){ //第二次当地主

    $sql = "UPDATE player SET identity='地主' WHERE id='$id'";
    if ($conn->query($sql) === TRUE) {
      $play1=true;
    } else {
        $play1=false;
        echo "Error: " . $sql . "<br>" . $conn->error;
    }
    $sql = "UPDATE player SET identity='农民' WHERE id='$id2' or id='$id3'";
    if ($conn->query($sql) === TRUE) {
      $play2=true;
      $play3=true;
    } else {
        $play2=false;
        $play3=false;
        echo "Error: " . $sql . "<br>" . $conn->error;
    }
    $sql = "UPDATE player SET weight='0' WHERE id='$id2' or id='$id3' or id='$id'";
        if ($conn->query($sql) === TRUE) {
          $play2=true;
          $play3=true;
        } else {
            $play2=false;
            $play3=false;
            echo "Error: " . $sql . "<br>" . $conn->error;
        }

//将底牌加到地主的手牌中
//获取手牌
 $url = 'http://1.13.182.150/ddz/api/show/?id='.$id;

 $json = file_get_contents($url);
 $re = json_decode($json,true);
 $arr_card= explode(',',$re["card"]);
 $over_card=explode(',',$re["over"]);


$new_card=array_merge($arr_card,$over_card);
sort($new_card);
$json1=json_encode($new_card);
$sql = "UPDATE player SET card='$json1' WHERE id='$id'";
if ($conn->query($sql) === TRUE) {
    $play1=true;
} else {
    $play1=false;
    echo "Error: " . $sql . "<br>" . $conn->error;
}
if($ex==0){
    $tmp="[3,0,0]";
}else if($ex==1){
    $tmp="[0,3,0]";
}else{
   $tmp="[0,0,3]"; 
}

$sql = "UPDATE room SET prepare='$tmp' WHERE id='$roomid'";
if ($conn->query($sql) === TRUE) {
    $play1=true;
} else {
    $play1=false;
    echo "Error: " . $sql . "<br>" . $conn->error;
}
echo json_encode(array("status"=>"200","weight"=>"2"));
}



?>