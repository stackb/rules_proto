goog.module('example.routeguide.ClientTest');

goog.setTestOnly('example.routeguide.ClientTest');

const Client = goog.require('example.routeguide.Client');
const Feature = goog.require('proto.routeguide.Feature');
const asserts = goog.require('goog.testing.asserts');
const testSuite = goog.require('goog.testing.testSuite');

goog.require('goog.testing.jsunit');

testSuite({

    testSerializeDeserialize: () => {
        const original = new Feature();
        original.setName("foo");
        const array = original.serializeBinary();
        const clone = Feature.deserializeBinary(array);
        asserts.assertEquals(original.getName(), clone.getName());
    },

    testAddFeature: () => {
        const feature = new Feature();
        feature.setName("foo");

        const client = new Client();
        client.addFeature(feature);

        asserts.assertEquals(1, client.getFeatures().length);
    },

});

