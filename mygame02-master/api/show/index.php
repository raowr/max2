<?php
//当存在get值id
if(isset($_GET["id"])){
    if(!empty($_GET["id"])){
            $id=$_GET["id"];
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
        
        
          $card=$row["card"];
        $over=$row["over_card"];
        
      
        
    }
    $num=Array();
    $cardnum=Array();
    $num=json_decode($card);
    sort($num);
    //print_r($num);
    for($i=0;$i<count($num);$i++){
      $cardnum[$i]=$arr[$num[$i]];
    }
    

    }
    $tmp=substr($over, 0, -1);
    $over=substr($tmp,1);
    echo json_encode(array("code"=>"200","card"=>implode(',',$num),"over"=> $over,"cardnum"=>json_encode($cardnum,JSON_UNESCAPED_UNICODE)),JSON_UNESCAPED_UNICODE); 
 

?>