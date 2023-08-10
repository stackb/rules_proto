goog.provide('example.routeguide.closure.main');

/**
 * Main entry point
 * @export
 */
example.routeguide.closure.main = function () {
    const Client = goog.require('example.routeguide.closure.GrpcJsClient');
    const client = new Client();
    client.run();
};