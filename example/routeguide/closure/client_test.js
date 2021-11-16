goog.module('example.routeguide.closure.GrpcJsClientTest');

goog.setTestOnly('example.routeguide.closure.GrpcJsClientTest');

const Client = goog.require('proto.example.routeguide.RouteguideClient');
const Rectangle = goog.require('proto.example.routeguide.Rectangle');

const asserts = goog.require('goog.testing.asserts');
const testSuite = goog.require('goog.testing.testSuite');

goog.require('goog.testing.jsunit');

testSuite({
    testListFeatures: () => {
        const client = new Client();
        return client.getRouteGuide().listFeatures(new Rectangle(), feature => {
            // no features should be returned so this should fail the test
            asserts.assertNull(feature);
        });
    },
});