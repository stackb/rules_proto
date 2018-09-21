goog.module('example.routeguide.grpc_web.ClientTest');

goog.setTestOnly('example.routeguide.grpc_web.ClientTest');

const Client = goog.require('example.routeguide.grpc_web.Client');
const Rectangle = goog.require('proto.routeguide.Rectangle');


const asserts = goog.require('goog.testing.asserts');
const testSuite = goog.require('goog.testing.testSuite');

goog.require('goog.testing.jsunit');

testSuite({

    testListFeatures: () => {
        const client = new Client("localhost", null, null);
        const rect = new Rectangle();
        
        // Currently this segfaults phantomjs!
        const stream = client.listFeatures(rect, {});
        asserts.assertNotNull(stream);
    },

});

