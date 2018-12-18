goog.module('example.routeguide.grpc_js.ClientTest');

goog.setTestOnly('example.routeguide.grpc_js.ClientTest');

const Client = goog.require('example.routeguide.grpc_js.Client');
const Rectangle = goog.require('proto.routeguide.Rectangle');

const asserts = goog.require('goog.testing.asserts');
const testSuite = goog.require('goog.testing.testSuite');

goog.require('goog.testing.jsunit');

testSuite({

    testListFeatures: () => {
        const client = new Client();
        const rect = new Rectangle();

        return client.listFeatures(rect).then(features => {
            asserts.assertEquals(0, features.length);
        });
    },

});
