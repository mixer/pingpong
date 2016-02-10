/**
 * Number of pings to make to each ingest server=.
 * @type {Number}
 */
var maxTests = 5;

/**
 * Creates a new PingPong tester against the list of ingests. The ingests list
 * is expected to be in the same format as the "nodes" list when recursively
 * getting the `ingests` directory from etcd.
 * @param {[]Object} ingests
 */
function PingPong(ingests) {
    this._ingests = ingests;
    this._tests = [];
}

/**
 * Tests latency to all ingest servers. The callback will be called multiple
 * times as we get updates from the servers. The first argument will be an
 * array of ingest records with an additional "latency" attribute.
 * @param {Function} callback
 */
PingPong.prototype.test = function (callback) {
    this._tests = this._ingests.map(function (ingest) {
        var data;
        try { data = JSON.parse(ingest.value); } catch (e) { return; }

        return new Test(data).start(function (err, latency) {
            data.err = err;
            data.latency = latency;
            callback(data);
        });
    });
};

/**
 * Halts all ongoing PingPong tests.
 */
PingPong.prototype.stop = function () {
    while (this._tests.length) {
        this._tests.pop().stop();
    }
};

/**
 * A Test is used to run pingpong tests against a single ingest instance.
 * @param {Object} ingest
 */
function Test(ingest) {
    this._ingest = ingest;
    this._lastPing = null;
    this._latencies = [];
    this._stopped = false;
}

/**
 * Starts pinging the socket, calling back to `update` with an err and latency.
 * @param  {Function} update
 * @return {Test}
 */
Test.prototype.start = function (update) {
    var cnx = this._websocket = new WebSocket(this._ingest.ping);
    var self = this;
    cnx.onopen = function () {
        self._ping();
    };

    cnx.onerror = function (err) {
        self.stop();
        update(err);
    };

    cnx.onmessage = function () {
        self._resolvePing(update);
    };

    return this;
};

Test.prototype._ping = function () {
    if (this._latencies.length >= maxTests) {
        this.stop();
    }

    if (this._stopped) {
        return;
    }

    this._lastPing = Date.now();
    this._websocket.send('ping');
};

Test.prototype._resolvePing = function (update) {
    var times = this._latencies;
    times.push((Date.now() - this._lastPing) / 2);
    var average = times.reduce(function (a, b) { return a + b; }, 0) / times.length;

    update(undefined, average);
    var self = this;

    setTimeout(function () {
        self._ping();
    }, 500)
};

/**
 * Halts ongoing tests.
 */
Test.prototype.stop = function () {
    this._stopped = true;
    this._websocket.close();
};

module.exports = PingPong;
