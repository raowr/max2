<?php session_start(); ?>
<!DOCTYPE html>
<html lang="en-US">
<title>斗地主 | madmic</title>
    <meta charset="UTF-8">

<head>
<script type="application/javascript" src="js/jquery-1.8.3.min.js"></script>
<script type="application/javascript" src="js/yinghua.js"></script>

<style>
   input::-webkit-input-placeholder { /* Chrome/Opera/Safari */
    color: white;
  }
</style>
</head>
        <body style="background: url('img/ui/blackbg.png')no-repeat; background-size: 100% 100%;height: 100%;background-attachment:fixed;overflow:hidden;" >   
          <div style="background: url('img/ui/LOGO_CHST.png'); background-size: 100% 100%;height: 50%;width: 45%;position: absolute;top: 10%;left: 15%;">
          </div>
          <div>
            <form method="post" action="">
            <input type="text" name="ddz_user" style="background-color: rgb(89, 102, 127);width: 10%;height: 3%;position: absolute;top: 30%;left: 55%;color:white;" placeholder="账号/邮箱" >
            <input type="password" name="ddz_pwd" style="background-color: rgb(89, 102, 127);width: 10%;height: 3%;position: absolute;top: 35%;left: 55%;color:white;" placeholder="密码" > 
            <button type="submit" name="ddz" style="background: url('img/ui/login.png')no-repeat;background-size: 100% 100%;width: 10%;height: 5%;position: absolute;top: 40%;left: 55%;"></button>
        </form>
          </div>

          <?php
          if(isset($_POST['ddz_user']) && $_POST['ddz_pwd']){
            $user = $_POST['ddz_user'];
            $pwd = $_POST['ddz_pwd'];
            include_once(__DIR__."/api/dbconn.php");
            $conn = new mysqli($servername, $username, $password, $dbname);
         // Check connection
             if ($conn->connect_error) {
    die("连接失败: " . $conn->connect_error);
                            } 
     $sql = "SELECT * FROM player where password='$pwd' and mail='$user'";
    $result = $conn->query($sql);
     if ($res=mysqli_query($conn,$sql))
    {
        // 返回记录数
        $sums=mysqli_num_rows($res);
    }
    
    if ($result->num_rows > 0) {
    // 输出数据
    while($row = $result->fetch_assoc()) {
        
       
          $_SESSION["name"]=$row["name"];
          $_SESSION["duanwei"]=$row["duanwei"];
          $_SESSION["chengwei"]=$row["chengwei"];
          $_SESSION["lihui"]=$row["lihui"];
          $_SESSION["bglihui"]=$row["bglihui"];
          $_SESSION["touxiang"]=$row["touxiang"];
          $_SESSION["jinbi"]=$row["jinbi"];
          $_SESSION["hunyu"]=$row["hunyu"];
          $_SESSION["id"]=$row["id"];
          header("Location: index.php");
        
        
    }



          }else{
            echo "<script type=\"text/javascript\">alert('账号或密码错误！');</script>";
          }

        }

          ?>
        </body>
</html>