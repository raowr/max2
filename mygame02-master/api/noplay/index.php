<?php
//当存在get值id
if(isset($_GET["roomid"]) && isset($_GET["prepare"]) ){
    
            $roomid=$_GET["roomid"];
            
            $prepare=$_GET["prepare"];
            

}else{
  header('HTTP/1.1 404 Not Found');

    echo json_encode(array("status"=>"404"));
    exit();
}

include_once(__DIR__."/../dbconn.php");
$conn = new mysqli($servername, $username, $password, $dbname);
  // Check connection
  if ($conn->connect_error) {
    die("连接失败: " . $conn->connect_error);
                            } 
     $sql = "SELECT * FROM room where id='$roomid'";
    $result = $conn->query($sql);
     
    
    if ($result->num_rows > 0) {
    // 输出数据
    while($row = $result->fetch_assoc()) {
          $noplaynum=$row["noplaynum"];
        
    }
     }
     if($noplaynum=="0" || $noplaynum==0){
         $noplaynum=1;
     }else if($noplaynum=="1" || $noplaynum==1){
         $noplaynum=2;
     }else{
        $noplaynum=1; 
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

    include_once(__DIR__."/../dbconn.php");
$conn = new mysqli($servername, $username, $password, $dbname);
  // Check connection
  if ($conn->connect_error) {
    die("连接失败: " . $conn->connect_error);
                            } 
    $sql = "UPDATE room SET prepare='$new_prepare',noplaynum='$noplaynum' WHERE id='$roomid'";
    if ($conn->query($sql) === TRUE) {
      echo json_encode(array("status"=>"200")); 
      
    } else {
        $play1=false;
        echo "Error: " . $sql . "<br>" . $conn->error;
    }
?>