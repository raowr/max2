<?php
session_start(); 
$arr_player=$_SESSION["player"];
$arr_prepare=$_SESSION["prepare"];
$roomid=$_SESSION["roomid"];
include_once(__DIR__."/api/dbconn.php");
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
    
            $data_player=$row["player"];
            $data_prepare=$row["prepare"];
            $begin=$row["begin"];

    }
}else{
    echo json_encode(array("code"=>"404"));
    exit();
}

if($begin=="ok"){
    $_SESSION["prepare"]=["1","0","0"];
    echo "2";
    exit();
}



if($data_player==$arr_player && $arr_prepare==$data_prepare){
    //echo "匹配";
    echo "1";
}else{
    //echo "不匹配！";
    echo "-1";
    $_SESSION["player"]= $data_player;
    $_SESSION["prepare"]= $data_prepare;

}

?>