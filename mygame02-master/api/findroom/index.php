<?php
//当存在get值id
if(isset($_GET["roomid"])  ){
    
            $roomid=$_GET["roomid"];
            
    
            

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
     if ($res=mysqli_query($conn,$sql))
    {
        // 返回记录数
        $sums=mysqli_num_rows($res);
    }
    
    if ($result->num_rows > 0) {
    // 输出数据
    while($row = $result->fetch_assoc()) {
        
        
          $player=$row["player"];
        $prepare=$row["prepare"];
          $begin=$row["begin"];
        $discard=$row["discard"];        
            $discard_player=$row["discard_player"];
        $discard_type=$row["discard_type"];         
        $discard_weight=$row["discard_weight"];    
    }
    echo json_encode(array("code"=>"200","player"=>$player,"begin"=> $begin,"prepare"=>$prepare,"discard"=> $discard,"discard_player"=> $discard_player,"discard_type"=> $discard_type,"discard_weight"=> $discard_weight)); 
    }else{
         echo json_encode(array("code"=>"500")); 
    }
  
    
    
?>