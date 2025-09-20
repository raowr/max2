<?php
session_start(); 
$id=$_SESSION["id"];
include_once(__DIR__."/api/dbconn.php");
$roomid=time();
$player="[".$id.",0,0]";
$_SESSION["roomid"]=$roomid;
$conn = new mysqli($servername, $username, $password, $dbname);
// Check connection
if ($conn->connect_error) {
die("连接失败: " . $conn->connect_error);
}        
$sql = "INSERT INTO room (id,player)
VALUES ('$roomid','$player')";
if ($conn->query($sql) === TRUE) {
     $add=true;
     echo "<script> window.location.href='room.php?id=".$roomid."';</script>";
} else {
     //数据库添加错误
     echo json_encode(array("code"=>"500"));
     exit();
} 



?>