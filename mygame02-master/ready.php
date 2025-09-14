<?php
session_start(); 
$roomid=$_SESSION["roomid"];
$id=$_SESSION["id"];
$arr_prepare=$_SESSION["prepare"];
$arr_player=$_SESSION["player"];

$player=json_decode($arr_player);
$prepare=json_decode($arr_prepare);
for($i=0;$i<3;$i++){
    if($player[$i]==$id){
        if( $prepare[$i]=="1"){
            $prepare[$i]="0";
        }else{
            $prepare[$i]="1";
        }
        break;
    }
}
$json1=json_encode($prepare);
include_once(__DIR__."/api/dbconn.php");
$conn = new mysqli($servername, $username, $password, $dbname);
 // Check connection
 if ($conn->connect_error) {
    die("连接失败: " . $conn->connect_error);
                            } 
$sql = "UPDATE room SET prepare='$json1' WHERE id='$roomid'";
if ($conn->query($sql) === TRUE) {
  $ok=true;
  $_SESSION["prepare"]= $json1;

} else {
    $ok=false;
    echo "Error: " . $sql . "<br>" . $conn->error;
}


?>