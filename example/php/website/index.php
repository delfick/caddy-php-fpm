<?php

require_once $_SERVER['DOCUMENT_ROOT'] . '/app.lib.php';
$app = Core::App();

$t = "";
if (isset($_GET["t"])) {
    $t = $_GET["t"];
}

$app->logs->website->info("Hit the index", ["t" => $t]);
?>

<html>
<body>
<p>hi</p>
</body>
</html>
