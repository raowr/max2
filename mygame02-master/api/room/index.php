<?php
session_start(); 
if(!isset($_SESSION["name"])){
    header("Location: login.php");
  }
if(isset($_GET['roomid'])){
    $id=$_SESSION['id'];
    $roomid=$_GET['roomid'];


}else{
    exit();
    
}
include_once(__DIR__."/../dbconn.php");
$conn = new mysqli($servername, $username, $password, $dbname);
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
    
            $arr_player=$row["player"]; 
            $_SESSION["player"]=$row["player"]; 
            $_SESSION["prepare"]=$row["prepare"]; 
       

    }
}else{
    echo json_encode(array("code"=>"404"));
    exit();
}


$player=json_decode($arr_player);

for($i=0;$i<3;$i++){
   
    
    if($player[$i]==$id){
        $_SESSION["roomid"]= $roomid;
       
        echo "<script>window.location.href='../../room.php'; </script>";
        exit();
    }
    if($player[$i]=="0"){
        $player[$i]=$id;
        break;
    }
}

$json1=json_encode($player);

$sql = "UPDATE room SET player='$json1' WHERE id='$roomid'";
if ($conn->query($sql) === TRUE) {
  $ok=true;
  $_SESSION["roomid"]= $roomid;
 echo "<script>window.location.href='../../room.php'; </script>";
} else {
    $ok=false;
    echo "Error: " . $sql . "<br>" . $conn->error;
}



?>