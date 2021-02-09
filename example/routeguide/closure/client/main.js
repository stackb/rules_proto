goog.provide('example.routeguide.client.run');

/**
 * Main entry point
 * @export
 */
example.routeguide.client.run = function() {
    const Client = goog.require('example.routeguide.Client');
    const client = new Client();
    client.run();
};
