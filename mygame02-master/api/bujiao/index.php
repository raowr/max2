<?php
//当存在get值id
if(isset($_GET["roomid"]) && isset($_GET["prepare"])  && isset($_GET["id"])){
    
            $roomid=$_GET["roomid"];
            
            $prepare=$_GET["prepare"];
            $id=$_GET["id"];

            

}else{
  header('HTTP/1.1 404 Not Found');

    echo json_encode(array("status"=>"404"));
    exit();
}
$url = 'http://1.13.182.150/ddz/api/state/?roomid='.$roomid."&id=".$id;
 $json = file_get_contents($url);
 $re = json_decode($json,true);
 $state=$re["state"];
 $ex=$re["ex"];
 $re2=ltrim(rtrim($prepare,"]"),"[");
 $arr_prepare= explode(',',$re2);
 if($ex==0){
     $next=1;
     $up=2;
 }else if($ex==1){
     $next=2;
     $up=0;
 }else{
     $next=0;
     $up=1;
 }
 if($state==2 || $state=="2"){
     $arr_prepare[$ex]=0;
     $arr_prepare[$next]=2;
      $arr_prepare[$up]=0;
 }else{
      $arr_prepare[$ex]=0;
     $arr_prepare[$next]=1;
      $arr_prepare[$up]=0;
 }
  //print_r($arr_prepare);
  $new_prepare=json_encode($arr_prepare);
 

    include_once(__DIR__."/../dbconn.php");
$conn = new mysqli($servername, $username, $password, $dbname);
  // Check connection
  if ($conn->connect_error) {
    die("连接失败: " . $conn->connect_error);
                            } 
    $sql = "UPDATE room SET prepare='$new_prepare' WHERE id='$roomid'";
    if ($conn->query($sql) === TRUE) {
      echo json_encode(array("status"=>"200")); 
      
    } else {
        $play1=false;
        echo "Error: " . $sql . "<br>" . $conn->error;
    }
?>