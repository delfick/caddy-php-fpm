<?php

require_once $_SERVER['DOCUMENT_ROOT'] . '/app.lib.php';
$app = Core::App();

$app->logs->website->error("Hit the 404");
?>

<html>
<body>
<p>nup</p>
</body>
</html>
