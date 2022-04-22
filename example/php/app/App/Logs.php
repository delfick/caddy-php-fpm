<?php

namespace Website\App;

use Monolog\Formatter\JsonFormatter;
use OpenCensus\Trace\Propagator\HttpHeaderPropagator;

class GoogleCloudJsonFormatter extends JsonFormatter {
    public function format(array $record): string {
        return json_encode(
            $this->translateRecordForGoogleCloudLoggingFormat($record),
            JSON_UNESCAPED_SLASHES
        ) . "\n";
    }

    protected function translateRecordForGoogleCloudLoggingFormat(array $record) {
        try {
            /** @var \Datetime $dt */
            $dt = $record['datetime'];

            $record = array_merge($record['extra'], $record['context'], [
                'message' => $record['message'],
                'severity' => $record['level_name'],
                'time' => $dt->format('Y-m-d\TH:i:s.u\Z'),
                'channel' => $record['channel'],
            ]);
        } catch (\Exception $e) {
            file_put_contents('php://stderr', 'Failed to make json log format: ' . serialize($e) . '\n', FILE_APPEND);
        }

        return $record;
    }
}

class TracingProcessor {
    private string $googleProject;
    private HttpHeaderPropagator $tracing;

    public function __construct(HttpHeaderPropagator $tracing, string $googleProject) {
        $this->tracing = $tracing;
        $this->googleProject = $googleProject;
    }

    public function __invoke(array $record): array {
        try {
            $info = $this->tracing->extract($_SERVER);

            if ($info->fromHeader()) {
                $record['extra']['logging.googleapis.com/trace'] = 'projects/' . $this->googleProject . '/traces/' . $info->traceId();
                $record['extra']['logging.googleapis.com/spanId'] = $info->spanId();
                $record['extra']['logging.googleapis.com/trace_sampled'] = $info->enabled() || false;
            }
        } catch (\Exception $e) {
            file_put_contents('php://stderr', 'Failed to make tracing: ' . serialize($e) . '\n', FILE_APPEND);
        }

        return $record;
    }
}

class Logs {
    public \Monolog\Logger $website;
    public \Monolog\Logger $error_handler;

    public function __construct(HttpHeaderPropagator $tracing, string $googleProject) {
        $this->website = new \Monolog\Logger('website');
        $this->error_handler = new \Monolog\Logger('error_handler');

        $handler = new \Monolog\Handler\StreamHandler('php://stderr', \MonoLog\Logger::INFO);
        $handler->setFormatter(new GoogleCloudJsonFormatter());

        $this->website->pushHandler($handler);
        $this->error_handler->pushHandler($handler);

        $this->website->pushProcessor(new TracingProcessor($tracing, $googleProject));
        $this->error_handler->pushProcessor(new TracingProcessor($tracing, $googleProject));

        register_shutdown_function(function () {
            $this->error_handler();
        });
    }

    private function error_handler() {
        $error = error_get_last();

        if (!is_null($error)) {
            $msg = $error['message'];
            unset($error['message']);
            $this->error_handler->error($msg, $error);
        }
    }
}
