<?php

namespace Website\App;

use OpenCensus\Trace\Propagator\HttpHeaderPropagator;

class App {
    public Logs $logs;

    private string $googleProject;
    private HttpHeaderPropagator $tracing;

    public function __construct() {
        $this->googleProject = getenv("GOOGlE_PROJECT");

        $this->tracing = new HttpHeaderPropagator();
        $this->logs = new Logs($this->tracing, $this->googleProject);
    }
}
