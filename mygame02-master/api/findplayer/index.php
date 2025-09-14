<?php
//当存在get值id
if(isset($_GET["id"])  ){
    
            $id=$_GET["id"];
            
    
            

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
     $sql = "SELECT * FROM player where id='$id'";
    $result = $conn->query($sql);
     if ($res=mysqli_query($conn,$sql))
    {
        // 返回记录数
        $sums=mysqli_num_rows($res);
    }
    
    if ($result->num_rows > 0) {
    // 输出数据
    while($row = $result->fetch_assoc()) {
        
        
          $name=$row["name"];
        $touxiang=$row["touxiang"];
          $card=$row["card"];
        $over_card=$row["over_card"];        
            $identity=$row["identity"];
   
    }
    echo json_encode(array("code"=>"200","name"=>$name,"touxiang"=> $touxiang,"card"=> $card,"overcard"=> $over_card,"identity"=> $identity),JSON_UNESCAPED_UNICODE); 
    }else{
         echo json_encode(array("code"=>"500")); 
    }
    
?>