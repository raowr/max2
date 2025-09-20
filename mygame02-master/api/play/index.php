<?php
//当存在get值id
if(isset($_GET["roomid"]) && isset($_GET["id"]) && isset($_GET["prepare"]) && isset($_GET["num"]) && isset($_GET["weight"]) && isset($_GET["type"])){
    
            $roomid=$_GET["roomid"];
            $id=$_GET["id"];
            $prepare=$_GET["prepare"];
            $arr_card_out=$_GET["num"];
   $weight=$_GET["weight"];
$type=$_GET["type"];

}else{
  header('HTTP/1.1 404 Not Found');

    echo json_encode(array("status"=>"404"));
    exit();
}



if($prepare=="[3,0,0]"){
$new_prepare="[0,4,0]";

}else if($prepare=="[0,3,0]"){
  $new_prepare="[0,0,4]";
}else if($prepare=="[0,0,3]"){
  $new_prepare="[4,0,0]";
}else if($prepare=="[4,0,0]"){
  $new_prepare="[0,4,0]";
}else if($prepare=="[0,4,0]"){
  $new_prepare="[0,0,4]";
}else if($prepare=="[0,0,4]"){
  $new_prepare="[4,0,0]";
}


$card_out=json_decode($arr_card_out);





//判断是否需要比牌
include_once(__DIR__."/../dbconn.php");
$conn = new mysqli($servername, $username, $password, $dbname);
$sql = "SELECT * FROM room where id='$roomid'";
    $result = $conn->query($sql);
     if ($res=mysqli_query($conn,$sql))
    {
        // 返回记录数
        $sums=mysqli_num_rows($res);
    }
    
    if ($result->num_rows > 0) {
    // 输出数据
    while($row = $result->fetch_assoc()) {
        $discard=$row['discard'];
        $discard_player=$row['discard_player'];
        $discard_type=$row['discard_type'];
        $discard_weight=$row['discard_weight'];
        
    }
  }

if($discard_player!="" && $discard_player!=$id){
  //需要跟弃牌堆进行比牌操作
  if($discard_type==$type){
    //牌型匹配
    if($weight>$discard_weight){ //比大小
      $canplay=true;
    }else{
      $canplay=false;
      echo json_encode(array("status"=>"500")); 
      exit();
    }


  }else{
    if(($type==4 && $discard_type!=0 ) || $type==0 ){
      $canplay=true;
    }else{
      $canplay=false;
   
     echo json_encode(array("status"=>"500")); 
      exit();
    }

  }


}else{
  $canplay=true;
}







//获取手牌
$url = 'http://1.13.182.150/ddz/api/show/?id='.$id;

 $json = file_get_contents($url);
 $re = json_decode($json,true);
 $arr_card= explode(",",$re["card"]);

 
 
 
 
//出牌功能
 


if($canplay==true){
//print_r($num);
$arr_discard=[];
for($i=0;$i<count($card_out);$i++){
        //获取要出的牌
    
     array_push($arr_discard,$arr_card[$card_out[$i]]);
    unset($arr_card[$card_out[$i]]);
 
}

$arr_card=array_values($arr_card);
$json1=json_encode($arr_card);

if ($conn->connect_error) {
    die("连接失败: " . $conn->connect_error);
    }

    $sql = "UPDATE player SET card='$json1' WHERE id='$id'";
    if ($conn->query($sql) === TRUE) {
      $play1=true;
    } else {
        $play1=false;
        echo "Error: " . $sql . "<br>" . $conn->error;
    }
    $json2=json_encode($arr_discard);
   
    $sql = "UPDATE room SET discard='$json2',discard_player='$id',prepare='$new_prepare',discard_type='$type',discard_weight='$weight' WHERE id='$roomid'";
    if ($conn->query($sql) === TRUE) {
      if($play1==true){
          echo json_encode(array("status"=>"200")); 
      }
      
    } else {
        $play1=false;
        echo "Error: " . $sql . "<br>" . $conn->error;
    }
  }
?>