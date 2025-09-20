<?php
//当存在get值id
if(isset($_GET["id"])  && isset($_GET["roomid"])){
    
            $id=$_GET["id"];
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
         $noplaynum=$row["noplaynum"];
         $discard_player=$row["discard_player"];
       
        
    }
     }else{
         echo json_encode(array("code"=>"500")); 
    }
$arr_player=json_decode($player);
$arr_prepare=json_decode($prepare);
for($i=0;$i<3;$i++){
    if($arr_player[$i]==$id){
        
        if(($id==$discard_player) && ($noplaynum==2 || $noplaynum=="2") && ($arr_prepare[$i]!=0 || $arr_prepare[$i]!="0")){
            
         echo json_encode(array("code"=>"200","state"=>3,"ex"=>$i));
             exit();
            
            
            
            
             
        }else{
            echo json_encode(array("code"=>"200","state"=>$arr_prepare[$i],"ex"=>$i));
             exit();
        }
        
       
    }
}

?>