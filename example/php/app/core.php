<?php

require_once __DIR__ . '/../vendor/autoload.php';

use Website\App\App;

class Core {
    private static ?App $instance = null;

    public static function App(): App {
        if (is_null(self::$instance)) {
            self::$instance = new App();
        }

        return self::$instance;
    }
}
