goog.provide('example.routeguide.grpc_web.client.run');

/**
 * Main entry point
 * @export
 */
example.routeguide.grpc_web.client.run = function() {
    const Client = goog.require('example.routeguide.grpc_web.Client');
    const client = new Client("localhost", null, null);
    client.run();
};
