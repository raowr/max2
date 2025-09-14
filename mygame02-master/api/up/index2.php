<?php
session_start();
//叫地主功能 默认id 为叫地主的人
if(isset($_GET['id']) && isset($_GET['id2']) && isset($_GET['id3']) )                   //判断是否存在page参数,获得页面值，否则取1
{
    
    $id = $_GET['id'];
    $id2 = $_GET['id2'];
    $id3 = $_GET['id3'];
    $next= $_GET['next'];

$roomid=$_SESSION["roomid"];
$arr_player=$_SESSION["player"];
$arr_prepare=$_SESSION["prepare"];

$player=json_decode($arr_player);
for($i=0;$i<3;$i++){
if($player[$i]==$id){
  $id_down=$i;
}
}



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

if($weight=="0"){
    $weight+=1;
    if($next==0){
      $new_prepare="[2,0,0]";
    }else if($next==1){
      $new_prepare="[0,2,0]";
    }else{
      $new_prepare="[0,0,2]";
    }
    $_SESSION["prepare"]=$new_prepare;
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
        
                $round=$row["round"]; 
              
           
        
        }
        }
        $round+=1;
        if($round==3){
          $round=0;
        }
        $sql = "UPDATE room SET round='$round',prepare='$new_prepare' where id = $roomid";
                if ($conn->query($sql) === TRUE) {
                  
                  
                } else {
                    
                    echo "Error: " . $sql . "<br>" . $conn->error;
                }
                exit();
}

if($weight=="1"){
  if($next==0){
    $new_prepare="[0,0,3]";
  }else if($next==1){
    $new_prepare="[3,0,0]";
  }else{
    $new_prepare="[0,3,0]";
  }
  $_SESSION["prepare"]=$new_prepare;
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
        $sql = "UPDATE room SET round='$id_down',prepare='$new_prepare' where id = $roomid";
        if ($conn->query($sql) === TRUE) {
          
          
        } else {
            
            echo "Error: " . $sql . "<br>" . $conn->error;
        }
}
        //底牌给地主
include_once(__DIR__."/../show/index.php");
$center=array_merge($num,json_decode($over));
sort($center);
print_r($center);
$json1=json_encode($center);
$sql = "UPDATE player SET card='$json1' WHERE id='$id'";
        if ($conn->query($sql) === TRUE) {
          $play1=true;
          
        } else {
            $play1=false;
            echo "Error: " . $sql . "<br>" . $conn->error;
        }

    

      }else{
        $next= $_GET['next'];
        $roomid=$_SESSION["roomid"];
        
        if($next==0){
          $new_prepare="[1,0,0]";
        }else if($next==1){
          $new_prepare="[0,1,0]";
        }else{
          $new_prepare="[0,0,1]";
        }
        $_SESSION["prepare"]=$new_prepare;
        include_once(__DIR__."/../dbconn.php");
$conn = new mysqli($servername, $username, $password, $dbname);
  // Check connection
  if ($conn->connect_error) {
    die("连接失败: " . $conn->connect_error);
                            } 
        $sql = "UPDATE room SET prepare='$new_prepare' where id = $roomid";
        if ($conn->query($sql) === TRUE) {
          
          
        } else {
            
            echo "Error: " . $sql . "<br>" . $conn->error;
        }



      }

?>