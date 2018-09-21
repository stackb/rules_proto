goog.provide('example.routeguide.grpc_js.client.run');

/**
 * Main entry point
 * @export
 */
example.routeguide.grpc_js.client.run = function() {
    const Client = goog.require('example.routeguide.grpc_js.Client');
    const client = new Client();
    client.run();
};
